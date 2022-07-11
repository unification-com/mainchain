package v045_test

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v045 "github.com/unification-com/mainchain/x/enterprise/legacy/v045"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	appparams "github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

const (
	modName       = "enterprise"
	paramsModName = "params"
)

var (
	paramsKey  = sdk.NewKVStoreKey(paramsModName)
	paramsTkey = sdk.NewTransientStoreKey("transient_params")
	entKey     = sdk.NewKVStoreKey(modName)

	v045TestParams = types.NewParams(test_helpers.TestDenomination, 1, 1, "")
)

func TestStoreMigrate(t *testing.T) {
	testAddrs := test_helpers.GenerateRandomTestAccounts(100)
	expectedTotalUsed := uint64(0)
	totalLockedCoin := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(0))

	encCfg := appparams.MakeTestEncodingConfig()
	cdc := encCfg.Marshaler

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTkey, sdk.StoreTypeTransient, db)
	cms.MountStoreWithDB(entKey, sdk.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	require.NoError(t, err)

	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())

	entStore := ctx.KVStore(entKey)
	paramsStore := ctx.KVStore(paramsKey)

	_ = prefix.NewStore(paramsStore, []byte(modName))
	v045ss := paramtypes.NewSubspace(cdc, encCfg.Amino, paramsKey, paramsTkey, modName)
	v045EntSs := v045ss.WithKeyTable(types.ParamKeyTable())

	// set params
	v045TestParams.EntSigners = testAddrs[0].String()
	v045EntSs.SetParamSet(ctx, &v045TestParams)

	// populate store
	for i := 0; i < len(testAddrs); i++ {
		poAmount := uint64(i+1) * 10
		toUnlock := uint64(i + 1)
		expectedTotalUsed += toUnlock
		poAmountCoin := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(poAmount))
		toUnlockCoin := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(toUnlock))
		lockedCoin := poAmountCoin.Sub(toUnlockCoin)
		totalLockedCoin = totalLockedCoin.Add(lockedCoin)
		lockedUnd := types.LockedUnd{
			Owner:  testAddrs[i].String(),
			Amount: lockedCoin,
		}

		completedId := uint64(i + 1)
		raisedId := uint64(i + len(testAddrs))
		acceptedId := uint64(i + (len(testAddrs) * 2))
		rejectedId := uint64(i + (len(testAddrs) * 3))

		newCompletedPo := types.EnterpriseUndPurchaseOrder{
			Id:        completedId,
			Purchaser: testAddrs[i].String(),
			Amount:    poAmountCoin,
			Status:    types.StatusCompleted,
		}

		// set a completed order
		entStore.Set(types.PurchaseOrderKey(newCompletedPo.Id), cdc.MustMarshal(&newCompletedPo))

		newRaisedPo := types.EnterpriseUndPurchaseOrder{
			Id:        raisedId,
			Purchaser: testAddrs[i].String(),
			Amount:    poAmountCoin,
			Status:    types.StatusRaised,
		}

		// set a raised order
		entStore.Set(types.PurchaseOrderKey(newRaisedPo.Id), cdc.MustMarshal(&newRaisedPo))

		newAcceptedPo := types.EnterpriseUndPurchaseOrder{
			Id:        acceptedId,
			Purchaser: testAddrs[i].String(),
			Amount:    poAmountCoin,
			Status:    types.StatusAccepted,
		}

		// set an accepted order
		entStore.Set(types.PurchaseOrderKey(newAcceptedPo.Id), cdc.MustMarshal(&newAcceptedPo))

		newRejectedPo := types.EnterpriseUndPurchaseOrder{
			Id:        rejectedId,
			Purchaser: testAddrs[i].String(),
			Amount:    poAmountCoin,
			Status:    types.StatusRejected,
		}

		// set a rejected order
		entStore.Set(types.PurchaseOrderKey(newRejectedPo.Id), cdc.MustMarshal(&newRejectedPo))

		// set amount locked, based on completed order
		entStore.Set(types.LockedUndAddressStoreKey(testAddrs[i]), cdc.MustMarshal(&lockedUnd))
	}

	// set total locked
	entStore.Set(types.TotalLockedUndKey, cdc.MustMarshal(&totalLockedCoin))

	// migrate
	err = v045.MigrateStore(ctx, entKey, v045EntSs, cdc)
	require.NoError(t, err)

	// check
	for i := 0; i < len(testAddrs); i++ {
		expectedUsed := uint64(i + 1)
		expectedUsedCoin := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(expectedUsed))

		bz := entStore.Get(types.SpentEFUNDAddressStoreKey(testAddrs[i]))
		var usedUnd sdk.Coin
		err := cdc.Unmarshal(bz, &usedUnd)
		require.NoError(t, err)
		require.Equal(t, expectedUsedCoin, usedUnd)
	}

	expectedTotalUsedCoin := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(expectedTotalUsed))
	bz := entStore.Get(types.TotalSpentEFUNDKey)
	var totalUsedCoin sdk.Coin
	err = cdc.Unmarshal(bz, &totalUsedCoin)
	require.NoError(t, err)
	require.Equal(t, expectedTotalUsedCoin, totalUsedCoin)
}
