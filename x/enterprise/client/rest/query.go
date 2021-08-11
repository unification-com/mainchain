package rest

import (
	"errors"
	"fmt"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// registerQueryRoutes - define REST query routes
func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/enterprise/params", enterpriseParamsHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/enterprise/locked", enterpriseTotalLockedHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/enterprise/unlocked", enterpriseTotalUnLockedHandler(clientCtx)).Methods("GET")

	r.HandleFunc("/enterprise/whitelist", enterpriseWhitelistHandler(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/enterprise/whitelisted/{%s}", RestWhitelistAddr), enterpriseWhitelistedHandler(clientCtx)).Methods("GET")

	r.HandleFunc("/enterprise/pos", enterprisePosWithParametersHandler(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/enterprise/po/{%s}", RestPurchaseOrderId), enterprisePurchaseOrderHandler(clientCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/enterprise/{%s}/locked", RestPurchaserAddr), enterpriseLockedForAddressHandler(clientCtx)).Methods("GET")
}

func registerEnterpriseAuthAccountOverride(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/auth/accounts/{address}", EnterpriseAuthAccountOverride(cliCtx)).Methods("GET")
}

func registerEnterpriseTotalSupplyOverride(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/supply/total", EnterpriseSupplyTotalOverride(cliCtx)).Methods("GET")
}

func registerEnterpriseSupplyByDenomOverride(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/supply/total/{denom}", EnterpriseSupplyByDenomOverride(cliCtx)).Methods("GET")
}

func enterpriseParamsHandler(cliCtx client.Context) http.HandlerFunc {
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

func enterpriseTotalLockedHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryTotalLocked)
		res, height, err := cliCtx.QueryWithData(route, nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func enterpriseTotalUnLockedHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryTotalUnlocked)
		res, height, err := cliCtx.QueryWithData(route, nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func enterprisePosWithParametersHandler(cliCtx client.Context) http.HandlerFunc {
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
			purchaserAddr sdk.AccAddress
			status        types.PurchaseOrderStatus
		)

		if v := r.URL.Query().Get(RestStatus); len(v) != 0 {
			status, err = types.PurchaseOrderStatusFromString(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if v := r.URL.Query().Get(RestPurchaserAddr); len(v) != 0 {
			purchaserAddr, err = sdk.AccAddressFromBech32(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.NewQueryPurchaseOrdersParams(page, limit, status, purchaserAddr)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, keeper.QueryPurchaseOrders)
		res, height, err := cliCtx.QueryWithData(route, bz)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func enterprisePurchaseOrderHandler(cliCtx client.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strPurchaseOrderID := vars[RestPurchaseOrderId]

		if len(strPurchaseOrderID) == 0 {
			err := errors.New("purchaseOrderID required but not given")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		poID, ok := rest.ParseUint64OrReturnBadRequest(w, strPurchaseOrderID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", types.ModuleName, keeper.QueryGetPurchaseOrder, poID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func enterpriseLockedForAddressHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strPurchaserAddr := vars[RestPurchaserAddr]

		if len(strPurchaserAddr) == 0 {
			err := errors.New("purchaserAddr required but not given")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		purchaserAddr, err := sdk.AccAddressFromBech32(strPurchaserAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, keeper.QueryGetLocked, purchaserAddr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func enterpriseWhitelistHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryWhitelist)
		res, height, err := cliCtx.QueryWithData(route, nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func enterpriseWhitelistedHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		strAddr := vars[RestWhitelistAddr]
		if len(strAddr) == 0 {
			err := errors.New("address required but not given")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		addr, err := sdk.AccAddressFromBech32(strAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, keeper.QueryWhitelisted, addr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func EnterpriseAuthAccountOverride(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := auth.AccountRetriever{}

		_, height, err := accGetter.GetAccountWithHeight(cliCtx, addr)
		lockedUndGetter := types.LockedRetriever{}

		if err != nil {
			// TODO: Handle more appropriately based on the error type.
			// Ref: https://github.com/cosmos/cosmos-sdk/issues/4923
			if err := accGetter.EnsureExists(cliCtx, addr); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				rest.PostProcessResponse(w, cliCtx, auth.BaseAccount{})
				return
			}

			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		lockedUnd, _, err := lockedUndGetter.GetLockedWithHeight(cliCtx, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// todo - this is a bit hackey
		//accountWithLocked := undtypes.NewAccountWithLocked()
		//entUnd := undtypes.NewEnterpriseUnd()
		//
		//entUnd.Locked = lockedUnd.Amount
		////entUnd.Available = account.GetCoins().Add(lockedUnd.Amount)
		//
		//accountWithLocked.Account = account
		//accountWithLocked.Enterprise = entUnd

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, lockedUnd)
	}
}

func EnterpriseSupplyTotalOverride(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		//totalSupplyGetter := keeper.NewTotalSupplyRetriever(cliCtx)
		//
		//totalSupply, height, err := totalSupplyGetter.GetTotalSupplyHeight()
		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		//	return
		//}
		//cliCtx = cliCtx.WithHeight(height)
		//rest.PostProcessResponse(w, cliCtx, totalSupply)
		rest.PostProcessResponse(w, cliCtx, nil)
	}
}

func EnterpriseSupplyByDenomOverride(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		//denom := mux.Vars(r)["denom"]
		//
		//// todo - get from params
		//if denom != types.DefaultDenomination {
		//	rest.WriteErrorResponse(w, http.StatusInternalServerError, "unknown denomination "+denom)
		//	return
		//}
		//
		//totalSupplyGetter := keeper.NewTotalSupplyRetriever(cliCtx)
		//
		//totalSupply, height, err := totalSupplyGetter.GetTotalSupplyHeight()
		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		//	return
		//}
		//
		//cliCtx = cliCtx.WithHeight(height)
		//rest.PostProcessResponse(w, cliCtx, totalSupply.Amount)
		rest.PostProcessResponse(w, cliCtx, nil)
	}
}
