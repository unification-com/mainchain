package beacon

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
	"strconv"
)

// NewHandler returns a handler for "beacon" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgRegisterBeacon:
			return handleMsgRegisterBeacon(ctx, keeper, msg)
		case MsgRecordBeaconTimestamp:
			return handleMsgRecordBeaconTimestamp(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// Handle a message to register a new BEACON
func handleMsgRegisterBeacon(ctx sdk.Context, keeper Keeper, msg MsgRegisterBeacon) (*sdk.Result, error) {

	if len(msg.BeaconName) > 128 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "moniker too big. 64 character limit")
	}

	if len(msg.Moniker) == 0 {
		return nil, sdkerrors.Wrap(types.ErrMissingData, "unable to register beacon - must have a moniker")
	}

	params := NewQueryBeaconParams(1, 1, msg.Moniker, sdk.AccAddress{})
	beacons := keeper.GetBeaconsFiltered(ctx, params)

	if (len(beacons)) > 0 {
		errMsg := fmt.Sprintf("beacon already registered with moniker '%s' - id: %d, owner: %s", msg.Moniker, beacons[0].BeaconID, beacons[0].Owner)
		return nil, sdkerrors.Wrap(ErrBeaconAlreadyRegistered, errMsg)
	}

	beaconID, err := keeper.RegisterBeacon(ctx, msg.Moniker, msg.BeaconName, msg.Owner) // register the BEACON

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeRegisterBeacon,
			sdk.NewAttribute(AttributeKeyBeaconId, strconv.FormatUint(beaconID, 10)),
			sdk.NewAttribute(AttributeKeyBeaconMoniker, msg.Moniker),
			sdk.NewAttribute(AttributeKeyBeaconName, msg.BeaconName),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   GetBeaconIDBytes(beaconID),
	}, nil
}

// Handle a message to record a new BEACON timestamp
func handleMsgRecordBeaconTimestamp(ctx sdk.Context, keeper Keeper, msg MsgRecordBeaconTimestamp) (*sdk.Result, error) {


	if len(msg.Hash) > 66 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "hash too big. 66 character limit")
	}

	if !keeper.IsBeaconRegistered(ctx, msg.BeaconID) { // Checks if the BEACON is registered
		return nil, sdkerrors.Wrap(ErrBeaconDoesNotExist, "beacon has not been registered yet") // If not, throw an error
	}

	if !keeper.IsAuthorisedToRecord(ctx, msg.BeaconID, msg.Owner) {
		return nil, sdkerrors.Wrap(ErrNotBeaconOwner, "you are not the owner of this beacon")
	}

	tsID, err := keeper.RecordBeaconTimestamp(ctx, msg.BeaconID, msg.Hash, msg.SubmitTime, msg.Owner)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeRecordBeaconTimestamp,
			sdk.NewAttribute(AttributeKeyBeaconId, strconv.FormatUint(msg.BeaconID, 10)),
			sdk.NewAttribute(AttributeKeyTimestampID, strconv.FormatUint(tsID, 10)),
			sdk.NewAttribute(AttributeKeyTimestampHash, msg.Hash),
			sdk.NewAttribute(AttributeKeyTimestampSubmitTime, strconv.FormatUint(msg.SubmitTime, 10)),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   GetTimestampIDBytes(tsID),
	}, nil
}
