package v040

import (
	v038 "github.com/unification-com/mainchain/x/enterprise/legacy/v038"
	v040 "github.com/unification-com/mainchain/x/enterprise/types"
)

// convertDecision convert an old byte decision to an enum
func convertDecision(oldDecision v038.PurchaseOrderStatus) v040.PurchaseOrderStatus {
	switch oldDecision {
	case v038.StatusNil:
		return v040.StatusNil
	case v038.StatusRaised:
		return v040.StatusRaised
	case v038.StatusAccepted:
		return v040.StatusAccepted
	case v038.StatusRejected:
		return v040.StatusRejected
	case v038.StatusCompleted:
		return v040.StatusCompleted
	default:
		return v040.StatusNil
	}
}

func Migrate(oldEnterpriseState v038.GenesisState) *v040.GenesisState {

	newPos := make(v040.EnterpriseUndPurchaseOrders, len(oldEnterpriseState.PurchaseOrders))

	for i, oldPo := range oldEnterpriseState.PurchaseOrders {
		newDecisions := make(v040.PurchaseOrderDecisions, len(oldPo.Decisions))
		for j, oldDecision := range oldPo.Decisions {
			newDecisions[j] = v040.PurchaseOrderDecision{
				Signer:       oldDecision.Signer.String(),
				Decision:     convertDecision(oldDecision.Decision),
				DecisionTime: uint64(oldDecision.DecisionTime),
			}
		}

		newPos[i] = v040.EnterpriseUndPurchaseOrder{
			Id:             oldPo.PurchaseOrderID,
			Purchaser:      oldPo.Purchaser.String(),
			Amount:         oldPo.Amount,
			Status:         convertDecision(oldPo.Status),
			RaiseTime:      uint64(oldPo.RaisedTime),
			CompletionTime: uint64(oldPo.CompletionTime),
			Decisions:      newDecisions,
		}
	}

	newLockedUnd := make(v040.LockedUnds, len(oldEnterpriseState.LockedUnds))

	for i, oldLockedUnd := range oldEnterpriseState.LockedUnds {
		newLockedUnd[i] = v040.LockedUnd{
			Owner:  oldLockedUnd.Owner.String(),
			Amount: oldLockedUnd.Amount,
		}
	}

	newWhiteList := make(v040.Whitelists, len(oldEnterpriseState.Whitelist))

	for i, oldWl := range oldEnterpriseState.Whitelist {
		newWhiteList[i] = oldWl.String()
	}

	return &v040.GenesisState{
		Params: v040.Params{
			EntSigners:        oldEnterpriseState.Params.EntSigners,
			Denom:             oldEnterpriseState.Params.Denom,
			MinAccepts:        oldEnterpriseState.Params.MinAccepts,
			DecisionTimeLimit: oldEnterpriseState.Params.DecisionLimit,
		},
		StartingPurchaseOrderId: oldEnterpriseState.StartingPurchaseOrderID,
		PurchaseOrders:          newPos,
		LockedUnd:               newLockedUnd,
		TotalLocked:             oldEnterpriseState.TotalLocked,
		Whitelist:               newWhiteList,
	}
}
