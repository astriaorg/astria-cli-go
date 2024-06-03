package sequencer

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"
)

// Account is the struct that holds the account information.
type Account struct {
	Address    string
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// NewAccountFromPrivKey creates a new Account struct from a given private key.
// It calculates the public key from the private key, generates the address from the public key,
// and returns a pointer to the new Account struct with the address, public key, and private key set.
func NewAccountFromPrivKey(privkey ed25519.PrivateKey) *Account {
	pub := privkey.Public().(ed25519.PublicKey)
	return &Account{
		Address:    addressFromPublicKey(pub),
		PublicKey:  pub,
		PrivateKey: privkey,
	}
}

// AccountJSON is for representing an `Account` as JSON
type AccountJSON struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// ToJSONStruct converts an Account into an AccountJSON struct for JSON representation.
func (a *Account) ToJSONStruct() *AccountJSON {
	return &AccountJSON{
		Address:    a.Address,
		PublicKey:  a.PublicKeyString(),
		PrivateKey: a.PrivateKeyString(),
	}
}

// PublicKeyString hex encodes the public key bytes.
func (a *Account) PublicKeyString() string {
	return hex.EncodeToString(a.PublicKey)
}

// PrivateKeyString hex encodes the last 32 bytes of the Private Key.
// FIXME - why last 32 bytes?
func (a *Account) PrivateKeyString() string {
	// NOTE - if the private key is empty we can assume we're not printing it for a reason
	if len(a.PrivateKey) == 0 {
		return "[REDACTED]"
	}
	return hex.EncodeToString(a.PrivateKey[:32])
}

func (a *Account) JSON() ([]byte, error) {
	accountJSON := a.ToJSONStruct()
	return json.MarshalIndent(accountJSON, "", "  ")
}

func (a *Account) TableHeader() []string {
	return []string{"Address", "Public Key", "Private Key"}
}

func (a *Account) TableRows() [][]string {
	return [][]string{
		{a.Address, a.PublicKeyString(), a.PrivateKeyString()},
	}
}

// BalancesResponse is the response of the GetBalances function.
type BalancesResponse []*Balance

// Balance is the balance of an asset.
type Balance struct {
	Denom   string   `json:"denom"`
	Balance *big.Int `json:"balance"`
}

func (br *BalancesResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(br, "", "  ")
}

func (br *BalancesResponse) TableHeader() []string {
	return []string{"Denom", "Balance"}
}

func (br *BalancesResponse) TableRows() [][]string {
	rows := make([][]string, len(*br))
	for i, balance := range *br {
		rows[i] = []string{balance.Denom, balance.Balance.String()}
	}
	return rows
}

// BlockheightResponse is the response of the GetBlockheight function.
type BlockheightResponse struct {
	Blockheight int64 `json:"blockheight"` // NOTE - cometbft returns int64 for this
}

func (br *BlockheightResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(br, "", "  ")

}

func (br *BlockheightResponse) TableHeader() []string {
	return []string{"Blockheight"}
}

func (br *BlockheightResponse) TableRows() [][]string {
	return [][]string{
		{strconv.Itoa(int(br.Blockheight))},
	}
}

// NonceResponse is the response of the GetNonce function.
type NonceResponse struct {
	Address string `json:"address"`
	Nonce   uint32 `json:"nonce"`
}

func (nr *NonceResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(nr, "", "  ")
}

func (nr *NonceResponse) TableHeader() []string {
	return []string{"Address", "Nonce"}
}

func (nr *NonceResponse) TableRows() [][]string {
	return [][]string{
		{nr.Address, strconv.Itoa(int(nr.Nonce))},
	}
}

// InitBridgeOpts are the options for the InitBridge function.
type InitBridgeOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// fromKey is the private key of the sender
	FromKey string
	// RollupID is the ID of the rollup to create the bridge account for
	RollupID string
	// SequencerChainID is the ID of the sequencer chain to create the bridge account on
	SequencerChainID string
	// AssetID is the name of the asset to bridge
	AssetID string
	// FeeAssetID is the name of the fee asset to use for the transaction fee
	FeeAssetID string
}
type InitBridgeResponse struct {
	RollupID string `json:"rollupID"`
	Nonce    uint32 `json:"nonce"`
	TxHash   string `json:"txHash"`
}

