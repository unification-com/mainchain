package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/unification-com/mainchain/x/stream/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateStream creates a new stream
func (k msgServer) CreateStream(goCtx context.Context, msg *types.MsgCreateStream) (*types.MsgCreateStreamResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	senderAddr, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return nil, accErr
	}
	receiverAddr, accErr := sdk.AccAddressFromBech32(msg.Receiver)
	if accErr != nil {
		return nil, accErr
	}

	if k.bankKeeper.BlockedAddr(receiverAddr) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.Receiver)
	}

	if k.IsStream(ctx, receiverAddr, senderAddr) {
		return nil, sdkerrors.Wrap(types.ErrStreamExists, "use update stream msg to modify existing stream")
	}

	if msg.Deposit.IsNil() || msg.Deposit.IsNegative() || msg.Deposit.IsZero() {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "deposit must be > zero")
	}

	if msg.FlowRate <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "flow rate must be > zero")
	}

	duration := types.CalculateDuration(msg.Deposit, msg.FlowRate)

	if duration < 60 {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "calculated duration too short. Must be > 1 minute")
	}

	// create the "empty" stream
	stream, err := k.CreateNewStream(ctx, receiverAddr, senderAddr, msg.Deposit, msg.FlowRate)

	if err != nil {
		return nil, err
	}

	// create lookup entry
	err = k.CreateIdLookup(ctx, receiverAddr, senderAddr, stream.StreamId)

	if err != nil {
		return nil, err
	}

	// add the deposit
	_, err = k.AddDeposit(ctx, receiverAddr, senderAddr, msg.Deposit)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, types.EventTypeCreateStreamAction)

	return &types.MsgCreateStreamResponse{
		StreamId: stream.StreamId,
	}, nil

}

// ClaimStream claims payment from a stream
func (k msgServer) ClaimStream(goCtx context.Context, msg *types.MsgClaimStream) (*types.MsgClaimStreamResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	receiverAddr, accErr := sdk.AccAddressFromBech32(msg.Receiver)
	if accErr != nil {
		return nil, accErr
	}

	if msg.StreamId <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream id must be > zero")
	}

	streamLookup, exists := k.GetIdLookup(ctx, msg.StreamId)

	if !exists {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream id does not exist")
	}

	if streamLookup.Receiver != msg.Receiver {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "you are not the receiver")
	}

	senderAddr, accErr := sdk.AccAddressFromBech32(streamLookup.Sender)
	if accErr != nil {
		return nil, accErr
	}

	finalClaimCoin, valFeeCoin, totalClaimValue, remainingDeposit, err := k.ClaimFromStream(ctx, receiverAddr, senderAddr)

	if err != nil {
		return nil, err
	}

	return &types.MsgClaimStreamResponse{
		StreamId:         msg.StreamId,
		TotalClaimed:     totalClaimValue,
		StreamPayment:    finalClaimCoin,
		ValidatorFee:     valFeeCoin,
		RemainingDeposit: remainingDeposit,
	}, nil
}

// TopUpDeposit adds more deposit to a stream
func (k msgServer) TopUpDeposit(goCtx context.Context, msg *types.MsgTopUpDeposit) (*types.MsgTopUpDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	senderAddr, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return nil, accErr
	}

	if msg.Deposit.IsNil() || msg.Deposit.IsNegative() || msg.Deposit.IsZero() {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "deposit must be > zero")
	}

	streamLookup, ok := k.GetIdLookup(ctx, msg.StreamId)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "lookup for stream id not found")
	}

	if msg.Sender != streamLookup.Sender {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "you are not the sender")
	}

	receiverAddr, accErr := sdk.AccAddressFromBech32(streamLookup.Receiver)
	if accErr != nil {
		return nil, accErr
	}

	if !k.IsStream(ctx, receiverAddr, senderAddr) {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream not found")
	}

	// Add the requested deposit
	_, err := k.AddDeposit(ctx, receiverAddr, senderAddr, msg.Deposit)

	if err != nil {
		return nil, err
	}

	// get updated stream data
	stream, _ := k.GetStream(ctx, receiverAddr, senderAddr)

	return &types.MsgTopUpDepositResponse{
		StreamId:        msg.StreamId,
		DepositAmount:   msg.Deposit,
		CurrentDeposit:  stream.Deposit,
		DepositZeroTime: stream.DepositZeroTime,
	}, nil

}

// UpdateFlowRate creates a new stream
func (k msgServer) UpdateFlowRate(goCtx context.Context, msg *types.MsgUpdateFlowRate) (*types.MsgUpdateFlowRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	senderAddr, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return nil, accErr
	}

	if msg.FlowRate <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "flow rate must be > zero")
	}

	streamLookup, ok := k.GetIdLookup(ctx, msg.StreamId)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "lookup for stream id not found")
	}

	if msg.Sender != streamLookup.Sender {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "you are not the sender")
	}

	receiverAddr, accErr := sdk.AccAddressFromBech32(streamLookup.Receiver)
	if accErr != nil {
		return nil, accErr
	}

	if !k.IsStream(ctx, receiverAddr, senderAddr) {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream not found")
	}

	// update the flow rate
	err := k.SetNewFlowRate(ctx, receiverAddr, senderAddr, msg.FlowRate)

	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateFlowRateResponse{
		StreamId: msg.StreamId,
		FlowRate: msg.FlowRate,
	}, nil
}

func (k msgServer) CancelStream(goCtx context.Context, msg *types.MsgCancelStream) (*types.MsgCancelStreamResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	senderAddr, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return nil, accErr
	}

	streamLookup, ok := k.GetIdLookup(ctx, msg.StreamId)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "lookup for stream id not found")
	}

	if msg.Sender != streamLookup.Sender {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "you are not the sender")
	}

	receiverAddr, accErr := sdk.AccAddressFromBech32(streamLookup.Receiver)
	if accErr != nil {
		return nil, accErr
	}

	_, ok = k.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "stream not found")
	}

	// cancel stream
	err := k.CancelStreamBySenderReceiver(ctx, receiverAddr, senderAddr)

	if err != nil {
		return nil, err
	}

	return &types.MsgCancelStreamResponse{
		StreamId: msg.StreamId,
	}, nil
}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
