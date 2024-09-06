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
	DefaultTargetNetwork             = "local"
	DefaultLocalNativeDenom          = "ntia"
	DefaultTUIConfigName             = "tui-config.toml"
	DefaultHighlightColor            = "blue"
	DefaultBorderColor               = "gray"

	// NOTE - do not include the 'v' at the beginning of the version number
	CometbftVersion        = "0.38.8"
	AstriaSequencerVersion = "0.17.0"
	AstriaComposerVersion  = "0.8.3"
	AstriaConductorVersion = "0.20.1"
)
