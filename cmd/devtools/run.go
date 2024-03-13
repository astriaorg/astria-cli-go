/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
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
func checkIfInitialized(path string) bool {
	// all paths that should exist
	filePaths := []string{
		"local-dev-astria/.env",
		"local-dev-astria/astria-sequencer",
		"local-dev-astria/astria-conductor",
		"local-dev-astria/astria-composer",
		"local-dev-astria/cometbft",
		"local-dev-astria/genesis.json",
		"local-dev-astria/priv_validator_key.json",
		"data",
	}
	status := true

	for _, fp := range filePaths {
		expandedPath := filepath.Join(path, fp)
		if !exists(expandedPath) {
			fmt.Println("no", fp, "found")
			status = false
		}
	}
	return status
}

func run() {
	sequencerStartComplete := make(chan bool)
	cometbftStartComplete := make(chan bool)
	composerStartComplete := make(chan bool)
	// conductorStartComplete := make(chan bool)

	// TODO: make the home dir name configuratble
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	defaultDir := filepath.Join(homePath, ".astria")
	// Load the .env file and get the environment variables
	envPath := filepath.Join(defaultDir, "local-dev-astria/.env")
	environment := loadAndGetEnvVariables(envPath)

	// FIXME: this is a temporary ignored for easier dev
	// TODO: remove if when actually done
	if false {
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
	}
	if !checkIfInitialized(defaultDir) {
		fmt.Println("Error: one or more required files not present. Did you run 'astria-dev init'?")
		return
	}

	app := tview.NewApplication()

	// create text view object for the sequencer
	sequencerTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	sequencerTextView.SetTitle(" Sequencer ").SetBorder(true)

	// create text view object for the cometbft
	cometbftTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	cometbftTextView.SetTitle(" Cometbft ").SetBorder(true)

	composerTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	composerTextView.SetTitle(" Composer ").SetBorder(true)

	conductorTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	conductorTextView.SetTitle(" Conductor ").SetBorder(true)

	// create the sequencer command
	sequencerBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-sequencer")
	seqCmd := exec.Command(sequencerBinPath)
	seqCmd.Env = environment

	// create the cometbft command
	cometbftDataPath := filepath.Join(homePath, ".astria/data/.cometbft")
	cometbftCmdPath := filepath.Join(homePath, ".astria/local-dev-astria/cometbft")
	nodeCmdArgs := []string{"node", "--home", cometbftDataPath}
	cometbftCmd := exec.Command(cometbftCmdPath, nodeCmdArgs...)
	cometbftCmd.Env = environment

	// create the composer command
	composerBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-composer")
	composerCmd := exec.Command(composerBinPath)
	composerCmd.Env = environment

	// create the conductor command
	conductorBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-conductor")
	conductorCmd := exec.Command(conductorBinPath)
	conductorCmd.Env = environment

	// Track the current word wrap status.
	wordWrapEnabled := true
	includeAnsiEscapeCharacters := false

	// set the input capture for the app
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			if err := seqCmd.Process.Signal(syscall.SIGINT); err != nil {
				fmt.Println("Failed to send SIGINT to the process:", err)
			}
			if err := cometbftCmd.Process.Signal(syscall.SIGINT); err != nil {
				fmt.Println("Failed to send SIGINT to the process:", err)
			}
			if err := composerCmd.Process.Signal(syscall.SIGINT); err != nil {
				fmt.Println("Failed to send SIGINT to the process:", err)
			}
			if err := conductorCmd.Process.Signal(syscall.SIGINT); err != nil {
				fmt.Println("Failed to send SIGINT to the process:", err)
			}
			app.Stop()
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'w' {
			// Toggle word wrap
			sequencerTextView.SetWrap(!wordWrapEnabled)
			cometbftTextView.SetWrap(!wordWrapEnabled)
			composerTextView.SetWrap(!wordWrapEnabled)
			conductorTextView.SetWrap(!wordWrapEnabled)
			wordWrapEnabled = !wordWrapEnabled
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'e' {
			includeAnsiEscapeCharacters = !includeAnsiEscapeCharacters
		}

		return event
	})
	// Run a command and stream its output to the TextView.
	// go func() for running the sequencer
	go func() {
		// Get a pipe to the command's output.
		// TODO: read both stdout and stderr
		// stderr, err := cmd.StderrPipe()
		// if err != nil {
		// 	panic(err)
		// }
		stdout, err := seqCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := seqCmd.Start(); err != nil {
			panic(err)
		}

		sequencerStartComplete <- true

		// Create a scanner to read the output line by line.
		// TODO: read both stdout and stderr
		// stderrScanner := bufio.NewScanner(stderr)
		stdoutScanner := bufio.NewScanner(stdout)
		// output := io.MultiReader(stdout, stderr)

		aWriter := tview.ANSIWriter(sequencerTextView)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				// sequencerTextView.Write([]byte(line + "\n"))
				aWriter.Write([]byte(line + "\n"))
				// sequencerTextView.Write([]byte("\x1b[31mThis should be red.\x1b[0m\n"))
				// sequencerTextView.Write([]byte("[red]This should be red.[-]\n"))

				sequencerTextView.ScrollToEnd()
			})
		}
		if err := stdoutScanner.Err(); err != nil {
			panic(err)
		}

		if err := seqCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	// go func() for running the cometbft
	go func() {
		<-sequencerStartComplete
		initCmdArgs := []string{"init", "--home", cometbftDataPath}
		initCmd := exec.Command(cometbftCmdPath, initCmdArgs...)
		initCmd.Env = environment

		aWriter := tview.ANSIWriter(cometbftTextView)

		p := fmt.Sprintf("Running command `%v %v %v %v`\n", initCmd, initCmdArgs[0], initCmdArgs[1], initCmdArgs[2])
		app.QueueUpdateDraw(func() {
			// cometbftTextView.Write([]byte(p))
			aWriter.Write([]byte(p))
			cometbftTextView.ScrollToEnd()
		})

		out, err := initCmd.CombinedOutput()
		if err != nil {
			p := fmt.Sprintf("Error executing command `%v`: %v\n", initCmd, err)
			app.QueueUpdateDraw(func() {
				// cometbftTextView.Write([]byte(p))
				aWriter.Write([]byte(p))
				cometbftTextView.ScrollToEnd()

			})
			return
		}
		app.QueueUpdateDraw(func() {
			// cometbftTextView.Write([]byte(out))
			aWriter.Write([]byte(out))
			cometbftTextView.ScrollToEnd()

		})

		// $ cp genesis.json ../data/.cometbft/config/genesis.json
		initGenesisJsonPath := filepath.Join(homePath, ".astria/local-dev-astria/genesis.json")
		endGenesisJsonPath := filepath.Join(homePath, ".astria/data/.cometbft/config/genesis.json")
		copyArgs := []string{initGenesisJsonPath, endGenesisJsonPath}
		copyCmd := exec.Command("cp", copyArgs...)
		copyCmd.Env = environment

		_, err = copyCmd.CombinedOutput()
		if err != nil {
			p := fmt.Sprintf("Error executing command `%v`: %v\n", copyCmd, err)
			app.QueueUpdateDraw(func() {
				// cometbftTextView.Write([]byte(p))
				aWriter.Write([]byte(p))
				cometbftTextView.ScrollToEnd()

			})
			return
		}
		p = fmt.Sprintf("Copied genesis.json to %s\n", endGenesisJsonPath)
		app.QueueUpdateDraw(func() {
			// cometbftTextView.Write([]byte(p))
			aWriter.Write([]byte(p))
			cometbftTextView.ScrollToEnd()

		})

		// // $ cp priv_validator_key.json
		initPrivValidatorJsonPath := filepath.Join(homePath, ".astria/local-dev-astria/priv_validator_key.json")
		endPrivValidatorJsonPath := filepath.Join(homePath, ".astria/data/.cometbft/config/priv_validator_key.json")
		copyArgs = []string{initPrivValidatorJsonPath, endPrivValidatorJsonPath}
		copyCmd = exec.Command("cp", copyArgs...)
		copyCmd.Env = environment

		_, err = copyCmd.CombinedOutput()
		if err != nil {
			p := fmt.Sprintf("Error executing command `%v`: %v\n", copyCmd, err)
			app.QueueUpdateDraw(func() {
				// cometbftTextView.Write([]byte(p))
				aWriter.Write([]byte(p))
				cometbftTextView.ScrollToEnd()

			})
			return
		}
		p = fmt.Sprintf("Copied priv_validator_key.json to %s\n", endPrivValidatorJsonPath)
		app.QueueUpdateDraw(func() {
			// cometbftTextView.Write([]byte(p))
			aWriter.Write([]byte(p))
			cometbftTextView.ScrollToEnd()

		})

		// go code for the following sed command
		// $ sed -i '.bak' 's/timeout_commit = \\\"1s\\\"/timeout_commit = \\\"2s\\\"/g' ../data/.cometbft/config/config.toml
		cometbftConfigPath := filepath.Join(homePath, ".astria/data/.cometbft/config/config.toml")
		oldValue := `timeout_commit = "1s"`
		newValue := `timeout_commit = "2s"`

		if err := replaceInFile(cometbftConfigPath, oldValue, newValue); err != nil {
			p := fmt.Sprintf("Error updating the file: %v : %v", cometbftConfigPath, err)
			app.QueueUpdateDraw(func() {
				// cometbftTextView.Write([]byte(p))
				aWriter.Write([]byte(p))
				cometbftTextView.ScrollToEnd()

			})
			return
		} else {
			p := fmt.Sprintf("Updated %v successfully", cometbftConfigPath)
			app.QueueUpdateDraw(func() {
				// cometbftTextView.Write([]byte(p))
				aWriter.Write([]byte(p))
				cometbftTextView.ScrollToEnd()

			})
		}

		stdout, err := cometbftCmd.StdoutPipe()
		// stdout, err := cometbftCmd.StderrPipe()
		if err != nil {
			panic(err)
		}

		if err := cometbftCmd.Start(); err != nil {
			panic(err)
		}

		cometbftStartComplete <- true

		stdoutScanner := bufio.NewScanner(stdout)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				// cometbftTextView.Write([]byte(line + "\n"))
				aWriter.Write([]byte(line + "\n"))
				cometbftTextView.ScrollToEnd()

			})
		}
		if err := stdoutScanner.Err(); err != nil {
			panic(err)
		}

		if err := cometbftCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	// go func() for running the composer
	go func() {
		<-cometbftStartComplete
		// Get a pipe to the command's output.
		// TODO: read both stdout and stderr
		// stderr, err := cmd.StderrPipe()
		// if err != nil {
		// 	panic(err)
		// }
		stdout, err := composerCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := composerCmd.Start(); err != nil {
			panic(err)
		}

		composerStartComplete <- true

		// Create a scanner to read the output line by line.
		// TODO: read both stdout and stderr
		// stderrScanner := bufio.NewScanner(stderr)
		stdoutScanner := bufio.NewScanner(stdout)
		// output := io.MultiReader(stdout, stderr)

		aWriter := tview.ANSIWriter(composerTextView)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				// composerTextView.Write([]byte(line + "\n"))
				aWriter.Write([]byte(line + "\n"))
				composerTextView.ScrollToEnd()
			})
		}
		if err := stdoutScanner.Err(); err != nil {
			panic(err)
		}

		if err := composerCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	// go func() for running the conductor
	go func() {
		<-composerStartComplete
		// Get a pipe to the command's output.
		// TODO: read both stdout and stderr
		// stderr, err := cmd.StderrPipe()
		// if err != nil {
		// 	panic(err)
		// }
		stdout, err := conductorCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := conductorCmd.Start(); err != nil {
			panic(err)
		}

		// Create a scanner to read the output line by line.
		// TODO: read both stdout and stderr
		// stderrScanner := bufio.NewScanner(stderr)
		stdoutScanner := bufio.NewScanner(stdout)
		// output := io.MultiReader(stdout, stderr)

		aWriter := tview.ANSIWriter(conductorTextView)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				// conductorTextView.Write([]byte(line + "\n"))
				aWriter.Write([]byte(line + "\n"))
				conductorTextView.ScrollToEnd()
			})
		}
		if err := stdoutScanner.Err(); err != nil {
			panic(err)
		}

		if err := conductorCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	// Create a new Flex layout.
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sequencerTextView, 0, 1, false).
		AddItem(cometbftTextView, 0, 1, false).
		AddItem(composerTextView, 0, 1, false).
		AddItem(conductorTextView, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

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

func replaceInFile(filename, oldValue, newValue string) error {
	// Read the original file.
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read the file: %w", err)
	}

	// Perform the replacement.
	modifiedContent := strings.ReplaceAll(string(content), oldValue, newValue)

	// Write the modified content to a new temporary file.
	tmpFilename := filename + ".tmp"
	if err := os.WriteFile(tmpFilename, []byte(modifiedContent), 0666); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Rename the original file to filename.bak.
	backupFilename := filename + ".bak"
	if err := os.Rename(filename, backupFilename); err != nil {
		return fmt.Errorf("failed to rename original file to backup: %w", err)
	}

	// Rename the temporary file to the original file name.
	if err := os.Rename(tmpFilename, filename); err != nil {
		// Attempt to restore the original file if renaming fails.
		os.Rename(backupFilename, filename)
		return fmt.Errorf("failed to rename temporary file to original: %w", err)
	}

	return nil
}
