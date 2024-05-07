package utilities

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// NetworksConfig is the struct that holds all Astria network configurations.
type NetworksConfig struct {
	Networks Networks `mapstructure:"networks" toml:"networks"`
}

// Networks is the struct that holds the configuration for all individual Astria networks.
type Networks struct {
	Local   Network `mapstructure:"local" toml:"local"`
	Dusk    Network `mapstructure:"dusk" toml:"dusk"`
	Dawn    Network `mapstructure:"dawn" toml:"dawn"`
	Mainnet Network `mapstructure:"mainnet" toml:"mainnet"`
}

// Network is the struct that holds the configuration for an individual Astria network.
type Network struct {
	SequencerChainId string `mapstructure:"sequencer_chain_id" toml:"sequencer_chain_id"`
	SequencerGRPC    string `mapstructure:"sequencer_grpc" toml:"sequencer_grpc"`
	SequencerRPC     string `mapstructure:"sequencer_rpc" toml:"sequencer_rpc"`
	RollupName       string `mapstructure:"rollup_name" toml:"rollup_name"`
}

// defaultNetworksConfig returns a NetworksConfig struct populated with all
// network defaults.
func defaultNetworksConfig() NetworksConfig {
	config := NetworksConfig{
		Networks: Networks{
			Local: Network{
				SequencerChainId: "sequencer-test-chain-0",
				SequencerGRPC:    "http://127.0.0.1:8080",
				SequencerRPC:     "http://127.0.0.1:26657",
				RollupName:       "local-rollup",
			},
			Dusk: Network{
				SequencerChainId: "astria-dusk-5",
				SequencerGRPC:    "https://grpc.sequencer.dusk-5.devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dusk-5.devnet.astria.org/",
				RollupName:       "",
			},
			Dawn: Network{
				SequencerChainId: "astria-dawn-0",
				SequencerGRPC:    "https://grpc.sequencer.dawn-0.devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dawn-0.devnet.astria.org/",
				RollupName:       "",
			},
			Mainnet: Network{
				SequencerChainId: "astria",
				SequencerGRPC:    "https://grpc.sequencer.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.astria.org/",
				RollupName:       "",
			},
		},
	}
	return config
}

// LoadNetworksConfig loads the NetworksConfig from the given path. If the file
// cannot be loaded or parsed, the function will panic.
func LoadNetworksConfig(path string) NetworksConfig {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}

	var config NetworksConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}

	// Marshal into JSON for printing
	jsonData, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
		panic(err)
	}

	// Print the JSON
	log.Debugf("Loaded Networks Configuration: %s", string(jsonData))

	return config
}

// CreateDefaultNetworksConfig creates a default networks configuration file at
// the given path, populating the file with the network defaults. It will skip
// initialization if the file already exists. It will panic if the file cannot
// be created or written to.
func CreateDefaultNetworksConfig(path string) {

	_, err := os.Stat(path)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", path)
		return
	}
	// Create an instance of the Config struct with some data
	config := defaultNetworksConfig()

	// Open a file for writing
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the struct to TOML and write to the file
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		panic(err)
	}
	log.Infof("New network config file created successfully: %s\n", path)

}

// GetEnvOverrides returns a slice of environment variables that can be used to
// override the default environment variables for the network configuration.
func (n Network) GetEnvOverrides() []string {
	return []string{
		"ASTRIA_COMPOSER_SEQUENCER_CHAIN_ID=" + n.SequencerChainId,
		"ASTRIA_CONDUCTOR_SEQUENCER_GRPC_URL=" + n.SequencerGRPC,
		"ASTRIA_CONDUCTOR_SEQUENCER_COMETBFT_URL=" + n.SequencerRPC,
		"ASTRIA_COMPOSER_SEQUENCER_URL=" + n.SequencerRPC,
		"ASTRIA_COMPOSER_ROLLUPS=" + n.RollupName + "::ws://127.0.0.1:8546",
	}
}
