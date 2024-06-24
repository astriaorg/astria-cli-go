package sequencer

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"strconv"
	"time"

	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	"buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria_vendored/tendermint/abci"
	"buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria_vendored/tendermint/crypto"
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

	addr, err := EncodeBech32M(prefix, address)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	log.Debug("Getting nonce for address: ", address)
	log.Debug("Creating CometBFT client with url: ", sequencerURL)

	c, err := client.NewClient(sequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &NonceResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
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
						To:         opts.ToAddress,
						Amount:     opts.Amount,
						AssetId:    opts.AssetID,
						FeeAssetId: opts.FeeAssetID,
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
		From:   addr.ToString(),
		To:     opts.ToAddress.Bech32M,
		Nonce:  nonce,
		Amount: opts.Amount.String(),
		TxHash: hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}

func InitBridgeAccount(opts InitBridgeOpts) (*InitBridgeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
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
						AssetId:           opts.AssetID,
						FeeAssetId:        opts.FeeAssetID,
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
	resp, err := c.BroadcastTxSync(ctx, signed)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
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
						AssetId:                 opts.AssetID,
						FeeAssetId:              opts.FeeAssetID,
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
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &BridgeLockResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &BridgeLockResponse{
		From:   addr.ToString(),
		To:     opts.ToAddress.Bech32M,
		Nonce:  nonce,
		Amount: opts.Amount.String(),
		TxHash: hash,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil

}

// AddFeeAsset adds a fee asset to the sequencer.
func AddFeeAsset(opts FeeAssetOpts) (*FeeAssetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Debugf("AddFeeAssetOpts: %v", opts)

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &FeeAssetResponse{}, err
	}

	// Get current address nonce
	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &FeeAssetResponse{}, err
	}

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_FeeAssetChangeAction{
					FeeAssetChangeAction: &txproto.FeeAssetChangeAction{
						Value: &txproto.FeeAssetChangeAction_Addition{
							Addition: assetIdFromDenom(opts.Asset),
						},
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &FeeAssetResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &FeeAssetResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &FeeAssetResponse{
		From:       addr.ToString(),
		Nonce:      nonce,
		TxHash:     hash,
		FeeAssetId: opts.Asset,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}

// RemoveFeeAsset removes a fee asset from the sequencer.
func RemoveFeeAsset(opts FeeAssetOpts) (*FeeAssetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Debugf("RemoveFeeAssetOpts: %v", opts)

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &FeeAssetResponse{}, err
	}

	// Get current address nonce
	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &FeeAssetResponse{}, err
	}

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_FeeAssetChangeAction{
					FeeAssetChangeAction: &txproto.FeeAssetChangeAction{
						Value: &txproto.FeeAssetChangeAction_Removal{
							Removal: assetIdFromDenom(opts.Asset),
						},
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &FeeAssetResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &FeeAssetResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &FeeAssetResponse{
		From:       addr.ToString(),
		Nonce:      nonce,
		TxHash:     hash,
		FeeAssetId: opts.Asset,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}

// AddIBCRelayer adds an IBC Relayer address to the sequencer.
func AddIBCRelayer(opts IBCRelayerOpts) (*IBCRelayerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Debugf("AddIBCRelayerOpts: %v", opts)

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &IBCRelayerResponse{}, err
	}

	// Get current address nonce
	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &IBCRelayerResponse{}, err
	}

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_IbcRelayerChangeAction{
					IbcRelayerChangeAction: &txproto.IbcRelayerChangeAction{
						Value: &txproto.IbcRelayerChangeAction_Addition{
							Addition: opts.IBCRelayerAddress,
						},
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &IBCRelayerResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &IBCRelayerResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &IBCRelayerResponse{
		From:              addr.ToString(),
		Nonce:             nonce,
		TxHash:            hash,
		IBCRelayerAddress: opts.IBCRelayerAddress.Bech32M,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}

// RemoveIBCRelayer removes an IBC Relayer address from the sequencer.
func RemoveIBCRelayer(opts IBCRelayerOpts) (*IBCRelayerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Debugf("RemoveIBCRelayerOpts: %v", opts)

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &IBCRelayerResponse{}, err
	}

	// Get current address nonce
	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &IBCRelayerResponse{}, err
	}

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_IbcRelayerChangeAction{
					IbcRelayerChangeAction: &txproto.IbcRelayerChangeAction{
						Value: &txproto.IbcRelayerChangeAction_Removal{
							Removal: opts.IBCRelayerAddress,
						},
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &IBCRelayerResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &IBCRelayerResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &IBCRelayerResponse{
		From:              addr.ToString(),
		Nonce:             nonce,
		TxHash:            hash,
		IBCRelayerAddress: opts.IBCRelayerAddress.Bech32M,
	}

	log.Debugf("Transfer hash: %v", hash)
	return tr, nil
}

// ChangeSudoAddress changes the sudo address.
func ChangeSudoAddress(opts ChangeSudoAddressOpts) (*ChangeSudoAddressResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Debugf("Change Sudo Address Opts: %v", opts)

	// client
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &ChangeSudoAddressResponse{}, err
	}

	// Get current address nonce
	signer := client.NewSigner(opts.FromKey)
	fromAddr := signer.Address()
	addr, err := EncodeBech32M(opts.AddressPrefix, fromAddr)
	if err != nil {
		log.WithError(err).Error("Failed to encode address")
		return nil, err
	}
	nonce, err := c.GetNonce(ctx, addr.ToString())
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &ChangeSudoAddressResponse{}, err
	}

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_SudoAddressChangeAction{
					SudoAddressChangeAction: &txproto.SudoAddressChangeAction{
						NewAddress: opts.UpdateAddress,
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &ChangeSudoAddressResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &ChangeSudoAddressResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &ChangeSudoAddressResponse{
		From:           addr.ToString(),
		Nonce:          nonce,
		NewSudoAddress: opts.UpdateAddress.Bech32M,
		TxHash:         hash,
	}

	log.Debugf("Change Sudo Address TX hash: %v", hash)
	return tr, nil
}

// UpdateValidator changes the power of a validator.
func UpdateValidator(opts UpdateValidatorOpts) (*UpdateValidatorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Debugf("Update Validator Opts: %v", opts)

	// client
	opts.SequencerURL = addPortToURL(opts.SequencerURL)
	log.Debug("Creating CometBFT client with url: ", opts.SequencerURL)
	c, err := client.NewClient(opts.SequencerURL)
	if err != nil {
		log.WithError(err).Error("Error creating sequencer client")
		return &UpdateValidatorResponse{}, err
	}

	// create signer
	from, err := privateKeyFromText(opts.FromKey)
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		return &UpdateValidatorResponse{}, err
	}
	signer := client.NewSigner(from)

	// Get current address nonce
	fromAddr := signer.Address()
	nonce, err := c.GetNonce(ctx, fromAddr)
	if err != nil {
		log.WithError(err).Error("Error getting nonce")
		return &UpdateValidatorResponse{}, err
	}

	// decode public key
	pk, err := publicKeyFromText(opts.PubKey)
	if err != nil {
		log.WithError(err).Errorf("Error decoding hex encoded public key %v", opts.PubKey)
		return &UpdateValidatorResponse{}, err
	}
	pubKey := &crypto.PublicKey{
		Sum: &crypto.PublicKey_Ed25519{
			Ed25519: pk,
		},
	}

	power, err := strconv.ParseInt(opts.Power, 10, 64)
	if err != nil {
		log.WithError(err).Errorf("Error decoding power string to int64 %v", opts.Power)
		return &UpdateValidatorResponse{}, err
	}

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: opts.SequencerChainID,
			Nonce:   nonce,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_ValidatorUpdateAction{
					ValidatorUpdateAction: &abci.ValidatorUpdate{
						PubKey: pubKey,
						Power:  power,
					},
				},
			},
		},
	}

	// sign transaction
	signed, err := signer.SignTransaction(tx)
	if err != nil {
		log.WithError(err).Error("Error signing transaction")
		return &UpdateValidatorResponse{}, err
	}

	// broadcast tx
	resp, err := c.BroadcastTxSync(ctx, signed)
	if err != nil {
		log.WithError(err).Error("Error broadcasting transaction")
		return &UpdateValidatorResponse{}, err
	}
	log.Debugf("Broadcast response: %v", resp)

	// response
	hash := hex.EncodeToString(resp.Hash)
	tr := &UpdateValidatorResponse{
		From:   hex.EncodeToString(fromAddr[:]),
		Nonce:  nonce,
		PubKey: opts.PubKey,
		Power:  opts.Power,
		TxHash: hash,
	}
	log.Debug(tr)

	log.Debugf("Update Validator TX hash: %v", hash)
	return tr, nil
}
