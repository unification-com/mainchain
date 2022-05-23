package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterBeacon(goCtx context.Context, msg *types.MsgRegisterBeacon) (*types.MsgRegisterBeaconResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ownerAddr, accErr := sdk.AccAddressFromBech32(msg.Owner)
	if accErr != nil {
		return nil, accErr
	}

	if len(msg.Name) > 128 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "moniker too big. 64 character limit")
	}

	if len(msg.Moniker) == 0 {
		return nil, sdkerrors.Wrap(types.ErrMissingData, "unable to register beacon - must have a moniker")
	}

	beacon := types.Beacon{
		Moniker: msg.Moniker,
		Name:    msg.Name,
		Owner:   ownerAddr.String(),
	}

	beaconID, err := k.RegisterNewBeacon(ctx, beacon) // register the BEACON

	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, types.RegisterAction)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterBeacon,
			sdk.NewAttribute(types.AttributeKeyBeaconId, strconv.FormatUint(beaconID, 10)),
			sdk.NewAttribute(types.AttributeKeyBeaconMoniker, msg.Moniker),
			sdk.NewAttribute(types.AttributeKeyBeaconName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyOwner, ownerAddr.String()),
		),
	)

	return &types.MsgRegisterBeaconResponse{
		BeaconId: beaconID,
	}, nil

}

func (k msgServer) RecordBeaconTimestamp(goCtx context.Context, msg *types.MsgRecordBeaconTimestamp) (*types.MsgRecordBeaconTimestampResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ownerAddr, accErr := sdk.AccAddressFromBech32(msg.Owner)
	if accErr != nil {
		return nil, accErr
	}

	if len(msg.Hash) > 66 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "hash too big. 66 character limit")
	}

	if !k.IsBeaconRegistered(ctx, msg.BeaconId) { // Checks if the BEACON is registered
		return nil, sdkerrors.Wrap(types.ErrBeaconDoesNotExist, "beacon has not been registered yet") // If not, throw an error
	}

	if !k.IsAuthorisedToRecord(ctx, msg.BeaconId, ownerAddr) {
		return nil, sdkerrors.Wrap(types.ErrNotBeaconOwner, "you are not the owner of this beacon")
	}

	subtime := msg.SubmitTime

	if subtime == 0 {
		subtime = uint64(time.Now().Unix())
	}

	tsID, err := k.RecordNewBeaconTimestamp(ctx, msg.BeaconId, msg.Hash, subtime)

	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, types.RecordAction)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecordBeaconTimestamp,
			sdk.NewAttribute(types.AttributeKeyBeaconId, strconv.FormatUint(msg.BeaconId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestampID, strconv.FormatUint(tsID, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestampHash, msg.Hash),
			sdk.NewAttribute(types.AttributeKeyTimestampSubmitTime, strconv.FormatUint(msg.SubmitTime, 10)),
			sdk.NewAttribute(types.AttributeKeyOwner, ownerAddr.String()),
		),
	)

	return &types.MsgRecordBeaconTimestampResponse{
		BeaconId:    msg.BeaconId,
		TimestampId: tsID,
	}, nil

}
