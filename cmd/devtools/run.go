package devtools

import (
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Run all the Astria services locally.",
	Long:   `Run all the Astria services locally. This will start the sequencer, cometbft, composer, and conductor.`,
	PreRun: cmd.SetLogLevel,
	Run:    runRun,
}

func init() {
	devCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("instance", "i", DefaultInstanceName, "Used as directory name in ~/.astria to enable running separate instances of the sequencer stack.")
	runCmd.Flags().BoolVarP(&isRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runCmd.Flags().BoolVarP(&isRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
	runCmd.MarkFlagsMutuallyExclusive("local", "remote")
}

func runRun(c *cobra.Command, args []string) {
	ctx := c.Context()

	instance := c.Flag("instance").Value.String()
	IsInstanceNameValidOrPanic(instance)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	defaultDir := filepath.Join(homePath, ".astria")
	instanceDir := filepath.Join(defaultDir, instance)

	runOpts := &runOpts{
		ctx:         ctx,
		instanceDir: instanceDir,
	}

	var runners []processrunner.ProcessRunner
	switch {
	case !isRunLocal && !isRunRemote:
		log.Debug("No --local or --remote flag provided. Defaulting to --local.")
		isRunLocal = true
		runners = runLocal(runOpts)
	case isRunLocal:
		log.Debug("--local flag provided. Running local sequencer.")
		runners = runLocal(runOpts)
	case isRunRemote:
		log.Debug("--remote flag provided. Connecting to remote sequencer.")
		runners = runRemote(runOpts)
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

func runLocal(opts *runOpts) []processrunner.ProcessRunner {
	instanceDir := opts.instanceDir
	ctx := opts.ctx
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(opts.instanceDir, LocalConfigDirName, ".env")
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

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err := seqRunner.Start(ctx, shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running sequencer")
	}
	err = cometRunner.Start(ctx, seqRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running cometbft")
	}
	err = compRunner.Start(ctx, cometRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running composer")
	}
	err = condRunner.Start(ctx, compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
	}

	runners := []processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}
	return runners
}

func runRemote(opts *runOpts) []processrunner.ProcessRunner {
	ctx := opts.ctx
	instanceDir := opts.instanceDir
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

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err := compRunner.Start(ctx, shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running composer")
	}
	err = condRunner.Start(ctx, compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
	}

	runners := []processrunner.ProcessRunner{compRunner, condRunner}
	return runners
}
