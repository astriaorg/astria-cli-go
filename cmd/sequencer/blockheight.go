package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// blockheightCmd represents the get-blockheight command
var blockheightCmd = &cobra.Command{
	Use:    "blockheight",
	Short:  "Retrieves the latest blockheight from the sequencer.",
	PreRun: cmd.SetLogLevel,
	Run:    blockheightCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(blockheightCmd)
	blockheightCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
	blockheightCmd.Flags().Bool("json", false, "Output the sequencer blockheight in JSON format.")
}

func blockheightCmdHandler(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"

	blockheight, err := sequencer.GetBlockheight(url)
	if err != nil {
		log.WithError(err)
		return
	}

	printer := ui.ResultsPrinter{
		Data:      blockheight,
		PrintJSON: printJSON,
	}
	printer.Render()
}
