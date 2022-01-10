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

	// legacy REST overrides
	registerEnterpriseTotalSupplyOverride(clientCtx, r)
	registerEnterpriseSupplyByDenomOverride(clientCtx, r)
}
