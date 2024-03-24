package sequencer

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/astriaorg/go-sequencer-client/client"
)

type Account struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

// CreateAccount creates a new account for the sequencer.
func CreateAccount() *Account {
	signer, err := client.GenerateSigner()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	address := signer.Address()
	seed := signer.Seed()
	return &Account{
		Address:    hex.EncodeToString(address[:]),
		PublicKey:  hex.EncodeToString(signer.PublicKey()),
		PrivateKey: hex.EncodeToString(seed[:]),
	}
}
