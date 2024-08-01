package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CliFlagHandler is a struct that handles the binding of flags and the retrieval of the flag's string value
type CliFlagHandler struct {
	Cmd           *cobra.Command
	EnvPrefix     string
	useConfigFlag string
	config        any
}

// CreateCliFlagHandler creates a new CliFlagHandler.
func CreateCliFlagHandler(c *cobra.Command, envPrefix string) *CliFlagHandler {
	return &CliFlagHandler{
		Cmd:       c,
		EnvPrefix: envPrefix,
	}
}

// CreateCliFlagHandlerWithUseConfigFlag creates a new CliFlagHandler with a
// specified config override flag. When the config flag is set, that will
// trigger the flag handler to use the config preset values instead of the flag
// defaults.
func CreateCliFlagHandlerWithUseConfigFlag(c *cobra.Command, envPrefix string, configOverrideFlag string) *CliFlagHandler {
	return &CliFlagHandler{
		Cmd:           c,
		EnvPrefix:     envPrefix,
		useConfigFlag: configOverrideFlag,
	}
}

// SetConfig sets the config for the CliFlagHandler.
func (f *CliFlagHandler) SetConfig(config any) {
	f.config = config
}

// BindStringFlag binds a string flag to a cobra flag and viper env var handler for a
// local command flag, and automatically creates the env var from the flag name.
func (f *CliFlagHandler) BindStringFlag(name string, defaultValue string, usage string) {
	envSuffix := strings.ToUpper(strings.ReplaceAll(name, "-", "_"))

	f.Cmd.Flags().String(name, defaultValue, usage)
	err := viper.BindPFlag(envSuffix, f.Cmd.Flags().Lookup(name))
	if err != nil {
		log.Fatalf("Error binding string flag: %s", err)
	}
}

// BindStringPFlag binds a string flag to a cobra flag and viper env var handler for a
// local command flag, and automatically creates the env var from the flag name.
func (f *CliFlagHandler) BindStringPFlag(name string, shorthand string, defaultValue string, usage string) {
	envSuffix := strings.ToUpper(strings.ReplaceAll(name, "-", "_"))

	f.Cmd.Flags().StringP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(envSuffix, f.Cmd.Flags().Lookup(name))
	if err != nil {
		log.Fatalf("Error binding string flag: %s", err)
	}
}

// BindBoolFlag binds a boolean flag to a cobra flag and viper env var handler for a
// local command flag, and automatically creates the env var from the flag name.
func (f *CliFlagHandler) BindBoolFlag(name string, defaultValue bool, usage string) {
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
func (f *CliFlagHandler) BindPersistentFlag(name string, defaultValue string, usage string) {
	envSuffix := strings.ToUpper(strings.ReplaceAll(name, "-", "_"))

	f.Cmd.PersistentFlags().String(name, defaultValue, usage)
	err := viper.BindPFlag(envSuffix, f.Cmd.PersistentFlags().Lookup(name))
	if err != nil {
		log.Fatalf("Error binding persistent flag: %s", err)
	}
}

// toEnvVarName returns the full env var name for a given flag name.
func (f *CliFlagHandler) toEnvVarName(flagName string) string {
	envSuffix := strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	fullEnvVar := strings.ToUpper(f.EnvPrefix) + "_" + strings.ToUpper(envSuffix)
	return fullEnvVar
}

// GetValue returns the value of a flag as a string and logs the source of the
// value. It will panic if the flag does not exist or if the flag cannot be read.
func (f *CliFlagHandler) GetValue(flagName string) string {
	envSuffix := strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	// confirm the flag exists
	flag := f.Cmd.Flags().Lookup(flagName)
	if flag == nil {
		log.Errorf("Flag '%s' doesn't exist. Has it been bound?", flagName)
		panic(fmt.Sprintf("getValue: flag doesn't exist: %s", flagName))
	}

	// get the value from viper based on type
	var value string
	valueKind := reflect.TypeOf(flag.Value).Elem().Kind()
	switch valueKind {
	case reflect.Bool:
		// we need to rebind the flag to viper when reading the value to
		// ensure that it is read correctly. Otherwise, viper will always
		// return the default value.
		err := viper.BindPFlag(envSuffix, flag)
		if err != nil {
			log.Fatalf("getValue: Error rebinding bool flag for reading: %s", flagName)
			panic(err)
		}
		value = fmt.Sprintf("%t", viper.GetBool(envSuffix))

	case reflect.String:
		err := viper.BindPFlag(envSuffix, flag)
		if err != nil {
			log.Fatalf("getValue: Error rebinding string flag for reading: %s", flagName)
			panic(err)
		}
		value = viper.GetString(envSuffix)

	default:
		log.Errorf("Flag '%s' has an unsupported type: %s", flagName, valueKind)
		panic(fmt.Sprintf("getValue: unsupported flag type: %s", valueKind))
	}

	if f.useConfigFlag != "" && f.Cmd.Flag(f.useConfigFlag).Changed {
		// check if value exists in config and return it
		if f.config != nil {
			configValue, found := getFieldByTag(f.config, "flag", flagName)
			if found {
				value := configValue.String()
				log.Debugf("%s flag is set via config file: %s", flagName, value)
				return value
			}
			log.Debugf("Config value for %s not found or invalid. Skipping.", flagName)
		}
	}

	if f.Cmd.Flag(flagName).Changed {
		log.Debugf("%s flag is set with value: %s", flagName, value)
		return value
	}
	_, envExists := os.LookupEnv(f.toEnvVarName(flagName))
	if envExists {
		log.Debugf("%s flag is set via env var to: %s", flagName, value)
		return value
	}

	log.Debugf("%s flag is not set, using default: %s", flagName, value)
	return value
}

// Helper function to get field by tag
func getFieldValueByTag(obj interface{}, tagName, tagValue string) (reflect.Value, bool) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == tagValue {
			return val.Field(i), true
		}
	}
	return reflect.Value{}, false
}
