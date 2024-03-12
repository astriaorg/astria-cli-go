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

// TODO: are there any other terminals that should be supported?
// var terminalEmulators = []struct {
// 	command string
// 	args    []string
// }{
// 	{"x-terminal-emulator", []string{"-e"}}, // Debian alternatives system
// 	{"gnome-terminal", []string{"--"}},      // GNOME
// 	{"konsole", []string{"-e"}},             // KDE
// 	{"xfce4-terminal", []string{"-e"}},      // XFCE
// 	{"lxterminal", []string{"-e"}},          // LXDE
// 	{"mate-terminal", []string{"-e"}},       // MATE
// 	{"terminator", []string{"-e"}},          // Terminator
// 	{"tilix", []string{"-e"}},               // Tilix
// 	{"xterm", []string{"-e"}},               // XTerm
// }

// openTerminal attempts to open a new terminal window running the specified command.
// func runLinuxCommand(command string) bool {
// 	for _, emulator := range terminalEmulators {
// 		if path, err := exec.LookPath(emulator.command); err == nil {
// 			// Command found, attempt to execute it
// 			args := append(emulator.args, command)
// 			cmd := exec.Command(path, args...)
// 			if err := cmd.Start(); err == nil {
// 				// Successfully started the terminal emulator
// 				return true
// 			}
// 		}
// 	}
// 	return false // No known terminal emulator found or succeeded in opening
// }

// func executeCommand(cmdIn string, env []string) {
// 	var cmd *exec.Cmd

// 	switch runtime.GOOS {
// 	case "darwin":
// 		// TODO: finish fixing the extra terminal window issue OR just move on
// 		// to the TUI
// 		// fullCmd := `tell application "Terminal"
// 		// 	if (count of windows) = 1 then
// 		// 		tell application "Terminal" to do script "` + cmdIn + `" in window 1
// 		// 	else
// 		// 		tell application "Terminal" to do script "` + cmdIn + `"
// 		// 	end if
// 		// end tell
// 		// `
// 		// cmd = exec.Command("osascript", "-e", fullCmd)
// 		cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "`+cmdIn+`"`)

// 	case "linux":
// 		didRun := runLinuxCommand(cmdIn)
// 		if !didRun {
// 			panic("No terminal emulator found")
// 		}

// 	default:
// 		panic("Unsupported OS")
// 	}
// 	cmd.Env = env

