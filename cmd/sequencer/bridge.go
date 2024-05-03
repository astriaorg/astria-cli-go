package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initBridgeCmd represents the init-bridge command
var initBridgeCmd = &cobra.Command{
	Use:    "initbridge [rollup-id] --privkey=[privkey]",
	Short:  "Initializing a bridge account",
	Args:   cobra.ExactArgs(1),
	PreRun: cmd.SetLogLevel,
	Run:    initBridgeCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(initBridgeCmd)
	initBridgeCmd.Flags().Bool("json", false, "Output bridge account as JSON")
	initBridgeCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account")

	initBridgeCmd.Flags().String("keyfile", "", "Path to secure keyfile for sender.")
	initBridgeCmd.Flags().String("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	initBridgeCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens")
	initBridgeCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	initBridgeCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

}

func initBridgeCmdHandler(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	rollupId := args[0]
	printJSON := cmd.Flag("json").Value.String() == "true"
	priv, err := GetPrivateKeyFromFlags(cmd)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	opts := sequencer.InitBridgeOpts{
		SequencerURL: url,
		FromKey:      priv,
		RollupID:     rollupId,
	}
	bridgeAccount, err := sequencer.InitBridgeAccount(opts)
	if err != nil {
		log.WithError(err).Error("Error creating account")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      bridgeAccount,
		PrintJSON: printJSON,
	}
	printer.Render()
}
