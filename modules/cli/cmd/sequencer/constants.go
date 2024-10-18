package sequencer

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
)

const (
	DefaultConfigDirName                   = ".astria"
	DefaultAddressPrefix                   = "astria"
	DefaultSequencerURL                    = "https://rpc.sequencer.dawn-" + cmd.DefaultDawnNum + ".astria.org"
	DefaultSequencerChainID                = "dawn-" + cmd.DefaultDawnNum
	DefaultAsset                           = "ntia"
	DefaultFeeAsset                        = "ntia"
	DefaultSequencerNetworksConfigFilename = "sequencer-networks-config.toml"
)
