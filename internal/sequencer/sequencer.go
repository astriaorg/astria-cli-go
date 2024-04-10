package sequencer

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	primitivev1 "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/primitive/v1"
	log "github.com/sirupsen/logrus"

	sqproto "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1"
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

// GetBlockheight returns the current blockheight of the sequencer.
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

// TransferOpts are the options for the Transfer function.
type TransferOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey []byte
	// ToAddress is the address of the receiver
	ToAddress []byte
	// Amount is the amount to be transferred
	Amount int
}

// Transfer transfers an amount from one address to another.
func Transfer(opts TransferOpts) error {
	opts.SequencerURL = addPortToURL(opts.SequencerURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// transaction
	// FIXME - support bigger numbers
	amount := &primitivev1.Uint128{
		Lo: uint64(opts.Amount),
		Hi: 0,
	}
	tx := &sqproto.UnsignedTransaction{
		Nonce: 1, // FIXME - do i need to get the nonce first? where from?
		Actions: []*sqproto.Action{
			{
				Value: &sqproto.Action_TransferAction{
					TransferAction: &sqproto.TransferAction{
						To:         opts.ToAddress,
						Amount:     amount,
						AssetId:    []byte("nria"), // FIXME - is this right? put in constant
						FeeAssetId: []byte("nria"), // FIXME - is this right? put in constant
					},
				},
			},
		},
	}

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return err
	}

	// signer
	signer := client.NewSigner(opts.FromKey)
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return err
	}

	// FIXME - return something else here and handle output in the cmd's code
	fmt.Println(resp)
	return nil
}
