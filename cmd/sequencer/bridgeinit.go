package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// bridgeInitCmd represents the init-bridge command
var bridgeInitCmd = &cobra.Command{
	Use:    "bridge init [rollup-id]",
	Short:  "Initializing a bridge account",
	Args:   cobra.ExactArgs(1),
	PreRun: cmd.SetLogLevel,
	Run:    bridgeInitCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(bridgeInitCmd)
	bridgeInitCmd.Flags().Bool("json", false, "Output bridge account as JSON")
	bridgeInitCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account")
	
	bridgeInitCmd.Flags().String("keyfile", "", "Path to secure keyfile for the bridge account.")
	bridgeInitCmd.Flags().String("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bridgeInitCmd.Flags().String("privkey", "", "The private key of the bridge account.")

	bridgeInitCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeInitCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

}

func bridgeInitCmdHandler(cmd *cobra.Command, args []string) {
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
