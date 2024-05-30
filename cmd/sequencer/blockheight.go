package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/cmd/sequencer/defaults"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// blockheightCmd represents the blockheight command
var blockheightCmd = &cobra.Command{
	Use:   "blockheight",
	Short: "Retrieves the latest blockheight from the sequencer.",
	Run:   blockheightCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(blockheightCmd)

	flagHandler := cmd.CreateCliFlagHandler(blockheightCmd, cmd.EnvPrefix)
	flagHandler.BindStringPFlag("sequencer-url", "u", defaults.DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
	flagHandler.BindBoolFlag("json", false, "Output an account's balances in JSON format.")
}

func blockheightCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	url := flagHandler.GetValue("sequencer-url")
	printJSON := flagHandler.GetValue("json") == "true"

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
