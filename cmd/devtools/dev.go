package devtools

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/cmd/devtools/config"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:    "dev",
	Short:  "Commands for running the Astria Shared Sequencer.",
	PreRun: cmd.SetLogLevel,
}

func init() {
	cmd.RootCmd.AddCommand(devCmd)

	flagHandler := cmd.CreateCliStringFlagHandler(devCmd, cmd.EnvPrefix)
	flagHandler.BindPersistentFlag("instance", config.DefaultInstanceName, "Choose the target instance for purging.")
	flagHandler.BindPersistentFlag("local-network-name", "sequencer-test-chain-0", "Set the local network name for the instance. This is used to set the chain ID in the CometBFT genesis.json file.")
	flagHandler.BindPersistentFlag("local-default-denom", "nria", "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")
}
