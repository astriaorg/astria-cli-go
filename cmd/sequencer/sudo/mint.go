package sudo

import (
	"github.com/astria/astria-cli-go/cmd"
	sequencercmd "github.com/astria/astria-cli-go/cmd/sequencer"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// mintCmd represents the mint command
var mintCmd = &cobra.Command{
	Use:   "mint [amount] [to] [--keyfile | --keyring-address | --privkey]",
	Short: "Mint tokens to an account.",
	Args:  cobra.ExactArgs(2),
	Run:   mintCmdHandler,
}

func mintCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	amount := args[0]
	to := args[1]

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.MintOpts{
		FromKey:          priv,
		ToAddress:        to,
		SequencerURL:     url,
		SequencerChainID: chainId,
		Amount:           amount,
	}
	tx, err := sequencer.Mint(opts)
	if err != nil {
		log.WithError(err).Error("Error minting tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

func init() {
	sudoCmd.AddCommand(mintCmd)

	flaghandler := cmd.CreateCliFlagHandler(mintCmd, cmd.EnvPrefix)
	flaghandler.BindBoolFlag("json", false, "Output the command result in JSON format.")
	flaghandler.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to add the relayer address to.")
	flaghandler.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	flaghandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flaghandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flaghandler.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
