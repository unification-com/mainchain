package v3_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/x/enterprise"
	"github.com/unification-com/mainchain/x/enterprise/exported"
	v3 "github.com/unification-com/mainchain/x/enterprise/migrations/v3"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

type mockSubspace struct {
	ps types.Params
}

func newMockSubspace(ps types.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*types.Params) = ms.ps
}

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(enterprise.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(v3.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	store := ctx.KVStore(storeKey)

	defaultParams := types.DefaultParams()

	pk := ed25519.GenPrivKey().PubKey()
	defaultParams.EntSigners = sdk.AccAddress(pk.Address()).String()
	legacySubspace := newMockSubspace(defaultParams)
	require.NoError(t, v3.Migrate(ctx, store, legacySubspace, cdc))

	var res types.Params
	bz := store.Get(v3.ParamsKey)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, legacySubspace.ps, res)
}
