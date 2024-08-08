package processrunner

import (
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type LogHandler struct {
	logPath        string
	exportLogs     bool
	fileDescriptor *os.File
}

func NewLogHandler(logPath string, exportLogs bool) *LogHandler {
	var fileDescriptor *os.File

	// conditionally create the log file
	if exportLogs {
		// Open the log file
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.Info("New log file created successfully:", logPath)
		fileDescriptor = logFile
	} else {
		fileDescriptor = nil
	}

	return &LogHandler{
		logPath:        logPath,
		exportLogs:     exportLogs,
		fileDescriptor: fileDescriptor,
	}
}

func (lh *LogHandler) Writeable() bool {
	return lh.exportLogs
}

func (lh *LogHandler) Write(data string) error {
	// Remove ANSI escape codes from data
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mGKH]`)
	cleanData := ansiRegex.ReplaceAllString(data, "")

	_, err := lh.fileDescriptor.Write([]byte(cleanData))
	if err != nil {
		log.Fatalf("error writing to logfile %s: %v", lh.logPath, err)
		return err
	}
	return nil
}

func (lh *LogHandler) Close() error {
	if lh.fileDescriptor != nil {
		return lh.fileDescriptor.Close()
	}
	return nil
}
