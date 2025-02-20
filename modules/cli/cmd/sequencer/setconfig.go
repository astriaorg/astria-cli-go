package sequencer

import (
	"os"
	"path/filepath"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// setConfigCmd represents the setconfig root command
var setConfigCmd = &cobra.Command{
	Use:   "setconfig",
	Short: "Update the configuration for the sequencer commands config.",
}

// setFeeAssetCmd represents the setconfig asset command
var setFeeAssetCmd = &cobra.Command{
	Use:   "feeasset [denom]",
	Short: "Sets the fee asset in the sequencer command configs.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setFeeAssetCmdHandler,
}

func setFeeAssetCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	sequencerNetworksConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", DefaultSequencerNetworksConfigFilename)
	// create the networks config file if it doesn't exist. Will skip if the
	// file already exists.
	CreateSequencerNetworkConfigs(sequencerNetworksConfigPath)
	sequencerNetworkConfigs := LoadSequencerNetworkConfigsOrPanic(sequencerNetworksConfigPath)

	network := flagHandler.GetValue("network")
	config.IsInstanceNameValidOrPanic(network)

	asset := args[0]
	config.IsInstanceNameValidOrPanic(asset)

	log.Info("Updating sequencer_chain_id to", asset, "in sequencer-networks-config.toml for network: ", network)
	config := sequencerNetworkConfigs.Configs[network]
	config.FeeAsset = asset
	sequencerNetworkConfigs.Configs[network] = config

	file, err := os.Create(sequencerNetworksConfigPath)
	if err != nil {
		log.Panicf("Error creating file: %s: %v", sequencerNetworksConfigPath, err)
	}
	defer file.Close()

	// Replace the networks config toml with the new config
	if err := toml.NewEncoder(file).Encode(sequencerNetworkConfigs); err != nil {
		panic(err)
	}
	log.Infof("Successfully updated networks-config.toml: %s\n", sequencerNetworksConfigPath)
}

// setAssetCmd represents the setconfig asset command
var setAssetCmd = &cobra.Command{
	Use:   "asset [denom]",
	Short: "Sets the asset in the sequencer command configs.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setAssetCmdHandler,
}

func setAssetCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	sequencerNetworksConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", DefaultSequencerNetworksConfigFilename)
	// create the networks config file if it doesn't exist. Will skip if the
	// file already exists.
	CreateSequencerNetworkConfigs(sequencerNetworksConfigPath)
	sequencerNetworkConfigs := LoadSequencerNetworkConfigsOrPanic(sequencerNetworksConfigPath)

	network := flagHandler.GetValue("network")
	config.IsInstanceNameValidOrPanic(network)

	asset := args[0]
	config.IsInstanceNameValidOrPanic(asset)

	log.Info("Updating sequencer_chain_id to", asset, "in sequencer-networks-config.toml for network: ", network)
	config := sequencerNetworkConfigs.Configs[network]
	config.Asset = asset
	sequencerNetworkConfigs.Configs[network] = config

	file, err := os.Create(sequencerNetworksConfigPath)
	if err != nil {
		log.Panicf("Error creating file: %s: %v", sequencerNetworksConfigPath, err)
	}
	defer file.Close()

	// Replace the networks config toml with the new config
	if err := toml.NewEncoder(file).Encode(sequencerNetworkConfigs); err != nil {
		panic(err)
	}
	log.Infof("Successfully updated networks-config.toml: %s\n", sequencerNetworksConfigPath)
}

// setSequencerChainIdCmd represents the setconfig sequencer-chain-id command
var setSequencerChainIdCmd = &cobra.Command{
	Use:   "sequencerchainid [chain-id]",
	Short: "Sets the sequencer chain id in the sequencer command configs.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   setSequencerChainIdCmdHandler,
}

func setSequencerChainIdCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	sequencerNetworksConfigPath := filepath.Join(cmd.GetUserHomeDirOrPanic(), ".astria", DefaultSequencerNetworksConfigFilename)
	// create the networks config file if it doesn't exist. Will skip if the
	// file already exists.
	CreateSequencerNetworkConfigs(sequencerNetworksConfigPath)
	sequencerNetworkConfigs := LoadSequencerNetworkConfigsOrPanic(sequencerNetworksConfigPath)

	network := flagHandler.GetValue("network")
	config.IsInstanceNameValidOrPanic(network)

	sequencerChainId := args[0]
	config.IsInstanceNameValidOrPanic(sequencerChainId)

	log.Info("Updating sequencer_chain_id to", sequencerChainId, "in sequencer-networks-config.toml for network: ", network)
	config := sequencerNetworkConfigs.Configs[network]
	config.SequencerChainId = sequencerChainId
	sequencerNetworkConfigs.Configs[network] = config

	file, err := os.Create(sequencerNetworksConfigPath)
	if err != nil {
		log.Panicf("Error creating file: %s: %v", sequencerNetworksConfigPath, err)
	}
	defer file.Close()

	// Replace the networks config toml with the new config
	if err := toml.NewEncoder(file).Encode(sequencerNetworkConfigs); err != nil {
		panic(err)
	}
	log.Infof("Successfully updated networks-config.toml: %s\n", sequencerNetworksConfigPath)
}

func init() {
	// root command
	SequencerCmd.AddCommand(setConfigCmd)

	// subcommands
	setConfigCmd.AddCommand(setAssetCmd)
	safh := cmd.CreateCliFlagHandler(setAssetCmd, cmd.EnvPrefix)
	safh.BindStringFlag("network", "local", "Specify the network that the sequencer chain id is being updated for.")

	setConfigCmd.AddCommand(setFeeAssetCmd)
	sfafh := cmd.CreateCliFlagHandler(setFeeAssetCmd, cmd.EnvPrefix)
	sfafh.BindStringFlag("network", "local", "Specify the network that the sequencer chain id is being updated for.")

	setConfigCmd.AddCommand(setSequencerChainIdCmd)
	sscifh := cmd.CreateCliFlagHandler(setSequencerChainIdCmd, cmd.EnvPrefix)
	sscifh.BindStringFlag("network", "local", "Specify the network that the sequencer chain id is being updated for.")

}
