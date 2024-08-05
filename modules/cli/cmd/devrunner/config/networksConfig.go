package config

import (
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// NetworkConfigs is the struct that holds the configuration for all individual Astria networks.
type NetworkConfigs struct {
	Configs map[string]NetworkConfig `mapstructure:"networks" toml:"networks"`
}

type ServiceConfig struct {
	Name        string `mapstructure:"name" toml:"name"`
	Version     string `mapstructure:"version" toml:"version"`
	DownloadURL string `mapstructure:"download_url" toml:"download_url"`
	LocalPath   string `mapstructure:"local_path" toml:"local_path"`
	// TODO: implement generic args?
}

// NetworkConfig is the struct that holds the configuration for an individual Astria network.
type NetworkConfig struct {
	SequencerChainId string                   `mapstructure:"sequencer_chain_id" toml:"sequencer_chain_id"`
	SequencerGRPC    string                   `mapstructure:"sequencer_grpc" toml:"sequencer_grpc"`
	SequencerRPC     string                   `mapstructure:"sequencer_rpc" toml:"sequencer_rpc"`
	RollupName       string                   `mapstructure:"rollup_name" toml:"rollup_name"`
	NativeDenom      string                   `mapstructure:"default_denom" toml:"default_denom"`
	Services         map[string]ServiceConfig `mapstructure:"services" toml:"services"`
}

// DefaultNetworksConfigs returns a NetworksConfig struct populated with all
// network defaults.
func DefaultNetworksConfigs(defaultBinDir string) NetworkConfigs {
	return NetworkConfigs{
		Configs: map[string]NetworkConfig{
			"local": {
				SequencerChainId: "sequencer-test-chain-0",
				SequencerGRPC:    "http://127.0.0.1:8080",
				SequencerRPC:     "http://127.0.0.1:26657",
				RollupName:       "astria-test-chain",
				NativeDenom:      "nria",
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + AstriaConductorVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-conductor").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-composer").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
					},
					"sequencer": {
						Name:        "astria-sequencer",
						Version:     "v" + AstriaSequencerVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-sequencer").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-sequencer-v"+AstriaSequencerVersion),
					},
					"cometbft": {
						Name:        "cometbft",
						Version:     "v" + CometbftVersion,
						DownloadURL: findBinaryByName(Binaries, "cometbft").Url,
						LocalPath:   filepath.Join(defaultBinDir, "cometbft-v"+CometbftVersion),
					},
				},
			},
			"dusk": {
				SequencerChainId: "astria-dusk-" + duskNum,
				SequencerGRPC:    "https://grpc.sequencer.dusk-" + duskNum + ".devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dusk-" + duskNum + ".devnet.astria.org/",
				RollupName:       "",
				NativeDenom:      "nria",
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + AstriaConductorVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-conductor").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-composer").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
					},
				},
			},
			"dawn": {
				SequencerChainId: "astria-dawn-" + dawnNum,
				SequencerGRPC:    "https://grpc.sequencer.dawn-" + dawnNum + ".devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dawn-" + dawnNum + ".devnet.astria.org/",
				RollupName:       "",
				NativeDenom:      "ibc/channel0/utia",
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + AstriaConductorVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-conductor").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-composer").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
					},
				},
			},
			"mainnet": {
				SequencerChainId: "astria",
				SequencerGRPC:    "https://grpc.sequencer.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.astria.org/",
				RollupName:       "",
				NativeDenom:      "ibc/channel0/utia",
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + AstriaConductorVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-conductor").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: findBinaryByName(Binaries, "astria-composer").Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
					},
				},
			},
		},
	}
}

// LoadNetworkConfigsOrPanic loads the NetworksConfig from the given path. If the file
// cannot be loaded or parsed, the function will panic.
func LoadNetworkConfigsOrPanic(path string) NetworkConfigs {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}

	var config NetworkConfigs
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}

	return config
}

// CreateNetworksConfig creates a networks configuration file at
// the given path, populating the file with the network defaults, and overriding
// the default local denom and local sequencer network chain id.
// It will skip initialization if the file already exists. It will panic if the
// file cannot be created or written to.
func CreateNetworksConfig(binPath, configPath, localSequencerChainId, localNativeDenom string) {

	_, err := os.Stat(configPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", configPath)
		return
	}
	// create an instance of the Config struct with some data
	config := DefaultNetworksConfigs(binPath)
	local := config.Configs["local"]
	local.NativeDenom = localNativeDenom
	local.SequencerChainId = localSequencerChainId
	config.Configs["local"] = local

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
	log.Infof("New network config file created successfully: %s\n", configPath)
}

// GetEndpointOverrides returns a slice of environment variables for supporting
// the ability to run against different Sequencer networks. It enables a way to
// dynamically configure endpoints for Conductor and Composer to override
// the default environment variables for the network configuration. It uses the
// BaseConfig to properly update the ASTRIA_COMPOSER_ROLLUPS env var.
func (n NetworkConfig) GetEndpointOverrides(bc BaseConfig) []string {
	rollupEndpoint := bc.Astria_composer_rollups
	// get the rollup ws endpoint
	pattern := `ws{1,2}:\/\/.*:\d+`
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Error("Error compiling regex")
		panic(err)
	}
	match := re.FindString(rollupEndpoint)

	return []string{
		"ASTRIA_COMPOSER_SEQUENCER_CHAIN_ID=" + n.SequencerChainId,
		"ASTRIA_CONDUCTOR_SEQUENCER_GRPC_URL=" + n.SequencerGRPC,
		"ASTRIA_CONDUCTOR_SEQUENCER_COMETBFT_URL=" + n.SequencerRPC,
		"ASTRIA_COMPOSER_SEQUENCER_URL=" + n.SequencerRPC,
		"ASTRIA_COMPOSER_ROLLUPS=" + n.RollupName + "::" + match,
	}
}
