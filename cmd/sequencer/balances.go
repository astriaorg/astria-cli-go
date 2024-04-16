package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
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

	printer := ui.ResultsPrinter{
		Data:      balances,
		PrintJSON: printJSON,
	}
	printer.Render()
}
