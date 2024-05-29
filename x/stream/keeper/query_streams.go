package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (q Keeper) AllStreamsForSender(c context.Context, req *types.QueryAllStreamsForSenderRequest) (*types.QueryAllStreamsForSenderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	_, err := sdk.AccAddressFromBech32(req.SenderAddr)
	if err != nil {
		return nil, err
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

		if req.SenderAddr == s.Sender {
			if accumulate {
				streams = append(streams, s)
			}

			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStreamsForSenderResponse{
		Streams: streams, Pagination: pageRes,
	}, nil
}

func (q Keeper) StreamByReceiverSender(c context.Context, req *types.QueryStreamByReceiverSenderRequest) (*types.QueryStreamByReceiverSenderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	receiverAddr, err := sdk.AccAddressFromBech32(req.ReceiverAddr)
	if err != nil {
		return nil, err
	}

	senderAddr, err := sdk.AccAddressFromBech32(req.SenderAddr)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	stream, ok := q.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream not found")
	}

	return &types.QueryStreamByReceiverSenderResponse{Stream: stream}, nil
}

func (q Keeper) StreamById(c context.Context, req *types.QueryStreamByIdRequest) (*types.QueryStreamByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StreamId == 0 {
		return nil, status.Error(codes.InvalidArgument, "stream id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	streamLookup, ok := q.GetIdLookup(ctx, req.StreamId)

	if !ok {
		return nil, status.Errorf(codes.NotFound, "stream %d doesn't exist", req.StreamId)
	}

	receiverAddr, err := sdk.AccAddressFromBech32(streamLookup.Receiver)
	if err != nil {
		return nil, err
	}

	senderAddr, err := sdk.AccAddressFromBech32(streamLookup.Sender)
	if err != nil {
		return nil, err
	}

	stream, ok := q.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream not found")
	}

	return &types.QueryStreamByIdResponse{Stream: stream}, nil
}

func (q Keeper) StreamByIdCurrentFlow(c context.Context, req *types.QueryStreamByIdCurrentFlowRequest) (*types.QueryStreamByIdCurrentFlowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StreamId == 0 {
		return nil, status.Error(codes.InvalidArgument, "stream id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	streamLookup, ok := q.GetIdLookup(ctx, req.StreamId)

	if !ok {
		return nil, status.Errorf(codes.NotFound, "stream %d doesn't exist", req.StreamId)
	}

	receiverAddr, err := sdk.AccAddressFromBech32(streamLookup.Receiver)
	if err != nil {
		return nil, err
	}

	senderAddr, err := sdk.AccAddressFromBech32(streamLookup.Sender)
	if err != nil {
		return nil, err
	}

	stream, ok := q.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream not found")
	}

	nowTime := ctx.BlockTime()
	currentFlow := stream.FlowRate

	if stream.DepositZeroTime.Before(nowTime) {
		currentFlow = 0
	}

	if stream.Deposit.IsNil() || stream.Deposit.IsZero() || stream.Deposit.IsNegative() {
		currentFlow = 0
	}

	return &types.QueryStreamByIdCurrentFlowResponse{
		ConfiguredFlowRate: stream.FlowRate,
		CurrentFlowRate:    currentFlow,
	}, nil
}

func (q Keeper) StreamReceiverSenderCurrentFlow(c context.Context, req *types.QueryStreamReceiverSenderCurrentFlowRequest) (*types.QueryStreamReceiverSenderCurrentFlowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	receiverAddr, err := sdk.AccAddressFromBech32(req.ReceiverAddr)
	if err != nil {
		return nil, err
	}

	senderAddr, err := sdk.AccAddressFromBech32(req.SenderAddr)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	stream, ok := q.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream not found")
	}

	nowTime := ctx.BlockTime()
	currentFlow := stream.FlowRate

	if stream.DepositZeroTime.Before(nowTime) {
		currentFlow = 0
	}

	if stream.Deposit.IsNil() || stream.Deposit.IsZero() || stream.Deposit.IsNegative() {
		currentFlow = 0
	}

	return &types.QueryStreamReceiverSenderCurrentFlowResponse{
		ConfiguredFlowRate: stream.FlowRate,
		CurrentFlowRate:    currentFlow,
	}, nil
}

func (q Keeper) AllStreamsForReceiver(c context.Context, req *types.QueryAllStreamsForReceiverRequest) (*types.QueryAllStreamsForReceiverResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	receiverAddr, err := sdk.AccAddressFromBech32(req.ReceiverAddr)
	if err != nil {
		return nil, err
	}

	var streams []types.Stream

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.storeKey)

	streamsStore := prefix.NewStore(store, types.GetStreamsByReceiverKey(receiverAddr))

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

	return &types.QueryAllStreamsForReceiverResponse{
		ReceiverAddr: req.ReceiverAddr,
		Streams:      streams,
		Pagination:   pageRes,
	}, nil
}
