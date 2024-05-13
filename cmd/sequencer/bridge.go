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
	Use:    "init [rollup-id]",
	Short:  "Initialize a bridge account",
	Args:   cobra.ExactArgs(1),
	PreRun: cmd.SetLogLevel,
	Run:    bridgeInitCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(bridgeCmd)

	bridgeCmd.AddCommand(bridgeInitCmd)
	bridgeInitCmd.Flags().String("chain-id", DefaultSequencerChainID, "The chain id of the sequencer")
	bridgeInitCmd.Flags().String("asset-id", DefaultBridgeAssetID, "The asset id of the asset we want to bridge")
	bridgeInitCmd.Flags().String("fee-asset-id", DefaultBridgeFeeAssetID, "The fee asset id of the asset used for fees")

	bridgeInitCmd.Flags().Bool("json", false, "Output bridge account as JSON")
	bridgeInitCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account")

	bridgeInitCmd.Flags().String("keyfile", "", "Path to secure keyfile for the bridge account.")
	bridgeInitCmd.Flags().String("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bridgeInitCmd.Flags().String("privkey", "", "The private key of the bridge account.")
	bridgeInitCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeInitCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	bridgeCmd.AddCommand(bridgeLockCmd)
	bridgeLockCmd.Flags().Bool("json", false, "Output bridge account as JSON")
	bridgeLockCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to init bridge account")

	bridgeLockCmd.Flags().String("keyfile", "", "Path to secure keyfile for the bridge account.")
	bridgeLockCmd.Flags().String("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bridgeLockCmd.Flags().String("privkey", "", "The private key of the bridge account.")
	bridgeLockCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeLockCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}

func bridgeInitCmdHandler(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	rollupID := args[0]
	chainID := cmd.Flag("chain-id").Value.String()
	assetID := cmd.Flag("asset-id").Value.String()
	feeAssetID := cmd.Flag("fee-asset-id").Value.String()

	printJSON := cmd.Flag("json").Value.String() == "true"
	priv, err := GetPrivateKeyFromFlags(cmd)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	opts := sequencer.InitBridgeOpts{
		SequencerURL: url,
		FromKey:      priv,
		RollupID:     rollupID,
		ChainID:      chainID,
		AssetId:      assetID,
		FeeAssetID:   feeAssetID,
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
	Use:    "lock [address] [amount] [destination-chain-address]",
	Short:  "Locks tokens on the bridge account",
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
	opts := sequencer.BridgeLockOpts{
		SequencerURL:     url,
		FromKey:          priv,
		ToAddress:        args[0],
		Amount:           args[1],
		DestinationChain: args[2],
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
