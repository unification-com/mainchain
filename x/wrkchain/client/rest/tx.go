package rest

//
//import (
//	"github.com/gorilla/mux"
//
//	"github.com/cosmos/cosmos-sdk/client"
//	"github.com/cosmos/cosmos-sdk/types/rest"
//)
//
//type registerWrkChainReq struct {
//	BaseReq      rest.BaseReq `json:"base_req"`
//	Moniker      string       `json:"moniker"`
//	WrkChainName string       `json:"name"`
//	GenesisHash  string       `json:"genesis"`
//	BaseType     string       `json:"base"`
//	Owner        string       `json:"owner"`
//}
//
//type recordWrkChainBlockReq struct {
//	BaseReq    rest.BaseReq `json:"base_req"`
//	WrkChainID uint64       `json:"id"`
//	Height     uint64       `json:"height"`
//	BlockHash  string       `json:"blockhash"`
//	ParentHash string       `json:"parenthash"`
//	Hash1      string       `json:"hash1"`
//	Hash2      string       `json:"hash2"`
//	Hash3      string       `json:"hash3"`
//	Owner      string       `json:"owner"`
//}
//
//// registerTxRoutes - define REST Tx routes
//func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
//	//r.HandleFunc("/wrkchain/reg", registerWrkChainHandler(cliCtx)).Methods("POST")
//	//
//	//r.HandleFunc("/wrkchain/rec", recordWrkChainBlockHandler(cliCtx)).Methods("POST")
//}
//
////func registerWrkChainHandler(cliCtx client.Context) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		var req registerWrkChainReq
////		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
////			return
////		}
////
////		baseReq := req.BaseReq.Sanitize()
////
////		// Todo -  automatically apply fees
////		//paramsRetriever := keeper.NewParamsRetriever(cliCtx)
////		//wrkchainParams, err := paramsRetriever.GetParams()
////		//if err != nil {
////		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////		//}
////		//
////		//fees := sdk.NewCoins(sdk.NewInt64Coin(wrkchainParams.Denom, int64(wrkchainParams.FeeRegister)))
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
////		msg := types.NewMsgRegisterWrkChain(req.Moniker, req.WrkChainName, req.GenesisHash, req.BaseType, addr)
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
////func recordWrkChainBlockHandler(cliCtx client.Context) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		var req recordWrkChainBlockReq
////		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
////			return
////		}
////
////		baseReq := req.BaseReq.Sanitize()
////
////		// todo - automatically apply fees
////		//paramsRetriever := keeper.NewParamsRetriever(cliCtx)
////		//wrkchainParams, err := paramsRetriever.GetParams()
////		//if err != nil {
////		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////		//}
////		//
////		//fees := sdk.NewCoins(sdk.NewInt64Coin(wrkchainParams.Denom, int64(wrkchainParams.FeeRecord)))
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
////		msg := types.NewMsgRecordWrkChainBlock(req.WrkChainID, req.Height, req.BlockHash, req.ParentHash, req.Hash1, req.Hash2, req.Hash3, addr)
////		err = msg.ValidateBasic()
////		if err != nil {
////			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
////			return
////		}
////
////		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
////	}
////}
