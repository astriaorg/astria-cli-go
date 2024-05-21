package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
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
	Use:   "init [rollup-id] [--keyfile | --keyring-address | --privkey]",
	Short: "Initialize a bridge account for the given rollup",
	Long: `Initialize a bridge account for the given rollup on the chain.
The sender of the transaction is used as the owner of the bridge account
and is the only actor authorized to transfer out of this account.`,
	Args:   cobra.ExactArgs(1),
	PreRun: cmd.SetLogLevel,
	Run:    bridgeInitCmdHandler,
}

func bridgeInitCmdHandler(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	rollupID := args[0]
	sequencerChainID := cmd.Flag("sequencer-chain-id").Value.String()
	assetID := cmd.Flag("asset-id").Value.String()
	feeAssetID := cmd.Flag("fee-asset-id").Value.String()

	printJSON := cmd.Flag("json").Value.String() == "true"
	priv, err := GetPrivateKeyFromFlags(cmd)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	opts := sequencer.InitBridgeOpts{
		SequencerURL:     url,
		FromKey:          priv,
		RollupID:         rollupID,
		SequencerChainID: sequencerChainID,
		AssetID:          assetID,
		FeeAssetID:       feeAssetID,
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
	Args:   cobra.ExactArgs(3),
	PreRun: cmd.SetLogLevel,
	Run:    bridgeLockCmdHandler,
}

func bridgeLockCmdHandler(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"
	priv, err := GetPrivateKeyFromFlags(cmd)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	amount := args[0]
	toAddress := args[1]
	destinationChainAddress := args[2]
	sequencerChainID := cmd.Flag("sequencer-chain-id").Value.String()
	assetID := cmd.Flag("asset-id").Value.String()
	feeAssetID := cmd.Flag("fee-asset-id").Value.String()
	opts := sequencer.BridgeLockOpts{
		SequencerURL:            url,
		FromKey:                 priv,
		ToAddress:               toAddress,
		Amount:                  amount,
		DestinationChainAddress: destinationChainAddress,
		SequencerChainID:        sequencerChainID,
		AssetID:                 assetID,
		FeeAssetID:              feeAssetID,
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
	sequencerCmd.AddCommand(bridgeCmd)

	bridgeCmd.AddCommand(bridgeInitCmd)
	bridgeInitCmd.Flags().String("sequencer-chain-id", DefaultSequencerChainID, "The chain id of the sequencer.")
	bridgeInitCmd.Flags().String("asset-id", DefaultBridgeAssetID, "The asset id of the asset we want to bridge.")
	bridgeInitCmd.Flags().String("fee-asset-id", DefaultBridgeFeeAssetID, "The fee asset id of the asset used for fees.")

	bridgeInitCmd.Flags().Bool("json", false, "Output bridge account as JSON.")
	bridgeInitCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account.")

	bridgeInitCmd.Flags().String("keyfile", "", "Path to secure keyfile for the bridge account.")
	bridgeInitCmd.Flags().String("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bridgeInitCmd.Flags().String("privkey", "", "The private key of the bridge account.")
	bridgeInitCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeInitCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	bridgeCmd.AddCommand(bridgeLockCmd)
	bridgeLockCmd.Flags().String("sequencer-chain-id", DefaultSequencerChainID, "The chain id of the sequencer.")
	bridgeLockCmd.Flags().String("asset-id", DefaultBridgeAssetID, "The asset to be locked and transferred.")
	bridgeLockCmd.Flags().String("fee-asset-id", DefaultBridgeFeeAssetID, "The asset used to pay the transaction fee.")

	bridgeLockCmd.Flags().Bool("json", false, "Output bridge account as JSON")
	bridgeLockCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account")

	bridgeLockCmd.Flags().String("keyfile", "", "Path to secure keyfile for the bridge account.")
	bridgeLockCmd.Flags().String("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bridgeLockCmd.Flags().String("privkey", "", "The private key of the bridge account.")
	bridgeLockCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeLockCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
