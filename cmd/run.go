/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func checkInstalled(cmds ...string) error {
	cmd := exec.Command(cmds[0], cmds[1:]...)

	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running command: %s", err)
	}

	return nil
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
	if !exists("./local-dev-astria/justfile") {
		fmt.Println("no justfile")
		status = false
	}
	if !exists("./local-dev-astria/mprocs.yaml") {
		fmt.Println("no mprocs.yaml")
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

func run() {
	// Check if just is installed
	err := checkInstalled("just")
	if err != nil {
		fmt.Println("Error: do you have 'just' installed?", err)
		return
	}

	// Check if mprocs is installed
	err = checkInstalled("mprocs", "-h")
	if err != nil {
		fmt.Println("Error: do you have 'mprocs' installed?", err)
		return
	}

	// Check if a rollup is running on the default port
	// TODO: make the port configurable
	rollupPort := 50051
	if !checkPortInUse(rollupPort) {
		fmt.Printf("Error: no rollup running on port %d\n", rollupPort)
		return

	}

	// TODO: check if the data and local-dev-astria directories exist
	if !checkIfInitialized() {
		fmt.Println("Error: one or more required files not present. Did you run 'astria-dev init'?")
		return
	}

	// Create the mprocs command
	cmd := exec.Command("mprocs")
	cmd.Dir = "local-dev-astria"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Stderr.WriteString("Error running command: " + err.Error() + "\n")
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
