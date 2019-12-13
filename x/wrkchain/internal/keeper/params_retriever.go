package keeper

import (
	"fmt"

	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

// NodeQuerier is an interface that is satisfied by types that provide the QueryWithData method
type NodeQuerier interface {
	// QueryWithData performs a query to a Tendermint node with the provided path
	// and a data payload. It returns the result and height of the query upon success
	// or an error if the query fails.
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

// ParamsRetriever defines the properties of a type that can be used to
// retrieve enterprise params.
type ParamsRetriever struct {
	querier NodeQuerier
}

// NewParamsRetriever initialises a new ParamsRetriever instance.
func NewParamsRetriever(querier NodeQuerier) ParamsRetriever {
	return ParamsRetriever{querier: querier}
}

// GetParams queries for parameters. An
// error is returned if the query or decoding fails.
func (ar ParamsRetriever) GetParams() (types.Params, error) {
	params, _, err := ar.GetParamsHeight()
	return params, err
}

// GetParamsHeight queries for parameters. Returns the
// height of the query with the params. An error is returned if the query
// or decoding fails.
func (ar ParamsRetriever) GetParamsHeight() (types.Params, int64, error) {

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, QueryParameters), nil)
	if err != nil {
		return types.NewParams(types.RegFee, types.RecordFee, types.FeeDenom), 0, err
	}

	var params types.Params
	if err := types.ModuleCdc.UnmarshalJSON(res, &params); err != nil {
		return types.NewParams(types.RegFee, types.RecordFee, types.FeeDenom), 0, err
	}

	return params, height, nil
}
