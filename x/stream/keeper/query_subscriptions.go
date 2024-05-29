package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/unification-com/mainchain/x/stream/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q Keeper) Streams(c context.Context, req *types.QueryStreamsRequest) (*types.QueryStreamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var streams []types.Stream

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.storeKey)

	streamsStore := prefix.NewStore(store, types.StreamKeyPrefix)

	pageRes, err := query.FilteredPaginate(streamsStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var s types.Stream

		if err := q.cdc.Unmarshal(value, &s); err != nil {
			return false, status.Error(codes.Internal, err.Error())
		}

		if accumulate {
			streams = append(streams, s)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryStreamsResponse{
		Streams: streams, Pagination: pageRes,
	}, nil
}
