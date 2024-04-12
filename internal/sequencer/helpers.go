package sequencer

import (
	"crypto/sha256"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// strip0xPrefix removes the 0x prefix from a string if present.
func strip0xPrefix(s string) string {
	return strings.TrimPrefix(s, "0x")
}

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

// AssetIdFromDenom returns a hash of a denom string
func AssetIdFromDenom(denom string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(denom))
	hash := hasher.Sum(nil)
	return hash
}
