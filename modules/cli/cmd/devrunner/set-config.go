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

// setConfigCmd represents the set-config root command
var setConfigCmd = &cobra.Command{
	Use:   "set-config",
	Short: "Update the configuration for the local development instance.",
}

// setRollupNameCmd represents the set-config rollup-name command
var setRollupNameCmd = &cobra.Command{
	Use:   "rollup-name [name]",
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
		log.Error("Error updating the file:", baseConfigPath, ":", err)
		return
	} else {
		log.Info("Successfully updated: ", baseConfigPath)
	}
	log.Info("Updated 'astria_composer_rollups' to: ", updatedRollups)

	// update networks-config.toml
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	log.Info("Updating rollup_name to", name, "in networks-config.toml for network: ", network)
	config := networkConfigs.Configs[network]
	config.RollupName = name
	networkConfigs.Configs[network] = config

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

// setDefaultDenomCmd represents the set-config default-denom command
var setDefaultDenomCmd = &cobra.Command{
	Use:   "default-denom [denom]",
	Short: "Set the default sequencer denom across all config for the instance.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setDefaultDenomCmdHandler,
}

func setDefaultDenomCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	network := flagHandler.GetValue("network")
	config.IsInstanceNameValidOrPanic(network)

	baseConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)
	networksConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultNetworksConfigName)
	cometbftGenesisPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", instance, config.DefaultConfigDirName, config.DefaultCometbftGenesisFilename)

	denom := args[0]
	denom = strings.ToLower(denom)
	config.IsValidDenomOrPanic(denom)

	// Update base-config.toml
	baseConfig := config.LoadBaseConfigOrPanic(baseConfigPath)
	astriaComposerFeeAsset := baseConfig["astria_composer_fee_asset"]

	if err := config.ReplaceInFile(baseConfigPath, astriaComposerFeeAsset, denom); err != nil {
		log.Error("Error updating the file:", baseConfigPath, ":", err)
		return
	} else {
		log.Info("Successfully updated: ", baseConfigPath)
	}
	log.Info("Updated 'astria_composer_fee_asset' to: ", denom)

	// Update networks-config.toml
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	log.Info("Updating default_denom to", denom, "in networks-config.toml for network: ", network)
	config := networkConfigs.Configs[network]
	config.NativeDenom = denom
	networkConfigs.Configs[network] = config

	file, err := os.Create(networksConfigPath)
	if err != nil {
		log.Error("Error creating file:", networksConfigPath, ":", err)
		panic(err)
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
		log.Error("Error reading file:", cometbftGenesisPath, ":", err)
		panic(err)
	}

	// Unmarshal into a map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		log.Error("Error unmarshalling JSON:", err)
		panic(err)
	}
	if appState, ok := jsonMap["app_state"].(map[string]interface{}); ok {
		appState["allowed_fee_assets"] = []string{denom}
		appState["native_asset_base_denomination"] = denom
	}

	// Marshal back to JSON with indentation
	updatedData, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		log.Error("Error marshalling JSON:", err)
		panic(err)
	}

	// Write back to file
	err = os.WriteFile(cometbftGenesisPath, updatedData, 0644)
	if err != nil {
		log.Error("Error writing file:", cometbftGenesisPath, ":", err)
		panic(err)

	}
}

// setSequencerChainIdCmd represents the set-config sequencer-chain-id command
var setSequencerChainIdCmd = &cobra.Command{
	Use:   "sequencer-chain-id [chain-id]",
	Short: "Set the default sequencer chain id across all config for the instance.",
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
		log.Error("Error updating the file:", baseConfigPath, ":", err)
		return
	} else {
		log.Info("Successfully updated: ", baseConfigPath)
	}
	log.Info("Updated 'astria_composer_sequencer_chain_id' to: ", chainId)
	log.Info("Updated 'astria_conductor_expected_sequencer_chain_id' to: ", chainId)

	// Update networks-config.toml
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	log.Info("Updating sequencer_chain_id to", chainId, "in networks-config.toml for network: ", network)
	config := networkConfigs.Configs[network]
	config.SequencerChainId = chainId
	networkConfigs.Configs[network] = config

	file, err := os.Create(networksConfigPath)
	if err != nil {
		log.Error("Error creating file:", networksConfigPath, ":", err)
		panic(err)
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
		log.Error("Error reading file:", cometbftGenesisPath, ":", err)
		panic(err)
	}

	// Unmarshal into a map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		log.Error("Error unmarshalling JSON:", err)
		panic(err)
	}

	// Update specific fields
	jsonMap["chain_id"] = chainId

	// Marshal back to JSON with indentation
	updatedData, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		log.Error("Error marshalling JSON:", err)
		panic(err)
	}

	// Write back to file
	err = os.WriteFile(cometbftGenesisPath, updatedData, 0644)
	if err != nil {
		log.Error("Error writing file:", cometbftGenesisPath, ":", err)
		panic(err)
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

	setConfigCmd.AddCommand(setDefaultDenomCmd)
	sddfh := cmd.CreateCliFlagHandler(setDefaultDenomCmd, cmd.EnvPrefix)
	sddfh.BindStringFlag("default-denom", config.DefaultLocalNativeDenom, "Set the default sequencer denom across all config for the instance.")
	sddfh.BindStringFlag("network", "local", "Specify the network that the default denom is being updated for.")

	setConfigCmd.AddCommand(setSequencerChainIdCmd)
	sscifh := cmd.CreateCliFlagHandler(setSequencerChainIdCmd, cmd.EnvPrefix)
	sscifh.BindStringFlag("sequencer-chain-id", config.DefaultLocalNetworkName, "Set the default sequencer chain id across all config for the instance.")
	sscifh.BindStringFlag("network", "local", "Specify the network that the sequencer chain id is being updated for.")

}
