package rest

//
//import (
//	"github.com/gorilla/mux"
//
//	"github.com/cosmos/cosmos-sdk/client"
//	"github.com/cosmos/cosmos-sdk/types/rest"
//)
//
//type registerBeaconReq struct {
//	BaseReq    rest.BaseReq `json:"base_req"`
//	Moniker    string       `json:"moniker"`
//	BeaconName string       `json:"name"`
//	Owner      string       `json:"owner"`
//}
//
//type recordBeaconTimestampReq struct {
//	BaseReq    rest.BaseReq `json:"base_req"`
//	BeaconID   uint64       `json:"id"`
//	SubmitTime uint64       `json:"subtime"`
//	Hash       string       `json:"hash"`
//	Owner      string       `json:"owner"`
//}
//
//// registerTxRoutes - define REST Tx routes
//func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
//	//r.HandleFunc("/beacon/reg", registerBeaconHandler(cliCtx)).Methods("POST")
//	//
//	//r.HandleFunc("/beacon/rec", recordBeaconTimestampHandler(cliCtx)).Methods("POST")
//}
//
////func registerBeaconHandler(cliCtx client.Context) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		var req registerBeaconReq
////		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
////			return
////		}
////
////		baseReq := req.BaseReq.Sanitize()
////
////		// Todo -  automatically apply fees
////		//paramsRetriever := keeper.NewParamsRetriever(cliCtx)
////		//beaconParams, err := paramsRetriever.GetParams()
////		//if err != nil {
////		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////		//}
////		//
////		//fees := sdk.NewCoins(sdk.NewInt64Coin(beaconParams.Denom, int64(beaconParams.FeeRegister)))
////		//
////		//baseReq.Fees = fees
////
////		if !baseReq.ValidateBasic(w) {
////			return
////		}
////
////		addr, err := sdk.AccAddressFromBech32(req.Owner)
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		// create the message
////		msg := types.NewMsgRegisterBeacon(req.Moniker, req.BeaconName, addr)
////		err = msg.ValidateBasic()
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
////	}
////}
////
////func recordBeaconTimestampHandler(cliCtx client.Context) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		var req recordBeaconTimestampReq
////		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
////			return
////		}
////
////		baseReq := req.BaseReq.Sanitize()
////
////		// Todo - automatically apply fees
////		//paramsRetriever := keeper.NewParamsRetriever(cliCtx)
////		//beaconParams, err := paramsRetriever.GetParams()
////		//if err != nil {
////		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////		//}
////		//
////		//fees := sdk.NewCoins(sdk.NewInt64Coin(beaconParams.Denom, int64(beaconParams.FeeRecord)))
////		//
////		//baseReq.Fees = fees
////
////		if !baseReq.ValidateBasic(w) {
////			return
////		}
////
////		addr, err := sdk.AccAddressFromBech32(req.Owner)
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		// create the message
////		msg := types.NewMsgRecordBeaconTimestamp(req.BeaconID, req.Hash, req.SubmitTime, addr)
////		err = msg.ValidateBasic()
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
////	}
////}
