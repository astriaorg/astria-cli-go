package sudo

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	sequencercmd "github.com/astriaorg/astria-cli-go/modules/cli/cmd/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
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
	sequencerURL := sequencercmd.AddPortToURL(url)

	chainId := flagHandler.GetValue("sequencer-chain-id")

	to := args[0]
	toAddress := sequencercmd.AddressFromText(to)

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := sequencercmd.PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		panic(err)
	}

	async := flagHandler.GetValue("async") == "true"

	opts := sequencer.ChangeSudoAddressOpts{
		Async:            async,
		AddressPrefix:    sequencercmd.DefaultAddressPrefix,
		FromKey:          from,
		UpdateAddress:    toAddress,
		SequencerURL:     sequencerURL,
		SequencerChainID: chainId,
	}
	tx, err := sequencer.ChangeSudoAddress(opts)
	if err != nil {
		log.WithError(err).Error("Error updating sudo address")
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
	flaghandler.BindBoolFlag("async", false, "If true, the function will return immediately. If false, the function will wait for the transaction to be seen on the network.")
	flaghandler.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to update the sudo address on.")
	flaghandler.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	flaghandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flaghandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flaghandler.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
