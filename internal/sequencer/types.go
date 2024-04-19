package sequencer

import (
	"encoding/json"
	"math/big"
	"strconv"
)

// Account is the struct that holds the account information.
type Account struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func (a *Account) JSON() ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}

func (a *Account) TableHeader() []string {
	return []string{"Address", "Public Key", "Private Key"}
}

func (a *Account) TableRows() [][]string {
	return [][]string{
		{a.Address, a.PublicKey, a.PrivateKey},
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

// InitbridgeOpts are the options for the InitBridge function.
type InitBridgeOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// fromKey is the private key of the sender
	FromKey string
	// RollupID is the ID of the rollup to creatte the bridge account for
	RollupID string
}
type InitBridgeResponse struct {
	Nonce  uint32 `json:"nonce"`
	TxHash string `json:"txHash"`
}

func (nr *InitBridgeResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(nr, "", "  ")
}

func (nr *InitBridgeResponse) TableHeader() []string {
	return []string{"Nonce", "TxHash"}
}

func (nr *InitBridgeResponse) TableRows() [][]string {
	return [][]string{
		{strconv.Itoa(int(nr.Nonce)), nr.TxHash},
	}
}

type BridgeLockOpts struct {
	// SequencerURL is the URL of the sequencer
	SequencerURL string
	// FromKey is the private key of the sender
	FromKey string
	// ToAddress is the address of the receiver
	ToAddress string
	// Amount is the amount to be locked
	Amount string
	// DestinationChain is the address on the destination chain
	DestinationChain string
}

type BridgeLockResponse struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Nonce  uint32 `json:"nonce"`
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
	// Amount is the amount to be transferred
	Amount string
}

// TransferResponse is the response of the Transfer function.
type TransferResponse struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Nonce  uint32 `json:"nonce"`
	Amount string `json:"amount"` // NOTE - string so we can support huge numbers
	TxHash string `json:"txHash"`
}

func (tr *TransferResponse) JSON() ([]byte, error) {
	return json.MarshalIndent(tr, "", "  ")
}

func (tr *TransferResponse) TableHeader() []string {
	return []string{"From", "To", "Nonce", "Amount", "TxHash"}
}

func (tr *TransferResponse) TableRows() [][]string {
	return [][]string{
		{tr.From, tr.To, strconv.Itoa(int(tr.Nonce)), tr.Amount, tr.TxHash},
	}
}
