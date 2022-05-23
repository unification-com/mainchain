package rest

import (
	"errors"
	"fmt"
	sdkquery "github.com/cosmos/cosmos-sdk/types/query"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/unification-com/mainchain/x/wrkchain/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// registerQueryRoutes - define REST query routes
func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/wrkchain/params", wrkChainParamsHandler(cliCtx)).Methods("GET")

	r.HandleFunc("/wrkchain/wrkchains", wrkChainsWithParametersHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/wrkchain/{%s}", RestWrkchainId), wrkChainHandler(cliCtx)).Methods("GET")

	// Block hashes
	r.HandleFunc(fmt.Sprintf("/wrkchain/{%s}/block/{%s}", RestWrkchainId, RestBlockHeight), wrkChainBlockHandler(cliCtx)).Methods("GET")
}

func wrkChainParamsHandler(cliCtx client.Context) http.HandlerFunc {
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

func wrkChainsWithParametersHandler(cliCtx client.Context) http.HandlerFunc {
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

		params := types.QueryWrkChainsFilteredRequest{
			Moniker: moniker,
			Owner:   ownerAddr.String(),
			Pagination: &sdkquery.PageRequest{
				Limit:  uint64(limit),
				Offset: uint64(page),
			},
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
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

func wrkChainHandler(cliCtx client.Context) http.HandlerFunc {

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

func wrkChainBlockHandler(cliCtx client.Context) http.HandlerFunc {
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
