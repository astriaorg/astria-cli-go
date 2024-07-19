//go:build linux && amd64

package config

type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v" + CometbftVersion + "/cometbft_" + CometbftVersion + "_linux_amd64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v" + AstriaSequencerVersion + "/astria-sequencer-x86_64-unknown-linux-gnu.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v" + AstriaComposerVersion + "/astria-composer-x86_64-unknown-linux-gnu.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v" + AstriaConductorVersion + "/astria-conductor-x86_64-unknown-linux-gnu.tar.gz"},
}
