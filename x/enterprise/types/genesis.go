package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	if data.StartingPurchaseOrderId == 0 {
		return fmt.Errorf("enterprise starting purchase order id should be greater than 0")
	}

	denom := data.Params.Denom

	// TotalLocked / TotalSpent denom must match Params.Denom
	if data.TotalLocked.Denom != denom {
		return fmt.Errorf("TotalLocked denom %q does not match Params.Denom %q", data.TotalLocked.Denom, denom)
	}
	if data.TotalSpent.Denom != denom {
		return fmt.Errorf("TotalSpent denom %q does not match Params.Denom %q", data.TotalSpent.Denom, denom)
	}

	seenPOIDs := make(map[uint64]struct{}, len(data.PurchaseOrders))
	for i, po := range data.PurchaseOrders {
		if po.Id == 0 {
			return fmt.Errorf("purchaseOrder[%d]: invalid PurchaseOrderID 0", i)
		}
		if _, dup := seenPOIDs[po.Id]; dup {
			return fmt.Errorf("purchaseOrder[%d]: duplicate PurchaseOrderID %d", i, po.Id)
		}
		seenPOIDs[po.Id] = struct{}{}
		if po.Id >= data.StartingPurchaseOrderId {
			return fmt.Errorf("purchaseOrder[%d]: Id %d must be < StartingPurchaseOrderId %d",
				i, po.Id, data.StartingPurchaseOrderId)
		}

		if _, err := sdk.AccAddressFromBech32(po.Purchaser); err != nil {
			return fmt.Errorf("purchaseOrder[%d]: invalid Purchaser %q: %w", i, po.Purchaser, err)
		}

		if !po.Amount.IsValid() {
			return fmt.Errorf("purchaseOrder[%d]: invalid Amount %s", i, po.Amount)
		}
		if po.Amount.IsZero() || po.Amount.IsNegative() {
			return fmt.Errorf("purchaseOrder[%d]: Amount must be > 0", i)
		}
		if po.Amount.Denom != denom {
			return fmt.Errorf("purchaseOrder[%d]: Amount denom %q does not match Params.Denom %q", i, po.Amount.Denom, denom)
		}
		if !ValidPurchaseOrderStatus(po.Status) {
			return fmt.Errorf("purchaseOrder[%d]: invalid Status %s", i, po.Status)
		}

		seenSigners := make(map[string]struct{}, len(po.Decisions))
		for j, decision := range po.Decisions {
			if _, err := sdk.AccAddressFromBech32(decision.Signer); err != nil {
				return fmt.Errorf("purchaseOrder[%d].decision[%d]: invalid Signer %q: %w",
					i, j, decision.Signer, err)
			}
			if _, dup := seenSigners[decision.Signer]; dup {
				return fmt.Errorf("purchaseOrder[%d].decision[%d]: duplicate Signer %q",
					i, j, decision.Signer)
			}
			seenSigners[decision.Signer] = struct{}{}
			if !ValidPurchaseOrderAcceptRejectStatus(decision.Decision) {
				return fmt.Errorf("purchaseOrder[%d].decision[%d]: invalid Decision %s",
					i, j, decision.Decision)
			}
		}
	}

	// Locked eFUND: validate each entry + invariant: sum == TotalLocked
	seenLockedOwners := make(map[string]struct{}, len(data.LockedUnd))
	lockedSum := sdk.NewInt64Coin(denom, 0)
	for i, locked := range data.LockedUnd {
		if _, err := sdk.AccAddressFromBech32(locked.Owner); err != nil {
			return fmt.Errorf("lockedUnd[%d]: invalid Owner %q: %w", i, locked.Owner, err)
		}
		if _, dup := seenLockedOwners[locked.Owner]; dup {
			return fmt.Errorf("lockedUnd[%d]: duplicate Owner %q", i, locked.Owner)
		}
		seenLockedOwners[locked.Owner] = struct{}{}
		if !locked.Amount.IsValid() {
			return fmt.Errorf("lockedUnd[%d]: invalid Amount %s", i, locked.Amount)
		}
		if locked.Amount.IsNegative() {
			return fmt.Errorf("lockedUnd[%d]: negative Amount %s", i, locked.Amount)
		}
		if locked.Amount.Denom != denom {
			return fmt.Errorf("lockedUnd[%d]: Amount denom %q does not match Params.Denom %q",
				i, locked.Amount.Denom, denom)
		}
		lockedSum = lockedSum.Add(locked.Amount)
	}
	if !lockedSum.IsEqual(data.TotalLocked) {
		return fmt.Errorf("sum of LockedUnd amounts (%s) does not equal TotalLocked (%s)",
			lockedSum, data.TotalLocked)
	}

	// Spent eFUND: same invariants
	seenSpentOwners := make(map[string]struct{}, len(data.SpentEfund))
	spentSum := sdk.NewInt64Coin(denom, 0)
	for i, spent := range data.SpentEfund {
		if _, err := sdk.AccAddressFromBech32(spent.Owner); err != nil {
			return fmt.Errorf("spentEfund[%d]: invalid Owner %q: %w", i, spent.Owner, err)
		}
		if _, dup := seenSpentOwners[spent.Owner]; dup {
			return fmt.Errorf("spentEfund[%d]: duplicate Owner %q", i, spent.Owner)
		}
		seenSpentOwners[spent.Owner] = struct{}{}
		if !spent.Amount.IsValid() {
			return fmt.Errorf("spentEfund[%d]: invalid Amount %s", i, spent.Amount)
		}
		if spent.Amount.IsNegative() {
			return fmt.Errorf("spentEfund[%d]: negative Amount %s", i, spent.Amount)
		}
		if spent.Amount.Denom != denom {
			return fmt.Errorf("spentEfund[%d]: Amount denom %q does not match Params.Denom %q",
				i, spent.Amount.Denom, denom)
		}
		spentSum = spentSum.Add(spent.Amount)
	}
	if !spentSum.IsEqual(data.TotalSpent) {
		return fmt.Errorf("sum of SpentEfund amounts (%s) does not equal TotalSpent (%s)",
			spentSum, data.TotalSpent)
	}

	// Whitelist: validate each address + dup check
	seenWhitelist := make(map[string]struct{}, len(data.Whitelist))
	for i, addr := range data.Whitelist {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("whitelist[%d]: invalid address %q: %w", i, addr, err)
		}
		if _, dup := seenWhitelist[addr]; dup {
			return fmt.Errorf("whitelist[%d]: duplicate address %q", i, addr)
		}
		seenWhitelist[addr] = struct{}{}
	}

	return nil
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingPurchaseOrderId uint64, totalLocked sdk.Coin,
	purchaseOrders EnterpriseUndPurchaseOrders, locked LockedUnds, whitelist Whitelists, totalSpent sdk.Coin, spentEFUND SpentEFUNDs) *GenesisState {
	return &GenesisState{
		Params:                  params,
		StartingPurchaseOrderId: startingPurchaseOrderId,
		PurchaseOrders:          purchaseOrders,
		LockedUnd:               locked,
		TotalLocked:             totalLocked,
		Whitelist:               whitelist,
		TotalSpent:              totalSpent,
		SpentEfund:              spentEFUND,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		DefaultParams(),
		1,
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
		nil, nil, nil,
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
		nil,
	)
}
