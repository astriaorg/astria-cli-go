package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:    "transfer [amount] [to] --privkey=[privkey]",
	Short:  "Transfer tokens from one account to another.",
	Args:   cobra.ExactArgs(2),
	PreRun: cmd.SetLogLevel,
	Run:    transferCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(transferCmd)

	transferCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens.")
	transferCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer.")
	// add chainId
	transferCmd.Flags().Bool("json", false, "Output in JSON format.")

	err := transferCmd.MarkFlagRequired("privkey")
	if err != nil {
		log.WithError(err).Error("Error marking flag as required")
		panic(err)
	}
}

func transferCmdHandler(cmd *cobra.Command, args []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"

	amount := args[0]
	to := args[1]

	url := cmd.Flag("url").Value.String()
	from := cmd.Flag("privkey").Value.String()

	opts := sequencer.TransferOpts{
		SequencerURL: url,
		FromKey:      from,
		ToAddress:    to,
		Amount:       amount,
	}
	tx, err := sequencer.Transfer(opts)
	if err != nil {
		log.WithError(err).Error("Error transferring tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}
