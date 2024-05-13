package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   storetypes.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace paramtypes.Subspace
	bankKeeper types.BankKeeper
	accKeeper  types.AccountKeeper
	cdc        codec.BinaryCodec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the enterprise Keeper
func NewKeeper(storeKey storetypes.StoreKey, bankKeeper types.BankKeeper,
	accKeeper types.AccountKeeper, paramSpace paramtypes.Subspace,
	cdc codec.BinaryCodec) Keeper {

	// ensure module account is set in SupplyKeeper
	if addr := accKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:   storeKey,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		bankKeeper: bankKeeper,
		accKeeper:  accKeeper,
		cdc:        cdc,
	}
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
