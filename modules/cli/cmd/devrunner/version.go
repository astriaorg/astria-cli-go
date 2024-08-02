package devrunner

import (
	"fmt"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
	"github.com/spf13/cobra"
)

// VersionCmd represents the sequencer command
var VersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version of the services used by the CLI.",
	Aliases: []string{"versions"},
	Run:     seqVersionCmdHandler,
}

func seqVersionCmdHandler(c *cobra.Command, _ []string) {
	// TODO: also load networks config and print the versions of the services
	// specified there
	fmt.Println("Default Service Versions:")
	fmt.Println("cometbft:        ", "v"+config.CometbftVersion)
	fmt.Println("astria-sequencer:", "v"+config.AstriaSequencerVersion)
	fmt.Println("astria-composer: ", "v"+config.AstriaComposerVersion)
	fmt.Println("astria-conductor:", "v"+config.AstriaConductorVersion)
}

func init() {
	devCmd.AddCommand(VersionCmd)
}
