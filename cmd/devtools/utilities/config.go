package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

const DefualtNetworksConfigName = "networks-config.toml"

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
	DefaultDenom     string `mapstructure:"default_denom" toml:"default_denom"`
	RollupName       string `mapstructure:"rollup_name" toml:"rollup_name"`
}

// defaultNetworksConfig returns a NetworksConfig struct populated with all
// network defaults.
func defaultNetworksConfig() NetworksConfig {
	config := NetworksConfig{
		Networks: Networks{
			Local: Network{
				SequencerChainId: "local-test-sequencer",
				SequencerGRPC:    "http://127.0.0.1:8080",
				SequencerRPC:     "http://127.0.0.1:26657",
				DefaultDenom:     "nria",
				RollupName:       "local-rollup",
			},
			Dusk: Network{
				SequencerChainId: "astria-dusk-5",
				SequencerGRPC:    "https://grpc.sequencer.dusk-5.devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dusk-5.devnet.astria.org/",
				DefaultDenom:     "nria",
				RollupName:       "",
			},
			Dawn: Network{
				SequencerChainId: "astria-dawn-0",
				SequencerGRPC:    "https://grpc.sequencer.dawn-0.devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dawn-0.devnet.astria.org/",
				DefaultDenom:     "ibc/channel0/utia",
				RollupName:       "",
			},
			Mainnet: Network{
				SequencerChainId: "astria",
				SequencerGRPC:    "https://grpc.sequencer.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.astria.org/",
				DefaultDenom:     "ibc/channel0/utia",
				RollupName:       "",
			},
		},
	}
	return config
}

// LoadNetworksConfig loads the NetworksConfig from the given path.
func LoadNetworksConfig(path string) NetworksConfig {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var config NetworksConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Marshal into JSON with indentation for printing
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Print the JSON
	fmt.Println("Loaded Networks Configuration:")
	fmt.Println(string(jsonData))

	return config
}

// CreateDefaultNetworksConfig creates a default networks configuration file at
// the given path, populating the file with the network defaults.
func CreateDefaultNetworksConfig(path string) {
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
}
