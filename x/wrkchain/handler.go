package wrkchain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
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
	if keeper.IsWrkChainRegistered(ctx, msg.WrkChainID) { // Checks if the WrkChain is already registered
		return sdk.ErrUnauthorized("WRKChain already registered").Result() // If so, throw an error
	}

	coins := keeper.BankKeeper.GetCoins(ctx, msg.Owner)
	fees := msg.Fee

	if fees.IsAllLT(FeesWrkChainRegistration) {
		return sdk.ErrUnauthorized("not enough fees sent to pay for registration").Result()
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; require %s, have %s", coins, fees),
		).Result()
	}

	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, types.FeeCollectorName, fees)
	if err != nil {
		return err.Result()
	}

	keeper.RegisterWrkChain(ctx, msg.WrkChainID, msg.WrkChainName, msg.GenesisHash, msg.Owner) // register the WRKChain

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeRegisterWrkChain,
			sdk.NewAttribute(AttributeKeyWrkChainId, msg.WrkChainID),
			sdk.NewAttribute(AttributeKeyWrkChainName, msg.WrkChainName),
			sdk.NewAttribute(AttributeKeyWrkChainGenesisHash, msg.GenesisHash),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
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

	coins := keeper.BankKeeper.GetCoins(ctx, msg.Owner)
	fees := msg.Fee

	if fees.IsAllLT(FeesWrkChainRecordHash) {
		return sdk.ErrUnauthorized("not enough fees sent to pay for hash submission").Result()
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; require %s, have %s", coins, fees),
		).Result()
	}

	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, types.FeeCollectorName, fees)
	if err != nil {
		return err.Result()
	}

	keeper.RecordWrkchainHashes(ctx, msg.WrkChainID, msg.Height, msg.BlockHash, msg.ParentHash, msg.Hash1, msg.Hash2, msg.Hash3, msg.Owner)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeRecordWrkChainBlock,
			sdk.NewAttribute(AttributeKeyWrkChainId, msg.WrkChainID),
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
