package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

// ToggleDebug toggles debug logging on or off.
func ToggleDebug(cmd *cobra.Command, args []string) {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
