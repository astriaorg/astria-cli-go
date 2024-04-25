package sequencer

import (
	"encoding/hex"
	"fmt"
	"strings"

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

	transferCmd.Flags().Bool("json", false, "Output in JSON format.")
	transferCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer.")

	transferCmd.Flags().String("keyfile", "", "Path to secure keyfile for sender.")
	transferCmd.Flags().String("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	transferCmd.Flags().String("privkey", "", "The private key of the sender.")
	transferCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
}

func transferCmdHandler(cmd *cobra.Command, args []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"

	amount := args[0]
	to := args[1]
	url := cmd.Flag("url").Value.String()

	priv, err := getPrivateKeyFromFlags(cmd)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.TransferOpts{
		SequencerURL: url,
		FromKey:      priv,
		ToAddress:    to,
		Amount:       amount,
	}
	tx, err := sequencer.Transfer(opts)
	if err != nil {
		log.WithError(err).Error("Error transferring tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

func getPrivateKeyFromFlags(cmd *cobra.Command) (string, error) {
	keyfile := cmd.Flag("keyfile").Value.String()
	keyringAddress := cmd.Flag("keyring-address").Value.String()
	priv := cmd.Flag("privkey").Value.String()

	if priv != "" {
		return priv, nil
	}

	if keyringAddress != "" {
		key, err := keys.GetKeyring(keyringAddress)
		if err != nil {
			log.WithError(err).Error("Error getting private key from keyring")
			return "", err
		}
		return key, nil
	}

	if keyfile == "" {
		kf, err := keys.ResolveKeyfilePath(keyfile)
		if err != nil {
			return "", err
		}

		// TODO - get password from user input
		password := "banana"

		privKey, err := keys.DecryptKeyfile(kf, strings.TrimRight(string(password), "\r\n"))
		if err != nil {
			log.WithError(err).Error("Error decrypting keyfile")
			return "", err
		}
		return hex.EncodeToString(privKey), nil
	}

	return "", fmt.Errorf("no private key specified")
}
