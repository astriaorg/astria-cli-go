package sequencer

import (
	"strconv"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// blockCmd represents the block command
var blockCmd = &cobra.Command{
	Use:   "block [height]",
	Short: "Get sequencer block at specified height.",
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
	printJSON := flagHandler.GetValue("json") == "true"

	url := flagHandler.GetValue("sequencer-url")
	sequencerURL := AddPortToURL(url)

	height, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.WithError(err).Error("Error parsing block height to int64")
		panic(err)
	}

	opts := sequencer.BlockOpts{
		SequencerURL: sequencerURL,
		BlockHeight:  height,
	}
	block, err := sequencer.GetBlock(opts)
	if err != nil {
		log.WithError(err).Error("Error getting block")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      block,
		PrintJSON: printJSON,
	}
	printer.Render()
}
