package devtools

import (
	"fmt"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// // ReadEnvVariables reads the environment variables from the src file as a map of key-value pairs.
// func ReadEnvVariables(src string) (map[string]string, error) {
// 	envMap, err := godotenv.Read(src)
// 	if err != nil {
// 		log.Fatalf("Error loading environment file: %v", err)
// 	}
// 	return envMap, err
// }

// // LoadEnvironment loads the environment variables from the file at filePath and
// // returns a list of environment variables in the form key=value. It will panic
// // if the file can't be loaded.
// func LoadEnvironment(filePath string) []string {
// 	envMap, err := ReadEnvVariables(filePath)
// 	if err != nil {
// 		panic(fmt.Sprintf("Error loading environment file: %v", err))
// 	}
// 	var envList []string
// 	for key, value := range envMap {
// 		envList = append(envList, key+"="+value)
// 	}
// 	return envList
// }

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

// CreateDirOrPanic creates a directory with the given name with 0755 permissions.
// If the directory can't be created, it will panic.
func CreateDirOrPanic(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.WithError(err).Error("Error creating data directory")
		panic(err)
	}
}

// PathExists checks if a file or directory exists at the given path.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Debug(err)
		return false
	}
	return err == nil
}