func (nr *InitBridgeResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(nr, "", "  ")
}

func (nr *InitBridgeResponse) TableHeader() []string {
	return []string{"RollupId", "Nonce", "TxHash"}
}

func (nr *InitBridgeResponse) TableRows() [][]string {
	return [][]string{
		{nr.RollupID, strconv.Itoa(int(nr.Nonce)), nr.TxHash},
	}
}

// BridgeLockOpts are the options for the BridgeLock function.
type BridgeLockOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey string
	// ToAddress is the address of the receiver
	ToAddress string
	// SequencerChainID is the ID of the sequencer chain to lock asset on
	SequencerChainID string
	// AssetID is the name of the asset to lock
	AssetID string
	// FeeAssetID is the name of the asset to use for the transaction fee
	FeeAssetID string
	// Amount is the amount to be locked
	Amount string
	// DestinationChainAddress is the address on the destination chain
	DestinationChainAddress string
}

// BridgeLockResponse is the response of the BridgeLock function.
type BridgeLockResponse struct {
	// From is the address of the sender
	From string `json:"from"`
	// To is the address of the receiver. For a bridge lock, this is the bridge account
	To string `json:"to"`
	// Amount is the amount locked
	Amount string `json:"amount"`
	// Nonce is the nonce of the transaction
	Nonce uint32 `json:"nonce"`
	// TxHash is the hash of the transaction
	TxHash string `json:"txHash"`
}

func (nr *BridgeLockResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(nr, "", "  ")
}

func (nr *BridgeLockResponse) TableHeader() []string {
	return []string{"From", "To", "Amount", "Nonce", "TxHash"}
}

func (nr *BridgeLockResponse) TableRows() [][]string {
	return [][]string{
		{nr.From, nr.To, nr.Amount, strconv.Itoa(int(nr.Nonce)), nr.TxHash},
	}
}

// TransferOpts are the options for the Transfer function.
type TransferOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey string
	// ToAddress is the address of the receiver
	ToAddress string
	// Amount is the amount to be transferred. Using string type to support huge numbers
	Amount string
	// SequencerChainID is the chain ID of the sequencer
	SequencerChainID string
}

// TransferResponse is the response of the Transfer function.
type TransferResponse struct {
	// From is the address of the sender
	From string `json:"from"`
	// To is the address of the receiver
	To string `json:"to"`
	// Amount is the amount transferred
	Amount string `json:"amount"`
	// Nonce is the nonce of the transaction
	Nonce uint32 `json:"nonce"`
	// TxHash is the hash of the transaction
	TxHash string `json:"txHash"`
}

func (tr *TransferResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(tr, "", "  ")
}

func (tr *TransferResponse) TableHeader() []string {
	return []string{"From", "To", "Amount", "Nonce", "TxHash"}
}

func (tr *TransferResponse) TableRows() [][]string {
	return [][]string{
		{tr.From, tr.To, strconv.Itoa(int(tr.Nonce)), tr.Amount, tr.TxHash},
	}
}

type FeeAssetOpts struct {
	// FromKey is the private key of the sender
	FromKey string
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// SequencerChainID is the chain ID of the sequencer
	SequencerChainID string
	// Asset is the fee asset that will be added or removed
	Asset string
}

type FeeAssetResponse struct {
	// From is the address of the sender
	From string `json:"from"`
	// Nonce is the nonce of the transaction
	Nonce uint32 `json:"nonce"`
	// TxHash is the hash of the transaction
	TxHash string `json:"txHash"`
	// FeeAssetId is the asset id of the fee asset
	FeeAssetId string `json:"feeAssetId"`
}

func (far *FeeAssetResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(far, "", "  ")
}

func (far *FeeAssetResponse) TableHeader() []string {
	return []string{"From", "Nonce", "TxHash", "FeeAssetId"}
}

func (far *FeeAssetResponse) TableRows() [][]string {
	return [][]string{
		{far.From, strconv.Itoa(int(far.Nonce)), far.TxHash, far.FeeAssetId},
	}
}

