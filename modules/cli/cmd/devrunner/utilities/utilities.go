package utilities

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// CopyFile copies the contents of the src file to dst.
// If dst does not exist, it will be created, and if it does, it will be overwritten.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// open the destination file for writing.
	// create the file if it does not exist, truncate it if it does.
	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// ensure that any writes made to the destination file are committed to stable storage.
	err = destFile.Sync()
	return err
}

// PathExists checks if the file or binary for the input path is a regular file
// and is executable. A regular file is one where no mode type bits are set.
func PathExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.WithError(err).Error("File does not exist")
		} else {
			log.WithError(err).Error("Error checking file")
		}
		return false
	}

	// check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		log.WithField("path", path).Error("The path is not a regular file")
		return false
	}

	// check if the file is executable
	if fileInfo.Mode().Perm()&0111 == 0 {
		log.WithField("path", path).Error("The file is not executable")

		return false
	}

	return true
}
