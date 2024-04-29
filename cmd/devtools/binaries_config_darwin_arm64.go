//go:build darwin && arm64

package devtools

type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v0.37.4/cometbft_0.37.4_darwin_arm64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v0.10.1/astria-sequencer-aarch64-apple-darwin.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.5.0/astria-composer-aarch64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.13.1/astria-conductor-aarch64-apple-darwin.tar.gz"},
}
