package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	util "github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/utilities"

	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// BaseConfig is a map of string key-value pairs that represent the base
// configuration for all services in the Astria stack.
//
// The key-values pairs are also parsed into environment variables for the
// services to consume. For example:
//
//	astria_var = 'value'
//
// in the toml, becomes:
//
//	ASTRIA_VAR='value'
//
// as an env var.
//
// A map was used here to allow for dynamically adding new configuration options
// to the config toml file.
type BaseConfig map[string]string

// NewBaseConfig creates a new BaseConfig.
func NewBaseConfig(instanceName, seqNetworkName, rollupName, feeAsset string) BaseConfig {
	return map[string]string{
		// conductor
		"astria_conductor_celestia_block_time_ms":        "1200",
		"astria_conductor_no_celestia_auth":              "true",
		"astria_conductor_celestia_bearer_token":         "<JWT Bearer token>",
		"astria_conductor_celestia_node_http_url":        "http://127.0.0.1:26658",
		"astria_conductor_execution_rpc_url":             "http://127.0.0.1:50051",
		"astria_conductor_execution_commit_level":        "SoftOnly",
		"astria_conductor_log":                           "astria_conductor=info",
		"astria_conductor_no_otel":                       "true",
		"astria_conductor_force_stdout":                  "true",
		"astria_conductor_pretty_print":                  "true",
		"astria_conductor_sequencer_grpc_url":            "http://127.0.0.1:8080",
		"astria_conductor_sequencer_cometbft_url":        "http://127.0.0.1:26657",
		"astria_conductor_sequencer_block_time_ms":       "2000",
		"astria_conductor_sequencer_requests_per_second": "500",
		"astria_conductor_no_metrics":                    "true",
		"astria_conductor_metrics_http_listener_addr":    "127.0.0.1:9000",
		"astria_conductor_expected_sequencer_chain_id":   seqNetworkName,
		"astria_conductor_expected_celestia_chain_id":    "",

		// sequencer
		"astria_sequencer_listen_addr":                 "127.0.0.1:26658",
		"astria_sequencer_db_filepath":                 filepath.Join("~", ".astria", instanceName, DataDirName, "astria_sequencer_db"),
		"astria_sequencer_mempool_parked_max_tx_count": "200",
		"astria_sequencer_grpc_addr":                   "127.0.0.1:8080",
		"astria_sequencer_log":                         "astria_sequencer=info",
		"astria_sequencer_no_otel":                     "true",
		"astria_sequencer_force_stdout":                "true",
		"astria_sequencer_no_metrics":                  "true",
		"astria_sequencer_metrics_http_listener_addr":  "127.0.0.1:9000",
		"astria_sequencer_pretty_print":                "true",

		// composer
		"astria_composer_log":                        "astria_composer=info",
		"astria_composer_no_otel":                    "true",
		"astria_composer_force_stdout":               "true",
		"astria_composer_pretty_print":               "true",
		"astria_composer_api_listen_addr":            "0.0.0.0:0",
		"astria_composer_sequencer_abci_endpoint":    "http://127.0.0.1:26657",
		"astria_composer_sequencer_grpc_endpoint":    "http://127.0.0.1:8080",
		"astria_composer_sequencer_chain_id":         seqNetworkName,
		"astria_composer_rollups":                    rollupName + "::ws://127.0.0.1:8546",
		"astria_composer_private_key_file":           filepath.Join("~", ".astria", instanceName, DefaultConfigDirName, "composer_dev_priv_key"),
		"astria_composer_sequencer_address_prefix":   "astria",
		"astria_composer_max_submit_interval_ms":     "2000",
		"astria_composer_max_bytes_per_bundle":       "200000",
		"astria_composer_bundle_queue_capacity":      "40000",
		"astria_composer_no_metrics":                 "true",
		"astria_composer_metrics_http_listener_addr": "127.0.0.1:9000",
		"astria_composer_grpc_addr":                  "0.0.0.0:0",
		"astria_composer_fee_asset":                  feeAsset,

		// ANSI
		"no_color": "",

		// otel
		"otel_exporter_otlp_endpoint":           "http://localhost:4317",
		"otel_exporter_otlp_traces_endpoint":    "http://localhost:4317/v1/traces",
		"otel_exporter_otlp_traces_timeout":     "10",
		"otel_exporter_otlp_traces_compression": "gzip",
		"otel_exporter_otlp_headers":            "key1=value1,key2=value2",
		"otel_exporter_otlp_trace_headers":      "key1=value1,key2=value2",
	}
}

// CreateBaseConfig creates a base configuration file at
// the given path, and populates the file.
//
// It will skip initialization if the file already exists.
//
// Panics if the file cannot be created or written to.
func CreateBaseConfig(path, instance, sequencerNetworkName, rollupName, feeAsset string) {
	_, err := os.Stat(path)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", path)
		return
	}
	// create an instance of the Config struct with some data
	config := NewBaseConfig(instance, sequencerNetworkName, rollupName, feeAsset)

	// open a file for writing
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

// LoadBaseConfigOrPanic loads the BaseConfig from the given path.
//
// Panics if the file cannot be loaded or parsed.
func LoadBaseConfigOrPanic(path string) BaseConfig {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}

	// var config BaseConfig
	config := make(map[string]string)
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}

	for key, value := range config {
		config[key] = util.ShellExpand(value)
	}

	return config
}

// ToSlice creates a []string of "key=value" pairs out of a BaseConfig.
//
// The variable name will become the env var key and that variable's value will
// be the value.
func (b BaseConfig) ToSlice() []string {
	var output []string

	for key, value := range b {
		output = append(output, fmt.Sprintf("%s=%s", strings.ToUpper(key), value))
	}

	return output
}
