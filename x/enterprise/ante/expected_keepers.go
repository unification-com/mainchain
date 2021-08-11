package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EnterpriseKeeper interface {
	GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin
	IsLocked(ctx sdk.Context, address sdk.AccAddress) bool
	UnlockCoinsForFees(ctx sdk.Context, feePayer sdk.AccAddress, feesToPay sdk.Coins) error
}
