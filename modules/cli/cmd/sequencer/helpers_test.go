package sequencer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPortToURL(t *testing.T) {
	t.Run("URL without port", func(t *testing.T) {
		url := "http://example.com"
		got := addPortToURL(url)
		want := "http://example.com:80"
		assert.Equal(t, want, got)
	})

	t.Run("URL with port", func(t *testing.T) {
		url := "http://example.com:80"
		got := addPortToURL(url)
		want := "http://example.com:80"
		assert.Equal(t, want, got)
	})

	t.Run("https URL without port", func(t *testing.T) {
		url := "https://example.com"
		got := addPortToURL(url)
		want := "https://example.com:443"
		assert.Equal(t, want, got)
	})

	t.Run("https URL with port", func(t *testing.T) {
		url := "https://example.com:443"
		got := addPortToURL(url)
		want := "https://example.com:443"
		assert.Equal(t, want, got)
	})
}
