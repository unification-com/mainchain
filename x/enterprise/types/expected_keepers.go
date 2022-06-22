package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetPaginatedTotalSupply(ctx sdk.Context, pagination *query.PageRequest) (sdk.Coins, *query.PageResponse, error)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error

	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error

	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}
