package devrunner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
	util "github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/utilities"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// setConfigCmd represents the setconfig root command
var setConfigCmd = &cobra.Command{
	Use:   "setconfig",
	Short: "Update the configuration for the local development instance.",
}

// setRollupNameCmd represents the setconfig rollupname command
var setRollupNameCmd = &cobra.Command{
	Use:   "rollupname [name]",
	Short: "Set the rollup name across all config for the instance.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setRollupNameCmdHandler,
}

func setRollupNameCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	network := flagHandler.GetValue("network")
	config.IsInstanceNameValidOrPanic(network)

	rollupPort := flagHandler.GetValue("rollup-port")
	util.IsValidPortOrPanic(rollupPort)

	baseConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)
	networksConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultNetworksConfigName)

	name := args[0]
	config.IsSequencerChainIdValidOrPanic(name) // rollup names follow the same rules as sequencer chain IDs
	log.Info("Setting rollup name to: ", name)

	// update base-config.toml
	baseConfig := config.LoadBaseConfigOrPanic(baseConfigPath)
	astriaComposerRollups := baseConfig["astria_composer_rollups"]

	// TODO: need more checking for port and name
	exp := `[^:]+::ws://127.0.0.1:` + rollupPort
	r := regexp.MustCompile(exp)
	updatedRollups := r.ReplaceAllString(astriaComposerRollups, name+"::ws://127.0.0.1:"+rollupPort)

	if updatedRollups == astriaComposerRollups {
		log.Warn("No changes made to base-config.toml.")
	}

	if err := config.ReplaceInFile(baseConfigPath, astriaComposerRollups, updatedRollups); err != nil {
		log.Errorf("Error updating the file: %s: %v", baseConfigPath, err)
		return
	} else {
		log.Info("Successfully updated: ", baseConfigPath)
	}
	log.Info("Updated 'astria_composer_rollups' to: ", updatedRollups)

	// update networks-config.toml
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	log.Info("Updating rollup_name to", name, "in networks-config.toml for network: ", network)
	tmpNetworkConfig := networkConfigs.Configs[network]
	tmpNetworkConfig.RollupName = name
	networkConfigs.Configs[network] = tmpNetworkConfig

	file, err := os.Create(networksConfigPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// replace the networks config toml with the new config
	if err := toml.NewEncoder(file).Encode(networkConfigs); err != nil {
		panic(err)
	}
	log.Infof("Successfully updated networks-config.toml: %s\n", networksConfigPath)
}

// setNativeAssetCmd represents the setconfig nativeasset command
var setNativeAssetCmd = &cobra.Command{
	Use:   "nativeasset [denom]",
	Short: "Set the netive asset for the sequencer across all config for the instance.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setNativeAssetCmdHandler,
}

func setNativeAssetCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	network := flagHandler.GetValue("network")
	config.IsInstanceNameValidOrPanic(network)

	networksConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultNetworksConfigName)
	cometbftGenesisPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultCometbftGenesisFilename)

	denom := args[0]
	denom = strings.ToLower(denom)
	config.IsValidDenomOrPanic(denom)

	// Update networks-config.toml
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	log.Info("Updating native_denom to", denom, "in networks-config.toml for network: ", network)
	tmpNetworkConfig := networkConfigs.Configs[network]
	tmpNetworkConfig.NativeDenom = denom
	networkConfigs.Configs[network] = tmpNetworkConfig

	file, err := os.Create(networksConfigPath)
	if err != nil {
		log.Panicf("Error creating file: %s: %v", networksConfigPath, err)
	}
	defer file.Close()

	// Replace the networks config toml with the new config
	if err := toml.NewEncoder(file).Encode(networkConfigs); err != nil {
		panic(err)
	}
	log.Infof("Successfully updated networks-config.toml: %s\n", networksConfigPath)

	// Update cometbft-genesis.json
	data, err := os.ReadFile(cometbftGenesisPath)
	if err != nil {
		log.Errorf("Error reading file: %s: %v", cometbftGenesisPath, err)
	}

	// Unmarshal into a map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		log.Panicf("Error unmarshalling JSON: %v", err)
	}
	if appState, ok := jsonMap["app_state"].(map[string]interface{}); ok {
		appState["native_asset_base_denomination"] = denom
	}

	// Marshal back to JSON with indentation
	updatedData, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		log.Panicf("Error marshalling JSON: %v", err)
	}

	// Write back to file
	err = os.WriteFile(cometbftGenesisPath, updatedData, 0644)
	if err != nil {
		log.Panicf("Error writing file: %s: %v", cometbftGenesisPath, err)
	}
}

// setFeeAssetCmd represents the setconfig feeasset command
var setFeeAssetCmd = &cobra.Command{
	Use:   "feeasset [denom]",
	Short: "Set the sequencer fee asset across all config for the instance.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setFeeAssetCmdHandler,
}

func setFeeAssetCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	baseConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)
	cometbftGenesisPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultCometbftGenesisFilename)

	denom := args[0]
	denom = strings.ToLower(denom)
	config.IsValidDenomOrPanic(denom)

	// Update base-config.toml
	baseConfig := config.LoadBaseConfigOrPanic(baseConfigPath)
	astriaComposerFeeAsset := baseConfig["astria_composer_fee_asset"]

	if err := config.ReplaceInFile(baseConfigPath, astriaComposerFeeAsset, denom); err != nil {
		log.Errorf("Error updating the file: %s: %v", baseConfigPath, err)
		return
	} else {
		log.Info("Successfully updated: ", baseConfigPath)
	}
	log.Info("Updated 'astria_composer_fee_asset' to: ", denom)

	// Update cometbft-genesis.json
	data, err := os.ReadFile(cometbftGenesisPath)
	if err != nil {
		log.Errorf("Error reading file: %s: %v", cometbftGenesisPath, err)
	}

	// Unmarshal into a map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		log.Panicf("Error unmarshalling JSON: %v", err)
	}
	if appState, ok := jsonMap["app_state"].(map[string]interface{}); ok {
		appState["allowed_fee_assets"] = []string{denom}
	}

	// Marshal back to JSON with indentation
	updatedData, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		log.Panicf("Error marshalling JSON: %v", err)
	}

	// Write back to file
	err = os.WriteFile(cometbftGenesisPath, updatedData, 0644)
	if err != nil {
		log.Panicf("Error writing file: %s: %v", cometbftGenesisPath, err)
	}
}

// setSequencerChainIdCmd represents the setconfig sequencerchainid command
var setSequencerChainIdCmd = &cobra.Command{
	Use:   "sequencerchainid [chain-id]",
	Short: "Set the sequencer chain id across all config for the instance.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setSequencerChainIdCmdHandler,
}

func setSequencerChainIdCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	network := flagHandler.GetValue("network")
	config.IsInstanceNameValidOrPanic(network)

	baseConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)
	networksConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultNetworksConfigName)
	cometbftGenesisPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultCometbftGenesisFilename)

	chainId := args[0]
	config.IsInstanceNameValidOrPanic(chainId)

	// Update base-config.toml
	baseConfig := config.LoadBaseConfigOrPanic(baseConfigPath)
	astriaComposerSequencerChainId := baseConfig["astria_composer_sequencer_chain_id"]

	// ReplaceInFile will update all instances of the chain ID in the file
	// It can be used here to update the astria_composer_sequencer_chain_id and
	// astria_conductor_expected_sequencer_chain_id fields at the same time
	if err := config.ReplaceInFile(baseConfigPath, astriaComposerSequencerChainId, chainId); err != nil {
		log.Errorf("Error updating the file: %s: %v", baseConfigPath, err)
		return
	} else {
		log.Info("Successfully updated: ", baseConfigPath)
	}
	log.Info("Updated 'astria_composer_sequencer_chain_id' to: ", chainId)
	log.Info("Updated 'astria_conductor_expected_sequencer_chain_id' to: ", chainId)

	// Update networks-config.toml
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	log.Info("Updating sequencer_chain_id to", chainId, "in networks-config.toml for network: ", network)
	tmpNetworkConfig := networkConfigs.Configs[network]
	tmpNetworkConfig.SequencerChainId = chainId
	networkConfigs.Configs[network] = tmpNetworkConfig

	file, err := os.Create(networksConfigPath)
	if err != nil {
		log.Panicf("Error creating file: %s: %v", networksConfigPath, err)
	}
	defer file.Close()

	// Replace the networks config toml with the new config
	if err := toml.NewEncoder(file).Encode(networkConfigs); err != nil {
		panic(err)
	}
	log.Infof("Successfully updated networks-config.toml: %s\n", networksConfigPath)

	// Update cometbft-genesis.json
	data, err := os.ReadFile(cometbftGenesisPath)
	if err != nil {
		log.Errorf("Error reading file: %s: %v", cometbftGenesisPath, err)
	}

	// Unmarshal into a map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		log.Panicf("Error unmarshalling JSON: %v", err)
	}

	// Update specific fields
	jsonMap["chain_id"] = chainId

	// Marshal back to JSON with indentation
	updatedData, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		log.Panicf("Error marshalling JSON: %v", err)
	}

	// Write back to file
	err = os.WriteFile(cometbftGenesisPath, updatedData, 0644)
	if err != nil {
		log.Errorf("Error writing file: %s: %v", cometbftGenesisPath, err)
	}
}

func init() {
	// root command
	devCmd.AddCommand(setConfigCmd)

	// subcommands
	setConfigCmd.AddCommand(setRollupNameCmd)
	srnfh := cmd.CreateCliFlagHandler(setRollupNameCmd, cmd.EnvPrefix)
	srnfh.BindStringFlag("rollup-port", config.DefaultRollupPort, "Select the localhost port that the rollup will be running on.")
	srnfh.BindStringFlag("network", "local", "Specify the network that the rollup name is being updated for.")

	setConfigCmd.AddCommand(setNativeAssetCmd)
	snafh := cmd.CreateCliFlagHandler(setNativeAssetCmd, cmd.EnvPrefix)
	snafh.BindStringFlag("network", "local", "Specify the network that the native asset is being updated for.")

	setConfigCmd.AddCommand(setFeeAssetCmd)

	setConfigCmd.AddCommand(setSequencerChainIdCmd)
	sscifh := cmd.CreateCliFlagHandler(setSequencerChainIdCmd, cmd.EnvPrefix)
	sscifh.BindStringFlag("network", "local", "Specify the network that the sequencer chain id is being updated for.")
}
