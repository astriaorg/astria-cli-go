package devtools

import (
	"os"
	"path/filepath"
	"time"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runMonoRepoCmd = &cobra.Command{
	Use:    "run-locally-built [mono-repo-path]",
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
	runMonoRepoCmd.Flags().BoolVar(&exportLogs, "export-logs", false, "Export logs to files.")
	runMonoRepoCmd.MarkFlagsMutuallyExclusive("local", "remote")
}

func runMonoRepoRun(c *cobra.Command, args []string) {
	monoRepoPath := args[0]

	currentTime := time.Now()
	appStartTime := currentTime.Format("20060102-150405") // YYYYMMDD-HHMMSS

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

	logsDir := filepath.Join(instanceDir, LogsDirName)

	// conditionally create a log file for the app
	var appLogFile *os.File
	if exportLogs {
		logPath := filepath.Join(logsDir, appStartTime+"-app.log")
		appLogFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.Info("New log file created successfully:", logPath)
		log.SetOutput(appLogFile)
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true, // Disable ANSI color codes
			FullTimestamp: true,
		})

	} else {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
		appLogFile = nil
	}

	runOpts := &runOpts{
		ctx:          ctx,
		instanceDir:  instanceDir,
		appStartTime: appStartTime,
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

	// close the log file
	if appLogFile != nil {
		err := appLogFile.Close()
		if err != nil {
			log.WithError(err).Error("Error closing app log file")
		}
	}
}

func runLocalUsingMonoRepo(opts *runOpts) []processrunner.ProcessRunner {
	instanceDir := opts.instanceDir
	runTime := opts.appStartTime
	ctx := opts.ctx
	monoRepoPath := opts.monoRepoPath
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	defaultEnvPath := filepath.Join(opts.instanceDir, LocalConfigDirName, ".env")
	log.Debug("defaultEnvPath:", defaultEnvPath)
	defaultEnvironment := loadAndGetEnvVariables(defaultEnvPath)

	logsDir := filepath.Join(opts.instanceDir, LogsDirName)

	// load the .env file from the mono-repo
	sequencerEnvPath := filepath.Join(monoRepoPath, "crates", "astria-sequencer", "local.env.example")
	log.Debug("sequencerEnvPath:", sequencerEnvPath)
	sequencerEnvironment := loadAndGetEnvVariables(sequencerEnvPath, defaultEnvPath)
	// TODO - set the db path for sequencer to use the instance data dir
	conductorEnvPath := filepath.Join(monoRepoPath, "crates", "astria-conductor", "local.env.example")
	log.Debug("conductorEnvPath:", conductorEnvPath)
	conductorEnvironment := loadAndGetEnvVariables(conductorEnvPath, defaultEnvPath)
	composerEnvPath := filepath.Join(monoRepoPath, "crates", "astria-composer", "local.env.example")
	log.Debug("composerEnvPath:", composerEnvPath)
	composerEnvironment := loadAndGetEnvVariables(composerEnvPath, defaultEnvPath)

	// create the binaries paths for the services within the mono-repo
	sequencerBinPath := filepath.Join(monoRepoPath, AstriaTargetDebugPath, "astria-sequencer")
	// composerBinPath := filepath.Join(monoRepoPath, AstriaTargetDebugPath, "astria-composer")
	composerBinPath := filepath.Join(instanceDir, BinariesDirName, "astria-composer")
	// conductorBinPath := filepath.Join(monoRepoPath, AstriaTargetDebugPath, "astria-conductor")
	conductorBinPath := filepath.Join(instanceDir, BinariesDirName, "astria-conductor")

	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Sequencer",
		BinPath:    sequencerBinPath,
		Env:        sequencerEnvironment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-sequencer.log"),
		ExportLogs: exportLogs,
	}
	seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

	// cometbft
	cometDataPath := filepath.Join(instanceDir, DataDirName, ".cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Comet BFT",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "cometbft"),
		Env:        defaultEnvironment,
		Args:       []string{"node", "--home", cometDataPath},
		LogPath:    filepath.Join(logsDir, runTime+"-cometbft.log"),
		ExportLogs: exportLogs,
	}
	cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Composer",
		BinPath:    composerBinPath,
		Env:        composerEnvironment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-composer.log"),
		ExportLogs: exportLogs,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Conductor",
		BinPath:    conductorBinPath,
		Env:        conductorEnvironment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-conductor.log"),
		ExportLogs: exportLogs,
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
	runTime := opts.appStartTime
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(instanceDir, RemoteConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	logsDir := filepath.Join(instanceDir, LogsDirName)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Composer",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-composer.log"),
		ExportLogs: exportLogs,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Conductor",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-conductor.log"),
		ExportLogs: exportLogs,
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
