//go:build darwin && arm64

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
		Url:     "https://github.com/cometbft/cometbft/releases/download/v" + CometbftVersion + "/cometbft_" + CometbftVersion + "_darwin_arm64.tar.gz",
	},
	AstriaSequencer: Binary{
		Name:    "astria-sequencer",
		Version: "v" + AstriaSequencerVersion,
		Url:     "https://github.com/astriaorg/astria/releases/download/sequencer-v" + AstriaSequencerVersion + "/astria-sequencer-aarch64-apple-darwin.tar.gz",
	},
	AstriaComposer: Binary{
		Name:    "astria-composer",
		Version: "v" + AstriaComposerVersion,
		Url:     "https://github.com/astriaorg/astria/releases/download/composer-v" + AstriaComposerVersion + "/astria-composer-aarch64-apple-darwin.tar.gz",
	},
	AstriaConductor: Binary{
		Name:    "astria-conductor",
		Version: "v" + AstriaConductorVersion,
		Url:     "https://github.com/astriaorg/astria/releases/download/conductor-v" + AstriaConductorVersion + "/astria-conductor-aarch64-apple-darwin.tar.gz",
	},
}
