package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	undtypes "github.com/unification-com/mainchain/types"
)

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	if data.StartingPurchaseOrderId == 0 {
		return fmt.Errorf("enterprise starting purchase order id should be greater than 0")
	}

	for _, po := range data.PurchaseOrders {
		if po.Id == 0 {
			return fmt.Errorf("invalid purchase order: PurchaseOrderID: %d. Error: Missing PurchaseOrderID", po.Id)
		}

		_, err := sdk.AccAddressFromBech32(po.Purchaser)
		if err != nil {
			return fmt.Errorf("invalid purchase order: Purchaser: %s. Error: Missing Purchaser", po.Purchaser)
		}

		if !po.Amount.IsValid() {
			return fmt.Errorf("invalid purchase order: Amount: %s. Error: Missing Amount", po.Amount.Amount)
		}
		if po.Amount.IsZero() || po.Amount.IsNegative() {
			return fmt.Errorf("invalid purchase order: Amount. Error: Amount must be greater than 0")
		}
		if !ValidPurchaseOrderStatus(po.Status) {
			return fmt.Errorf("invalid purchase order: Status: %s. Error: Invalid Status", po.Status)
		}

		for _, decision := range po.Decisions {
			_, err := sdk.AccAddressFromBech32(decision.Signer)
			if err != nil {
				return fmt.Errorf("invalid purchase order: Purchaser: %s. Error: Missing Purchaser", po.Purchaser)
			}
			if !ValidPurchaseOrderAcceptRejectStatus(decision.Decision) {
				return fmt.Errorf("invalid decision: Decision: %s. Error: Invalid Decision", decision.Decision)
			}
		}
	}

	return nil
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingPurchaseOrderId uint64, totalLocked sdk.Coin,
	purchaseOrders EnterpriseUndPurchaseOrders, locked LockedUnds, whitelist Whitelists) *GenesisState {
	return &GenesisState{
		Params:                  params,
		StartingPurchaseOrderId: startingPurchaseOrderId,
		PurchaseOrders:          purchaseOrders,
		LockedUnd:               locked,
		TotalLocked:             totalLocked,
		Whitelist:               whitelist,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		DefaultParams(),
		1,
		sdk.NewInt64Coin(undtypes.DefaultDenomination, 0),
		nil, nil, nil,
	)
}
