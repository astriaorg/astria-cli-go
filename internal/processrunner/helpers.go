package processrunner

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// GetEnvironment reads the environment variables from the file at filePath and
// returns a list of environment variables in the form key=value. It will panic
// if the file can't be loaded.
func GetEnvironment(filePath string) []string {
	envMap, err := godotenv.Read(filePath)
	if err != nil {
		log.Fatalf("Error loading environment file: %v", err)
	}
	if err != nil {
		panic(fmt.Sprintf("Error loading environment file: %v", err))
	}
	var envList []string
	for key, value := range envMap {
		envList = append(envList, key+"="+value)
	}
	return envList
}
