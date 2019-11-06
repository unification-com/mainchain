package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
// retrieve accounts.
type LockedUndRetriever struct {
	querier NodeQuerier
}

// NewLockedUndRetriever initialises a new LockedUndRetriever instance.
func NewLockedUndRetriever(querier NodeQuerier) LockedUndRetriever {
	return LockedUndRetriever{querier: querier}
}

// GetLockedUnd queries for locked UND given an address. An
// error is returned if the query or decoding fails.
func (ar LockedUndRetriever) GetLockedUnd(addr sdk.AccAddress) (types.LockedUnd, error) {
	lockedUnd, _, err := ar.GetLockedUndHeight(addr)
	return lockedUnd, err
}

// GetLockedUndHeight queries for locked UND  given an address. Returns the
// height of the query with the account. An error is returned if the query
// or decoding fails.
func (ar LockedUndRetriever) GetLockedUndHeight(addr sdk.AccAddress) (types.LockedUnd, int64, error) {

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, QueryGetLocked, addr.String()), nil)
	if err != nil {
		return types.NewLockedUnd(addr), 0, err
	}

	var lockedUnd types.LockedUnd
	if err := types.ModuleCdc.UnmarshalJSON(res, &lockedUnd); err != nil {
		return types.NewLockedUnd(addr), height, err
	}

	return lockedUnd, height, nil
}
