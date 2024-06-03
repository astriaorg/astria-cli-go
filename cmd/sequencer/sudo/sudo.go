package sudo

import (
	"github.com/astria/astria-cli-go/cmd/sequencer"
	"github.com/spf13/cobra"
)

// sudoCmd represents the root sequencer sudo command
var sudoCmd = &cobra.Command{
	Use:   "sudo",
	Short: "The root command for all sudo commands for interacting with the sequencer.",
}

func init() {
	sequencer.SequencerCmd.AddCommand(sudoCmd)
}
