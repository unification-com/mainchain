package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper Methods imported from account should be defined here
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
	GetModuleAddress(module string) sdk.AccAddress
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}

// BankKeeper Methods imported from bank should be defined here
type BankKeeper interface {
	BlockedAddr(recipient sdk.AccAddress) bool
	GetBlockedAddresses() map[string]bool
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoin(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, fromModule string, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, fromAddr sdk.AccAddress, toModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
