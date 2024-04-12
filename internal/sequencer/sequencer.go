package sequencer

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
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

type NonceResponse struct {
	Address string `json:"address"`
	Nonce   uint32 `json:"nonce"`
}

// GetNonce returns the nonce of an address.
func GetNonce(address string, sequencerURL string) (*NonceResponse, error) {
	address = strip0xPrefix(address)
	sequencerURL = addPortToURL(sequencerURL)

	log.Debug("Getting nonce for address: ", address)
	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &NonceResponse{}, err
	}

	a, err := hex.DecodeString(address)
	if err != nil {
		log.WithError(err).Error("Error decoding hex encoded address")
		return &NonceResponse{}, err
	}

	var address20 [20]byte
	copy(address20[:], a)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nonce, err := c.GetNonce(ctx, address20)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &NonceResponse{}, err
	}

	log.Debug("Nonce: ", nonce)
	return &NonceResponse{
		Address: address,
		Nonce:   nonce,
	}, nil
}

// TransferOpts are the options for the Transfer function.
type TransferOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey string
	// ToAddress is the address of the receiver
	ToAddress string
	// Amount is the amount to be transferred
	Amount int
}

// TransferResponse is the response of the Transfer function.
type TransferResponse struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Nonce  uint32 `json:"nonce"`
	Amount int    `json:"amount"`
	TxHash string `json:"txHash"`
}

// Transfer transfers an amount from one address to another.
// It returns the hash of the transaction.
func Transfer(opts TransferOpts) (*TransferResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// client
	opts.SequencerURL = addPortToURL(opts.SequencerURL)
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &TransferResponse{}, err
	}

	// create signer
	privateKeyBytes, err := hex.DecodeString(opts.FromKey)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return &TransferResponse{}, err
	}
	from := ed25519.NewKeyFromSeed(privateKeyBytes)
	signer := client.NewSigner(from)

	// create transaction
	// FIXME - support bigger numbers
	amount := &primitivev1.Uint128{
		Lo: uint64(opts.Amount),
		Hi: 0,
	}
	opts.ToAddress = strip0xPrefix(opts.ToAddress)
	to, err := hex.DecodeString(opts.ToAddress)
	if err != nil {
		log.WithError(err).Errorf("Error decoding hex encoded 'to' address %v", opts.ToAddress)
		return &TransferResponse{}, err
	}
	log.Debugf("Transferring %v to %v", opts.Amount, opts.ToAddress)
	fromAddr := signer.Address()
	nonce, err := c.GetNonce(ctx, fromAddr)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &TransferResponse{}, err
	}
	log.Debugf("Nonce: %v", nonce)
	tx := &sqproto.UnsignedTransaction{
		Nonce: nonce,
		Actions: []*sqproto.Action{
			{
				Value: &sqproto.Action_TransferAction{
					TransferAction: &sqproto.TransferAction{
						To:         to,
						Amount:     amount,
						AssetId:    AssetIdFromDenom("nria"),
						FeeAssetId: AssetIdFromDenom("nria"),
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &TransferResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &TransferResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &TransferResponse{
		From:   hex.EncodeToString(fromAddr[:]),
		To:     opts.ToAddress,
		Nonce:  nonce,
		Amount: opts.Amount,
		TxHash: hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}
