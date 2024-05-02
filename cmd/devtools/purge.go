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
var purgeCmd = &cobra.Command{
	Use:    "purge",
	Short:  "The root command for purging the local development instance data.",
	Long:   `The root command for purging the local development instance data. Whenever a purge command is run, it will delete the specified data. You will need to re-initialize the instance to replace the data.`,
	PreRun: cmd.SetLogLevel,
}

func init() {
	// top level command
	devCmd.AddCommand(purgeCmd)
	purgeCmd.PersistentFlags().StringP("instance", "i", DefaultInstanceName, "Choose the target instance for resetting.")

	// subcommands
	purgeCmd.AddCommand(purgeBinariesCmd)
	purgeCmd.AddCommand(purgeAllCmd)
}

// purgeBinariesCmd represents the 'purge binaries' command
var purgeBinariesCmd = &cobra.Command{
	Use:    "binaries",
	Short:  "Delete all binaries for a given instance.",
	Long:   "Delete all binaries for a given instance.",
	PreRun: cmd.SetLogLevel,
	Run:    purgeBinariesCmdHandler,
}

func purgeBinariesCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	binDir := filepath.Join(homePath, ".astria", instance, BinariesDirName)

	log.Infof("Deleting binaries for instance '%s'", instance)

	// Remove the state files for sequencer and Cometbft
	err = os.RemoveAll(binDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	cmd.CreateDirOrPanic(binDir)

	log.Infof("Successfully deleted binaries for instance '%s'", instance)

}

// purgeAllCmd represents the 'purge all' command
var purgeAllCmd = &cobra.Command{
	Use:    "all",
	Short:  "Delete the entire instance.",
	Long:   "Delete the entire instance directory. This will remove all data, binaries, and configuration files for the specified instance.",
	PreRun: cmd.SetLogLevel,
	Run:    purgeAllCmdHandler,
}

func purgeAllCmdHandler(c *cobra.Command, _ []string) {
	// Get the instance name from the -i flag or use the default
	instance, _ := c.Parent().PersistentFlags().GetString("instance")
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	instanceDir := filepath.Join(homePath, ".astria", instance)

	log.Infof("Deleting instance '%s'", instance)

	// Remove the state files for sequencer and Cometbft
	err = os.RemoveAll(instanceDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	log.Infof("Successfully deleted instance '%s'", instance)
}
