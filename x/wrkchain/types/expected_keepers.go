package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	GetSupply(ctx sdk.Context, denom string) sdk.Coin

	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation
}
