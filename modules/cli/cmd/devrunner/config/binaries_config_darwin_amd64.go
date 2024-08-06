//go:build darwin && amd64

package config

type Binary struct {
	Name    string
	Version string
	Url     string
}

type Binaries struct {
	CometBFT        Binary
	AstriaSequencer Binary
	AstriaComposer  Binary
	AstriaConductor Binary
}

var KnownBinaries = Binaries{
	CometBFT: Binary{
		Name:    "cometbft",
		Version: "v" + CometbftVersion,
		Url:     "https://github.com/cometbft/cometbft/releases/download/v" + CometbftVersion + "/cometbft_" + CometbftVersion + "_darwin_amd64.tar.gz",
	},
	AstriaSequencer: Binary{
		Name:    "astria-sequencer",
		Version: "v" + AstriaSequencerVersion,
		Url:     "https://github.com/astriaorg/astria/releases/download/sequencer-v" + AstriaSequencerVersion + "/astria-sequencer-x86_64-apple-darwin.tar.gz",
	},
	AstriaComposer: Binary{
		Name:    "astria-composer",
		Version: "v" + AstriaComposerVersion,
		Url:     "https://github.com/astriaorg/astria/releases/download/composer-v" + AstriaComposerVersion + "/astria-composer-x86_64-apple-darwin.tar.gz",
	},
	AstriaConductor: Binary{
		Name:    "astria-conductor",
		Version: "v" + AstriaConductorVersion,
		Url:     "https://github.com/astriaorg/astria/releases/download/conductor-v" + AstriaConductorVersion + "/astria-conductor-x86_64-apple-darwin.tar.gz",
	},
}
