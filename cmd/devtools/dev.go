package devtools

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/cmd/devtools/config"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Commands for running the Astria Shared Sequencer.",
}

func init() {
	cmd.RootCmd.AddCommand(devCmd)
	devCmd.PersistentFlags().StringP("instance", "i", config.DefaultInstanceName, "Choose the target instance for purging.")
	devCmd.PersistentFlags().String("local-network-name", "sequencer-test-chain-0", "Set the local network name for the instance. This is used to set the chain ID in the CometBFT genesis.json file.")
	devCmd.PersistentFlags().String("local-default-denom", "nria", "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")

}
