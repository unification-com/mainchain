package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/stream/types"
	"strconv"
	"time"
)

// GetTotalDeposits gets the total deposits - just a wrapper for getting the module account's balances from the bank
func (k Keeper) GetTotalDeposits(ctx sdk.Context) sdk.Coins {
	moduleAcc := k.GetStreamModuleAccount(ctx)
	totalDeposits := k.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	return totalDeposits
}

// SetStream Sets the stream
func (k Keeper) SetStream(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress, stream types.Stream) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetStreamKey(receiverAddr, senderAddr), k.cdc.MustMarshal(&stream))

	return nil
}

// IsStream Checks if the stream is present in the store or not
func (k Keeper) IsStream(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetStreamKey(receiverAddr, senderAddr))
}

// GetStream Gets the stream data
func (k Keeper) GetStream(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress) (types.Stream, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStreamKey(receiverAddr, senderAddr))
	if bz == nil {
		// return a new empty stream struct
		return types.Stream{}, false
	}
	var stream types.Stream
	k.cdc.MustUnmarshal(bz, &stream)
	return stream, true
}

// IterateAllStreams iterates over all the Streams of all accounts
// that are provided to a callback. If true is returned from the
// callback, iteration is halted. Potentially expensive, and only intended
// for use during genesis export etc.
func (k Keeper) IterateAllStreams(ctx sdk.Context, cb func(sdk.AccAddress, sdk.AccAddress, types.Stream) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StreamKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		receiverAddr, senderAddr := types.AddressesFromStreamKey(iterator.Key())

		var stream types.Stream
		err := k.cdc.Unmarshal(iterator.Value(), &stream)

		if err != nil {
			panic(err)
		}

		if cb(receiverAddr, senderAddr, stream) {
			break
		}
	}
}

func (k Keeper) ClaimFromStream(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress) (sdk.Coin, sdk.Coin, sdk.Coin, sdk.Coin, error) {
	stream, ok := k.GetStream(ctx, receiverAddr, senderAddr)
	params := k.GetParams(ctx)

	if !ok {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "stream does not exist")
	}

	// 1. check current stream deposit > 0
	if stream.Deposit.IsNil() || stream.Deposit.IsNegative() || stream.Deposit.IsZero() {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "stream deposit is zero")
	}

	// 2. calculate amount to claim
	nowTime := ctx.BlockTime()
	claimTotal, remainingDeposit := types.CalculateAmountToClaim(nowTime, stream.DepositZeroTime, stream.LastOutflowTime, stream.Deposit, stream.FlowRate)

	// 3.1 sanity check 1: claimTotal is not negative or nil
	if claimTotal.IsNil() || claimTotal.IsNegative() {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "claim must cannot be nil or negative")
	}

	// 3.2 sanity check 2: claimTotal <= deposit
	if stream.Deposit.IsLT(claimTotal) {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "not enough deposit to claim")
	}

	// 4. calculate validator fee and deduct from claim amount
	receiverAmount, valFee := types.CalculateValidatorFee(params.ValidatorFee, claimTotal)

	if valFee.Amount.GT(sdk.NewIntFromUint64(0)) {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(valFee))

		if err != nil {
			return sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, err
		}
	}

	// 5. send modified amount from module account to receiver
	if receiverAmount.Amount.GT(sdk.NewIntFromUint64(0)) {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddr, sdk.NewCoins(receiverAmount))

		if err != nil {
			return sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, err
		}
	}

	// 6. update & save stream
	stream.Deposit = remainingDeposit
	stream.LastOutflowTime = nowTime
	err := k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, sdk.Coin{}, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimStreamAction,
			sdk.NewAttribute(types.AttributeKeyStreamSender, senderAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStreamReceiver, receiverAddr.String()),
			sdk.NewAttribute(types.AttributeKeyClaimTotal, claimTotal.String()),
			sdk.NewAttribute(types.AttributeKeyClaimAmountReceived, receiverAmount.String()),
			sdk.NewAttribute(types.AttributeKeyClaimValidatorFee, valFee.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingDeposit, remainingDeposit.String()),
		),
	)

	return receiverAmount, valFee, claimTotal, remainingDeposit, nil
}

