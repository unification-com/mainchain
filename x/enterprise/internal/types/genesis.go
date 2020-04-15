package types

import (
	"bytes"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/types"
)

// GenesisState - enterprise state
type GenesisState struct {
	Params                  Params             `json:"params" yaml:"params"`                                         // enterprise params
	StartingPurchaseOrderID uint64             `json:"starting_purchase_order_id" yaml:"starting_purchase_order_id"` // should be 1
	PurchaseOrders          PurchaseOrders     `json:"purchase_orders" yaml:"purchase_orders"`
	LockedUnds              LockedUnds         `json:"locked_und" yaml:"locked_und"`
	TotalLocked             sdk.Coin           `json:"total_locked" yaml:"total_locked"`
	Whitelist               WhitelistAddresses `json:"whitelist" yaml:"whitelist"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingPurchaseOrderID uint64, totalLocked sdk.Coin) GenesisState {
	return GenesisState{
		Params:                  params,
		StartingPurchaseOrderID: startingPurchaseOrderID,
		TotalLocked:             totalLocked,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		DefaultParams(),
		DefaultStartingPurchaseOrderID,
		sdk.NewInt64Coin(types.DefaultDenomination, 0),
	)
}

// Equal checks whether two enterprise GenesisState structs are equivalent
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	if data.StartingPurchaseOrderID == 0 {
		return fmt.Errorf("enterprise starting purchase order id should be greater than 0")
	}

	for _, po := range data.PurchaseOrders {
		if po.PurchaseOrderID == 0 {
			return fmt.Errorf("invalid purchase order: PurchaseOrderID: %d. Error: Missing PurchaseOrderID", po.PurchaseOrderID)
		}
		if po.Purchaser == nil {
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
			if decision.Signer.Empty() {
				return fmt.Errorf("invalid decision: Signer cannot be empty")
			}
			if !ValidPurchaseOrderAcceptRejectStatus(decision.Decision) {
				return fmt.Errorf("invalid decision: Decision: %s. Error: Invalid Decision", decision.Decision)
			}
		}
	}

	return nil
}
