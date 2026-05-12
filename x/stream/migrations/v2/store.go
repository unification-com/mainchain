// Package v2 migrates the x/stream KV store from the v1 key shape
// (0x11 | len(receiver) | receiver | len(sender) | sender) to the v2 key shape
// (0x11 | len(receiver) | receiver | len(sender) | sender | len(denom) | denom)
// by reading each row's serialised Stream value, extracting Deposit.Denom, and
// rewriting the key. Values are byte-identical.
package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/unification-com/mainchain/x/stream/types"
)

// MigrateStore performs the v1 → v2 in-place migration for x/stream.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.StreamKeyPrefix)

	type rewrite struct {
		oldKey []byte
		newKey []byte
		value  []byte
	}
	var batch []rewrite

	for ; iter.Valid(); iter.Next() {
		var s types.Stream
		cdc.MustUnmarshal(iter.Value(), &s)

		receiver, sender := AddressesFromStreamKeyV1(iter.Key())
		newKey := newStreamKeyV2(receiver, sender, s.Deposit.Denom)

		batch = append(batch, rewrite{
			oldKey: append([]byte(nil), iter.Key()...),
			newKey: newKey,
			value:  append([]byte(nil), iter.Value()...),
		})
	}
	// Close before mutating the store.
	if err := iter.Close(); err != nil {
		return err
	}

	for _, r := range batch {
		store.Delete(r.oldKey)
		store.Set(r.newKey, r.value)
	}
	return nil
}

// AddressesFromStreamKeyV1 is a frozen copy of the pre-v2 key parser.
// Never call from live code after the upgrade. Format:
//
//	0x11 | len(receiver) | receiver | len(sender) | sender
func AddressesFromStreamKeyV1(key []byte) (sdk.AccAddress, sdk.AccAddress) {
	receiverAddrLen, receiverAddrLenEndIndex := sdk.ParseLengthPrefixedBytes(key, 1, 1) // skip 0x11 prefix
	receiverAddr, receiverAddrEndIndex := sdk.ParseLengthPrefixedBytes(key, receiverAddrLenEndIndex+1, int(receiverAddrLen[0]))

	senderAddrLen, senderAddrLenEndIndex := sdk.ParseLengthPrefixedBytes(key, receiverAddrEndIndex+1, 1)
	senderAddr, senderAddrEndIndex := sdk.ParseLengthPrefixedBytes(key, senderAddrLenEndIndex+1, int(senderAddrLen[0]))

	kv.AssertKeyAtLeastLength(key, senderAddrEndIndex+1)
	return receiverAddr, senderAddr
}

// newStreamKeyV2 mirrors types.GetStreamKey but is duplicated locally so the migration
// is insulated from any future drift in the live key constructor.
func newStreamKeyV2(receiver sdk.AccAddress, sender sdk.AccAddress, denom string) []byte {
	out := append([]byte{}, types.StreamKeyPrefix...)
	out = append(out, address.MustLengthPrefix(receiver)...)
	out = append(out, address.MustLengthPrefix(sender)...)
	denomBytes := []byte(denom)
	if len(denomBytes) > 255 {
		panic("denom length must be <= 255")
	}
	out = append(out, byte(len(denomBytes)))
	out = append(out, denomBytes...)
	return out
}
