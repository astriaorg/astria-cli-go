package sudo

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sudoAddressChangeCmd represents the sudo address change command
var sudoAddressChangeCmd = &cobra.Command{
	Use:   "sudo-address-change",
	Short: "Update the sequencer's sudo address to a new address.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sudo-address-change called")
	},
}

func init() {
	SudoCmd.AddCommand(sudoAddressChangeCmd)
}
