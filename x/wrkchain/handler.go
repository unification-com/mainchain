package wrkchain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// NewHandler returns a handler for "wrkchain" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRegisterWrkChain:
			return handleMsgRegisterWrkChain(ctx, keeper, msg)
		case MsgRecordWrkChainBlock:
			return handleMsgRecordWrkChainBlock(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized WRKChain Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to register a new WRKChain
func handleMsgRegisterWrkChain(ctx sdk.Context, keeper Keeper, msg MsgRegisterWrkChain) sdk.Result {
	//if keeper.IsWrkChainRegistered(ctx, msg.Moniker) { // Checks if the WrkChain is already registered
	//	return sdk.ErrUnauthorized("WRKChain already registered").Result() // If so, throw an error
	//}

	// todo - search by monker to see if already registered


	wrkChainID, err := keeper.RegisterWrkChain(ctx, msg.Moniker, msg.WrkChainName, msg.GenesisHash, msg.Owner) // register the WRKChain

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeRegisterWrkChain,
			sdk.NewAttribute(AttributeKeyWrkChainId, strconv.FormatUint(wrkChainID, 10)),
			sdk.NewAttribute(AttributeKeyWrkChainMoniker, msg.Moniker),
			sdk.NewAttribute(AttributeKeyWrkChainName, msg.WrkChainName),
			sdk.NewAttribute(AttributeKeyWrkChainGenesisHash, msg.GenesisHash),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
		Data: GetWrkChainIDBytes(wrkChainID),
	}
}

func handleMsgRecordWrkChainBlock(ctx sdk.Context, keeper Keeper, msg MsgRecordWrkChainBlock) sdk.Result {
	if !keeper.IsWrkChainRegistered(ctx, msg.WrkChainID) { // Checks if the WrkChain is already registered
		return sdk.ErrUnauthorized("WRKChain has not been registered yet").Result() // If not, throw an error
	}

	if !keeper.IsAuthorisedToRecord(ctx, msg.WrkChainID, msg.Owner) {
		return sdk.ErrUnauthorized("you are not the owner of this WRKChain").Result()
	}

	if keeper.IsWrkChainBlockRecorded(ctx, msg.WrkChainID, msg.Height) {
		return sdk.ErrUnauthorized("WRKChain block hashes have already been recorded for this height").Result()
	}

	keeper.RecordWrkchainHashes(ctx, msg.WrkChainID, msg.Height, msg.BlockHash, msg.ParentHash, msg.Hash1, msg.Hash2, msg.Hash3, msg.Owner)

	ctx.EventManager().EmitEvents(sdk.Events{
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
	})
	return sdk.Result{Events: ctx.EventManager().Events()}

}
