//go:build darwin && amd64

package config

type Binary struct {
	Name    string
	Version string
	Url     string
}

var Binaries = []Binary{
	{"cometbft", "v" + CometbftVersion, "https://github.com/cometbft/cometbft/releases/download/v" + CometbftVersion + "/cometbft_" + CometbftVersion + "_darwin_amd64.tar.gz"},
	{"astria-sequencer", "v" + AstriaSequencerVersion, "https://github.com/astriaorg/astria/releases/download/sequencer-v" + AstriaSequencerVersion + "/astria-sequencer-x86_64-apple-darwin.tar.gz"},
	{"astria-composer", "v" + AstriaComposerVersion, "https://github.com/astriaorg/astria/releases/download/composer-v" + AstriaComposerVersion + "/astria-composer-x86_64-apple-darwin.tar.gz"},
	{"astria-conductor", "v" + AstriaConductorVersion, "https://github.com/astriaorg/astria/releases/download/conductor-v" + AstriaConductorVersion + "/astria-conductor-x86_64-apple-darwin.tar.gz"},
}
