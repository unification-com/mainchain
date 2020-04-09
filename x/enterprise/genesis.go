package enterprise

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// InitGenesis new enterprise UND genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, supplyKeeper types.SupplyKeeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestPurchaseOrderID(ctx, data.StartingPurchaseOrderID)

	moduleAcc := keeper.GetEnterpriseAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", ModuleName))
	}

	if data.Whitelist != nil {
		for _, wlAddr := range data.Whitelist {
			err := keeper.AddAddressToWhitelist(ctx, wlAddr)
			if err != nil {
				panic(err)
			}
		}
	}

	err := keeper.SetTotalLockedUnd(ctx, data.TotalLocked)
	if err != nil {
		panic(err)
	}

	for _, po := range data.PurchaseOrders {
		err = keeper.SetPurchaseOrder(ctx, po)
		if err != nil {
			panic(err)
		}
	}

	for _, lund := range data.LockedUnds {
		err = keeper.SetLockedUndForAccount(ctx, lund)
		if err != nil {
			panic(err)
		}
	}

	// ensure locked UND is registered with supply keeper
	if moduleAcc.GetCoins().IsZero() {
		var moduleHoldings sdk.Coins
		moduleHoldings = moduleHoldings.Add(data.TotalLocked)
		if err := moduleAcc.SetCoins(moduleHoldings); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	purchaseOrderId, _ := keeper.GetHighestPurchaseOrderID(ctx)
	purchaseOrders := keeper.GetAllPurchaseOrders(ctx)
	lockedUnds := keeper.GetAllLockedUnds(ctx)
	totalLocked := keeper.GetTotalLockedUnd(ctx)
	whitelist := keeper.GetAllWhitelistedAddresses(ctx)

	return GenesisState{
		Params:                  params,
		StartingPurchaseOrderID: purchaseOrderId,
		PurchaseOrders:          purchaseOrders,
		LockedUnds:              lockedUnds,
		TotalLocked:             totalLocked,
		Whitelist:               whitelist,
	}
}
