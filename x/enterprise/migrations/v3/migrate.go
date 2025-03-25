package v3

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/enterprise/exported"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

const (
	ModuleName = "enterprise"
)

var (
	ParamsKey = types.ParamsKey
)

// Migrate migrates the x/enterprise module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/mint
// module state.
func Migrate(
	ctx sdk.Context,
	store storetypes.KVStore,
	legacySubspace exported.Subspace,
	cdc codec.BinaryCodec,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	store.Set(types.ParamsKey, bz)

	return nil
}
