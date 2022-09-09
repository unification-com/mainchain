package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding wrkchain type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.RegisteredWrkChainPrefix):
			var wcA, wcB types.WrkChain
			cdc.MustUnmarshal(kvA.Value, &wcA)
			cdc.MustUnmarshal(kvB.Value, &wcB)
			return fmt.Sprintf("%v\n%v", wcA, wcB)
		case bytes.Equal(kvA.Key[:1], types.RecordedWrkChainBlockHashPrefix):
			var wcbA, wcbB types.WrkChainBlock
			cdc.MustUnmarshal(kvA.Value, &wcbA)
			cdc.MustUnmarshal(kvB.Value, &wcbB)
			return fmt.Sprintf("%v\n%v", wcbA, wcbB)
		case bytes.Equal(kvA.Key[:1], types.HighestWrkChainIDKey):
			kA := binary.BigEndian.Uint64(kvA.Value)
			kB := binary.BigEndian.Uint64(kvB.Value)
			return fmt.Sprintf("%v\n%v", kA, kB)
		default:
			panic(fmt.Sprintf("invalid wrkchain key prefix %X", kvA.Key[:1]))
		}
	}
}
