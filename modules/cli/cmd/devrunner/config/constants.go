package config

const (
	duskNum                          = "9"
	dawnNum                          = "0"
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
	LocalNativeDenom                 = "nria"

	// NOTE - do not include the 'v' at the beginning of the version number
	CometbftVersion        = "0.38.8"
	AstriaSequencerVersion = "0.15.0"
	AstriaComposerVersion  = "0.8.1"
	AstriaConductorVersion = "0.19.0"
)
