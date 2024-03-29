package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logLevel string

func init() {
	// default log level
	log.SetLevel(log.ErrorLevel)
}

// SetLogLevel sets the log level based on the flag passed in
func SetLogLevel(cmd *cobra.Command, args []string) {
	lowercased := strings.ToLower(logLevel)
	switch lowercased {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}
}
