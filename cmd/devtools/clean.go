package devtools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans the local development environment data.",
	Long:  `Cleans the local development environment data. Does not remove the binaries or the config files.`,
	Run: func(cmd *cobra.Command, args []string) {
		runClean()
	},
}

func runClean() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	defaultDataDir := filepath.Join(homePath, ".astria/data")

	cleanCmd := exec.Command("rm", "-rf", defaultDataDir)
	if err := cleanCmd.Run(); err != nil {
		panic(err)
	}

	err = os.MkdirAll(defaultDataDir, 0755) // Read & execute by everyone, write by owner.
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Clean all local data including binaries and config files.",
	Long:  "Clean all local data including binaries and config files. `dev init` will need to be run again to get the binaries and config files back.",
	Run: func(cmd *cobra.Command, args []string) {
		runCleanAll()
	},
}

func runCleanAll() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}

	// TODO: allow for configuration of this directory
	defaultDataDir := filepath.Join(homePath, ".astria")

	cleanCmd := exec.Command("rm", "-rf", defaultDataDir)
	if err := cleanCmd.Run(); err != nil {
		panic(err)
	}
}

func init() {
	// top level command
	devCmd.AddCommand(cleanCmd)

	// subcommands
	cleanCmd.AddCommand(allCmd)
}