func (k Keeper) AddDeposit(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress, topUpDeposit sdk.Coin) (bool, error) {

	stream, ok := k.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return false, sdkerrors.Wrapf(types.ErrStreamDoesNotExist, "sender: %s, receiver %s", senderAddr.String(), receiverAddr.String())
	}

	nowTime := ctx.BlockTime()
	// calculate duration and deposit time to zero extension
	durationExtension := types.CalculateDuration(topUpDeposit, stream.FlowRate)
	// should be added to the current deposit zero time if the stream has not expired, or from
	// "now" if it has.
	var depositZeroTime time.Time

	if stream.DepositZeroTime.Before(nowTime) || stream.DepositZeroTime.Equal(nowTime) {
		// In the case of expired, ClaimFromStream is called first to "reset" deposit to zero and forward
		// remaining payment to the receiver wallet, effectively creating a new stream
		if stream.Deposit.Amount.GT(sdk.NewIntFromUint64(0)) {
			// only if stream has deposit
			_, _, _, _, err := k.ClaimFromStream(ctx, receiverAddr, senderAddr)
			if err != nil {
				return false, err
			}
			// refresh stream data, since deposits and total streamed may have changed
			// after claim stream call
			stream, _ = k.GetStream(ctx, receiverAddr, senderAddr)
		}

		// stream expired or new. Calculate from now
		depositZeroTime = nowTime.Add(time.Second * time.Duration(durationExtension))
	} else {
		// stream not expired. Add to current deposit zero time
		depositZeroTime = stream.DepositZeroTime.Add(time.Second * time.Duration(durationExtension))
	}

	// Send topUpDeposit from user acc to module acc
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, sdk.NewCoins(topUpDeposit))

	if err != nil {
		return false, err
	}

	// set and save new stream data
	// add topUpDeposit to current stream deposit
	newDeposit := stream.Deposit.Add(topUpDeposit) // may have been refreshed above for expired streams
	stream.Deposit = newDeposit
	stream.DepositZeroTime = depositZeroTime

	err = k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return false, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositToStream,
			sdk.NewAttribute(types.AttributeKeyStreamSender, senderAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStreamReceiver, receiverAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmountDeposited, topUpDeposit.String()),
			sdk.NewAttribute(types.AttributeKeyDepositDuration, strconv.FormatInt(durationExtension, 10)),
			sdk.NewAttribute(types.AttributeKeyDepositZeroTime, depositZeroTime.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingDeposit, newDeposit.String()),
		),
	)

	return true, nil
}

func (k Keeper) SetNewFlowRate(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress, newFlowRate int64) error {
	stream, ok := k.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return sdkerrors.Wrapf(types.ErrStreamDoesNotExist, "sender: %s, receiver %s", senderAddr.String(), receiverAddr.String())
	}

	// for event emission
	oldFlowRate := stream.FlowRate

	// default to now, and 0 duration. If there is no remaining deposit, then
	// updating the flow won't have any effect on the deposit zero time, since there
	// are no outstanding payments in the stream
	nowTime := ctx.BlockTime()
	depositZeroTime := nowTime
	duration := int64(0)

	// Check if the stream still has deposit value.
	if stream.Deposit.Amount.GT(sdk.NewIntFromUint64(0)) {
		// still has deposit. Claim unpaid deposits with the old flow rate first
		_, _, _, _, err := k.ClaimFromStream(ctx, receiverAddr, senderAddr)
		if err != nil {
			return err
		}

		// refresh stream data
		stream, _ = k.GetStream(ctx, receiverAddr, senderAddr)

		// Calculate new duration & deposit zero time based on new flow rate & remaining deposit.
		// Calculation is from "now", since the Claim function has been called
		// above. We're effectively creating a "new" stream, based on existing deposit value
		// and the new flow rate
		duration = types.CalculateDuration(stream.Deposit, newFlowRate)
		depositZeroTime = nowTime.Add(time.Second * time.Duration(duration))
	}

	// save new stream data
	stream.FlowRate = newFlowRate
	stream.DepositZeroTime = depositZeroTime

	err := k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateFlowRate,
			sdk.NewAttribute(types.AttributeKeyStreamSender, senderAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStreamReceiver, receiverAddr.String()),
			sdk.NewAttribute(types.AttributeKeyOldFlowRate, strconv.FormatInt(oldFlowRate, 10)),
			sdk.NewAttribute(types.AttributeKeyNewFlowRate, strconv.FormatInt(newFlowRate, 10)),
			sdk.NewAttribute(types.AttributeKeyDepositDuration, strconv.FormatInt(duration, 10)),
			sdk.NewAttribute(types.AttributeKeyDepositZeroTime, depositZeroTime.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingDeposit, stream.Deposit.String()),
		),
	)

	return nil
}

