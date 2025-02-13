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
	DefaultRollupName                = "astria-test-chain-0"
	DefaultNetworksConfigName        = "networks-config.toml"
	DefaultServiceLogLevel           = "info"
	DefaultLocalNativeDenom          = "ntia"
	DefaultTUIConfigName             = "tui-config.toml"
	DefaultHighlightColor            = "blue"
	DefaultBorderColor               = "gray"
	DefaultMaxUiLogLines             = 1000
	DefaultRollupPort                = "8546"

	// NOTE - do not include the 'v' at the beginning of the version number
	// Service versions matched to live networks
	MainnetCometbftVersion        = "0.38.11"
	MainnetAstriaSequencerVersion = "1.0.0"
	MainnetAstriaComposerVersion  = "1.0.0"
	MainnetAstriaConductorVersion = "1.0.0"

	// Testnet service versions
	TestnetCometbftVersion  = "0.38.11"
	TestnetSequencerVersion = "1.0.0-rc.2"
	TestnetComposerVersion  = "1.0.0-rc.2"
	TestnetConductorVersion = "1.0.0-rc.2"

	// Devnet service versions
	DevnetCometbftVersion  = "0.38.11"
	DevnetSequencerVersion = "1.0.0-rc.2"
	DevnetComposerVersion  = "1.0.0-rc.2"
	DevnetConductorVersion = "1.0.0-rc.2"
)
