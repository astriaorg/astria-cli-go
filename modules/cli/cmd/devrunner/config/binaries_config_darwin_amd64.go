//go:build darwin && amd64

package config

var ServiceUrls = AvailableServiceReleaseUrls{
	CometBFT: Url{
		Prefix: "https://github.com/cometbft/cometbft/releases/download/v",
		Middle: "/cometbft_",
		Suffix: "_darwin_amd64.tar.gz",
	},
	AstriaSequencer: Url{
		Prefix: "https://github.com/astriaorg/astria/releases/download/sequencer-v",
		Middle: "",
		Suffix: "/astria-sequencer-x86_64-apple-darwin.tar.gz",
	},
	AstriaComposer: Url{
		Prefix: "https://github.com/astriaorg/astria/releases/download/composer-v",
		Middle: "",
		Suffix: "/astria-composer-x86_64-apple-darwin.tar.gz",
	},
	AstriaConductor: Url{
		Prefix: "https://github.com/astriaorg/astria/releases/download/conductor-v",
		Middle: "",
		Suffix: "/astria-conductor-x86_64-apple-darwin.tar.gz",
	},
}
