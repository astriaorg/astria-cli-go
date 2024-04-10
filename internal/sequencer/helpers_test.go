package sequencer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrip0xPrefix(t *testing.T) {
	t.Run("with '0x' prefix", func(t *testing.T) {
		str := "0x1234abcd"
		got := strip0xPrefix(str)
		want := "1234abcd"
		assert.Equal(t, got, want)
	})

	t.Run("without '0x' prefix", func(t *testing.T) {
		str := "abcd1234"
		got := strip0xPrefix(str)
		want := "abcd1234"
		assert.Equal(t, got, want)
	})
}

func TestAddPortToURL(t *testing.T) {
	t.Run("URL without port", func(t *testing.T) {
		url := "http://example.com"
		got := addPortToURL(url)
		want := "http://example.com:80"
		assert.Equal(t, got, want)
	})

	t.Run("URL with port", func(t *testing.T) {
		url := "http://example.com:80"
		got := addPortToURL(url)
		want := "http://example.com:80"
		assert.Equal(t, got, want)
	})

	t.Run("https URL without port", func(t *testing.T) {
		url := "https://example.com"
		got := addPortToURL(url)
		want := "https://example.com:443"
		assert.Equal(t, got, want)
	})

	t.Run("https URL with port", func(t *testing.T) {
		url := "https://example.com:443"
		got := addPortToURL(url)
		want := "https://example.com:443"
		assert.Equal(t, got, want)
	})
}
