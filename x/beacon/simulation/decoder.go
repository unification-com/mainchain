package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/unification-com/mainchain/x/beacon/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding beacon type
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.RegisteredBeaconPrefix):
			var bA, bB types.Beacon
			cdc.MustUnmarshal(kvA.Value, &bA)
			cdc.MustUnmarshal(kvB.Value, &bB)
			return fmt.Sprintf("%v\n%v", bA, bB)
		case bytes.Equal(kvA.Key[:1], types.RecordedBeaconTimestampPrefix):
			var btsA, btsB types.BeaconTimestamp
			cdc.MustUnmarshal(kvA.Value, &btsA)
			cdc.MustUnmarshal(kvB.Value, &btsB)
			return fmt.Sprintf("%v\n%v", btsA, btsB)
		case bytes.Equal(kvA.Key[:1], types.HighestBeaconIDKey):
			kA := binary.BigEndian.Uint64(kvA.Value)
			kB := binary.BigEndian.Uint64(kvB.Value)
			return fmt.Sprintf("%v\n%v", kA, kB)
		default:
			panic(fmt.Sprintf("invalid beacon key prefix %X", kvA.Key[:1]))
		}
	}
}
