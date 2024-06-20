package cmd

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

var cliLogLevel string

const DefaultCliLogLevel = log.InfoLevel

func init() {
	// default log level
	log.SetLevel(DefaultCliLogLevel)

}

// SetLogLevel sets the log level based on the flag passed in
func SetLogLevel(cliLogLevel string) {
	lowercased := strings.ToLower(cliLogLevel)
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
		log.SetLevel(DefaultCliLogLevel)
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
	log.Debug("New log file created successfully: ", logPath)
}
