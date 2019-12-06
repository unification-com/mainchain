package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// nolint
const (
	RestBeaconId          = "beaconid"
	RestBeaconTimestampId = "timestampid"
	RestMoniker           = "moniker"
	RestOwnerAddr         = "owner"
	RestSubmitTime        = "subtime"
	RestHash              = "hash"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
