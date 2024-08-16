package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TUIConfig is the struct that holds the configuration and start state for the
// TUI.
type TUIConfig struct {
	// Service settings
	AutoScroll bool `mapstructure:"auto_scroll" toml:"auto_scroll"`
	WrapLines  bool `mapstructure:"wrap_lines" toml:"wrap_lines"`
	Borderless bool `mapstructure:"borderless" toml:"borderless"`

	// Override value for the Default Instance Name
	OverrideInstanceName string `mapstructure:"override_instance_name" toml:"override_instance_name"`

	// Known services start minimized
	CometBFTStartsMinimized  bool `mapstructure:"cometbft_starts_minimized" toml:"cometbft_starts_minimized"`
	ConductorStartsMinimized bool `mapstructure:"conductor_starts_minimized" toml:"conductor_starts_minimized"`
	ComposerStartsMinimized  bool `mapstructure:"composer_starts_minimized" toml:"composer_starts_minimized"`
	SequencerStartsMinimized bool `mapstructure:"sequencer_starts_minimized" toml:"sequencer_starts_minimized"`
	// Generic services start minimized
	GenericStartsMinimized bool `mapstructure:"generic_starts_minimized" toml:"generic_starts_minimized"`
}

// DefaultTUIConfig returns a TUIConfig struct populated with all default
// values.
func DefaultTUIConfig() TUIConfig {
	return TUIConfig{
		AutoScroll:               true,
		WrapLines:                true,
		Borderless:               false,
		OverrideInstanceName:     "default",
		CometBFTStartsMinimized:  false,
		ConductorStartsMinimized: false,
		ComposerStartsMinimized:  false,
		SequencerStartsMinimized: false,
		GenericStartsMinimized:   true,
	}
}

// LoadTUIConfigsOrPanic loads the TUIConfigs from the given path. If the file
// cannot be loaded or parsed, the function will panic.
func LoadTUIConfigsOrPanic(path string) TUIConfig {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}

	var config TUIConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}

	return config
}

// CreateNetworksConfig creates a TUI configuration file and populates it
// with the defaults for the devrunner TUI. It will panic if the file
// cannot be created or written to.
func CreateTUIConfig(configPath string) {
	_, err := os.Stat(configPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", configPath)
		return
	}
	// create an instance of the Config struct with some data
	config := DefaultTUIConfig()

	// open a file for writing
	file, err := os.Create(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// encode the struct to TOML and write to the file
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		panic(err)
	}
	log.Infof("New TUI config file created successfully: %s\n", configPath)
}
