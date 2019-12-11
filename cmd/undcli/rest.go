package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// RegisterQueryRestApiEndpoints registers the /api REST endpoint which outputs a simple list
// of available REST endpoints
func RegisterQueryRestApiEndpoints(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/api"), queryRestApiEndpoints(r)).Methods("GET")
}

func queryRestApiEndpoints(rtr *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := rtr.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			pathTemplate, err := route.GetPathTemplate()

			outp := ""

			methods, err := route.GetMethods()
			if err == nil {
				outp = strings.Join(methods, ",") + ": "
			}

			if err == nil {
				outp = outp + pathTemplate + "\n"
			}

			_, _ = w.Write([]byte(outp))
			return nil
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}