package sequencer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astria/astria-cli-go/internal/sequencer"
)

func TestCreateAccount(t *testing.T) {
	account, err := sequencer.CreateAccount()
	assert.NoError(t, err, "CreateAccount should not return an error on success")

	assert.NotEmpty(t, account.Address, "Address should not be empty")
	assert.NotEmpty(t, account.PublicKey, "Public Key should not be empty")
	assert.NotEmpty(t, account.PrivateKey, "Private Key should not be empty")
}

func TestTransfer(t *testing.T) {
	// TODO
	assert.Equal(t, 1, 1)
}
