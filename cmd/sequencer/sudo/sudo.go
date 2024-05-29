package sudo

import (
	"github.com/spf13/cobra"
)

// sudoCmd represents the root sequencer sudo command
var SudoCmd = &cobra.Command{
	Use:   "sudo",
	Short: "The root command for all sudo commands for interacting with the sequencer.",
}
