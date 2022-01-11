package keeper

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/unification-com/mainchain/x/wrkchain/types"
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

func (k msgServer) RegisterWrkChain(goCtx context.Context, msg *types.MsgRegisterWrkChain) (*types.MsgRegisterWrkChainResponse, error) {
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
		return nil, sdkerrors.Wrap(types.ErrMissingData, "unable to register wrkchain - must have a moniker")
	}

	wrkchainId, err := k.RegisterNewWrkChain(ctx, msg.Moniker, msg.Name, msg.GenesisHash, msg.BaseType, ownerAddr) // register the WrkChain

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
			types.EventTypeRegisterWrkChain,
			sdk.NewAttribute(types.AttributeKeyWrkChainId, strconv.FormatUint(wrkchainId, 10)),
			sdk.NewAttribute(types.AttributeKeyWrkChainMoniker, msg.Moniker),
			sdk.NewAttribute(types.AttributeKeyWrkChainName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyWrkChainGenesisHash, msg.GenesisHash),
			sdk.NewAttribute(types.AttributeKeyBaseType, msg.BaseType),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	)

	return &types.MsgRegisterWrkChainResponse{
		WrkchainId: wrkchainId,
	}, nil

}

func (k msgServer) RecordWrkChainBlock(goCtx context.Context, msg *types.MsgRecordWrkChainBlock) (*types.MsgRecordWrkChainBlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ownerAddr, accErr := sdk.AccAddressFromBech32(msg.Owner)
	if accErr != nil {
		return nil, accErr
	}

	if msg.Height == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "height must be > 0")
	}

	if len(msg.BlockHash) > 66 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "block hash too big. 66 character limit")
	}
	if len(msg.ParentHash) > 66 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "parent hash too big. 66 character limit")
	}
	if len(msg.Hash1) > 66 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "hash1 too big. 66 character limit")
	}
	if len(msg.Hash2) > 66 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "hash2 too big. 66 character limit")
	}
	if len(msg.Hash3) > 66 {
		return nil, sdkerrors.Wrap(types.ErrContentTooLarge, "hash3 too big. 66 character limit")
	}

	if !k.IsWrkChainRegistered(ctx, msg.WrkchainId) { // Checks if the WrkChain is already registered
		return nil, sdkerrors.Wrap(types.ErrWrkChainDoesNotExist, "wrkchain has not been registered yet") // If not, throw an error
	}

	if !k.IsAuthorisedToRecord(ctx, msg.WrkchainId, ownerAddr) {
		return nil, sdkerrors.Wrap(types.ErrNotWrkChainOwner, "you are not the owner of this wrkchain")
	}

	if k.QuickCheckHeightIsRecorded(ctx, msg.WrkchainId, msg.Height) {
		return nil, sdkerrors.Wrap(types.ErrWrkChainBlockAlreadyRecorded, "wrkchain block hashes have already been recorded for this height")
	}

	err := k.RecordNewWrkchainHashes(ctx, msg.WrkchainId, msg.Height, msg.BlockHash, msg.ParentHash, msg.Hash1, msg.Hash2, msg.Hash3, ownerAddr)

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
			types.EventTypeRecordWrkChainBlock,
			sdk.NewAttribute(types.AttributeKeyWrkChainId, strconv.FormatUint(msg.WrkchainId, 10)),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatUint(msg.Height, 10)),
			sdk.NewAttribute(types.AttributeKeyBlockHash, msg.BlockHash),
			sdk.NewAttribute(types.AttributeKeyParentHash, msg.ParentHash),
			sdk.NewAttribute(types.AttributeKeyHash1, msg.Hash1),
			sdk.NewAttribute(types.AttributeKeyHash2, msg.Hash2),
			sdk.NewAttribute(types.AttributeKeyHash3, msg.Hash3),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	)

	return &types.MsgRecordWrkChainBlockResponse{
		WrkchainId: msg.WrkchainId,
		Height:     msg.Height,
	}, nil

}
