package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EnterpriseKeeper defines the expected enterprise keeper
type EnterpriseKeeper interface {
	GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin
}

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin

	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin

	LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI // only used for simulation
}
