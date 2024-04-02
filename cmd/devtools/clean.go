package devtools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans the local development instance data.",
	Long:  `Cleans the local development instance data. Does not remove the binaries or the config files.`,
	Run:   runClean,
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

	cleanCmd := exec.Command("rm", "-rf", dataDir)
	if err := cleanCmd.Run(); err != nil {
		log.WithError(err).Error("Error running rm")
		panic(err)
	}

	err = os.MkdirAll(dataDir, 0755) // Read & execute by everyone, write by owner.
	if err != nil {
		log.WithError(err).Error("Error creating data directory")
		panic(err)
	}
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Delete everything in the ~/.astria directory.",
	Long:  "Clean all local data including binaries and config files. `dev init` will need to be run again to get the binaries and config files back.",
	Run: func(cmd *cobra.Command, args []string) {
		runCleanAll()
	},
}

func runCleanAll() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}

	// TODO: allow for configuration of this directory
	defaultDataDir := filepath.Join(homePath, ".astria")

	cleanCmd := exec.Command("rm", "-rf", defaultDataDir)
	if err := cleanCmd.Run(); err != nil {
		log.WithError(err).Error("Error running rm")
		panic(err)
	}
}

func init() {
	// top level command
	devCmd.AddCommand(cleanCmd)
	instanceFlagUsage := fmt.Sprintf("Choose the instance that will be cleaned.", DefaultInstanceName)
	cleanCmd.Flags().StringP("instance", "i", DefaultInstanceName, instanceFlagUsage)

	// subcommands
	cleanCmd.AddCommand(allCmd)
}
