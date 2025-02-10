package devrunner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// resetCmd represents the root reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "The root command for resetting the local development instance data.",
	Long:  "The root command for resetting the local development instance data. The specified data will be reset to its initial state as though initialization was just run.",
}

// resetConfigCmd represents the 'reset config' command
var resetConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Reset the base config files to their initial default state.",
	Long:  "Reset the base config files to their initial default state. This will return all files in the ~/.astria/<instance>/config directory to their default state as though initially created.",
	Run:   resetConfigCmdHandler,
}

func resetConfigCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	homeDir := cmd.GetUserHomeDirOrPanic()
	tuiConfigPath := filepath.Join(homeDir, ".astria", config.DefaultTUIConfigName)
	tuiConfig := config.LoadTUIConfigOrPanic(tuiConfigPath)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	if !flagHandler.GetChanged("instance") {
		instance = tuiConfig.OverrideInstanceName
	}

	localNetworkName := flagHandler.GetValue("local-network-name")
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	localDenom := flagHandler.GetValue("local-native-denom")

	localConfigDir := filepath.Join(homeDir, ".astria", instance, config.DefaultConfigDirName)
	baseConfigPath := filepath.Join(homeDir, ".astria", instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)

	log.Infof("Resetting config for instance '%s'", instance)

	// remove the config files
	err := os.Remove(filepath.Join(localConfigDir, config.DefaultCometbftGenesisFilename))
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	err = os.Remove(filepath.Join(localConfigDir, config.DefaultCometbftValidatorFilename))
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	err = os.Remove(baseConfigPath)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	config.RecreateCometbftAndSequencerGenesisData(localConfigDir, localNetworkName, localDenom)
	config.CreateBaseConfig(baseConfigPath, instance)

	log.Infof("Successfully reset config files for instance '%s'", instance)
}

// resetNetworksCmd represents the 'reset networks' command
var resetNetworksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Reset the networks config to its default values.",
	Long:  "Reset the networks config to its default values. This command only resets the ~/.astria/<instance>/networks-config.toml file.",
	Run:   resetNetworksCmdHandler,
}

func resetNetworksCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	homeDir := cmd.GetUserHomeDirOrPanic()
	tuiConfigPath := filepath.Join(homeDir, ".astria", config.DefaultTUIConfigName)
	tuiConfig := config.LoadTUIConfigOrPanic(tuiConfigPath)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	if !flagHandler.GetChanged("instance") {
		instance = tuiConfig.OverrideInstanceName
	}

	localNetworkName := flagHandler.GetValue("local-network-name")
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	localDenom := flagHandler.GetValue("local-native-denom")

	networksConfigPath := filepath.Join(homeDir, ".astria", instance, config.DefaultNetworksConfigName)

	err := os.Remove(networksConfigPath)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	genericBinariesDir := filepath.Join("~", ".astria", instance, config.BinariesDirName)
	config.CreateNetworksConfig(networksConfigPath, genericBinariesDir, localNetworkName, localDenom)
}

// resetStateCmd represents the 'reset state' command
var resetStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Reset Sequencer chain state to its initial state.",
	Long:  "Reset Sequencer chain state to its initial state. This will reset both the Astria Sequencer and Cometbft chain state to their initial value.",
	Run:   resetStateCmdHandler,
}

func resetStateCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	homeDir := cmd.GetUserHomeDirOrPanic()
	tuiConfigPath := filepath.Join(homeDir, ".astria", config.DefaultTUIConfigName)
	tuiConfig := config.LoadTUIConfigOrPanic(tuiConfigPath)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	if !flagHandler.GetChanged("instance") {
		instance = tuiConfig.OverrideInstanceName
	}

	instanceDir := filepath.Join(homeDir, ".astria", instance)
	dataDir := filepath.Join(homeDir, ".astria", instance, config.DataDirName)

	log.Infof("Resetting state for instance '%s'", instance)

	// remove the state files for sequencer and Cometbft
	err := os.RemoveAll(dataDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	cmd.CreateDirOrPanic(dataDir)
	config.InitCometbft(instanceDir, config.DataDirName, config.BinariesDirName, config.MainnetCometbftVersion, config.DefaultConfigDirName)

	log.Infof("Successfully reset state for instance '%s'", instance)
}

func init() {
	// top level command
	devCmd.AddCommand(resetCmd)

	// subcommands
	resetCmd.AddCommand(resetConfigCmd)
	rcfh := cmd.CreateCliFlagHandler(resetConfigCmd, cmd.EnvPrefix)
	rcfh.BindStringFlag("local-network-name", config.DefaultLocalNetworkName, "Set the local network name for the instance. This is used to set the chain ID in the CometBFT genesis.json file.")
	rcfh.BindStringFlag("local-native-denom", config.DefaultLocalNativeDenom, "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")

	resetCmd.AddCommand(resetNetworksCmd)
	rnfh := cmd.CreateCliFlagHandler(resetNetworksCmd, cmd.EnvPrefix)
	rnfh.BindStringFlag("local-network-name", config.DefaultLocalNetworkName, "Set the local network name for the instance. This is used to set the chain ID in the CometBFT genesis.json file.")
	rnfh.BindStringFlag("local-native-denom", config.DefaultLocalNativeDenom, "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")

	resetCmd.AddCommand(resetStateCmd)
}
