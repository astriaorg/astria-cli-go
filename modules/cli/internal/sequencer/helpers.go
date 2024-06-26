package sequencer

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"

	"github.com/btcsuite/btcd/btcutil/bech32"
	log "github.com/sirupsen/logrus"
)

// rollupIdFromText converts a string to a RollupId protobuf.
func rollupIdFromText(rollup string) *primproto.RollupId {
	hash := sha256.Sum256([]byte(rollup))
	return &primproto.RollupId{
		Inner: hash[:],
	}
}

// addressFromPublicKey converts an ed25519 public key to a hexadecimal string representation of its address.
func addressFromPublicKey(prefix string, pubkey ed25519.PublicKey) (*Bech32MAddress, error) {
	hash := sha256.Sum256(pubkey)
	var addr [20]byte
	copy(addr[:], hash[:20])
	address, err := Bech32MFromBytes(prefix, addr)
	if err != nil {
		log.WithError(err).Error("Error encoding address")
		return nil, err
	}
	return address, nil
}

// privateKeyFromText converts a string representation of a private key to an ed25519.PrivateKey.
// It decodes the private key from hex string format and creates a new ed25519.PrivateKey.
func privateKeyFromText(privkey string) (ed25519.PrivateKey, error) {
	privKeyBytes, err := hex.DecodeString(privkey)
	if err != nil {
		return nil, err
	}
	from := ed25519.NewKeyFromSeed(privKeyBytes)
	return from, nil
}

// Bech32MFromBytes creates a bech32m address from a [20]byte array and string
// prefix.
func Bech32MFromBytes(prefix string, data [20]byte) (*Bech32MAddress, error) {
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
		Address: address,
		Prefix:  prefix,
		Bytes:   data,
	}, nil
}

// assetIdFromDenom returns a hash of a denom string
func assetIdFromDenom(denom string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(denom))
	hash := hasher.Sum(nil)
	return hash
}
