package sudo

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/cmd/sequencer/defaults"
	util "github.com/astria/astria-cli-go/cmd/sequencer/key-utils"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sudoAddressChangeCmd represents the sudo address change command
var sudoAddressChangeCmd = &cobra.Command{
	Use:   "sudo-address-change [address] [--keyfile | --keyring-address | --privkey]",
	Short: "Update the sequencer's sudo address to a new address.",
	Args:  cobra.ExactArgs(1),
	Run:   sudoAddressChangeCmdHandler,
}

func sudoAddressChangeCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	to := args[0]

	priv, err := util.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.ChangeSudoAddressOpts{
		FromKey:          priv,
		UpdateAddress:    to,
		SequencerURL:     url,
		SequencerChainID: chainId,
	}
	tx, err := sequencer.ChangeSudoAddress(opts)
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
	SudoCmd.AddCommand(sudoAddressChangeCmd)

	flaghandler := cmd.CreateCliFlagHandler(sudoAddressChangeCmd, cmd.EnvPrefix)
	flaghandler.BindBoolFlag("json", false, "Output the command result in JSON format.")
	flaghandler.BindStringPFlag("sequencer-url", "u", defaults.DefaultSequencerURL, "The URL of the sequencer to add the relayer address to.")
	flaghandler.BindStringPFlag("sequencer-chain-id", "c", defaults.DefaultSequencerChainID, "The chain ID of the sequencer.")
	flaghandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flaghandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flaghandler.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
