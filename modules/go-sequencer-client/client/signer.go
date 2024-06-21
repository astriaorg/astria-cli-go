package client

import (
	"crypto/ed25519"
	"crypto/sha256"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
)

const DefaultAstriaAsset = "nria"

var (
	DefaultAstriaAssetID = sha256.Sum256([]byte(DefaultAstriaAsset))
)

type Signer struct {
	private ed25519.PrivateKey
}

func NewSigner(private ed25519.PrivateKey) *Signer {
	return &Signer{
		private: private,
	}
}

func GenerateSigner() (*Signer, error) {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	return &Signer{
		private: priv,
	}, nil
}

func (s *Signer) SignTransaction(tx *txproto.UnsignedTransaction) (*txproto.SignedTransaction, error) {
	for _, action := range tx.Actions {
		switch v := action.Value.(type) {
		case *txproto.Action_TransferAction:
			if len(v.TransferAction.FeeAssetId) == 0 {
				v.TransferAction.FeeAssetId = DefaultAstriaAssetID[:]
			}
		case *txproto.Action_SequenceAction:
			if len(v.SequenceAction.FeeAssetId) == 0 {
				v.SequenceAction.FeeAssetId = DefaultAstriaAssetID[:]
			}
		}
	}

	bytes, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}

	transaction := &anypb.Any{
		TypeUrl: "/astria.protocol.transactions.v1alpha1.UnsignedTransaction",
		Value:   bytes,
	}

	sig := ed25519.Sign(s.private, bytes)
	return &txproto.SignedTransaction{
		Transaction: transaction,
		Signature:   sig,
		PublicKey:   s.private.Public().(ed25519.PublicKey),
	}, nil
}

// Seed returns the 32-byte "seed" for the key, which is used as the
// input to generate a private key in the rust implementation, ie:
// `ed25519_consensus::SigningKey::from(seed)`
func (s *Signer) Seed() [ed25519.SeedSize]byte {
	return [ed25519.SeedSize]byte(s.private.Seed())
}

func (s *Signer) PublicKey() ed25519.PublicKey {
	return s.private.Public().(ed25519.PublicKey)
}

func (s *Signer) Address() [20]byte {
	hash := sha256.Sum256(s.PublicKey())
	var addr [20]byte
	copy(addr[:], hash[:20])
	return addr
}
