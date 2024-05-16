package cmd

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var EnvPrefix = "ASTRIA_GO"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "astria-go",
	Short:            "A CLI to run Astria, interact with the Sequencer, deploy rollups, and more.",
	PersistentPreRun: validateLogLevels,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := RootCmd.ExecuteContext(ctx)
	if err != nil {
		log.WithError(err).Error("Error executing root command")
		panic(err)
	}
}

func init() {
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	// disabling the completion command for now
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	// TODO - update flags to the new flag handler
	RootCmd.PersistentFlags().StringVar(&cliLogLevel, "log-level", "info", "cli log level (debug, info, warn, error, fatal, panic)")
	err := viper.BindPFlag("log_level", RootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		log.Fatalf("Error binding flag: %s", err)
	}
	cliLogLevel = viper.GetString("log_level")
	RootCmd.PersistentFlags().StringVar(&serviceLogLevel, "service-log-level", "info", "services log level (debug, info, error)")
	err = viper.BindPFlag("service_log_level", RootCmd.PersistentFlags().Lookup("service-log-level"))
	if err != nil {
		log.Fatalf("Error binding flag: %s", err)
	}
	serviceLogLevel = viper.GetString("service_log_level")
}

// validateLogLevels validates the log levels passed in by the user for both the
// cli and the services.
func validateLogLevels(_ *cobra.Command, _ []string) {
	switch cliLogLevel {
	case "debug", "info", "warn", "error", "fatal", "panic":
	default:
		log.WithField("log-level", cliLogLevel).Fatal("Invalid cli log level. Must be one of: 'debug', 'info', 'warn', 'error', 'fatal', 'panic'")
		panic("invalid cli log level")

	}
	switch serviceLogLevel {
	case "debug", "info", "error":
	default:
		log.WithField("service-log-level", serviceLogLevel).Fatal("Invalid services log level. Must be one of: 'debug', 'info', 'error'")
		panic("invalid service log level")
	}
}
