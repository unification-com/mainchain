package enterprise

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis new enterprise UND genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestPurchaseOrderID(ctx, data.StartingPurchaseOrderID)

	moduleAcc := keeper.GetEnterpriseAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", ModuleName))
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

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	purchaseOrderId, _ := keeper.GetHighestPurchaseOrderID(ctx)
	purchaseOrders := keeper.GetAllPurchaseOrders(ctx)
	lockedUnds := keeper.GetAllLockedUnds(ctx)
	totalLocked := keeper.GetTotalLockedUnd(ctx)

	return GenesisState{
		Params:                  params,
		StartingPurchaseOrderID: purchaseOrderId,
		PurchaseOrders:          purchaseOrders,
		LockedUnds:              lockedUnds,
		TotalLocked:             totalLocked,
	}
}
