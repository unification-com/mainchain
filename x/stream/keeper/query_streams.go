package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/unification-com/mainchain/x/stream/types"
)

func (q Keeper) Streams(c context.Context, req *types.QueryStreamsRequest) (*types.QueryStreamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(q.storeKey), types.StreamKeyPrefix)

	streams, pageRes, err := query.GenericFilteredPaginate(q.cdc, store, req.Pagination, func(key []byte, stream *types.Stream) (*types.StreamResult, error) {

		// need to prefix the StreamKeyPrefix 0x11 to the returned key as AddressesFromStreamKey expects it
		receiverAddr, senderAddr := types.AddressesFromStreamKey(append(types.StreamKeyPrefix, key...))

		return &types.StreamResult{
			Receiver: receiverAddr.String(),
			Sender:   senderAddr.String(),
			Stream:   stream,
		}, nil
	}, func() *types.Stream { return &types.Stream{} })

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

	senderAddr, err := sdk.AccAddressFromBech32(req.SenderAddr)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(q.storeKey), types.StreamKeyPrefix)

	streams, pageRes, err := query.GenericFilteredPaginate(q.cdc, store, req.Pagination, func(key []byte, stream *types.Stream) (*types.StreamResult, error) {

		// need to prefix the StreamKeyPrefix 0x11 to the returned key as AddressesFromStreamKey expects it
		receiverAddr, s := types.AddressesFromStreamKey(append(types.StreamKeyPrefix, key...))

		// filter by sender address
		if !s.Equals(senderAddr) {
			return nil, nil
		}

		return &types.StreamResult{
			Receiver: receiverAddr.String(),
			Sender:   senderAddr.String(),
			Stream:   stream,
		}, nil
	}, func() *types.Stream { return &types.Stream{} })

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

	return &types.QueryStreamByReceiverSenderResponse{
		Stream: types.StreamResult{
			Receiver: req.ReceiverAddr,
			Sender:   req.SenderAddr,
			Stream:   &stream,
		},
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

	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(q.storeKey), types.GetStreamsByReceiverKey(receiverAddr))

	streams, pageRes, err := query.GenericFilteredPaginate(q.cdc, store, req.Pagination, func(key []byte, stream *types.Stream) (*types.StreamResult, error) {
		senderAddr := types.FirstAddressFromStreamStoreKey(key)

		return &types.StreamResult{
			Receiver: receiverAddr.String(),
			Sender:   senderAddr.String(),
			Stream:   stream,
		}, nil
	}, func() *types.Stream { return &types.Stream{} })

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStreamsForReceiverResponse{
		ReceiverAddr: req.ReceiverAddr,
		Streams:      streams,
		Pagination:   pageRes,
	}, nil
}
