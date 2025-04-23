package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/unification-com/mainchain/x/enterprise/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   storetypes.StoreKey // Unexposed key to access store from sdk.Context
	bankKeeper types.BankKeeper
	accKeeper  types.AccountKeeper
	cdc        codec.BinaryCodec // The wire codec for binary encoding/decoding.
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates new instances of the enterprise Keeper
func NewKeeper(storeKey storetypes.StoreKey, bankKeeper types.BankKeeper,
	accKeeper types.AccountKeeper, cdc codec.BinaryCodec, authority string) Keeper {

	// ensure module account is set in SupplyKeeper
	if addr := accKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		accKeeper:  accKeeper,
		cdc:        cdc,
		authority:  authority,
	}
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetEnterpriseAccount returns the enterprise ModuleAccount
func (k Keeper) GetEnterpriseAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.bankKeeper.GetAllBalances(ctx, addr)
}
