package processrunner

import (
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// LogHandler is a struct that manages the writing of process logs to a file.
type LogHandler struct {
	logPath        string
	exportLogs     bool
	fileDescriptor *os.File
}

// NewLogHandler creates a new LogHandler to be used by a ProcessRunner to write
// the process logs to a file.
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

// Writeable reports if the data sent to the LogHandler.Write() function will be
// written to the log file. If Writeable() returns false, the data will not be
// written to a log file. If Writeable() returns true, the data will be written
// to the log file when the Write() function is called.
func (lh *LogHandler) Writeable() bool {
	return lh.exportLogs
}

// Write writes the data to the log file managed by the LogHandler.
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

// Close closes the log file within the LogHandler.
func (lh *LogHandler) Close() error {
	if lh.fileDescriptor != nil {
		return lh.fileDescriptor.Close()
	}
	return nil
}
