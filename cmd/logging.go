package cmd

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logLevel string

const DefaultLogLevel = log.InfoLevel

func init() {
	// default log level
	log.SetLevel(DefaultLogLevel)
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
		log.SetLevel(DefaultLogLevel)
	}
}

// CreateUILog creates a log file in the provided `destDir` for the UI app. It will panic if the log file cannot be created.
func CreateUILog(destDir string) {
	// create log file for the UI app
	logPath := filepath.Join(destDir, "astria-go.log")
	appLogFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		panic(err)
	}
	log.SetOutput(appLogFile)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true, // Disable ANSI color codes
		FullTimestamp: true,
	})
	log.Debug("New log file created successfully:", logPath)
}
