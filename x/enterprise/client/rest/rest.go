package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
)

// nolint
const (
	RestPurchaseOrderId = "poid"
	RestPurchaserAddr   = "purchaser"
	RestStatus          = "status"
	RestWhitelistAddr   = "address"
)

// RegisterLegacyRESTRoutes - Central function to define routes that get registered by the main application
func RegisterLegacyRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)

	registerQueryRoutes(clientCtx, r)
	registerTxRoutes(clientCtx, r)

	// legacy overrides
	registerEnterpriseAuthAccountOverride(clientCtx, r)
	registerEnterpriseTotalSupplyOverride(clientCtx, r)
	registerEnterpriseSupplyByDenomOverride(clientCtx, r)
}

// RegisterAuthAccountOverride registers a REST route to override the default /auth/accounts/{address} path
// and additionally return Enterprise Locked FUND data
//func RegisterAuthAccountOverride(cliCtx client.Context, r *mux.Router) {
//	registerEnterpriseAuthAccountOverride(cliCtx, r)
//}
//
//func RegisterTotalSupplyOverride(cliCtx client.Context, r *mux.Router) {
//	registerEnterpriseTotalSupplyOverride(cliCtx, r)
//	registerEnterpriseSupplyByDenomOverride(cliCtx, r)
//}
