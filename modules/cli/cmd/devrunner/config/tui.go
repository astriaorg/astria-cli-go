package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TUIConfig is the struct that holds the configuration and start state for the
// TUI.
type TUIConfig struct {
	// Log viewer settings for services
	AutoScroll bool `mapstructure:"auto_scroll" toml:"auto_scroll"`
	WrapLines  bool `mapstructure:"wrap_lines" toml:"wrap_lines"`
	Borderless bool `mapstructure:"borderless" toml:"borderless"`

	// Override value for the Instance Name
	OverrideInstanceName string `mapstructure:"override_instance_name" toml:"override_instance_name"`

	// Known services start minimized
	CometBFTStartsMinimized  bool `mapstructure:"cometbft_starts_minimized" toml:"cometbft_starts_minimized"`
	ConductorStartsMinimized bool `mapstructure:"conductor_starts_minimized" toml:"conductor_starts_minimized"`
	ComposerStartsMinimized  bool `mapstructure:"composer_starts_minimized" toml:"composer_starts_minimized"`
	SequencerStartsMinimized bool `mapstructure:"sequencer_starts_minimized" toml:"sequencer_starts_minimized"`
	// Generic services start minimized
	GenericStartsMinimized bool `mapstructure:"generic_starts_minimized" toml:"generic_starts_minimized"`

	// Generic services start position relative to known services
	GenericStartPosition string `mapstructure:"generic_start_position" toml:"generic_start_position"`

	// Accessibility settings
	HighlightColor string `mapstructure:"highlight_color" toml:"highlight_color"`
	BorderColor    string `mapstructure:"border_color" toml:"border_color"`

	// Max number of lines to display in the TUI log viewer for all services
	MaxUiLogLines int `mapstructure:"max_ui_log_lines" toml:"max_ui_log_lines"`
}

// DefaultTUIConfig returns a default TUIConfig struct.
func DefaultTUIConfig() TUIConfig {
	return TUIConfig{
		AutoScroll:               true,
		WrapLines:                false,
		Borderless:               false,
		OverrideInstanceName:     DefaultInstanceName,
		CometBFTStartsMinimized:  false,
		ConductorStartsMinimized: false,
		ComposerStartsMinimized:  false,
		SequencerStartsMinimized: false,
		GenericStartsMinimized:   true,
		GenericStartPosition:     "after",
		HighlightColor:           DefaultHighlightColor,
		BorderColor:              DefaultBorderColor,
		MaxUiLogLines:            DefaultMaxUiLogLines,
	}
}

// String returns a string representation of the TUIConfig struct.
func (c TUIConfig) String() string {
	output := "TUI Config: {"
	output += fmt.Sprintf("AutoScroll: %v, ", c.AutoScroll)
	output += fmt.Sprintf("WrapLines: %v, ", c.WrapLines)
	output += fmt.Sprintf("Borderless: %v, ", c.Borderless)
	output += fmt.Sprintf("OverrideInstanceName: %s, ", c.OverrideInstanceName)
	output += fmt.Sprintf("CometBFTStartsMinimized: %v, ", c.CometBFTStartsMinimized)
	output += fmt.Sprintf("ConductorStartsMinimized: %v, ", c.ConductorStartsMinimized)
	output += fmt.Sprintf("ComposerStartsMinimized: %v, ", c.ComposerStartsMinimized)
	output += fmt.Sprintf("SequencerStartsMinimized: %v, ", c.SequencerStartsMinimized)
	output += fmt.Sprintf("GenericStartsMinimized: %v", c.GenericStartsMinimized)
	output += fmt.Sprintf("GenericStartPosition: %v", c.GenericStartPosition)
	output += fmt.Sprintf("HighlightColor: %s, ", c.HighlightColor)
	output += fmt.Sprintf("BorderColor: %s, ", c.BorderColor)
	output += fmt.Sprintf("MaxUiLogLines: %d", c.MaxUiLogLines)
	output += "}"
	return output
}

// LoadTUIConfigOrPanic loads the TUIConfigs from the given path.
//
// Panics if the file cannot be loaded or parsed.
func LoadTUIConfigOrPanic(path string) TUIConfig {
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

	// validate the generic start position value
	switch config.GenericStartPosition {
	case "before", "after", "default":
		// valid values; do nothing
	default:
		log.Warnf("Invalid value for generic_start_position: %q. Valid values are: 'before', 'after', 'default'", config.GenericStartPosition)
		log.Warnf("Setting generic_start_position to 'default'")
		config.GenericStartPosition = "default"
	}

	return config
}

// CreateTUIConfig creates a TUI configuration file and populates it
// with the defaults for the devrunner TUI.
//
// Panics if the file cannot be created or written to.
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
