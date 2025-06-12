package ante

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the contract needed for AccountKeeper related APIs.
// Interface provides support to use non-sdk AccountKeeper for AnteHandler's decorators.
type AccountKeeper interface {
	GetParams(ctx context.Context) (params types.Params)
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	GetModuleAddress(moduleName string) sdk.AccAddress
	AddressCodec() address.Codec
}

type BankKeeper interface {
	GetAllBalances(ctx context.Context, address sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx context.Context, address sdk.AccAddress) sdk.Coins
}

type EnterpriseKeeper interface {
	GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin
}

type BeaconKeeper interface {
	GetZeroFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetParamDenom(ctx sdk.Context) string
	GetRegistrationFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetRecordFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetPurchaseStorageFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetMaxPurchasableSlots(ctx sdk.Context, beaconId uint64) uint64
}
