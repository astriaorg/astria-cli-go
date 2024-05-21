package sequencer

import (
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/keys"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createaccountCmd represents the createaccount command
var createaccountCmd = &cobra.Command{
	Use:   "createaccount",
	Short: "Create a new account for the sequencer.",
	Long: `Create an account for the sequencer. The account will be used to sign
transactions and blocks. The account will be created with a private key, public key, and address.`,
	Run: createaccountCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(createaccountCmd)

	flagHandler := cmd.CreateCliFlagHandler(createaccountCmd, cmd.EnvPrefix)
	flagHandler.BindBoolFlag("json", false, "Output the account information in JSON format.")
	flagHandler.BindBoolFlag("insecure", false, "Print the account private key to terminal instead of storing securely.")
	// user has multiple options for storing private key
	flagHandler.BindBoolFlag("keyfile", false, "Store the account private key in a keyfile.")
	flagHandler.BindBoolFlag("keyring", false, "Store the account private key in the system keyring.")

	// you can't print private key AND store securely
	createaccountCmd.MarkFlagsMutuallyExclusive("insecure", "keyring", "keyfile")
}

func createaccountCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	isInsecure := flagHandler.GetValue("insecure") == "true"
	useKeyfile := flagHandler.GetValue("keyfile") == "true"
	useKeyring := flagHandler.GetValue("keyring") == "true"

	if !isInsecure && !useKeyring && !useKeyfile {
		// useKeyfile is the default if nothing is set
		useKeyfile = true
	}

	account, err := sequencer.CreateAccount()
	if err != nil {
		log.WithError(err).Error("Error creating account")
		panic(err)
	}

	if !isInsecure {
		if useKeyfile {
			pwIn := pterm.DefaultInteractiveTextInput.WithMask("*")
			pw, _ := pwIn.Show("Your new account is locked with a password. Please give a password. Do not forget this password.\nPassword:")

			ks, err := keys.NewEncryptedKeyStore(pw, account.Address, account.PrivateKey)
			if err != nil {
				log.WithError(err).Error("Error storing private key")
				panic(err)
			}
			homePath, err := os.UserHomeDir()
			if err != nil {
				log.WithError(err).Error("Error getting home dir")
				panic(err)
			}
			astriaDir := filepath.Join(homePath, ".astria")
			keydir := filepath.Join(astriaDir, "keyfiles")
			cmd.CreateDirOrPanic(keydir)

			filename, err := keys.SaveKeystoreToFile(keydir, ks)
			if err != nil {
				log.WithError(err).Error("Error storing private key")
				panic(err)
			}

			log.Infof("Storing private key in keyfile at %s", filename)
		}
		if useKeyring {
			err = keys.StoreKeyring(account.Address, account.PrivateKeyString())
			if err != nil {
				log.WithError(err).Error("Error storing private key")
				panic(err)
			}
			log.Infof("Private key for %s stored in keychain", account.Address)
		}

		// clear the private key. we don't want to print it since we are secure here
		account.PrivateKey = nil
	}

	printer := ui.ResultsPrinter{
		Data:      account,
		PrintJSON: printJSON,
	}
	printer.Render()
}
