package simulation

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding enterprise type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PurchaseOrderIDKeyPrefix):
			var poA, poB types.EnterpriseUndPurchaseOrder
			cdc.MustUnmarshalBinaryBare(kvA.Value, &poA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &poB)
			return fmt.Sprintf("%v\n%v", poA, poB)

		case bytes.Equal(kvA.Key[:1], types.LockedUndAddressKeyPrefix):
			var lundA, lundB types.LockedUnd
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lundA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &lundB)
			return fmt.Sprintf("%v\n%v", lundA, lundB)

		case bytes.Equal(kvA.Key[:1], types.TotalLockedUndKey):
			var tlA, tlB sdk.Coin
			cdc.MustUnmarshalBinaryBare(kvA.Value, &tlA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &tlB)
			return fmt.Sprintf("%v\n%v", tlA, tlB)

		default:
			panic(fmt.Sprintf("invalid enterprise key prefix %X", kvA.Key[:1]))
		}
	}
}
