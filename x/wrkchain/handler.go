package wrkchain

import (
	"fmt"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "wrkchain" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgRegisterWrkChain:
			return handleMsgRegisterWrkChain(ctx, keeper, msg)
		case MsgRecordWrkChainBlock:
			return handleMsgRecordWrkChainBlock(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// Handle a message to register a new WRKChain
func handleMsgRegisterWrkChain(ctx sdk.Context, keeper Keeper, msg MsgRegisterWrkChain) (*sdk.Result, error) {

	if len(msg.WrkChainName) > 128 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "moniker too big. 64 character limit")
	}

	if len(msg.Moniker) == 0 {
		return nil, sdkerrors.Wrap(types.ErrMissingData, "unable to set WRKChain - must have a moniker")
	}

	if len(msg.GenesisHash) > 66 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "genesis hash too big. 66 character limit")
	}

	params := NewQueryWrkChainParams(1, 1, msg.Moniker, sdk.AccAddress{})
	wrkChains := keeper.GetWrkChainsFiltered(ctx, params)

	if (len(wrkChains)) > 0 {
		errMsg := fmt.Sprintf("wrkchain already registered with moniker '%s' - id: %d, owner: %s", msg.Moniker, wrkChains[0].WrkChainID, wrkChains[0].Owner)
		return nil, sdkerrors.Wrap(ErrWrkChainAlreadyRegistered, errMsg)
	}

	wrkChainID, err := keeper.RegisterWrkChain(ctx, msg.Moniker, msg.WrkChainName, msg.GenesisHash, msg.BaseType, msg.Owner) // register the WRKChain

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
			EventTypeRegisterWrkChain,
			sdk.NewAttribute(AttributeKeyWrkChainId, strconv.FormatUint(wrkChainID, 10)),
			sdk.NewAttribute(AttributeKeyWrkChainMoniker, msg.Moniker),
			sdk.NewAttribute(AttributeKeyWrkChainName, msg.WrkChainName),
			sdk.NewAttribute(AttributeKeyWrkChainGenesisHash, msg.GenesisHash),
			sdk.NewAttribute(AttributeKeyBaseType, msg.BaseType),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   GetWrkChainIDBytes(wrkChainID),
	}, nil
}

func handleMsgRecordWrkChainBlock(ctx sdk.Context, keeper Keeper, msg MsgRecordWrkChainBlock) (*sdk.Result, error) {

	if len(msg.BlockHash) > 66 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "block hash too big. 66 character limit")
	}
	if len(msg.ParentHash) > 66 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "parent hash too big. 66 character limit")
	}
	if len(msg.Hash1) > 66 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "hash1 too big. 66 character limit")
	}
	if len(msg.Hash2) > 66 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "hash2 too big. 66 character limit")
	}
	if len(msg.Hash3) > 66 {
		return nil, sdkerrors.Wrap(ErrContentTooLarge, "hash3 too big. 66 character limit")
	}

	if !keeper.IsWrkChainRegistered(ctx, msg.WrkChainID) { // Checks if the WrkChain is already registered
		return nil, sdkerrors.Wrap(ErrWrkChainDoesNotExist, "WRKChain has not been registered yet") // If not, throw an error
	}

	if !keeper.IsAuthorisedToRecord(ctx, msg.WrkChainID, msg.Owner) {
		return nil, sdkerrors.Wrap(ErrNotWrkChainOwner, "you are not the owner of this WRKChain")
	}

	if keeper.QuickCheckHeightIsRecorded(ctx, msg.WrkChainID, msg.Height) {
		return nil, sdkerrors.Wrap(ErrWrkChainBlockAlreadyRecorded, "WRKChain block hashes have already been recorded for this height")
	}

	err := keeper.RecordWrkchainHashes(ctx, msg.WrkChainID, msg.Height, msg.BlockHash, msg.ParentHash, msg.Hash1, msg.Hash2, msg.Hash3, msg.Owner)

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
			EventTypeRecordWrkChainBlock,
			sdk.NewAttribute(AttributeKeyWrkChainId, strconv.FormatUint(msg.WrkChainID, 10)),
			sdk.NewAttribute(AttributeKeyBlockHeight, strconv.FormatUint(msg.Height, 10)),
			sdk.NewAttribute(AttributeKeyBlockHash, msg.BlockHash),
			sdk.NewAttribute(AttributeKeyParentHash, msg.ParentHash),
			sdk.NewAttribute(AttributeKeyHash1, msg.Hash1),
			sdk.NewAttribute(AttributeKeyHash2, msg.Hash2),
			sdk.NewAttribute(AttributeKeyHash3, msg.Hash3),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}
