package devtools

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Commands for running the Astria Shared Sequencer.",
}

func init() {
	cmd.RootCmd.AddCommand(devCmd)
}
