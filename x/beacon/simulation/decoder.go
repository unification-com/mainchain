package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	tmkv "github.com/tendermint/tendermint/libs/kv"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding wrkchain type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.RegisteredBeaconPrefix):
		var bA, bB types.Beacon
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &bA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &bB)
		return fmt.Sprintf("%v\n%v", bA, bB)

	case bytes.Equal(kvA.Key[:1], types.RecordedBeaconTimestampPrefix):
		var btsA, btsB types.BeaconTimestamp
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &btsA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &btsB)
		return fmt.Sprintf("%v\n%v", btsA, btsB)

	default:
		panic(fmt.Sprintf("invalid beacon key prefix %X", kvA.Key[:1]))
	}
}
