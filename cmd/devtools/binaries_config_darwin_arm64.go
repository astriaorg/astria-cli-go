//go:build darwin && arm64

package cmd

// Your code here
type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.3.1/astria-composer-aarch64-apple-darwin.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.11.1/astria-conductor-aarch64-apple-darwin.tar.gz"},
}
