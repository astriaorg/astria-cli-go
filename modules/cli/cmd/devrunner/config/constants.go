package config

const (
	BinariesDirName                  = "bin"
	LogsDirName                      = "logs"
	DataDirName                      = "data"
	DefaultBaseConfigName            = "base-config.toml"
	DefaultCometbftGenesisFilename   = "genesis.json"
	DefaultCometbftValidatorFilename = "priv_validator_key.json"
	DefaultConfigDirName             = "config"
	DefaultInstanceName              = "default"
	DefaultLocalNetworkName          = "sequencer-test-chain-0"
	DefaultNetworksConfigName        = "networks-config.toml"
	DefaultServiceLogLevel           = "info"
	DefaultLocalNativeDenom          = "ntia"
	DefaultTUIConfigName             = "tui-config.toml"
	DefaultHighlightColor            = "blue"
	DefaultBorderColor               = "gray"

	// NOTE - do not include the 'v' at the beginning of the version number
	// Service versions matched to live networks
	CometbftVersion        = "0.38.8"
	AstriaSequencerVersion = "1.0.0-rc.1"
	AstriaComposerVersion  = "1.0.0-rc.1"
	AstriaConductorVersion = "1.0.0-rc.1"

	// Local service versions in case of differences between networks
	LocalCometbftVersion  = "0.38.8"
	LocalSequencerVersion = "1.0.0-rc.1"
	LocalComposerVersion  = "1.0.0-rc.1"
	LocalConductorVersion = "1.0.0-rc.1"
)
