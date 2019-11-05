package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RouterKey = ModuleName // defined in keys.go file

	PurchaseAction = "purchase"
)

// --- Enterprise UND Purcahse Order Msg ---

// MsgRaiseUndPurchaseOrder defines a PurchaseUnd message
type MsgRaiseUndPurchaseOrder struct {
	Purchaser  sdk.AccAddress `json:"purchaser"`
	Amount     sdk.Coin       `json:"amount"`
}

// NewMsgRaiseUndPurchaseOrder is a constructor function for MsgRaiseUndPurchaseOrder
func NewMsgRaiseUndPurchaseOrder(purchaser sdk.AccAddress, amount sdk.Coin) MsgRaiseUndPurchaseOrder {
	return MsgRaiseUndPurchaseOrder{
		Purchaser:  purchaser,
		Amount:     amount,
	}
}

// Route should return the name of the module
func (msg MsgRaiseUndPurchaseOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRaiseUndPurchaseOrder) Type() string { return PurchaseAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRaiseUndPurchaseOrder) ValidateBasic() sdk.Error {
	if msg.Purchaser.Empty() {
		return sdk.ErrInvalidAddress(msg.Purchaser.String())
	}
	if msg.Amount.IsZero() {
		return sdk.ErrInvalidCoins("amount must be greater than zero")
	}
	if msg.Amount.IsNegative() {
		return sdk.ErrInvalidCoins("amount must be a positive value")
	}
	if msg.Amount.Denom != "nund" { // Todo - take from global app types/denom.go
		return sdk.ErrInvalidCoins("denomination must be in nano UND - nund")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRaiseUndPurchaseOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRaiseUndPurchaseOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Purchaser} // ToDo: see if we can get this from genesis/keeper params
}
