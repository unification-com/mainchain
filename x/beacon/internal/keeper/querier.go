package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the beacon Querier
const (
	QueryParameters      = "params"
	QueryBeacon          = "beacon"
	QueryBeacons         = "beacons"
	QueryBeaconTimestamp = "timestamp"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, keeper)
		case QueryBeacon:
			return queryBeacon(ctx, path[1:], req, keeper)
		case QueryBeacons:
			return queryBeaconsFiltered(ctx, path[1:], req, keeper)
		case QueryBeaconTimestamp:
			return queryBeaconTimestamp(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func queryBeacon(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

	beaconID, err := strconv.Atoi(path[0])

	if err != nil {
		beaconID = 0
	}

	beacon := keeper.GetBeacon(ctx, uint64(beaconID))

	res, err := codec.MarshalJSONIndent(keeper.cdc, beacon)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBeaconTimestamp(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

	beaconID, err := strconv.Atoi(path[0])

	if err != nil {
		beaconID = 0
	}

	timestampID, err := strconv.Atoi(path[1])

	if err != nil {
		timestampID = 0
	}

	timestamp := keeper.GetBeaconTimestampByID(ctx, uint64(beaconID), uint64(timestampID))

	res, err := codec.MarshalJSONIndent(keeper.cdc, timestamp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBeaconsFiltered(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {

	var queryParams types.QueryBeaconParams

	err := k.cdc.UnmarshalJSON(req.Data, &queryParams)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	filteredBeacons := k.GetBeaconsFiltered(ctx, queryParams)

	if filteredBeacons == nil {
		filteredBeacons = types.Beacons{}
	}

	res, err := codec.MarshalJSONIndent(k.cdc, filteredBeacons)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
