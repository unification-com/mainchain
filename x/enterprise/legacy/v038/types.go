package v038

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// module name
	ModuleName = "enterprise"

	RouterKey = ModuleName // defined in keys.go file

	PurchaseAction         = "raise_enterprise_purchase_order"
	ProcessAction          = "process_enterprise_purchase_order"
	WhitelistAddressAction = "whitelist_enterprise_purchase_order_address"

	// Valid Purchase Order statuses
	StatusNil       PurchaseOrderStatus = 0x00
	StatusRaised    PurchaseOrderStatus = 0x01
	StatusAccepted  PurchaseOrderStatus = 0x02
	StatusRejected  PurchaseOrderStatus = 0x03
	StatusCompleted PurchaseOrderStatus = 0x04

	WhitelistActionAdd    WhitelistAction = 0x01
	WhitelistActionRemove WhitelistAction = 0x02
)

type (
	PurchaseOrderStatus byte
	WhitelistAction     byte

	PurchaseOrders     []EnterpriseUndPurchaseOrder
	LockedUnds         []LockedUnd
	UndSupplies        []UndSupply
	WhitelistAddresses []sdk.AccAddress

	Params struct {
		EntSigners    string `json:"ent_signers" yaml:"ent_signers"` // Accounts allowed to sign decisions on FUND purchase orders
		Denom         string `json:"denom" yaml:"denom"`
		MinAccepts    uint64 `json:"min_Accepts" yaml:"min_Accepts"`                 // must be <= len(EntSigners)
		DecisionLimit uint64 `json:"decision_time_limit" yaml:"decision_time_limit"` // num seconds elapsed before auto-reject
	}

	EnterpriseUndPurchaseOrder struct {
		PurchaseOrderID uint64                  `json:"id"`
		Purchaser       sdk.AccAddress          `json:"purchaser"`
		Amount          sdk.Coin                `json:"amount"`
		Status          PurchaseOrderStatus     `json:"status"`
		RaisedTime      int64                   `json:"raise_time"`
		Decisions       []PurchaseOrderDecision `json:"decisions"`
		CompletionTime  int64                   `json:"completion_time"`
	}

	PurchaseOrderDecision struct {
		Signer       sdk.AccAddress      `json:"signer"`
		Decision     PurchaseOrderStatus `json:"decision"`
		DecisionTime int64               `json:"decision_time"`
	}

	LockedUnd struct {
		Owner  sdk.AccAddress `json:"owner"`
		Amount sdk.Coin       `json:"amount"`
	}

	UndSupply struct {
		Denom  string `json:"denom"`
		Amount int64  `json:"amount"`
		Locked int64  `json:"locked"`
		Total  int64  `json:"total"`
	}

	GenesisState struct {
		Params                  Params             `json:"params" yaml:"params"`                                         // enterprise params
		StartingPurchaseOrderID uint64             `json:"starting_purchase_order_id" yaml:"starting_purchase_order_id"` // should be 1
		PurchaseOrders          PurchaseOrders     `json:"purchase_orders" yaml:"purchase_orders"`
		LockedUnds              LockedUnds         `json:"locked_und" yaml:"locked_und"`
		TotalLocked             sdk.Coin           `json:"total_locked" yaml:"total_locked"`
		Whitelist               WhitelistAddresses `json:"whitelist" yaml:"whitelist"`
	}
)

// PurchaseOrderStatusFromString turns a string into a ProposalStatus
func PurchaseOrderStatusFromString(str string) (PurchaseOrderStatus, error) {
	switch str {
	case "accept":
		return StatusAccepted, nil

	case "reject":
		return StatusRejected, nil

	case "raised":
		return StatusRaised, nil

	case "complete":
		return StatusCompleted, nil

	case "":
		return StatusNil, nil

	default:
		return PurchaseOrderStatus(0xff), fmt.Errorf("'%s' is not a valid purchase order status", str)
	}
}

