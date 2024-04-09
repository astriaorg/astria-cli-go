package sequencer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
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
	getBlockheightCmd.Flags().Bool("json", false, "Output the account information in JSON format.")

}

func runGetBlockheight(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"

	blockheight, err := sequencer.GetBlockheight(url)
	if err != nil {
		log.WithError(err)
		return
	}

	// TODO - abstract table and json printing logic to helper functions
	if printJSON {
		obj := map[string]int64{
			"blockheight": blockheight,
		}
		j, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			log.WithError(err).Error("Error marshalling account to JSON")
			os.Exit(1)
		}
		fmt.Println(string(j))
	} else {
		header := []string{"Blockheight"}
		strBH := strconv.FormatInt(blockheight, 10)
		row := []string{strBH}
		data := pterm.TableData{header, row}
		output, err := pterm.DefaultTable.WithHasHeader().WithSeparator(" ").WithData(data).Srender()
		if err != nil {
			log.WithError(err).Error("Error rendering table")
			os.Exit(1)
		}
		pterm.Println(output)
	}
}
