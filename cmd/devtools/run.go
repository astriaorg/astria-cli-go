/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "TODO: short description",
	Long:  `TODO: long description`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func loadEnvVariables(filePath string) {
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

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

func loadAndGetEnvVariables(filePath string) []string {
	loadEnvVariables(filePath)
	return getEnvList()
}

func checkPortInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		// If we cannot listen on the port, it is likely in use
		return true
	}
	// Don't forget to close the listener if the port is not in use
	ln.Close()
	return false
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false // The file or directory does not exist
		}
	}
	return true // The file or directory exists
}

// checkIfInitialized checks if the files required for local development are present
func checkIfInitialized() bool {
	status := true
	if !exists("./local-dev-astria/.env") {
		fmt.Println("no .env file")
		status = false
	}
	if !exists("./local-dev-astria/astria-sequencer") {
		fmt.Println("no astria-sequencer")
		status = false
	}
	if !exists("./local-dev-astria/astria-conductor") {
		fmt.Println("no astria-conductor")
		status = false
	}
	if !exists("./local-dev-astria/astria-composer") {
		fmt.Println("no astria-composer")
		status = false
	}
	if !exists("./local-dev-astria/cometbft") {
		fmt.Println("no cometbft")
		status = false
	}
	if !exists("./local-dev-astria/genesis.json") {
		fmt.Println("no genesis.json")
		status = false
	}
	if !exists("./local-dev-astria/priv_validator_key.json") {
		fmt.Println("no priv_validator_key.json")
		status = false
	}
	if !exists("./data") {
		fmt.Println("no data directory")
		status = false
	}

	if status {
		return true
	} else {
		return false
	}
}

func executeCommand(cmdIn string, env []string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "`+cmdIn+`"`)
	// TODO: add support for windows
	// case "windows":
	// 	cmd = exec.Command("cmd", "/C", "start", "cmd", "/C", cmdIn)
	case "linux":
		// TODO: using gnome-terminal for now, but need to add support for other terminals?
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", cmdIn)
	default:
		panic("Unsupported OS")
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func run() {
	// TODO: make the dir name configuratble
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting cwd:", err)
		return
	}
	// Load the .env file and get the environment variables
	envPath := filepath.Join(cwd, "local-dev-astria/.env")
	environment := loadAndGetEnvVariables(envPath)

	// Check if a rollup is running on the default port
	// TODO: make the port configurable
	rollupPort := 50051
	if !checkPortInUse(rollupPort) {
		fmt.Printf("Error: no rollup running on port %d\n", rollupPort)
		return
	}
	if !checkIfInitialized() {
		fmt.Println("Error: one or more required files not present. Did you run 'astria-dev init'?")
		return
	}

	path := "cd " + filepath.Join(cwd, "local-dev-astria")

	// launch sequencer in new terminal
	cmdIn := path + " && ./astria-sequencer"
	executeCommand(cmdIn, environment)
	// launch cometbft in new terminal
	cmdIn = path + " && ./cometbft init --home ../data/.cometbft && cp genesis.json ../data/.cometbft/config/genesis.json && cp priv_validator_key.json ../data/.cometbft/config/priv_validator_key.json && sed -i '.bak' 's/timeout_commit = \\\"1s\\\"/timeout_commit = \\\"2s\\\"/g' ../data/.cometbft/config/config.toml && ./cometbft node --home ../data/.cometbft"
	executeCommand(cmdIn, environment)
	// launch composer in new terminal
	cmdIn = path + " && ./astria-composer"
	executeCommand(cmdIn, environment)
	// launch conductor in new terminal
	cmdIn = path + " && ./astria-conductor"
	executeCommand(cmdIn, environment)

}

func init() {
	devCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
