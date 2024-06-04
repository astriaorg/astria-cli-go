package config

import "github.com/astria/astria-cli-go/cmd"

const (
	BinariesDirName                  = "bin"
	DataDirName                      = "data"
	DefaultBaseConfigName            = "base-config.toml"
	DefaultCometbftGenesisFilename   = "genesis.json"
	DefaultCometbftValidatorFilename = "priv_validator_key.json"
	DefaultConfigDirName             = "config"
	DefaultInstanceName              = "default"
	DefaultLocalNetworkName          = cmd.DefaultLocalSequencerChainID
	DefaultNetworksConfigName        = "networks-config.toml"
	DefaultServiceLogLevel           = "info"
	DefaultTargetNetwork             = "local"
	LocalNativeDenom                 = "nria"
)
