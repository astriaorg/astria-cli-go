//go:build darwin && arm64

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
		UrlSuffix: "_darwin_arm64.tar.gz",
	},
	AstriaSequencer: Binary{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/sequencer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-sequencer-aarch64-apple-darwin.tar.gz",
	},
	AstriaComposer: Binary{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/composer-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-composer-aarch64-apple-darwin.tar.gz",
	},
	AstriaConductor: Binary{
		UrlPrefix: "https://github.com/astriaorg/astria/releases/download/conductor-v",
		UrlMiddle: "",
		UrlSuffix: "/astria-conductor-aarch64-apple-darwin.tar.gz",
	},
}
