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

// BuildSequencerNetworkConfigsFilepath returns the path to the sequencer
// networks configuration file. The file is located in the user's home directory
// in the default Astria config directory (~/.astria/).
func BuildSequencerNetworkConfigsFilepath() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	return filepath.Join(homePath, DefaultConfigDirName, DefaultSequencerNetworksConfigFilename)
}

// CreateSequencerNetworkConfigs creates a sequencer networks configuration file at the
// given path. It will skip initialization if the file already exists. It will
// panic if the file cannot be created or if there is an error encoding the
// NetworksConfigs struct to a file.
func CreateSequencerNetworkConfigs(path string) {
	_, err := os.Stat(path)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", path)
		return
	}

	config := DefaultNetworksConfigs()

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// encode the struct to TOML and write to the file
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		panic(err)
	}
	log.Infof("New network config file created successfully: %s\n", path)
}

// GetSequencerNetworkSettingsFromConfig returns the SequencerNetworkConfig for
// the given network. The function automatically checks if the sequencer network
// config file exists, creates it if it does not, and then loads the config.
// It will panic if the network is not one of 'local', 'dusk', 'dawn', or
// 'mainnet', or if the config file cannot be created or loaded.
func GetSequencerNetworkSettingsFromConfig(network, path string) SequencerNetworkConfig {
	sequencerConfig := LoadSequencerNetworkConfigsOrPanic(path)

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

// ChooseFlagValue returns the value of the flag based on the usage of the
// specified flag and the usage of the network config flag.
func ChooseFlagValue(networksChange bool, flagChange bool, configValue string, flagValue string) string {
	// There are four possible scenarios for setting a flag's value:
	// 1. network flag hasn't changed & flag hasn't changed
	//    -> return the flag default value
	// 2. network flag hasn't changed & flag has changed
	//    -> return the flag value
	// 3. network flag has changed & flag hasn't changed
	//    -> return the network config value
	// 4. network flag has changed & flag has changed
	//    -> return the flag value
	//
	// Using Cobra, situations 1, 2, and 4 are already handled.
	// If situation 3 occurs, the network config value needs to be handled
	// specifically.
	// The logic below will return the config value if situation 3 occurs,
	// otherwise it will return the flag value output by Cobra.
	if networksChange && !flagChange {
		return configValue
	} else {
		return flagValue
	}
}
