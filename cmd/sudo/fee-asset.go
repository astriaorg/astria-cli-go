package sudo

import (
	"github.com/astria/astria-cli-go/cmd"
	sequencercmd "github.com/astria/astria-cli-go/cmd/sequencer"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// feeAssetCmd represents the root command for interacting with fee assets on
// the sequencer.
var feeAssetCmd = &cobra.Command{
	Use:   "fee-asset",
	Short: "Interact with fee assets on the sequencer.",
}

// addFeeAssetCmd represents the add fee asset command
var addFeeAssetCmd = &cobra.Command{
	Use:   "add [asset] [--keyfile | --keyring-address | --privkey]",
	Short: "Add a fee asset to the sequencer.",
	Args:  cobra.ExactArgs(1),
	Run:   addFeeAssetCmdHandler,
}

func addFeeAssetCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	asset := args[0]

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.FeeAssetOpts{
		FromKey:          priv,
		SequencerURL:     url,
		SequencerChainID: chainId,
		Asset:            asset,
	}
	tx, err := sequencer.AddFeeAsset(opts)
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

// removeFeeAssetCmd represents the remove fee asset command
var removeFeeAssetCmd = &cobra.Command{
	Use:   "remove [asset] [--keyfile | --keyring-address | --privkey]",
	Short: "Remove a fee asset from the sequencer.",
	Args:  cobra.ExactArgs(1),
	Run:   removeFeeAssetCmdHandler,
}

func removeFeeAssetCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	asset := args[0]

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.FeeAssetOpts{
		FromKey:          priv,
		SequencerURL:     url,
		SequencerChainID: chainId,
		Asset:            asset,
	}
	tx, err := sequencer.RemoveFeeAsset(opts)
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

func init() {
	sudoCmd.AddCommand(feeAssetCmd)
	feeAssetCmd.AddCommand(addFeeAssetCmd)

	afafh := cmd.CreateCliFlagHandler(addFeeAssetCmd, cmd.EnvPrefix)
	afafh.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to add fee asset to.")
	afafh.BindBoolFlag("json", false, "Output the command result in JSON format.")
	afafh.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	afafh.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	afafh.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	afafh.BindStringFlag("privkey", "", "The private key of the sender.")
	addFeeAssetCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addFeeAssetCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	feeAssetCmd.AddCommand(removeFeeAssetCmd)

	rfafh := cmd.CreateCliFlagHandler(removeFeeAssetCmd, cmd.EnvPrefix)
	rfafh.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to remove fee asset from.")
	rfafh.BindBoolFlag("json", false, "Output the command result in JSON format.")
	rfafh.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	rfafh.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	rfafh.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	rfafh.BindStringFlag("privkey", "", "The private key of the sender.")
	removeFeeAssetCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	removeFeeAssetCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
