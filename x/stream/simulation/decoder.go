package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/unification-com/mainchain/x/stream/types"
)

// NewDecodeStore returns a decoder function closure that umarshals the KVPair's
// Value to the corresponding stream type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.StreamKeyPrefix):
			var streamA, streamB types.Stream
			cdc.MustUnmarshal(kvA.Value, &streamA)
			cdc.MustUnmarshal(kvB.Value, &streamB)
			return fmt.Sprintf("%v\n%v", streamA, streamB)
		default:
			panic(fmt.Sprintf("invalid stream key %X", kvA.Key))
		}
	}
}
