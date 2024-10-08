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

// purgeCmd represents the root purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "The root command for deleting data and files for the local development instance.",
	Long:  "The root command for deleting data and files for the local development instance. Whenever a purge command is run, it will delete the specified data.",
}

// purgeBinariesCmd represents the 'purge binaries' command
var purgeBinariesCmd = &cobra.Command{
	Use:   "binaries",
	Short: "Delete all locally downloaded service binaries for a given instance.",
	Long:  "Delete all locally downloaded service binaries for a given instance. Reinitializing is required after using this command to redownload the service binaries.",
	Run:   purgeBinariesCmdHandler,
}

func purgeBinariesCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	homeDir := cmd.GetUserHomeDirOrPanic()
	binDir := filepath.Join(homeDir, ".astria", instance, config.BinariesDirName)

	log.Infof("Deleting binaries for instance '%s'", instance)

	// remove the state files for sequencer and Cometbft
	err := os.RemoveAll(binDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	cmd.CreateDirOrPanic(binDir)

	log.Infof("Successfully deleted binaries for instance '%s'", instance)

}

// purgeAllCmd represents the 'purge all' command
var purgeAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Delete an instance. This includes chain state, configuration, and the service binaries for that instance.",
	Run:   purgeAllCmdHandler,
}

func purgeAllCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	homeDir := cmd.GetUserHomeDirOrPanic()
	instanceDir := filepath.Join(homeDir, ".astria", instance)

	log.Infof("Deleting instance '%s'", instance)

	// remove the state files for sequencer and Cometbft
	err := os.RemoveAll(instanceDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	log.Infof("Successfully deleted instance '%s'", instance)
}

// purgeLogsCmd represents the 'purge logs' command
var purgeLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Delete all logs for a given instance.",
	Long:  "Delete all logs for a given instance. Reinitializing is NOT required after using this command.",
	Run:   purgeLogsCmdHandler,
}

func purgeLogsCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	homeDir := cmd.GetUserHomeDirOrPanic()
	logDir := filepath.Join(homeDir, ".astria", instance, config.LogsDirName)

	log.Infof("Deleting logs for instance '%s'", instance)

	err := os.RemoveAll(logDir)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
	cmd.CreateDirOrPanic(logDir)

	log.Infof("Successfully deleted logs for instance '%s'", instance)
}

func init() {
	// top level command
	devCmd.AddCommand(purgeCmd)

	// subcommands
	purgeCmd.AddCommand(purgeBinariesCmd)
	purgeCmd.AddCommand(purgeAllCmd)
	purgeCmd.AddCommand(purgeLogsCmd)
}
