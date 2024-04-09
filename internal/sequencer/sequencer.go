package sequencer

import (
	"context"
	"encoding/hex"
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

	log.Debug("Created account with address: ", hex.EncodeToString(address[:]))
	return &Account{
		Address:    hex.EncodeToString(address[:]),
		PublicKey:  hex.EncodeToString(signer.PublicKey()),
		PrivateKey: hex.EncodeToString(seed[:]),
	}, nil
}

// GetBalances returns the balances of an address.
func GetBalances(address string, sequencerURL string) ([]*client.BalanceResponse, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balances, err := c.GetBalances(ctx, address20)
	if err != nil {
		log.WithError(err).Error("Error getting balance")
		return nil, err
	}

	for _, b := range balances {
		log.Debug("Denom:", b.Denom, "Balance:", b.Balance.String())
	}
	return balances, nil
}

func GetBlockheight(sequencerURL string) (int64, error) {
	sequencerURL = addPortToURL(sequencerURL)

	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blockheight, err := c.GetBlockHeight(ctx)
	if err != nil {
		log.WithError(err).Error("Error getting blockheight")
		return 0, err
	}

	log.Debug("Blockheight: ", blockheight)
	return blockheight, nil
}
