package sequencer_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
)

func TestCreateAccount(t *testing.T) {
	account, err := sequencer.CreateAccount("astria")
	assert.NoError(t, err, "CreateAccount should not return an error on success")

	assert.NotEmpty(t, account.Address, "Address should not be empty")
	assert.NotEmpty(t, account.PublicKey, "Public Key should not be empty")
	assert.NotEmpty(t, account.PrivateKey, "Private Key should not be empty")

	assert.Equal(t, account.Address, account.ToJSONStruct().Address, "Address should match JSON representation")
	assert.Equal(t, account.PublicKeyString(), account.ToJSONStruct().PublicKey, "Public Key should match JSON representation")
	assert.Equal(t, account.PrivateKeyString(), account.ToJSONStruct().PrivateKey, "Private Key should match JSON representation")

	// test our logic even though they're one line conversions
	assert.Equal(t, hex.EncodeToString(account.PrivateKey[:32]), account.PrivateKeyString(), "Private Key string should be hex encoded last 32 bytes of PrivateKey")
	assert.Equal(t, hex.EncodeToString(account.PublicKey), account.PublicKeyString(), "Public Key string should be hex encoded bytes of PublicKey")
}
