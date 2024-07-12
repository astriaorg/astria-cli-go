package bundler

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// CreateBasicTransferAction creates a basic transfer action.
func CreateBasicTransferAction(toAddr, amount, asset, feeAsset string) *txproto.Action_TransferAction {
	a, err := convertToUint128(amount)
	if err != nil {
		log.WithError(err).Error("Error converting amount to Uint128 proto")
		return nil
	}
	return &txproto.Action_TransferAction{
		TransferAction: &txproto.TransferAction{
			To: &primproto.Address{
				Bech32M: toAddr,
			},
			Amount:   a,
			Asset:    asset,
			FeeAsset: feeAsset,
		},
	}
}

// CreateBasicSequenceAction creates a basic sequence action.
func CreateBasicSequenceAction(rollupId, data, feeAsset string) *txproto.Action_SequenceAction {
	hash := sha256.Sum256([]byte(rollupId))
	dataBytes := []byte(data)
	return &txproto.Action_SequenceAction{
		SequenceAction: &txproto.SequenceAction{
			RollupId: &primproto.RollupId{
				Inner: hash[:],
			},
			Data:     dataBytes,
			FeeAsset: feeAsset,
		},
	}
}

// [  ] *Action_InitBridgeAccountAction
// [  ] *Action_BridgeLockAction
// [  ] *Action_BridgeUnlockAction
// [  ] *Action_BridgeSudoChangeAction
// [  ] *Action_IbcAction
// [  ] *Action_Ics20Withdrawal
// [  ] *Action_SudoAddressChangeAction
// [  ] *Action_ValidatorUpdateAction
// [  ] *Action_IbcRelayerChangeAction
// [  ] *Action_FeeAssetChangeAction
// [  ] *Action_FeeChangeAction

// convertToUint128 converts a string to an Uint128 protobuf
func convertToUint128(numStr string) (*primproto.Uint128, error) {
	bigInt := new(big.Int)

	// convert the string to a big.Int
	_, ok := bigInt.SetString(numStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert string to big.Int")
	}

	// check if the number is negative or overflows Uint128
	if bigInt.Sign() < 0 {
		return nil, fmt.Errorf("negative number not allowed")
	} else if bigInt.BitLen() > 128 {
		return nil, fmt.Errorf("value overflows Uint128")
	}

	// split the big.Int into two uint64s
	// convert the big.Int to uint64, which will drop the higher 64 bits
	lo := bigInt.Uint64()
	// shift the big.Int to the right by 64 bits and convert to uint64
	hi := bigInt.Rsh(bigInt, 64).Uint64()
	uint128 := &primproto.Uint128{
		Lo: lo,
		Hi: hi,
	}

	return uint128, nil
}
