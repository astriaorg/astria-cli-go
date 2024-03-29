package sequencer

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// strip0xPrefix removes the 0x prefix from a string if present.
func strip0xPrefix(s string) string {
	return strings.TrimPrefix(s, "0x")
}

// addPortToURL adds a port to a URL if it doesn't already have one.
// The port is needed for the
func addPortToURL(url string) string {
	if strings.Contains(url, "http:") {
		return url + ":80"
	}
	if strings.Contains(url, "https:") {
		return url + ":443"
	}
	log.Debug("No port added to URL", url)
	return url
}
