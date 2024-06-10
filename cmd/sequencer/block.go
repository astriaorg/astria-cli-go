package sequencer

import (
	"strconv"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// blockCmd represents the getblock command
var blockCmd = &cobra.Command{
	Use:   "block [height]",
	Short: "Get the specific block from the sequencer.",
	Args:  cobra.ExactArgs(1),
	Run:   blockCmdHandler,
}

func init() {
	SequencerCmd.AddCommand(blockCmd)

	flagHandler := cmd.CreateCliFlagHandler(blockCmd, cmd.EnvPrefix)
	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to retrieve the block from.")
	flagHandler.BindBoolFlag("json", false, "Output the block in JSON format.")
}

func blockCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	url := flagHandler.GetValue("sequencer-url")
	printJSON := flagHandler.GetValue("json") == "true"

	h := args[0]
	height, err := strconv.ParseInt(h, 10, 64)
	if err != nil {
		log.WithError(err).Error("Error parsing block height to int64")
		panic(err)
	}

	opts := sequencer.BlockOpts{
		SequencerURL: url,
		BlockHeight:  height,
	}
	block, err := sequencer.GetBlock(opts)
	if err != nil {
		log.WithError(err)
		return
	}

	printer := ui.ResultsPrinter{
		Data:      block,
		PrintJSON: printJSON,
	}
	printer.Render()
}
