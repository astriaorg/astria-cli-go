//go:build darwin && arm64

package devtools

type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v0.38.6/cometbft_0.38.6_darwin_arm64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v0.11.0/astria-sequencer-aarch64-apple-darwin.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.6.0/astria-composer-aarch64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.14.0/astria-conductor-aarch64-apple-darwin.tar.gz"},
}
