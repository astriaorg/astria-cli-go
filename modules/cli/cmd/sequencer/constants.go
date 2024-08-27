package sequencer

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
)

const (
	DefaultConfigDirName                   = ".astria"
	DefaultAddressPrefix                   = "astria"
	DefaultSequencerURL                    = "https://rpc.sequencer.dusk-" + config.DuskNum + ".devnet.astria.org"
	DefaultSequencerChainID                = "astria-dusk-" + config.DuskNum
	DefaultAsset                           = "nria"
	DefaultFeeAsset                        = "nria"
	DefaultSequencerNetworksConfigFilename = "sequencer-networks-config.toml"
	DefaultTargetNetwork                   = "dusk"
)
