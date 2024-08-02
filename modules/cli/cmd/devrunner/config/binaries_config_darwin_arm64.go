//go:build darwin && arm64

package config

type Binary struct {
	Name    string
	Version string
	Url     string
}

var Binaries = []Binary{
	{"cometbft", "v" + CometbftVersion, "https://github.com/cometbft/cometbft/releases/download/v" + CometbftVersion + "/cometbft_" + CometbftVersion + "_darwin_arm64.tar.gz"},
	{"astria-sequencer", "v" + AstriaSequencerVersion, "https://github.com/astriaorg/astria/releases/download/sequencer-v" + AstriaSequencerVersion + "/astria-sequencer-aarch64-apple-darwin.tar.gz"},
	{"astria-composer", "v" + AstriaComposerVersion, "https://github.com/astriaorg/astria/releases/download/composer-v" + AstriaComposerVersion + "/astria-composer-aarch64-apple-darwin.tar.gz"},
	{"astria-conductor", "v" + AstriaConductorVersion, "https://github.com/astriaorg/astria/releases/download/conductor-v" + AstriaConductorVersion + "/astria-conductor-aarch64-apple-darwin.tar.gz"},
}
