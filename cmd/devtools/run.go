package devtools

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var IsRunLocal bool
var IsRunRemote bool
var ExportLogs bool

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
	runCmd.Flags().BoolVarP(&IsRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runCmd.Flags().BoolVarP(&IsRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
	runCmd.Flags().BoolVar(&ExportLogs, "export-logs", false, "Export logs to files.")
	runCmd.MarkFlagsMutuallyExclusive("local", "remote")
}

func runRun(c *cobra.Command, args []string) {
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
	if ExportLogs {
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

	var runners []processrunner.ProcessRunner
	switch {
	case !IsRunLocal && !IsRunRemote:
		log.Debug("No --local or --remote flag provided. Defaulting to --local.")
		IsRunLocal = true
		runners = runLocal(ctx, instanceDir, appStartTime)
	case IsRunLocal:
		log.Debug("--local flag provided. Running local sequencer.")
		runners = runLocal(ctx, instanceDir, appStartTime)
	case IsRunRemote:
		log.Debug("--remote flag provided. Connecting to remote sequencer.")
		runners = runRemote(ctx, instanceDir, appStartTime)
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

func runLocal(ctx context.Context, instanceDir string, runTime string) []processrunner.ProcessRunner {
	// load the .env file and get the environment variables
	// TODO - move config to own package w/ structs w/ defaults. still use .env for overrides.
	envPath := filepath.Join(instanceDir, LocalConfigDirName, ".env")
	environment := loadAndGetEnvVariables(envPath)

	logsDir := filepath.Join(instanceDir, LogsDirName)

	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Sequencer",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-sequencer"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-sequencer.log"),
		ExportLogs: ExportLogs,
	}
	seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

	// cometbft
	cometDataPath := filepath.Join(instanceDir, DataDirName, ".cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Comet BFT",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "cometbft"),
		Env:        environment,
		Args:       []string{"node", "--home", cometDataPath},
		LogPath:    filepath.Join(logsDir, runTime+"-cometbft.log"),
		ExportLogs: ExportLogs,
	}
	cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Composer",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-composer"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-composer.log"),
		ExportLogs: ExportLogs,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Conductor",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-conductor.log"),
		ExportLogs: ExportLogs,
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

func runRemote(ctx context.Context, instanceDir string, runTime string) []processrunner.ProcessRunner {
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
		ExportLogs: ExportLogs,
	}
	compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:      "Conductor",
		BinPath:    filepath.Join(instanceDir, BinariesDirName, "astria-conductor"),
		Env:        environment,
		Args:       nil,
		LogPath:    filepath.Join(logsDir, runTime+"-astria-conductor.log"),
		ExportLogs: ExportLogs,
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
