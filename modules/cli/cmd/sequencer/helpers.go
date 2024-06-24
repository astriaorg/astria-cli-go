package sequencer

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"

	log "github.com/sirupsen/logrus"
)

// addPortToURL adds a port to a URL if it doesn't already have one.
// The port is needed for use with the Sequencer Client.
func addPortToURL(url string) string {
	// Check if the URL already has a port
	matched, err := regexp.MatchString(`:\d+$`, url)
	if err != nil {
		log.WithError(err).Error("Error matching string")
		return url
	}
	if matched {
		log.Debug("Port already present in URL: ", url)
		return url
	}
	if strings.Contains(url, "http:") {
		log.Debug("http url detected without a port. Adding port :80 to url: ", url)
		return url + ":80"
	}
	if strings.Contains(url, "https:") {
		log.Debug("https url detected without a port. Adding port :443 to url: ", url)
		return url + ":443"
	}
	return url
}

// privateKeyFromText converts a string representation of a private key to an ed25519.PrivateKey.
// It decodes the private key from hex string format and creates a new ed25519.PrivateKey.
func privateKeyFromText(privkey string) (ed25519.PrivateKey, error) {
	privKeyBytes, err := hex.DecodeString(privkey)
	if err != nil {
		return nil, err
	}
	from := ed25519.NewKeyFromSeed(privKeyBytes)
	return from, nil
}

// addressFromText converts a bech32m string representation of an address to an
// Address protobuf. No validation is done on the input string.
func addressFromText(addr string) *primproto.Address {
	return &primproto.Address{
		Bech32M: addr,
	}
}

// assetIdFromDenom returns a hash of a denom string
func assetIdFromDenom(denom string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(denom))
	hash := hasher.Sum(nil)
	return hash
}

// rollupIdFromText converts a string to a RollupId protobuf.
func rollupIdFromText(rollup string) *primproto.RollupId {
	hash := sha256.Sum256([]byte(rollup))
	return &primproto.RollupId{
		Inner: hash[:],
	}
}
