package rest

//
//import (
//	"github.com/gorilla/mux"
//
//	"github.com/cosmos/cosmos-sdk/client"
//)
//
//// registerQueryRoutes - define REST query routes
//func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
//	// Deprecated
//
//	//r.HandleFunc("/beacon/params", beaconParamsHandler(cliCtx)).Methods("GET")
//	//
//	//r.HandleFunc("/beacon/beacons", beaconsWithParametersHandler(cliCtx)).Methods("GET")
//	//r.HandleFunc(fmt.Sprintf("/beacon/{%s}", RestBeaconId), beaconHandler(cliCtx)).Methods("GET")
//	//
//	//// Timestamps
//	//r.HandleFunc(fmt.Sprintf("/beacon/{%s}/timestamp/{%s}", RestBeaconId, RestBeaconTimestampId), beaconTimestampHandler(cliCtx)).Methods("GET")
//}
//
////func beaconParamsHandler(cliCtx client.Context) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		cliCtx, _ := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
////		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryParameters)
////		res, height, err := cliCtx.QueryWithData(route, nil)
////
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
////			return
////		}
////
////		cliCtx = cliCtx.WithHeight(height)
////		rest.PostProcessResponse(w, cliCtx, res)
////	}
////}
////
////func beaconsWithParametersHandler(cliCtx client.Context) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
////		if !ok {
////			return
////		}
////
////		var (
////			ownerAddr sdk.AccAddress
////			moniker   string
////		)
////
////		if v := r.URL.Query().Get(RestMoniker); len(v) != 0 {
////			moniker = v
////		}
////
////		if v := r.URL.Query().Get(RestOwnerAddr); len(v) != 0 {
////			ownerAddr, err = sdk.AccAddressFromBech32(v)
////			if err != nil {
////				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////				return
////			}
////		}
////
////		params := types.QueryBeaconsFilteredRequest{
////			Moniker: moniker,
////			Owner:   ownerAddr.String(),
////			Pagination: &sdkquery.PageRequest{
////				Limit:  uint64(limit),
////				Offset: uint64(page),
////			},
////		}
////
////		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, keeper.QueryBeacons)
////		res, height, err := cliCtx.QueryWithData(route, bz)
////
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
////			return
////		}
////
////		cliCtx = cliCtx.WithHeight(height)
////		rest.PostProcessResponse(w, cliCtx, res)
////	}
////}
////
////func beaconHandler(cliCtx client.Context) http.HandlerFunc {
////
////	return func(w http.ResponseWriter, r *http.Request) {
////		vars := mux.Vars(r)
////		strBeaconID := vars[RestBeaconId]
////
////		if len(strBeaconID) == 0 {
////			err := errors.New("beaconID required but not given")
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		beaconID, ok := rest.ParseUint64OrReturnBadRequest(w, strBeaconID)
////		if !ok {
////			return
////		}
////
////		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
////		if !ok {
////			return
////		}
////
////		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", types.ModuleName, keeper.QueryBeacon, beaconID), nil)
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
////			return
////		}
////
////		cliCtx = cliCtx.WithHeight(height)
////		rest.PostProcessResponse(w, cliCtx, res)
////	}
////}
////
////func beaconTimestampHandler(cliCtx client.Context) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		vars := mux.Vars(r)
////		strBeaconID := vars[RestBeaconId]
////		strBeaconTimestampID := vars[RestBeaconTimestampId]
////
////		if len(strBeaconID) == 0 {
////			err := errors.New("beaconID required but not given")
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////		if len(strBeaconTimestampID) == 0 {
////			err := errors.New("block height required but not given")
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		beaconID, ok := rest.ParseUint64OrReturnBadRequest(w, strBeaconID)
////		if !ok {
////			return
////		}
////		beaconTimestampID, ok := rest.ParseUint64OrReturnBadRequest(w, strBeaconTimestampID)
////		if !ok {
////			return
////		}
////		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
////		if !ok {
////			return
////		}
////
////		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d/%d", types.ModuleName, keeper.QueryBeaconTimestamp, beaconID, beaconTimestampID), nil)
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
////			return
////		}
////
////		cliCtx = cliCtx.WithHeight(height)
////		rest.PostProcessResponse(w, cliCtx, res)
////	}
////}
