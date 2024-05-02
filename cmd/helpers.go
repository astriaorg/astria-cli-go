package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// CreateDirOrPanic creates a directory with the given name with 0755 permissions.
// If the directory can't be created, it will panic.
func CreateDirOrPanic(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.WithError(err).Error("Error creating data directory")
		panic(err)
	}
}
