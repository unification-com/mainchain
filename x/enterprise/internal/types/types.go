package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	undtypes "github.com/unification-com/mainchain/types"
	"strings"
)

type (
	PurchaseOrderStatus byte
	WhitelistAction     byte
)

const (
	DefaultDenomination                   = undtypes.DefaultDenomination
	DefaultStartingPurchaseOrderID uint64 = 1 // used in init genesis

	// Valid Purchase Order statuses
	StatusNil       PurchaseOrderStatus = 0x00
	StatusRaised    PurchaseOrderStatus = 0x01
	StatusAccepted  PurchaseOrderStatus = 0x02
	StatusRejected  PurchaseOrderStatus = 0x03
	StatusCompleted PurchaseOrderStatus = 0x04

	WhitelistActionAdd    WhitelistAction = 0x01
	WhitelistActionRemove WhitelistAction = 0x02
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

// ValidPurchaseOrderStatus returns true if the purchase order status is valid and false
// otherwise.
func ValidPurchaseOrderStatus(status PurchaseOrderStatus) bool {
	if status == StatusRaised ||
		status == StatusAccepted ||
		status == StatusRejected ||
		status == StatusCompleted {
		return true
	}
	return false
}

// ValidPurchaseOrderAcceptRejectStatus checks the decision - returns true if accept/reject.
func ValidPurchaseOrderAcceptRejectStatus(status PurchaseOrderStatus) bool {
	if status == StatusAccepted || status == StatusRejected {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status PurchaseOrderStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *PurchaseOrderStatus) Unmarshal(data []byte) error {
	*status = PurchaseOrderStatus(data[0])
	return nil
}

// MarshalJSON Marshals to JSON using string representation of the status
func (status PurchaseOrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// UnmarshalJSON Unmarshals from JSON assuming Bech32 encoding
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

// String implements the Stringer interface.
func (status PurchaseOrderStatus) String() string {
	switch status {
	case StatusAccepted:
		return "accept"

	case StatusRejected:
		return "reject"

	case StatusRaised:
		return "raised"

	case StatusCompleted:
		return "complete"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (status PurchaseOrderStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

// WhitelistActionFromString turns a string into a ProposalStatus
func WhitelistActionFromString(str string) (WhitelistAction, error) {
	switch str {
	case "add":
		return WhitelistActionAdd, nil

	case "remove":
		return WhitelistActionRemove, nil

	default:
		return WhitelistAction(0xff), fmt.Errorf("'%s' is not a valid whitelist action", str)
	}
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

// Marshal needed for protobuf compatibility
func (action WhitelistAction) Marshal() ([]byte, error) {
	return []byte{byte(action)}, nil
}

// Unmarshal needed for protobuf compatibility
func (action *WhitelistAction) Unmarshal(data []byte) error {
	*action = WhitelistAction(data[0])
	return nil
}

// MarshalJSON Marshals to JSON using string representation of the status
func (action WhitelistAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(action.String())
}

// UnmarshalJSON Unmarshals from JSON assuming Bech32 encoding
func (action *WhitelistAction) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := WhitelistActionFromString(s)
	if err != nil {
		return err
	}

	*action = bz2
	return nil
}

// String implements the Stringer interface.
func (action WhitelistAction) String() string {
	switch action {
	case WhitelistActionAdd:
		return "add"

	case WhitelistActionRemove:
		return "remove"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (action WhitelistAction) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(action.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(action))))
	}
}

// PurchaseOrders is an array of purchase orders
type PurchaseOrders []EnterpriseUndPurchaseOrder

// String implements stringer interface
func (p PurchaseOrders) String() string {
	out := "ID - [Purchaser] Amount (Status) {Raised Time} <Decision Time>\n"
	for _, po := range p {
		out += fmt.Sprintf("%d - [%s] %s (%s) {%d} <%v>\n",
			po.PurchaseOrderID, po.Amount,
			po.Purchaser, po.Status, po.RaisedTime, po.Decisions)
	}
	return strings.TrimSpace(out)
}

type PurchaseOrderDecision struct {
	Signer       sdk.AccAddress      `json:"signer"`
	Decision     PurchaseOrderStatus `json:"decision"`
	DecisionTime int64               `json:"decision_time"`
}

func NewPurchaseOrderDecision(signer sdk.AccAddress, decision PurchaseOrderStatus) PurchaseOrderDecision {
	return PurchaseOrderDecision{
		Signer:   signer,
		Decision: decision,
	}
}

// EnterpriseUndPurchaseOrder is a struct that contains information on Enterprise UND purchase orders and their status
type EnterpriseUndPurchaseOrder struct {
	PurchaseOrderID uint64                  `json:"id"`
	Purchaser       sdk.AccAddress          `json:"purchaser"`
	Amount          sdk.Coin                `json:"amount"`
	Status          PurchaseOrderStatus     `json:"status"`
	RaisedTime      int64                   `json:"raise_time"`
	Decisions       []PurchaseOrderDecision `json:"decisions"`
	CompletionTime  int64                   `json:"completion_time"`
}

// NewEnterpriseUndPurchaseOrder returns a new EnterpriseUndPurchaseOrder struct
func NewEnterpriseUndPurchaseOrder() EnterpriseUndPurchaseOrder {
	return EnterpriseUndPurchaseOrder{
		Status:         StatusNil,
		RaisedTime:     0,
		CompletionTime: 0,
	}
}

// implement fmt.Stringer
func (po EnterpriseUndPurchaseOrder) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %d
Purchaser: %s
Amount: %s
RaisedTime: %d
Decisions: %v
Status: %b
`, po.PurchaseOrderID, po.Purchaser, po.Amount, po.RaisedTime, po.Decisions, po.Status))
}

// LockedUnds is an array of locked UND
type LockedUnds []LockedUnd

// String implements stringer interface
func (lund LockedUnds) String() string {
	out := "Purchaser [Amount]\n"
	for _, l := range lund {
		out += fmt.Sprintf("%s [%s]\n",
			l.Owner, l.Amount)
	}
	return strings.TrimSpace(out)
}

// LockedUnd is a struct that is used to track "Locked" Enterprise purchased UND
type LockedUnd struct {
	Owner  sdk.AccAddress `json:"owner"`
	Amount sdk.Coin       `json:"amount"`
}

func NewLockedUnd(owner sdk.AccAddress, denom string) LockedUnd {
	return LockedUnd{
		Owner:  owner,
		Amount: sdk.NewInt64Coin(denom, 0),
	}
}

func (l LockedUnd) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Amount: %s
`, l.Owner, l.Amount))
}

type UndSupplies []UndSupply

type UndSupply struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
	Locked int64  `json:"locked"`
	Total  int64  `json:"total"`
}

func NewUndSupply(denom string) UndSupply {
	return UndSupply{
		Denom:  denom,
		Amount: 0, // current unlocked, liquid UND
		Locked: 0,
		Total:  0,
	}
}

func (u UndSupply) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Denom: %s
Amount: %d
Locked: %d
Total: %d
`, u.Denom, u.Amount, u.Locked, u.Total))
}

type WhitelistAddresses []sdk.AccAddress
