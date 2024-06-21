//go:build darwin && amd64

package config

type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v" + cometbftVersion + "/cometbft_" + cometbftVersion + "_darwin_amd64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v" + sequencerVersion + "/astria-sequencer-x86_64-apple-darwin.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v" + composerVersion + "/astria-composer-x86_64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v" + conductorVersion + "/astria-conductor-x86_64-apple-darwin.tar.gz"},
}
