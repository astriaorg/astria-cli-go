package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags in the build-for-release workflow
var Version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	if Version == "" {
		Version = "development"
	}
	fmt.Println(Version)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
