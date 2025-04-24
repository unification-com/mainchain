package v4

import (
	storetypes "cosmossdk.io/store/types"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	ModuleName = "wrkchain"
)

var (
	RegisteredWrkChainPrefix = []byte{0x01}
)

// WrkChainKey gets a specific purchase order ID key for use in the store
func WrkChainKey(wrkChainID uint64) []byte {
	return append(RegisteredWrkChainPrefix, GetWrkChainIDBytes(wrkChainID)...)
}

// GetWrkChainIDBytes returns the byte representation of the wrkChainID
// used for getting the highest WRKChain ID from the database
func GetWrkChainIDBytes(wrkChainID uint64) (wrkChainIDBz []byte) {
	wrkChainIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(wrkChainIDBz, wrkChainID)
	return
}

// migrateRegisteredWrkChainsToNewType migrates all WrkChains to use the new BaseType field name for Type
func migrateRegisteredWrkChainsToNewType(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec) error {
	iterator := storetypes.KVStorePrefixIterator(store, RegisteredWrkChainPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oldWc V3WrkChain

		err := cdc.Unmarshal(iterator.Value(), &oldWc)
		if err != nil {
			return err
		}

		newWc := types.WrkChain{
			WrkchainId:   oldWc.WrkchainId,
			Moniker:      oldWc.Moniker,
			Name:         oldWc.Name,
			Genesis:      oldWc.Genesis,
			BaseType:     oldWc.Type,
			Lastblock:    oldWc.Lastblock,
			NumBlocks:    oldWc.NumBlocks,
			LowestHeight: oldWc.LowestHeight,
			RegTime:      oldWc.RegTime,
			Owner:        oldWc.Owner,
		}

		bz, err := cdc.Marshal(&newWc)
		if err != nil {
			return err
		}

		store.Set(iterator.Key(), bz)
	}

	return nil

}

// Migrate performs in-place store migrations from v3 to v4.
func Migrate(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec) error {
	if err := migrateRegisteredWrkChainsToNewType(ctx, store, cdc); err != nil {
		return err
	}
	return nil
}
