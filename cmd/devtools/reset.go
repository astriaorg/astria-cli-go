package devtools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// resetCmd represents the root reset command
var resetCmd = &cobra.Command{
	Use:    "reset",
	Short:  "The root command for resetting the local development instance data.",
	Long:   `The root command for resetting the local development instance data.`,
	PreRun: cmd.SetLogLevel,
}

func init() {
	// top level command
	devCmd.AddCommand(resetCmd)
	resetCmd.PersistentFlags().StringP("instance", "i", DefaultInstanceName, "Choose the target instance for resetting.")

	// subcommands
	resetCmd.AddCommand(resetConfigCmd)
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
	Short:  "Reset the Cometbft config files.",
	Long:   "Reset the Cometbft config files. This will return the config files to their default state as though initially created.",
	PreRun: cmd.SetLogLevel,
	Run:    resetConfigCmdHandler,
}

func resetConfigCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	localConfigDir := filepath.Join(homePath, ".astria", instance, LocalConfigDirName)

	log.Infof("Resetting config for instance '%s'", instance)

	// Remove the config files
	err = os.Remove(filepath.Join(localConfigDir, DefaultCometbftGenesisFilename))
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	err = os.Remove(filepath.Join(localConfigDir, DefaultCometbftValidatorFilename))
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	recreateCometbftAndSequencerGenesisData(localConfigDir)

	log.Infof("Successfully reset config files for instance '%s'", instance)
}

// resetEnvCmd represents the 'reset env' command
var resetEnvCmd = &cobra.Command{
	Use:    "env",
	Short:  "Reset the environment files.",
	Long:   `Reset the environtment files. By default this will revert all environment files to their default state as though initially created. To select a specific environment file to reset, use the --local or --remote flags.`,
	PreRun: cmd.SetLogLevel,
	Run:    resetEnvCmdHandler,
}

func resetEnvCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	instanceDir := filepath.Join(homePath, ".astria", instance)
	localConfigDir := filepath.Join(instanceDir, LocalConfigDirName)
	remoteConfigDir := filepath.Join(instanceDir, RemoteConfigDirName)

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
		recreateLocalEnvFile(instanceDir, localConfigDir)
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
		recreateRemoteEnvFile(instanceDir, remoteConfigDir)
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
		recreateLocalEnvFile(instanceDir, localConfigDir)

		_, err = os.Stat(remoteEnvPath)
		if err == nil {
			err = os.Remove(remoteEnvPath)
			if err != nil {
				fmt.Println("Error removing file:", err)
				return
			}
		}
		recreateRemoteEnvFile(instanceDir, remoteConfigDir)
		log.Infof("Successfully reset environment files for instance '%s'", instance)
	}
}

// resetStateCmd represents the 'reset state' command
var resetStateCmd = &cobra.Command{
	Use:    "state",
	Short:  "Reset the seqeuencer state.",
	Long:   "Reset the seqeuencer state. This will reset both the sequencer and Cometbft data to their initial state.",
	PreRun: cmd.SetLogLevel,
	Run:    resetStateCmdHandler,
}

func resetStateCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	instanceDir := filepath.Join(homePath, ".astria", instance)
	dataDir := filepath.Join(instanceDir, DataDirName)

	log.Infof("Resetting state for instance '%s'", instance)

	// Remove the state files for sequencer and Cometbft
	err = os.RemoveAll(dataDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	cmd.CreateDirOrPanic(dataDir)
	initCometbft(instanceDir, DataDirName, BinariesDirName, LocalConfigDirName)

	log.Infof("Successfully reset state for instance '%s'", instance)
}
