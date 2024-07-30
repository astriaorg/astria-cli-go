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

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Bridge commands.",
}

// bridgeInitCmd represents the `bridge init` command
var bridgeInitCmd = &cobra.Command{
	Use:   "init [rollup-name] [--keyfile | --keyring-address | --privkey]",
	Short: "Initialize a bridge account for the given rollup",
	Long: `Initialize a bridge account for the given rollup on the chain.
The sender of the transaction is used as the owner of the bridge account
and is the only actor authorized to transfer out of this account.`,
	Args: cobra.ExactArgs(1),
	Run:  bridgeInitCmdHandler,
}

func bridgeInitCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	networkDefaultsUsed := flagHandler.GetChanged("network")
	var networkSettings SequencerNetworkConfig
	if networkDefaultsUsed {
		network := flagHandler.GetValue("network")
		networksConfigPath := BuildSequencerNetworkConfigsFilepath()
		CreateSequencerNetworkConfigs(networksConfigPath)
		networkSettings = GetSequencerNetworkSettingsFromConfig(network, networksConfigPath)
	} else {
		log.Info("Target network not specified. Using flag values.")
	}

	printJSON := flagHandler.GetValue("json") == "true"

	url := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("sequencer-url"),
		networkSettings.SequencerURL,
		flagHandler.GetValue("sequencer-url"),
	)
	sequencerURL := AddPortToURL(url)

	priv, err := GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		panic(err)
	}
	fromAccount, err := sequencer.NewAccountFromPrivKey(DefaultAddressPrefix, from)
	if err != nil {
		log.WithError(err).Error("Error constructing address from private key")
	}

	rollupName := args[0]

	sequencerChainID := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("sequencer-chain-id"),
		networkSettings.SequencerChainId,
		flagHandler.GetValue("sequencer-chain-id"),
	)

	asset := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("asset"),
		networkSettings.Asset,
		flagHandler.GetValue("asset"),
	)

	feeAsset := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("fee-asset"),
		networkSettings.FeeAsset,
		flagHandler.GetValue("fee-asset"),
	)

	sa := flagHandler.GetValue("sudo-address")
	if sa == "" {
		sa = fromAccount.Address.String()
	}
	if !strings.HasPrefix(sa, DefaultAddressPrefix) && sa != "" {
		log.Errorf("sudo address does not have the expected prefix: %s, address: %s", DefaultAddressPrefix, sa)
		panic(fmt.Errorf("sudo address does not have the expected prefix: %s", DefaultAddressPrefix))
	}
	sudoAddress := AddressFromText(sa)

	wa := flagHandler.GetValue("withdrawer-address")
	if wa == "" {
		wa = fromAccount.Address.String()
	}
	if !strings.HasPrefix(wa, DefaultAddressPrefix) && wa != "" {
		log.Errorf("withdrawer address does not have the expected prefix: %s, address: %s", DefaultAddressPrefix, wa)
		panic(fmt.Errorf("withdrawer address does not have the expected prefix: %s", DefaultAddressPrefix))
	}
	withdrawerAddress := AddressFromText(wa)

	isAsync := flagHandler.GetValue("async") == "true"

	opts := sequencer.InitBridgeOpts{
		IsAsync:           isAsync,
		AddressPrefix:     DefaultAddressPrefix,
		SequencerURL:      sequencerURL,
		FromKey:           from,
		RollupName:        rollupName,
		SequencerChainID:  sequencerChainID,
		Asset:             asset,
		FeeAsset:          feeAsset,
		SudoAddress:       sudoAddress,
		WithdrawerAddress: withdrawerAddress,
	}
	bridgeAccount, err := sequencer.InitBridgeAccount(opts)
	if err != nil {
		log.WithError(err).Error("Error initializing bridge account")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      bridgeAccount,
		PrintJSON: printJSON,
	}
	printer.Render()
}

// bridgeLockCmd represents the `bridge lock` command
var bridgeLockCmd = &cobra.Command{
	Use:   "lock [amount] [to] [destination-chain-address] [--keyfile | --keyring-address | --privkey]",
	Short: "Lock tokens on the bridge account",
	Long: `A bridge lock is a transfer of tokens from the signing Sequencer
account to a Sequencer bridge account. These tokens will then be
bridged to a destination chain address if an IBC relayer is running.`,
	Args: cobra.ExactArgs(3),
	Run:  bridgeLockCmdHandler,
}

func bridgeLockCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	networkDefaultsUsed := flagHandler.GetChanged("network")
	var networkSettings SequencerNetworkConfig
	if networkDefaultsUsed {
		network := flagHandler.GetValue("network")
		networksConfigPath := BuildSequencerNetworkConfigsFilepath()
		CreateSequencerNetworkConfigs(networksConfigPath)
		networkSettings = GetSequencerNetworkSettingsFromConfig(network, networksConfigPath)
	} else {
		log.Info("Target network not specified. Using flag values.")
	}

	printJSON := flagHandler.GetValue("json") == "true"

	url := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("sequencer-url"),
		networkSettings.SequencerURL,
		flagHandler.GetValue("sequencer-url"),
	)
	sequencerURL := AddPortToURL(url)

	priv, err := GetPrivateKeyFromFlags(c)
	if err != nil {
		log.WithError(err).Error("Could not get private key from flags")
		panic(err)
	}
	from, err := PrivateKeyFromText(priv)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		panic(err)
	}

	amount, err := convertToUint128(args[0])
	if err != nil {
		log.WithError(err).Error("Error converting amount to Uint128 proto")
		panic(err)
	}

	to := args[1]
	toAddress := AddressFromText(to)

	sequencerChainID := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("sequencer-chain-id"),
		networkSettings.SequencerChainId,
		flagHandler.GetValue("sequencer-chain-id"),
	)

	asset := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("asset"),
		networkSettings.Asset,
		flagHandler.GetValue("asset"),
	)

	feeAsset := ChooseFlagValue(
		networkDefaultsUsed,
		flagHandler.GetChanged("fee-asset"),
		networkSettings.FeeAsset,
		flagHandler.GetValue("fee-asset"),
	)

	destinationChainAddress := args[2]

	isAsync := flagHandler.GetValue("async") == "true"

	opts := sequencer.BridgeLockOpts{
		IsAsync:                 isAsync,
		AddressPrefix:           DefaultAddressPrefix,
		SequencerURL:            sequencerURL,
		FromKey:                 from,
		Amount:                  amount,
		ToAddress:               toAddress,
		SequencerChainID:        sequencerChainID,
		Asset:                   asset,
		FeeAsset:                feeAsset,
		DestinationChainAddress: destinationChainAddress,
	}
	tx, err := sequencer.BridgeLock(opts)
	if err != nil {
		log.WithError(err).Error("Error locking tokens")
		panic(err)
	}

	printer := ui.ResultsPrinter{
		Data:      tx,
		PrintJSON: printJSON,
	}
	printer.Render()
}

func init() {
	SequencerCmd.AddCommand(bridgeCmd)

	bridgeCmd.AddCommand(bridgeInitCmd)
	bifh := cmd.CreateCliFlagHandler(bridgeInitCmd, cmd.EnvPrefix)
	bifh.BindStringPFlag("sequencer-chain-id", "c", DefaultDuskSequencerChainID, "The chain ID of the sequencer.")
	bifh.BindStringFlag("asset", DefaultAsset, "The name of the asset we want to bridge.")
	bifh.BindStringFlag("fee-asset", DefaultFeeAsset, "The name of the asset used for fees.")
	bifh.BindStringFlag("sudo-address", "", "Set the sudo address to use for the bridge account. The address of the sender is used if this is not set.")
	bifh.BindStringFlag("withdrawer-address", "", "Set the withdrawer address to use for the bridge account. The address of the sender is used if this is not set.")

	bifh.BindBoolFlag("json", false, "Output bridge account as JSON.")
	bifh.BindBoolFlag("async", false, "If true, the function will return immediately. If false, the function will wait for the transaction to be seen on the network.")

	bifh.BindStringPFlag("sequencer-url", "u", DefaultDuskSequencerURL, "The URL of the sequencer to init bridge account on.")

	bifh.BindStringFlag("keyfile", "", "Path to secure keyfile for the bridge account.")
	bifh.BindStringFlag("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	bifh.BindStringFlag("privkey", "", "The private key of the bridge account.")
	bridgeInitCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeInitCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")

	bridgeCmd.AddCommand(bridgeLockCmd)
	blfh := cmd.CreateCliFlagHandler(bridgeLockCmd, cmd.EnvPrefix)
	blfh.BindStringFlag("sequencer-chain-id", DefaultDuskSequencerChainID, "The chain ID of the sequencer.")
	blfh.BindStringFlag("asset", DefaultAsset, "The asset to be locked and transferred.")
	blfh.BindStringFlag("fee-asset", DefaultFeeAsset, "The asset used to pay the transaction fee.")

	blfh.BindBoolFlag("json", false, "Output bridge account as JSON")
	blfh.BindBoolFlag("async", false, "If true, the function will return immediately. If false, the function will wait for the transaction to be seen on the network.")
	blfh.BindStringPFlag("sequencer-url", "u", DefaultDuskSequencerURL, "The URL of the sequencer to lock assets on.")

	blfh.BindStringFlag("keyfile", "", "Path to secure keyfile for the bridge account.")
	blfh.BindStringFlag("keyring-address", "", "The address of the bridge account. Requires private key be stored in keyring.")
	blfh.BindStringFlag("privkey", "", "The private key of the bridge account.")
	bridgeLockCmd.MarkFlagsOneRequired("keyfile", "keyring-address", "privkey")
	bridgeLockCmd.MarkFlagsMutuallyExclusive("keyfile", "keyring-address", "privkey")
}
