package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

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

// BeaconTimestampsByHash returns every timestamp of a beacon that recorded a given hash (#129). The hash
// index is FORWARD-ONLY — timestamps recorded before the upgrade are not indexed and won't appear here.
func (q Keeper) BeaconTimestampsByHash(c context.Context, req *types.QueryBeaconTimestampsByHashRequest) (*types.QueryBeaconTimestampsByHashResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.BeaconId == 0 {
		return nil, status.Error(codes.InvalidArgument, "beacon id can not be 0")
	}

	if len(req.Hash) == 0 {
		return nil, status.Error(codes.InvalidArgument, "hash can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if _, found := q.GetBeacon(ctx, req.BeaconId); !found {
		return nil, status.Errorf(codes.NotFound, "beacon %d doesn't exist in state", req.BeaconId)
	}

	var timestamps []types.BeaconTimestamp
	store := ctx.KVStore(q.storeKey)

	// Keys within this prefix store are the timestampId(8); the value is empty (the index is key-only).
	hashStore := prefix.NewStore(store, types.BeaconHashIndexHashKey(req.BeaconId, req.Hash))

	pageRes, err := query.Paginate(hashStore, req.Pagination, func(key []byte, _ []byte) error {
		timestampID := types.GetTimestampIDFromBytes(key)
		ts, found := q.GetBeaconTimestampByID(ctx, req.BeaconId, timestampID)
		if found {
			timestamps = append(timestamps, ts)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBeaconTimestampsByHashResponse{
		BeaconId:   req.BeaconId,
		Timestamps: timestamps,
		Pagination: pageRes,
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

	beaconStorage, _ := q.GetBeaconStorageLimit(ctx, req.BeaconId)

	maxStorageLimit := q.GetParamMaxStorageLimit(ctx)

	return &types.QueryBeaconStorageResponse{
		BeaconId:       beacon.BeaconId,
		Owner:          beacon.Owner,
		CurrentLimit:   beaconStorage.InStateLimit,
		CurrentUsed:    beacon.NumInState,
		Max:            maxStorageLimit,
		MaxPurchasable: maxStorageLimit - beaconStorage.InStateLimit,
	}, nil
}
