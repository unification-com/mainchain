package v045_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	appparams "github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/app/test_helpers"
	v040 "github.com/unification-com/mainchain/x/wrkchain/legacy/v040"
	v045 "github.com/unification-com/mainchain/x/wrkchain/legacy/v045"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	modName       = "wrkchain"
	paramsModName = "params"
)

var (
	paramsKey   = sdk.NewKVStoreKey(paramsModName)
	paramsTkey  = sdk.NewTransientStoreKey("transient_params")
	wrkchainKey = sdk.NewKVStoreKey(modName)

	v045TestParams = types.NewParams(24, 1, 1, test_helpers.TestDenomination, 1, 1)
)

func TestStoreMigrate(t *testing.T) {
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)
	lastBlockHeight := uint64(69000)
	expectedLowestHeight := (lastBlockHeight - types.DefaultStorageLimit) + 1

	encCfg := appparams.MakeTestEncodingConfig()
	cdc := encCfg.Marshaler

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTkey, sdk.StoreTypeTransient, db)
	cms.MountStoreWithDB(wrkchainKey, sdk.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	require.NoError(t, err)

	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())

	wrkchainStore := ctx.KVStore(wrkchainKey)
	paramsStore := ctx.KVStore(paramsKey)

	_ = prefix.NewStore(paramsStore, []byte(modName))
	v045ss := paramtypes.NewSubspace(cdc, encCfg.Amino, paramsKey, paramsTkey, modName)
	v045WrkchainSs := v045ss.WithKeyTable(types.ParamKeyTable())

	// set old params
	v045WrkchainSs.SetParamSet(ctx, &v045TestParams)

	wcRegTime := uint64(ctx.BlockTime().Unix())
	oldWrkchain := v040.WrkChain{
		WrkchainId:   1,
		Moniker:      "wc1",
		Name:         "wc1",
		Genesis:      "wc1_genesis",
		Type:         "cosmos",
		Lastblock:    0,
		NumBlocks:    0,
		LowestHeight: 0,
		RegTime:      wcRegTime,
		Owner:        testAddrs[0].String(),
	}

	wrkchainStore.Set(v040.WrkChainKey(1), cdc.MustMarshal(&oldWrkchain))

	for i := uint64(1); i <= lastBlockHeight; i++ {
		subTime := uint64(time.Now().Unix())
		hash := test_helpers.GenerateRandomString(32)

		wcBlock := v040.WrkChainBlock{
			Height:     i,
			Blockhash:  hash,
			Parenthash: "",
			Hash1:      "",
			Hash2:      "",
			Hash3:      "",
			SubTime:    subTime,
		}

		wrkchainStore.Set(v040.WrkChainBlockKey(1, i), cdc.MustMarshal(&wcBlock))
	}

	oldWrkchain.Lastblock = lastBlockHeight
	oldWrkchain.LowestHeight = 1
	oldWrkchain.NumBlocks = lastBlockHeight

	wrkchainStore.Set(v040.WrkChainKey(1), cdc.MustMarshal(&oldWrkchain))

	err = v045.MigrateStore(ctx, wrkchainKey, v045WrkchainSs, cdc)
	require.NoError(t, err)

	var migratedParams types.Params

	v045WrkchainSs.GetParamSet(ctx, &migratedParams)

	// should remain the same
	require.True(t, migratedParams.FeeRegister == 24)
	require.True(t, migratedParams.FeeRecord == 1)
	require.True(t, migratedParams.Denom == test_helpers.TestDenomination)

	// should be added
	require.True(t, migratedParams.DefaultStorageLimit == types.DefaultStorageLimit)
	require.True(t, migratedParams.MaxStorageLimit == types.DefaultMaxStorageLimit)
	require.True(t, migratedParams.FeePurchaseStorage == types.PurchaseStorageFee)

	var newWrkchain types.WrkChain

	newWrkchainBz := wrkchainStore.Get(types.WrkChainKey(1))

	err = cdc.Unmarshal(newWrkchainBz, &newWrkchain)

	require.NoError(t, err)
	// should be migrated
	require.True(t, newWrkchain.Lastblock == lastBlockHeight)
	require.True(t, newWrkchain.NumBlocks == types.DefaultStorageLimit)
	require.True(t, newWrkchain.LowestHeight == expectedLowestHeight)

	// should remain the same
	require.True(t, newWrkchain.WrkchainId == 1)
	require.True(t, newWrkchain.Name == "wc1")
	require.True(t, newWrkchain.Moniker == "wc1")
	require.True(t, newWrkchain.Genesis == "wc1_genesis")
	require.True(t, newWrkchain.Type == "cosmos")
	require.True(t, newWrkchain.RegTime == wcRegTime)
	require.True(t, newWrkchain.Owner == testAddrs[0].String())

	// check expected prunes
	for i := uint64(1); i <= lastBlockHeight; i++ {
		has := wrkchainStore.Has(types.WrkChainBlockKey(1, i))
		if i < expectedLowestHeight {
			require.False(t, has)
		} else {
			require.True(t, has)
		}
	}
}
