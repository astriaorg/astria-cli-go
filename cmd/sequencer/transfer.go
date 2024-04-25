package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/keys"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:    "transfer [amount] [to] --privkey=[privkey]",
	Short:  "Transfer tokens from one account to another.",
	Args:   cobra.ExactArgs(2),
	PreRun: cmd.SetLogLevel,
	Run:    transferCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(transferCmd)

	transferCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer.")
	transferCmd.Flags().Bool("json", false, "Output in JSON format.")

	transferCmd.Flags().String("address", "", "The address of the account from which to transfer tokens. Requires keystore or keyfile.")
	transferCmd.Flags().String("keyfile", "", "The private key of the account from which to transfer tokens.")
	transferCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens.")
	transferCmd.MarkFlagsOneRequired("address", "privkey")
}

func transferCmdHandler(cmd *cobra.Command, args []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"

	amount := args[0]
	to := args[1]

	url := cmd.Flag("url").Value.String()
	priv := cmd.Flag("privkey").Value.String()
	address := cmd.Flag("address").Value.String()

	if address != "" {
		key, err := keys.GetKeyring(address)
		if err != nil {
			log.WithError(err).Error("error getting private key from keyring")
			panic(err)
		}
		priv = key
	}

	opts := sequencer.TransferOpts{
		SequencerURL: url,
		FromKey:      priv,
		ToAddress:    to,
		Amount:       amount,
	}
	tx, err := sequencer.Transfer(opts)
	if err != nil {
		log.WithError(err).Error("error transferring tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

//func getPrivateKeyFromFlags(cmd *cobra.Command) (string, error) {
//	privkey := cmd.Flag("privkey").Value.String()
//	keyfilePath := cmd.Flag("keyfile").Value.String()
//	keyfilePass := cmd.Flag("keyfile-pass").Value.String()
//
//	if privkey != "" {
//		return privkey, nil
//	}
//
//	keyfile, err := keys.ResolveKeyfilePath(keyfilePath)
//	if err != nil {
//		return "", err
//	}
//	password, err := os.ReadFile(keyfilePass)
//	if err != nil {
//		return "", err
//	}
//
//	return keys.DecryptKeyfile(keyfile, strings.TrimRight(string(password), "\r\n"))
//}
