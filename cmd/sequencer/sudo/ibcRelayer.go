package sudo

import (
	"github.com/astria/astria-cli-go/cmd"
	sequencercmd "github.com/astria/astria-cli-go/cmd/sequencer"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astria/astria-cli-go/internal/ui"
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
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	address := args[0]

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.IBCRelayerOpts{
		FromKey:           priv,
		SequencerURL:      url,
		SequencerChainID:  chainId,
		IBCRelayerAddress: address,
	}
	tx, err := sequencer.AddIBCRelayer(opts)
	if err != nil {
		log.WithError(err).Error("Error transferring tokens")
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
	printJSON := flagHandler.GetValue("json") == "true"
	url := flagHandler.GetValue("sequencer-url")
	chainId := flagHandler.GetValue("sequencer-chain-id")

	address := args[0]

	priv, err := sequencercmd.GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}

	opts := sequencer.IBCRelayerOpts{
		FromKey:           priv,
		SequencerURL:      url,
		SequencerChainID:  chainId,
		IBCRelayerAddress: address,
	}
	tx, err := sequencer.RemoveIBCRelayer(opts)
	if err != nil {
		log.WithError(err).Error("Error transferring tokens")
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
	aibfh.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to add the relayer address to.")
	aibfh.BindBoolFlag("json", false, "Output the command result in JSON format.")
	aibfh.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	aibfh.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	aibfh.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	aibfh.BindStringFlag("privkey", "", "The private key of the sender.")
	addIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	addIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	IBCRelayerCmd.AddCommand(removeIBCRelayerCmd)

	ribfh := cmd.CreateCliFlagHandler(removeIBCRelayerCmd, cmd.EnvPrefix)
	ribfh.BindStringPFlag("sequencer-url", "u", sequencercmd.DefaultSequencerURL, "The URL of the sequencer to remove the relayer address from.")
	ribfh.BindBoolFlag("json", false, "Output the command result in JSON format.")
	ribfh.BindStringPFlag("sequencer-chain-id", "c", sequencercmd.DefaultSequencerChainID, "The chain ID of the sequencer.")
	ribfh.BindStringFlag("keyfile", "", "Path to secure keyfile for sender.")
	ribfh.BindStringFlag("keyring-address", "", "The address of the sender. Requires private key be stored in keyring.")
	ribfh.BindStringFlag("privkey", "", "The private key of the sender.")
	removeIBCRelayerCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	removeIBCRelayerCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
