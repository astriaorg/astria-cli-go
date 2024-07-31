package sequencer

import (
	"os"
	"path/filepath"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	log "github.com/sirupsen/logrus"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// SequencerNetworkConfig is the struct that holds the configuration for
// interacting with a given Astria sequencer network.
type SequencerNetworkConfig struct {
	SequencerChainId string `mapstructure:"sequencer_chain_id" toml:"sequencer_chain_id"`
	SequencerURL     string `mapstructure:"sequencer_url" toml:"sequencer_url"`
	Asset            string `mapstructure:"asset" toml:"asset"`
	FeeAsset         string `mapstructure:"fee_asset" toml:"fee_asset"`
}

// SequencerNetworkConfigs is a map of SequencerNetworkConfig structs.
// Using a map here to allow for user to add new networks by adding a new section in the toml
type SequencerNetworkConfigs struct {
	Configs map[string]SequencerNetworkConfig `mapstructure:"networks" toml:"networks"`
}

// GetSequencerNetworkConfigsPresets returns a map of all sequencer network presets.
// Used to generate the initial config file.
func GetSequencerNetworkConfigsPresets() SequencerNetworkConfigs {
	return SequencerNetworkConfigs{
		Configs: map[string]SequencerNetworkConfig{
			"Local": {
				SequencerChainId: "sequencer-test-chain-0",
				SequencerURL:     "http://127.0.0.1:26657",
				Asset:            "nria",
				FeeAsset:         "nria",
			},
			"Dusk": {
				SequencerChainId: DefaultSequencerChainID,
				SequencerURL:     DefaultSequencerURL,
				Asset:            "nria",
				FeeAsset:         "nria",
			},
			"Dawn": {
				SequencerChainId: "astria-dawn-0",
				SequencerURL:     "https://rpc.sequencer.dawn-0.devnet.astria.org",
				Asset:            "ibc/channel0/utia",
				FeeAsset:         "ibc/channel0/utia",
			},
			"Mainnet": {
				SequencerChainId: "astria",
				SequencerURL:     "https://rpc.sequencer.astria.org/",
				Asset:            "ibc/channel0/utia",
				FeeAsset:         "ibc/channel0/utia",
			},
		},
	}
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
		log.Fatalf("Unable to decode toml into struct, %v", err)
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
		log.Debugf("%s already exists. Skipping initialization.\n", path)
		return
	}

	config := GetSequencerNetworkConfigsPresets()

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
func GetSequencerNetworkSettingsFromConfig(network, path string) SequencerNetworkConfig {
	sequencerConfig := LoadSequencerNetworkConfigsOrPanic(path)

	if _, ok := sequencerConfig.Configs[network]; !ok {
		log.Fatalf("Network %s not found in config file at %s", network, path)
		panic("Network not found in config file")
	}

	return sequencerConfig.Configs[network]
}

// GetNetworkConfigFromFlags returns a SequencerNetworkConfig based on the
// network flag value. It will create the network config file if it does not
// exist, and then load the config.
func GetNetworkConfigFromFlags(flagHandler *cmd.CliFlagHandler) SequencerNetworkConfig {
	network := flagHandler.GetValue("network")
	networksConfigPath := BuildSequencerNetworkConfigsFilepath()
	CreateSequencerNetworkConfigs(networksConfigPath)
	networkSettings := GetSequencerNetworkSettingsFromConfig(network, networksConfigPath)

	return networkSettings
}

// ChooseFlagValue returns the value of the flag based on the usage of the
// specified flag and the usage of the network config flag.
// TODO - delete after all the commands are refactored
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
