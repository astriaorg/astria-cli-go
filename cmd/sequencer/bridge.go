package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/bech32m"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Bridge commands.",
}

// bridgeInitCmd represents the `bridge init` command
var bridgeInitCmd = &cobra.Command{
	Use:   "init [rollup-name] [--keyfile | --keyring-address | --privkey]",
	Short: "Initialize a bridge account for the given rollup",
	Long: `Initialize a bridge account for the given rollup on the chain.
The sender of the transaction is used as the owner of the bridge account
and is the only actor authorized to transfer out of this account.`,
	Args: cobra.ExactArgs(1),
	Run:  bridgeInitCmdHandler,
}

func bridgeInitCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	url := flagHandler.GetValue("sequencer-url")
	sequencerURL := cmd.AddPortToURL(url)

	printJSON := flagHandler.GetValue("json") == "true"
	sequencerChainID := flagHandler.GetValue("sequencer-chain-id")
	assetID := flagHandler.GetValue("asset-id")
	feeAssetID := flagHandler.GetValue("fee-asset-id")

	sudoAddress := flagHandler.GetValue("sudo-address")
	bech32MSudoAddress, err := bech32m.DecodeAndValidateBech32M(sudoAddress, "astria")
	if err != nil {
		log.WithError(err).Error("Error decoding sudo address")
		return
	}
	withdrawerAddress := flagHandler.GetValue("withdrawer-address")
	bech32MWithdrawerAddress, err := bech32m.DecodeAndValidateBech32M(withdrawerAddress, "astria")
	if err != nil {
		log.WithError(err).Error("Error decoding withdrawer address")
		return
	}

	rollupName := args[0]

	priv, err := GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := cmd.PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return
	}

	opts := sequencer.InitBridgeOpts{
		SequencerURL:      sequencerURL,
		FromKey:           from,
		RollupName:        rollupName,
		SequencerChainID:  sequencerChainID,
		AssetID:           cmd.AssetIdFromDenom(assetID),
		FeeAssetID:        cmd.AssetIdFromDenom(feeAssetID),
		SudoAddress:       bech32MSudoAddress.AsProtoAddress(),
		WithdrawerAddress: bech32MWithdrawerAddress.AsProtoAddress(),
	}
	bridgeAccount, err := sequencer.InitBridgeAccount(opts)
	if err != nil {
		log.WithError(err).Error("Error creating account")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      bridgeAccount,
		PrintJSON: printJSON,
	}
	printer.Render()
}

// bridgeLockCmd represents the `bridge lock` command
var bridgeLockCmd = &cobra.Command{
	Use:   "lock [amount] [to] [destination-chain-address] [--keyfile | --keyring-address | --privkey]",
	Short: "Lock tokens on the bridge account",
	Long: `A bridge lock is a transfer of tokens from the signing Sequencer
account to a Sequencer bridge account. These tokens will then be
bridged to a destination chain address if an IBC relayer is running.`,
	Args: cobra.ExactArgs(3),
	Run:  bridgeLockCmdHandler,
}

func bridgeLockCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	url := flagHandler.GetValue("sequencer-url")
	sequencerURL := cmd.AddPortToURL(url)

	printJSON := flagHandler.GetValue("json") == "true"
	sequencerChainID := flagHandler.GetValue("sequencer-chain-id")
	assetID := flagHandler.GetValue("asset-id")
	feeAssetID := flagHandler.GetValue("fee-asset-id")

	priv, err := GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := cmd.PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return
	}

	amount, err := cmd.ConvertToUint128(args[0])
	if err != nil {
		log.WithError(err).Error("Error converting amount to Uint128 proto")
		return
	}

	toAddress := args[1]
	bech32mToAddress, err := bech32m.DecodeAndValidateBech32M(toAddress, "astria")
	if err != nil {
		log.WithError(err).Error("Error decoding address")
		return
	}

	destinationChainAddress := args[2]

	opts := sequencer.BridgeLockOpts{
		SequencerURL:            sequencerURL,
		FromKey:                 from,
		ToAddress:               bech32mToAddress.AsProtoAddress(),
		Amount:                  amount,
		DestinationChainAddress: destinationChainAddress,
		SequencerChainID:        sequencerChainID,
		AssetID:                 cmd.AssetIdFromDenom(assetID),
		FeeAssetID:              cmd.AssetIdFromDenom(feeAssetID),
	}
	tx, err := sequencer.BridgeLock(opts)
	if err != nil {
		log.WithError(err).Error("Error locking tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

func init() {
	SequencerCmd.AddCommand(bridgeCmd)

	bridgeCmd.AddCommand(bridgeInitCmd)
	bifh := cmd.CreateCliFlagHandler(bridgeInitCmd, cmd.EnvPrefix)
	bifh.BindStringPFlag("sequencer-chain-id", "c", DefaultSequencerChainID, "The chain ID of the sequencer.")
	bifh.BindStringFlag("asset-id", DefaultBridgeAssetID, "The asset id of the asset we want to bridge.")
	bifh.BindStringFlag("fee-asset-id", DefaultBridgeFeeAssetID, "The fee asset id of the asset used for fees.")
	bifh.BindStringFlag("sudo-address", "", "Set the sudo address to use for the bridge account. The address of the sender is used if this is not set.")
	bifh.BindStringFlag("withdrawer-address", "", "Set the withdrawer address to use for the bridge account. The address of the sender is used if this is not set.")

	bifh.BindBoolFlag("json", false, "Output bridge account as JSON.")
	bifh.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to init bridge account on.")

	bifh.BindStringFlag("keyfile", "", "Path to secure keyfile for the bridge account.")
	bifh.BindStringFlag("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bifh.BindStringFlag("privkey", "", "The private key of the bridge account.")
	bridgeInitCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeInitCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	bridgeCmd.AddCommand(bridgeLockCmd)
	blfh := cmd.CreateCliFlagHandler(bridgeLockCmd, cmd.EnvPrefix)
	blfh.BindStringFlag("sequencer-chain-id", DefaultSequencerChainID, "The chain ID of the sequencer.")
	blfh.BindStringFlag("asset-id", DefaultBridgeAssetID, "The asset to be locked and transferred.")
	blfh.BindStringFlag("fee-asset-id", DefaultBridgeFeeAssetID, "The asset used to pay the transaction fee.")

	blfh.BindBoolFlag("json", false, "Output bridge account as JSON")
	blfh.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to lock assets on.")

	blfh.BindStringFlag("keyfile", "", "Path to secure keyfile for the bridge account.")
	blfh.BindStringFlag("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	blfh.BindStringFlag("privkey", "", "The private key of the bridge account.")
	bridgeLockCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeLockCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
