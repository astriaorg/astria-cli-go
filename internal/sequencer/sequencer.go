package sequencer

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"time"

	log "github.com/sirupsen/logrus"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"

	"github.com/astriaorg/go-sequencer-client/client"
	log "github.com/sirupsen/logrus"
)

const (
	// DefaultSequencerNetworkId is the default network id for the sequencer.
	DefaultSequencerNetworkId = "astria-dusk-5"
)

// CreateAccount creates a new account for the sequencer.
func CreateAccount() (*Account, error) {
	signer, err := client.GenerateSigner()
	if err != nil {
		log.WithError(err).Error("Failed to generate signer")
		return nil, err
	}
	address := signer.Address()
	seed := signer.Seed()

	addr := hex.EncodeToString(address[:])
	priv := ed25519.NewKeyFromSeed(seed[:])
	pub := priv.Public().(ed25519.PublicKey)

	log.Debug("Created account with address: ", addr)
	return &Account{
		Address:    addr,
		PublicKey:  pub,
		PrivateKey: priv,
	}, nil
}

// GetBalances returns the balances of an address.
func GetBalances(address string, sequencerURL string) (*BalancesResponse, error) {
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

	// convert to our BalancesResponse type
	b := make(BalancesResponse, len(balances))
	for i, balance := range balances {
		b[i] = &Balance{
			Denom:   balance.Denom,
			Balance: balance.Balance,
		}
	}
	return &b, nil
}

// GetBlockheight returns the current blockheight of the sequencer.
func GetBlockheight(sequencerURL string) (*BlockheightResponse, error) {
	sequencerURL = addPortToURL(sequencerURL)

	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &BlockheightResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blockheight, err := c.GetBlockHeight(ctx)
	if err != nil {
		log.WithError(err).Error("Error getting blockheight")
		return &BlockheightResponse{}, err
	}

	log.Debug("Blockheight: ", blockheight)
	return &BlockheightResponse{
		Blockheight: blockheight,
	}, nil
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
	amount, err := convertToUint128(opts.Amount)
	if err != nil {
		log.WithError(err).Error("Error converting amount to Uint128 proto")
		return &TransferResponse{}, err
	}
	opts.ToAddress = strip0xPrefix(opts.ToAddress)
	toBytes, err := hex.DecodeString(opts.ToAddress)
	if err != nil {
		log.WithError(err).Errorf("Error decoding hex encoded 'to' address %v", opts.ToAddress)
		return &TransferResponse{}, err
	}
	to := &primproto.Address{
		Inner: toBytes,
	}
	log.Debugf("Transferring %v to %v", opts.Amount, opts.ToAddress)
	fromAddr := signer.Address()
	nonce, err := c.GetNonce(ctx, fromAddr)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &TransferResponse{}, err
	}
	log.Debugf("Nonce: %v", nonce)
	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: DefaultSequencerNetworkId,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_TransferAction{
					TransferAction: &txproto.TransferAction{
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

func InitBridgeAccount(opts InitBridgeOpts) (*InitBridgeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// sha256 hash of rollupID
	hashedRollupID := RollupIdFromText(opts.RollupID)

	// client
	opts.SequencerURL = addPortToURL(opts.SequencerURL)
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &InitBridgeResponse{}, err
	}

	// create signer
	privateKeyBytes, err := hex.DecodeString(opts.FromKey)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return &InitBridgeResponse{}, err
	}
	from := ed25519.NewKeyFromSeed(privateKeyBytes)
	signer := client.NewSigner(from)

	// Get current address nonce
	fromAddr := signer.Address()
	nonce, err := c.GetNonce(ctx, fromAddr)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &InitBridgeResponse{}, err
	}

	// build transaction
	tx := &sqproto.UnsignedTransaction{
		Nonce: nonce,
		Actions: []*sqproto.Action{
			{
				Value: &sqproto.Action_InitBridgeAccountAction{
					InitBridgeAccountAction: &sqproto.InitBridgeAccountAction{
						RollupId:   hashedRollupID,
						AssetIds:   [][]byte{AssetIdFromDenom("nria")},
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
		return &InitBridgeResponse{}, err
	}

	// broadcast transaction
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &InitBridgeResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &InitBridgeResponse{
		RollupID: opts.RollupID,
		Nonce:    nonce,
		TxHash:   hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil

}

func BridgeLock(opts BridgeLockOpts) (*BridgeLockResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// client
	opts.SequencerURL = addPortToURL(opts.SequencerURL)
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &BridgeLockResponse{}, err
	}

	// create signer
	privateKeyBytes, err := hex.DecodeString(opts.FromKey)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return &BridgeLockResponse{}, err
	}
	from := ed25519.NewKeyFromSeed(privateKeyBytes)
	signer := client.NewSigner(from)

	// Get current address nonce
	fromAddr := signer.Address()
	nonce, err := c.GetNonce(ctx, fromAddr)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &BridgeLockResponse{}, err
	}

	// create transaction
	amount, err := convertToUint128(opts.Amount)
	if err != nil {
		log.WithError(err).Error("Error converting amount to Uint128 proto")
		return &BridgeLockResponse{}, err
	}
	opts.ToAddress = strip0xPrefix(opts.ToAddress)
	to, err := hex.DecodeString(opts.ToAddress)
	if err != nil {
		log.WithError(err).Errorf("Error decoding hex encoded 'to' address %v", opts.ToAddress)
		return &BridgeLockResponse{}, err
	}

	tx := &sqproto.UnsignedTransaction{
		Nonce: nonce,
		Actions: []*sqproto.Action{
			{
				Value: &sqproto.Action_BridgeLockAction{
					BridgeLockAction: &sqproto.BridgeLockAction{
						To:                      to,
						Amount:                  amount,
						AssetId:                 AssetIdFromDenom("nria"),
						FeeAssetId:              AssetIdFromDenom("nria"),
						DestinationChainAddress: opts.DestinationChain,
					},
				},
			},
		},
	}
	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &BridgeLockResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &BridgeLockResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &BridgeLockResponse{
		From:   hex.EncodeToString(fromAddr[:]),
		To:     opts.ToAddress,
		Nonce:  nonce,
		Amount: opts.Amount,
		TxHash: hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil

}
