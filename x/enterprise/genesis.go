package enterprise

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// InitGenesis new enterprise FUND genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, data types.GenesisState) []abci.ValidatorUpdate {
	moduleAcc := keeper.GetEnterpriseAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestPurchaseOrderID(ctx, data.StartingPurchaseOrderId)

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

	err = keeper.SetTotalSpentEFUND(ctx, data.TotalSpent)
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

		// if export occured during decision making, needs to be
		// set in keeper
		if po.Status == types.StatusRaised {
			keeper.AddPoToRaisedQueue(ctx, po.Id)
		}
		if po.Status == types.StatusAccepted {
			keeper.AddPoToAcceptedQueue(ctx, po.Id)
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

	for _, spent := range data.SpentEfund {
		err = keeper.SetSpentEFUNDForAccount(ctx, spent)
		if err != nil {
			panic(err)
		}
	}

	// ensure locked FUND is registered with supply keeper
	var moduleHoldings sdk.Coins
	moduleHoldings = moduleHoldings.Add(data.TotalLocked)

	balances := bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		accountKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	if !balances.Equal(moduleHoldings) {
		panic(fmt.Sprintf("enterprise module balance does not match the module holdings: %s <-> %s", balances, moduleHoldings))
	}

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
	totalSpent := keeper.GetTotalSpentEFUND(ctx)
	spentEFUNDs := keeper.GetAllSpentEFUNDs(ctx)

	return &types.GenesisState{
		Params:                  params,
		StartingPurchaseOrderId: purchaseOrderId,
		PurchaseOrders:          purchaseOrders,
		LockedUnd:               lockedUnds,
		TotalLocked:             totalLocked,
		Whitelist:               whitelist,
		TotalSpent:              totalSpent,
		SpentEfund:              spentEFUNDs,
	}
}
