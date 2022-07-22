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
	v040 "github.com/unification-com/mainchain/x/beacon/legacy/v040"
	v045 "github.com/unification-com/mainchain/x/beacon/legacy/v045"
	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	modName       = "beacon"
	paramsModName = "params"
)

var (
	paramsKey  = sdk.NewKVStoreKey(paramsModName)
	paramsTkey = sdk.NewTransientStoreKey("transient_params")
	beaconKey  = sdk.NewKVStoreKey(modName)

	v045TestParams = types.NewParams(24, 1, 1, test_helpers.TestDenomination, 1, 1)
)

func TestStoreMigrate(t *testing.T) {
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)
	lastTimestampId := uint64(69000)
	expectedFirstInState := (lastTimestampId - types.DefaultStorageLimit) + 1

	encCfg := appparams.MakeTestEncodingConfig()
	cdc := encCfg.Marshaler

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTkey, sdk.StoreTypeTransient, db)
	cms.MountStoreWithDB(beaconKey, sdk.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	require.NoError(t, err)

	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())

	beaconStore := ctx.KVStore(beaconKey)
	paramsStore := ctx.KVStore(paramsKey)

	_ = prefix.NewStore(paramsStore, []byte(modName))
	v045ss := paramtypes.NewSubspace(cdc, encCfg.Amino, paramsKey, paramsTkey, modName)
	v045BeaconSs := v045ss.WithKeyTable(types.ParamKeyTable())

	// set old params
	v045BeaconSs.SetParamSet(ctx, &v045TestParams)

	bRegTime := uint64(ctx.BlockTime().Unix())

	oldBeacon := v040.Beacon{
		BeaconId:        1,
		Moniker:         "tb1",
		Name:            "tb1",
		LastTimestampId: 0,
		FirstIdInState:  0,
		NumInState:      0,
		RegTime:         bRegTime,
		Owner:           testAddrs[0].String(),
	}

	beaconStore.Set(v040.BeaconKey(1), cdc.MustMarshal(&oldBeacon))

	for i := uint64(1); i <= lastTimestampId; i++ {
		subTime := uint64(time.Now().Unix())
		hash := test_helpers.GenerateRandomString(32)

		timestamp := v040.BeaconTimestamp{}
		timestamp.TimestampId = i
		timestamp.Hash = hash
		timestamp.SubmitTime = subTime

		beaconStore.Set(v040.BeaconTimestampKey(1, i), cdc.MustMarshal(&timestamp))
	}

	oldBeacon.LastTimestampId = lastTimestampId
	oldBeacon.FirstIdInState = 1
	oldBeacon.NumInState = lastTimestampId

	beaconStore.Set(v040.BeaconKey(1), cdc.MustMarshal(&oldBeacon))

	err = v045.MigrateStore(ctx, beaconKey, v045BeaconSs, cdc)
	require.NoError(t, err)

	var migratedParams types.Params

	v045BeaconSs.GetParamSet(ctx, &migratedParams)

	// should remain the same
	require.True(t, migratedParams.FeeRegister == 24)
	require.True(t, migratedParams.FeeRecord == 1)
	require.True(t, migratedParams.Denom == test_helpers.TestDenomination)

	// should be added
	require.True(t, migratedParams.DefaultStorageLimit == types.DefaultStorageLimit)
	require.True(t, migratedParams.MaxStorageLimit == types.DefaultMaxStorageLimit)
	require.True(t, migratedParams.FeePurchaseStorage == types.PurchaseStorageFee)

	var newBeacon types.Beacon

	newBeaconBz := beaconStore.Get(types.BeaconKey(1))

	err = cdc.Unmarshal(newBeaconBz, &newBeacon)

	require.NoError(t, err)
	// should be migrated
	require.True(t, newBeacon.LastTimestampId == lastTimestampId)
	require.True(t, newBeacon.NumInState == types.DefaultStorageLimit)
	require.True(t, newBeacon.FirstIdInState == expectedFirstInState)

	// should remain the same
	require.True(t, newBeacon.BeaconId == 1)
	require.True(t, newBeacon.Name == "tb1")
	require.True(t, newBeacon.Moniker == "tb1")
	require.True(t, newBeacon.RegTime == bRegTime)
	require.True(t, newBeacon.Owner == testAddrs[0].String())

	// should have new storage data
	var newBeaconStorage types.BeaconStorageLimit

	newBeaconStorageBz := beaconStore.Get(types.BeaconStorageLimitKey(1))

	err = cdc.Unmarshal(newBeaconStorageBz, &newBeaconStorage)
	require.NoError(t, err)
	require.True(t, newBeaconStorage.InStateLimit == types.DefaultStorageLimit)

	// check expected prunes
	for i := uint64(1); i <= lastTimestampId; i++ {
		has := beaconStore.Has(types.BeaconTimestampKey(1, i))
		if i < expectedFirstInState {
			require.False(t, has)
		} else {
			require.True(t, has)
		}
	}
}
