package beacon


import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// NewHandler returns a handler for "beacon" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRegisterBeacon:
			return handleMsgRegisterBeacon(ctx, keeper, msg)
		case MsgRecordBeaconTimestamp:
			return handleMsgRecordBeaconTimestamp(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognised beacon Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to register a new BEACON
func handleMsgRegisterBeacon(ctx sdk.Context, keeper Keeper, msg MsgRegisterBeacon) sdk.Result {

	params := NewQueryBeaconParams(1, 1, msg.Moniker, sdk.AccAddress{})
	beacons := keeper.GetBeaconsFiltered(ctx, params)

	if (len(beacons)) > 0 {
		errMsg := fmt.Sprintf("beacon already registered with moniker '%s' - id: %d, owner: %s", msg.Moniker, beacons[0].BeaconID, beacons[0].Owner)
		return ErrBeaconAlreadyRegistered(keeper.Codespace(), errMsg).Result()
	}

	beaconID, err := keeper.RegisterBeacon(ctx, msg.Moniker, msg.BeaconName, msg.Owner) // register the BEACON

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeRegisterBeacon,
			sdk.NewAttribute(AttributeKeyBeaconId, strconv.FormatUint(beaconID, 10)),
			sdk.NewAttribute(AttributeKeyBeaconMoniker, msg.Moniker),
			sdk.NewAttribute(AttributeKeyBeaconName, msg.BeaconName),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   GetBeaconIDBytes(beaconID),
	}
}

// Handle a message to record a new BEACON timestamp
func handleMsgRecordBeaconTimestamp(ctx sdk.Context, keeper Keeper, msg MsgRecordBeaconTimestamp) sdk.Result {
	if !keeper.IsBeaconRegistered(ctx, msg.BeaconID) { // Checks if the BEACON is registered
		return ErrBeaconDoesNotExist(keeper.Codespace(), "beacon has not been registered yet").Result() // If not, throw an error
	}

	if !keeper.IsAuthorisedToRecord(ctx, msg.BeaconID, msg.Owner) {
		return ErrNotBeaconOwner(keeper.Codespace(), "you are not the owner of this beacon").Result()
	}

	if keeper.IsBeaconTimestampRecordedByHashTime(ctx, msg.BeaconID, msg.Hash, msg.SubmitTime) {
		return ErrBeaconTimestampAlreadyRecorded(keeper.Codespace(), fmt.Sprintf("timestamp hash %s already recorded at time %d", msg.Hash, msg.SubmitTime)).Result()
	}

	tsID, err := keeper.RecordBeaconTimestamp(ctx, msg.BeaconID, msg.Hash, msg.SubmitTime, msg.Owner)

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeRecordBeaconTimestamp,
			sdk.NewAttribute(AttributeKeyBeaconId, strconv.FormatUint(msg.BeaconID, 10)),
			sdk.NewAttribute(AttributeKeyTimestampID, strconv.FormatUint(tsID, 10)),
			sdk.NewAttribute(AttributeKeyTimestampHash, msg.Hash),
			sdk.NewAttribute(AttributeKeyTimestampSubmitTime, strconv.FormatUint(msg.SubmitTime, 10)),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
		Data: GetTimestampIDBytes(tsID),
	}
}
