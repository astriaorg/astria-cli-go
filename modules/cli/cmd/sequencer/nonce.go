package sequencer

import (
	"fmt"
	"strings"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// nonceCmd represents the nonce command
var nonceCmd = &cobra.Command{
	Use:   "nonce [address]",
	Short: "Retrieves and prints the nonce of an account.",
	Args:  cobra.ExactArgs(1),
	Run:   nonceCmdHandler,
}

func init() {
	SequencerCmd.AddCommand(nonceCmd)

	flagHandler := cmd.CreateCliFlagHandler(nonceCmd, cmd.EnvPrefix)

	flagHandler.BindStringPFlag("sequencer-url", "u", DefaultSequencerURL, "The URL of the sequencer.")
	flagHandler.BindBoolFlag("json", false, "Output in JSON format.")
}

func nonceCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"

	url := flagHandler.GetValue("sequencer-url")
	sequencerURL := AddPortToURL(url)

	address := args[0]
	if !strings.HasPrefix(address, DefaultAddressPrefix) {
		log.Errorf("Address does not have the expected prefix: %s, address: %s", DefaultAddressPrefix, address)
		panic(fmt.Errorf("address does not have the expected prefix: %s", DefaultAddressPrefix))
	}

	nonce, err := sequencer.GetNonce(address, sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      nonce,
		PrintJSON: printJSON,
	}
	printer.Render()
}
