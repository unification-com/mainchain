package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// supplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) supplyexported.SupplyI

	GetModuleAddress(name string) sdk.AccAddress

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error

	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error

	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}
