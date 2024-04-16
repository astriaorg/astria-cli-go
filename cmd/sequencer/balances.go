package sequencer

import (
	"encoding/json"
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// balancesCmd represents the balances command
var balancesCmd = &cobra.Command{
	Use:    "balances [address]",
	Short:  "Retrieves and prints the balances of an account.",
	Args:   cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	PreRun: cmd.SetLogLevel,
	Run:    balancesCmdHandler,
}

func init() {
	sequencerCmd.AddCommand(balancesCmd)
	balancesCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
	balancesCmd.Flags().Bool("json", false, "Output an account's balances in JSON format.")

	viper.RegisterAlias("balance", "balances")
}

func balancesCmdHandler(cmd *cobra.Command, args []string) {
	address := args[0]
	url := cmd.Flag("url").Value.String()
	printJSON := cmd.Flag("json").Value.String() == "true"

	balances, err := sequencer.GetBalances(address, url)
	if err != nil {
		log.WithError(err).Error("Error getting balance")
		return
	}

	// TODO - abstract table and json printing logic to helper functions
	if printJSON {
		j, err := json.MarshalIndent(balances, "", "  ")
		if err != nil {
			log.WithError(err).Error("Error marshalling account to JSON")
			panic(err)
		}
		fmt.Println(string(j))
	} else {
		header := []string{"Denom", "Balance"}
		var rows [][]string
		for _, balance := range balances {
			rows = append(rows, []string{balance.Denom, fmt.Sprintf("%d", balance.Balance)})
		}
		data := append([][]string{header}, rows...)
		output, err := pterm.DefaultTable.WithHasHeader().WithSeparator(" ").WithData(data).Srender()
		if err != nil {
			log.WithError(err).Error("Error rendering table")
			panic(err)
		}
		pterm.Println(output)
	}
}