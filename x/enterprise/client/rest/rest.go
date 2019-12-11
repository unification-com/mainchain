package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// nolint
const (
	RestPurchaseOrderId = "poid"
	RestPurchaserAddr   = "purchaser"
	RestStatus          = "status"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// RegisterAuthAccountOverride registers a REST route to override the default /auth/accounts/{address} path
// and additionally return Enterprise Locked UND data
func RegisterAuthAccountOverride(cliCtx context.CLIContext, r *mux.Router) {
	registerEnterpriseAuthAccountOverride(cliCtx, r)
}

func RegisterTotalSupplyOverride(cliCtx context.CLIContext, r *mux.Router) {
	registerEnterpriseTotalSupplyOverride(cliCtx, r)
	registerEnterpriseSupplyByDenomOverride(cliCtx, r)
}