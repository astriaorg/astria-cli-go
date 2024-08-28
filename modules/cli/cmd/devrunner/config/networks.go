package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	log "github.com/sirupsen/logrus"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// NetworkConfigs is the struct that holds the configuration for all individual Astria networks.
type NetworkConfigs struct {
	Configs map[string]NetworkConfig `mapstructure:"networks" toml:"networks"`
}

type ServiceConfig struct {
	Name        string   `mapstructure:"name" toml:"name"`
	Version     string   `mapstructure:"version" toml:"version"`
	DownloadURL string   `mapstructure:"download_url" toml:"download_url"`
	LocalPath   string   `mapstructure:"local_path" toml:"local_path"`
	Args        []string `mapstructure:"args" toml:"args"`
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
				RollupName:       "astria-test-chain-1",
				NativeDenom:      DefaultLocalNativeDenom,
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + AstriaConductorVersion,
						DownloadURL: KnownBinaries.AstriaConductor.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: KnownBinaries.AstriaComposer.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
						Args:        nil,
					},
					"sequencer": {
						Name:        "astria-sequencer",
						Version:     "v" + AstriaSequencerVersion,
						DownloadURL: KnownBinaries.AstriaSequencer.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-sequencer-v"+AstriaSequencerVersion),
						Args:        nil,
					},
					"cometbft": {
						Name:        "cometbft",
						Version:     "v" + CometbftVersion,
						DownloadURL: KnownBinaries.CometBFT.Url,
						LocalPath:   filepath.Join(defaultBinDir, "cometbft-v"+CometbftVersion),
						Args:        nil,
					},
				},
			},
			"dusk": {
				SequencerChainId: "astria-dusk-" + cmd.DefaultDuskNum,
				SequencerGRPC:    "https://grpc.sequencer.dusk-" + cmd.DefaultDuskNum + ".devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dusk-" + cmd.DefaultDuskNum + ".devnet.astria.org/",
				RollupName:       "",
				NativeDenom:      DefaultLocalNativeDenom,
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + AstriaConductorVersion,
						DownloadURL: KnownBinaries.AstriaConductor.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: KnownBinaries.AstriaComposer.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
						Args:        nil,
					},
				},
			},
			"dawn": {
				SequencerChainId: "astria-dawn-" + cmd.DefaultDawnNum,
				SequencerGRPC:    "https://grpc.sequencer.dawn-" + cmd.DefaultDawnNum + ".devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dawn-" + cmd.DefaultDawnNum + ".devnet.astria.org/",
				RollupName:       "",
				NativeDenom:      "ibc/channel0/utia",
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + AstriaConductorVersion,
						DownloadURL: KnownBinaries.AstriaConductor.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: KnownBinaries.AstriaComposer.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
						Args:        nil,
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
						DownloadURL: KnownBinaries.AstriaConductor.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+AstriaConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + AstriaComposerVersion,
						DownloadURL: KnownBinaries.AstriaComposer.Url,
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+AstriaComposerVersion),
						Args:        nil,
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

// CreateNetworksConfig creates a networks configuration file and populates it
// with the network defaults. The binPath is required to accommodate which CLI
// instance this particular networks config is for and to build the proper paths
// to the binaries that will be used for the given instance. The configPath is
// provided for the same reason; which instance is this config file for and
// where to put it. This function will also override the default local denom and local
// sequencer network chain id based on the command line flags provided. It will
// skip initialization if the file already exists. It will panic if the file
// cannot be created or written to.
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
	rollupEndpoint, exists := bc["astria_composer_rollups"]
	if !exists {
		log.Error("ASTRIA_COMPOSER_ROLLUPS not found in BaseConfig")
		panic(fmt.Errorf("ASTRIA_COMPOSER_ROLLUPS not found in BaseConfig"))
	}
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
