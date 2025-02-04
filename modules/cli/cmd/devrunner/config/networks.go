package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	util "github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/utilities"
	log "github.com/sirupsen/logrus"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// NetworkConfigs is the struct that holds the configuration for all individual Astria networks.
type NetworkConfigs struct {
	Configs map[string]NetworkConfig `mapstructure:"networks" toml:"networks"`
}

// Expand shell expands all the fields in the NetworkConfigs struct.
func (n NetworkConfigs) Expand() NetworkConfigs {
	for networkName, networkConfig := range n.Configs {
		n.Configs[networkName] = networkConfig.Expand()
	}

	return n
}

type ServiceConfig struct {
	Name        string   `mapstructure:"name" toml:"name"`
	Version     string   `mapstructure:"version" toml:"version"`
	DownloadURL string   `mapstructure:"download_url" toml:"download_url"`
	LocalPath   string   `mapstructure:"local_path" toml:"local_path"`
	Args        []string `mapstructure:"args" toml:"args"`
}

// Expand shell expands all the fields in the ServiceConfig struct.
func (s ServiceConfig) Expand() ServiceConfig {
	s.Name = util.ShellExpand(s.Name)
	s.Version = util.ShellExpand(s.Version)
	s.DownloadURL = util.ShellExpand(s.DownloadURL)
	s.LocalPath = util.ShellExpand(s.LocalPath)

	for i, arg := range s.Args {
		s.Args[i] = util.ShellExpand(arg)
	}

	return s
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

// Expand shell expands all the fields in the NetworkConfig struct.
func (n NetworkConfig) Expand() NetworkConfig {
	n.SequencerChainId = util.ShellExpand(n.SequencerChainId)
	n.SequencerGRPC = util.ShellExpand(n.SequencerGRPC)
	n.SequencerRPC = util.ShellExpand(n.SequencerRPC)
	n.RollupName = util.ShellExpand(n.RollupName)
	n.NativeDenom = util.ShellExpand(n.NativeDenom)

	for serviceName, serviceConfig := range n.Services {
		n.Services[serviceName] = serviceConfig.Expand()
	}

	return n
}

// DefaultNetworksConfigs returns a NetworksConfig struct populated with all
// network defaults.
func DefaultNetworksConfigs(defaultBinDir, sequencerNetworkName, rollupName, defaultDenom string) NetworkConfigs {
	return NetworkConfigs{
		Configs: map[string]NetworkConfig{
			"local": {
				SequencerChainId: sequencerNetworkName,
				SequencerGRPC:    "http://127.0.0.1:8080",
				SequencerRPC:     "http://127.0.0.1:26657",
				RollupName:       rollupName,
				NativeDenom:      defaultDenom,
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + MainnetAstriaConductorVersion,
						DownloadURL: ServiceUrls.AstriaConductorReleaseUrl(MainnetAstriaConductorVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+MainnetAstriaConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + MainnetAstriaComposerVersion,
						DownloadURL: ServiceUrls.AstriaComposerReleaseUrl(MainnetAstriaComposerVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+MainnetAstriaComposerVersion),
						Args:        nil,
					},
					"sequencer": {
						Name:        "astria-sequencer",
						Version:     "v" + MainnetAstriaSequencerVersion,
						DownloadURL: ServiceUrls.AstriaSequencerReleaseUrl(MainnetAstriaSequencerVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-sequencer-v"+MainnetAstriaSequencerVersion),
						Args:        nil,
					},
					"cometbft": {
						Name:        "cometbft",
						Version:     "v" + MainnetCometbftVersion,
						DownloadURL: ServiceUrls.CometBftReleaseUrl(MainnetCometbftVersion),
						LocalPath:   filepath.Join(defaultBinDir, "cometbft-v"+MainnetCometbftVersion),
						Args:        nil,
					},
				},
			},
			"dusk": {
				SequencerChainId: "astria-dusk-" + cmd.DefaultDuskNum,
				SequencerGRPC:    "https://grpc.sequencer.dusk-" + cmd.DefaultDuskNum + ".devnet.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dusk-" + cmd.DefaultDuskNum + ".devnet.astria.org/",
				RollupName:       rollupName,
				NativeDenom:      defaultDenom,
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + DevnetConductorVersion,
						DownloadURL: ServiceUrls.AstriaConductorReleaseUrl(DevnetConductorVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+DevnetConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + DevnetComposerVersion,
						DownloadURL: ServiceUrls.AstriaComposerReleaseUrl(DevnetComposerVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+DevnetComposerVersion),
						Args:        nil,
					},
				},
			},
			"dawn": {
				SequencerChainId: "dawn-" + cmd.DefaultDawnNum,
				SequencerGRPC:    "https://grpc.sequencer.dawn-" + cmd.DefaultDawnNum + ".astria.org/",
				SequencerRPC:     "https://rpc.sequencer.dawn-" + cmd.DefaultDawnNum + ".astria.org/",
				RollupName:       rollupName,
				NativeDenom:      "ibc/channel-0/utia",
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + TestnetConductorVersion,
						DownloadURL: ServiceUrls.AstriaConductorReleaseUrl(TestnetConductorVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+TestnetConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + TestnetComposerVersion,
						DownloadURL: ServiceUrls.AstriaComposerReleaseUrl(TestnetComposerVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+TestnetComposerVersion),
						Args:        nil,
					},
				},
			},
			"mainnet": {
				SequencerChainId: "astria",
				SequencerGRPC:    "https://grpc.sequencer.astria.org/",
				SequencerRPC:     "https://rpc.sequencer.astria.org/",
				RollupName:       rollupName,
				NativeDenom:      "ibc/channel-0/utia",
				Services: map[string]ServiceConfig{
					"conductor": {
						Name:        "astria-conductor",
						Version:     "v" + MainnetAstriaConductorVersion,
						DownloadURL: ServiceUrls.AstriaConductorReleaseUrl(MainnetAstriaConductorVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-conductor-v"+MainnetAstriaConductorVersion),
						Args:        nil,
					},
					"composer": {
						Name:        "astria-composer",
						Version:     "v" + MainnetAstriaComposerVersion,
						DownloadURL: ServiceUrls.AstriaComposerReleaseUrl(MainnetAstriaComposerVersion),
						LocalPath:   filepath.Join(defaultBinDir, "astria-composer-v"+MainnetAstriaComposerVersion),
						Args:        nil,
					},
				},
			},
		},
	}
}

