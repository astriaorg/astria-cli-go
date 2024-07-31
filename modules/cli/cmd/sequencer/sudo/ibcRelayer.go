package sudo

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	sequencercmd "github.com/astriaorg/astria-cli-go/modules/cli/cmd/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// IBCRelayerCmd represents the root command for interacting with IBC Relayers
// on the sequencer.
var IBCRelayerCmd = &cobra.Command{
	Use:   "ibc-relayer",
	Short: "Interact with IBC Relayers on the sequencer.",
}

// addIBCRelayerCmd represents the add ibc relayer command
var addIBCRelayerCmd = &cobra.Command{
	Use:   "add [address] [--keyfile | --keyring-address | --privkey]",
	Short: "Add an address to the IBC Relayer set on the sequencer.",
	Args:  cobra.ExactArgs(1),
	Run:   addIBCRelayerCmdHandler,
}

func addIBCRelayerCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	useNetworkPreset := flagHandler.GetChanged("network")
	var networkSettings sequencercmd.SequencerNetworkConfig
	if useNetworkPreset {
		network := flagHandler.GetValue("network")
		networksConfigPath := sequencercmd.BuildSequencerNetworkConfigsFilepath()
		sequencercmd.CreateSequencerNetworkConfigs(networksConfigPath)
		networkSettings = sequencercmd.GetSequencerNetworkSettingsFromConfig(network, networksConfigPath)
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

	ibcAdd := args[0]
	addIbcAddress := sequencercmd.AddressFromText(ibcAdd)

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

	opts := sequencer.IBCRelayerOpts{
		IsAsync:           isAsync,
		AddressPrefix:     sequencercmd.DefaultAddressPrefix,
		FromKey:           from,
		SequencerURL:      sequencerURL,
		SequencerChainID:  chainId,
		IBCRelayerAddress: addIbcAddress,
	}
	tx, err := sequencer.AddIBCRelayer(opts)
	if err != nil {
		log.WithError(err).Error("Error adding IBC relayer")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

// removeIBCRelayerCmd represents the remove ibc relayer command
var removeIBCRelayerCmd = &cobra.Command{
	Use:   "remove [address] [--keyfile | --keyring-address | --privkey]",
	Short: "Remove an address from the IBC Relayer set on the sequencer.",
	Args:  cobra.ExactArgs(1),
	Run:   removeIBCRelayerCmdHandler,
}

func removeIBCRelayerCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	useNetworkPreset := flagHandler.GetChanged("network")
	var networkSettings sequencercmd.SequencerNetworkConfig
	if useNetworkPreset {
		network := flagHandler.GetValue("network")
		networksConfigPath := sequencercmd.BuildSequencerNetworkConfigsFilepath()
		sequencercmd.CreateSequencerNetworkConfigs(networksConfigPath)
		networkSettings = sequencercmd.GetSequencerNetworkSettingsFromConfig(network, networksConfigPath)
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

	ibcRmv := args[0]
	rmvIbcAddress := sequencercmd.AddressFromText(ibcRmv)

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

	opts := sequencer.IBCRelayerOpts{
		IsAsync:           isAsync,
		AddressPrefix:     sequencercmd.DefaultAddressPrefix,
		FromKey:           from,
		SequencerURL:      sequencerURL,
		SequencerChainID:  chainId,
		IBCRelayerAddress: rmvIbcAddress,
	}
	tx, err := sequencer.RemoveIBCRelayer(opts)
	if err != nil {
		log.WithError(err).Error("Error removing IBC relayer")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

func init() {
	sudoCmd.AddCommand(IBCRelayerCmd)
	IBCRelayerCmd.AddCommand(addIBCRelayerCmd)

	aibfh := cmd.CreateCliFlagHandler(addIBCRelayerCmd, cmd.EnvPrefix)
	aibfh.BindStringFlag("network", sequencercmd.DefaultTargetNetwork, "Configure the values to target a specific network.")
	aibfh.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to add the relayer address to.")
	aibfh.BindBoolFlag("json", false, "Output the command result in JSON format.")
	aibfh.BindBoolFlag("async", false, "If true, the function will return immediately. If false, the function will wait for the transaction to be seen on the network.")
	aibfh.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	aibfh.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	aibfh.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	aibfh.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	IBCRelayerCmd.AddCommand(removeIBCRelayerCmd)

	ribfh := cmd.CreateCliFlagHandler(removeIBCRelayerCmd, cmd.EnvPrefix)
	ribfh.BindStringFlag("network", sequencercmd.DefaultTargetNetwork, "Configure the values to target a specific network.")
	ribfh.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to remove the relayer address from.")
	ribfh.BindBoolFlag("json", false, "Output the command result in JSON format.")
	ribfh.BindBoolFlag("async", false, "If true, the function will return immediately. If false, the function will wait for the transaction to be seen on the network.")
	ribfh.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	ribfh.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	ribfh.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	ribfh.BindStringFlag("privkey", "", "The private key of the sender.")
	removeIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	removeIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
