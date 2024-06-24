package sudo

import (
	"strconv"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	sequencercmd "github.com/astriaorg/astria-cli-go/modules/cli/cmd/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
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
	sequencerURL := sequencercmd.AddPortToURL(url)

	chainId := flagHandler.GetValue("sequencer-chain-id")

	pk := args[0]
	pubKey, err := sequencercmd.PublicKeyFromText(pk)
	if err != nil {
		log.WithError(err).Error("Error decoding public key")
		return
	}

	pow := args[1]
	power, err := strconv.ParseInt(pow, 10, 64)
	if err != nil {
		log.WithError(err).Errorf("Error decoding power string to int64 %v", pow)
		return
	}

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := sequencercmd.PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return
	}

	opts := sequencer.UpdateValidatorOpts{
		AddressPrefix:    sequencercmd.DefaultAccountPrefix,
		FromKey:          from,
		PubKey:           pubKey,
		Power:            power,
		SequencerURL:     sequencerURL,
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
	flaghandler.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to update the validator on.")
	flaghandler.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	flaghandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flaghandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flaghandler.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
