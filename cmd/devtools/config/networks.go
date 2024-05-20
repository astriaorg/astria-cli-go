package config

import (
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	// conductor
	Astria_conductor_celestia_block_time_ms     int    `mapstructure:"astria_conductor_celestia_block_time_ms" toml:"astria_conductor_celestia_block_time_ms"`
	Astria_conductor_celestia_bearer_token      string `mapstructure:"astria_conductor_celestia_bearer_token" toml:"astria_conductor_celestia_bearer_token"`
	Astria_conductor_celestia_node_http_url     string `mapstructure:"astria_conductor_celestia_node_http_url" toml:"astria_conductor_celestia_node_http_url"`
	Astria_conductor_execution_rpc_url          string `mapstructure:"astria_conductor_execution_rpc_url" toml:"astria_conductor_execution_rpc_url"`
	Astria_conductor_execution_commit_level     string `mapstructure:"astria_conductor_execution_commit_level" toml:"astria_conductor_execution_commit_level"`
	Astria_conductor_log                        string `mapstructure:"astria_conductor_log" toml:"astria_conductor_log"`
	Astria_conductor_no_otel                    bool   `mapstructure:"astria_conductor_no_otel" toml:"astria_conductor_no_otel"`
	Astria_conductor_force_stdout               bool   `mapstructure:"astria_conductor_force_stdout" toml:"astria_conductor_force_stdout"`
	Astria_conductor_pretty_print               bool   `mapstructure:"astria_conductor_pretty_print" toml:"astria_conductor_pretty_print"`
	Astria_conductor_sequencer_grpc_url         string `mapstructure:"astria_conductor_sequencer_grpc_url" toml:"astria_conductor_sequencer_grpc_url"`
	Astria_conductor_sequencer_cometbft_url     string `mapstructure:"astria_conductor_sequencer_cometbft_url" toml:"astria_conductor_sequencer_cometbft_url"`
	Astria_conductor_sequencer_block_time_ms    int    `mapstructure:"astria_conductor_sequencer_block_time_ms" toml:"astria_conductor_sequencer_block_time_ms"`
	Astria_conductor_no_metrics                 bool   `mapstructure:"astria_conductor_no_metrics" toml:"astria_conductor_no_metrics"`
	Astria_conductor_metrics_http_listener_addr string `mapstructure:"astria_conductor_metrics_http_listener_addr" toml:"astria_conductor_metrics_http_listener_addr"`

	// sequencer
	Astria_sequencer_listen_addr                string `mapstructure:"astria_sequencer_listen_addr" toml:"astria_sequencer_listen_addr"`
	Astria_sequencer_db_filepath                string `mapstructure:"astria_sequencer_db_filepath" toml:"astria_sequencer_db_filepath"`
	Astria_sequencer_enable_mint                bool   `mapstructure:"astria_sequencer_enable_mint" toml:"astria_sequencer_enable_mint"`
	Astria_sequencer_grpc_addr                  string `mapstructure:"astria_sequencer_grpc_addr" toml:"astria_sequencer_grpc_addr"`
	Astria_sequencer_log                        string `mapstructure:"astria_sequencer_log" toml:"astria_sequencer_log"`
	Astria_sequencer_no_otel                    bool   `mapstructure:"astria_sequencer_no_otel" toml:"astria_sequencer_no_otel"`
	Astria_sequencer_force_stdout               bool   `mapstructure:"astria_sequencer_force_stdout" toml:"astria_sequencer_force_stdout"`
	Astria_sequencer_no_metrics                 bool   `mapstructure:"astria_sequencer_no_metrics" toml:"astria_sequencer_no_metrics"`
	Astria_sequencer_metrics_http_listener_addr string `mapstructure:"astria_sequencer_metrics_http_listener_addr" toml:"astria_sequencer_metrics_http_listener_addr"`
	Astria_sequencer_pretty_print               bool   `mapstructure:"astria_sequencer_pretty_print" toml:"astria_sequencer_pretty_print"`

	// composer
	Astria_composer_log                        string `mapstructure:"astria_composer_log" toml:"astria_composer_log"`
	Astria_composer_no_otel                    bool   `mapstructure:"astria_composer_no_otel" toml:"astria_composer_no_otel"`
	Astria_composer_force_stdout               bool   `mapstructure:"astria_composer_force_stdout" toml:"astria_composer_force_stdout"`
	Astria_composer_pretty_print               bool   `mapstructure:"astria_composer_pretty_print" toml:"astria_composer_pretty_print"`
	Astria_composer_api_listen_addr            string `mapstructure:"astria_composer_api_listen_addr" toml:"astria_composer_api_listen_addr"`
	Astria_composer_sequencer_url              string `mapstructure:"astria_composer_sequencer_url" toml:"astria_composer_sequencer_url"`
	Astria_composer_sequencer_chain_id         string `mapstructure:"astria_composer_sequencer_chain_id" toml:"astria_composer_sequencer_chain_id"`
	Astria_composer_rollups                    string `mapstructure:"astria_composer_rollups" toml:"astria_composer_rollups"`
	Astria_composer_private_key                string `mapstructure:"astria_composer_private_key" toml:"astria_composer_private_key"`
	Astria_composer_max_submit_interval_ms     int    `mapstructure:"astria_composer_max_submit_interval_ms" toml:"astria_composer_max_submit_interval_ms"`
	Astria_composer_max_bytes_per_bundle       int    `mapstructure:"astria_composer_max_bytes_per_bundle" toml:"astria_composer_max_bytes_per_bundle"`
	Astria_composer_bundle_queue_capacity      int    `mapstructure:"astria_composer_bundle_queue_capacity" toml:"astria_composer_bundle_queue_capacity"`
	Astria_composer_no_metrics                 bool   `mapstructure:"astria_composer_no_metrics" toml:"astria_composer_no_metrics"`
	Astria_composer_metrics_http_listener_addr string `mapstructure:"astria_composer_metrics_http_listener_addr" toml:"astria_composer_metrics_http_listener_addr"`
	Astria_composer_grpc_addr                  string `mapstructure:"astria_composer_grpc_addr" toml:"astria_composer_grpc_addr"`

	// global
	No_color string `mapstructure:"no_color" toml:"no_color"`

	// otel
	Otel_exporter_otlp_endpoint           string `mapstructure:"otel_exporter_otlp_endpoint" toml:"otel_exporter_otlp_endpoint"`
	Otel_exporter_otlp_traces_endpoint    string `mapstructure:"otel_exporter_otlp_traces_endpoint" toml:"otel_exporter_otlp_traces_endpoint"`
	Otel_exporter_otlp_traces_timeout     int    `mapstructure:"otel_exporter_otlp_traces_timeout" toml:"otel_exporter_otlp_traces_timeout"`
	Otel_exporter_otlp_traces_compression string `mapstructure:"otel_exporter_otlp_traces_compression" toml:"otel_exporter_otlp_traces_compression"`
	Otel_exporter_otlp_headers            string `mapstructure:"otel_exporter_otlp_headers" toml:"otel_exporter_otlp_headers"`
	Otel_exporter_otlp_trace_headers      string `mapstructure:"otel_exporter_otlp_trace_headers" toml:"otel_exporter_otlp_trace_headers"`
}

