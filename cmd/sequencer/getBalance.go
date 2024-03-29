package sequencer

import (
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getBalanceCmd represents the get-balance command
var getBalanceCmd = &cobra.Command{
	Use:    "get-balance [address]",
	Short:  "Retrieves and prints the balance of an account.",
	Args:   cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	PreRun: cmd.ToggleDebug,
	Run:    runGetBalance,
}

const DefaultSequencerURL = "http://127.0.0.1:26657"

func init() {
	sequencerCmd.AddCommand(getBalanceCmd)
	getBalanceCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
}

func runGetBalance(cmd *cobra.Command, args []string) {
	address := args[0]
	url := cmd.Flag("url").Value.String()

	balance, err := sequencer.GetBalance(address, url)
	if err != nil {
		log.WithError(err).Error("Error getting balance")
		return
	}

	fmt.Println(balance)
}
