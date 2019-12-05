package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
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

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
