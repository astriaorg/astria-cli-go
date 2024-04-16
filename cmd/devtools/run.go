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
	runCmd.Flags().String("environment", "", "Provide an override path to a specific environment file.")
	runCmd.Flags().String("conductor", "", "Provide an override path to a specific conductor binary.")
	runCmd.Flags().String("cometbft", "", "Provide an override path to a specific cometbft binary.")
	runCmd.Flags().String("composer", "", "Provide an override path to a specific composer binary.")
	runCmd.Flags().String("sequencer", "", "Provide an override path to a specific sequencer binary.")
}

func runRun(c *cobra.Command, args []string) {
	// parse the input
	runOpts := parseInput(c)

	// generate the process runners
	var runners []processrunner.ProcessRunner
	switch runOpts.runConfiguration {
	case Local:
		runners = getLocalRunners(runOpts)
	case Remote:
		runners = getRemoteRunners(runOpts)
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

type runConfiguration int

const (
	Local runConfiguration = iota
	Remote
)

type runOpts struct {
	ctx              context.Context
	instance         string
	runConfiguration runConfiguration
	environment      []string
	conductorBinPath string
	cometBFTBinPath  string
	composerBinPath  string
	sequencerBinPath string
}

// getRunConfigureation returns the run configuration based on the flags provided.
func getRunConfigureation() runConfiguration {
	switch {
	case !isRunLocal && !isRunRemote:
		log.Debug("No --local or --remote flag provided. Defaulting to --local.")
		return Local
	case isRunLocal:
		log.Debug("--local flag provided. Running local sequencer.")
		return Local
	case isRunRemote:
		log.Debug("--remote flag provided. Connecting to remote sequencer.")
		return Remote
	default:
		// this should never happen
		log.Debug("Unkown run configuration found. Defaulting to --local.")
		return Local
	}
}

// getFlagPathOrPanic gets the override path from the flag, returning a default
// value if the flag was unused, or panics if the provided path does not exist.
func getFlagPathOrPanic(c *cobra.Command, flagName string, defaultValue string) string {
	flag := c.Flags().Lookup(flagName)
	if flag != nil && flag.Changed {
		path := flag.Value.String()
		if pathExists(path) {
			log.Info(fmt.Sprintf("Path provided for %s binary updated to: %s", flagName, path))
			return path
		} else {
			panic(fmt.Sprintf("Path provided for input %s does not exist.", flagName))
		}
	} else {
		log.Debug(fmt.Sprintf("No path provided for %s binary. Using default path: %s", flagName, defaultValue))
		return defaultValue
	}
}

// parseInput creates a runOpts struct from the input flags and returns it.
func parseInput(c *cobra.Command) *runOpts {
	// get the command context
	ctx := c.Context()

	// get the home directory
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	// set the default directory
	defaultDir := filepath.Join(homePath, ".astria")

	// check the instance name
	instance := c.Flag("instance").Value.String()
	IsInstanceNameValidOrPanic(instance)
	instanceDir := filepath.Join(defaultDir, instance)

	// get the run configuration
	runConfiguration := getRunConfigureation()

	// create default paths
	defaultEnvPath := filepath.Join(instanceDir, LocalConfigDirName, ".env")
	defaultBinPath := filepath.Join(instanceDir, BinariesDirName)
	// update the env path based on the run configuration
	switch runConfiguration {
	case Local:
		defaultEnvPath = filepath.Join(instanceDir, LocalConfigDirName, ".env")
	case Remote:
		defaultEnvPath = filepath.Join(instanceDir, RemoteConfigDirName, ".env")
	}

	// get the environment file path
	envPath := getFlagPathOrPanic(c, "environment", defaultEnvPath)
	log.Debug("Using environment file:", envPath)

	// get the binary paths
	conductorPath := getFlagPathOrPanic(c, "conductor", filepath.Join(defaultBinPath, "astria-conductor"))
	cometbftPath := getFlagPathOrPanic(c, "cometbft", filepath.Join(defaultBinPath, "cometbft"))
	composerPath := getFlagPathOrPanic(c, "composer", filepath.Join(defaultBinPath, "astria-composer"))
	sequencerPath := getFlagPathOrPanic(c, "sequencer", filepath.Join(defaultBinPath, "astria-sequencer"))
	log.Debug("Using conductor binary:", conductorPath)
	log.Debug("Using cometbft binary:", cometbftPath)
	log.Debug("Using composer binary:", composerPath)
	log.Debug("Using sequencer binary:", sequencerPath)

	// create the runOpts struct
	runOpts := &runOpts{
		ctx:              ctx,
		instance:         instance,
		runConfiguration: runConfiguration,
		environment:      loadEnvironment(envPath),
		conductorBinPath: conductorPath,
		cometBFTBinPath:  cometbftPath,
		composerBinPath:  composerPath,
		sequencerBinPath: sequencerPath,
	}
	return runOpts
}

func getLocalRunners(opts *runOpts) []processrunner.ProcessRunner {
	// sequencer
	seqOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Sequencer",
		BinPath: opts.sequencerBinPath,
		Env:     opts.environment,
		Args:    nil,
	}
	seqRunner := processrunner.NewProcessRunner(opts.ctx, seqOpts)

	// cometbft
	cometDataPath := filepath.Join(opts.instance, DataDirName, ".cometbft")
	cometOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Comet BFT",
		BinPath: opts.cometBFTBinPath,
		Env:     opts.environment,
		Args:    []string{"node", "--home", cometDataPath},
	}
	cometRunner := processrunner.NewProcessRunner(opts.ctx, cometOpts)

	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: opts.composerBinPath,
		Env:     opts.environment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(opts.ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: opts.conductorBinPath,
		Env:     opts.environment,
		Args:    nil,
	}
	condRunner := processrunner.NewProcessRunner(opts.ctx, conductorOpts)

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err := seqRunner.Start(opts.ctx, shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running sequencer")
	}
	err = cometRunner.Start(opts.ctx, seqRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running cometbft")
	}
	err = compRunner.Start(opts.ctx, cometRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running composer")
	}
	err = condRunner.Start(opts.ctx, compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
	}

	runners := []processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}
	return runners
}

func getRemoteRunners(opts *runOpts) []processrunner.ProcessRunner {
	// composer
	composerOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Composer",
		BinPath: opts.composerBinPath,
		Env:     opts.environment,
		Args:    nil,
	}
	compRunner := processrunner.NewProcessRunner(opts.ctx, composerOpts)

	// conductor
	conductorOpts := processrunner.NewProcessRunnerOpts{
		Title:   "Conductor",
		BinPath: opts.conductorBinPath,
		Env:     opts.environment,
		Args:    nil,
	}
	condRunner := processrunner.NewProcessRunner(opts.ctx, conductorOpts)

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err := compRunner.Start(opts.ctx, shouldStart)
	if err != nil {
		log.WithError(err).Error("Error running composer")
	}
	err = condRunner.Start(opts.ctx, compRunner.GetDidStart())
	if err != nil {
		log.WithError(err).Error("Error running conductor")
	}

	runners := []processrunner.ProcessRunner{compRunner, condRunner}
	return runners
}
