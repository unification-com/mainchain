package v4_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/wrkchain/types"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/unification-com/mainchain/x/wrkchain"
	v4 "github.com/unification-com/mainchain/x/wrkchain/migrations/v4"
)

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(wrkchain.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(v4.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	store := ctx.KVStore(storeKey)

	oldWc1 := v4.V3WrkChain{
		WrkchainId:   1,
		Moniker:      "Moniker1",
		Name:         "name1",
		Genesis:      "genesis1",
		Type:         "cosmos1",
		Lastblock:    201,
		NumBlocks:    100,
		LowestHeight: 101,
		RegTime:      777,
		Owner:        "addr1",
	}

	oldWc2 := v4.V3WrkChain{
		WrkchainId:   2,
		Moniker:      "Moniker2",
		Name:         "name2",
		Genesis:      "genesis2",
		Type:         "cosmos2",
		Lastblock:    301,
		NumBlocks:    200,
		LowestHeight: 201,
		RegTime:      888,
		Owner:        "addr2",
	}

	oldWc3 := v4.V3WrkChain{
		WrkchainId:   3,
		Moniker:      "Moniker3",
		Name:         "name3",
		Genesis:      "genesis3",
		Type:         "cosmos3",
		Lastblock:    401,
		NumBlocks:    300,
		LowestHeight: 301,
		RegTime:      999,
		Owner:        "addr3",
	}

	store.Set(v4.WrkChainKey(1), cdc.MustMarshal(&oldWc1))
	store.Set(v4.WrkChainKey(2), cdc.MustMarshal(&oldWc2))
	store.Set(v4.WrkChainKey(3), cdc.MustMarshal(&oldWc3))

	require.NoError(t, v4.Migrate(ctx, store, cdc))

	newWc1Bz := store.Get(v4.WrkChainKey(1))
	require.NotNil(t, newWc1Bz)
	var newWc1 types.WrkChain
	cdc.MustUnmarshal(newWc1Bz, &newWc1)

	require.Equal(t, oldWc1.WrkchainId, newWc1.WrkchainId)
	require.Equal(t, oldWc1.Moniker, newWc1.Moniker)
	require.Equal(t, oldWc1.Name, newWc1.Name)
	require.Equal(t, oldWc1.Genesis, newWc1.Genesis)
	require.Equal(t, oldWc1.Type, newWc1.BaseType)
	require.Equal(t, oldWc1.Lastblock, newWc1.Lastblock)
	require.Equal(t, oldWc1.NumBlocks, newWc1.NumBlocks)
	require.Equal(t, oldWc1.LowestHeight, newWc1.LowestHeight)
	require.Equal(t, oldWc1.RegTime, newWc1.RegTime)
	require.Equal(t, oldWc1.Owner, newWc1.Owner)

	newWc2Bz := store.Get(v4.WrkChainKey(2))
	require.NotNil(t, newWc2Bz)
	var newWc2 types.WrkChain
	cdc.MustUnmarshal(newWc2Bz, &newWc2)

	require.Equal(t, oldWc2.WrkchainId, newWc2.WrkchainId)
	require.Equal(t, oldWc2.Moniker, newWc2.Moniker)
	require.Equal(t, oldWc2.Name, newWc2.Name)
	require.Equal(t, oldWc2.Genesis, newWc2.Genesis)
	require.Equal(t, oldWc2.Type, newWc2.BaseType)
	require.Equal(t, oldWc2.Lastblock, newWc2.Lastblock)
	require.Equal(t, oldWc2.NumBlocks, newWc2.NumBlocks)
	require.Equal(t, oldWc2.LowestHeight, newWc2.LowestHeight)
	require.Equal(t, oldWc2.RegTime, newWc2.RegTime)
	require.Equal(t, oldWc2.Owner, newWc2.Owner)

	newWc3Bz := store.Get(v4.WrkChainKey(3))
	require.NotNil(t, newWc3Bz)
	var newWc3 types.WrkChain
	cdc.MustUnmarshal(newWc3Bz, &newWc3)

	require.Equal(t, oldWc3.WrkchainId, newWc3.WrkchainId)
	require.Equal(t, oldWc3.Moniker, newWc3.Moniker)
	require.Equal(t, oldWc3.Name, newWc3.Name)
	require.Equal(t, oldWc3.Genesis, newWc3.Genesis)
	require.Equal(t, oldWc3.Type, newWc3.BaseType)
	require.Equal(t, oldWc3.Lastblock, newWc3.Lastblock)
	require.Equal(t, oldWc3.NumBlocks, newWc3.NumBlocks)
	require.Equal(t, oldWc3.LowestHeight, newWc3.LowestHeight)
	require.Equal(t, oldWc3.RegTime, newWc3.RegTime)
	require.Equal(t, oldWc3.Owner, newWc3.Owner)
}
