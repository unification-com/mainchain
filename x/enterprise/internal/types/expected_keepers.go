package types

import (
sdk "github.com/cosmos/cosmos-sdk/types"
"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, exported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}