func NewBaseConfig(instanceName string) BaseConfig {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("Error getting home dir")
		panic(err)
	}
	return BaseConfig{
		Astria_conductor_celestia_block_time_ms:     1200,
		Astria_conductor_celestia_bearer_token:      "<JWT Bearer token>",
		Astria_conductor_celestia_node_http_url:     "http://127.0.0.1:26658",
		Astria_conductor_execution_rpc_url:          "http://127.0.0.1:50051",
		Astria_conductor_execution_commit_level:     "SoftOnly",
		Astria_conductor_log:                        "astria_conductor=info",
		Astria_conductor_no_otel:                    true,
		Astria_conductor_force_stdout:               true,
		Astria_conductor_pretty_print:               true,
		Astria_conductor_sequencer_grpc_url:         "http://127.0.0.1:8080",
		Astria_conductor_sequencer_cometbft_url:     "http://127.0.0.1:26657",
		Astria_conductor_sequencer_block_time_ms:    2000,
		Astria_conductor_no_metrics:                 true,
		Astria_conductor_metrics_http_listener_addr: "127.0.0.1:9000",

		Astria_sequencer_listen_addr:                "127.0.0.1:26658",
		Astria_sequencer_db_filepath:                filepath.Join(homePath, ".astria", instanceName, DataDirName, "astria_sequencer_db"),
		Astria_sequencer_enable_mint:                false,
		Astria_sequencer_grpc_addr:                  "127.0.0.1:8080",
		Astria_sequencer_log:                        "astria_sequencer=info",
		Astria_sequencer_no_otel:                    true,
		Astria_sequencer_force_stdout:               true,
		Astria_sequencer_no_metrics:                 true,
		Astria_sequencer_metrics_http_listener_addr: "127.0.0.1:9000",
		Astria_sequencer_pretty_print:               true,

		Astria_composer_log:                        "astria_composer=info",
		Astria_composer_no_otel:                    true,
		Astria_composer_force_stdout:               true,
		Astria_composer_pretty_print:               true,
		Astria_composer_api_listen_addr:            "0.0.0.0:0",
		Astria_composer_sequencer_url:              "http://127.0.0.1:26657",
		Astria_composer_sequencer_chain_id:         "astria-dusk-5",
		Astria_composer_rollups:                    "astriachain::ws://127.0.0.1:8546",
		Astria_composer_private_key:                "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90",
		Astria_composer_max_submit_interval_ms:     2000,
		Astria_composer_max_bytes_per_bundle:       200000,
		Astria_composer_bundle_queue_capacity:      40000,
		Astria_composer_no_metrics:                 true,
		Astria_composer_metrics_http_listener_addr: "127.0.0.1:9000",
		Astria_composer_grpc_addr:                  "0.0.0.0:0",

		No_color: "",

		Otel_exporter_otlp_endpoint:           "http://localhost:4317",
		Otel_exporter_otlp_traces_endpoint:    "http://localhost:4317/v1/traces",
		Otel_exporter_otlp_traces_timeout:     10,
		Otel_exporter_otlp_traces_compression: "gzip",
		Otel_exporter_otlp_headers:            "key1=value1,key2=value2",
		Otel_exporter_otlp_trace_headers:      "key1=value1,key2=value2",
	}
}

