package client

import (
	"math/big"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
)

func ProtoU128ToBigInt(u128 *primproto.Uint128) *big.Int {
	lo := big.NewInt(0).SetUint64(u128.Lo)
	hi := big.NewInt(0).SetUint64(u128.Hi)
	hi.Lsh(hi, 64)
	return lo.Add(lo, hi)
}
