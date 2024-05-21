package devtools

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/astria/astria-cli-go/cmd/devtools/config"
	util "github.com/astria/astria-cli-go/cmd/devtools/utilities"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all the Astria services locally.",
	Long:  `Run all the Astria services locally. This will start the sequencer, cometbft, composer, and conductor.`,
	Run:   runCmdHandler,
}

func init() {
	devCmd.AddCommand(runCmd)

	flagHandler := cmd.CreateCliFlagHandler(runCmd, cmd.EnvPrefix)
	flagHandler.BindStringFlag("service-log-level", "info", "Set the log level for services (debug, info, error)")
	flagHandler.BindStringFlag("network", "local", "Provide an override path to a specific environment file.")
	flagHandler.BindStringFlag("conductor-path", "", "Provide an override path to a specific conductor binary.")
	flagHandler.BindStringFlag("cometbft-path", "", "Provide an override path to a specific cometbft binary.")
	flagHandler.BindStringFlag("composer-path", "", "Provide an override path to a specific composer binary.")
	flagHandler.BindStringFlag("sequencer-path", "", "Provide an override path to a specific sequencer binary.")
}

func runCmdHandler(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	ctx := c.Context()

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}

	// astriaDir is the directory where all the astria instances data is stored
	astriaDir := filepath.Join(homePath, ".astria")

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	cmd.CreateUILog(filepath.Join(astriaDir, instance))

	network := flagHandler.GetValue("network")

	baseConfigPath := filepath.Join(astriaDir, instance, config.DefaultConfigDirName, config.DefualtBaseConfigName)
	baseConfig := config.LoadBaseConfigOrPanic(baseConfigPath)
	baseConfigEnvVars := config.ConvertStructToEnvArray(baseConfig)

	networksConfigPath := filepath.Join(astriaDir, instance, config.DefualtNetworksConfigName)
	networkConfigs := config.LoadNetworksConfigsOrPanic(networksConfigPath)

	// update the log level for the Astria Services using override env vars.
	// The log level for Cometbft is updated via command line flags and is set
	// in the ProcessRunnerOpts for the Cometbft ProcessRunner
	serviceLogLevel := flagHandler.GetValue("service-log-level")
	ValidateServiceLogLevelOrPanic(serviceLogLevel)
	serviceLogLevelOverrides := []string{
		"ASTRIA_SEQUENCER_LOG=\"astria_sequencer=" + serviceLogLevel + "\"",
		"ASTRIA_COMPOSER_LOG=\"astria_composer=" + serviceLogLevel + "\"",
		"ASTRIA_CONDUCTOR_LOG=\"astria_conductor=" + serviceLogLevel + "\"",
	}

	// we will set runners after we decide which binaries we need to run
	var runners []processrunner.ProcessRunner

	// setup services based on network config
	switch network {
	case "local":
		networkOverrides := networkConfigs.Local.GetEnvOverrides(baseConfig)
		networkOverrides = config.MergeConfig(baseConfigEnvVars, networkOverrides)
		networkOverrides = config.MergeConfig(networkOverrides, serviceLogLevelOverrides)
		config.LogConfig(networkOverrides)

		log.Debug("Running local sequencer")
		dataDir := filepath.Join(astriaDir, instance, config.DataDirName)
		binDir := filepath.Join(astriaDir, instance, config.BinariesDirName)

		// get the binary paths
		conductorPath := getFlagPath(c, "conductor-path", "conductor", filepath.Join(binDir, "astria-conductor"))
		cometbftPath := getFlagPath(c, "cometbft-path", "cometbft", filepath.Join(binDir, "cometbft"))
		composerPath := getFlagPath(c, "composer-path", "composer", filepath.Join(binDir, "astria-composer"))
		sequencerPath := getFlagPath(c, "sequencer-path", "sequencer", filepath.Join(binDir, "astria-sequencer"))
		log.Debugf("Using binaries from %s", binDir)

		// sequencer
		seqRCOpts := processrunner.ReadyCheckerOpts{
			CallBackName:  "Sequencer gRPC server is OK",
			Callback:      getSequencerOKCallback(networkOverrides),
			RetryCount:    10,
			RetryInterval: 100 * time.Millisecond,
			HaltIfFailed:  false,
		}
		seqReadinessCheck := processrunner.NewReadyChecker(seqRCOpts)
		seqOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Sequencer",
			BinPath:    sequencerPath,
			Config:     networkOverrides,
			Args:       nil,
			ReadyCheck: &seqReadinessCheck,
		}
		seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

		// cometbft
		cometRCOpts := processrunner.ReadyCheckerOpts{
			CallBackName:  "CometBFT rpc server is OK",
			Callback:      getCometbftOKCallback(networkOverrides),
			RetryCount:    10,
			RetryInterval: 100 * time.Millisecond,
			HaltIfFailed:  false,
		}
		cometReadinessCheck := processrunner.NewReadyChecker(cometRCOpts)
		cometDataPath := filepath.Join(dataDir, ".cometbft")
		cometOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Comet BFT",
			BinPath:    cometbftPath,
			Config:     networkOverrides,
			Args:       []string{"node", "--home", cometDataPath, "--log_level", serviceLogLevel},
			ReadyCheck: &cometReadinessCheck,
		}
		cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

		// composer
		composerOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Composer",
			BinPath:    composerPath,
			Config:     networkOverrides,
			Args:       nil,
			ReadyCheck: nil,
		}
		compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

		// conductor
		conductorOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Conductor",
			BinPath:    conductorPath,
			Config:     networkOverrides,
			Args:       nil,
			ReadyCheck: nil,
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

	case "dusk", "dawn", "mainnet":
		var networkOverrides []string
		if network == "dusk" {
			networkOverrides = networkConfigs.Dusk.GetEnvOverrides(baseConfig)
		} else if network == "dawn" {
			networkOverrides = networkConfigs.Dawn.GetEnvOverrides(baseConfig)
		} else {
			networkOverrides = networkConfigs.Mainnet.GetEnvOverrides(baseConfig)
		}
		networkOverrides = config.MergeConfig(baseConfigEnvVars, networkOverrides)
		networkOverrides = config.MergeConfig(networkOverrides, serviceLogLevelOverrides)
		config.LogConfig(networkOverrides)

		log.Debug("Running remote sequencer")
		binDir := filepath.Join(astriaDir, instance, config.BinariesDirName)

		// get the binary paths
		conductorPath := getFlagPath(c, "conductor_bin_path", "conductor", filepath.Join(binDir, "astria-conductor"))
		composerPath := getFlagPath(c, "composer_bin_path", "composer", filepath.Join(binDir, "astria-composer"))

		// composer
		composerOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Composer",
			BinPath:    composerPath,
			Config:     networkOverrides,
			Args:       nil,
			ReadyCheck: nil,
		}
		compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

		// conductor
		conductorOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Conductor",
			BinPath:    conductorPath,
			Config:     networkOverrides,
			Args:       nil,
			ReadyCheck: nil,
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

	default:
		log.Fatalf("Invalid network provided: %s", network)
		log.Fatalf("Valid networks are: local, dusk, dawn, mainnet")
		panic("Invalid network provided")
	}

	// create and start ui app
	app := ui.NewApp(runners)
	app.Start()
}

