package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	undtypes "github.com/unification-com/mainchain/types"
	"github.com/unification-com/mainchain/x/enterprise/internal/keeper"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// registerQueryRoutes - define REST query routes
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/enterprise/params"), enterpriseParamsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/enterprise/locked"), enterpriseTotalLockedHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/enterprise/unlocked"), enterpriseTotalUnLockedHandler(cliCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/enterprise/pos"), enterprisePosWithParametersHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/enterprise/po/{%s}", RestPurchaseOrderId), enterprisePurchaseOrderHandler(cliCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/enterprise/{%s}/locked", RestPurchaserAddr), enterpriseLockedForAddressHandler(cliCtx)).Methods("GET")
}

func registerEnterpriseAuthAccountOverride(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/auth/accounts/{address}"), EnterpriseAuthAccountOverride(cliCtx)).Methods("GET")
}

func registerEnterpriseTotalSupplyOverride(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/supply/total"), EnterpriseSupplyTotalOverride(cliCtx)).Methods("GET")
}

func registerEnterpriseSupplyByDenomOverride(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/supply/total/{denom}"), EnterpriseSupplyByDenomOverride(cliCtx)).Methods("GET")
}

func enterpriseParamsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func enterpriseTotalLockedHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func enterpriseTotalUnLockedHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func enterprisePosWithParametersHandler(cliCtx context.CLIContext) http.HandlerFunc {
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

		bz, err := cliCtx.Codec.MarshalJSON(params)
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

func enterprisePurchaseOrderHandler(cliCtx context.CLIContext) http.HandlerFunc {

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

func enterpriseLockedForAddressHandler(cliCtx context.CLIContext) http.HandlerFunc {
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

func EnterpriseAuthAccountOverride(cliCtx context.CLIContext) http.HandlerFunc {
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

		accGetter := auth.NewAccountRetriever(cliCtx)
		lockedUndGetter := keeper.NewLockedUndRetriever(cliCtx)

		account, height, err := accGetter.GetAccountWithHeight(addr)
		if err != nil {
			// TODO: Handle more appropriately based on the error type.
			// Ref: https://github.com/cosmos/cosmos-sdk/issues/4923
			if err := accGetter.EnsureExists(addr); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				rest.PostProcessResponse(w, cliCtx, auth.BaseAccount{})
				return
			}

			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		lockedUnd, err := lockedUndGetter.GetLockedUnd(addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// todo - this is a bit hackey
		accountWithLocked := undtypes.NewAccountWithLocked()
		entUnd := undtypes.NewEnterpriseUnd()

		entUnd.Locked = lockedUnd.Amount
		entUnd.Available = account.GetCoins().Add(lockedUnd.Amount)

		accountWithLocked.Account = account
		accountWithLocked.Enterprise = entUnd

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, accountWithLocked)
	}
}

func EnterpriseSupplyTotalOverride(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		totalSupplyGetter := keeper.NewTotalSupplyRetriever(cliCtx)

		totalSupply, height, err := totalSupplyGetter.GetTotalSupplyHeight()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, totalSupply)
	}
}

func EnterpriseSupplyByDenomOverride(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		denom := mux.Vars(r)["denom"]

		// todo - get from params
		if denom != types.DefaultDenomination {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "unknown denomination "+denom)
			return
		}

		totalSupplyGetter := keeper.NewTotalSupplyRetriever(cliCtx)

		totalSupply, height, err := totalSupplyGetter.GetTotalSupplyHeight()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, totalSupply.Amount)
	}
}
