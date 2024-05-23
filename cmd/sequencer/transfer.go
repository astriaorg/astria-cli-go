package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer [amount] [to]",
	Short: "Transfer tokens from one account to another.",
	Args:  cobra.ExactArgs(2),
	Run:   transferCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(transferCmd)

	flagHandler := cmd.CreateCliFlagHandler(transferCmd, cmd.EnvPrefix)
	flagHandler.BindBoolFlag("json", false, "Output in JSON format.")
	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer.")
	flagHandler.BindStringPFlag("sequencer-chain-id", "c", DefaultSequencerChainID, "The chain ID of the sequencer.")
	flagHandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flagHandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flagHandler.BindStringFlag("privkey", "", "The private key of the sender.")

	transferCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	transferCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}

func transferCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	amount := args[0]
	to := args[1]

	priv, err := GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.TransferOpts{
		SequencerURL:     url,
		FromKey:          priv,
		ToAddress:        to,
		Amount:           amount,
		SequencerChainID: chainId,
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
