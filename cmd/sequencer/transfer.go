package sequencer

import (
	"strconv"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:    "transfer [amount] [to] --privkey=[privkey]",
	Short:  "Transfer tokens from one account to another.",
	Args:   cobra.ExactArgs(2),
	PreRun: cmd.SetLogLevel,
	Run:    runTransfer,
}

func init() {
	sequencerCmd.AddCommand(transferCmd)

	transferCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens.")
	transferCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer.")

	err := transferCmd.MarkFlagRequired("privkey")
	if err != nil {
		log.WithError(err).Error("Error marking flag as required")
		panic(err)
	}
}

func runTransfer(cmd *cobra.Command, args []string) {
	amount, err := strconv.Atoi(args[0])
	if err != nil {
		log.WithError(err).Error("Error converting amount to integer")
		panic(err)
	}
	to := args[1]

	url := cmd.Flag("url").Value.String()
	from := cmd.Flag("privkey").Value.String()

	opts := sequencer.TransferOpts{
		SequencerURL: url,
		FromKey:      from,
		ToAddress:    to,
		Amount:       amount,
	}
	hash, err := sequencer.Transfer(opts)
	if err != nil {
		log.WithError(err).Error("Error transferring tokens")
		panic(err)
	}

	pterm.Printfln("Transferred %d tokens to %s", amount, to)
	pterm.Printfln("Transaction hash: %s", hash)
}
