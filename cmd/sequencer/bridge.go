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
	bifh := cmd.CreateCliFlagHandler(bridgeInitCmd, cmd.EnvPrefix)
	bifh.BindStringPFlag("sequencer-chain-id", "c", DefaultSequencerChainID, "The chain ID of the sequencer.")
	bifh.BindStringFlag("asset-id", DefaultBridgeAssetID, "The asset id of the asset we want to bridge")
	bifh.BindStringFlag("fee-asset-id", DefaultBridgeFeeAssetID, "The fee asset id of the asset used for fees")

	bifh.BindBoolFlag("json", false, "Output bridge account as JSON")
	bifh.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to init bridge account")

	bifh.BindStringFlag("keyfile", "", "Path to secure keyfile for the bridge account.")
	bifh.BindStringFlag("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bifh.BindStringFlag("privkey", "", "The private key of the bridge account.")
	bridgeInitCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeInitCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	bridgeCmd.AddCommand(bridgeLockCmd)
	blfh := cmd.CreateCliFlagHandler(bridgeLockCmd, cmd.EnvPrefix)
	blfh.BindBoolFlag("json", false, "Output bridge account as JSON")
	blfh.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to init bridge account")

	blfh.BindStringFlag("keyfile", "", "Path to secure keyfile for the bridge account.")
	blfh.BindStringFlag("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	blfh.BindStringFlag("privkey", "", "The private key of the bridge account.")
	bridgeLockCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeLockCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}

func bridgeInitCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	url := flagHandler.GetValue("sequencer-url")
	chainID := flagHandler.GetValue("chain-id")
	assetID := flagHandler.GetValue("asset-id")
	feeAssetID := flagHandler.GetValue("fee-asset-id")
	printJSON := flagHandler.GetValue("json") == "true"

	rollupID := args[0]

	priv, err := GetPrivateKeyFromFlags(c)
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

func bridgeLockCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	url := flagHandler.GetValue("sequencer-url")
	printJSON := flagHandler.GetValue("json") == "true"

	priv, err := GetPrivateKeyFromFlags(c)
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
