//go:build darwin && arm64

package config

type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v" + CometbftVersion + "/cometbft_" + CometbftVersion + "_darwin_arm64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v" + AstriaSequencerVersion + "/astria-sequencer-aarch64-apple-darwin.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v" + AstriaComposerVersion + "/astria-composer-aarch64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v" + AstriaConductorVersion + "/astria-conductor-aarch64-apple-darwin.tar.gz"},
}
