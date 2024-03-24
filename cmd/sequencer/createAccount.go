package cmd

import (
	"fmt"
	"os"

	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/spf13/cobra"
)

// createAccountCmd represents the createAccount command
var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account for the sequencer.",
	Long: `Create an account for the sequencer. The account will be used to sign
transactions and blocks. The account will be created with a private key, public key, and address.`,
	Run: func(cmd *cobra.Command, args []string) {
		account, err := sequencer.CreateAccount()
		if err != nil {
			fmt.Println("Error creating account:", err)
			os.Exit(1)
		}
		// FIXME - improve output. this is just a placeholder for now.
		fmt.Println("Created account:")
		fmt.Println("  Private Key:", account.PrivateKey)
		fmt.Println("  Public Key: ", account.PublicKey)
		fmt.Println("  Address:    ", account.Address)
		os.Exit(0)
	},
}

func init() {
	sequencerCmd.AddCommand(createAccountCmd)
}
