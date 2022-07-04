package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

//__PARAMS______________________________________________________________

// GetParams returns the total set of Beacon parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	return types.NewParams(
		k.GetParamRegistrationFee(ctx),
		k.GetParamRecordFee(ctx),
		k.GetParamPurchaseStorageFee(ctx),
		k.GetParamDenom(ctx),
		k.GetParamDefaultStorageLimit(ctx),
		k.GetParamMaxStorageLimit(ctx),
	)
}

// SetParams sets the total set of Beacon parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) GetParamDenom(ctx sdk.Context) string {
	var denomParams string
	k.paramSpace.Get(ctx, types.KeyDenom, &denomParams)
	return denomParams
}

func (k Keeper) GetParamRegistrationFee(ctx sdk.Context) uint64 {
	var feeRegParams uint64
	k.paramSpace.Get(ctx, types.KeyFeeRegister, &feeRegParams)
	return feeRegParams
}

func (k Keeper) GetParamRecordFee(ctx sdk.Context) uint64 {
	var feeRecordParams uint64
	k.paramSpace.Get(ctx, types.KeyFeeRecord, &feeRecordParams)
	return feeRecordParams
}

func (k Keeper) GetParamPurchaseStorageFee(ctx sdk.Context) uint64 {
	var feePurchaseStorageParams uint64
	k.paramSpace.Get(ctx, types.KeyFeePurchaseStorage, &feePurchaseStorageParams)
	return feePurchaseStorageParams
}

func (k Keeper) GetParamDefaultStorageLimit(ctx sdk.Context) uint64 {
	var defaultStorageLimitParams uint64
	k.paramSpace.Get(ctx, types.KeyDefaultStorageLimit, &defaultStorageLimitParams)
	return defaultStorageLimitParams
}

func (k Keeper) GetParamMaxStorageLimit(ctx sdk.Context) uint64 {
	var maxStorageLimitParams uint64
	k.paramSpace.Get(ctx, types.KeyMaxStorageLimit, &maxStorageLimitParams)
	return maxStorageLimitParams
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
