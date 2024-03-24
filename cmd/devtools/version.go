package devtools

import (
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/spf13/cobra"
)

var version = "v0.0.1"

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
	cmd.RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("version", "v", "print the version")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("version", "v", false, "Print the version of the CLI tool")
}
