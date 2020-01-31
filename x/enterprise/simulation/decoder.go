package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmkv "github.com/tendermint/tendermint/libs/kv"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding enterprise type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.PurchaseOrderIDKeyPrefix):
		var poA, poB types.EnterpriseUndPurchaseOrder
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &poA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &poB)
		return fmt.Sprintf("%v\n%v", poA, poB)

	case bytes.Equal(kvA.Key[:1], types.LockedUndAddressKeyPrefix):
		var lundA, lundB types.LockedUnd
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &lundA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &lundB)
		return fmt.Sprintf("%v\n%v", lundA, lundB)

	case bytes.Equal(kvA.Key[:1], types.TotalLockedUndKey):
		var tlA, tlB sdk.Coin
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &tlA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &tlB)
		return fmt.Sprintf("%v\n%v", tlA, tlB)

	default:
		panic(fmt.Sprintf("invalid enterprise key prefix %X", kvA.Key[:1]))
	}
}