// 	err := cmd.Start()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func run() {
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

	// path := "cd " + filepath.Join(defaultDir, "local-dev-astria")

	// // launch sequencer in new terminal
	// cmdIn := path + " && ./astria-sequencer"
	// executeCommand(cmdIn, environment)
	// // launch cometbft in new terminal
	// // TODO: think about the relative vs absolute path for this command
	// cmdIn = path + " && ./cometbft init --home ../data/.cometbft && cp genesis.json ../data/.cometbft/config/genesis.json && cp priv_validator_key.json ../data/.cometbft/config/priv_validator_key.json && sed -i '.bak' 's/timeout_commit = \\\"1s\\\"/timeout_commit = \\\"2s\\\"/g' ../data/.cometbft/config/config.toml && ./cometbft node --home ../data/.cometbft"
	// executeCommand(cmdIn, environment)
	// // launch composer in new terminal
	// cmdIn = path + " && ./astria-composer"
	// executeCommand(cmdIn, environment)
	// // launch conductor in new terminal
	// cmdIn = path + " && ./astria-conductor"
	// executeCommand(cmdIn, environment)
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

	// get the HOME path
	// TODO: make this configureable
	// homePath, err := os.UserHomeDir()
	// if err != nil {
	// 	fmt.Println("could not get home dir", err)
	// 	return
	// }

	// envPath := filepath.Join(homePath, ".astria/local-dev-astria/.env")
	// environment := loadAndGetEnvVariables(envPath)

	// Run a command and stream its output to the TextView.
	// go func() for running the sequencer
	go func() {
		binPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-sequencer")

		// for _, env := range environment {
		// 	fmt.Println(env)
		// }
		seqCmd := exec.Command(binPath)
		seqCmd.Env = environment

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlC {
				if err := seqCmd.Process.Signal(syscall.SIGINT); err != nil {
					fmt.Println("Failed to send SIGINT to the process:", err)
				}
				app.Stop()
				return nil
			}
			return event
		})

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

		// Create a scanner to read the output line by line.
		// TODO: read both stdout and stderr
		// output := io.MultiReader(stdout, stderr)
		// scanner := bufio.NewScanner(output)
		// scanner := bufio.NewScanner(stderr)
		// stderrScanner := bufio.NewScanner(stderr)
		stdoutScanner := bufio.NewScanner(stdout)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				sequencerTextView.Write([]byte(line + "\n"))
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
		// app.QueueUpdateDraw(func() {
		// 	cometbftTextView.Write([]byte("entered\n"))
		// })
		cometbftDataPath := filepath.Join(homePath, ".astria/data/.cometbft")
		cometbftCmdPath := filepath.Join(homePath, ".astria/local-dev-astria/cometbft")
		initCmdArgs := []string{"init", "--home", cometbftDataPath}
		// p := fmt.Sprintf("Running: `%v`\n", initCmdPath)
		// app.QueueUpdateDraw(func() {
		// 	cometbftTextView.Write([]byte(p))
		// })
		initCmd := exec.Command(cometbftCmdPath, initCmdArgs...)
		// initCmd := exec.Command(initCmdPath)
		initCmd.Env = environment

		out, err := initCmd.CombinedOutput()
		if err != nil {
			p := fmt.Sprintf("Error executing command `%v`: %v\n", initCmd, err)
			app.QueueUpdateDraw(func() {
				cometbftTextView.Write([]byte(p))
			})
			return
		}
		// p := fmt.Sprintf("Output of command `%v`: %s\n", initCmd, out)
		app.QueueUpdateDraw(func() {
			cometbftTextView.Write([]byte(out))
		})

		// app.QueueUpdateDraw(func() {
		// 	cometbftTextView.Write([]byte("command created\n"))
		// })

		// stdout, err := initCmd.StdoutPipe()
		// if err != nil {
		// 	panic(err)
		// }

		// app.QueueUpdateDraw(func() {
		// 	cometbftTextView.Write([]byte("scanner created\n"))
		// })
		// stdoutScanner := bufio.NewScanner(stdout)

		// if err := initCmd.Run(); err != nil {
		// 	panic(err)
		// }
		// if err := initCmd.Start(); err != nil {
		// 	panic(err)
		// }

		// for stdoutScanner.Scan() {
		// 	line := stdoutScanner.Text()
		// 	app.QueueUpdateDraw(func() {
		// 		cometbftTextView.Write([]byte(line + "\n"))
		// 	})
		// }
		// if err := stdoutScanner.Err(); err != nil {
		// 	panic(err)
		// }
		// if err := initCmd.Wait(); err != nil {
		// 	panic(err)
		// }
		// app.QueueUpdateDraw(func() {
		// 	cometbftTextView.Write([]byte("done\n"))
		// })

		// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 	if event.Key() == tcell.KeyCtrlC {
		// 		if err := initCmd.Process.Signal(syscall.SIGINT); err != nil {
		// 			fmt.Println("Failed to send SIGINT to the process:", err)
		// 		}
		// 		// app.Stop()
		// 		return nil
		// 	}
		// 	return event
		// })

		// stderr, err := initCmd.StderrPipe()
		// if err != nil {
		// 	panic(err)
		// }
		// stdout, err := initCmd.StdoutPipe()
		// if err != nil {
		// 	panic(err)
		// }

		// if err := initCmd.Start(); err != nil {
		// 	panic(err)
		// }
		// stdoutScanner := bufio.NewScanner(stdout)
		// stdoutScanner := bufio.NewScanner(stderr)

		// for stdoutScanner.Scan() {
		// for stdoutScanner.Scan() {
		// line := stdoutScanner.Text()
		// app.QueueUpdateDraw(func() {
		// 	cometbftTextView.SetText("testing\n")
		// 	// cometbftTextView.Write([]byte(line + "\n"))
		// })
		// }
		// if err := stdoutScanner.Err(); err != nil {
		// 	panic(err)
		// }

		// if err := initCmd.Run(); err != nil {
		// 	panic(err)
		// }
		// 	executeFiniteCommand(app, initCmd, cometbftTextView)

		// // $ cp genesis.json ../data/.cometbft/config/genesis.json
		// initGenesisJsonPath := filepath.Join(homePath, ".astria/local-dev-astria/genesis.json")
		// endGenesisJsonPath := filepath.Join(homePath, ".astria/data/.cometbft/config/genesis.json")
		// copyGenesisJsonPath := "cp " + initGenesisJsonPath + " " + endGenesisJsonPath
		// // run cp command here

		// executeFiniteCommand(app, cmd*exec.Cmd, textView*tview.TextView)

		// // $ cp priv_validator_key.json
		// initPrivValidatorJsonPath := filepath.Join(homePath, ".astria/local-dev-astria/priv_validator_key.json")
		// endPrivValidatorJsonPath := filepath.Join(homePath, ".astria/data/.cometbft/config/priv_validator_key.json")
		// copyPrivValidatorJsonPath := "cp " + initPrivValidatorJsonPath + " " + endPrivValidatorJsonPath
		// // run cp command here

		// // $ sed -i '.bak' 's/timeout_commit = \\\"1s\\\"/timeout_commit = \\\"2s\\\"/g' ../data/.cometbft/config/config.toml
		// cometbftConfigPath := filepath.Join(homePath, ".astria/data/.cometbft/config/config.toml")
		// cometbftSedCommand := "sed -i '.bak' 's/timeout_commit = \\\"1s\\\"/timeout_commit = \\\"2s\\\"/g' " + cometbftConfigPath
		// // run sed command here

		// // actual run command - ./cometbft node --home ../data/.cometbft
		// binPath := filepath.Join(homePath, ".astria/local-dev-astria/cometbft")
		// // run the actual command

		// initCmd := exec.Command(initCmdPath)

		// envPath := filepath.Join(homePath, ".astria/local-dev-astria/.env")
		// environment := loadAndGetEnvVariables(envPath)

		// for _, env := range environment {
		// 	fmt.Println(env)
		// }
		// cmd = exec.Command(binPath)
		// cmd.Env = environment

		// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 	if event.Key() == tcell.KeyCtrlC {
		// 		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		// 			fmt.Println("Failed to send SIGINT to the process:", err)
		// 		}
		// 		app.Stop()
		// 		return nil
		// 	}
		// 	return event
		// })

		// stdout, err := cmd.StdoutPipe()
		// if err != nil {
		// 	panic(err)
		// }

		// if err := cmd.Start(); err != nil {
		// 	panic(err)
		// }

		// stdoutScanner := bufio.NewScanner(stdout)

		// for stdoutScanner.Scan() {
		// 	line := stdoutScanner.Text()
		// 	app.QueueUpdateDraw(func() {
		// 		sequencerTextView.Write([]byte(line + "\n"))
		// 	})
		// }

		// app.QueueUpdateDraw(func() {
		// 	sequencerTextView.Write([]byte("we finished" + "\n"))
		// })
		// if err := stdoutScanner.Err(); err != nil {
		// 	panic(err)
		// }

		// if err := cmd.Wait(); err != nil {
		// 	panic(err)
		// }
	}()

	// Create a new Flex layout.
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sequencerTextView, 0, 1, false).
		AddItem(cometbftTextView, 0, 1, false).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle(" 2 "), 0, 1, false).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle(" 3 "), 0, 1, false).
		AddItem(composerTextView, 0, 1, false).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle(" 4 "), 0, 1, false)
		AddItem(conductorTextView, 0, 1, false)
		// textview.SetBorder(true).SetTitle("Astria")
		// if err := box.SetRoot(textView, true).Run(); err != nil {
		// 	panic(err)
		// }

		// frame := tview.NewFrame(flex).
		// 	SetBorders(0, 0, 0, 0, 0, 0).
		// 	SetBorder(true).
		// 	SetTitle(" Astria Dev ")

		// if err := app.SetRoot(frame, true).SetFocus(frame).Run(); err != nil {
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
