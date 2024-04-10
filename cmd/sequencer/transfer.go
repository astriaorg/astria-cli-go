package sequencer

import (
	"strconv"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:    "transfer [from] [to] [amount]",
	Short:  "Transfer tokens from one account to another.",
	Args:   cobra.ExactArgs(3),
	PreRun: cmd.SetLogLevel,
	Run:    runTransfer,
}

func init() {
	sequencerCmd.AddCommand(transferCmd)
	transferCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer.")
}

func runTransfer(cmd *cobra.Command, args []string) {
	from := args[0]
	to := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		log.WithError(err).Error("Error converting amount to integer")
		panic(err)
	}

	url := cmd.Flag("url").Value.String()

	opts := sequencer.TransferOpts{
		SequencerURL: url,
		FromKey:      from,
		ToAddress:    to,
		Amount:       amount,
	}
	err = sequencer.Transfer(opts)
	if err != nil {
		log.WithError(err).Error("Error transferring tokens")
		panic(err)
	}
}
