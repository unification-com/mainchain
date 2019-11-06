package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RouterKey = ModuleName // defined in keys.go file

	PurchaseAction = "purchase"
	ProcessAction  = "process"
)

// __Enterprise_UND_Purchase_Order_Msg__________________________________

// MsgUndPurchaseOrder defines a UndPurchaseOrder message
type MsgUndPurchaseOrder struct {
	Purchaser  sdk.AccAddress `json:"purchaser"`
	Amount     sdk.Coin       `json:"amount"`
}

// NewMsgUndPurchaseOrder is a constructor function for MsgUndPurchaseOrder
func NewMsgUndPurchaseOrder(purchaser sdk.AccAddress, amount sdk.Coin) MsgUndPurchaseOrder {
	return MsgUndPurchaseOrder{
		Purchaser:  purchaser,
		Amount:     amount,
	}
}

// Route should return the name of the module
func (msg MsgUndPurchaseOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUndPurchaseOrder) Type() string { return PurchaseAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgUndPurchaseOrder) ValidateBasic() sdk.Error {
	if msg.Purchaser.Empty() {
		return sdk.ErrInvalidAddress(msg.Purchaser.String())
	}
	if msg.Amount.IsZero() {
		return sdk.ErrInvalidCoins("amount must be greater than zero")
	}
	if msg.Amount.IsNegative() {
		return sdk.ErrInvalidCoins("amount must be a positive value")
	}
	if msg.Amount.Denom != DefaultDenomination {
		return sdk.ErrInvalidCoins(fmt.Sprintf("denomination must be in %s", DefaultDenomination))
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUndPurchaseOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUndPurchaseOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Purchaser}
}

// __Enterprise_UND_Process_Purchase_Order_Msg__________________________

// MsgProcessUndPurchaseOrder defines a ProcessUndPurchaseOrder message - used to accept/reject a PO
type MsgProcessUndPurchaseOrder struct {
	PurchaseOrderID uint64              `json:"id"`
	Decision        PurchaseOrderStatus `json:"status"`
	Signer          sdk.AccAddress      `json:"signer"`
}

// NewMsgProcessUndPurchaseOrder is a constructor function for MsgProcessUndPurchaseOrder
func NewMsgProcessUndPurchaseOrder(purchaseOrderID uint64, decision PurchaseOrderStatus, signer sdk.AccAddress) MsgProcessUndPurchaseOrder {
	return MsgProcessUndPurchaseOrder{
		PurchaseOrderID: purchaseOrderID,
		Decision:        decision,
		Signer:          signer,
	}
}

// Route should return the name of the module
func (msg MsgProcessUndPurchaseOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgProcessUndPurchaseOrder) Type() string { return ProcessAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgProcessUndPurchaseOrder) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.PurchaseOrderID <= 0 {
		return sdk.ErrUnknownRequest("purchase order id must be greater than zero")
	}
	if !ValidPurchaseOrderAcceptRejectStatus(msg.Decision) {
		return sdk.ErrUnknownRequest("status must be accept or reject")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgProcessUndPurchaseOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProcessUndPurchaseOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
