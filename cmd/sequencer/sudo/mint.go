package sudo

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mintCmd represents the mint command
var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint native assets to an account on the sequencer.",
	Run:   mintCmdHandler,
}

func mintCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("mint called")
}

func init() {
	SudoCmd.AddCommand(mintCmd)
}
