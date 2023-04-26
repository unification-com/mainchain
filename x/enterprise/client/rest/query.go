package rest

//
//import (
//	"errors"
//	"fmt"
//	"net/http"
//
//	"github.com/gorilla/mux"
//
//	"github.com/cosmos/cosmos-sdk/client"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/types/rest"
//	"github.com/unification-com/mainchain/x/enterprise/keeper"
//	"github.com/unification-com/mainchain/x/enterprise/types"
//)
//
//// registerQueryRoutes - define REST query routes
//func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
//	//r.HandleFunc("/enterprise/params", enterpriseParamsHandler(clientCtx)).Methods("GET")
//	//r.HandleFunc("/enterprise/locked", enterpriseTotalLockedHandler(clientCtx)).Methods("GET")
//	//r.HandleFunc("/enterprise/unlocked", enterpriseTotalUnLockedHandler(clientCtx)).Methods("GET")
//	//r.HandleFunc("/enterprise/ent_supply", enterpriseEnterpriseSupplyHandler(clientCtx)).Methods("GET")
//	//
//	//r.HandleFunc("/enterprise/whitelist", enterpriseWhitelistHandler(clientCtx)).Methods("GET")
//	//r.HandleFunc(fmt.Sprintf("/enterprise/whitelisted/{%s}", RestWhitelistAddr), enterpriseWhitelistedHandler(clientCtx)).Methods("GET")
//	//
//	//r.HandleFunc("/enterprise/pos", enterprisePosWithParametersHandler(clientCtx)).Methods("GET")
//	//r.HandleFunc(fmt.Sprintf("/enterprise/po/{%s}", RestPurchaseOrderId), enterprisePurchaseOrderHandler(clientCtx)).Methods("GET")
//	//
//	//r.HandleFunc(fmt.Sprintf("/enterprise/{%s}/locked", RestPurchaserAddr), enterpriseLockedForAddressHandler(clientCtx)).Methods("GET")
//}
//
//func registerEnterpriseTotalSupplyOverride(cliCtx client.Context, r *mux.Router) {
//	r.HandleFunc("/supply/total", EnterpriseSupplyTotalOverride(cliCtx)).Methods("GET")
//}
//
//func registerEnterpriseSupplyByDenomOverride(cliCtx client.Context, r *mux.Router) {
//	r.HandleFunc("/supply/total/{denom}", EnterpriseSupplyByDenomOverride(cliCtx)).Methods("GET")
//}
//
//func enterpriseParamsHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryParameters)
//		res, height, err := cliCtx.QueryWithData(route, nil)
//
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterpriseTotalLockedHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryTotalLocked)
//		res, height, err := cliCtx.QueryWithData(route, nil)
//
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterpriseTotalUnLockedHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryTotalUnlocked)
//		res, height, err := cliCtx.QueryWithData(route, nil)
//
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterpriseEnterpriseSupplyHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryEnterpriseSupply)
//		res, height, err := cliCtx.QueryWithData(route, nil)
//
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterprisePosWithParametersHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		if !ok {
//			return
//		}
//
//		var (
//			purchaserAddr sdk.AccAddress
//			status        types.PurchaseOrderStatus
//		)
//
//		if v := r.URL.Query().Get(RestStatus); len(v) != 0 {
//			status, err = types.PurchaseOrderStatusFromString(v)
//			if err != nil {
//				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//				return
//			}
//		}
//
//		if v := r.URL.Query().Get(RestPurchaserAddr); len(v) != 0 {
//			purchaserAddr, err = sdk.AccAddressFromBech32(v)
//			if err != nil {
//				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//				return
//			}
//		}
//
//		params := types.NewQueryPurchaseOrdersParams(page, limit, status, purchaserAddr)
//
//		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, keeper.QueryPurchaseOrders)
//		res, height, err := cliCtx.QueryWithData(route, bz)
//
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterprisePurchaseOrderHandler(cliCtx client.Context) http.HandlerFunc {
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r)
//		strPurchaseOrderID := vars[RestPurchaseOrderId]
//
//		if len(strPurchaseOrderID) == 0 {
//			err := errors.New("purchaseOrderID required but not given")
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		poID, ok := rest.ParseUint64OrReturnBadRequest(w, strPurchaseOrderID)
//		if !ok {
//			return
//		}
//
//		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		if !ok {
//			return
//		}
//
//		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", types.ModuleName, keeper.QueryGetPurchaseOrder, poID), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterpriseLockedForAddressHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r)
//		strPurchaserAddr := vars[RestPurchaserAddr]
//
//		if len(strPurchaserAddr) == 0 {
//			err := errors.New("purchaserAddr required but not given")
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		purchaserAddr, err := sdk.AccAddressFromBech32(strPurchaserAddr)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		if !ok {
//			return
//		}
//
//		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, keeper.QueryGetLocked, purchaserAddr), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterpriseWhitelistHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryWhitelist)
//		res, height, err := cliCtx.QueryWithData(route, nil)
//
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func enterpriseWhitelistedHandler(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		vars := mux.Vars(r)
//		strAddr := vars[RestWhitelistAddr]
//		if len(strAddr) == 0 {
//			err := errors.New("address required but not given")
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//		addr, err := sdk.AccAddressFromBech32(strAddr)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		if !ok {
//			return
//		}
//
//		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, keeper.QueryWhitelisted, addr), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func EnterpriseSupplyTotalOverride(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.ModuleName, keeper.QueryTotalSupply), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func EnterpriseSupplyByDenomOverride(cliCtx client.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		denom := mux.Vars(r)["denom"]
//
//		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
//
//		if !ok {
//			return
//		}
//
//		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, keeper.QueryTotalSupplyOf, denom), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//
//		cliCtx = cliCtx.WithHeight(height)
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
