package simulation

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding wrkchain type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.RegisteredWrkChainPrefix):
			var wcA, wcB types.WrkChain
			cdc.MustUnmarshalBinaryBare(kvA.Value, &wcA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &wcB)
			return fmt.Sprintf("%v\n%v", wcA, wcB)

		case bytes.Equal(kvA.Key[:1], types.RecordedWrkChainBlockHashPrefix):
			var wcbA, wcbB types.WrkChainBlock
			cdc.MustUnmarshalBinaryBare(kvA.Value, &wcbA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &wcbB)
			return fmt.Sprintf("%v\n%v", wcbA, wcbB)

		default:
			panic(fmt.Sprintf("invalid wrkchain key prefix %X", kvA.Key[:1]))
		}
	}
}
