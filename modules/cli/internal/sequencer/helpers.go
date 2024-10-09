package sequencer

import (
	"crypto/sha256"
	"time"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
)

// rollupIdFromText converts a string to a RollupId protobuf.
func rollupIdFromText(rollup string) *primproto.RollupId {
	hash := sha256.Sum256([]byte(rollup))
	return &primproto.RollupId{
		Inner: hash[:],
	}
}

// nowPlusFiveMinutes returns the current time plus five minutes in nanoseconds.
func nowPlusFiveMinutes() uint64 {
	return uint64(time.Now().UnixNano() + 5*60*1e9)
}
