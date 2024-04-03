//go:build linux && amd64

package devtools

// Your code here
type Binary struct {
	Name string
	Url  string
}

var LocalBinaries = []Binary{
	{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v0.37.4/cometbft_0.37.4_linux_amd64.tar.gz"},
	{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v0.9.0/astria-sequencer-x86_64-unknown-linux-gnu.tar.gz"},
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.4.0/astria-composer-x86_64-unknown-linux-gnu.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.12.0/astria-conductor-x86_64-unknown-linux-gnu.tar.gz"},
}

var RemoteBinaries = []Binary{
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.3.1/astria-composer-x86_64-unknown-linux-gnu.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.11.1/astria-conductor-x86_64-unknown-linux-gnu.tar.gz"},
}
