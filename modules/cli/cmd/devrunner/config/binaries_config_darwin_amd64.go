//go:build darwin && amd64

package config

type Binary struct {
	UrlPrefix string
	UrlMiddle string
	UrlSuffix string
}

type BinariesInfo struct {
	CometBFT        Binary
	AstriaSequencer Binary
	AstriaComposer  Binary
	AstriaConductor Binary
}

var DownloadUrlParts = BinariesInfo{
	CometBFT: Binary{
		UrlPrefix: "https://github.com/cometbft/cometbft/releases/download/v",
		UrlMiddle: "/cometbft_",
		UrlSuffix: "_darwin_amd64.tar.gz",
	},
	AstriaSequencer: Binary{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/sequencer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-sequencer-x86_64-apple-darwin.tar.gz",
	},
	AstriaComposer: Binary{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/composer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-composer-x86_64-apple-darwin.tar.gz",
	},
	AstriaConductor: Binary{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/conductor-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-conductor-x86_64-apple-darwin.tar.gz",
	},
}
