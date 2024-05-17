package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CliStringFlagHandler is a struct that handles the binding and retrieval of
// string flag values.
type CliStringFlagHandler struct {
	Cmd       *cobra.Command
	EnvPrefix string
}

// CreateCliStringFlagHandler creates a new CliStringFlagHandler.
func CreateCliStringFlagHandler(c *cobra.Command, envPrefix string) *CliStringFlagHandler {
	return &CliStringFlagHandler{
		Cmd:       c,
		EnvPrefix: envPrefix,
	}
}

// BindStringFlag binds a string flag to a cobra flag and viper env var handler for a
// local command flag, and automatically creates the env var from the flag name.
func (f *CliStringFlagHandler) BindStringFlag(name string, defaultValue string, usage string) {
	envSuffix := strings.ToUpper(strings.ReplaceAll(name, "-", "_"))

	f.Cmd.Flags().String(name, defaultValue, usage)
	err := viper.BindPFlag(envSuffix, f.Cmd.Flags().Lookup(name))
	if err != nil {
		log.Fatalf("Error binding string flag: %s", err)
	}
}

// BindBoolFlag binds a boolean flag to a cobra flag and viper env var handler for a
// local command flag, and automatically creates the env var from the flag name.
func (f *CliStringFlagHandler) BindBoolFlag(name string, defaultValue bool, usage string) {
	envSuffix := strings.ToUpper(strings.ReplaceAll(name, "-", "_"))

	f.Cmd.Flags().Bool(name, defaultValue, usage)
	err := viper.BindPFlag(envSuffix, f.Cmd.Flags().Lookup(name))
	if err != nil {
		log.Fatalf("Error binding bool flag: %s", err)
	}
}

// BindPersistentFlag binds a string flag to a cobra flag and viper env var
// handler for a persistent command flag shared by a command and its
// subcommands, and automatically creates the env var from the flag name.
func (f *CliStringFlagHandler) BindPersistentFlag(name string, defaultValue string, usage string) {
	envSuffix := strings.ToUpper(strings.ReplaceAll(name, "-", "_"))

	f.Cmd.PersistentFlags().String(name, defaultValue, usage)
	err := viper.BindPFlag(envSuffix, f.Cmd.PersistentFlags().Lookup(name))
	if err != nil {
		log.Fatalf("Error binding persistent flag: %s", err)
	}
}

// getEnvVar returns the full env var name for a given flag name.
func (f *CliStringFlagHandler) getEnvVar(flagName string) string {
	envSuffix := strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	fullEnvVar := strings.ToUpper(f.EnvPrefix) + "_" + strings.ToUpper(envSuffix)
	return fullEnvVar
}

// GetValue returns the value of a flag and logs the source of the value. It
// will panic if the flag does not exist.
func (f *CliStringFlagHandler) GetValue(flagName string) string {
	envSuffix := strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	value := viper.GetString(envSuffix)
	exists := f.Cmd.Flag(flagName)
	if exists == nil {
		log.Errorf("Flag '%s' doesn't exist. Has it been bound?", flagName)
		panic(fmt.Sprintf("getValue: flag doesn't exist: %s", flagName))
	}

	if f.Cmd.Flag(flagName).Changed {
		log.Debugf("%s flag is set with value: %s", flagName, value)
		return value
	}
	_, envExists := os.LookupEnv(f.getEnvVar(flagName))
	if envExists {
		log.Debugf("%s flag is set via env var to: %s", flagName, value)
		return value
	}

	log.Debugf("%s flag is not set, using default: %s", flagName, value)
	return value
}