// NetworkConfigs is the struct that holds the configuration for all individual Astria networks.
type NetworkConfigs struct {
	Local   NetworkConfig `mapstructure:"local" toml:"local"`
	Dusk    NetworkConfig `mapstructure:"dusk" toml:"dusk"`
	Dawn    NetworkConfig `mapstructure:"dawn" toml:"dawn"`
	Mainnet NetworkConfig `mapstructure:"mainnet" toml:"mainnet"`
}

// NetworkConfig is the struct that holds the configuration for an individual Astria network.
type NetworkConfig struct {
	SequencerChainId string `mapstructure:"sequencer_chain_id" toml:"sequencer_chain_id"`
	SequencerGRPC    string `mapstructure:"sequencer_grpc" toml:"sequencer_grpc"`
	SequencerRPC     string `mapstructure:"sequencer_rpc" toml:"sequencer_rpc"`
	RollupName       string `mapstructure:"rollup_name" toml:"rollup_name"`
	DefaultDenom     string `mapstructure:"default_denom" toml:"default_denom"`
}

// DefaultNetworksConfig returns a NetworksConfig struct populated with all
// network defaults.
func DefaultNetworksConfigs() NetworkConfigs {
	config := NetworkConfigs{
		Local: NetworkConfig{
			SequencerChainId: "sequencer-test-chain-0",
			SequencerGRPC:    "http://127.0.0.1:8080",
			SequencerRPC:     "http://127.0.0.1:26657",
			RollupName:       "astria-test-chain",
			DefaultDenom:     "nria",
		},
		Dusk: NetworkConfig{
			SequencerChainId: "astria-dusk-5",
			SequencerGRPC:    "https://grpc.sequencer.dusk-5.devnet.astria.org/",
			SequencerRPC:     "https://rpc.sequencer.dusk-5.devnet.astria.org/",
			RollupName:       "",
			DefaultDenom:     "nria",
		},
		Dawn: NetworkConfig{
			SequencerChainId: "astria-dawn-0",
			SequencerGRPC:    "https://grpc.sequencer.dawn-0.devnet.astria.org/",
			SequencerRPC:     "https://rpc.sequencer.dawn-0.devnet.astria.org/",
			RollupName:       "",
			DefaultDenom:     "ibc/channel0/utia",
		},
		Mainnet: NetworkConfig{
			SequencerChainId: "astria",
			SequencerGRPC:    "https://grpc.sequencer.astria.org/",
			SequencerRPC:     "https://rpc.sequencer.astria.org/",
			RollupName:       "",
			DefaultDenom:     "ibc/channel0/utia",
		},
	}
	// }
	return config
}

// LoadNetworksConfigsOrPanic loads the NetworksConfig from the given path. If the file
// cannot be loaded or parsed, the function will panic.
func LoadNetworksConfigsOrPanic(path string) NetworkConfigs {
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
// the default local denom and local sequencer network chain id .
// It will skip initialization if the file already exists. It will panic if the
// file cannot be created or written to.
func CreateNetworksConfig(path, localSequencerChainId, localDefaultDenom string) {

	_, err := os.Stat(path)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", path)
		return
	}
	// Create an instance of the Config struct with some data
	config := DefaultNetworksConfigs()
	config.Local.DefaultDenom = localDefaultDenom
	config.Local.SequencerChainId = localSequencerChainId

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

// CreateBaseConfig creates a networks configuration file at
// the given path, populating the file with the network defaults, and overriding
// the default local denom and local sequencer network chain id .
// It will skip initialization if the file already exists. It will panic if the
// file cannot be created or written to.
func CreateBaseConfig(path, instance string) {

	_, err := os.Stat(path)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", path)
		return
	}
	// Create an instance of the Config struct with some data
	config := NewBaseConfig(instance)

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

// LoadBaseConfigOrPanic loads the BaseConfig from the given path. If the file
// cannot be loaded or parsed, the function will panic.
func LoadBaseConfigOrPanic(path string) BaseConfig {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}

	var config BaseConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}

	return config
}

// GetEnvOverrides returns a slice of environment variables that can be used to
// override the default environment variables for the network configuration. It
// uses the BaseConfig to properly update the ASTRIA_COMPOSER_ROLLUPS env var.
func (n NetworkConfig) GetEnvOverrides(bc BaseConfig) []string {
	rollupEndpoint := bc.Astria_composer_rollups
	// find the ip:port from the rollup endpoint
	pattern := `(\d+\.\d+\.\d+\.\d+:\d+)`
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
		"ASTRIA_COMPOSER_ROLLUPS=" + n.RollupName + "::ws://" + match,
	}
}
