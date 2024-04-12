package sequencer

import (
	"encoding/json"
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createAccountCmd represents the create-account command
var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account for the sequencer.",
	Long: `Create an account for the sequencer. The account will be used to sign
transactions and blocks. The account will be created with a private key, public key, and address.`,
	PreRun: cmd.SetLogLevel,
	Run:    runCreateAccountCmd,
}

func init() {
	sequencerCmd.AddCommand(createAccountCmd)
	createAccountCmd.Flags().Bool("json", false, "Output the account information in JSON format.")
}

func runCreateAccountCmd(cmd *cobra.Command, args []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"

	account, err := sequencer.CreateAccount()
	if err != nil {
		log.WithError(err).Error("Error creating account")
		panic(err)
	}

	// TODO - abstract table and json printing logic to helper functions
	if printJSON {
		obj := map[string]string{
			"address":     account.Address,
			"public_key":  account.PublicKey,
			"private_key": account.PrivateKey,
		}
		j, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			log.WithError(err).Error("Error marshalling account to JSON")
			panic(err)
		}
		fmt.Println(string(j))
	} else {
		header := []string{"Address", "Public Key", "Private Key"}
		row := []string{account.Address, account.PublicKey, account.PrivateKey}
		data := pterm.TableData{header, row}
		output, err := pterm.DefaultTable.WithHasHeader().WithSeparator(" ").WithData(data).Srender()
		if err != nil {
			log.WithError(err).Error("Error rendering table")
			panic(err)
		}
		pterm.Println(output)
	}
}
