package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var bridgeLockCmd = &cobra.Command{
	Use:    "bridgelock [address] [amount] [destination-chain-address] --privkey=[privkey]",
	Short:  "Locks tokens on the bridge account",
	Args:   cobra.ExactArgs(3),
	PreRun: cmd.SetLogLevel,
	Run:    bridgeLockCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(bridgeLockCmd)

	bridgeLockCmd.Flags().Bool("json", false, "Output bridge account as JSON")
	bridgeLockCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account")

	bridgeLockCmd.Flags().String("keyfile", "", "Path to secure keyfile for sender.")
	bridgeLockCmd.Flags().String("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	bridgeLockCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens")
	bridgeLockCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeLockCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}

func bridgeLockCmdHandler(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"
	priv, err := GetPrivateKeyFromFlags(cmd)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	opts := sequencer.BridgeLockOpts{
		SequencerURL:     url,
		FromKey:          priv,
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
