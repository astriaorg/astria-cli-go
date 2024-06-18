package bech32m

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/bech32"
)

// ValidateBech32M checks if the address is a valid Bech32m address and has the
// correct prefix. It returns nil if the address is valid, otherwise it returns
// an error.
func ValidateBech32M(address, expectedPrefix string) error {
	// Decode the Bech32m address
	hrp, _, version, err := bech32.DecodeGeneric(address)
	if err != nil {
		return fmt.Errorf("failed to decode address as bech32m: %v", err)
	}

	// Check if the version is Bech32m and not a different bech32 version
	if version != bech32.VersionM {
		return fmt.Errorf("bech32 address is not a bech32m address")
	}

	// expected prefix should be all lowercase
	if expectedPrefix != strings.ToLower(expectedPrefix) {
		return fmt.Errorf("expected prefix should be all lowercase: got %s", expectedPrefix)
	}

	// Check if the human-readable prefix matches the expected prefix
	if hrp != expectedPrefix {
		return fmt.Errorf("invalid address prefix: got %s, want %s", hrp, expectedPrefix)
	}

	return nil
}

// DecodeBech32MAsBytes converts a Bech32m address to a byte slice.
func DecodeBech32M(address string, prefix string) ([20]byte, error) {
	err := ValidateBech32M(address, prefix)
	if err != nil {
		return [20]byte{}, fmt.Errorf("failed to validate addres as bech32m: %v", err)
	}

	// can ignore the version here because we already validated the address
	_, data, _, err := bech32.DecodeGeneric(address)
	if err != nil {
		return [20]byte{}, fmt.Errorf("failed to decode address: %v", err)
	}

	// Convert the data from 5-bit groups back to 8-bit
	decoded, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return [20]byte{}, fmt.Errorf("failed to convert bits from 5-bit groups to 8-bit groups: %v", err)
	}

	var address20 [20]byte
	copy(address20[:], decoded)

	return address20, nil
}

// EncodeBech32M creates a bech32m address from a byte slice and string
// prefix.
func EncodeBech32M(prefix string, data []byte) (string, error) {
	// Convert the data from 8-bit groups to 5-bit
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("failed to convert bits from 8-bit groups to 5-bit groups: %v", err)
	}

	// Encode the data as Bech32m
	address, err := bech32.EncodeM(prefix, converted)
	if err != nil {
		return "", fmt.Errorf("failed to encode address as bech32m: %v", err)
	}

	return address, nil
}
