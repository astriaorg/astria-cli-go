package sudo

import (
	"github.com/astria/astria-cli-go/cmd"
	sequencercmd "github.com/astria/astria-cli-go/cmd/sequencer"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// validatorUpdateCmd represents the validator update command
var validatorUpdateCmd = &cobra.Command{
	Use:   "validator-update [public key] [power] [--keyfile | --keyring-address | --privkey]",
	Short: "Update a validator on the sequencer.",
	Long:  `Update a validator on the sequencer. The user needs to provide a public key and a power level for the validator.`,
	Args:  cobra.ExactArgs(2),
	Run:   validatorUpdateCmdHandler,
}

func validatorUpdateCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	pubKey := args[0]
	power := args[1]

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.UpdateValidatorOpts{
		FromKey:          priv,
		PubKey:           pubKey,
		Power:            power,
		SequencerURL:     url,
		SequencerChainID: chainId,
	}
	tx, err := sequencer.UpdateValidator(opts)
	if err != nil {
		log.WithError(err).Error("Error minting tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

func init() {
	sudoCmd.AddCommand(validatorUpdateCmd)

	flaghandler := cmd.CreateCliFlagHandler(validatorUpdateCmd, cmd.EnvPrefix)
	flaghandler.BindBoolFlag("json", false, "Output the command result in JSON format.")
	flaghandler.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to add the relayer address to.")
	flaghandler.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	flaghandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flaghandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flaghandler.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
