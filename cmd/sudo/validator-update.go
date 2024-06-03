package sudo

import (
	"fmt"

	"github.com/spf13/cobra"
)

// validatorUpdateCmd represents the validator update command
var validatorUpdateCmd = &cobra.Command{
	Use:   "validator-update",
	Short: "Update the validator set on the sequencer.",
	Run:   validatorUpdateCmdHandler,
}

func validatorUpdateCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("validator-update called")
}

func init() {
	sudoCmd.AddCommand(validatorUpdateCmd)
}
