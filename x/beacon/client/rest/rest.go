package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
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

// RegisterLegacyRESRRoutes - Central function to define routes that get registered by the main application
func RegisterLegacyRESRRoutes(clientCtx client.Context, rtr *mux.Router) {
	//r := rest.WithHTTPDeprecationHeaders(rtr)
	//
	//registerQueryRoutes(clientCtx, r)
	//registerTxRoutes(clientCtx, r)
}
