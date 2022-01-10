package enterprise

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// InitGenesis new enterprise FUND genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestPurchaseOrderID(ctx, data.StartingPurchaseOrderId)

	moduleAcc := keeper.GetEnterpriseAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if data.Whitelist != nil {
		for _, wlAddr := range data.Whitelist {
			addr, err := sdk.AccAddressFromBech32(wlAddr)
			if err != nil {
				panic(err)
			}
			err = keeper.AddAddressToWhitelist(ctx, addr)
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
		epo := types.EnterpriseUndPurchaseOrder{
			Id:             po.Id,
			Purchaser:      po.Purchaser,
			Amount:         po.Amount,
			Status:         po.Status,
			RaiseTime:      po.RaiseTime,
			CompletionTime: po.CompletionTime,
			Decisions:      po.Decisions,
		}
		err = keeper.SetPurchaseOrder(ctx, epo)
		if err != nil {
			panic(err)
		}
	}

	for _, lund := range data.LockedUnd {
		locked := types.LockedUnd{
			Owner:  lund.Owner,
			Amount: lund.Amount,
		}
		err = keeper.SetLockedUndForAccount(ctx, locked)
		if err != nil {
			panic(err)
		}
	}

	// ensure locked FUND is registered with supply keeper
	balances := bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		var moduleHoldings sdk.Coins
		moduleHoldings = moduleHoldings.Add(data.TotalLocked)
		if err := bankKeeper.SetBalances(ctx, moduleAcc.GetAddress(), moduleHoldings); err != nil {
			panic(err)
		}

		accountKeeper.SetModuleAccount(ctx, moduleAcc)
	}
	//if moduleAcc.GetCoins().IsZero() {
	//	var moduleHoldings sdk.Coins
	//	moduleHoldings = moduleHoldings.Add(data.TotalLocked)
	//	if err := moduleAcc.SetCoins(moduleHoldings); err != nil {
	//		panic(err)
	//	}
	//	bankKeeper.SetModuleAccount(ctx, moduleAcc)
	//}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	params := keeper.GetParams(ctx)
	purchaseOrderId, _ := keeper.GetHighestPurchaseOrderID(ctx)
	purchaseOrders := keeper.GetAllPurchaseOrders(ctx)
	lockedUnds := keeper.GetAllLockedUnds(ctx)
	totalLocked := keeper.GetTotalLockedUnd(ctx)
	whitelist := keeper.GetAllWhitelistedAddresses(ctx)

	return &types.GenesisState{
		Params:                  params,
		StartingPurchaseOrderId: purchaseOrderId,
		PurchaseOrders:          purchaseOrders,
		LockedUnd:               lockedUnds,
		TotalLocked:             totalLocked,
		Whitelist:               whitelist,
	}
}
