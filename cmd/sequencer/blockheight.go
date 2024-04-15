package sequencer

import (
	"encoding/json"
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
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

	// TODO - abstract table and json printing logic to helper functions
	if printJSON {
		j, err := json.MarshalIndent(blockheight, "", "  ")
		if err != nil {
			log.WithError(err).Error("Error marshalling account to JSON")
			panic(err)
		}
		fmt.Println(string(j))
	} else {
		header := []string{"Blockheight"}
		//strBH := strconv.FormatInt(blockheight, 10)
		//row := []string{strBH}
		data := pterm.TableData{header, []string{fmt.Sprintf("%d", blockheight.Blockheight)}}
		output, err := pterm.DefaultTable.WithHasHeader().WithSeparator(" ").WithData(data).Srender()
		if err != nil {
			log.WithError(err).Error("Error rendering table")
			panic(err)
		}
		pterm.Println(output)
	}
}
