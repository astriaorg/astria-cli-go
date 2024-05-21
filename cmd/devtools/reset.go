package devtools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd/devtools/config"

	"github.com/astria/astria-cli-go/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// resetCmd represents the root reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "The root command for resetting the local development instance data.",
	Long:  `The root command for resetting the local development instance data. The specified data will be reset to its initial state as though initialization was just run.`,
}

// resetConfigCmd represents the 'reset config' command
var resetConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Reset config files.",
	Long:  "Reset config files. This will return all files in the config directory to their default state as though initially created.",
	Run:   resetConfigCmdHandler,
}

func resetConfigCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	localNetworkName := flagHandler.GetValue("local-network-name")
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	localDefaultDenom := flagHandler.GetValue("local-default-denom")

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	localConfigDir := filepath.Join(homePath, ".astria", instance, config.DefaultConfigDirName)
	baseConfigPath := filepath.Join(localConfigDir, config.DefualtBaseConfigName)

	log.Infof("Resetting config for instance '%s'", instance)

	// Remove the config files
	err = os.Remove(filepath.Join(localConfigDir, config.DefaultCometbftGenesisFilename))
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

	config.RecreateCometbftAndSequencerGenesisData(localConfigDir, localNetworkName, localDefaultDenom)
	config.CreateBaseConfig(baseConfigPath, instance)

	log.Infof("Successfully reset config files for instance '%s'", instance)
}

// resetNetworksCmd represents the 'reset networks' command
var resetNetworksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Reset the networks config.",
	Long:  `Reset the networks config for the cli. This command only resets the networks-config.toml file. No other config files are affected.`,
	Run:   resetNetworksCmdHandler,
}

func resetNetworksCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	localNetworkName := flagHandler.GetValue("local-network-name")
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	localDefaultDenom := flagHandler.GetValue("local-default-denom")

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	networksConfigPath := filepath.Join(homePath, ".astria", instance, config.DefualtNetworksConfigName)

	err = os.Remove(networksConfigPath)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	config.CreateNetworksConfig(networksConfigPath, localNetworkName, localDefaultDenom)
}

// resetStateCmd represents the 'reset state' command
var resetStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Reset Sequencer state.",
	Long:  "Reset Sequencer state. This will reset both the sequencer and Cometbft data to their initial state.",
	Run:   resetStateCmdHandler,
}

func resetStateCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	instanceDir := filepath.Join(homePath, ".astria", instance)
	dataDir := filepath.Join(instanceDir, config.DataDirName)

	log.Infof("Resetting state for instance '%s'", instance)

	// Remove the state files for sequencer and Cometbft
	err = os.RemoveAll(dataDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	cmd.CreateDirOrPanic(dataDir)
	config.InitCometbft(instanceDir, config.DataDirName, config.BinariesDirName, config.DefaultConfigDirName)

	log.Infof("Successfully reset state for instance '%s'", instance)
}

func init() {
	// top level command
	devCmd.AddCommand(resetCmd)

	// subcommands
	resetCmd.AddCommand(resetConfigCmd)
	resetCmd.AddCommand(resetNetworksCmd)
	resetCmd.AddCommand(resetStateCmd)
}
