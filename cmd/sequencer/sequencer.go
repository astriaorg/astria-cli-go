package cmd

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/spf13/cobra"
)

// sequencerCmd represents the sequencer command
var sequencerCmd = &cobra.Command{
	Use:   "sequencer",
	Short: "Interact with the Astria Shared Sequencer.",
	Long: `Use this command to interact with the Astria Shared Sequencer.
Generate accounts, get account balances, transfer tokens, and more.`,
	// TODO - could be neat to have this base sequencer command print out details
	//  about the sequencer and its current state, like block height, etc. should
	//  take flag like `--sequencer-url`, where default would be url for local sequencer.
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("sequencer called")
	//},
}

func init() {
	cmd.RootCmd.AddCommand(sequencerCmd)
}
