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
var runMonoRepoCmd = &cobra.Command{
	Use:    "run-mono-repo [mono-repo-path]",
	Short:  "Run all the Astria services locally using locally built binaries.",
	Long:   `Run all the Astria services locally using the binaries built in the Astria mono-repo. This will start the sequencer, composer, and conductor from the mono-repo, but still use cometbft binary downloaded by the cli.`,
	Args:   cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	PreRun: cmd.SetLogLevel,
	Run:    runMonoRepoRun,
}

func init() {
	devCmd.AddCommand(runMonoRepoCmd)
	runMonoRepoCmd.Flags().StringP("instance", "i", DefaultInstanceName, "Used as directory name in ~/.astria to enable running separate instances of the sequencer stack.")
	runMonoRepoCmd.Flags().BoolVarP(&isRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runMonoRepoCmd.Flags().BoolVarP(&isRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
	runMonoRepoCmd.MarkFlagsMutuallyExclusive("local", "remote")
}

func runMonoRepoRun(c *cobra.Command, args []string) {
	monoRepoPath := args[0]

	ctx := c.Context()

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	defaultDir := filepath.Join(homePath, ".astria")

	instance := c.Flag("instance").Value.String()
	IsInstanceNameValidOrPanic(instance)
	instanceDir := filepath.Join(defaultDir, instance)

	runOpts := &runOpts{
		ctx:         ctx,
		instanceDir: instanceDir,
		// TODO - add validation for the mono-repo path
		monoRepoPath: monoRepoPath,
	}

	var runners []processrunner.ProcessRunner
	switch {
	case !isRunLocal && !isRunRemote:
		log.Debug("No --local or --remote flag provided. Defaulting to --local.")
		isRunLocal = true
		runners = runLocalUsingMonoRepo(runOpts)
	case isRunLocal:
		log.Debug("--local flag provided. Running local sequencer.")
		runners = runLocalUsingMonoRepo(runOpts)
	case isRunRemote:
		log.Debug("--remote flag provided. Connecting to remote sequencer.")
		runners = runRemoteUsingMonoRepo(runOpts)
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

func runLocalUsingMonoRepo(opts *runOpts) []processrunner.ProcessRunner {
	instanceDir := opts.instanceDir
	ctx := opts.ctx
	monoRepoPath := opts.monoRepoPath
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	defaultEnvPath := filepath.Join(opts.instanceDir, LocalConfigDirName, ".env")
	log.Debug("defaultEnvPath:", defaultEnvPath)
	defaultEnvironment := loadAndGetEnvVariables(defaultEnvPath)

	// load the .env file from the mono-repo
	// sequencerEnvPath := filepath.Join(monoRepoPath, "crates", "astria-sequencer", "local.env.example")
	// log.Debug("sequencerEnvPath:", sequencerEnvPath)
	// sequencerEnvironment := loadAndGetEnvVariables(sequencerEnvPath)
	// log.Debug("sequencerEnvironment:", sequencerEnvironment)
	// TODO - set the db path for sequencer to use the instance data dir
	// conductorEnvPath := filepath.Join(monoRepoPath, "crates", "astria-conductor", "local.env.example")
	// log.Debug("conductorEnvPath:", conductorEnvPath)
	// conductorEnvironment := loadAndGetEnvVariables(conductorEnvPath)
	// composerEnvPath := filepath.Join(monoRepoPath, "crates", "astria-composer", "local.env.example")
	// log.Debug("composerEnvPath:", composerEnvPath)
	// composerEnvironment := loadAndGetEnvVariables(composerEnvPath)

	// create the binaries paths for the services within the mono-repo
	sequencerBinPath := filepath.Join(monoRepoPath, AstriaTargetDebugPath, "astria-sequencer")
	log.Debug("sequencerBinPath:", sequencerBinPath)

	composerBinPath := filepath.Join(monoRepoPath, AstriaTargetDebugPath, "astria-composer")
	// composerBinPath := filepath.Join(instanceDir, BinariesDirName, "astria-composer")
	log.Debug("composerBinPath:", composerBinPath)

	conductorBinPath := filepath.Join(monoRepoPath, AstriaTargetDebugPath, "astria-conductor")
	// conductorBinPath := filepath.Join(instanceDir, BinariesDirName, "astria-conductor")
	log.Debug("conductorBinPath:", conductorBinPath)

	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Sequencer",
		BinPath: sequencerBinPath,
		Env:     defaultEnvironment,
		Args:    nil,
	}
	seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

	// cometbft
	cometDataPath := filepath.Join(instanceDir, DataDirName, ".cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Comet BFT",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "cometbft"),
		Env:     defaultEnvironment,
		Args:    []string{"node", "--home", cometDataPath},
	}
	cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: composerBinPath,
		Env:     defaultEnvironment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: conductorBinPath,
		Env:     defaultEnvironment,
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

func runRemoteUsingMonoRepo(opts *runOpts) []processrunner.ProcessRunner {
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
