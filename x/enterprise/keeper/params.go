package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/enterprise/types"
)

//__PARAMS______________________________________________________________

func (k Keeper) GetParamDenom(ctx sdk.Context) string {
	return k.GetParams(ctx).Denom
}

func (k Keeper) GetParamMinAccepts(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).MinAccepts
}

func (k Keeper) GetParamDecisionLimit(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).DecisionTimeLimit
}

func (k Keeper) GetParamEntSigners(ctx sdk.Context) string {
	return k.GetParams(ctx).EntSigners
}

func (k Keeper) GetParamEntSignersAsAddressArray(ctx sdk.Context) []sdk.AccAddress {
	var entSignersArray []sdk.AccAddress
	paramEntSigners := k.GetParamEntSigners(ctx)
	entSigners := strings.Split(paramEntSigners, ",")
	for _, authAddr := range entSigners {
		addr, err := sdk.AccAddressFromBech32(authAddr)
		if err == nil && !addr.Empty() {
			entSignersArray = append(entSignersArray, addr)
		}
	}
	return entSignersArray
}

// GetParams returns the total set of Enterprise FUND parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)

	return params
}

// SetParams sets the total set of Enterprise FUND parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}
