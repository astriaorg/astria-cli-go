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
		Address:    AddressFromPublicKey(pub),
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
