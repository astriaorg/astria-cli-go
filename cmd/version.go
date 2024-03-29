package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v0.2.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Long:  `Print the version of the CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Println(version)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