// LoadNetworkConfigsOrPanic loads the NetworksConfig from the given path.
//
// Panics if the file cannot be loaded or parsed.
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

	// shell expand all the fields in the config
	config = config.Expand()

	return config
}

// CreateNetworksConfig creates a networks configuration file and populates it
// with the network defaults.
//   - configPath: the path to the networks configuration file
//   - binPath: the path to the binaries directory within a given instance
//   - localSequencerChainId: the chain id for the local sequencer
//   - rollupName: the name of the rollup
//   - localNativeDenom: the native denom for the local sequencer
//
// Note: The configPath and binPath should be part of the same instance.
//
// This function will also override the default local denom and local
// sequencer network chain id based on the command line flags provided. It will
// skip initialization if the file already exists.
//
// Panic if the file cannot be created or written to.
func CreateNetworksConfig(configPath, binPath, localSequencerChainId, rollupName, localNativeDenom string) {
	_, err := os.Stat(configPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", configPath)
		return
	}
	// create an instance of the Config struct with some data
	config := DefaultNetworksConfigs(binPath, localSequencerChainId, rollupName, localNativeDenom)

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
//
// The overrides this function returns are:
//   - ASTRIA_CONDUCTOR_SEQUENCER_GRPC_URL
//   - ASTRIA_CONDUCTOR_SEQUENCER_COMETBFT_URL
//   - ASTRIA_CONDUCTOR_EXPECTED_SEQUENCER_CHAIN_ID
//   - ASTRIA_COMPOSER_SEQUENCER_CHAIN_ID
//   - ASTRIA_COMPOSER_SEQUENCER_ABCI_ENDPOINT
//   - ASTRIA_COMPOSER_SEQUENCER_GRPC_ENDPOINT
//   - ASTRIA_COMPOSER_ROLLUPS
//
// Panics if the ASTRIA_COMPOSER_ROLLUPS env var is not found.
func (n NetworkConfig) GetEndpointOverrides(bc BaseConfig) []string {
	rollups, exists := bc["astria_composer_rollups"]
	if !exists {
		log.Error("ASTRIA_COMPOSER_ROLLUPS not found in BaseConfig")
		panic(fmt.Errorf("ASTRIA_COMPOSER_ROLLUPS not found in BaseConfig"))
	}

	return []string{
		"ASTRIA_CONDUCTOR_SEQUENCER_GRPC_URL=" + n.SequencerGRPC,
		"ASTRIA_CONDUCTOR_SEQUENCER_COMETBFT_URL=" + n.SequencerRPC,
		"ASTRIA_CONDUCTOR_EXPECTED_SEQUENCER_CHAIN_ID" + "=" + n.SequencerChainId,
		"ASTRIA_COMPOSER_SEQUENCER_CHAIN_ID=" + n.SequencerChainId,
		"ASTRIA_COMPOSER_SEQUENCER_ABCI_ENDPOINT=" + n.SequencerRPC,
		"ASTRIA_COMPOSER_SEQUENCER_GRPC_ENDPOINT=" + n.SequencerGRPC,
		"ASTRIA_COMPOSER_ROLLUPS=" + rollups,
	}
}
