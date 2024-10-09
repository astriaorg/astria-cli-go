package sequencer

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ibctransferCmd = &cobra.Command{
	Use:   "ibctransfer [amount] [to] [src-channel] [--keyfile | --keyring-address | --privkey]",
	Short: "Ibc Transfer tokens from a sequencer account to another chain account.",
	Args:  cobra.ExactArgs(3),
	Run:   ibctransferCmdHandler,
}

func init() {
	SequencerCmd.AddCommand(ibctransferCmd)

	flagHandler := cmd.CreateCliFlagHandler(ibctransferCmd, cmd.EnvPrefix)
	flagHandler.BindBoolFlag("json", false, "Output in JSON format.")
	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer.")
	flagHandler.BindStringPFlag("sequencer-chain-id", "c", DefaultSequencerChainID, "The chain ID of the sequencer.")
	flagHandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flagHandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flagHandler.BindStringFlag("privkey", "", "The private key of the sender.")
	flagHandler.BindStringFlag("asset", DefaultAsset, "The asset to be transferred.")
	flagHandler.BindStringFlag("fee-asset", DefaultFeeAsset, "The asset used for paying fees.")
	flagHandler.BindStringFlag("network", DefaultTargetNetwork, "Configure the values to target a specific network.")
	flagHandler.BindBoolFlag("async", false, "If true, the function will return immediately. If false, the function will wait for the transaction to be seen on the network.")

	ibctransferCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	ibctransferCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}

func ibctransferCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandlerWithUseConfigFlag(c, cmd.EnvPrefix, "network")
	networkConfig := GetNetworkConfigFromFlags(flagHandler)
	flagHandler.SetConfig(networkConfig)

	sequencerURL := flagHandler.GetValue("sequencer-url")
	sequencerURL = AddPortToURL(sequencerURL)
	asset := flagHandler.GetValue("asset")
	feeAsset := flagHandler.GetValue("fee-asset")
	sequencerChainID := flagHandler.GetValue("sequencer-chain-id")
	sourceChannelID := args[2]
	destinationChainAddress := args[1]
	returnAddress := "astria12n3yqgdt92kmgmrwj6vzu7lvvsq7wn4yh94403"
	returnAddr := AddressFromText(returnAddress)
	printJSON := flagHandler.GetValue("json") == "true"

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
	amount, err := convertToUint128(args[0])
	if err != nil {
		log.WithError(err).Error("Error converting amount to Uint128 proto")
		panic(err)
	}

	isAsync := flagHandler.GetValue("async") == "true"

	opts := sequencer.IbcTransferOpts{
		IsAsync:                        isAsync,
		AddressPrefix:                  DefaultAddressPrefix,
		SequencerURL:                   sequencerURL,
		FromKey:                        from,
		DestinationChainAddressAddress: destinationChainAddress,
		ReturnAddress:                  returnAddr,
		Amount:                         amount,
		Asset:                          asset,
		FeeAsset:                       feeAsset,
		SequencerChainID:               sequencerChainID,
		SourceChannelID:                sourceChannelID,
	}
	tx, err := sequencer.IbcTransfer(opts)
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
