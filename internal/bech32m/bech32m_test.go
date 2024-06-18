package bech32m

import (
	"testing"

	"encoding/hex"

	"github.com/astriaorg/go-sequencer-client/client"
	"github.com/stretchr/testify/assert"
)

// Bech32Address and Bech32MAddress were encoded from the LegacyAddress using
// this bech32 encoding tool: https://slowli.github.io/bech32-buffer/
const LegacyAddress = "1c0c490f1b5528d8173c5de46d131160e4b2c0c3"
const Bech32Address = "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrweg9de"
const Bech32MAddress = "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm"
const BechAddressPrefix = "astria"

func TestBech32MValidate(t *testing.T) {
	// bech32m address should work
	err := ValidateBech32M(Bech32MAddress, BechAddressPrefix)
	if err != nil {
		t.Fatalf("failed to validate bech32m address: %v", err)
	}

	// checking legitimate bech32m address against a different prefix should fail
	err = ValidateBech32M(Bech32MAddress, "differentprefix")
	if err == nil {
		t.Fatalf("failed to validate bech32m address: %v", err)
	}

	// bech32 address should fail
	err = ValidateBech32M(Bech32Address, BechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32 address as bech32m")
	}

	// non bech32m addresses should fail
	err = ValidateBech32M(LegacyAddress, BechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated non-bech32 address as bech32m")
	}

	// bech32m address with typo in prefix should fail
	err = ValidateBech32M("astri1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm", BechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32m address with typo in prefix")
	}

	// bech32m address with missing characters should fail
	// full address "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9c[fgm]"
	// chars in [] are removed
	err = ValidateBech32M("astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9c", BechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32m address with missing characters")
	}

	// bech32m address with different prefix should fail
	err = ValidateBech32M("otherp1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm", BechAddressPrefix)
	if err == nil {
		t.Fatalf("incorrectly validated bech32m address with different prefix")
	}
}

func TestBech32MEncode(t *testing.T) {
	bytes, err := hex.DecodeString(LegacyAddress)
	if err != nil {
		t.Fatalf("could not decode address as hex: %v", err)
	}
	bechString, err := EncodeBech32M(BechAddressPrefix, bytes)
	if err != nil {
		t.Fatalf("could not encode address as bech32m: %v", err)
	}

	assert.Equal(t, Bech32MAddress, bechString)
}

func TestBech32MDecode(t *testing.T) {
	bytes, err := DecodeBech32M(Bech32MAddress, BechAddressPrefix)
	if err != nil {
		t.Fatalf("could not decode bech32m address: %v", err)
	}

	assert.Equal(t, LegacyAddress, hex.EncodeToString(bytes[:]))
}

func TestBech32MEncodeDecode(t *testing.T) {
	signer, err := client.GenerateSigner()
	if err != nil {
		t.Fatalf("failed to generate signer: %v", err)
	}
	addressBytes := signer.Address()

	encoded, err := EncodeBech32M(BechAddressPrefix, addressBytes[:])
	if err != nil {
		t.Fatalf("failed to encode bech32m address: %v", err)
	}

	address, err := DecodeBech32M(encoded, BechAddressPrefix)
	if err != nil {
		t.Fatalf("failed to decode bech32m address: %v", err)
	}

	assert.Equal(t, addressBytes, address)
}
