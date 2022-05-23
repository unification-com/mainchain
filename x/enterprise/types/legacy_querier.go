package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NOTE: Will be deprecated. Used for legacy REST queries

// QueryResPurchaseOrders Queries raised Enterprise FUND purchase orders
type QueryResPurchaseOrders []EnterpriseUndPurchaseOrder

// QueryPurchaseOrdersParams Params for query 'custom/enterprise/pos'
type QueryPurchaseOrdersParams struct {
	Page                int
	Limit               int
	PurchaseOrderStatus PurchaseOrderStatus
	Purchaser           sdk.AccAddress
}

// NewQueryPurchaseOrdersParams creates a new instance of QueryPurchaseOrdersParams
func NewQueryPurchaseOrdersParams(page, limit int, status PurchaseOrderStatus, purchaser sdk.AccAddress) QueryPurchaseOrdersParams {
	return QueryPurchaseOrdersParams{
		Page:                page,
		Limit:               limit,
		PurchaseOrderStatus: status,
		Purchaser:           purchaser,
	}
}
