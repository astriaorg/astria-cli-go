package devrunner

import (
	"context"
	"fmt"
	"net/http"
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
	flagHandler.BindBoolFlag("export-logs", false, "Export logs to files.")
}

func runCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)
	ctx := c.Context()

	homeDir := cmd.GetUserHomeDirOrPanic()

	astriaDir := filepath.Join(homeDir, ".astria")

	tuiConfigPath := filepath.Join(astriaDir, config.DefaultTUIConfigName)
	tuiConfig := config.LoadTUIConfigsOrPanic(tuiConfigPath)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	// set the instance name from the correct source so that logging is handled
	// correctly
	if !flagHandler.GetChanged("instance") {
		instance = tuiConfig.OverrideInstanceName
	}

	exportLogs := flagHandler.GetValue("export-logs") == "true"
	logsDir := filepath.Join(astriaDir, instance, config.LogsDirName)
	currentTime := time.Now()
	appStartTime := currentTime.Format("20060102-150405") // YYYYMMDD-HHMMSS

	cmd.CreateUILog(filepath.Join(astriaDir, instance))

	// log the instance name in the tui logs once they are created
	if !flagHandler.GetChanged("instance") {
		log.Debug("Using overridden default instance name: ", instance)
	} else {
		log.Debug("Instance name: ", instance)
	}
	log.Debug(tuiConfig)

	network := flagHandler.GetValue("network")

	baseConfigPath := filepath.Join(astriaDir, instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)
	baseConfig := config.LoadBaseConfigOrPanic(baseConfigPath)
	baseConfigEnvVars := baseConfig.ToSlice()

	networksConfigPath := filepath.Join(astriaDir, instance, config.DefaultNetworksConfigName)
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	// check if the network exists in the networks config
	if _, ok := networkConfigs.Configs[network]; !ok {
		log.Fatalf("Network %s not found in config file at %s", network, networksConfigPath)
		panic("Network not found in config file")
	}

	// get the log level for the Astria Services using override env vars.
	// The log level for Cometbft is updated via command line flags and is set
	// in the ProcessRunnerOpts for the Cometbft ProcessRunner
	serviceLogLevel := flagHandler.GetValue("service-log-level")
	serviceLogLevelOverrides := config.GetServiceLogLevelOverrides(serviceLogLevel)

	networkOverrides := networkConfigs.Configs[network].GetEndpointOverrides(baseConfig)

	environment := config.MergeConfigs(baseConfigEnvVars, networkOverrides, serviceLogLevelOverrides)
	config.LogEnv(environment)

	// known services
	var seqRunner processrunner.ProcessRunner
	var cometRunner processrunner.ProcessRunner
	var compRunner processrunner.ProcessRunner
	var condRunner processrunner.ProcessRunner
	// generic services
	var genericRunners []processrunner.ProcessRunner

	// load the services from the networks config and build the process runners
	// for each service, with special treatment for "known" services like
	// sequencer, composer, conductor, and cometbft
	for label, service := range networkConfigs.Configs[network].Services {
		switch label {
		case "sequencer":
			sequencerPath := getFlagPath(c, "sequencer-path", "sequencer", service.LocalPath)
			seqRCOpts := processrunner.ReadyCheckerOpts{
				CallBackName:  "Sequencer gRPC server is OK",
				Callback:      getSequencerOKCallback(environment),
				RetryCount:    10,
				RetryInterval: 100 * time.Millisecond,
				HaltIfFailed:  false,
			}
			seqReadinessCheck := processrunner.NewReadyChecker(seqRCOpts)
			log.Debugf("arguments for sequencer service: %v", service.Args)
			seqOpts := processrunner.NewProcessRunnerOpts{
				Title:          "Sequencer",
				BinPath:        sequencerPath,
				Env:            environment,
				Args:           service.Args,
				ReadyCheck:     &seqReadinessCheck,
				LogPath:        filepath.Join(logsDir, appStartTime+"-astria-sequencer.log"),
				ExportLogs:     exportLogs,
				StartMinimized: tuiConfig.SequencerStartsMinimized,
			}
			seqRunner = processrunner.NewProcessRunner(ctx, seqOpts)
		case "composer":
			composerPath := getFlagPath(c, "composer-path", "composer", service.LocalPath)
			log.Debugf("arguments for composer service: %v", service.Args)
			composerOpts := processrunner.NewProcessRunnerOpts{
				Title:          "Composer",
				BinPath:        composerPath,
				Env:            environment,
				Args:           service.Args,
				ReadyCheck:     nil,
				LogPath:        filepath.Join(logsDir, appStartTime+"-astria-composer.log"),
				ExportLogs:     exportLogs,
				StartMinimized: tuiConfig.ComposerStartsMinimized,
			}
			compRunner = processrunner.NewProcessRunner(ctx, composerOpts)
		case "conductor":
			conductorPath := getFlagPath(c, "conductor-path", "conductor", service.LocalPath)
			log.Debugf("arguments for conductor service: %v", service.Args)
			conductorOpts := processrunner.NewProcessRunnerOpts{
				Title:          "Conductor",
				BinPath:        conductorPath,
				Env:            environment,
				Args:           service.Args,
				ReadyCheck:     nil,
				LogPath:        filepath.Join(logsDir, appStartTime+"-astria-conductor.log"),
				ExportLogs:     exportLogs,
				StartMinimized: tuiConfig.ConductorStartsMinimized,
			}
			condRunner = processrunner.NewProcessRunner(ctx, conductorOpts)
		case "cometbft":
			cometbftPath := getFlagPath(c, "cometbft-path", "cometbft", service.LocalPath)
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
			args := append([]string{"node", "--home", cometDataPath, "--log_level", serviceLogLevel}, service.Args...)
			log.Debugf("arguments for cometbft service: %v", args)
			cometOpts := processrunner.NewProcessRunnerOpts{
				Title:          "Comet BFT",
				BinPath:        cometbftPath,
				Env:            environment,
				Args:           args,
				ReadyCheck:     &cometReadinessCheck,
				LogPath:        filepath.Join(logsDir, appStartTime+"-cometbft.log"),
				ExportLogs:     exportLogs,
				StartMinimized: tuiConfig.CometBFTStartsMinimized,
			}
			cometRunner = processrunner.NewProcessRunner(ctx, cometOpts)
		default:
			log.Debugf("arguments for %s service: %v", label, service.Args)
			genericOpts := processrunner.NewProcessRunnerOpts{
				Title:          service.Name,
				BinPath:        service.LocalPath,
				Env:            environment,
				Args:           service.Args,
				ReadyCheck:     nil,
				LogPath:        filepath.Join(logsDir, appStartTime+"-"+service.Name+".log"),
				ExportLogs:     exportLogs,
				StartMinimized: tuiConfig.GenericStartsMinimized,
			}
			genericRunner := processrunner.NewProcessRunner(ctx, genericOpts)
			genericRunners = append(genericRunners, genericRunner)
		}
	}

	// set the order of the services to start
	allRunners := append([]processrunner.ProcessRunner{seqRunner, cometRunner, compRunner, condRunner}, genericRunners...)

	runners, err := startProcessInOrder(ctx, allRunners...)
	if err != nil {
		log.WithError(err).Error("Error starting services")
	}

	// create and start ui app
	app := ui.NewApp(runners)
	// start the app with initial setting from the tui config, the border will
	// always start on
	appStartState := ui.BuildStateStore(tuiConfig.AutoScroll, tuiConfig.WrapLines, tuiConfig.Borderless)
	app.Start(appStartState)
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

