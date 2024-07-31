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

	useNetworkPreset := flagHandler.GetChanged("network")
	var networkSettings sequencercmd.SequencerNetworkConfig
	if useNetworkPreset {
		network := flagHandler.GetValue("network")
		networksConfigPath := sequencercmd.BuildSequencerNetworkConfigsFilepath()
		sequencercmd.CreateSequencerNetworkConfigs(networksConfigPath)
		networkSettings = sequencercmd.GetSequencerNetworkSettingsFromConfig(network, networksConfigPath)
	} else {
		log.Info("Target network not specified. Using flag values.")
	}

	printJSON := flagHandler.GetValue("json") == "true"

	url := sequencercmd.ChooseFlagValue(
		useNetworkPreset,
		flagHandler.GetChanged("sequencer-url"),
		networkSettings.SequencerURL,
		flagHandler.GetValue("sequencer-url"),
	)
	sequencerURL := sequencercmd.AddPortToURL(url)

	chainId := sequencercmd.ChooseFlagValue(
		useNetworkPreset,
		flagHandler.GetChanged("sequencer-chain-id"),
		networkSettings.SequencerChainId,
		flagHandler.GetValue("sequencer-chain-id"),
	)

	pk := args[0]
	pubKey, err := sequencercmd.PublicKeyFromText(pk)
	if err != nil {
		log.WithError(err).Error("Error decoding public key")
		panic(err)
	}

	pow := args[1]
	power, err := strconv.ParseInt(pow, 10, 64)
	if err != nil {
		log.WithError(err).Errorf("Error decoding power string to int64 %v", pow)
		panic(err)
	}

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := sequencercmd.PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		panic(err)
	}

	isAsync := flagHandler.GetValue("async") == "true"

	opts := sequencer.UpdateValidatorOpts{
		IsAsync:          isAsync,
		AddressPrefix:    sequencercmd.DefaultAddressPrefix,
		FromKey:          from,
		PubKey:           pubKey,
		Power:            power,
		SequencerURL:     sequencerURL,
		SequencerChainID: chainId,
	}
	tx, err := sequencer.UpdateValidator(opts)
	if err != nil {
		log.WithError(err).Error("Error updating validator")
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
	flaghandler.BindStringFlag("network", sequencercmd.DefaultTargetNetwork, "Configure the values to target a specific network.")
	flaghandler.BindBoolFlag("json", false, "Output the command result in JSON format.")
	flaghandler.BindBoolFlag("async", false, "If true, the function will return immediately. If false, the function will wait for the transaction to be seen on the network.")
	flaghandler.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to update the validator on.")
	flaghandler.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	flaghandler.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	flaghandler.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	flaghandler.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
