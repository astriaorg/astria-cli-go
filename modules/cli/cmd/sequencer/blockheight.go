package sequencer

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// blockheightCmd represents the blockheight command
var blockheightCmd = &cobra.Command{
	Use:   "blockheight",
	Short: "Retrieves the latest blockheight from the sequencer.",
	Run:   blockheightCmdHandler,
}

func init() {
	SequencerCmd.AddCommand(blockheightCmd)

	flagHandler := cmd.CreateCliFlagHandler(blockheightCmd, cmd.EnvPrefix)
	flagHandler.BindStringFlag("network", DefaultTargetNetwork, "Configure the values to target a specific network.")
	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
	flagHandler.BindBoolFlag("json", false, "Output an account's balances in JSON format.")
}

func blockheightCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandlerWithUseConfigFlag(c, cmd.EnvPrefix, "network")
	networkConfig := GetNetworkConfigFromFlags(flagHandler)
	flagHandler.SetConfig(networkConfig)

	printJSON := flagHandler.GetValue("json") == "true"
	sequencerURL := flagHandler.GetValue("sequencer-url")
	sequencerURL = AddPortToURL(sequencerURL)

	blockheight, err := sequencer.GetBlockheight(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error getting blockheight")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      blockheight,
		PrintJSON: printJSON,
	}
	printer.Render()
}
