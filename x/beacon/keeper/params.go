package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon/types"
)

//__PARAMS______________________________________________________________

// GetParams returns the total set of Beacon parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)

	return params
}

// SetParams sets the total set of Beacon parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

func (k Keeper) GetParamDenom(ctx sdk.Context) string {
	return k.GetParams(ctx).Denom
}

func (k Keeper) GetParamRegistrationFee(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).FeeRegister
}

func (k Keeper) GetParamRecordFee(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).FeeRecord
}

func (k Keeper) GetParamPurchaseStorageFee(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).FeePurchaseStorage
}

func (k Keeper) GetParamDefaultStorageLimit(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).DefaultStorageLimit
}

func (k Keeper) GetParamMaxStorageLimit(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).MaxStorageLimit
}

func (k Keeper) GetZeroFeeAsCoin(ctx sdk.Context) sdk.Coin {
	return sdk.NewInt64Coin(k.GetParamDenom(ctx), 0)
}

func (k Keeper) GetRegistrationFeeAsCoin(ctx sdk.Context) sdk.Coin {
	return sdk.NewInt64Coin(k.GetParamDenom(ctx), int64(k.GetParamRegistrationFee(ctx)))
}

func (k Keeper) GetRecordFeeAsCoin(ctx sdk.Context) sdk.Coin {
	return sdk.NewInt64Coin(k.GetParamDenom(ctx), int64(k.GetParamRecordFee(ctx)))
}

func (k Keeper) GetPurchaseStorageFeeAsCoin(ctx sdk.Context) sdk.Coin {
	return sdk.NewInt64Coin(k.GetParamDenom(ctx), int64(k.GetParamPurchaseStorageFee(ctx)))
}

func (k Keeper) GetZeroFeeAsCoins(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{k.GetZeroFeeAsCoin(ctx)}
}

func (k Keeper) GetRegistrationFeeAsCoins(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{k.GetRegistrationFeeAsCoin(ctx)}
}

func (k Keeper) GetRecordFeeAsCoins(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{k.GetRecordFeeAsCoin(ctx)}
}

func (k Keeper) GetPurchaseStorageFeeAsCoins(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{k.GetPurchaseStorageFeeAsCoin(ctx)}
}
