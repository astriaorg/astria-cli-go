package devtools

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DeleteLogs bool
var DeleteData bool

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:    "clean",
	Short:  "Delete the local development instance data and logs, excluding binaries and config.",
	Long:   `Delete the local development instance data and logs. Does not delete the binaries or the config files. If no flags are provided, the data and logs will be deleted. If flags are provided, only the selected items will be deleted.`,
	PreRun: cmd.SetLogLevel,
	Run:    runClean,
}

func init() {
	// top level command
	devCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringP("instance", "i", DefaultInstanceName, "Choose the instance that will be cleaned.")
	cleanCmd.Flags().BoolVar(&DeleteLogs, "logs", false, "Delete log files.")
	cleanCmd.Flags().BoolVar(&DeleteData, "data", false, "Delete local data.")

	// subcommands
	cleanCmd.AddCommand(allCmd)
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
	logsDir := filepath.Join(instanceDir, LogsDirName)

	// If no flags are set, delete everything
	if !DeleteLogs && !DeleteData {
		DeleteLogs = true
		DeleteData = true
	}

	// If selected, delete logs
	if DeleteLogs {
		log.Infof("Deleting logs for instance '%s'", instance)
		rmCmd := exec.Command("rm", "-rf", logsDir)
		if err := rmCmd.Run(); err != nil {
			log.WithError(err).Error("Error running rm")
			panic(err)
		}
		log.Infof("Recreating data dir for instance '%s'", instance)
		CreateDirOrPanic(logsDir)
	}

	// If selected, delete data
	if DeleteData {
		log.Infof("Deleting data for instance '%s'", instance)
		rmCmd := exec.Command("rm", "-rf", dataDir)
		if err := rmCmd.Run(); err != nil {
			log.WithError(err).Error("Error running rm")
			panic(err)
		}
		log.Debugf("Recreating data dir for instance '%s'", instance)
		CreateDirOrPanic(dataDir)
		// Reinitialize the cometbft instance after deleting the data to allow
		// user to start fresh without needing to run initialization again
		initCometbft(instanceDir, DataDirName, BinariesDirName, LocalConfigDirName)
	}
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
