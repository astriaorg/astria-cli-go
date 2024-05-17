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

	flagHandler := cmd.CreateCliFlagHandler(balancesCmd, cmd.EnvPrefix)
	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
	flagHandler.BindBoolFlag("json", false, "Output an account's balances in JSON format.")

	viper.RegisterAlias("balance", "balances")
}

func balancesCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	url := flagHandler.GetValue("sequencer-url")
	printJSON := flagHandler.GetValue("json") == "true"

	address := args[0]

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
