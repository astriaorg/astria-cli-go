package devtools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var IsRunLocal bool
var IsRunRemote bool

type RunConfiguration int

const (
	Local RunConfiguration = iota
	Remote
	Error
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all the Astria services locally.",
	Long:  `Run all the Astria services locally. This will start the sequencer, cometbft, composer, and conductor.`,
	Run:   runall,
}

func init() {
	devCmd.AddCommand(runCmd)
	instanceFlagUsage := fmt.Sprintf("Choose where the local-dev-astria directory will be created. Defaults to \"%s\" if not provided.", DefaultInstanceName)
	runCmd.Flags().StringP("instance", "i", DefaultInstanceName, instanceFlagUsage)
	runCmd.Flags().BoolVarP(&IsRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runCmd.Flags().BoolVarP(&IsRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
}

func runall(c *cobra.Command, args []string) {
	ctx := cmd.RootCmd.Context()

	instance := c.Flag("instance").Value.String()
	err := IsInstanceNameValid(instance)
	if err != nil {
		log.WithError(err).Error("Error getting --instance flag")
		return
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	defaultDir := filepath.Join(homePath, ".astria")
	instanceDir := filepath.Join(defaultDir, instance)

	var runners []processrunner.ProcessRunner
	runConfiguration := determineRunConfiguration()
	switch runConfiguration {
	case Local:
		runners = runLocal(ctx, instanceDir)
	case Remote:
		runners = runRemote(ctx, instanceDir)
	case Error:
		log.Error("Error determining run configuration")
		return
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

func determineRunConfiguration() RunConfiguration {
	switch {
	case IsRunLocal && IsRunRemote:
		fmt.Println("Can only run one of --local or --remote, not both. Exiting.")
		return Error
	case !IsRunLocal && !IsRunRemote:
		fmt.Println("No --local or --remote flag provided. Defaulting to --local.")
		IsRunLocal = true
		return Local
	case IsRunLocal:
		fmt.Println("--local flag provided. Running local sequencer.")
		return Local
	case IsRunRemote:
		fmt.Println("--remote flag provided. Connecting to remote sequencer.")
		return Remote
	}
	return Error
}

func runLocal(ctx context.Context, instanceDir string) []processrunner.ProcessRunner {
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(instanceDir, LocalConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Sequencer",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-sequencer"),
		Env:     environment,
		Args:    nil,
	}
	seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

	// cometbft
	cometDataPath := filepath.Join(instanceDir, DataDirName, ".cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Comet BFT",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "cometbft"),
		Env:     environment,
		Args:    []string{"node", "--home", cometDataPath},
	}
	cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:     environment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:     environment,
		Args:    nil,
	}
	condRunner := processrunner.NewProcessRunner(ctx, conductorOpts)

	// cleanup function to stop processes if there is an error starting another process
	// FIXME - this isn't good enough. need to use context to stop all processes.
	cleanup := func() {
		seqRunner.Stop()
		cometRunner.Stop()
		compRunner.Stop()
		condRunner.Stop()
	}

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err := seqRunner.Start(shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running sequencer")
		cleanup()
		panic(err)
	}
	err = cometRunner.Start(seqRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running cometbft")
		cleanup()
		panic(err)
	}
	err = compRunner.Start(cometRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running composer")
		cleanup()
		panic(err)
	}
	err = condRunner.Start(compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
		cleanup()
		panic(err)
	}

	runners := []processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}
	return runners
}

func runRemote(ctx context.Context, instanceDir string) []processrunner.ProcessRunner {
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(instanceDir, RemoteConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:     environment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:     environment,
		Args:    nil,
	}
	condRunner := processrunner.NewProcessRunner(ctx, conductorOpts)

	// cleanup function to stop processes if there is an error starting another process
	// FIXME - this isn't good enough. need to use context to stop all processes.
	cleanup := func() {
		compRunner.Stop()
		condRunner.Stop()
	}

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err := compRunner.Start(shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running composer")
		cleanup()
		panic(err)
	}
	err = condRunner.Start(compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
		cleanup()
		panic(err)
	}

	runners := []processrunner.ProcessRunner{compRunner, condRunner}
	return runners
}
