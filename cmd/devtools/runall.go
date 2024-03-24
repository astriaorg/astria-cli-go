package devtools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	"github.com/spf13/cobra"
)

// runallCmd represents the runall command
var runallCmd = &cobra.Command{
	Use:   "runall",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runall()
	},
}

func init() {
	cmd.RootCmd.AddCommand(runallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runall() {
	ctx := cmd.RootCmd.Context()

	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	defaultDir := filepath.Join(homePath, ".astria")

	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(defaultDir, "local-dev-astria/.env")
	environment := loadAndGetEnvVariables(envPath)

	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Sequencer",
		BinPath: filepath.Join(homePath, ".astria/local-dev-astria/astria-sequencer"),
		Env:     environment,
		Args:    nil,
	}
	seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err = seqRunner.Start(shouldStart)
	if err != nil {
		fmt.Println("Error running sequencer:", err)
		panic(err)
	}

	// cometbft
	cometDataPath := filepath.Join(homePath, ".astria/data/.cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Comet BFT",
		BinPath: filepath.Join(homePath, ".astria/local-dev-astria/cometbft"),
		Env:     environment,
		Args:    []string{"node", "--home", cometDataPath},
	}
	cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)
	err = cometRunner.Start(seqRunner.GetDidStart())
	if err != nil {
		fmt.Println("Error running composer:", err)
		panic(err)
	}

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: filepath.Join(homePath, ".astria/local-dev-astria/astria-composer"),
		Env:     environment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)
	err = compRunner.Start(cometRunner.GetDidStart())
	if err != nil {
		fmt.Println("Error running composer:", err)
		panic(err)
	}

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: filepath.Join(homePath, ".astria/local-dev-astria/astria-conductor"),
		Env:     environment,
		Args:    nil,
	}
	condRunner := processrunner.NewProcessRunner(ctx, conductorOpts)
	err = condRunner.Start(compRunner.GetDidStart())
	if err != nil {
		fmt.Println("Error running conductor:", err)
		panic(err)
	}

	runners := []*processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}
