package keeper

import (
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/stream/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeKey         storetypes.StoreKey
		bankKeeper       types.BankKeeper
		accKeeper        types.AccountKeeper
		feeCollectorName string
		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	accKeeper types.AccountKeeper,
	cdc codec.BinaryCodec,
	feeCollectorName string,
	authority string,
) Keeper {

	// ensure module account is set in SupplyKeeper
	if addr := accKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		authority:        authority,
		bankKeeper:       bankKeeper,
		accKeeper:        accKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// GetAuthority returns the x/stream module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Cdc() codec.BinaryCodec {
	return k.cdc
}

// GetStreamModuleAccount returns the stream ModuleAccount
func (k Keeper) GetStreamModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accKeeper.GetModuleAccount(ctx, types.ModuleName)
}
