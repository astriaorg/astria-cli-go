package sudo

import (
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/cmd/sequencer/defaults"
	util "github.com/astria/astria-cli-go/cmd/sequencer/key-utils"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// addFeeAssetCmd represents the add fee asset command
var addFeeAssetCmd = &cobra.Command{
	Use:   "add-fee-asset",
	Short: "Add a fee asset to the sequencer.",
	Run:   addFeeAssetCmdHandler,
}

func addFeeAssetCmdHandler(c *cobra.Command, args []string) {
	fmt.Println("add-fee-asset called")
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")
	asset := flagHandler.GetValue("asset")

	priv, err := util.GetPrivateKeyFromFlags(c)
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
	Use:   "remove-fee-asset",
	Short: "Remove a fee asset from the sequencer.",
	Run:   removeFeeAssetCmdHandler,
}

func removeFeeAssetCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("remove-fee-asset called")
}

func init() {
	SudoCmd.AddCommand(addFeeAssetCmd)
	flagHandler := cmd.CreateCliFlagHandler(addFeeAssetCmd, cmd.EnvPrefix)
	flagHandler.BindStringPFlag("sequencer-url", "u", defaults.DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
	flagHandler.BindBoolFlag("json", false, "Output an account's balances in JSON format.")

	SudoCmd.AddCommand(removeFeeAssetCmd)
}
