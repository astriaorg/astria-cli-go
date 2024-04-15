package devtools

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type runOpts struct {
	ctx          context.Context
	instanceDir  string
	monoRepoPath string
}

// loadEnvVariables loads the environment variables from the src file
func loadEnvVariables(src ...string) {
	err := godotenv.Load(src...)
	if err != nil {
		log.Fatalf("Error loading environment file: %v", err)
	}
}

// getEnvList returns a list of environment variables in the form key=value
func getEnvList() []string {
	var envList []string
	for _, env := range os.Environ() {
		// Each string is in the form key=value
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		envList = append(envList, key+"="+os.Getenv(key))
	}
	return envList
}

// loadAndGetEnvVariables loads the environment variables from the src file and returns a list of environment variables in the form key=value
func loadAndGetEnvVariables(filePath ...string) []string {
	loadEnvVariables(filePath...)
	return getEnvList()
}

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

// CreateDirOrPanic creates a directory with the given name with 0755 permissions.
// If the directory can't be created, it will panic.
func CreateDirOrPanic(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.WithError(err).Error("Error creating data directory")
		panic(err)
	}

}
