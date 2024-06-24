package sequencer

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
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

func TestAddressFromText(t *testing.T) {
	expectedOutput := &primproto.Address{
		Inner: []byte{0x1c, 0xc, 0x49, 0xf, 0x1b, 0x55, 0x28, 0xd8, 0x17, 0x3c, 0x5d, 0xe4, 0x6d, 0x13, 0x11, 0x60, 0xe4, 0xb2, 0xc0, 0xc3},
	}

	// Test table
	tests := []struct {
		name           string
		input          string
		expectedOutput *primproto.Address
		expectError    bool
	}{
		{
			name:           "Valid hex address",
			input:          "0x1c0c490f1b5528d8173c5de46d131160e4b2c0c3",
			expectedOutput: expectedOutput,
			expectError:    false,
		},
		{
			name:           "Valid hex address without '0x'",
			input:          "1c0c490f1b5528d8173c5de46d131160e4b2c0c3",
			expectedOutput: expectedOutput,
			expectError:    false,
		},
		{
			name:           "Invalid hex address",
			input:          "0xXYZ",
			expectedOutput: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := addressFromText(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but did not get one")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got one: %v", err)
				}
				assert.Equal(t, tt.expectedOutput, actual)
			}
		})
	}
}

func TestRollupIdFromText(t *testing.T) {
	rollupID := "steezeburger"
	actual := rollupIdFromText(rollupID)
	expected := &primproto.RollupId{
		Inner: []uint8{0x18, 0x88, 0x7, 0x48, 0xea, 0xe, 0x3c, 0xff, 0xd1, 0xcd, 0x64, 0xc1, 0xc, 0x23, 0x59, 0x31, 0xf4, 0xce, 0x4, 0x0, 0xa5, 0xae, 0xd6, 0x9c, 0x5f, 0x15, 0x57, 0x58, 0x82, 0x29, 0x9a, 0x3d},
	}
	assert.Equal(t, expected, actual)
}

func TestAddressFromPublicKey(t *testing.T) {
	// bech32m address encoded from 1c0c490f1b5528d8173c5de46d131160e4b2c0c3 bytes
	expected := "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm"

	testFromPrivKey := "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90"

	from, _ := privateKeyFromText(testFromPrivKey)
	pub := from.Public().(ed25519.PublicKey)
	actual, err := addressFromPublicKey("astria", pub)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestPrivateKeyFromText(t *testing.T) {
	privkey := "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90"
	bytes, _ := hex.DecodeString(privkey)
	expected := ed25519.NewKeyFromSeed(bytes)
	actual, err := privateKeyFromText(privkey)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
