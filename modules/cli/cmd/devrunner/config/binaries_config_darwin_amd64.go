//go:build darwin && amd64

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
		UrlSuffix: "_darwin_amd64.tar.gz",
	},
	AstriaSequencer: UrlParts{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/sequencer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-sequencer-x86_64-apple-darwin.tar.gz",
	},
	AstriaComposer: UrlParts{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/composer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-composer-x86_64-apple-darwin.tar.gz",
	},
	AstriaConductor: UrlParts{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/conductor-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-conductor-x86_64-apple-darwin.tar.gz",
	},
}
