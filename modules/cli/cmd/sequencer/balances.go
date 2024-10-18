package sequencer

import (
	"fmt"
	"strings"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// balancesCmd represents the balances command
var balancesCmd = &cobra.Command{
	Use:   "balances [address]",
	Short: "Retrieves and prints the balances of an account.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   balancesCmdHandler,
}

func init() {
	SequencerCmd.AddCommand(balancesCmd)

	flagHandler := cmd.CreateCliFlagHandler(balancesCmd, cmd.EnvPrefix)
	flagHandler.BindStringFlag("network", cmd.DefaultTargetNetwork, "Configure the values to target a specific network.")
	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
	flagHandler.BindBoolFlag("json", false, "Output an account's balances in JSON format.")

	viper.RegisterAlias("balance", "balances")
}

func balancesCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandlerWithUseConfigFlag(c, cmd.EnvPrefix, "network")
	networkConfig := GetNetworkConfigFromFlags(flagHandler)
	flagHandler.SetConfig(networkConfig)

	printJSON := flagHandler.GetValue("json") == "true"
	sequencerURL := flagHandler.GetValue("sequencer-url")
	sequencerURL = AddPortToURL(sequencerURL)

	address := args[0]
	if !strings.HasPrefix(address, DefaultAddressPrefix) {
		log.Errorf("Address does not have the expected prefix: %s, address: %s", DefaultAddressPrefix, address)
		panic(fmt.Errorf("address does not have the expected prefix: %s", DefaultAddressPrefix))
	}

	balances, err := sequencer.GetBalances(address, sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error getting balances")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      balances,
		PrintJSON: printJSON,
	}
	printer.Render()
}
