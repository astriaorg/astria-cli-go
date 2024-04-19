package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var bridgeLockCmd = &cobra.Command{
	Use:    "bridge-lock [address] [amount] [destination-chain-address] --privkey=[privkey]",
	Short:  "Lock tokens on the bridge",
	Args:   cobra.ExactArgs(3),
	PreRun: cmd.SetLogLevel,
	Run:    bridgeLockCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(bridgeLockCmd)
	bridgeLockCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens")
	bridgeLockCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to lock tokens")
	bridgeLockCmd.Flags().Bool("json", false, "Output bridge lock transaction as JSON")

	err := bridgeLockCmd.MarkFlagRequired("privkey")
	if err != nil {
		log.WithError(err).Error("Error marking flag as required")
		panic(err)
	}
}

func bridgeLockCmdHandler(cmd *cobra.Command, args []string) {
	from := cmd.Flag("privkey").Value.String()
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"
	opts := sequencer.BridgeLockOpts{
		SequencerURL:     url,
		FromKey:          from,
		ToAddress:        args[0],
		Amount:           args[1],
		DestinationChain: args[2],
	}
	tx, err := sequencer.BridgeLock(opts)
	if err != nil {
		log.WithError(err).Error("Error locking tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}
