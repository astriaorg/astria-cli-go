package sequencer

import (
	"context"
	"encoding/hex"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

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

// GetBalance returns the balance of an address.
func GetBalance(address string, sequencerURL string) (*big.Int, error) {
	address = strip0xPrefix(address)
	sequencerURL = addPortToURL(sequencerURL)

	log.Debug("Getting balance for address: ", address)
	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return nil, err
	}

	a, err := hex.DecodeString(address)
	if err != nil {
		log.WithError(err).Error("Error decoding hex encoded address")
		return nil, err
	}

	var address20 [20]byte
	copy(address20[:], a)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	balance, err := c.GetBalance(ctx, address20)
	if err != nil {
		log.WithError(err).Error("Error getting balance")
		return nil, err
	}

	log.Debug("Balance: ", balance.String())
	return balance, nil
}
