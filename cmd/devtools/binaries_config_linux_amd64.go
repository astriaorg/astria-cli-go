//go:build linux && amd64

package cmd

// Your code here
type Binary struct {
	Name string
	Url  string
}

var Binaries = []Binary{
	{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.3.1/astria-composer-x86_64-unknown-linux-gnu.tar.gz"},
	{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.11.1/astria-conductor-x86_64-unknown-linux-gnu.tar.gz"},
}
