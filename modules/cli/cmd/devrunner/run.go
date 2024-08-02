package devrunner

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
	util "github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/utilities"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/processrunner"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/ui"
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
	flagHandler.BindStringFlag("service-log-level", config.DefaultServiceLogLevel, "Set the log level for services (debug, info, error)")
	flagHandler.BindStringFlag("network", config.DefaultTargetNetwork, "Select the network to run the services against. Valid networks are: local, dusk, dawn, mainnet")
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

	baseConfigPath := filepath.Join(astriaDir, instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)
	baseConfig := config.LoadBaseConfigOrPanic(baseConfigPath)
	baseConfigEnvVars := baseConfig.ToSlice()

	networksConfigPath := filepath.Join(astriaDir, instance, config.DefaultNetworksConfigName)
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	// get the log level for the Astria Services using override env vars.
	// The log level for Cometbft is updated via command line flags and is set
	// in the ProcessRunnerOpts for the Cometbft ProcessRunner
	serviceLogLevel := flagHandler.GetValue("service-log-level")
	serviceLogLevelOverrides := config.GetServiceLogLevelOverrides(serviceLogLevel)

	networkOverrides := networkConfigs.Configs[network].GetEndpointOverrides(baseConfig)

	environment := config.MergeConfigs(baseConfigEnvVars, networkOverrides, serviceLogLevelOverrides)
	config.LogEnv(environment)

	binDir := filepath.Join(astriaDir, instance, config.BinariesDirName)

	// we will set runners after we decide which binaries we need to run
	var runners []processrunner.ProcessRunner
	// known services
	var seqRunner processrunner.ProcessRunner
	var cometRunner processrunner.ProcessRunner
	var compRunner processrunner.ProcessRunner
	var condRunner processrunner.ProcessRunner
	// generic services
	// TODO: add to default section
	// var genericRunners []processrunner.ProcessRunner

	// TODO: add generic print to say, "running blah services for network blah"

	// TODO: load the services from the networks config based on network name
	// and build the process runners for each service, with special treatment
	// for "known" services like sequencer, composer, conductor, and cometbft
	// this should all be dynamic and based on the network config
	// also need to order known services then anything else at the end
	// TODO: revert just label to "label, service" service for use in default section
	for label := range networkConfigs.Configs[network].Services {
		// perform specific setup for known services
		switch label {
		case "sequencer":
			sequencerPath := getFlagPath(c, "sequencer-path", "sequencer", filepath.Join(binDir, "astria-sequencer-v"+config.AstriaSequencerVersion))
			seqRCOpts := processrunner.ReadyCheckerOpts{
				CallBackName:  "Sequencer gRPC server is OK",
				Callback:      getSequencerOKCallback(environment),
				RetryCount:    10,
				RetryInterval: 100 * time.Millisecond,
				HaltIfFailed:  false,
			}
			seqReadinessCheck := processrunner.NewReadyChecker(seqRCOpts)
			seqOpts := processrunner.NewProcessRunnerOpts{
				Title:      "Sequencer",
				BinPath:    sequencerPath,
				Env:        environment,
				Args:       nil,
				ReadyCheck: &seqReadinessCheck,
			}
			seqRunner = processrunner.NewProcessRunner(ctx, seqOpts)
		case "composer":
			composerPath := getFlagPath(c, "composer-path", "composer", filepath.Join(binDir, "astria-composer-v"+config.AstriaComposerVersion))
			composerOpts := processrunner.NewProcessRunnerOpts{
				Title:      "Composer",
				BinPath:    composerPath,
				Env:        environment,
				Args:       nil,
				ReadyCheck: nil,
			}
			compRunner = processrunner.NewProcessRunner(ctx, composerOpts)
		case "conductor":
			conductorPath := getFlagPath(c, "conductor-path", "conductor", filepath.Join(binDir, "astria-conductor-v"+config.AstriaConductorVersion))
			conductorOpts := processrunner.NewProcessRunnerOpts{
				Title:      "Conductor",
				BinPath:    conductorPath,
				Env:        environment,
				Args:       nil,
				ReadyCheck: nil,
			}
			condRunner = processrunner.NewProcessRunner(ctx, conductorOpts)
		case "cometbft":
			cometbftPath := getFlagPath(c, "cometbft-path", "cometbft", filepath.Join(binDir, "cometbft-v"+config.CometbftVersion))
			cometRCOpts := processrunner.ReadyCheckerOpts{
				CallBackName:  "CometBFT rpc server is OK",
				Callback:      getCometbftOKCallback(environment),
				RetryCount:    10,
				RetryInterval: 100 * time.Millisecond,
				HaltIfFailed:  false,
			}
			cometReadinessCheck := processrunner.NewReadyChecker(cometRCOpts)
			dataDir := filepath.Join(astriaDir, instance, config.DataDirName)
			cometDataPath := filepath.Join(dataDir, ".cometbft")
			cometOpts := processrunner.NewProcessRunnerOpts{
				Title:      "Comet BFT",
				BinPath:    cometbftPath,
				Env:        environment,
				Args:       []string{"node", "--home", cometDataPath, "--log_level", serviceLogLevel},
				ReadyCheck: &cometReadinessCheck,
			}
			cometRunner = processrunner.NewProcessRunner(ctx, cometOpts)
		default:
			// anything else
			log.Info("generic process runner not implemented yet")
		}
	}

	// TODO: make this dynamic based on the services in the network
	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	close(shouldStart)
	err = seqRunner.Start(ctx, shouldStart)
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

	// TODO: fill runners with all known services in the correct order, then add
	// all generic services at the end
	runners = []processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}

	// setup services based on network config
	switch network {
	case "local":
		// networkOverrides := networkConfigs.Configs["local"].GetEndpointOverrides(baseConfig)
		// serviceLogLevelOverrides := config.GetServiceLogLevelOverrides(serviceLogLevel)
		// environment := config.MergeConfigs(baseConfigEnvVars, networkOverrides, serviceLogLevelOverrides)
		// config.LogEnv(environment)

		// log.Debug("Running local sequencer")
		// dataDir := filepath.Join(astriaDir, instance, config.DataDirName)
		// binDir := filepath.Join(astriaDir, instance, config.BinariesDirName)

		// get the binary paths
		// TODO: known services should be loaded from the networks config
		// conductorPath := getFlagPath(c, "conductor-path", "conductor", filepath.Join(binDir, "astria-conductor-v"+config.AstriaConductorVersion))
		// cometbftPath := getFlagPath(c, "cometbft-path", "cometbft", filepath.Join(binDir, "cometbft-v"+config.CometbftVersion))
		// composerPath := getFlagPath(c, "composer-path", "composer", filepath.Join(binDir, "astria-composer-v"+config.AstriaComposerVersion))
		// sequencerPath := getFlagPath(c, "sequencer-path", "sequencer", filepath.Join(binDir, "astria-sequencer-v"+config.AstriaSequencerVersion))
		// log.Debugf("Using binaries from %s", binDir)

		// // sequencer
		// seqRCOpts := processrunner.ReadyCheckerOpts{
		// 	CallBackName:  "Sequencer gRPC server is OK",
		// 	Callback:      getSequencerOKCallback(environment),
		// 	RetryCount:    10,
		// 	RetryInterval: 100 * time.Millisecond,
		// 	HaltIfFailed:  false,
		// }
		// seqReadinessCheck := processrunner.NewReadyChecker(seqRCOpts)
		// seqOpts := processrunner.NewProcessRunnerOpts{
		// 	Title:      "Sequencer",
		// 	BinPath:    sequencerPath,
		// 	Env:        environment,
		// 	Args:       nil,
		// 	ReadyCheck: &seqReadinessCheck,
		// }
		// seqRunner := processrunner.NewProcessRunner(ctx, seqOpts)

		// // cometbft
		// cometRCOpts := processrunner.ReadyCheckerOpts{
		// 	CallBackName:  "CometBFT rpc server is OK",
		// 	Callback:      getCometbftOKCallback(environment),
		// 	RetryCount:    10,
		// 	RetryInterval: 100 * time.Millisecond,
		// 	HaltIfFailed:  false,
		// }
		// cometReadinessCheck := processrunner.NewReadyChecker(cometRCOpts)
		// cometDataPath := filepath.Join(dataDir, ".cometbft")
		// cometOpts := processrunner.NewProcessRunnerOpts{
		// 	Title:      "Comet BFT",
		// 	BinPath:    cometbftPath,
		// 	Env:        environment,
		// 	Args:       []string{"node", "--home", cometDataPath, "--log_level", serviceLogLevel},
		// 	ReadyCheck: &cometReadinessCheck,
		// }
		// cometRunner := processrunner.NewProcessRunner(ctx, cometOpts)

		// // composer
		// composerOpts := processrunner.NewProcessRunnerOpts{
		// 	Title:      "Composer",
		// 	BinPath:    composerPath,
		// 	Env:        environment,
		// 	Args:       nil,
		// 	ReadyCheck: nil,
		// }
		// compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

		// // conductor
		// conductorOpts := processrunner.NewProcessRunnerOpts{
		// 	Title:      "Conductor",
		// 	BinPath:    conductorPath,
		// 	Env:        environment,
		// 	Args:       nil,
		// 	ReadyCheck: nil,
		// }
		// condRunner := processrunner.NewProcessRunner(ctx, conductorOpts)

		// // shouldStart acts as a control channel to start this first process
		// shouldStart := make(chan bool)
		// close(shouldStart)
		// err := seqRunner.Start(ctx, shouldStart)
		// if err != nil {
		// 	log.WithError(err).Error("Error running sequencer")
		// }
		// err = cometRunner.Start(ctx, seqRunner.GetDidStart())
		// if err != nil {
		// 	log.WithError(err).Error("Error running cometbft")
		// }
		// err = compRunner.Start(ctx, cometRunner.GetDidStart())
		// if err != nil {
		// 	log.WithError(err).Error("Error running composer")
		// }
		// err = condRunner.Start(ctx, compRunner.GetDidStart())
		// if err != nil {
		// 	log.WithError(err).Error("Error running conductor")
		// }

		// runners = []processrunner.ProcessRunner{seqRunner, cometRunner,
		// compRunner, condRunner}
		log.Info("old local")

	case "dusk", "dawn", "mainnet":
		networkOverrides := networkConfigs.Configs[network].GetEndpointOverrides(baseConfig)
		serviceLogLevelOverrides := config.GetServiceLogLevelOverrides(serviceLogLevel)
		environment := config.MergeConfigs(baseConfigEnvVars, networkOverrides, serviceLogLevelOverrides)
		config.LogEnv(environment)

		log.Debug("Running remote sequencer")
		binDir := filepath.Join(astriaDir, instance, config.BinariesDirName)

		// get the binary paths
		conductorPath := getFlagPath(c, "conductor-path", "conductor", filepath.Join(binDir, "astria-conductor"))
		composerPath := getFlagPath(c, "composer-path", "composer", filepath.Join(binDir, "astria-composer"))

		// composer
		composerOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Composer",
			BinPath:    composerPath,
			Env:        environment,
			Args:       nil,
			ReadyCheck: nil,
		}
		compRunner := processrunner.NewProcessRunner(ctx, composerOpts)

		// conductor
		conductorOpts := processrunner.NewProcessRunnerOpts{
			Title:      "Conductor",
			BinPath:    conductorPath,
			Env:        environment,
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
	// get the sequencer gRPC address from the environment
	var seqGRPCAddr string
	for _, envVar := range config {
		if strings.HasPrefix(envVar, "ASTRIA_SEQUENCER_GRPC_ADDR") {
			seqGRPCAddr = strings.Split(envVar, "=")[1]
			break
		}
	}
	seqGRPCHealthURL := "http://" + seqGRPCAddr + "/health"

	// return the anonymous callback function
	return func() bool {
		// make the HTTP request
		seqResp, err := http.Get(seqGRPCHealthURL)
		if err != nil {
			log.WithError(err).Debug("Startup callback check to sequencer gRPC /health did not succeed")
			return false
		}
		defer seqResp.Body.Close()

		// check status code
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
	// get the CometBFT rpc address from the environment
	var seqRPCAddr string
	for _, envVar := range config {
		if strings.HasPrefix(envVar, "ASTRIA_CONDUCTOR_SEQUENCER_COMETBFT_URL") {
			seqRPCAddr = strings.Split(envVar, "=")[1]
			break
		}
	}
	cometbftRPCHealthURL := seqRPCAddr + "/health"

	// return the anonymous callback function
	return func() bool {
		// make the HTTP request
		cometbftResp, err := http.Get(cometbftRPCHealthURL)
		if err != nil {
			log.WithError(err).Debug("Startup callback check to CometBFT rpc /health did not succeed")
			return false
		}
		defer cometbftResp.Body.Close()

		// check status code
		if cometbftResp.StatusCode == 200 {
			log.Debug("CometBFT rpc server started")
			return true
		} else {
			log.Debugf("CometBFT rpc status code: %d", cometbftResp.StatusCode)
			return false
		}
	}
}
