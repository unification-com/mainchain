package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	PurchaseAction         = "raise_ent_po"
	ProcessAction          = "proc_ent_po"
	WhitelistAddressAction = "ent_whitelist"
)

// __Enterprise_UND_Purchase_Order_Msg__________________________________

var (
	_ sdk.Msg = &MsgUndPurchaseOrder{}
	_ sdk.Msg = &MsgProcessUndPurchaseOrder{}
	_ sdk.Msg = &MsgWhitelistAddress{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// NewMsgUndPurchaseOrder is a constructor function for MsgUndPurchaseOrder
//
//nolint:interfacer
func NewMsgUndPurchaseOrder(purchaser sdk.AccAddress, amount sdk.Coin) *MsgUndPurchaseOrder {
	return &MsgUndPurchaseOrder{Purchaser: purchaser.String(), Amount: amount}
}

// Route should return the name of the module
func (msg MsgUndPurchaseOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUndPurchaseOrder) Type() string { return PurchaseAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgUndPurchaseOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Purchaser)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid purchaser address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount must be greater than zero")
	}
	if msg.Amount.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount must be a positive value")
	}
	return nil
}

// __Enterprise_UND_Process_Purchase_Order_Msg__________________________

// MsgProcessUndPurchaseOrder defines a ProcessUndPurchaseOrder message - used to accept/reject a PO

// NewMsgProcessUndPurchaseOrder is a constructor function for MsgProcessUndPurchaseOrder
func NewMsgProcessUndPurchaseOrder(purchaseOrderID uint64, decision PurchaseOrderStatus, signer sdk.AccAddress) *MsgProcessUndPurchaseOrder {
	return &MsgProcessUndPurchaseOrder{
		PurchaseOrderId: purchaseOrderID,
		Decision:        decision,
		Signer:          signer.String(),
	}
}

// Route should return the name of the module
func (msg MsgProcessUndPurchaseOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgProcessUndPurchaseOrder) Type() string { return ProcessAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgProcessUndPurchaseOrder) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid signer address (%s)", err)
	}

	if msg.PurchaseOrderId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "purchase order id must be greater than zero")
	}
	if !ValidPurchaseOrderAcceptRejectStatus(msg.Decision) {
		return errorsmod.Wrap(ErrInvalidStatus, "status must be accept or reject")
	}
	return nil
}

// __Enterprise_UND_Whitelist_Msg__________________________

// MsgWhitelistAddress defines a WhitelistAddress message - used to add/remove addresses from PO whitelist
// and determine which addresses are allowed to raise purchase orders

// NewMsgWhitelistAddress is a constructor function for MsgWhitelistAddress
func NewMsgWhitelistAddress(address sdk.AccAddress, action WhitelistAction, signer sdk.AccAddress) *MsgWhitelistAddress {
	return &MsgWhitelistAddress{
		Address: address.String(),
		Signer:  signer.String(),
		Action:  action,
	}
}

// Route should return the name of the module
func (msg MsgWhitelistAddress) Route() string { return RouterKey }

// Type should return the action
func (msg MsgWhitelistAddress) Type() string { return WhitelistAddressAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgWhitelistAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid signer address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}
	if !ValidWhitelistAction(msg.Action) {
		return errorsmod.Wrap(ErrInvalidWhitelistAction, "action must be add or remove")
	}
	return nil
}

// --- Modify Params Msg Type ---

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := m.Params.Validate(); err != nil {
		return err
	}

	return nil
}
