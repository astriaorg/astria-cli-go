package devtools

import (
	"context"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var IsRunLocal bool
var IsRunRemote bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all the Astria services locally.",
	Long:  `Run all the Astria services locally. This will start the sequencer, cometbft, composer, and conductor.`,
	Run:   runRun,
}

func init() {
	devCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("instance", "i", DefaultInstanceName, "Used as directory name in ~/.astria to enable running separate instances of the sequencer stack.")
	runCmd.Flags().BoolVarP(&IsRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runCmd.Flags().BoolVarP(&IsRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
	runCmd.MarkFlagsMutuallyExclusive("local", "remote")
}

func runRun(c *cobra.Command, args []string) {
	ctx := c.Context()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	instance := c.Flag("instance").Value.String()
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	defaultDir := filepath.Join(homePath, ".astria")
	instanceDir := filepath.Join(defaultDir, instance)

	var runners []processrunner.ProcessRunner
	switch {
	case !IsRunLocal && !IsRunRemote:
		log.Info("No --local or --remote flag provided. Defaulting to --local.")
		IsRunLocal = true
		runners = runLocal(ctx, instanceDir)
	case IsRunLocal:
		log.Info("--local flag provided. Running local sequencer.")
		runners = runLocal(ctx, instanceDir)
	case IsRunRemote:
		log.Info("--remote flag provided. Connecting to remote sequencer.")
		runners = runRemote(ctx, instanceDir)
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

func runLocal(ctx context.Context, instanceDir string) []processrunner.ProcessRunner {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(instanceDir, LocalConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Sequencer",
		BinPath: filepath.Join(instanceDir, LocalBinariesDirName, "astria-sequencer"),
		Env:     environment,
		Args:    nil,
	}
	seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

	// cometbft
	cometDataPath := filepath.Join(instanceDir, DataDirName, ".cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Comet BFT",
		BinPath: filepath.Join(instanceDir, LocalBinariesDirName, "cometbft"),
		Env:     environment,
		Args:    []string{"node", "--home", cometDataPath},
	}
	cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: filepath.Join(instanceDir, LocalBinariesDirName, "astria-composer"),
		Env:     environment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: filepath.Join(instanceDir, LocalBinariesDirName, "astria-conductor"),
		Env:     environment,
		Args:    nil,
	}
	condRunner := processrunner.NewProcessRunner(ctx, conductorOpts)

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err := seqRunner.Start(ctx, shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running sequencer")
		cancel()
	}
	err = cometRunner.Start(ctx, seqRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running cometbft")
		cancel()
	}
	err = compRunner.Start(ctx, cometRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running composer")
		cancel()
	}
	err = condRunner.Start(ctx, compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
		cancel()
	}

	runners := []processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}
	return runners
}

func runRemote(ctx context.Context, instanceDir string) []processrunner.ProcessRunner {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(instanceDir, RemoteConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: filepath.Join(instanceDir, RemoteBinariesDirName, "astria-composer"),
		Env:     environment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: filepath.Join(instanceDir, RemoteBinariesDirName, "astria-conductor"),
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
	err := compRunner.Start(ctx, shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running composer")
		cleanup()
		panic(err)
	}
	err = condRunner.Start(ctx, compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
		cleanup()
		panic(err)
	}

	runners := []processrunner.ProcessRunner{compRunner, condRunner}
	return runners
}