type IBCRelayerOpts struct {
	// FromKey is the private key of the sender
	FromKey string
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// SequencerChainID is the chain ID of the sequencer
	SequencerChainID string
	// IBCRelayerAddress is the ibc relayer address that will be added or removed
	IBCRelayerAddress string
}

type IBCRelayerResponse struct {
	// From is the address of the sender
	From string `json:"from"`
	// Nonce is the nonce of the transaction
	Nonce uint32 `json:"nonce"`
	// TxHash is the hash of the transaction
	TxHash string `json:"txHash"`
	// IBCRelayerAddress is the asset id of the fee asset
	IBCRelayerAddress string `json:"ibcRelayerAddress"`
}

func (i *IBCRelayerResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(i, "", "  ")
}

func (i *IBCRelayerResponse) TableHeader() []string {
	return []string{"From", "Nonce", "TxHash", "IBCRelayerAddress"}
}

func (i *IBCRelayerResponse) TableRows() [][]string {
	return [][]string{
		{i.From, strconv.Itoa(int(i.Nonce)), i.TxHash, i.IBCRelayerAddress},
	}
}

type MintOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey string
	// ToAddress is the address of the receiver
	ToAddress string
	// Amount is the amount to be transferred. Using string type to support huge numbers
	Amount string
	// SequencerChainID is the chain ID of the sequencer
	SequencerChainID string
}

type MintResponse struct {
	// From is the address of the sender
	From string `json:"from"`
	// Nonce is the nonce of the transaction
	Nonce uint32 `json:"nonce"`
	// To is the address of the receiver
	To string `json:"to"`
	// Amount is the amount transferred
	Amount string `json:"amount"`
	// TxHash is the hash of the transaction
	TxHash string `json:"txHash"`
}

func (m *MintResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

func (m *MintResponse) TableHeader() []string {
	return []string{"From", "Nonce", "To", "Amount", "TxHash"}
}

func (m *MintResponse) TableRows() [][]string {
	return [][]string{
		{m.From, strconv.Itoa(int(m.Nonce)), m.To, m.Amount, m.TxHash},
	}
}

type ChangeSudoAddressOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey string
	// ToAddress is the address of the receiver
	UpdateAddress string
	// SequencerChainID is the chain ID of the sequencer
	SequencerChainID string
}

type ChangeSudoAddressResponse struct {
	// From is the address of the sender
	From string `json:"from"`
	// Nonce is the nonce of the transaction
	Nonce uint32 `json:"nonce"`
	// To is the address of the receiver
	NewSudoAddress string `json:"newSudoAddress"`
	// TxHash is the hash of the transaction
	TxHash string `json:"txHash"`
}

func (c *ChangeSudoAddressResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *ChangeSudoAddressResponse) TableHeader() []string {
	return []string{"From", "Nonce", "NewSudoAddress", "TxHash"}
}

func (c *ChangeSudoAddressResponse) TableRows() [][]string {
	return [][]string{
		{c.From, strconv.Itoa(int(c.Nonce)), c.NewSudoAddress, c.TxHash},
	}
}

type UpdateValidatorOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey string
	// PubKey is the public key of the validtor being updated
	PubKey string
	// Power is the new power of the validator
	Power string
	// SequencerChainID is the chain ID of the sequencer
	SequencerChainID string
}

type UpdateValidatorResponse struct {
	// From is the address of the sender
	From string `json:"from"`
	// Nonce is the nonce of the transaction
	Nonce uint32 `json:"nonce"`
	// Is the public key of the validator being updated
	PubKey string `json:"pubKey"`
	// Power is the new power of the validator
	Power string `json:"power"`
	// TxHash is the hash of the transaction
	TxHash string `json:"txHash"`
}

func (uv *UpdateValidatorResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(uv, "", "  ")
}

func (uv *UpdateValidatorResponse) TableHeader() []string {
	return []string{"From", "Nonce", "PubKey", "Power", "TxHash"}
}

func (uv *UpdateValidatorResponse) TableRows() [][]string {
	return [][]string{
		{uv.From, strconv.Itoa(int(uv.Nonce)), uv.PubKey, uv.Power, uv.TxHash},
	}
}
