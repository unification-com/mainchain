package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

const (
	restWrkchainId  = "wrkchainid"
	restBlockHeight = "height"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/wrkchain", storeName), registerWrkChainHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/wrkchain/{%s}/get", storeName, restWrkchainId), wrkChainHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/wrkchain", storeName), recordWrkChainBlockHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/wrkchain/{%s}/{%s}/get-block", storeName, restWrkchainId, restBlockHeight), wrkChainBlockHandler(cliCtx, storeName)).Methods("GET")
}
