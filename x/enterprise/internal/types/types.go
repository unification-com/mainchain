package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	undtypes "github.com/unification-com/mainchain-cosmos/types"
	"strings"
)

type (
	PurchaseOrderStatus byte
)

const (
	DefaultDenomination                   = undtypes.DefaultDenomination
	BaseDenomination                      = undtypes.BaseDenomination
	DefaultStartingPurchaseOrderID uint64 = 1 // used in init genesis

	// Valid Purchase Order statuses
	StatusNil       PurchaseOrderStatus = 0x00
	StatusRaised    PurchaseOrderStatus = 0x01
	StatusAccepted  PurchaseOrderStatus = 0x02
	StatusRejected  PurchaseOrderStatus = 0x03
	StatusCompleted PurchaseOrderStatus = 0x04
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

// Proposals is an array of proposal
type PurchaseOrders []EnterpriseUndPurchaseOrder

// String implements stringer interface
func (p PurchaseOrders) String() string {
	out := "ID - [Purchaser] Amount (Status)\n"
	for _, po := range p {
		out += fmt.Sprintf("%d - (%s) [%s] %s\n",
			po.PurchaseOrderID, po.Amount,
			po.Purchaser, po.Status)
	}
	return strings.TrimSpace(out)
}

// EnterpriseUndPurchaseOrder is a struct that contains information on Enterprise UND purchase orders and their status
type EnterpriseUndPurchaseOrder struct {
	PurchaseOrderID uint64              `json:"id"`
	Purchaser       sdk.AccAddress      `json:"purchaser"`
	Amount          sdk.Coin            `json:"amount"`
	Status          PurchaseOrderStatus `json:"status"`
}

// NewEnterpriseUnd returns a new EnterpriseUndPurchaseOrder struct
func NewEnterpriseUnd() EnterpriseUndPurchaseOrder {
	return EnterpriseUndPurchaseOrder{
		Status: StatusNil,
	}
}

// implement fmt.Stringer
func (po EnterpriseUndPurchaseOrder) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %d
Purchaser: %s
Amount: %s
Decision: %b
`, po.PurchaseOrderID, po.Purchaser, po.Amount, po.Status))
}

// LockedUnd is a struct that is used to track "Locked" Enterprise purchased UND
type LockedUnd struct {
	Owner  sdk.AccAddress `json:"owner"`
	Amount sdk.Coin       `json:"amount"`
}

func NewLockedUnd(owner sdk.AccAddress) LockedUnd {
	return LockedUnd{
		Owner:  owner,
		Amount: sdk.NewInt64Coin(DefaultDenomination, 0),
	}
}

func (l LockedUnd) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Amount: %s
`, l.Owner, l.Amount))
}

type UndSupply struct {
	Locked   sdk.Coin `json:"locked"`
	Unlocked sdk.Coin `json:"unlocked"`
	Total    sdk.Coin `json:"total_supply"`
}

func NewUndSupply() UndSupply {
	return UndSupply{
		Locked:   sdk.NewInt64Coin(DefaultDenomination, 0),
		Unlocked: sdk.NewInt64Coin(DefaultDenomination, 0),
		Total:    sdk.NewInt64Coin(DefaultDenomination, 0),
	}
}

func (u UndSupply) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Locked: %s
Unlocked: %s
Total: %s
`, u.Locked, u.Unlocked, u.Total))
}
