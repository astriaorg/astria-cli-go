package devtools

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:    "clean",
	Short:  "Delete the local development instance data, excluding binaries.",
	Long:   `Delete the local development instance data. Does not delete the binaries or the config files.`,
	PreRun: cmd.SetLogLevel,
	Run:    runClean,
}

func runClean(cmd *cobra.Command, args []string) {
	// Get the instance name from the -i flag or use the default
	instance := cmd.Flag("instance").Value.String()
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	defaultDir := filepath.Join(homePath, ".astria")
	instanceDir := filepath.Join(defaultDir, instance)
	dataDir := filepath.Join(instanceDir, DataDirName)

	log.Infof("Deleting data for instance '%s'", instance)
	rmCmd := exec.Command("rm", "-rf", dataDir)
	if err := rmCmd.Run(); err != nil {
		log.WithError(err).Error("Error running rm")
		panic(err)
	}

	log.Infof("Recreating data dir for instance '%s'", instance)
	CreateDirOrPanic(dataDir)
}

var allCmd = &cobra.Command{
	Use:    "all",
	Short:  "Delete everything in the ~/.astria directory.",
	Long:   "Clean all local data including binaries and config files. `dev init` will need to be run again to get the binaries and config files back.",
	PreRun: cmd.SetLogLevel,
	Run:    runCleanAll,
}

func runCleanAll(cmd *cobra.Command, args []string) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}

	// TODO: allow for configuration of this directory
	defaultDataDir := filepath.Join(homePath, ".astria")

	log.Infof("Deleting all data in %s", defaultDataDir)
	rmCmd := exec.Command("rm", "-rf", defaultDataDir)
	if err := rmCmd.Run(); err != nil {
		log.WithError(err).Error("Error running rm")
		panic(err)
	}
}

func init() {
	// top level command
	devCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringP("instance", "i", DefaultInstanceName, "Choose the instance that will be cleaned.")

	// subcommands
	cleanCmd.AddCommand(allCmd)
}
