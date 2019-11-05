package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type (
	PurchaseOrderStatus byte
)

// Valid Purchase Order statuses
const (

	DefaultStartingPurchaseOrderID uint64 = 1

	StatusNil       PurchaseOrderStatus = 0x00
	StatusRaised    PurchaseOrderStatus = 0x01
	StatusProcessed PurchaseOrderStatus = 0x02
	StatusRejected  PurchaseOrderStatus = 0x03
)

// EnterpriseUndPurchaseOrder is a struct that contains information on which account has purchased Enterprise UND
// and how much UND is currently locked for WRKChain Tx only use
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
Status: %b
`, po.PurchaseOrderID, po.Purchaser, po.Amount, po.Status))
}
