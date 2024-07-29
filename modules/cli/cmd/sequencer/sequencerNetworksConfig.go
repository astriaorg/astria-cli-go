package sequencer

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// SequencerNetworkConfigs is the struct that holds the configurations for all
// individual Astria networks.
type SequencerNetworkConfigs struct {
	Local   SequencerNetworkConfig `mapstructure:"local" toml:"local"`
	Dusk    SequencerNetworkConfig `mapstructure:"dusk" toml:"dusk"`
	Dawn    SequencerNetworkConfig `mapstructure:"dawn" toml:"dawn"`
	Mainnet SequencerNetworkConfig `mapstructure:"mainnet" toml:"mainnet"`
}

// SequencerNetworkConfig is the struct that holds the configuration for
// interacting with a given Astria sequencer network.
type SequencerNetworkConfig struct {
	SequencerChainId string `mapstructure:"sequencer_chain_id" toml:"sequencer_chain_id"`
	SequencerURL     string `mapstructure:"sequencer_url" toml:"sequencer_url"`
	Asset            string `mapstructure:"asset" toml:"asset"`
	FeeAsset         string `mapstructure:"fee_asset" toml:"fee_asset"`
}

// DefaultNetworksConfigs returns a SequencerNetworkConfigs struct populated
// with all sequencer network defaults.
func DefaultNetworksConfigs() SequencerNetworkConfigs {
	config := SequencerNetworkConfigs{
		Local: SequencerNetworkConfig{
			SequencerChainId: "sequencer-test-chain-0",
			SequencerURL:     "http://127.0.0.1:26657",
			Asset:            "nria",
			FeeAsset:         "nria",
		},
		Dusk: SequencerNetworkConfig{
			SequencerChainId: DefaultDuskSequencerChainID,
			SequencerURL:     DefaultDuskSequencerURL,
			Asset:            "nria",
			FeeAsset:         "nria",
		},
		Dawn: SequencerNetworkConfig{
			SequencerChainId: DefaultDawnSequencerChainID,
			SequencerURL:     DefaultDawnSequencerURL,
			Asset:            "ibc/channel0/utia",
			FeeAsset:         "ibc/channel0/utia",
		},
		Mainnet: SequencerNetworkConfig{
			SequencerChainId: DefaultMainnetSequencerChainID,
			SequencerURL:     DefaultMainnetSequencerURL,
			Asset:            "ibc/channel0/utia",
			FeeAsset:         "ibc/channel0/utia",
		},
	}
	return config
}

// LoadSequencerNetworkConfigsOrPanic loads the SequencerNetworkConfigs from the given
// path. If the file cannot be loaded or parsed, the function will panic.
func LoadSequencerNetworkConfigsOrPanic(path string) SequencerNetworkConfigs {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}

	var config SequencerNetworkConfigs
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}

	return config
}

// CreateSequencerNetworkConfigs creates a networks configuration file at the
// given path. It will skip initialization if the file already exists. It will
// return an error if the file cannot be created or written to.
func CreateSequencerNetworkConfigs(path string) error {

	_, err := os.Stat(path)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", path)
		return nil
	}

	config := DefaultNetworksConfigs()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// encode the struct to TOML and write to the file
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}
	log.Infof("New network config file created successfully: %s\n", path)
	return nil
}

// GetSequencerNetworkSettingsFromConfig returns the SequencerNetworkConfig for
// the given network. The function automatically checks if the sequencer network
// config file exists, creates it if it does not, and then loads the config.
// It will panic if the network is not one of 'local', 'dusk', 'dawn', or
// 'mainnet', or if the config file cannot be created or loaded.
func GetSequencerNetworkSettingsFromConfig(network string) SequencerNetworkConfig {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	sequencerConfigPath := filepath.Join(homePath, DefaultConfigDirName, DefaultSequencerNetworksConfigFilename)

	log.Info("Network flag changed")

	err = CreateSequencerNetworkConfigs(sequencerConfigPath)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer networks config file")
		panic(err)
	}

	sequencerConfig := LoadSequencerNetworkConfigsOrPanic(sequencerConfigPath)

	var networkSettings SequencerNetworkConfig
	switch network {
	case "local":
		networkSettings = sequencerConfig.Local
	case "dusk":
		networkSettings = sequencerConfig.Dusk
	case "dawn":
		networkSettings = sequencerConfig.Dawn
	case "mainnet":
		networkSettings = sequencerConfig.Mainnet
	default:
		panic("Invalid network selected: Must be one of 'local', 'dusk', 'dawn', or 'mainnet'.")
	}

	return networkSettings
}

// TODO: Add a description for this function.
func ChooseFlagValue(networksChange bool, flagChange bool, configValue string, flagValue string) string {
	if networksChange && !flagChange {
		return configValue
	} else {
		return flagValue
	}
}
