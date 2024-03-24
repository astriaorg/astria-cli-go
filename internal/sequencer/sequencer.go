package sequencer

import (
	"encoding/hex"

	"github.com/astriaorg/go-sequencer-client/client"
)

type Account struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

// CreateAccount creates a new account for the sequencer.
func CreateAccount() (*Account, error) {
	signer, err := client.GenerateSigner()
	if err != nil {
		return nil, err
	}
	address := signer.Address()
	seed := signer.Seed()
	return &Account{
		Address:    hex.EncodeToString(address[:]),
		PublicKey:  hex.EncodeToString(signer.PublicKey()),
		PrivateKey: hex.EncodeToString(seed[:]),
	}, nil
}
