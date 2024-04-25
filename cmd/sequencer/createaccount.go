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

var useKeyfile = true

func init() {
	sequencerCmd.AddCommand(createaccountCmd)
	createaccountCmd.Flags().Bool("json", false, "Output the account information in JSON format.")

	createaccountCmd.Flags().Bool("insecure", false, "Print the account private key to terminal instead of storing securely.")
	// user has multiple options for storing private key
	createaccountCmd.Flags().BoolVar(&useKeyfile, "keyfile", true, "Store the account private key in a keyfile.")
	createaccountCmd.Flags().Bool("keyring", false, "Store the account private key in the system keyring.")

	// user can't print private key AND store securely.
	createaccountCmd.MarkFlagsMutuallyExclusive("insecure", "keyring", "keyfile")
}

func createaccountCmdHandler(cmd *cobra.Command, _ []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"

	insecure := cmd.Flag("insecure").Value.String() == "true"
	keyring := cmd.Flag("keyring").Value.String() == "true"

	account, err := sequencer.CreateAccount()
	if err != nil {
		log.WithError(err).Error("error creating account")
		panic(err)
	}

	if !insecure {
		// clear the private key since we are "secure" here
		account.PrivateKey = ""

		if keyring {
			err = keys.StoreKeyring(account.Address, account.PrivateKey)
			if err != nil {
				log.WithError(err).Error("error storing private key")
				panic(err)
			}
			log.Infof("Private key for %s stored in keychain", account.Address)
		}
		if useKeyfile {
			// TODO
			log.Infof("Storing private key in keyfile %s", "/fake/path")
		}
	}

	printer := ui.ResultsPrinter{
		Data:      account,
		PrintJSON: printJSON,
	}
	printer.Render()
}
