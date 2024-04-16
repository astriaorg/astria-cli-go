package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createaccountCmd represents the createaccount command
var createaccountCmd = &cobra.Command{
	Use:   "createaccount",
	Short: "Create a new account for the sequencer.",
	Long: `Create an account for the sequencer. The account will be used to sign
transactions and blocks. The account will be created with a private key, public key, and address.`,
	PreRun: cmd.SetLogLevel,
	Run:    createaccountCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(createaccountCmd)
	createaccountCmd.Flags().Bool("json", false, "Output the account information in JSON format.")
}

func createaccountCmdHandler(cmd *cobra.Command, args []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"

	account, err := sequencer.CreateAccount()
	if err != nil {
		log.WithError(err).Error("Error creating account")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      account,
		PrintJSON: printJSON,
	}
	printer.Render()
}
