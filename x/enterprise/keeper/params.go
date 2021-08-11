package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"strings"
)

//__PARAMS______________________________________________________________

func (k Keeper) GetParamDenom(ctx sdk.Context) string {
	var paramDenom string
	k.paramSpace.Get(ctx, types.KeyDenom, &paramDenom)
	return paramDenom
}

func (k Keeper) GetParamMinAccepts(ctx sdk.Context) uint64 {
	var paramMinAccepts uint64
	k.paramSpace.Get(ctx, types.KeyMinAccepts, &paramMinAccepts)
	return paramMinAccepts
}

func (k Keeper) GetParamDecisionLimit(ctx sdk.Context) uint64 {
	var paramDecisionLimit uint64
	k.paramSpace.Get(ctx, types.KeyDecisionLimit, &paramDecisionLimit)
	return paramDecisionLimit
}

func (k Keeper) GetParamEntSigners(ctx sdk.Context) string {
	var paramEntSigners string
	k.paramSpace.Get(ctx, types.KeyEntSigners, &paramEntSigners)
	return paramEntSigners
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
	return types.NewParams(
		k.GetParamDenom(ctx),
		k.GetParamMinAccepts(ctx),
		k.GetParamDecisionLimit(ctx),
		k.GetParamEntSigners(ctx),
	)
}

// SetParams sets the total set of Enterprise FUND parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
