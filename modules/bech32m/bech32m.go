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

// Bech32MAddressFromString converts a bech32m string into a Bech32MAddress struct.
func Bech32MAddressFromString(address string) (*Bech32MAddress, error) {
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

// VerifyBech32MAddress verifies that a bech32m string in a valid address. It
// will return nil if the address is valid, otherwise it will return an error.
func VerifyBech32MAddress(address string) error {
	prefix, byteAddress, err := bech32.Decode(address)
	if err != nil {
		return fmt.Errorf("address must be a bech32 encoded string")
	}
	byteAddress, err = bech32.ConvertBits(byteAddress, 5, 8, false)
	if err != nil {
		return fmt.Errorf("failed to convert address to 8 bit")
	}
	if prefix == "" {
		return fmt.Errorf("address must have prefix")
	}
	if len(byteAddress) != 20 {
		return fmt.Errorf("address must have resolve to 20 byte address, got %d", len(byteAddress))
	}

	return nil
}

// Bech32MAddressFromBytes creates a bech32m address from a [20]byte array and string
// prefix.
func Bech32MAddressFromBytes(prefix string, data [20]byte) (*Bech32MAddress, error) {
	// Convert the data from 8-bit groups to 5-bit
	convertedBytes, err := bech32.ConvertBits(data[:], 8, 5, true)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bits from 8-bit groups to 5-bit groups: %v", err)
	}

	// Encode the data as Bech32m
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

// Bech32MAddressFromPublicKey converts an ed25519 public key to a hexadecimal string representation of its address.
func Bech32MAddressFromPublicKey(prefix string, pubkey ed25519.PublicKey) (*Bech32MAddress, error) {
	hash := sha256.Sum256(pubkey)
	var addr [20]byte
	copy(addr[:], hash[:20])
	address, err := Bech32MAddressFromBytes(prefix, addr)
	if err != nil {
		log.WithError(err).Error("Error encoding address")
		return nil, err
	}
	return address, nil
}

// Bech32MAddressFromPrivateKey converts an ed25519 private key into a Bech32MAddress.
func Bech32MAddressFromPrivateKey(prefix string, privkey ed25519.PrivateKey) (*Bech32MAddress, error) {
	from := ed25519.NewKeyFromSeed(privkey)
	pubkey := from.Public().(ed25519.PublicKey)
	hash := sha256.Sum256(pubkey)
	var addr [20]byte
	copy(addr[:], hash[:20])
	address, err := Bech32MAddressFromBytes(prefix, addr)
	if err != nil {
		log.WithError(err).Error("Error encoding address")
		return nil, err
	}
	return address, nil
}
