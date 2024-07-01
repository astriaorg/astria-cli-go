package sequencer

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer [amount] [to] [--keyfile | --keyring-address | --privkey]",
	Short: "Transfer tokens from one account to another.",
	Args:  cobra.ExactArgs(2),
	Run:   transferCmdHandler,
}

func init() {
	SequencerCmd.AddCommand(transferCmd)

	flagHandler := cmd.CreateCliFlagHandler(transferCmd, cmd.EnvPrefix)
	flagHandler.BindBoolFlag("json", false, "Output in JSON format.")
	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer.")
	flagHandler.BindStringPFlag("sequencer-chain-id", "c", DefaultSequencerChainID, "The chain ID of the sequencer.")
	flagHandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flagHandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flagHandler.BindStringFlag("privkey", "", "The private key of the sender.")
	flagHandler.BindStringFlag("asset", DefaultAsset, "The asset to be transferred.")
	flagHandler.BindStringFlag("fee-asset", DefaultFeeAsset, "The asset used for paying fees.")

	transferCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	transferCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}

func transferCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"

	url := flagHandler.GetValue("sequencer-url")
	sequencerURL := AddPortToURL(url)

	priv, err := GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		panic(err)
	}

	to := args[1]
	toAddress := AddressFromText(to)

	amount, err := convertToUint128(args[0])
	if err != nil {
		log.WithError(err).Error("Error converting amount to Uint128 proto")
		panic(err)
	}

	asset := flagHandler.GetValue("asset")
	feeAsset := flagHandler.GetValue("fee-asset")

	chainId := flagHandler.GetValue("sequencer-chain-id")

	opts := sequencer.TransferOpts{
		AddressPrefix:    DefaultAddressPrefix,
		SequencerURL:     sequencerURL,
		FromKey:          from,
		ToAddress:        toAddress,
		Amount:           amount,
		Asset:            asset,
		FeeAsset:         feeAsset,
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
