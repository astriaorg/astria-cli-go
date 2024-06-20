package bech32m

import (
	"fmt"
	"strings"

	primitives "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	"github.com/btcsuite/btcd/btcutil/bech32"
)

type Bech32MAddress struct {
	// Address is the original address used to build the type.
	address string
	// Prefix is the extracted prefix from the address.
	prefix string
	// Bytes is the decoded bytes from the address.
	bytes [20]byte
}

// ToString returns the address of the Bech32MAddress type as a string.
func (b *Bech32MAddress) ToString() string {
	return b.address
}

// Prefix returns the prefix of the Bech32MAddress type as a string.
func (b *Bech32MAddress) Prefix() string {
	return b.prefix
}

// AsBytes returns the bytes of the Bech32MAddress as a [20]byte array.
func (b *Bech32MAddress) AsBytes() [20]byte {
	return b.bytes
}

// AsProtoAddress returns the bech32m address as an Astria protobuf address.
func (b *Bech32MAddress) AsProtoAddress() *primitives.Address {
	return &primitives.Address{
		Bech32M: b.address,
	}
}

// DecodeBech32MAsBytes decodes and validates a Bech32m address into a
// *Bech32MAddress type for handling of address data.
func DecodeAndValidateBech32M(address string, expectedPrefix string) (*Bech32MAddress, error) {
	// Decode the Bech32m address
	hrp, data, version, err := bech32.DecodeGeneric(address)
	if err != nil {
		return nil, fmt.Errorf("failed to decode address as bech32m: %v", err)
	}

	// Check if the version is Bech32m and not a different bech32 version
	if version != bech32.VersionM {
		return nil, fmt.Errorf("bech32 address is not a bech32m address")
	}

	// expected prefix should be all lowercase
	if expectedPrefix != strings.ToLower(expectedPrefix) {
		return nil, fmt.Errorf("expected prefix should be all lowercase: got %s", expectedPrefix)
	}

	// Check if the human-readable prefix matches the expected prefix
	if hrp != expectedPrefix {
		return nil, fmt.Errorf("invalid address prefix: got %s, want %s", hrp, expectedPrefix)
	}

	// Convert the data from 5-bit groups back to 8-bit
	decoded, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bits from 5-bit groups to 8-bit groups: %v", err)
	}

	var bytes [20]byte
	copy(bytes[:], decoded)

	return &Bech32MAddress{
		address: address,
		prefix:  hrp,
		bytes:   bytes,
	}, nil

}

// EncodeBech32M creates a bech32m address from a [20]byte slice and string
// prefix.
func EncodeBech32M(prefix string, data [20]byte) (*Bech32MAddress, error) {
	// Convert the data from 8-bit groups to 5-bit
	converted, err := bech32.ConvertBits(data[:], 8, 5, true)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bits from 8-bit groups to 5-bit groups: %v", err)
	}

	// Encode the data as Bech32m
	address, err := bech32.EncodeM(prefix, converted)
	if err != nil {
		return nil, fmt.Errorf("failed to encode address as bech32m: %v", err)
	}

	return &Bech32MAddress{
		address: address,
		prefix:  prefix,
		bytes:   data,
	}, nil
}
