package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewEnterpriseUndPurchaseOrder is currently only used in the simulation decoder unit tests
func NewEnterpriseUndPurchaseOrder(poId uint64, purchaser string, amount sdk.Coin, status PurchaseOrderStatus, raisedTime, completionTime uint64) (EnterpriseUndPurchaseOrder, error) {
	po := EnterpriseUndPurchaseOrder{
		Id:             poId,
		Purchaser:      purchaser,
		Amount:         amount,
		Status:         status,
		RaiseTime:      raisedTime,
		CompletionTime: completionTime,
		Decisions:      nil,
	}

	return po, nil
}

// NewLockedUnd is currently only used in the simulation decoder unit tests
func NewLockedUnd(owner string, amount sdk.Coin) (LockedUnd, error) {
	l := LockedUnd{
		Owner:  owner,
		Amount: amount,
	}

	return l, nil
}
