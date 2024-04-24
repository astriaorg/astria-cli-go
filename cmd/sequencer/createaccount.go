package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/keys"
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

	// user has multiple options for storing private key
	createaccountCmd.Flags().Bool("keyring", false, "Store the account private key in the system keyring.")
	createaccountCmd.Flags().Bool("keyfile", false, "Store the account private key to an encrypted keyfile.")
	createaccountCmd.MarkFlagsMutuallyExclusive("keyring", "keyfile")
}

func createaccountCmdHandler(cmd *cobra.Command, _ []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"
	keyring := cmd.Flag("keyring").Value.String() == "true"
	keyfile := cmd.Flag("keyfile").Value.String() == "true"

	account, err := sequencer.CreateAccount()
	if err != nil {
		log.WithError(err).Error("error creating account")
		panic(err)
	}

	if keyring {
		err = keys.StoreKeyring(account.Address, account.PrivateKey)
		if err != nil {
			log.WithError(err).Error("error storing private key")
			panic(err)
		}
		// don't print private key if they choose to store in keyring or file
		account.PrivateKey = ""
	}

	if keyfile {
		// TODO
		// don't print private key if they choose to store in keyring or file
		account.PrivateKey = ""
	}

	printer := ui.ResultsPrinter{
		Data:      account,
		PrintJSON: printJSON,
	}
	printer.Render()
}
