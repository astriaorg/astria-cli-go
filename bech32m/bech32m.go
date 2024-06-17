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
	hrp, _, err := bech32.Decode(address)
	if err != nil {
		return fmt.Errorf("failed to decode address: %v", err)
	}

	// Check if the human-readable prefix matches the expected prefix
	if !strings.EqualFold(hrp, expectedPrefix) {
		return fmt.Errorf("invalid prefix: got %s, want %s", hrp, expectedPrefix)
	}

	return nil
}

// DecodeBech32MAsBytes converts a Bech32m address to a byte slice.
func DecodeBech32M(address string, prefix string) ([20]byte, error) {
	err := ValidateBech32M(address, prefix)
	if err != nil {
		return [20]byte{}, fmt.Errorf("failed to validate addres as bech32m: %v", err)
	}

	_, data, err := bech32.Decode(address)
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
func EncodeBech32M(prefix string, data [20]byte) (string, error) {
	// Convert the data from 8-bit groups to 5-bit
	converted, err := bech32.ConvertBits(data[:], 8, 5, true)
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
