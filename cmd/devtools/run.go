package devtools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// boolean flags
var (
	isRunLocal  bool
	isRunRemote bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Run all the Astria services locally.",
	Long:   `Run all the Astria services locally. This will start the sequencer, cometbft, composer, and conductor.`,
	PreRun: cmd.SetLogLevel,
	Run:    runCmdHandler,
}

func init() {
	devCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("instance", "i", DefaultInstanceName, "Used as directory name in ~/.astria to enable running separate instances of the sequencer stack.")
	runCmd.Flags().BoolVarP(&isRunLocal, "local", "l", false, "Run the Astria stack using a locally running sequencer.")
	runCmd.Flags().BoolVarP(&isRunRemote, "remote", "r", false, "Run the Astria stack using a remote sequencer.")
	runCmd.MarkFlagsMutuallyExclusive("local", "remote")

	runCmd.Flags().String("environment-path", "", "Provide an override path to a specific environment file.")
	runCmd.Flags().String("conductor-path", "", "Provide an override path to a specific conductor binary.")
	runCmd.Flags().String("cometbft-path", "", "Provide an override path to a specific cometbft binary.")
	runCmd.Flags().String("composer-path", "", "Provide an override path to a specific composer binary.")
	runCmd.Flags().String("sequencer-path", "", "Provide an override path to a specific sequencer binary.")
}

func runCmdHandler(c *cobra.Command, args []string) {
	ctx := c.Context()

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}

	// astriaDir is the directory where all the astria instances data is stored
	astriaDir := filepath.Join(homePath, ".astria")

	// get instance name and check if it's valid
	instance := c.Flag("instance").Value.String()
	IsInstanceNameValidOrPanic(instance)

	// we will set runners after we decide which binaries we need to run
	var runners []processrunner.ProcessRunner

	// check if running local or remote sequencer.
	isLocalSequencer := isLocalSequencer()
	if isLocalSequencer {
		log.Debug("Running local sequencer")
		confDir := filepath.Join(astriaDir, instance, LocalConfigDirName)
		dataDir := filepath.Join(astriaDir, instance, DataDirName)
		binDir := filepath.Join(astriaDir, instance, BinariesDirName)
		// env path
		envPath := getFlagPathOrPanic(c, "environment-path", filepath.Join(confDir, ".env"))
		env := LoadEnvironment(envPath)

		// get the binary paths
		conductorPath := getFlagPathOrPanic(c, "conductor-path", filepath.Join(binDir, "astria-conductor"))
		cometbftPath := getFlagPathOrPanic(c, "cometbft-path", filepath.Join(binDir, "cometbft"))
		composerPath := getFlagPathOrPanic(c, "composer-path", filepath.Join(binDir, "astria-composer"))
		sequencerPath := getFlagPathOrPanic(c, "sequencer-path", filepath.Join(binDir, "astria-sequencer"))
		log.Debugf("Using binaries from %s", binDir)

		// sequencer
		seqOpts := processrunner.NewProcessRunnerOpts{
			Title:   "Sequencer",
			BinPath: sequencerPath,
			Env:     env,
			Args:    nil,
		}
		seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

		// cometbft
		cometDataPath := filepath.Join(dataDir, ".cometbft")
		cometOpts := processrunner.NewProcessRunnerOpts{
			Title:   "Comet BFT",
			BinPath: cometbftPath,
			Env:     env,
			Args:    []string{"node", "--home", cometDataPath},
		}
		cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

		// composer
		composerOpts := processrunner.NewProcessRunnerOpts{
			Title:   "Composer",
			BinPath: composerPath,
			Env:     env,
			Args:    nil,
		}
		compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

		// conductor
		conductorOpts := processrunner.NewProcessRunnerOpts{
			Title:   "Conductor",
			BinPath: conductorPath,
			Env:     env,
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

		runners = []processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}
	} else {
		log.Debug("Running remote sequencer")
		confDir := filepath.Join(astriaDir, instance, RemoteConfigDirName)
		binDir := filepath.Join(astriaDir, instance, BinariesDirName)
		// env path
		envPath := getFlagPathOrPanic(c, "environment-path", filepath.Join(confDir, ".env"))
		env := LoadEnvironment(envPath)

		// get the binary paths
		conductorPath := getFlagPathOrPanic(c, "conductor-path", filepath.Join(binDir, "astria-conductor"))
		composerPath := getFlagPathOrPanic(c, "composer-path", filepath.Join(binDir, "astria-composer"))

		// composer
		composerOpts := processrunner.NewProcessRunnerOpts{
			Title:   "Composer",
			BinPath: composerPath,
			Env:     env,
			Args:    nil,
		}
		compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

		// conductor
		conductorOpts := processrunner.NewProcessRunnerOpts{
			Title:   "Conductor",
			BinPath: conductorPath,
			Env:     env,
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
		runners = []processrunner.ProcessRunner{compRunner, condRunner}
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

// isLocalSequencer returns true if we should run the local sequencer
func isLocalSequencer() bool {
	switch {
	case !isRunLocal && !isRunRemote:
		log.Debug("No --local or --remote flag provided. Defaulting to --local.")
		return true
	case isRunLocal:
		log.Debug("--local flag provided. Running local sequencer.")
		return true
	case isRunRemote:
		log.Debug("--remote flag provided. Connecting to remote sequencer.")
		return false
	default:
		// this should never happen
		log.Debug("Unknown run configuration found. Defaulting to --local.")
		return true
	}
}

// getFlagPathOrPanic gets the override path from the flag, returning a default
// value if the flag was unused, or panics if the provided path does not exist.
func getFlagPathOrPanic(c *cobra.Command, flagName string, defaultValue string) string {
	flag := c.Flags().Lookup(flagName)
	if flag != nil && flag.Changed {
		path := flag.Value.String()
		if PathExists(path) {
			log.Info(fmt.Sprintf("Override path provided for %s binary: %s", flagName, path))
			return path
		} else {
			panic(fmt.Sprintf("Invalid input path provided for --%s flag", flagName))
		}
	} else {
		log.Debug(fmt.Sprintf("No path provided for %s binary. Using default path: %s", flagName, defaultValue))
		return defaultValue
	}
}
