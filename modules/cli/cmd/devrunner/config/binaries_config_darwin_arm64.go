//go:build darwin && arm64

package config

type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v" + cometbft_version + "/cometbft_" + cometbft_version + "_darwin_arm64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v" + astria_sequencer_version + "/astria-sequencer-aarch64-apple-darwin.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v" + astria_composer_version + "/astria-composer-aarch64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v" + astria_conductor_version + "/astria-conductor-aarch64-apple-darwin.tar.gz"},
}
