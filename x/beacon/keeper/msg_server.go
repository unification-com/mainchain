package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

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
		return nil, errorsmod.Wrap(types.ErrContentTooLarge, "name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return nil, errorsmod.Wrap(types.ErrContentTooLarge, "moniker too big. 64 character limit")
	}

	if len(msg.Moniker) == 0 {
		return nil, errorsmod.Wrap(types.ErrMissingData, "unable to register beacon - must have a moniker")
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
		return nil, errorsmod.Wrap(types.ErrContentTooLarge, "hash too big. 66 character limit")
	}

	if !k.IsBeaconRegistered(ctx, msg.BeaconId) { // Checks if the BEACON is registered
		return nil, errorsmod.Wrap(types.ErrBeaconDoesNotExist, "beacon has not been registered yet") // If not, throw an error
	}

	if !k.IsAuthorisedToRecord(ctx, msg.BeaconId, ownerAddr) {
		return nil, errorsmod.Wrap(types.ErrNotBeaconOwner, "you are not the owner of this beacon")
	}

	subtime := msg.SubmitTime

	if subtime == 0 {
		subtime = uint64(time.Now().Unix())
	}

	tsID, deleteTimestampId, err := k.RecordNewBeaconTimestamp(ctx, msg.BeaconId, msg.Hash, subtime)

	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, types.RecordAction)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecordBeaconTimestamp,
			sdk.NewAttribute(types.AttributeKeyBeaconId, strconv.FormatUint(msg.BeaconId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestampID, strconv.FormatUint(tsID, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestampIdPruned, strconv.FormatUint(deleteTimestampId, 10)),
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

func (k msgServer) PurchaseBeaconStateStorage(goCtx context.Context, msg *types.MsgPurchaseBeaconStateStorage) (*types.MsgPurchaseBeaconStateStorageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ownerAddr, accErr := sdk.AccAddressFromBech32(msg.Owner)
	if accErr != nil {
		return nil, accErr
	}

	if msg.Number == 0 {
		return nil, errorsmod.Wrap(types.ErrContentTooLarge, "cannot purchase zero")
	}

	_, found := k.GetBeacon(ctx, msg.BeaconId)

	if !found { // Checks if the BEACON is registered
		return nil, errorsmod.Wrap(types.ErrBeaconDoesNotExist, "beacon has not been registered yet") // If not, throw an error
	}

	if !k.IsAuthorisedToRecord(ctx, msg.BeaconId, ownerAddr) {
		return nil, errorsmod.Wrap(types.ErrNotBeaconOwner, "you are not the owner of this beacon")
	}

	beaconStorage, _ := k.GetBeaconStorageLimit(ctx, msg.BeaconId)

	// check not exceeding max
	maxParam := k.GetParamMaxStorageLimit(ctx)
	beaconStorageAfter := beaconStorage.InStateLimit + msg.Number

	if beaconStorageAfter > maxParam {
		return nil, errorsmod.Wrap(types.ErrExceedsMaxStorage, fmt.Sprintf("%d will exceed max storage of %d", beaconStorageAfter, maxParam))
	}

	err := k.IncreaseInStateStorage(ctx, msg.BeaconId, msg.Number)

	if err != nil {
		return nil, err
	}

	// get remianing can purchase
	numCanPurchase := k.GetMaxPurchasableSlots(ctx, msg.BeaconId)

	defer telemetry.IncrCounter(1, types.ModuleName, types.PurchaseStorageAction)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePurchaseStorage,
			sdk.NewAttribute(types.AttributeKeyBeaconId, strconv.FormatUint(msg.BeaconId, 10)),
			sdk.NewAttribute(types.AttributeKeyBeaconStorageNumPurchased, strconv.FormatUint(msg.Number, 10)),
			sdk.NewAttribute(types.AttributeKeyBeaconStorageNumCanPurchase, strconv.FormatUint(numCanPurchase, 10)),
			sdk.NewAttribute(types.AttributeKeyOwner, ownerAddr.String()),
		),
	)

	return &types.MsgPurchaseBeaconStateStorageResponse{
		BeaconId:        msg.BeaconId,
		NumberPurchased: msg.Number,
		NumCanPurchase:  numCanPurchase,
	}, nil

}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