// UnmarshalJSON Unmarshals from JSON
func (status *PurchaseOrderStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := PurchaseOrderStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// __Enterprise_UND_Purchase_Order_Msg__________________________________

// MsgUndPurchaseOrder defines a UndPurchaseOrder message
type MsgUndPurchaseOrder struct {
	Purchaser sdk.AccAddress `json:"purchaser"`
	Amount    sdk.Coin       `json:"amount"`
}

// NewMsgUndPurchaseOrder is a constructor function for MsgUndPurchaseOrder
func NewMsgUndPurchaseOrder(purchaser sdk.AccAddress, amount sdk.Coin) MsgUndPurchaseOrder {
	return MsgUndPurchaseOrder{
		Purchaser: purchaser,
		Amount:    amount,
	}
}

// Route should return the name of the module
func (msg MsgUndPurchaseOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUndPurchaseOrder) Type() string { return PurchaseAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgUndPurchaseOrder) ValidateBasic() error {
	if msg.Purchaser.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Purchaser.String())
	}
	if msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be greater than zero")
	}
	if msg.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be a positive value")
	}
	return nil
}

// __Enterprise_UND_Process_Purchase_Order_Msg__________________________

// MsgProcessUndPurchaseOrder defines a ProcessUndPurchaseOrder message - used to accept/reject a PO
type MsgProcessUndPurchaseOrder struct {
	PurchaseOrderID uint64              `json:"id"`
	Decision        PurchaseOrderStatus `json:"decision"`
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
func (msg MsgProcessUndPurchaseOrder) ValidateBasic() error {
	if msg.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer.String())
	}
	if msg.PurchaseOrderID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "purchase order id must be greater than zero")
	}
	if !ValidPurchaseOrderAcceptRejectStatus(msg.Decision) {
		return fmt.Errorf("status must be accept or reject")
	}
	return nil
}

// __Enterprise_UND_Whitelist_Msg__________________________

// MsgWhitelistAddress defines a WhitelistAddress message - used to add/remove addresses from PO whitelist
// and determine which addresses are allowed to raise purchase orders
type MsgWhitelistAddress struct {
	Address sdk.AccAddress  `json:"address"`
	Signer  sdk.AccAddress  `json:"signer"`
	Action  WhitelistAction `json:"action"`
}

// NewMsgWhitelistAddress is a constructor function for MsgWhitelistAddress
func NewMsgWhitelistAddress(address sdk.AccAddress, action WhitelistAction, signer sdk.AccAddress) MsgWhitelistAddress {
	return MsgWhitelistAddress{
		Address: address,
		Signer:  signer,
		Action:  action,
	}
}

// Route should return the name of the module
func (msg MsgWhitelistAddress) Route() string { return RouterKey }

// Type should return the action
func (msg MsgWhitelistAddress) Type() string { return WhitelistAddressAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgWhitelistAddress) ValidateBasic() error {
	if msg.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer.String())
	}
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Address.String())
	}
	if !ValidWhitelistAction(msg.Action) {
		return fmt.Errorf("action must be add or remove")
	}
	return nil
}

// ValidPurchaseOrderAcceptRejectStatus checks the decision - returns true if accept/reject.
func ValidPurchaseOrderAcceptRejectStatus(status PurchaseOrderStatus) bool {
	if status == StatusAccepted || status == StatusRejected {
		return true
	}
	return false
}

// ValidWhitelistAction returns true if the purchase order status is valid and false
// otherwise.
func ValidWhitelistAction(action WhitelistAction) bool {
	if action == WhitelistActionAdd ||
		action == WhitelistActionRemove {
		return true
	}
	return false
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgUndPurchaseOrder{}, "enterprise/PurchaseUnd", nil)
	cdc.RegisterConcrete(MsgProcessUndPurchaseOrder{}, "enterprise/ProcessUndPurchaseOrder", nil)
	cdc.RegisterConcrete(MsgWhitelistAddress{}, "enterprise/WhitelistAddress", nil)
}
