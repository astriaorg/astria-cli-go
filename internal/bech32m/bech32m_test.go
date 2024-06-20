package bech32m

import (
	"testing"

	"encoding/hex"

	"github.com/astriaorg/go-sequencer-client/client"
	"github.com/stretchr/testify/assert"
)

// Bech32Address and Bech32MAddress were encoded from the LegacyAddress using
// this bech32 encoding tool: https://slowli.github.io/bech32-buffer/
const legacyAddress = "1c0c490f1b5528d8173c5de46d131160e4b2c0c3"
const bech32Address = "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrweg9de"
const bech32MAddress = "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm"
const bechAddressPrefix = "astria"

func TestBech32MDecodeAndValidate(t *testing.T) {
	// bech32m address should work
	_, err := DecodeAndValidateBech32M(bech32MAddress, bechAddressPrefix)
	if err != nil {
		t.Fatalf("failed to validate bech32m address: %v", err)
	}

	// checking legitimate bech32m address against a different prefix should fail
	_, err = DecodeAndValidateBech32M(bech32MAddress, "differentprefix")
	if err == nil {
		t.Fatalf("incorrectly validated bech32m address with incorrect prefix: %v", err)
	}

	// bech32 address should fail
	_, err = DecodeAndValidateBech32M(bech32Address, bechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32 address as bech32m")
	}

	// non bech32m addresses should fail
	_, err = DecodeAndValidateBech32M(legacyAddress, bechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated non-bech32 address as bech32m")
	}

	// bech32m address with typo in prefix should fail
	_, err = DecodeAndValidateBech32M("astri1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm", bechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32m address with typo in prefix")
	}

	// bech32m address with missing characters should fail
	// full address "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9c[fgm]"
	// chars in [] are removed
	_, err = DecodeAndValidateBech32M("astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9c", bechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32m address with missing characters")
	}

	// bech32m address with different prefix should fail
	_, err = DecodeAndValidateBech32M("otherp1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm", bechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32m address with different prefix")
	}
}

func TestBech32MEncode(t *testing.T) {
	b, err := hex.DecodeString(legacyAddress)
	if err != nil {
		t.Fatalf("could not decode address as hex: %v", err)
	}

	var bytes [20]byte
	copy(bytes[:], b)

	bech32m, err := EncodeBech32M(bechAddressPrefix, bytes)
	if err != nil {
		t.Fatalf("could not encode address as bech32m: %v", err)
	}

	assert.Equal(t, bech32MAddress, bech32m.ToString())
}

func TestBech32MDecode(t *testing.T) {
	bech32m, err := DecodeAndValidateBech32M(bech32MAddress, bechAddressPrefix)
	if err != nil {
		t.Fatalf("could not decode bech32m address: %v", err)
	}

	bytes := bech32m.AsBytes()
	assert.Equal(t, legacyAddress, hex.EncodeToString(bytes[:]))
}

func TestBech32MEncodeDecode(t *testing.T) {
	signer, err := client.GenerateSigner()
	if err != nil {
		t.Fatalf("failed to generate signer: %v", err)
	}
	addressBytes := signer.Address()

	bech32m, err := EncodeBech32M(bechAddressPrefix, addressBytes)
	if err != nil {
		t.Fatalf("failed to encode bech32m address: %v", err)
	}

	address, err := DecodeAndValidateBech32M(bech32m.ToString(), bechAddressPrefix)
	if err != nil {
		t.Fatalf("failed to decode bech32m address: %v", err)
	}

	assert.Equal(t, addressBytes, address.AsBytes())
}
