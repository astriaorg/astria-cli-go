package sequencer

import "github.com/astria/astria-cli-go/cmd"

const (
	DefaultSequencerURL     = "http://127.0.0.1:26657"
	DefaultSequencerChainID = cmd.DefaultLocalSequencerChainID
	DefaultBridgeAssetID    = "transfer/channel-0/utia"
	DefaultBridgeFeeAssetID = "nria"
)
