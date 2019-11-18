package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	undtypes "github.com/unification-com/mainchain-cosmos/types"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

// NodeQuerier is an interface that is satisfied by types that provide the QueryWithData method
type NodeQuerier interface {
	// QueryWithData performs a query to a Tendermint node with the provided path
	// and a data payload. It returns the result and height of the query upon success
	// or an error if the query fails.
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

// LockedUndRetriever defines the properties of a type that can be used to
// retrieve locked UND.
type LockedUndRetriever struct {
	querier NodeQuerier
}

// NewLockedUndRetriever initialises a new LockedUndRetriever instance.
func NewLockedUndRetriever(querier NodeQuerier) LockedUndRetriever {
	return LockedUndRetriever{querier: querier}
}

// GetLockedUndForAccount queries for locked UND given an address. An
// error is returned if the query or decoding fails.
func (ar LockedUndRetriever) GetLockedUnd(addr sdk.AccAddress) (types.LockedUnd, error) {
	lockedUnd, _, err := ar.GetLockedUndHeight(addr)
	return lockedUnd, err
}

// GetLockedUndHeight queries for locked UND  given an address. Returns the
// height of the query with the account. An error is returned if the query
// or decoding fails.
func (ar LockedUndRetriever) GetLockedUndHeight(addr sdk.AccAddress) (types.LockedUnd, int64, error) {

	params, _, err := ar.getParams()
	if err != nil {
		return types.NewLockedUnd(addr, undtypes.DefaultDenomination), 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, QueryGetLocked, addr.String()), nil)
	if err != nil {
		return types.NewLockedUnd(addr, params.Denom), 0, err
	}

	var lockedUnd types.LockedUnd
	if err := types.ModuleCdc.UnmarshalJSON(res, &lockedUnd); err != nil {
		return types.NewLockedUnd(addr, params.Denom), height, err
	}

	return lockedUnd, height, nil
}

func (ar LockedUndRetriever) getParams() (types.Params, int64, error) {
	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, QueryParameters), nil)
	if err != nil {
		return types.NewParams(sdk.AccAddress{}, undtypes.DefaultDenomination), 0, err
	}

	var params types.Params
	if err := types.ModuleCdc.UnmarshalJSON(res, &params); err != nil {
		return types.NewParams(sdk.AccAddress{}, undtypes.DefaultDenomination), 0, err
	}
	return params, height, nil
}

// TotalSupplyRetriever defines the properties of a type that can be used to
// retrieve total UND supply.
type TotalSupplyRetriever struct {
	querier NodeQuerier
}

// NewTotalSupplyRetriever initialises a new TotalSupplyRetriever instance.
func NewTotalSupplyRetriever(querier NodeQuerier) TotalSupplyRetriever {
	return TotalSupplyRetriever{querier: querier}
}

// GetLockedUndForAccount queries for locked UND given an address. An
// error is returned if the query or decoding fails.
func (ar TotalSupplyRetriever) GetTotalSupply() (types.UndSupply, error) {
	totalSupply, _, err := ar.GetTotalSupplyHeight()
	return totalSupply, err
}

// GetLockedUndHeight queries for locked UND  given an address. Returns the
// height of the query with the account. An error is returned if the query
// or decoding fails.
func (ar TotalSupplyRetriever) GetTotalSupplyHeight() (types.UndSupply, int64, error) {

	params, _, err := ar.getParams()
	if err != nil {
		return types.NewUndSupply(undtypes.DefaultDenomination), 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, QueryTotalSupply), nil)
	if err != nil {
		return types.NewUndSupply(params.Denom), 0, err
	}

	var totalSupply types.UndSupply
	if err := types.ModuleCdc.UnmarshalJSON(res, &totalSupply); err != nil {
		return types.NewUndSupply(params.Denom), height, err
	}

	return totalSupply, height, nil
}

func (ar TotalSupplyRetriever) getParams() (types.Params, int64, error) {
	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, QueryParameters), nil)
	if err != nil {
		return types.NewParams(sdk.AccAddress{}, undtypes.DefaultDenomination), 0, err
	}

	var params types.Params
	if err := types.ModuleCdc.UnmarshalJSON(res, &params); err != nil {
		return types.NewParams(sdk.AccAddress{}, undtypes.DefaultDenomination), 0, err
	}
	return params, height, nil
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
		return types.NewParams(sdk.AccAddress{}, undtypes.DefaultDenomination), 0, err
	}

	var params types.Params
	if err := types.ModuleCdc.UnmarshalJSON(res, &params); err != nil {
		return types.NewParams(sdk.AccAddress{}, undtypes.DefaultDenomination), 0, err
	}

	return params, height, nil
}
