package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
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

func (k msgServer) UndPurchaseOrder(goCtx context.Context, msg *types.MsgUndPurchaseOrder) (*types.MsgUndPurchaseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	accAddr, accErr := sdk.AccAddressFromBech32(msg.Purchaser)
	if accErr != nil {
		return nil, accErr
	}

	if msg.Amount.Denom != k.GetParamDenom(ctx) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDenomination, fmt.Sprintf("denomination must be %s", k.GetParamDenom(ctx)))
	}

	if !msg.Amount.IsPositive() {
		return nil, sdkerrors.Wrap(types.ErrInvalidData, "amount must be > 0")
	}

	if !k.AddressIsWhitelisted(ctx, accAddr) {
		return nil, sdkerrors.Wrap(types.ErrNotAuthorisedToRaisePO, fmt.Sprintf("%s is not whitelisted to raise purchase orders", msg.Purchaser))
	}

	po := types.EnterpriseUndPurchaseOrder{
		Purchaser: msg.Purchaser,
		Amount:    msg.Amount,
	}

	purchaseOrderId, err := k.RaiseNewPurchaseOrder(ctx, po)

	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, types.PurchaseAction)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRaisePurchaseOrder,
			sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(purchaseOrderId, 10)),
			sdk.NewAttribute(types.AttributeKeyPurchaser, msg.Purchaser),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &types.MsgUndPurchaseOrderResponse{
		PurchaseOrderId: purchaseOrderId,
	}, nil
}

func (k msgServer) ProcessUndPurchaseOrder(goCtx context.Context, msg *types.MsgProcessUndPurchaseOrder) (*types.MsgProcessUndPurchaseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signer, accErr := sdk.AccAddressFromBech32(msg.Signer)
	if accErr != nil {
		return nil, accErr
	}

	if !k.IsAuthorisedToDecide(ctx, signer) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")
	}

	if !k.PurchaseOrderExists(ctx, msg.PurchaseOrderId) {
		return nil, sdkerrors.Wrapf(types.ErrPurchaseOrderDoesNotExist, "id: %d", msg.PurchaseOrderId)
	}

	if !types.ValidPurchaseOrderAcceptRejectStatus(msg.Decision) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")
	}

	purchaseOrder, found := k.GetPurchaseOrder(ctx, msg.PurchaseOrderId)

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrPurchaseOrderDoesNotExist, "purchase order id %d does not exist", msg.PurchaseOrderId)
	}

	if purchaseOrder.Status == types.StatusNil {
		return nil, sdkerrors.Wrapf(types.ErrPurchaseOrderNotRaised, "purchase order %d not raised!", msg.PurchaseOrderId)
	}

	if purchaseOrder.Status != types.StatusRaised {
		return nil, sdkerrors.Wrapf(types.ErrPurchaseOrderAlreadyProcessed, "id %d already processed: %s", msg.PurchaseOrderId, purchaseOrder.Status.String())
	}

	currentDecisions := purchaseOrder.Decisions
	for _, d := range currentDecisions {
		if msg.Signer == d.Signer {
			return nil, sdkerrors.Wrapf(types.ErrSignerAlreadyMadeDecision, "signer %s already decided: %s", msg.Signer, d.Decision.String())
		}
	}

	err := k.ProcessPurchaseOrderDecision(ctx, msg.PurchaseOrderId, msg.Decision, signer)

	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, types.ProcessAction, msg.Decision.String())

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProcessPurchaseOrderDecision,
			sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(msg.PurchaseOrderId, 10)),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyDecision, msg.Decision.String()),
		),
	)

	return &types.MsgProcessUndPurchaseOrderResponse{}, nil
}

func (k msgServer) WhitelistAddress(goCtx context.Context, msg *types.MsgWhitelistAddress) (*types.MsgWhitelistAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signer, accErr := sdk.AccAddressFromBech32(msg.Signer)
	if accErr != nil {
		return nil, sdkerrors.Wrap(accErr, "signer address")
	}

	addr, accErr := sdk.AccAddressFromBech32(msg.Address)
	if accErr != nil {
		return nil, sdkerrors.Wrap(accErr, "whitelist address")
	}

	if !k.IsAuthorisedToDecide(ctx, signer) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer modifying whitelist")
	}

	if !types.ValidWhitelistAction(msg.Action) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDecision, "action should be add or remove")
	}

	err := k.ProcessWhitelistAction(ctx, addr, msg.Action, signer)

	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, types.WhitelistAddressAction, msg.Action.String())

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWhitelistAddress,
			sdk.NewAttribute(types.AttributeWhitelistAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyWhitelistAction, msg.Action.String()),
		),
	)

	return &types.MsgWhitelistAddressResponse{}, nil
}
