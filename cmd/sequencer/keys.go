package sequencer

import (
	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/keys"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var setKeyCmd = &cobra.Command{
	Use:    "setkey [address] [private key]",
	Short:  "Set private key for an address in system keyring.",
	Args:   cobra.ExactArgs(2),
	PreRun: cmd.SetLogLevel,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		val := args[1]

		err := keys.StoreKeyring(key, val)
		if err != nil {
			panic(err)
		}
	},
}

var getKeyCmd = &cobra.Command{
	Use:    "getkey [address]",
	Short:  "Get private key for an address in system keyring.",
	Args:   cobra.ExactArgs(1),
	PreRun: cmd.SetLogLevel,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		val, err := keys.GetKeyring(key)
		if err != nil {
			panic(err)
		}
		log.Infof("value: %s", val)
	},
}

func init() {
	sequencerCmd.AddCommand(setKeyCmd)
	sequencerCmd.AddCommand(getKeyCmd)
}
