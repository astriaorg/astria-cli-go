package bech32m

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const bech32MAddress = "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm"
const bech32MAddressBytes = "1c0c490f1b5528d8173c5de46d131160e4b2c0c3"
const fromPubAddress = "astria1x66v8ph5x8z95vxw6uxmyg5xahkfg0tk8lvrvf"
const pubKey = "88787e29db8d5247c6adfac9909b56e6b2705c3120b2e3885e8ec8aa416a10f1"
const testPrefix = "astria"

func TestValidate(t *testing.T) {
	err := Validate(bech32MAddress)
	assert.Nil(t, err)
}
func TestEncode(t *testing.T) {
	bytes, _ := hex.DecodeString(bech32MAddressBytes)
	var addrBytes [20]byte
	copy(addrBytes[:], bytes)
	addr, _ := EncodeFromBytes(testPrefix, addrBytes)

	assert.Equal(t, bech32MAddress, addr.String())
	assert.Equal(t, testPrefix, addr.Prefix())
	assert.Equal(t, addrBytes, addr.Bytes())
}

func TestEncodeFromPublicKey(t *testing.T) {
	bytes, _ := hex.DecodeString(pubKey)
	addr, _ := EncodeFromPublicKey(testPrefix, bytes)

	assert.Equal(t, fromPubAddress, addr.String())
	assert.Equal(t, testPrefix, addr.Prefix())
}

func TestDecode(t *testing.T) {
	prefix, addr, err := Decode(bech32MAddress)
	assert.Nil(t, err)
	assert.Equal(t, prefix, testPrefix)
	assert.Equal(t, bech32MAddressBytes, hex.EncodeToString(addr[:]))
}
