package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EnterpriseKeeper defines the expected enterprise keeper
type EnterpriseKeeper interface {
	GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin
}