// startProcessInOrder starts the ProcessRunners in order they are provided, and
// returns an array of all successfully started services. It will skip any
// ProcessRunners that are nil. It will return an error if any of the
// ProcessRunners fail to start.
func startProcessInOrder(ctx context.Context, runners ...processrunner.ProcessRunner) ([]processrunner.ProcessRunner, error) {
	if len(runners) < 1 {
		return nil, fmt.Errorf("no runners provided. Nothing to start")
	}

	var returnRunners []processrunner.ProcessRunner

	previousRunner := runners[0]
	remainingRunners := runners[1:]

	var err error

	// start the first runner
	shouldStart := make(chan bool)
	close(shouldStart)
	if previousRunner != nil {
		err = previousRunner.Start(ctx, shouldStart)
		if err != nil {
			log.WithError(err).Errorf("Error running %s", previousRunner.GetTitle())
			return nil, err
		}
		returnRunners = append(returnRunners, previousRunner)
	}

	if len(remainingRunners) == 0 {
		return returnRunners, nil
	}

	// start the remaining runners
	for _, runner := range remainingRunners {
		if runner != nil {
			if previousRunner == nil {
				err = runner.Start(ctx, shouldStart)
			} else {
				err = runner.Start(ctx, previousRunner.GetDidStart())
			}
			if err != nil {
				log.WithError(err).Errorf("Error running %s", runner.GetTitle())
				return nil, err
			}
			returnRunners = append(returnRunners, runner)
			previousRunner = runner
		}
	}

	return returnRunners, nil
}
