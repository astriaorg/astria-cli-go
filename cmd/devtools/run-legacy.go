package devtools

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

// runLegacyCmd represents the run-legacy command
var runLegacyCmd = &cobra.Command{
	Use:   "run-legacy",
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
	sequencerTVRowCount := 0

	// create text view object for the cometbft
	cometbftTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	cometbftTextView.SetTitle(" Cometbft ").SetBorder(true)
	cometbftTVRowCount := 0

	// create text view object for the composer
	composerTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	composerTextView.SetTitle(" Composer ").SetBorder(true)
	composerTVRowCount := 0

	// create text view object for the conductor
	conductorTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	conductorTextView.SetTitle(" Conductor ").SetBorder(true)
	conductorTVRowCount := 0

	// app settings
	isFullscreen := false    // controlled by the 'enter' and 'esc' keys
	isAutoScrolling := true  // controlled by the 's' key
	wordWrapEnabled := true  // controlled by the 'w' key
	isBorderlessLog := false // controlled by the 'b' key
	var focusedItem tview.Primitive = sequencerTextView

	helpTextHelp := "(h)elp"
	helpTextQuit := "(q)uit"
	helpTextFocus := "(up/down) arrows to select app focus"
	helpTextEnterFullscreen := "(enter) to go fullscreen on focused app"
	helpTextExitFullscreen := "(esc) to exit fullscreen"
	helpTextWordWrap := "(w)ord wrap"
	helpTextAutoScroll := "(a)utoscroll"
	helpTextBorderless := "(b)oarderless"
	helpTextLogScroll := "if not auto scrolling: (up/down) arrows or mousewheel to scroll"
	helpTextHead := "(0) jump to head"
	helpTextTail := "(1) jump to tail"

	appendStatus := func(text string, status bool) string {
		output := ""
		output += text
		if status {
			output += " - [black:green]  on [-:-]"
		} else {
			output += " - [white:darkred] off [-:-]"
		}
		return output
	}

	buildMainWindowHelpInfo := func() string {
		output := " "
		output += helpTextHelp + " | "
		output += helpTextQuit + " | "
		output += helpTextFocus + " | "
		output += helpTextEnterFullscreen + " | "
		output += appendStatus(helpTextWordWrap, wordWrapEnabled) + " | "
		output += appendStatus(helpTextAutoScroll, isAutoScrolling)
		return output
	}

	buildFullscreenHelpInfo := func() string {
		output := " "
		output += helpTextHelp + " | "
		output += helpTextQuit + " | "
		output += helpTextExitFullscreen + " | "
		output += appendStatus(helpTextWordWrap, wordWrapEnabled) + " | "
		output += appendStatus(helpTextAutoScroll, isAutoScrolling) + " | "
		output += appendStatus(helpTextBorderless, isBorderlessLog) + " | "
		output += helpTextTail + " | "
		output += helpTextHead + " | "
		output += helpTextLogScroll
		return output
	}

	buildMainHelpScreenText := func() string {
		output := "Navigation:\t\n"
		output += "\t[:darkslategray]tab[:-]: Cycle the focus to the next app.\n"
		output += "\t[:darkslategray]up/down[:-] arrows: [yellow:][When in main window][-:] Change focus to the previous or next app.\n"
		output += "\t[:darkslategray]up/down[:-] arrows: [yellow:][When in focued window with autoscroll OFF][-:] Go up or down a line in the focued logs.\n"
		output += "\t[:darkslategray]mouse scroll[:-]: [yellow:][When in focued window with autoscroll OFF][-:] Scroll up or down in the focued logs.\n\n"
		output += "Focus Control:\n"
		output += "\t[:darkslategray]enter[:-]: Go from the main screen to fullscreen on the focused app's logs.\n"
		output += "\t[:darkslategray]esc[:-]:   Go from the fullscreened log view back to the main window.\n\n"
		output += "Word wrap:\n"
		output += "\t[:darkslategray]w[:-]: Toggle if word wrap is on or off.\n\n"
		output += "Logs Controls:\n"
		output += "\t[:darkslategray]a[:-]: Toggle if autoscroll is on or off.\n"
		output += "\t[:darkslategray]1[:-]: [yellow:][When in focued window][-:] Jump to tail of logs and disable autoscrolling.\n"
		output += "\t[:darkslategray]0[:-]: [yellow:][When in focued window][-:] Jump to head of logs and disable autoscrolling.\n\n"
		output += "Borderless:\n"
		output += "\t[:darkslategray]b[:-]: [yellow:][When in focued window][-:] Toggle the border around the logs on or off.\n\n"
		output += "Quitting:\n"
		output += "\t[:darkslategray]q[:-]:      Quit the app.\n"
		output += "\t[:darkslategray]ctrl-c[:-]: Quit the app.\n\n"
		output += "Help:\n"
		output += "\t[:darkslategray]h[:-]: Show this help screen or return to previous window.\n\n"
		return output
	}

	helpscreenHelpInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetText("(q)uit | (h) to return to previous window").
		SetChangedFunc(func() {
			app.Draw()
		})

	helpMainWindowInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetText(buildMainHelpScreenText()).
		SetWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	helpMainWindowInfo.SetTitle(" Astria CLI Help ").SetBorder(true).SetBorderColor(tcell.ColorBlue).SetBorderPadding(0, 0, 1, 0)
	helpScreenFlex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(helpMainWindowInfo, 0, 1, true).
		AddItem(helpscreenHelpInfo, 1, 0, false), 0, 4, false)

	mainWindowHelpInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetText(buildMainWindowHelpInfo()).
		SetChangedFunc(func() {
			app.Draw()
		})

	fullscreenHelpInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetText(buildFullscreenHelpInfo()).
		SetChangedFunc(func() {
			app.Draw()
		})

	mainWindow := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(sequencerTextView, 0, 1, true).
			AddItem(cometbftTextView, 0, 1, false).
			AddItem(composerTextView, 0, 1, false).
			AddItem(conductorTextView, 0, 1, false), 0, 4, false).SetDirection(tview.FlexRow).
		AddItem(mainWindowHelpInfo, 1, 0, false)
	mainWindow.SetTitle(" Astria Dev ").SetBorder(true)

	// prevWindow is used for toggling in and out of the help window
	var prevWindow tview.Primitive = mainWindow

	// Create ANSI writers for the text views
	aWriterSequencerTextView := tview.ANSIWriter(sequencerTextView)
	aWriterCometbftTextView := tview.ANSIWriter(cometbftTextView)
	aWriterComposerTextView := tview.ANSIWriter(composerTextView)
	aWriterConductorTextView := tview.ANSIWriter(conductorTextView)

	// start the app with auto scrolling enabled
	sequencerTextView.ScrollToEnd()
	cometbftTextView.ScrollToEnd()
	composerTextView.ScrollToEnd()
	conductorTextView.ScrollToEnd()
	appendText := func(text string, writer io.Writer) {
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

	var mainWindowInputCapture, focusWindowInputCapture, helpWindowInputCapture func(event *tcell.EventKey) *tcell.EventKey
	var prevInputCapture func(event *tcell.EventKey) *tcell.EventKey

	// create the input capture for the app in fullscreen
	mainWindowInputCapture = func(event *tcell.EventKey) *tcell.EventKey {
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
		if event.Key() == tcell.KeyRune && (event.Rune() == 'q' || event.Rune() == 'Q') {
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
		if event.Key() == tcell.KeyRune && (event.Rune() == 'w' || event.Rune() == 'W') {
			sequencerTextView.SetWrap(!wordWrapEnabled)
			cometbftTextView.SetWrap(!wordWrapEnabled)
			composerTextView.SetWrap(!wordWrapEnabled)
			conductorTextView.SetWrap(!wordWrapEnabled)
			wordWrapEnabled = !wordWrapEnabled

			mainWindowHelpInfo.SetText(buildMainWindowHelpInfo())
			fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())

			return nil
		}
		// set 'a' to toggle auto scrolling
		if event.Key() == tcell.KeyRune && (event.Rune() == 'a' || event.Rune() == 'A') {
			isAutoScrolling = !isAutoScrolling
			if isAutoScrolling {
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

			mainWindowHelpInfo.SetText(buildMainWindowHelpInfo())
			fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())

			return nil
		}

		// TODO: add a key to just to head of logs (automatically turn off auto scrolling)

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
			prevWindow = fullscreenFlex
			prevInputCapture = focusWindowInputCapture
			app.SetRoot(fullscreenFlex, true)
			app.SetInputCapture(nil)
			app.SetInputCapture(focusWindowInputCapture)
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

		if event.Key() == tcell.KeyRune && (event.Rune() == 'h' || event.Rune() == 'H') {
			prevWindow = mainWindow
			prevInputCapture = mainWindowInputCapture
			app.SetInputCapture(nil)
			app.SetInputCapture(helpWindowInputCapture)
			app.SetRoot(helpScreenFlex, true)
		}
		return event
	}
	prevInputCapture = mainWindowInputCapture

	// set the input capture for the app in app focus mode
	focusWindowInputCapture = func(event *tcell.EventKey) *tcell.EventKey {
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
		if event.Key() == tcell.KeyRune && (event.Rune() == 'q' || event.Rune() == 'Q') {
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
		if event.Key() == tcell.KeyRune && (event.Rune() == 'w' || event.Rune() == 'W') {
			sequencerTextView.SetWrap(!wordWrapEnabled)
			cometbftTextView.SetWrap(!wordWrapEnabled)
			composerTextView.SetWrap(!wordWrapEnabled)
			conductorTextView.SetWrap(!wordWrapEnabled)
			wordWrapEnabled = !wordWrapEnabled

			mainWindowHelpInfo.SetText(buildMainWindowHelpInfo())
			fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())

			return nil
		}
		// set 'a' to toggle auto scrolling
		if event.Key() == tcell.KeyRune && (event.Rune() == 'a' || event.Rune() == 'A') {
			isAutoScrolling = !isAutoScrolling
			if isAutoScrolling {
				sequencerTextView.ScrollToEnd()
				cometbftTextView.ScrollToEnd()
				composerTextView.ScrollToEnd()
				conductorTextView.ScrollToEnd()
			} else {
				// stop auto scrolling and allow the user to scroll manually
				currentOffset, _ := sequencerTextView.GetScrollOffset()
				sequencerTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = cometbftTextView.GetScrollOffset()
				cometbftTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = composerTextView.GetScrollOffset()
				composerTextView.ScrollTo(currentOffset, 0)
				currentOffset, _ = conductorTextView.GetScrollOffset()
				conductorTextView.ScrollTo(currentOffset, 0)
			}

			mainWindowHelpInfo.SetText(buildMainWindowHelpInfo())
			fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())

			return nil
		}

		// TODO: add a key to just to head of logs (automatically turn off auto scrolling)

		// set 'esc' to exit fullscreen
		if event.Key() == tcell.KeyEscape && isFullscreen {
			isFullscreen = false
			frame, ok := items[currentIndex].(*tview.TextView)
			if !ok {
				return event
			}
			// clear settings on the focused item
			frame.SetInputCapture(nil)
			frame.SetMouseCapture(nil)
			// reenable the border on the focused item so it shows up in the main window
			frame.SetBorder(true)
			isBorderlessLog = false
			fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())
			// set the app back to the main window and update the prev input and
			// window for working with the help window
			prevInputCapture = mainWindowInputCapture
			prevWindow = mainWindow
			app.SetRoot(mainWindow, true)
			app.SetInputCapture(mainWindowInputCapture)
			return nil
		}

		switch event.Key() {
		case tcell.KeyUp:
			if !isAutoScrolling {
				row, _ := frame.GetScrollOffset()
				frame.ScrollTo(row-1, 0)
				return nil
			}
		case tcell.KeyDown:
			if !isAutoScrolling {
				row, _ := frame.GetScrollOffset()
				frame.ScrollTo(row+1, 0)
				return nil
			}
		}
		frame.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			switch action {
			case tview.MouseScrollUp:
				if !isAutoScrolling {
					row, _ := frame.GetScrollOffset()
					frame.ScrollTo(row-1, 0)
				}
				return action, event
			case tview.MouseScrollDown:
				if !isAutoScrolling {
					row, _ := frame.GetScrollOffset()
					frame.ScrollTo(row+1, 0)
				}
				return action, event
			}
			return action, event
		})
		// jump to help screen
		if event.Key() == tcell.KeyRune && (event.Rune() == 'h' || event.Rune() == 'H') {
			app.SetInputCapture(nil)
			app.SetInputCapture(helpWindowInputCapture)
			app.SetRoot(helpScreenFlex, true)
		}
		// toggle the border on the longs with 'b'
		if event.Key() == tcell.KeyRune && (event.Rune() == 'b' || event.Rune() == 'B') {
			isBorderlessLog = !isBorderlessLog
			if isBorderlessLog {
				frame, ok := items[currentIndex].(*tview.TextView)
				if !ok {
					return event
				}
				// TODO: make the set focus stuff below into a function
				fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())
				frame.SetBorder(false)
				fullscreenFlex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(frame, 0, 1, true).
					AddItem(fullscreenHelpInfo, 1, 0, false), 0, 4, false)
				prevWindow = fullscreenFlex
				prevInputCapture = focusWindowInputCapture
				app.SetRoot(fullscreenFlex, true)
			} else {
				frame, ok := items[currentIndex].(*tview.TextView)
				if !ok {
					return event
				}
				// TODO: make the set focus stuff below into a function
				fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())
				frame.SetBorder(true)
				fullscreenFlex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(frame, 0, 1, true).
					AddItem(fullscreenHelpInfo, 1, 0, false), 0, 4, false)
				prevWindow = fullscreenFlex
				prevInputCapture = focusWindowInputCapture
				app.SetRoot(fullscreenFlex, true)
			}
			return nil
		}
		// using '0' for head, 'h' already in use for help
		if event.Key() == tcell.KeyRune && (event.Rune() == '0' || event.Rune() == ')') {
			// disable auto scrolling and jump to the head of the logs
			isAutoScrolling = false
			sequencerTextView.ScrollToBeginning()
			cometbftTextView.ScrollToBeginning()
			composerTextView.ScrollToBeginning()
			conductorTextView.ScrollToBeginning()

			mainWindowHelpInfo.SetText(buildMainWindowHelpInfo())
			fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())
			return nil
		}
		if event.Key() == tcell.KeyRune && (event.Rune() == '1' || event.Rune() == '!') {
			// disable auto scrolling and jump to the tail of the logs
			// stop auto scrolling and allow the user to scroll manually
			isAutoScrolling = false

			sequencerTextView.ScrollTo(sequencerTVRowCount, 0)
			cometbftTextView.ScrollTo(cometbftTVRowCount, 0)
			composerTextView.ScrollTo(composerTVRowCount, 0)
			conductorTextView.ScrollTo(conductorTVRowCount, 0)

			mainWindowHelpInfo.SetText(buildMainWindowHelpInfo())
			fullscreenHelpInfo.SetText(buildFullscreenHelpInfo())
			return nil
		}
		return event
	}

	helpWindowInputCapture = func(event *tcell.EventKey) *tcell.EventKey {
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
		if event.Key() == tcell.KeyRune && (event.Rune() == 'q' || event.Rune() == 'Q') {
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
		if event.Key() == tcell.KeyRune && (event.Rune() == 'h' || event.Rune() == 'H') {
			app.SetInputCapture(nil)
			app.SetInputCapture(prevInputCapture)
			app.SetRoot(prevWindow, true)
		}
		return event
	}

	// set the input capture for the app
	app.SetInputCapture(mainWindowInputCapture)

	// go routine for running the sequencer
	go func() {
		// Get a pipe to the command's output.
		stdout, err := seqCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		stderr, err := seqCmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		// Create a scanner to read the output line by line.
		output := io.MultiReader(stdout, stderr)
		outputScanner := bufio.NewScanner(output)

		if err := seqCmd.Start(); err != nil {
			panic(err)
		}

		// let the cometbft go routine know that it can start
		sequencerStartComplete <- true

		for outputScanner.Scan() {
			line := outputScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterSequencerTextView)
				sequencerTVRowCount++
			})
		}
		if err := outputScanner.Err(); err != nil {
			panic(err)
		}
		if err := seqCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	// go routine for running cometbft
	go func() {
		<-sequencerStartComplete

		stdout, err := cometbftCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		stderr, err := cometbftCmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		// Create a scanner to read the output line by line.
		output := io.MultiReader(stdout, stderr)
		outputScanner := bufio.NewScanner(output)

		if err := cometbftCmd.Start(); err != nil {
			panic(err)
		}

		// let the composer go routine know that it can start
		cometbftStartComplete <- true

		for outputScanner.Scan() {
			line := outputScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterCometbftTextView)
				cometbftTVRowCount++

			})
		}
		if err := outputScanner.Err(); err != nil {
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
		stderr, err := composerCmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		// Create a scanner to read the output line by line.
		output := io.MultiReader(stdout, stderr)
		outputScanner := bufio.NewScanner(output)

		if err := composerCmd.Start(); err != nil {
			panic(err)
		}

		// let the conductor go routine know that it can start
		composerStartComplete <- true

		for outputScanner.Scan() {
			line := outputScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterComposerTextView)
				composerTVRowCount++
			})
		}
		if err := outputScanner.Err(); err != nil {
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
		stderr, err := conductorCmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		// Create a scanner to read the output line by line.
		output := io.MultiReader(stdout, stderr)
		outputScanner := bufio.NewScanner(output)

		if err := conductorCmd.Start(); err != nil {
			panic(err)
		}

		for outputScanner.Scan() {
			line := outputScanner.Text()
			app.QueueUpdateDraw(func() {
				appendText(line, aWriterConductorTextView)
				conductorTVRowCount++
			})
		}
		if err := outputScanner.Err(); err != nil {
			panic(err)
		}
		if err := conductorCmd.Wait(); err != nil {
			panic(err)
		}
	}()

	prevWindow = mainWindow
	prevInputCapture = mainWindowInputCapture
	app.SetFocus(focusedItem)
	if err := app.SetRoot(mainWindow, true).Run(); err != nil {
		panic(err)
	}
}

func init() {
	devCmd.AddCommand(runLegacyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runLegacyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runLegacyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
