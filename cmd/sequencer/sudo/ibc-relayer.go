package sudo

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addIBCRelayerCmd represents the add ibc relayer command
var addIBCRelayerCmd = &cobra.Command{
	Use:   "add-ibc-relayer",
	Short: "Add an address to the IBC Relayer set on the sequencer.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add-ibc-relayer called")
	},
}

// removeIBCRelayerCmd represents the remove ibc relayer command
var removeIBCRelayerCmd = &cobra.Command{
	Use:   "remove-ibc-relayer",
	Short: "Remove an address to the IBC Relayer set on the sequencer.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remove-ibc-relayer called")
	},
}

func init() {
	SudoCmd.AddCommand(addIBCRelayerCmd)
	SudoCmd.AddCommand(removeIBCRelayerCmd)
}
