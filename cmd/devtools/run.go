package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

// loadEnvVariables loads the environment variables from the src file
func loadEnvVariables(src string) {
	err := godotenv.Load(src)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

// exists checks if a file or directory exists
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
	// all dirs/files that should exist
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
	// Create channels to properly control the start order of all processes
	sequencerStartComplete := make(chan bool)
	cometbftStartComplete := make(chan bool)
	composerStartComplete := make(chan bool)

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
	// check if the `dev init` command has been run
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

	// create text view object for the composer
	composerTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	composerTextView.SetTitle(" Composer ").SetBorder(true)

	// create text view object for the conductor
	conductorTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	conductorTextView.SetTitle(" Conductor ").SetBorder(true)

	mainWindowHelpInfo := tview.NewTextView().
		SetText(" Press 'Ctrl-C' or 'q' to exit | 'w' to toggle word wrap | 'tab' or 'up/down' arrows to select app focus | 'enter' to go fullscreen on selected app")

	fullscreenHelpInfo := tview.NewTextView().
		SetText(" Press 'Ctrl-C' or 'q' to exit | 'w' to toggle word wrap | 'esc' to exit fullscreen")

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(sequencerTextView, 0, 1, true).
			AddItem(cometbftTextView, 0, 1, false).
			AddItem(composerTextView, 0, 1, false).
			AddItem(conductorTextView, 0, 1, false), 0, 4, false).SetDirection(tview.FlexRow).
		AddItem(mainWindowHelpInfo, 1, 0, false)
	flex.SetTitle(" Astria Dev ").SetBorder(true)

	// Create ANSI writers for the text views
	aWriterSequencerTextView := tview.ANSIWriter(sequencerTextView)
	aWriterCometbftTextView := tview.ANSIWriter(cometbftTextView)
	aWriterComposerTextView := tview.ANSIWriter(composerTextView)
	aWriterConductorTextView := tview.ANSIWriter(conductorTextView)

	isFullscreen := false   // controlled by the 'enter' and 'esc' keys
	isAutoScrolling := true // controlled by the 's' key
	var focusedItem tview.Primitive = sequencerTextView

	// start the app with auto scrolling enabled
	sequencerTextView.ScrollToEnd()
	cometbftTextView.ScrollToEnd()
	composerTextView.ScrollToEnd()
	conductorTextView.ScrollToEnd()
	appendText := func(text string, writer io.Writer, textView *tview.TextView) {
		writer.Write([]byte(text + "\n"))
	}

	// Create a list of items to cycle through
	items := []tview.Primitive{sequencerTextView, cometbftTextView, composerTextView, conductorTextView}
	currentIndex := 0

	setFocus := func(index int) {
		currentIndex = index
		for i, item := range items {
			// Use a type assertion to convert the tview.Primitive back to *tview.TextView
			frame, ok := item.(*tview.TextView)
			if !ok {
				// The item is not a *tview.TextView, so skip it.
				continue
			}
			if i == index {
				title := frame.GetTitle()
				title = "[black:green]" + title + "[::-]"
				frame.SetBorderColor(tcell.ColorGreen).SetTitle(title)
			} else {
				title := frame.GetTitle()
				regexPattern := `\[.*?\]`
				re, err := regexp.Compile(regexPattern)
				if err != nil {
					fmt.Println("Error compiling regex:", err)
					return
				}
				title = re.ReplaceAllString(title, "")
				frame.SetBorderColor(tcell.ColorGray).SetTitle(title)
			}
		}

		app.SetFocus(items[index])
	}
	setFocus(currentIndex)

	// create the sequencer run command
	sequencerBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-sequencer")
	seqCmd := exec.Command(sequencerBinPath)
	seqCmd.Env = environment

	// create the cometbft run command
	cometbftDataPath := filepath.Join(homePath, ".astria/data/.cometbft")
	cometbftCmdPath := filepath.Join(homePath, ".astria/local-dev-astria/cometbft")
	nodeCmdArgs := []string{"node", "--home", cometbftDataPath}
	cometbftCmd := exec.Command(cometbftCmdPath, nodeCmdArgs...)
	cometbftCmd.Env = environment

	// create the composer run command
	composerBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-composer")
	composerCmd := exec.Command(composerBinPath)
	composerCmd.Env = environment

	// create the conductor run command
	conductorBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-conductor")
	conductorCmd := exec.Command(conductorBinPath)
	conductorCmd.Env = environment

	// Track the current word wrap status.
	wordWrapEnabled := true

	var fullscreenInputCapture, focusModeInputCapture func(event *tcell.EventKey) *tcell.EventKey

	// create the input capture for the app in fullscreen
	fullscreenInputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		// properly handle ctrl-c and pass SIGINT to the running processes
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
		// set 'q' to exit the app and pass SIGINT to the running processes
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
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
		// set 'w' to toggle word wrap
		if event.Key() == tcell.KeyRune && event.Rune() == 'w' {
			sequencerTextView.SetWrap(!wordWrapEnabled)
			cometbftTextView.SetWrap(!wordWrapEnabled)
			composerTextView.SetWrap(!wordWrapEnabled)
			conductorTextView.SetWrap(!wordWrapEnabled)
			wordWrapEnabled = !wordWrapEnabled
			return nil
		}
		// set 's' to toggle auto scrolling
		if event.Key() == tcell.KeyRune && event.Rune() == 's' {
			if !isAutoScrolling {
				sequencerTextView.ScrollToEnd()
				cometbftTextView.ScrollToEnd()
				composerTextView.ScrollToEnd()
				conductorTextView.ScrollToEnd()
			} else {
				currentOffset, _ := sequencerTextView.GetScrollOffset()
				sequencerTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = cometbftTextView.GetScrollOffset()
				cometbftTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = composerTextView.GetScrollOffset()
				composerTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = conductorTextView.GetScrollOffset()
				conductorTextView.ScrollTo(currentOffset, 0)
			}
			isAutoScrolling = !isAutoScrolling
		}
		// set 'tab' to cycle through the apps
		if event.Key() == tcell.KeyTab && !isFullscreen {
			newIndex := (currentIndex + 1) % len(items)
			setFocus(newIndex)
			return nil
		}
		// set 'enter' to go fullscreen on the selected app
		if event.Key() == tcell.KeyEnter && !isFullscreen {
			isFullscreen = true
			frame, ok := items[currentIndex].(*tview.TextView)
			if !ok {
				return event
			}
			fullscreenFlex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(frame, 0, 1, true).
				AddItem(fullscreenHelpInfo, 1, 0, false), 0, 4, false)
			app.SetRoot(fullscreenFlex, true)
			app.SetInputCapture(focusModeInputCapture)
			return nil
		}

		// set 'up' arrow to cycle through the apps
		if event.Key() == tcell.KeyUp && !isFullscreen {
			if currentIndex > 0 {
				newIndex := (currentIndex - 1) % len(items)
				setFocus(newIndex)
			} else {
				setFocus(0)
			}
			return nil
		}
		// set 'down' arrow to cycle through the apps
		if event.Key() == tcell.KeyDown && !isFullscreen {
			if currentIndex == len(items)-1 {
				setFocus(currentIndex)
			} else {
				newIndex := (currentIndex + 1) % len(items)
				setFocus(newIndex)
			}
			return nil

		}
		return event
	}

	// set the input capture for the app in app focus mode
	focusModeInputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		// get the focused item
		frame, ok := items[currentIndex].(*tview.TextView)
		if !ok {
			return event
		}
		// properly handle ctrl-c and pass SIGINT to the running processes
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
		// set 'q' to exit the app and pass SIGINT to the running processes
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
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
		// set 'w' to toggle word wrap
		if event.Key() == tcell.KeyRune && event.Rune() == 'w' {
			sequencerTextView.SetWrap(!wordWrapEnabled)
			cometbftTextView.SetWrap(!wordWrapEnabled)
			composerTextView.SetWrap(!wordWrapEnabled)
			conductorTextView.SetWrap(!wordWrapEnabled)
			wordWrapEnabled = !wordWrapEnabled
			return nil
		}
		// set 's' to toggle auto scrolling
		if event.Key() == tcell.KeyRune && event.Rune() == 's' {
			if !isAutoScrolling {
				sequencerTextView.ScrollToEnd()
				cometbftTextView.ScrollToEnd()
				composerTextView.ScrollToEnd()
				conductorTextView.ScrollToEnd()
			} else {
				currentOffset, _ := sequencerTextView.GetScrollOffset()
				sequencerTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = cometbftTextView.GetScrollOffset()
				cometbftTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = composerTextView.GetScrollOffset()
				composerTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = conductorTextView.GetScrollOffset()
				conductorTextView.ScrollTo(currentOffset, 0)
			}
			isAutoScrolling = !isAutoScrolling
		}
		// set 'esc' to exit fullscreen
		if event.Key() == tcell.KeyEscape && isFullscreen {
			isFullscreen = false
			frame, ok := items[currentIndex].(*tview.TextView)
			if !ok {
				return event
			}
			frame.SetInputCapture(nil)
			frame.SetMouseCapture(nil)
			app.SetRoot(flex, true)
			app.SetInputCapture(fullscreenInputCapture)
			return nil
		}

		switch event.Key() {
		case tcell.KeyUp:
			row, _ := frame.GetScrollOffset()
			frame.ScrollTo(row-1, 0)
		case tcell.KeyDown:
			row, _ := frame.GetScrollOffset()
			frame.ScrollTo(row+1, 0)
		}

		frame.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			switch action {
			case tview.MouseScrollUp:
				row, _ := frame.GetScrollOffset()
				frame.ScrollTo(row-1, 0)
			case tview.MouseScrollDown:
				row, _ := frame.GetScrollOffset()
				frame.ScrollTo(row+1, 0)
			}
			return action, event
		})

		return event
	}

	// set the input capture for the app
	app.SetInputCapture(fullscreenInputCapture)

	// go routine for running the sequencer
	go func() {
		// Get a pipe to the command's output.
		stdout, err := seqCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := seqCmd.Start(); err != nil {
			panic(err)
		}

		// let the cometbft go routine know that it can start
		sequencerStartComplete <- true

		// Create a scanner to read the output line by line.
		stdoutScanner := bufio.NewScanner(stdout)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterSequencerTextView, sequencerTextView)
			})
		}
		if err := stdoutScanner.Err(); err != nil {
			panic(err)
		}
		if err := seqCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	// go routine for running cometbft
	go func() {
		<-sequencerStartComplete
		initCmdArgs := []string{"init", "--home", cometbftDataPath}
		initCmd := exec.Command(cometbftCmdPath, initCmdArgs...)
		initCmd.Env = environment

		p := fmt.Sprintf("Running command `%v %v %v %v`\n", initCmd, initCmdArgs[0], initCmdArgs[1], initCmdArgs[2])
		app.QueueUpdateDraw(func() {
			appendText(p, aWriterCometbftTextView, cometbftTextView)
		})

		out, err := initCmd.CombinedOutput()
		if err != nil {
			p := fmt.Sprintf("Error executing command `%v`: %v\n", initCmd, err)
			app.QueueUpdateDraw(func() {
				appendText(p, aWriterCometbftTextView, cometbftTextView)

			})
			return
		}
		app.QueueUpdateDraw(func() {
			appendText(string(out), aWriterCometbftTextView, cometbftTextView)

		})

		// create the comand to replace the defualt genesis.json with the
		// configured one
		initGenesisJsonPath := filepath.Join(homePath, ".astria/local-dev-astria/genesis.json")
		endGenesisJsonPath := filepath.Join(homePath, ".astria/data/.cometbft/config/genesis.json")
		copyArgs := []string{initGenesisJsonPath, endGenesisJsonPath}
		copyCmd := exec.Command("cp", copyArgs...)
		copyCmd.Env = environment

		_, err = copyCmd.CombinedOutput()
		if err != nil {
			p := fmt.Sprintf("Error executing command `%v`: %v\n", copyCmd, err)
			app.QueueUpdateDraw(func() {
				appendText(p, aWriterCometbftTextView, cometbftTextView)

			})
			return
		}
		p = fmt.Sprintf("Copied genesis.json to %s\n", endGenesisJsonPath)
		app.QueueUpdateDraw(func() {
			appendText(p, aWriterCometbftTextView, cometbftTextView)

		})

		// create the comand to replace the defualt priv_validator_key.json with the
		// configured one
		initPrivValidatorJsonPath := filepath.Join(homePath, ".astria/local-dev-astria/priv_validator_key.json")
		endPrivValidatorJsonPath := filepath.Join(homePath, ".astria/data/.cometbft/config/priv_validator_key.json")
		copyArgs = []string{initPrivValidatorJsonPath, endPrivValidatorJsonPath}
		copyCmd = exec.Command("cp", copyArgs...)
		copyCmd.Env = environment

		_, err = copyCmd.CombinedOutput()
		if err != nil {
			p := fmt.Sprintf("Error executing command `%v`: %v\n", copyCmd, err)
			app.QueueUpdateDraw(func() {
				appendText(p, aWriterCometbftTextView, cometbftTextView)

			})
			return
		}
		p = fmt.Sprintf("Copied priv_validator_key.json to %s\n", endPrivValidatorJsonPath)
		app.QueueUpdateDraw(func() {
			appendText(p, aWriterCometbftTextView, cometbftTextView)

		})

		// update the cometbft config.toml file to have the proper block time
		cometbftConfigPath := filepath.Join(homePath, ".astria/data/.cometbft/config/config.toml")
		oldValue := `timeout_commit = "1s"`
		newValue := `timeout_commit = "2s"`

		if err := replaceInFile(cometbftConfigPath, oldValue, newValue); err != nil {
			p := fmt.Sprintf("Error updating the file: %v : %v", cometbftConfigPath, err)
			app.QueueUpdateDraw(func() {
				appendText(p, aWriterCometbftTextView, cometbftTextView)

			})
			return
		} else {
			p := fmt.Sprintf("Updated %v successfully", cometbftConfigPath)
			app.QueueUpdateDraw(func() {
				appendText(p, aWriterCometbftTextView, cometbftTextView)

			})
		}

		stdout, err := cometbftCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		if err := cometbftCmd.Start(); err != nil {
			panic(err)
		}

		// let the composer go routine know that it can start
		cometbftStartComplete <- true

		stdoutScanner := bufio.NewScanner(stdout)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterCometbftTextView, cometbftTextView)

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

		stdout, err := composerCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := composerCmd.Start(); err != nil {
			panic(err)
		}

		// let the conductor go routine know that it can start
		composerStartComplete <- true

		// Create a scanner to read the output line by line.
		stdoutScanner := bufio.NewScanner(stdout)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterComposerTextView, composerTextView)
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
		stdout, err := conductorCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := conductorCmd.Start(); err != nil {
			panic(err)
		}

		// Create a scanner to read the output line by line.
		stdoutScanner := bufio.NewScanner(stdout)

		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterConductorTextView, conductorTextView)

			})
		}
		if err := stdoutScanner.Err(); err != nil {
			panic(err)
		}
		if err := conductorCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	app.SetFocus(focusedItem)
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

// replaceInFile replaces oldValue with newValue in the file at filename.
// it is used here to update the block time in the cometbft config.toml file.
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
