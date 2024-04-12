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
		assert.Equal(t, want, got)
	})

	t.Run("without '0x' prefix", func(t *testing.T) {
		str := "abcd1234"
		got := strip0xPrefix(str)
		want := "abcd1234"
		assert.Equal(t, want, got)
	})
}

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

func TestConvertToUint128(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedLo  uint64
		expectedHi  uint64
		expectError bool
	}{
		{"Zero", "0", 0, 0, false},
		{"Boundary64", "18446744073709551615", 18446744073709551615, 0, false}, // 2^64 - 1
		// {lo: 0, hi: 1} = 0b00000000000000000000000000000000000000000000000000000000000000001 (1 in 65th) = 0d18446744073709551616
		{"Boundary64PlusOne", "18446744073709551616", 0, 1, false}, // 2^64
		// 0d18446744073709551615 = 0b1111111111111111111111111111111111111111111111111111111111111111 (64 bits),
		{"Boundary128", "340282366920938463463374607431768211455", 18446744073709551615, 18446744073709551615, false}, // 2^128 - 1
		// this number is too big to handle!
		{"Boundary128PlusOne", "340282366920938463463374607431768211456", 0, 0, true}, // 2^128
		{"Small", "1234", 1234, 0, false},
		{"NonNumeric", "nonnumeric", 0, 0, true},
		{"Negative", "-1", 0, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uint128, err := convertToUint128(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error but did not get one")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got one: %v", err)
				}
				if uint128.Lo != tc.expectedLo || uint128.Hi != tc.expectedHi {
					t.Errorf("Expected (%d, %d), got (%d, %d)", tc.expectedLo, tc.expectedHi, uint128.Lo, uint128.Hi)
				}
			}
		})
	}
}
