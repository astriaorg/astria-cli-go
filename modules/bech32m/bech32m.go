package bech32m

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/bech32"
	log "github.com/sirupsen/logrus"
)

type Bech32MAddress struct {
	address string
	prefix  string
	bytes   [20]byte
}

// String returns the bech32m address as a string
func (a *Bech32MAddress) String() string {
	return a.address
}

// Prefix returns the prefix of the bech32m address
func (a *Bech32MAddress) Prefix() string {
	return a.prefix
}

// Bytes returns the underlying bytes for the bech32m address as a [20]byte array
func (a *Bech32MAddress) Bytes() [20]byte {
	return a.bytes
}

// Verify verifies that a string in a valid bech32m address. It
// will return nil if the address is valid, otherwise it will return an error.
func Verify(address string) error {
	prefix, byteAddress, version, err := bech32.DecodeGeneric(address)
	if err != nil {
		return fmt.Errorf("address must be a bech32 encoded string")
	}
	if version != bech32.VersionM {
		return fmt.Errorf("address must be a bech32m address")
	}
	byteAddress, err = bech32.ConvertBits(byteAddress, 5, 8, false)
	if err != nil {
		return fmt.Errorf("failed to convert address to 8 bit")
	}
	if prefix == "" {
		return fmt.Errorf("address must have prefix")
	}
	if len(byteAddress) != 20 {
		return fmt.Errorf("address must decode to a 20 length byte array: got len %d", len(byteAddress))
	}

	return nil
}

// Encode creates a *Bech32MAddress from a [20]byte array and string
// prefix.
func Encode(prefix string, data [20]byte) (*Bech32MAddress, error) {
	// Convert the data from 8-bit groups to 5-bit
	convertedBytes, err := bech32.ConvertBits(data[:], 8, 5, true)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bits from 8-bit groups to 5-bit groups: %v", err)
	}

	// Encode the data as bech32m
	address, err := bech32.EncodeM(prefix, convertedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode address as bech32m: %v", err)
	}

	return &Bech32MAddress{
		address: address,
		prefix:  prefix,
		bytes:   data,
	}, nil
}

// EncodeFromString converts a bech32m string address into a *Bech32MAddress.
func EncodeFromString(address string) (*Bech32MAddress, error) {
	prefix, bytes, err := bech32.Decode(address)
	if err != nil {
		return nil, fmt.Errorf("input address must be a bech32 encoded string")
	}
	convertedBytes, err := bech32.ConvertBits(bytes, 5, 8, false)
	if err != nil {
		return nil, fmt.Errorf("failed to convert address bytes to 8 bit")
	}

	var addrBytes [20]byte
	copy(addrBytes[:], convertedBytes)

	return &Bech32MAddress{
		address: address,
		prefix:  prefix,
		bytes:   addrBytes,
	}, nil
}

// EncodeFromPublicKey takes an ed25519 public key and string prefix and encodes
// them into a *Bech32MAddress.
func EncodeFromPublicKey(prefix string, pubkey ed25519.PublicKey) (*Bech32MAddress, error) {
	hash := sha256.Sum256(pubkey)
	var addr [20]byte
	copy(addr[:], hash[:20])
	address, err := Encode(prefix, addr)
	if err != nil {
		log.WithError(err).Error("Error encoding address")
		return nil, err
	}
	return address, nil
}

// EncodeFromPrivateKey takes an ed25519 private key and string prefix and
// encodes them into a *Bech32MAddress.
func EncodeFromPrivateKey(prefix string, privkey ed25519.PrivateKey) (*Bech32MAddress, error) {
	pubkey := privkey.Public().(ed25519.PublicKey)
	hash := sha256.Sum256(pubkey)
	var addr [20]byte
	copy(addr[:], hash[:20])
	address, err := Encode(prefix, addr)
	if err != nil {
		log.WithError(err).Error("Error encoding address")
		return nil, err
	}
	return address, nil
}

// Decode decodes a bech32m string into a string prefix and the underlying
// [20]byte array originally used to encode the address. It also checks if the
// address is a bech32m address and not a different bech32 version.
func Decode(address string) (string, [20]byte, error) {
	prefix, bytes, version, err := bech32.DecodeGeneric(address)
	if err != nil {
		var defaultBytes [20]byte
		copy(defaultBytes[:], bytes)
		return prefix, defaultBytes, fmt.Errorf("failed to decode address")
	}

	if version != bech32.VersionM {
		var defaultBytes [20]byte
		copy(defaultBytes[:], bytes)
		return prefix, defaultBytes, fmt.Errorf("address must be a bech32m address")
	}

	convertedBytes, err := bech32.ConvertBits(bytes, 5, 8, false)
	if err != nil {
		var defaultBytes [20]byte
		copy(defaultBytes[:], convertedBytes)
		return prefix, defaultBytes, fmt.Errorf("failed to convert address bytes to 8 bit")
	}

	var addrBytes [20]byte
	copy(addrBytes[:], convertedBytes)

	return prefix, addrBytes, nil
}
