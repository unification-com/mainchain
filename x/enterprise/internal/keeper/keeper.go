package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace   params.Subspace
	codespace    sdk.CodespaceType
	supplyKeeper types.SupplyKeeper
	accKeeper    auth.AccountKeeper
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the enterprise Keeper
func NewKeeper(storeKey sdk.StoreKey, supplyKeeper types.SupplyKeeper, accKeeper auth.AccountKeeper, paramSpace params.Subspace, codespace sdk.CodespaceType, cdc *codec.Codec) Keeper {

	// ensure module account is set in SupplyKeeper
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     storeKey,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:    codespace,
		supplyKeeper: supplyKeeper,
		accKeeper:    accKeeper,
		cdc:          cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) Cdc() *codec.Codec {
	return k.cdc
}

// GetEnterpriseAccount returns the enterprise ModuleAccount
func (k Keeper) GetEnterpriseAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