// getFlagPath gets the override path from the flag. It returns the default
// value if the flag was not set.
func getFlagPath(c *cobra.Command, flag string, serviceName string, defaultValue string) string {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	path := flagHandler.GetValue(flag)

	if util.PathExists(path) && path != "" {
		log.Info(fmt.Sprintf("getFlagPath: Override path provided for %s binary: %s", serviceName, path))
		return path
	} else {
		log.Info(fmt.Sprintf("getFlagPath: Invalid input path provided for --%s flag. Using default: %s", flag, defaultValue))
		return defaultValue
	}
}

// getSequencerOKCallback builds an anonymous function for use in a ProcessRunner
// ReadyChecker callback. The anonymous function checks if the gRPC server that
// is started by the sequencer is OK by making an HTTP request to the health
// endpoint. Being able to connect to the gRPC server is a requirement for both
// the Conductor and Composer services.
func getSequencerOKCallback(config []string) func() bool {
	// Get the sequencer gRPC address from the environment
	var seqGRPCAddr string
	for _, envVar := range config {
		if strings.HasPrefix(envVar, "ASTRIA_SEQUENCER_GRPC_ADDR") {
			seqGRPCAddr = strings.Split(envVar, "=")[1]
			break
		}
	}
	seqGRPCHealthURL := "http://" + seqGRPCAddr + "/health"

	// Return the anonymous callback function
	return func() bool {
		// Make the HTTP request
		seqResp, err := http.Get(seqGRPCHealthURL)
		if err != nil {
			log.WithError(err).Debug("Startup callback check to sequencer gRPC /health did not succeed")
			return false
		}
		defer seqResp.Body.Close()

		// Check status code
		if seqResp.StatusCode == 200 {
			log.Debug("Sequencer gRPC server started")
			return true
		} else {
			log.Debugf("Sequencer gRPC status code: %d", seqResp.StatusCode)
			return false
		}
	}
}

// getCometbftOKCallback builds an anonymous function for use in a ProcessRunner
// ReadyChecker callback. The anonymous function checks if the rpc server that
// is started by Cometbft is OK by making an HTTP request to the health
// endpoint. Being able to connect to the rpc server is a requirement for both
// the Conductor and Composer services.
func getCometbftOKCallback(config []string) func() bool {
	// Get the CometBFT rpc address from the environment
	var seqRPCAddr string
	for _, envVar := range config {
		if strings.HasPrefix(envVar, "ASTRIA_CONDUCTOR_SEQUENCER_COMETBFT_URL") {
			seqRPCAddr = strings.Split(envVar, "=")[1]
			break
		}
	}
	cometbftRPCHealthURL := seqRPCAddr + "/health"

	// Return the anonymous callback function
	return func() bool {
		// Make the HTTP request
		cometbftResp, err := http.Get(cometbftRPCHealthURL)
		if err != nil {
			log.WithError(err).Debug("Startup callback check to CometBFT rpc /health did not succeed")
			return false
		}
		defer cometbftResp.Body.Close()

		// Check status code
		if cometbftResp.StatusCode == 200 {
			log.Debug("CometBFT rpc server started")
			return true
		} else {
			log.Debugf("CometBFT rpc status code: %d", cometbftResp.StatusCode)
			return false
		}
	}
}

// ValidateServiceLogLevelOrPanic validates the service log level and panics if
// it is invalid. The valid log levels are: debug, info, error.
func ValidateServiceLogLevelOrPanic(logLevel string) {
	switch logLevel {
	case "debug", "info", "error":
		return
	default:
		log.WithField("service-log-level", logLevel).Fatal("Invalid service log level. Must be one of: 'debug', 'info', 'error'")
		panic("Invalid service log level")
	}

}
