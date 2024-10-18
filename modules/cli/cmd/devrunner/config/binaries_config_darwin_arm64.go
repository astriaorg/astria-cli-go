//go:build darwin && arm64

package config

type UrlParts struct {
	UrlPrefix string
	UrlMiddle string
	UrlSuffix string
}

type KnownServiceReleaseUrlsParts struct {
	CometBFT        UrlParts
	AstriaSequencer UrlParts
	AstriaComposer  UrlParts
	AstriaConductor UrlParts
}

var ServiceUrls = KnownServiceReleaseUrlsParts{
	CometBFT: UrlParts{
		UrlPrefix: "https://github.com/cometbft/cometbft/releases/download/v",
		UrlMiddle: "/cometbft_",
		UrlSuffix: "_darwin_arm64.tar.gz",
	},
	AstriaSequencer: UrlParts{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/sequencer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-sequencer-aarch64-apple-darwin.tar.gz",
	},
	AstriaComposer: UrlParts{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/composer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-composer-aarch64-apple-darwin.tar.gz",
	},
	AstriaConductor: UrlParts{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/conductor-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-conductor-aarch64-apple-darwin.tar.gz",
	},
}
