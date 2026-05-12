package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		receiverAddr, senderAddr, denom := types.AddressesFromStreamKey(append(types.StreamKeyPrefix, key...))

		return &types.StreamResult{
			Receiver: receiverAddr.String(),
			Sender:   senderAddr.String(),
			Stream:   stream,
			Denom:    denom,
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
		receiverAddr, s, denom := types.AddressesFromStreamKey(append(types.StreamKeyPrefix, key...))

		// filter by sender address
		if !s.Equals(senderAddr) {
			return nil, nil
		}

		return &types.StreamResult{
			Receiver: receiverAddr.String(),
			Sender:   senderAddr.String(),
			Stream:   stream,
			Denom:    denom,
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

	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	stream, ok := q.GetStream(ctx, receiverAddr, senderAddr, req.Denom)

	if !ok {
		return nil, errorsmod.Wrap(types.ErrInvalidData, "stream not found")
	}

	return &types.QueryStreamByReceiverSenderResponse{
		Stream: types.StreamResult{
			Receiver: req.ReceiverAddr,
			Sender:   req.SenderAddr,
			Stream:   &stream,
			Denom:    req.Denom,
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

	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	stream, ok := q.GetStream(ctx, receiverAddr, senderAddr, req.Denom)

	if !ok {
		return nil, errorsmod.Wrap(types.ErrInvalidData, "stream not found")
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
		senderAddr, denom := senderAndDenomFromPairRemainder(key)

		return &types.StreamResult{
			Receiver: receiverAddr.String(),
			Sender:   senderAddr.String(),
			Stream:   stream,
			Denom:    denom,
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

func (q Keeper) AllStreamsByPair(c context.Context, req *types.QueryAllStreamsByPairRequest) (*types.QueryAllStreamsByPairResponse, error) {
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

	store := prefix.NewStore(ctx.KVStore(q.storeKey), types.GetStreamsByPairKey(receiverAddr, senderAddr))

	streams, pageRes, err := query.GenericFilteredPaginate(q.cdc, store, req.Pagination, func(key []byte, stream *types.Stream) (*types.StreamResult, error) {
		// key inside the pair-prefix store is just len(denom)|denom
		denom := denomFromLengthPrefixed(key)

		return &types.StreamResult{
			Receiver: receiverAddr.String(),
			Sender:   senderAddr.String(),
			Stream:   stream,
			Denom:    denom,
		}, nil
	}, func() *types.Stream { return &types.Stream{} })

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStreamsByPairResponse{
		ReceiverAddr: req.ReceiverAddr,
		SenderAddr:   req.SenderAddr,
		Streams:      streams,
		Pagination:   pageRes,
	}, nil
}

// senderAndDenomFromPairRemainder parses (sender, denom) out of a key that has had its
// 0x11|len(receiver)|receiver prefix stripped — i.e. just len(sender)|sender|len(denom)|denom.
func senderAndDenomFromPairRemainder(key []byte) (sdk.AccAddress, string) {
	if len(key) == 0 {
		return nil, ""
	}
	senderLen := int(key[0])
	if len(key) < 1+senderLen+1 {
		return nil, ""
	}
	senderAddr := sdk.AccAddress(key[1 : 1+senderLen])
	denomLen := int(key[1+senderLen])
	if len(key) < 1+senderLen+1+denomLen {
		return senderAddr, ""
	}
	denom := string(key[1+senderLen+1 : 1+senderLen+1+denomLen])
	return senderAddr, denom
}

// denomFromLengthPrefixed parses a denom out of a single length-prefixed entry: len(denom)|denom
func denomFromLengthPrefixed(key []byte) string {
	if len(key) == 0 {
		return ""
	}
	denomLen := int(key[0])
	if len(key) < 1+denomLen {
		return ""
	}
	return string(key[1 : 1+denomLen])
}
