package sudo

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addFeeAssetCmd represents the add fee asset command
var addFeeAssetCmd = &cobra.Command{
	Use:   "add-fee-asset",
	Short: "Add a fee asset to the sequencer.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add-fee-asset called")
	},
}

// removeFeeAssetCmd represents the remove fee asset command
var removeFeeAssetCmd = &cobra.Command{
	Use:   "remove-fee-asset",
	Short: "Remove a fee asset from the sequencer.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remove-fee-asset called")
	},
}

func init() {
	SudoCmd.AddCommand(addFeeAssetCmd)
	SudoCmd.AddCommand(removeFeeAssetCmd)
}
