package sequencer

import (
	"crypto/ed25519"
	"testing"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	"github.com/stretchr/testify/assert"
)

func TestStrip0xPrefix(t *testing.T) {
	t.Run("with '0x' prefix", func(t *testing.T) {
		str := "0x1234abcd"
		got := strip0xPrefix(str)
		want := "1234abcd"
		assert.Equal(t, want, got)
	})

	t.Run("without '0x' prefix", func(t *testing.T) {
		str := "abcd1234"
		got := strip0xPrefix(str)
		want := "abcd1234"
		assert.Equal(t, want, got)
	})
}

func TestRollupIdFromText(t *testing.T) {
	rollupID := "steezeburger"
	actual := rollupIdFromText(rollupID)
	expected := &primproto.RollupId{
		Inner: []uint8{0x18, 0x88, 0x7, 0x48, 0xea, 0xe, 0x3c, 0xff, 0xd1, 0xcd, 0x64, 0xc1, 0xc, 0x23, 0x59, 0x31, 0xf4, 0xce, 0x4, 0x0, 0xa5, 0xae, 0xd6, 0x9c, 0x5f, 0x15, 0x57, 0x58, 0x82, 0x29, 0x9a, 0x3d},
	}
	assert.Equal(t, expected, actual)
}

func TestAddressFromPublicKey(t *testing.T) {
	// bech32m address encoded from 1c0c490f1b5528d8173c5de46d131160e4b2c0c3 bytes
	expected := "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm"

	testFromPrivKey := "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90"

	from, _ := privateKeyFromText(testFromPrivKey)
	pub := from.Public().(ed25519.PublicKey)
	actual, err := addressFromPublicKey("astria", pub)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual.Address)
}
