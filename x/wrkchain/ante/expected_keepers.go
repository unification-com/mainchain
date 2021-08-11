package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the contract needed for AccountKeeper related APIs.
// Interface provides support to use non-sdk AccountKeeper for AnteHandler's decorators.
type AccountKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
}

type EnterpriseKeeper interface {
	GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin
}

type WrkchainKeeper interface {
	GetZeroFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetParamDenom(ctx sdk.Context) string
	GetRegistrationFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetRecordFeeAsCoin(ctx sdk.Context) sdk.Coin
}
