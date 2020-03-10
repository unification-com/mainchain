package rest

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/unification-com/mainchain/x/wrkchain/internal/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
	"net/http"
)

// registerQueryRoutes - define REST query routes
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/wrkchain/params"), wrkChainParamsHandler(cliCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/wrkchain/wrkchains"), wrkChainsWithParametersHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/wrkchain/{%s}", RestWrkchainId), wrkChainHandler(cliCtx)).Methods("GET")

	// Block hashes
	r.HandleFunc(fmt.Sprintf("/wrkchain/{%s}/block/{%s}", RestWrkchainId, RestBlockHeight), wrkChainBlockHandler(cliCtx)).Methods("GET")
}

func wrkChainParamsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryParameters)
		res, height, err := cliCtx.QueryWithData(route, nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func wrkChainsWithParametersHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var (
			ownerAddr sdk.AccAddress
			moniker   string
		)

		if v := r.URL.Query().Get(RestMoniker); len(v) != 0 {
			moniker = v
		}

		if v := r.URL.Query().Get(RestOwnerAddr); len(v) != 0 {
			ownerAddr, err = sdk.AccAddressFromBech32(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.NewQueryWrkChainParams(page, limit, moniker, ownerAddr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, keeper.QueryWrkChainsFiltered)
		res, height, err := cliCtx.QueryWithData(route, bz)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func wrkChainHandler(cliCtx context.CLIContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strWrkchainID := vars[RestWrkchainId]

		if len(strWrkchainID) == 0 {
			err := errors.New("wrkchainID required but not given")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		wrkchainID, ok := rest.ParseUint64OrReturnBadRequest(w, strWrkchainID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", types.ModuleName, keeper.QueryWrkChain, wrkchainID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func wrkChainBlockHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strWrkchainID := vars[RestWrkchainId]
		strBlockHeight := vars[RestBlockHeight]

		if len(strWrkchainID) == 0 {
			err := errors.New("wrkchainID required but not given")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(strBlockHeight) == 0 {
			err := errors.New("block height required but not given")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		wrkchainID, ok := rest.ParseUint64OrReturnBadRequest(w, strWrkchainID)
		if !ok {
			return
		}
		wrkchainBlockHeight, ok := rest.ParseUint64OrReturnBadRequest(w, strBlockHeight)
		if !ok {
			return
		}
		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d/%d", types.ModuleName, keeper.QueryWrkChainBlock, wrkchainID, wrkchainBlockHeight), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
