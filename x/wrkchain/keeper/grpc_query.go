package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/wrkchain/types"
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

func (q Keeper) WrkChain(c context.Context, req *types.QueryWrkChainRequest) (*types.QueryWrkChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.WrkchainId == 0 {
		return nil, status.Error(codes.InvalidArgument, "wrkchain id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	wrkchain, found := q.GetWrkChain(ctx, req.WrkchainId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "wrkchain %d doesn't exist", req.WrkchainId)
	}

	return &types.QueryWrkChainResponse{Wrkchain: &wrkchain}, nil
}

func (q Keeper) WrkChainBlock(c context.Context, req *types.QueryWrkChainBlockRequest) (*types.QueryWrkChainBlockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.WrkchainId == 0 {
		return nil, status.Error(codes.InvalidArgument, "wrkchain id can not be 0")
	}

	// todo - if height == 0, return genesis
	if req.Height == 0 {
		return nil, status.Error(codes.InvalidArgument, "height can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	wrkchainBlock := q.GetWrkChainBlock(ctx, req.WrkchainId, req.Height)

	return &types.QueryWrkChainBlockResponse{Block: &wrkchainBlock}, nil
}

func (q Keeper) WrkChainsFiltered(c context.Context, req *types.QueryWrkChainsFilteredRequest) (*types.QueryWrkChainsFilteredResponse, error) {
	var wrkchains []types.WrkChain

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.storeKey)

	wrkchainStore := prefix.NewStore(store, types.RegisteredWrkChainPrefix)

	pageRes, err := query.FilteredPaginate(wrkchainStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var wc types.WrkChain
		if err := q.cdc.UnmarshalBinaryBare(value, &wc); err != nil {
			return false, status.Error(codes.Internal, err.Error())
		}

		matchOwner, matchMoniker := true, true

		if len(req.Owner) > 0 {
			_, err := sdk.AccAddressFromBech32(req.Owner)
			if err != nil {
				return false, err
			}

			matchOwner = wc.Owner == req.Owner
		}

		if len(req.Moniker) > 0 {
			matchMoniker = wc.Moniker == req.Moniker
		}

		if matchOwner && matchMoniker {
			if accumulate {
				wrkchains = append(wrkchains, wc)
			}

			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWrkChainsFilteredResponse{
		Wrkchains: wrkchains, Pagination: pageRes,
	}, nil
}
