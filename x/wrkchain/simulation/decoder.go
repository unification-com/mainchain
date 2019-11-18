package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding wrkchain type
func DecodeStore(cdc *codec.Codec, kvA, kvB cmn.KVPair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.RegisteredWrkChainPrefix):
		var wcA, wcB types.WrkChain
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &wcA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &wcB)
		return fmt.Sprintf("%v\n%v", wcA, wcB)

	case bytes.Equal(kvA.Key[:1], types.RecordedWrkChainBlockHashPrefix):
		var wcbA, wcbB types.WrkChainBlock
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &wcbA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &wcbB)
		return fmt.Sprintf("%v\n%v", wcbA, wcbB)

	default:
		panic(fmt.Sprintf("invalid wrkchain key prefix %X", kvA.Key[:1]))
	}
}

