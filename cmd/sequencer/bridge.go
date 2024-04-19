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
	Use:    "init-bridge [rollup-id] --privkey=[privkey]",
	Short:  "Initializing a bridge account",
	Args:   cobra.ExactArgs(1),
	PreRun: cmd.SetLogLevel,
	Run:    initBridgeCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(initBridgeCmd)
	initBridgeCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens")
	initBridgeCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account")
	initBridgeCmd.Flags().Bool("json", false, "Output bridge account as JSON")

	err := initBridgeCmd.MarkFlagRequired("privkey")
	if err != nil {
		log.WithError(err).Error("Error marking flag as required")
		panic(err)
	}

}

func initBridgeCmdHandler(cmd *cobra.Command, args []string) {
	from := cmd.Flag("privkey").Value.String()
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"
	opts := sequencer.InitBridgeOpts{
		SequencerURL: url,
		FromKey:      from,
		RollupID:     args[0],
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
