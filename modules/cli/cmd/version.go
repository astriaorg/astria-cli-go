package cmd

import (
	"fmt"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
	"github.com/spf13/cobra"
)

// version is set at build time via ldflags in the build-for-release workflow
var version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	if version == "" {
		version = "development"
	}
	fmt.Println("astria-go:", version)
	fmt.Println("  Default Service Versions:")
	fmt.Println("  cometbft:        ", "v"+config.Cometbft_version)
	fmt.Println("  astria-sequencer:", "v"+config.Astria_sequencer_version)
	fmt.Println("  astria-composer: ", "v"+config.Astria_composer_version)
	fmt.Println("  astria-conductor:", "v"+config.Astria_conductor_version)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
