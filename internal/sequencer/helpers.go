package sequencer

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"

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

// addPortToURL adds a port to a URL if it doesn't already have one.
// The port is needed for use with the Sequencer Client.
func addPortToURL(url string) string {
	// Check if the URL already has a port
	matched, err := regexp.MatchString(`:\d+$`, url)
	if err != nil {
		log.WithError(err).Error("Error matching string")
		return url
	}
	if matched {
		log.Debug("Port already present in URL: ", url)
		return url
	}
	if strings.Contains(url, "http:") {
		log.Debug("http url detected without a port. Adding port :80 to url: ", url)
		return url + ":80"
	}
	if strings.Contains(url, "https:") {
		log.Debug("https url detected without a port. Adding port :443 to url: ", url)
		return url + ":443"
	}
	return url
}

// AssetIdFromDenom returns a hash of a denom string
func AssetIdFromDenom(denom string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(denom))
	hash := hasher.Sum(nil)
	return hash
}

// RollupIdFromText converts a string to a RollupId protobuf.
func RollupIdFromText(rollup string) *primproto.RollupId {
	hasher := sha256.New()
	hasher.Write([]byte(rollup))
	hash := hasher.Sum(nil)
	return &primproto.RollupId{
		Inner: hash,
	}
}

func AddressFromPublicKey(pubkey ed25519.PublicKey) string {
	hash := sha256.Sum256(pubkey)
	var addr [20]byte
	copy(addr[:], hash[:20])
	return hex.EncodeToString(addr[:])
}
