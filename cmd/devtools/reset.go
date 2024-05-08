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
	Use:    "reset",
	Short:  "The root command for resetting the local development instance data.",
	Long:   `The root command for resetting the local development instance data. The specified data will be reset to its initial state as though initialization was just run.`,
	PreRun: cmd.SetLogLevel,
}

func init() {
	// top level command
	devCmd.AddCommand(resetCmd)
	resetCmd.PersistentFlags().StringP("instance", "i", config.DefaultInstanceName, "Choose the target instance for purging.")

	// subcommands
	resetCmd.AddCommand(resetConfigCmd)
	resetConfigCmd.Flags().String("local-network-name", "sequencer-test-chain-0", "Set the local network name for the instance. This is used to set the chain ID in the CometBFT genesis.json file.")
	resetConfigCmd.Flags().String("local-default-denom", "nria", "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")
	resetCmd.AddCommand(resetEnvCmd)
	resetCmd.AddCommand(resetStateCmd)

	// flags for resetting specific env files
	resetEnvCmd.Flags().Bool("local", false, "Reset the local environment file.")
	resetEnvCmd.Flags().Bool("remote", false, "Reset the remote environment file.")
	resetEnvCmd.MarkFlagsMutuallyExclusive("local", "remote")
}

// resetConfigCmd represents the 'reset config' command
var resetConfigCmd = &cobra.Command{
	Use:    "config",
	Short:  "Reset config files.",
	Long:   "Reset config files. This will return the config files to their default state as though initially created.",
	PreRun: cmd.SetLogLevel,
	Run:    resetConfigCmdHandler,
}

func resetConfigCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
	config.IsInstanceNameValidOrPanic(instance)

	localNetworkName := c.Flag("local-network-name").Value.String()
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	localDefaultDenom := c.Flag("local-default-denom").Value.String()

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	localConfigDir := filepath.Join(homePath, ".astria", instance, config.LocalConfigDirName)
	networksConfigPath := filepath.Join(homePath, ".astria", instance, config.DefualtNetworksConfigName)

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
	err = os.Remove(networksConfigPath)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	config.RecreateCometbftAndSequencerGenesisData(localConfigDir, localNetworkName, localDefaultDenom)
	config.CreateNetworksConfig(networksConfigPath, localNetworkName, localDefaultDenom)

	log.Infof("Successfully reset config files for instance '%s'", instance)
}

// resetEnvCmd represents the 'reset env' command
var resetEnvCmd = &cobra.Command{
	Use:    "env",
	Short:  "Reset environment variable files.",
	Long:   `Reset environment variable files. By default this will revert all environment files to their default state as though initially created. To select a specific environment file to reset, use the --local or --remote flags.`,
	PreRun: cmd.SetLogLevel,
	Run:    resetEnvCmdHandler,
}

func resetEnvCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
	config.IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	instanceDir := filepath.Join(homePath, ".astria", instance)
	localConfigDir := filepath.Join(instanceDir, config.LocalConfigDirName)
	remoteConfigDir := filepath.Join(instanceDir, config.RemoteConfigDirName)

	// Check if we are resetting the local or remote environment files
	isLocal, _ := c.Flags().GetBool("local")
	isRemote, _ := c.Flags().GetBool("remote")

	if isLocal {
		localEnvPath := filepath.Join(localConfigDir, ".env")
		log.Infof("Resetting local environment file for instance '%s'", instance)
		_, err = os.Stat(localEnvPath)
		if err == nil {
			err = os.Remove(localEnvPath)
			if err != nil {
				fmt.Println("Error removing file:", err)
				return
			}
		}
		config.RecreateLocalEnvFile(instanceDir, localConfigDir)
		log.Infof("Successfully reset local environment file for instance '%s'", instance)

	} else if isRemote {
		remoteEnvPath := filepath.Join(remoteConfigDir, ".env")
		log.Infof("Resetting remote environment file for instance '%s'", instance)
		_, err = os.Stat(remoteEnvPath)
		if err == nil {
			err = os.Remove(remoteEnvPath)
			if err != nil {
				fmt.Println("Error removing file:", err)
				return
			}
		}
		config.RecreateRemoteEnvFile(instanceDir, remoteConfigDir)
		log.Infof("Successfully reset remote environment file for instance '%s'", instance)

	} else {
		localEnvPath := filepath.Join(localConfigDir, ".env")
		remoteEnvPath := filepath.Join(remoteConfigDir, ".env")

		log.Infof("Resetting all environment files for instance '%s'", instance)
		_, err = os.Stat(localEnvPath)
		if err == nil {
			err = os.Remove(localEnvPath)
			if err != nil {
				fmt.Println("Error removing file:", err)
				return
			}
		}
		config.RecreateLocalEnvFile(instanceDir, localConfigDir)

		_, err = os.Stat(remoteEnvPath)
		if err == nil {
			err = os.Remove(remoteEnvPath)
			if err != nil {
				fmt.Println("Error removing file:", err)
				return
			}
		}
		config.RecreateRemoteEnvFile(instanceDir, remoteConfigDir)
		log.Infof("Successfully reset environment files for instance '%s'", instance)
	}
}

// resetStateCmd represents the 'reset state' command
var resetStateCmd = &cobra.Command{
	Use:    "state",
	Short:  "Reset Sequencer state.",
	Long:   "Reset Sequencer state. This will reset both the sequencer and Cometbft data to their initial state.",
	PreRun: cmd.SetLogLevel,
	Run:    resetStateCmdHandler,
}

func resetStateCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
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
	config.InitCometbft(instanceDir, config.DataDirName, config.BinariesDirName, config.LocalConfigDirName)

	log.Infof("Successfully reset state for instance '%s'", instance)
}
