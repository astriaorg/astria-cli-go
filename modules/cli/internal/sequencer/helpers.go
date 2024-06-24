package sequencer

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"

	"github.com/btcsuite/btcd/btcutil/bech32"
	log "github.com/sirupsen/logrus"
)

// convertToUint128 converts a string to an Uint128 protobuf
func convertToUint128(numStr string) (*primproto.Uint128, error) {
	bigInt := new(big.Int)

	// convert the string to a big.Int
	_, ok := bigInt.SetString(numStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert string to big.Int")
	}

	// check if the number is negative or overflows Uint128
	if bigInt.Sign() < 0 {
		return nil, fmt.Errorf("negative number not allowed")
	} else if bigInt.BitLen() > 128 {
		return nil, fmt.Errorf("value overflows Uint128")
	}

	// split the big.Int into two uint64s
	// convert the big.Int to uint64, which will drop the higher 64 bits
	lo := bigInt.Uint64()
	// shift the big.Int to the right by 64 bits and convert to uint64
	hi := bigInt.Rsh(bigInt, 64).Uint64()
	uint128 := &primproto.Uint128{
		Lo: lo,
		Hi: hi,
	}

	return uint128, nil
}

// strip0xPrefix removes the 0x prefix from a string if present.
func strip0xPrefix(s string) string {
	return strings.TrimPrefix(s, "0x")
}

// // addPortToURL adds a port to a URL if it doesn't already have one.
// // The port is needed for use with the Sequencer Client.
// func addPortToURL(url string) string {
// 	// Check if the URL already has a port
// 	matched, err := regexp.MatchString(`:\d+$`, url)
// 	if err != nil {
// 		log.WithError(err).Error("Error matching string")
// 		return url
// 	}
// 	if matched {
// 		log.Debug("Port already present in URL: ", url)
// 		return url
// 	}
// 	if strings.Contains(url, "http:") {
// 		log.Debug("http url detected without a port. Adding port :80 to url: ", url)
// 		return url + ":80"
// 	}
// 	if strings.Contains(url, "https:") {
// 		log.Debug("https url detected without a port. Adding port :443 to url: ", url)
// 		return url + ":443"
// 	}
// 	return url
// }

// assetIdFromDenom returns a hash of a denom string
func assetIdFromDenom(denom string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(denom))
	hash := hasher.Sum(nil)
	return hash
}

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
	address, err := EncodeBech32M(prefix, addr)
	if err != nil {
		log.WithError(err).Error("Error encoding address")
		return nil, err
	}
	return address, nil
}

// addressFromText converts a hexadecimal string representation of an address to an Address protobuf
// The input address string is expected to have the "0x" prefix stripped before being passed to this function.
// If the input string is not a valid hexadecimal string, an error will be returned.
func addressFromText(addr string) (*primproto.Address, error) {
	addr = strip0xPrefix(addr)
	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	return &primproto.Address{
		Inner: bytes,
	}, nil
}

// publicKeyFromText converts a hexadecimal string representation of a public
// key to an ed25519.PublicKey. If the input string is not a valid hexadecimal
// string, an error will be returned.
func publicKeyFromText(addr string) (ed25519.PublicKey, error) {
	addr = strip0xPrefix(addr)
	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
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

// EncodeBech32M creates a bech32m address from a [20]byte array and string
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
		Address: address,
		Prefix:  prefix,
		Bytes:   data,
	}, nil
}
