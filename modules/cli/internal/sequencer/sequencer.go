package sequencer

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	"github.com/astriaorg/astria-cli-go/modules/bech32m"
	"github.com/astriaorg/astria-cli-go/modules/go-sequencer-client/client"
	log "github.com/sirupsen/logrus"
)

// CreateAccount creates a new account for the sequencer. The address will be a
// bech32m encoded string, which is created using the prefix provided.
func CreateAccount(prefix string) (*Account, error) {
	signer, err := client.GenerateSigner()
	if err != nil {
		log.WithError(err).Error("Failed to generate signer")
		return nil, err
	}
	address := signer.Address()
	seed := signer.Seed()

	log.Debugf("Address as hex: %x", address[:])

	addr, err := bech32m.EncodeFromBytes(prefix, address)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
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
	log.Debug("Getting balance for address: ", address)
	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	balances, err := c.GetBalances(ctx, address)
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

// GetBlock returns the specific block from the sequencer.
func GetBlock(opts BlockOpts) (*BlockResponse, error) {
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)

	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &BlockResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	block, err := c.GetBlock(ctx, &opts.BlockHeight)
	if err != nil {
		log.WithError(err).Error("Error getting blockheight")
		return &BlockResponse{}, err
	}

	log.Debug("Retrieved Block at block height: ", opts.BlockHeight)
	return &BlockResponse{
		Block: block,
	}, nil
}

// GetBlockheight returns the current blockheight of the sequencer.
func GetBlockheight(sequencerURL string) (*BlockheightResponse, error) {
	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &BlockheightResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	log.Debug("Getting nonce for address: ", address)
	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &NonceResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nonce, err := c.GetNonce(ctx, address)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &TransferResponse{}, err
	}

	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := bech32m.EncodeFromBytes(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.String())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &TransferResponse{}, err
	}
	log.Debugf("Nonce: %v", nonce)

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_TransferAction{
					TransferAction: &txproto.TransferAction{
						To:       opts.ToAddress,
						Amount:   opts.Amount,
						Asset:    opts.Asset,
						FeeAsset: opts.FeeAsset,
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
	resp, err := c.BroadcastTx(ctx, signed, opts.IsAsync)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &TransferResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)
	// response
	hash := hex.EncodeToString(resp.Hash)
	amount := fmt.Sprint(client.ProtoU128ToBigInt(opts.Amount))
	tr := &TransferResponse{
		From:   addr.String(),
		To:     opts.ToAddress.Bech32M,
		Nonce:  nonce,
		Amount: amount,
		TxHash: hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}

// IbcTransfer performs an ICS20 withdrawal from the sequencer to a recipient on another chain.
func IbcTransfer(opts IbcTransferOpts) (*IbcTransferResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &IbcTransferResponse{}, err
	}

	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := bech32m.EncodeFromBytes(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.String())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &IbcTransferResponse{}, err
	}
	log.Debugf("Nonce: %v", nonce)

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_Ics20Withdrawal{
					Ics20Withdrawal: &txproto.Ics20Withdrawal{
						Amount:                  opts.Amount,
						Denom:                   opts.Asset,
						DestinationChainAddress: opts.DestinationChainAddressAddress,
						ReturnAddress:           opts.ReturnAddress,
						TimeoutHeight: &txproto.IbcHeight{
							RevisionNumber: math.MaxUint64,
							RevisionHeight: math.MaxUint64,
						},
						TimeoutTime:   nowPlusFiveMinutes(),
						SourceChannel: opts.SourceChannelID,
						FeeAsset:      opts.FeeAsset,
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &IbcTransferResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTx(ctx, signed, opts.IsAsync)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &IbcTransferResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)
	// response
	hash := hex.EncodeToString(resp.Hash)
	amount := fmt.Sprint(client.ProtoU128ToBigInt(opts.Amount))
	tr := &IbcTransferResponse{
		From:   addr.String(),
		To:     opts.DestinationChainAddressAddress,
		Nonce:  nonce,
		Amount: amount,
		TxHash: hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}

func InitBridgeAccount(opts InitBridgeOpts) (*InitBridgeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &InitBridgeResponse{}, err
	}

	// Get current address nonce
	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := bech32m.EncodeFromBytes(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.String())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &InitBridgeResponse{}, err
	}

	// build transaction
	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_InitBridgeAccountAction{
					InitBridgeAccountAction: &txproto.InitBridgeAccountAction{
						RollupId:          rollupIdFromText(opts.RollupName),
						Asset:             opts.Asset,
						FeeAsset:          opts.FeeAsset,
						SudoAddress:       opts.SudoAddress,
						WithdrawerAddress: opts.WithdrawerAddress,
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
	resp, err := c.BroadcastTx(ctx, signed, opts.IsAsync)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &InitBridgeResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &InitBridgeResponse{
		RollupID: opts.RollupName,
		Nonce:    nonce,
		TxHash:   hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil

}

// BridgeLock locks tokens on the source chain and initiates a cross-chain transfer to the destination chain.
func BridgeLock(opts BridgeLockOpts) (*BridgeLockResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Debugf("BridgeLockOpts: %v", opts)

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &BridgeLockResponse{}, err
	}

	// Get current address nonce
	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := bech32m.EncodeFromBytes(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.String())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &BridgeLockResponse{}, err
	}

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_BridgeLockAction{
					BridgeLockAction: &txproto.BridgeLockAction{
						To:                      opts.ToAddress,
						Amount:                  opts.Amount,
						Asset:                   opts.Asset,
						FeeAsset:                opts.FeeAsset,
						DestinationChainAddress: opts.DestinationChainAddress,
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
	resp, err := c.BroadcastTx(ctx, signed, opts.IsAsync)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &BridgeLockResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &BridgeLockResponse{
		From:   addr.String(),
		To:     opts.ToAddress.Bech32M,
		Nonce:  nonce,
		Amount: opts.Amount.String(),
		TxHash: hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil

}
