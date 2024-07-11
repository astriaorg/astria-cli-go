package bech32m

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const bech32MAddress = "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm"
const bech32MAddressPrivKey = "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90"
const bech32MAddressBytes = "1c0c490f1b5528d8173c5de46d131160e4b2c0c3"
const fromPubAddress = "astria1x66v8ph5x8z95vxw6uxmyg5xahkfg0tk8lvrvf"
const pubKey = "88787e29db8d5247c6adfac9909b56e6b2705c3120b2e3885e8ec8aa416a10f1"
const prefix = "astria"

func TestEncodeFromString(t *testing.T) {
	addr, _ := EncodeFromString(bech32MAddress)
	bytes, _ := hex.DecodeString(bech32MAddressBytes)

	var len20Bytes [20]byte
	copy(len20Bytes[:], bytes)

	assert.Equal(t, bech32MAddress, addr.String())
	assert.Equal(t, prefix, addr.Prefix())
	assert.Equal(t, len20Bytes, addr.Bytes())
}

func TestVerify(t *testing.T) {
	verify := Verify(bech32MAddress)
	assert.Nil(t, verify)
}
func TestEncode(t *testing.T) {
	bytes, _ := hex.DecodeString(bech32MAddressBytes)
	var addrBytes [20]byte
	copy(addrBytes[:], bytes)
	addr, _ := Encode(prefix, addrBytes)

	assert.Equal(t, bech32MAddress, addr.String())
	assert.Equal(t, prefix, addr.Prefix())
	assert.Equal(t, addrBytes, addr.Bytes())
}

func TestEncodeFromPublicKey(t *testing.T) {
	bytes, _ := hex.DecodeString(pubKey)
	addr, _ := EncodeFromPublicKey(prefix, bytes)

	assert.Equal(t, fromPubAddress, addr.String())
	assert.Equal(t, prefix, addr.Prefix())
}

func TestEncodeFromPrivateKey(t *testing.T) {
	privBytes, _ := hex.DecodeString(bech32MAddressPrivKey)
	addr, _ := EncodeFromPrivateKey(prefix, privBytes)

	assert.Equal(t, bech32MAddress, addr.String())
	assert.Equal(t, prefix, addr.Prefix())

	bytes, _ := hex.DecodeString(bech32MAddressBytes)
	var len20Bytes [20]byte
	copy(len20Bytes[:], bytes)

	assert.Equal(t, len20Bytes, addr.Bytes())
}
