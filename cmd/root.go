package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "astria-cli",
	Short: "A CLI to run Astria, interact with the Sequencer, deploy rollups, and more.",
	Long:  `TODO: longer description`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// disabling the completion command for now
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}
