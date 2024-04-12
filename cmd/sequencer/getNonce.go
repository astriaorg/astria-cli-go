package sequencer

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getNonceCmd represents the get-nonce command
var getNonceCmd = &cobra.Command{
	Use:    "get-nonce [address]",
	Short:  "Retrieves and prints the nonce of an account.",
	Args:   cobra.ExactArgs(1),
	PreRun: cmd.SetLogLevel,
	Run:    runGetNonce,
}

func init() {
	sequencerCmd.AddCommand(getNonceCmd)
	getNonceCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer.")
	getNonceCmd.Flags().Bool("json", false, "Output in JSON format.")
}

func runGetNonce(cmd *cobra.Command, args []string) {
	address := args[0]
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"

	nonce, err := sequencer.GetNonce(address, url)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		panic(err)
	}

	if printJSON {
		j, err := json.MarshalIndent(nonce, "", "  ")
		if err != nil {
			log.WithError(err).Error("Error marshalling nonce to JSON")
			os.Exit(1)
		}
		pterm.Println(string(j))
	} else {
		header := []string{"Nonce"}
		data := append([][]string{header}, [][]string{{strconv.Itoa(int(nonce.Nonce))}}...)
		output, err := pterm.DefaultTable.WithHasHeader().WithSeparator(" ").WithData(data).Srender()
		if err != nil {
			log.WithError(err).Error("Error rendering table")
			os.Exit(1)
		}
		pterm.Println(output)
	}
}
