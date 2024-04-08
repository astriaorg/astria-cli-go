package sequencer

import (
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getBlockheightCmd represents the get-blockheight command
var getBlockheightCmd = &cobra.Command{
	Use:    "get-blockheight",
	Short:  "Retrieves the latest blockheight from the sequencer.",
	PreRun: cmd.SetLogLevel,
	Run:    runGetBlockheight,
}

func init() {
	sequencerCmd.AddCommand(getBlockheightCmd)
	getBlockheightCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
}

func runGetBlockheight(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()

	blockheight, err := sequencer.GetBlockheight(url)
	if err != nil {
		log.WithError(err)
		return
	}

	fmt.Println(blockheight)
}
