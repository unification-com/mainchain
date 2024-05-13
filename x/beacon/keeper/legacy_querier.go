package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/beacon/types"
	"strconv"
)

// query endpoints supported by the beacon Querier
const (
	QueryParameters      = "params"
	QueryBeacon          = "beacon"
	QueryBeacons         = "beacons"
	QueryBeaconTimestamp = "timestamp"
)

// NewQuerier is the module level router for state queries
func NewLegacyQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, keeper, legacyQuerierCdc)
		case QueryBeacon:
			return queryBeacon(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case QueryBeacons:
			return queryBeaconsFiltered(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case QueryBeaconTimestamp:
			return queryBeaconTimestamp(ctx, path[1:], req, keeper, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func queryBeacon(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	beaconID, err := strconv.Atoi(path[0])

	if err != nil {
		beaconID = 0
	}

	beacon, found := keeper.GetBeacon(ctx, uint64(beaconID))

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBeaconDoesNotExist, "beacon %d not found", beaconID)
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, beacon)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBeaconTimestamp(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	beaconID, err := strconv.Atoi(path[0])

	if err != nil {
		beaconID = 0
	}

	timestampID, err := strconv.Atoi(path[1])

	if err != nil {
		timestampID = 0
	}

	beacon, found := keeper.GetBeacon(ctx, uint64(beaconID))

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBeaconDoesNotExist, "beacon %d not found", beaconID)
	}

	timestamp, found := keeper.GetBeaconTimestampByID(ctx, uint64(beaconID), uint64(timestampID))

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBeaconDoesNotExist, "timestamp %d not found for beacon %d", timestampID, beaconID)
	}

	legactTs := types.BeaconTimestampLegacy{
		BeaconID:    beacon.BeaconId,
		TimestampID: timestamp.TimestampId,
		SubmitTime:  timestamp.SubmitTime,
		Hash:        timestamp.Hash,
		Owner:       beacon.Owner,
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, legactTs)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBeaconsFiltered(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var queryParams types.QueryBeaconsFilteredRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &queryParams)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	filteredBeacons := k.GetBeaconsFiltered(ctx, queryParams)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, filteredBeacons)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil

}
