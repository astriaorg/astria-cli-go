package sudo

import (
	"github.com/astria/astria-cli-go/cmd"
	sequencercmd "github.com/astria/astria-cli-go/cmd/sequencer"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sudoAddressChangeCmd represents the sudo address change command
var sudoAddressChangeCmd = &cobra.Command{
	Use:   "sudo-address-change [address] [--keyfile | --keyring-address | --privkey]",
	Short: "Update the sequencer's sudo address to a new address.",
	Long:  `Update the sequencer's sudo address to a new address. The provided address must be a valid address on the chain and will become the new sudo address.`,
	Args:  cobra.ExactArgs(1),
	Run:   sudoAddressChangeCmdHandler,
}

func sudoAddressChangeCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"

	url := flagHandler.GetValue("sequencer-url")
	sequencerURL := cmd.AddPortToURL(url)

	chainId := flagHandler.GetValue("sequencer-chain-id")

	to := args[0]
	// bech32mAddress, err := bech32m.DecodeAndValidateBech32M(to, "astria")
	// if err != nil {
	// 	log.WithError(err).Error("Error decoding address")
	// 	return
	// }

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := cmd.PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return
	}

	opts := sequencer.ChangeSudoAddressOpts{
		FromKey:          from,
		UpdateAddress:    to,
		SequencerURL:     sequencerURL,
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
	sudoCmd.AddCommand(sudoAddressChangeCmd)

	flaghandler := cmd.CreateCliFlagHandler(sudoAddressChangeCmd, cmd.EnvPrefix)
	flaghandler.BindBoolFlag("json", false, "Output the command result in JSON format.")
	flaghandler.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to update the sudo address on.")
	flaghandler.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	flaghandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flaghandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flaghandler.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
