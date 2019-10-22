package wrkchain

import (
	"fmt"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgRegisterWrkChain:
			return handleMsgRegisterWrkChain(ctx, keeper, msg)
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
	// todo - subtract registration fee

	keeper.RegisterWrkChain(ctx, msg.WrkChainID, msg.WrkChainName, msg.GenesisHash, msg.Owner) // If so, register the WRKChain

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterWrkChain,
			sdk.NewAttribute(types.AttributeKeyWrkChainId, msg.WrkChainID),
			sdk.NewAttribute(types.AttributeKeyWrkChainName, msg.WrkChainName),
			sdk.NewAttribute(types.AttributeKeyWrkChainGenesisHash, msg.GenesisHash),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
