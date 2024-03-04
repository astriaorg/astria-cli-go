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
	// TODO: make this check more sophisticated then just "is the port in use"
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
	// all paths that should exist
	paths := []string{
		"./local-dev-astria/.env",
		"./local-dev-astria/astria-sequencer",
		"./local-dev-astria/astria-conductor",
		"./local-dev-astria/astria-composer",
		"./local-dev-astria/cometbft",
		"./local-dev-astria/genesis.json",
		"./local-dev-astria/priv_validator_key.json",
		"./data",
	}
	status := true

	for _, path := range paths {
		if !exists(path) {
			fmt.Println("no", path)
			status = false
		}
	}
	return status
}

// TODO: are there any other terminals that should be supported?
var terminalEmulators = []struct {
	command string
	args    []string
}{
	{"x-terminal-emulator", []string{"-e"}}, // Debian alternatives system
	{"gnome-terminal", []string{"--"}},      // GNOME
	{"konsole", []string{"-e"}},             // KDE
	{"xfce4-terminal", []string{"-e"}},      // XFCE
	{"lxterminal", []string{"-e"}},          // LXDE
	{"mate-terminal", []string{"-e"}},       // MATE
	{"terminator", []string{"-e"}},          // Terminator
	{"tilix", []string{"-e"}},               // Tilix
	{"xterm", []string{"-e"}},               // XTerm
}

// openTerminal attempts to open a new terminal window running the specified command.
func runLinuxCommand(command string) bool {
	for _, emulator := range terminalEmulators {
		if path, err := exec.LookPath(emulator.command); err == nil {
			// Command found, attempt to execute it
			args := append(emulator.args, command)
			cmd := exec.Command(path, args...)
			if err := cmd.Start(); err == nil {
				// Successfully started the terminal emulator
				return true
			}
		}
	}
	return false // No known terminal emulator found or succeeded in opening
}

func executeCommand(cmdIn string, env []string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// TODO: finish fixing the extra terminal window issue
		// fullCmd := `tell application "Terminal"
		// 	if (count of windows) = 1 then
		// 		tell application "Terminal" to do script "` + cmdIn + `" in window 1
		// 	else
		// 		tell application "Terminal" to do script "` + cmdIn + `"
		// 	end if
		// end tell
		// `
		// cmd = exec.Command("osascript", "-e", fullCmd)
		cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "`+cmdIn+`"`)

	case "linux":
		didRun := runLinuxCommand(cmdIn)
		if !didRun {
			panic("No terminal emulator found")
		}
		// TODO: using gnome-terminal for now, but need to add support for other terminals?
		// cmd = exec.Command("gnome-terminal", "--", "bash", "-c", cmdIn)

	default:
		panic("Unsupported OS")
	}
	cmd.Env = env

	err := cmd.Start()
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
	rollupExecutionPort := 50051
	if !checkPortInUse(rollupExecutionPort) {
		fmt.Printf("Error: no rollup execution rpc detected on port %d\n", rollupExecutionPort)
		return
	}
	// TODO: make the port configurable
	rollupRpcPort := 8546
	if !checkPortInUse(rollupRpcPort) {
		fmt.Printf("Error: no rollup rpc detected on port %d\n", rollupRpcPort)
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
