//go:build darwin && arm64

package config

var ServiceUrls = AvailableServiceReleaseUrls{
	CometBFT: Url{
		Prefix: "https://github.com/cometbft/cometbft/releases/download/v",
		Middle: "/cometbft_",
		Suffix: "_darwin_arm64.tar.gz",
	},
	AstriaSequencer: Url{
		Prefix: "https://github.com/astriaorg/astria/releases/download/sequencer-v",
		Middle: "",
		Suffix: "/astria-sequencer-aarch64-apple-darwin.tar.gz",
	},
	AstriaComposer: Url{
		Prefix: "https://github.com/astriaorg/astria/releases/download/composer-v",
		Middle: "",
		Suffix: "/astria-composer-aarch64-apple-darwin.tar.gz",
	},
	AstriaConductor: Url{
		Prefix: "https://github.com/astriaorg/astria/releases/download/conductor-v",
		Middle: "",
		Suffix: "/astria-conductor-aarch64-apple-darwin.tar.gz",
	},
}
