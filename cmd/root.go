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
	PersistentPreRun: validateAndSetCliLogLevel,
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

	flagHandler := CreateCliFlagHandler(RootCmd, EnvPrefix)
	flagHandler.BindPersistentFlag("log-level", "info", "log level (debug, info, warn, error, fatal, panic)")
}

// validateAndSetCliLogLevel validates the log levels passed in by the user for both the
// cli and the services.
func validateAndSetCliLogLevel(c *cobra.Command, _ []string) {
	flagHandler := CreateCliFlagHandler(c, EnvPrefix)
	cliLogLevel = flagHandler.GetValue("log-level")
	switch cliLogLevel {
	case "debug", "info", "warn", "error", "fatal", "panic":
	default:
		log.WithField("log-level", cliLogLevel).Fatal("Invalid cli log level. Must be one of: 'debug', 'info', 'warn', 'error', 'fatal', 'panic'")
		panic("invalid cli log level")
	}
	SetLogLevel(cliLogLevel)
}
