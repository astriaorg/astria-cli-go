package sequencer

import (
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getBlockHeightCmd represents the get-blockheight command
var getBlockHeightCmd = &cobra.Command{
	Use:    "get-blockheight",
	Short:  "Retrieves the latest blockheight from the sequencer.",
	PreRun: cmd.SetLogLevel,
	Run:    runGetBlockHeight,
}

func init() {
	sequencerCmd.AddCommand(getBlockHeightCmd)
	getBlockHeightCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
}

func runGetBlockHeight(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()

	blockheight, err := sequencer.GetBlockHeight(url)
	if err != nil {
		log.WithError(err).Error("Error getting block height")
		return
	}

	fmt.Println(blockheight)
}
