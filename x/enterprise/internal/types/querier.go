package types

import "strings"

// QueryResRaisedPurchaseOrders Queries raised Enterprise UND purchase orders
type QueryResRaisedPurchaseOrders []EnterpriseUndPurchaseOrder

// implement fmt.Stringer
func (po QueryResRaisedPurchaseOrders) String() (out string) {
	for _, val := range po {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
