package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/beacon/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (q Keeper) Beacon(c context.Context, req *types.QueryBeaconRequest) (*types.QueryBeaconResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.BeaconId == 0 {
		return nil, status.Error(codes.InvalidArgument, "beacon id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	beacon, found := q.GetBeacon(ctx, req.BeaconId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "beacon %d doesn't exist", req.BeaconId)
	}

	return &types.QueryBeaconResponse{Beacon: &beacon}, nil
}

func (q Keeper) BeaconTimestamp(c context.Context, req *types.QueryBeaconTimestampRequest) (*types.QueryBeaconTimestampResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.BeaconId == 0 {
		return nil, status.Error(codes.InvalidArgument, "beacon id can not be 0")
	}

	if req.TimestampId == 0 {
		return nil, status.Error(codes.InvalidArgument, "timestamp id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	beacon, found := q.GetBeacon(ctx, req.BeaconId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "beacon %d doesn't exist in state", req.BeaconId)
	}

	beaconTimestamp, found := q.GetBeaconTimestampByID(ctx, req.BeaconId, req.TimestampId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "timestamp %d doesn't exist in state for beacon %d", req.TimestampId, req.BeaconId)
	}

	return &types.QueryBeaconTimestampResponse{
		Timestamp: &beaconTimestamp,
		BeaconId:  beacon.BeaconId,
		Owner:     beacon.Owner,
	}, nil
}

func (q Keeper) BeaconsFiltered(c context.Context, req *types.QueryBeaconsFilteredRequest) (*types.QueryBeaconsFilteredResponse, error) {
	var beacons []types.Beacon

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.storeKey)

	beaconsStore := prefix.NewStore(store, types.RegisteredBeaconPrefix)

	pageRes, err := query.FilteredPaginate(beaconsStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var b types.Beacon
		if err := q.cdc.Unmarshal(value, &b); err != nil {
			return false, status.Error(codes.Internal, err.Error())
		}

		matchOwner, matchMoniker := true, true

		if len(req.Owner) > 0 {
			_, err := sdk.AccAddressFromBech32(req.Owner)
			if err != nil {
				return false, err
			}

			matchOwner = b.Owner == req.Owner
		}

		if len(req.Moniker) > 0 {
			matchMoniker = b.Moniker == req.Moniker
		}

		if matchOwner && matchMoniker {
			if accumulate {
				beacons = append(beacons, b)
			}

			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBeaconsFilteredResponse{
		Beacons: beacons, Pagination: pageRes,
	}, nil
}

func (q Keeper) BeaconStorage(c context.Context, req *types.QueryBeaconStorageRequest) (*types.QueryBeaconStorageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.BeaconId == 0 {
		return nil, status.Error(codes.InvalidArgument, "beacon id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	beacon, found := q.GetBeacon(ctx, req.BeaconId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "beacon %d doesn't exist in state", req.BeaconId)
	}

	maxStorageLimit := q.GetParamMaxStorageLimit(ctx)

	return &types.QueryBeaconStorageResponse{
		BeaconId:       beacon.BeaconId,
		Owner:          beacon.Owner,
		CurrentLimit:   beacon.InStateLimit,
		CurrentUsed:    beacon.NumInState,
		Max:            maxStorageLimit,
		MaxPurchasable: maxStorageLimit - beacon.InStateLimit,
	}, nil
}
