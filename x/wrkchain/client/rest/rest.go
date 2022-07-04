package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

// nolint
const (
	RestWrkchainId  = "wrkchainid"
	RestBlockHeight = "height"
	RestMoniker     = "moniker"
	RestOwnerAddr   = "owner"
	RestMinHeight   = "min"
	RestMaxHeight   = "max"
	RestMinDate     = "after"
	RestMaxDate     = "before"
)

// RegisterLegacyRESTRoutes - Central function to define routes that get registered by the main application
func RegisterLegacyRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	//r := rest.WithHTTPDeprecationHeaders(rtr)
	//
	//registerQueryRoutes(clientCtx, r)
	//registerTxRoutes(clientCtx, r)
}
