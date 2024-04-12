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

<<<<<<< HEAD
var IsRunLocal bool
var IsRunRemote bool

=======
>>>>>>> f04e4d3 (add run-mono-repo command)
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
<<<<<<< HEAD
	runCmd.Flags().BoolVarP(&IsRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runCmd.Flags().BoolVarP(&IsRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
=======
	runCmd.Flags().BoolVarP(&isRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runCmd.Flags().BoolVarP(&isRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
	runCmd.Flags().BoolVar(&exportLogs, "export-logs", false, "Export logs to files.")
>>>>>>> f04e4d3 (add run-mono-repo command)
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

<<<<<<< HEAD
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
=======
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
>>>>>>> f04e4d3 (add run-mono-repo command)
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

<<<<<<< HEAD
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
=======
func runLocal(opts *runOpts) []processrunner.ProcessRunner {
	instanceDir := opts.instanceDir
	runTime := opts.appStartTime
	ctx := opts.ctx
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(opts.instanceDir, LocalConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	logsDir := filepath.Join(opts.instanceDir, LogsDirName)

	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Sequencer",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-sequencer"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-sequencer.log"),
		ExportLogs: exportLogs,
>>>>>>> f04e4d3 (add run-mono-repo command)
	}
	seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

	// cometbft
	cometDataPath := filepath.Join(instanceDir, DataDirName, ".cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
<<<<<<< HEAD
		Title:   "Comet BFT",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "cometbft"),
		Env:     environment,
		Args:    []string{"node", "--home", cometDataPath},
=======
		Title:      "Comet BFT",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "cometbft"),
		Env:        environment,
		Args:       []string{"node", "--home", cometDataPath},
		LogPath:    filepath.Join(logsDir, runTime+"-cometbft.log"),
		ExportLogs: exportLogs,
>>>>>>> f04e4d3 (add run-mono-repo command)
	}
	cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
<<<<<<< HEAD
		Title:   "Composer",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:     environment,
		Args:    nil,
=======
		Title:      "Composer",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-composer.log"),
		ExportLogs: exportLogs,
>>>>>>> f04e4d3 (add run-mono-repo command)
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
<<<<<<< HEAD
		Title:   "Conductor",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:     environment,
		Args:    nil,
=======
		Title:      "Conductor",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-conductor.log"),
		ExportLogs: exportLogs,
>>>>>>> f04e4d3 (add run-mono-repo command)
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

<<<<<<< HEAD
func runRemote(ctx context.Context, instanceDir string) []processrunner.ProcessRunner {
=======
func runRemote(opts *runOpts) []processrunner.ProcessRunner {
	ctx := opts.ctx
	instanceDir := opts.instanceDir
	runTime := opts.appStartTime
>>>>>>> f04e4d3 (add run-mono-repo command)
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(instanceDir, RemoteConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
<<<<<<< HEAD
		Title:   "Composer",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:     environment,
		Args:    nil,
=======
		Title:      "Composer",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-composer.log"),
		ExportLogs: exportLogs,
>>>>>>> f04e4d3 (add run-mono-repo command)
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
<<<<<<< HEAD
		Title:   "Conductor",
		BinPath: filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:     environment,
		Args:    nil,
=======
		Title:      "Conductor",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-conductor.log"),
		ExportLogs: exportLogs,
>>>>>>> f04e4d3 (add run-mono-repo command)
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
