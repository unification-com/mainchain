package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

type raisePurchaseOrderReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Amount    sdk.Coin     `json:"amount"`
	Purchaser string       `json:"purchaser"`
}

type processPurchaseOrderReq struct {
	BaseReq         rest.BaseReq `json:"base_req"`
	PurchaseOrderID uint64       `json:"poid"`
	Decision        string       `json:"decision"`
	Signer          string       `json:"signer"`
}

type processWhitelistActionReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Address string       `json:"address"`
	Action  string       `json:"action"`
	Signer  string       `json:"signer"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/enterprise/purchase", raisePurchaseOrderHandler(cliCtx)).Methods("POST")

	r.HandleFunc("/enterprise/process", processPurchaseOrderHandler(cliCtx)).Methods("POST")

	r.HandleFunc("/enterprise/whitelist", processWhitelistActionHandler(cliCtx)).Methods("POST")
}

func raisePurchaseOrderHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req raisePurchaseOrderReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Purchaser)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgUndPurchaseOrder(addr, req.Amount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func processPurchaseOrderHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req processPurchaseOrderReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		decision, err := types.PurchaseOrderStatusFromString(req.Decision)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !types.ValidPurchaseOrderAcceptRejectStatus(decision) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "decision should be accept or reject")
			return
		}

		// create the message
		msg := types.NewMsgProcessUndPurchaseOrder(req.PurchaseOrderID, decision, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func processWhitelistActionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req processWhitelistActionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		action, err := types.WhitelistActionFromString(req.Action)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !types.ValidWhitelistAction(action) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "action should be add or remove")
			return
		}

		// create the message
		msg := types.NewMsgWhitelistAddress(addr, action, signer)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
