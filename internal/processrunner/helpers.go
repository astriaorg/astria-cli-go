package processrunner

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// ReadEnvVariables reads the environment variables from the src file as a map of key-value pairs.
func ReadEnvVariables(src string) (map[string]string, error) {
	envMap, err := godotenv.Read(src)
	if err != nil {
		log.Fatalf("Error loading environment file: %v", err)
	}
	return envMap, err
}

// LoadEnvironment loads the environment variables from the file at filePath and
// returns a list of environment variables in the form key=value. It will panic
// if the file can't be loaded.
func LoadEnvironment(filePath string) []string {
	envMap, err := ReadEnvVariables(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error loading environment file: %v", err))
	}
	var envList []string
	for key, value := range envMap {
		envList = append(envList, key+"="+value)
	}
	return envList
}
