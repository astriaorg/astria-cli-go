package devtools

import (
	"fmt"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// IsInstanceNameValidOrPanic checks if the instance name is valid and panics if it's not.
func IsInstanceNameValidOrPanic(instance string) {
	re, err := regexp.Compile(`^[a-z]+[a-z0-9]*(-[a-z0-9]+)*$`)
	if err != nil {
		log.WithError(err).Error("Error compiling regex")
		panic(err)
	}
	if !re.MatchString(instance) {
		log.Errorf("Invalid instance name: %s", instance)
		err := fmt.Errorf(`
Invalid instance name: '%s'. Instance names must be lowercase, alphanumeric, 
and may contain dashes. It can't begin or end with a dash. No repeating dashes.
`, instance)
		panic(err)
	}
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

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		log.WithField("path", path).Error("The path is not a regular file")
		return false
	}

	// Check if the file is executable
	if fileInfo.Mode().Perm()&0111 == 0 {
		log.WithField("path", path).Error("The file is not executable")

		return false
	}

	return true
}
