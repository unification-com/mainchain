package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace params.Subspace
	codespace  sdk.CodespaceType
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the wrkchain Keeper
func NewKeeper(storeKey sdk.StoreKey, paramSpace params.Subspace, codespace sdk.CodespaceType, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:   storeKey,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:  codespace,
		cdc:        cdc,
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

//__WRKCHAIN_ID_________________________________________________________

// GetHighestWrkChainID gets the highest WRKChain ID
func (k Keeper) GetHighestWrkChainID(ctx sdk.Context) (wrkChainID uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestWrkChainIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis(k.codespace, "initial wrkchain ID hasn't been set")
	}
	// convert from bytes to uint64
	wrkChainID = types.GetWrkChainIDFromBytes(bz)
	return wrkChainID, nil
}

// SetHighestWrkChainID sets the new highest WRKChain ID to the store
func (k Keeper) SetHighestWrkChainID(ctx sdk.Context, wrkChainID uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert from uint64 to bytes for storage
	wrkChainIDbz := types.GetWrkChainIDBytes(wrkChainID)
	store.Set(types.HighestWrkChainIDKey, wrkChainIDbz)
}
