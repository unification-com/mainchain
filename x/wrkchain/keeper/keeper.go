package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace paramtypes.Subspace
	cdc        codec.BinaryCodec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the wrkchain Keeper
func NewKeeper(storeKey sdk.StoreKey, paramSpace paramtypes.Subspace, cdc codec.BinaryCodec) Keeper {
	return Keeper{
		storeKey:   storeKey,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		cdc:        cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Cdc() codec.BinaryCodec {
	return k.cdc
}