func (k Keeper) CancelStreamBySenderReceiver(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress) error {

	stream, ok := k.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return sdkerrors.Wrapf(types.ErrStreamDoesNotExist, "sender: %s, receiver %s", senderAddr.String(), receiverAddr.String())
	}

	// claim any outstanding flow
	if stream.Deposit.Amount.GT(sdk.NewIntFromUint64(0)) {
		_, _, _, _, err := k.ClaimFromStream(ctx, receiverAddr, senderAddr)
		if err != nil {
			return err
		}
		// refresh stream data
		stream, _ = k.GetStream(ctx, receiverAddr, senderAddr)
	}

	refundCoin := stream.Deposit
	// return any existing deposit to the sender
	if refundCoin.Amount.GT(sdk.NewIntFromUint64(0)) {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, sdk.NewCoins(refundCoin))
		if err != nil {
			return err
		}
	}

	// set all to zero etc.
	// ToDo - delete instead of set to zero
	nowTime := ctx.BlockTime()
	stream.Deposit = sdk.NewCoin(refundCoin.Denom, sdk.NewInt(0))
	stream.FlowRate = 0
	stream.DepositZeroTime = nowTime

	err := k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStreamCancelled,
			sdk.NewAttribute(types.AttributeKeyStreamSender, senderAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStreamReceiver, receiverAddr.String()),
			sdk.NewAttribute(types.AttributeKeyRefundAmount, refundCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingDeposit, stream.Deposit.String()),
		),
	)

	return nil
}

// CreateNewStream creates a new "empty" stream for a sender/receiver pair.
// Deposit and Deposit Zero Time are handled by the AddDeposit function.
// The value passed in the deposit var is only used to determine the denomination of the deposit.
func (k Keeper) CreateNewStream(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress, deposit sdk.Coin, flowRate int64) (types.Stream, error) {

	if k.IsStream(ctx, receiverAddr, senderAddr) {
		return types.Stream{}, sdkerrors.Wrap(types.ErrStreamExists, "stream exists")
	}

	nowTime := ctx.BlockTime()

	stream := types.Stream{
		Deposit:         sdk.NewCoin(deposit.Denom, sdk.NewInt(0)), // set to zero for correct calculation in AddDeposit
		FlowRate:        flowRate,
		LastOutflowTime: nowTime,
		DepositZeroTime: time.Unix(0, 0).UTC(), // set to past, so deposit zero time correctly calculated in AddDeposit
		Cancellable:     true,                  // default to true for now. Eventually, using eFUND will set to false
	}

	err := k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return types.Stream{}, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateStreamAction,
			sdk.NewAttribute(types.AttributeKeyStreamSender, senderAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStreamReceiver, receiverAddr.String()),
			sdk.NewAttribute(types.AttributeKeyFlowRate, strconv.FormatInt(flowRate, 10)),
		),
	)

	return stream, nil
}
