package v045

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v040 "github.com/unification-com/mainchain/x/enterprise/legacy/v040"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// MigrateJSON accepts exported v0.40 x/enterprise genesis state and migrates it to
// v0.45 (0.43) x/enterprise genesis state.
func MigrateJSON(oldEnterpriseState *v040.GenesisState) *types.GenesisState {

	// prefill with 0nund
	completedByAcc := make(map[string]sdk.Coin)
	for _, u := range oldEnterpriseState.LockedUnd {
		completedByAcc[u.Owner] = sdk.NewInt64Coin(oldEnterpriseState.Params.Denom, int64(0))
	}

	// total spent & POs
	newPos := make(types.EnterpriseUndPurchaseOrders, len(oldEnterpriseState.PurchaseOrders))
	totalCompleted := sdk.NewInt64Coin(oldEnterpriseState.Params.Denom, int64(0))
	for i, oldPo := range oldEnterpriseState.PurchaseOrders {

		newPoDecisions := make(types.PurchaseOrderDecisions, len(oldPo.Decisions))

		for j, oldDec := range oldPo.Decisions {
			newPoDecisions[j] = types.PurchaseOrderDecision{
				Signer:       oldDec.Signer,
				Decision:     types.PurchaseOrderStatus(oldDec.Decision),
				DecisionTime: oldDec.DecisionTime,
			}
		}

		newPos[i] = types.EnterpriseUndPurchaseOrder{
			Id:             oldPo.Id,
			Purchaser:      oldPo.Purchaser,
			Amount:         oldPo.Amount,
			Status:         types.PurchaseOrderStatus(oldPo.Status),
			RaiseTime:      oldPo.RaiseTime,
			CompletionTime: oldPo.CompletionTime,
			Decisions:      newPoDecisions,
		}
		if types.PurchaseOrderStatus(oldPo.Status) == types.StatusCompleted {
			// PO was completed and purchased amount minted
			totalCompleted = totalCompleted.Add(oldPo.Amount)

			completedByAcc[oldPo.Purchaser] = completedByAcc[oldPo.Purchaser].Add(oldPo.Amount)

		}
	}

	totalSpent := totalCompleted.Sub(oldEnterpriseState.TotalLocked)

	// whitelist
	newWl := make(types.Whitelists, len(oldEnterpriseState.Whitelist))
	for i, wl := range oldEnterpriseState.Whitelist {
		newWl[i] = wl
	}

	// locked & spent by user
	newLocked := make(types.LockedUnds, len(oldEnterpriseState.LockedUnd))
	newSpent := make(types.SpentEFUNDs, len(oldEnterpriseState.LockedUnd))
	for i, oldLocked := range oldEnterpriseState.LockedUnd {
		newLocked[i] = types.LockedUnd{
			Owner:  oldLocked.Owner,
			Amount: oldLocked.Amount,
		}
		newSpent[i] = types.SpentEFUND{
			Owner:  oldLocked.Owner,
			Amount: completedByAcc[oldLocked.Owner].Sub(oldLocked.Amount),
		}
	}

	return &types.GenesisState{
		Params: types.Params{
			EntSigners:        oldEnterpriseState.Params.EntSigners,
			Denom:             oldEnterpriseState.Params.Denom,
			MinAccepts:        oldEnterpriseState.Params.MinAccepts,
			DecisionTimeLimit: oldEnterpriseState.Params.DecisionTimeLimit,
		},
		StartingPurchaseOrderId: oldEnterpriseState.StartingPurchaseOrderId,
		PurchaseOrders:          newPos,
		LockedUnd:               newLocked,
		TotalLocked:             oldEnterpriseState.TotalLocked,
		Whitelist:               newWl,
		SpentEfund:              newSpent,
		TotalSpent:              totalSpent,
	}
}
