//go:build darwin && amd64

package devtools

// Your code here
type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v0.37.4/cometbft_0.37.4_darwin_amd64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v0.9.0/astria-sequencer-x86_64-apple-darwin.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.4.0/astria-composer-x86_64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.12.0/astria-conductor-x86_64-apple-darwin.tar.gz"},
}
