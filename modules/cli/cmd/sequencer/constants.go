package sequencer

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
)

const (
	DefaultConfigDirName                   = ".astria"
	DefaultAddressPrefix                   = "astria"
	DefaultSequencerURL                    = "https://rpc.sequencer.dusk-" + cmd.DefaultDuskNum + ".devnet.astria.org"
	DefaultSequencerChainID                = "astria-dusk-" + cmd.DefaultDuskNum
	DefaultAsset                           = "ntia"
	DefaultFeeAsset                        = "ntia"
	DefaultSequencerNetworksConfigFilename = "sequencer-networks-config.toml"
	DefaultTargetNetwork                   = "dusk"
)
