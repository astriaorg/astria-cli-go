//go:build darwin && amd64

package config

type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v0.38.6/cometbft_0.38.6_darwin_amd64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v0.13.0/astria-sequencer-x86_64-apple-darwin.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.7.0/astria-composer-x86_64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.16.0/astria-conductor-x86_64-apple-darwin.tar.gz"},
}
