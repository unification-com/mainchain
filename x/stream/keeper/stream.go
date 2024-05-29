package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/stream/types"
	"strconv"
	"time"
)

// GetHighestStreamId gets the highest BEACON ID
func (k Keeper) GetHighestStreamId(ctx sdk.Context) (StreamId uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestStreamIdKey)
	if bz == nil {
		return 1, nil
	}
	// convert from bytes to uint64
	StreamId = types.GetStreamIdFromBytes(bz)
	return StreamId, nil
}

// SetHighestStreamId sets the new highest BEACON ID to the store
func (k Keeper) SetHighestStreamId(ctx sdk.Context, StreamId uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert from uint64 to bytes for storage
	StreamIdbz := types.GetStreamIdBytes(StreamId)
	store.Set(types.HighestStreamIdKey, StreamIdbz)
}

// GetTotalDeposits gets the total deposits
func (k Keeper) GetTotalDeposits(ctx sdk.Context) (types.TotalDeposits, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TotalDepositsKey)
	if bz == nil {
		return types.TotalDeposits{}, false
	}
	var totalDeposits types.TotalDeposits
	k.cdc.MustUnmarshal(bz, &totalDeposits)
	return totalDeposits, true
}

// SetTotalDeposits sets the total deposits
func (k Keeper) SetTotalDeposits(ctx sdk.Context, totalDeposits types.TotalDeposits) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TotalDepositsKey, k.cdc.MustMarshal(&totalDeposits))
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

// SetUuidLookup Sets the uuid lookup
func (k Keeper) SetUuidLookup(ctx sdk.Context, streamId uint64, idLookup types.StreamIdLookup) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetStreamIdLookupKey(streamId), k.cdc.MustMarshal(&idLookup))
	return nil
}

// GetIdLookup Gets the uuid lookup
func (k Keeper) GetIdLookup(ctx sdk.Context, streamId uint64) (types.StreamIdLookup, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetStreamIdLookupKey(streamId))

	if bz == nil {
		// return a new empty stream lookup struct
		return types.StreamIdLookup{}, false
	}

	var idLookup types.StreamIdLookup
	k.cdc.MustUnmarshal(bz, &idLookup)
	return idLookup, true
}

func (k Keeper) ClaimFromStream(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress) (sdk.Coin, sdk.Coin, error) {
	stream, ok := k.GetStream(ctx, receiverAddr, senderAddr)
	params := k.GetParams(ctx)

	if !ok {
		return sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "stream does not exist")
	}

	totalDeposits, ok := k.GetTotalDeposits(ctx)

	if !ok {
		return sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "total deposits does not exist")
	}

	// 1. check current stream deposit > 0
	if stream.Deposit.IsNil() || stream.Deposit.IsNegative() || stream.Deposit.IsZero() {
		return sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "stream deposit is zero")
	}

	// 2. calculate amount to claim
	nowTime := ctx.BlockTime()
	amountToClaim, remainingDepositValue := types.CalculateAmountToClaim(nowTime, stream.DepositZeroTime, stream.LastOutflowTime, stream.Deposit, stream.FlowRate)

	// 3. sanity check: amount <= deposit
	if stream.Deposit.IsLT(amountToClaim) {
		return sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidData, "not enough deposit to claim")
	}

	// 4. deduct amount from module's total deposits and set in keeper
	// ToDo - use SafeSub
	newTotalDeposits := totalDeposits.Total.Sub(amountToClaim)
	totalDeposits.Total = newTotalDeposits
	k.SetTotalDeposits(ctx, totalDeposits)

	// 5. calculate validator fee and deduct from claim amount
	finalClaimCoin, valFeeCoin := types.CalculateValidatorFee(params.ValidatorFee, amountToClaim)

	if valFeeCoin.Amount.GT(sdk.NewIntFromUint64(0)) {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(valFeeCoin))

		if err != nil {
			return sdk.Coin{}, sdk.Coin{}, err
		}
	}

	// 6. send modified amount from module account to receiver
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddr, sdk.NewCoins(finalClaimCoin))

	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	// 7. update & save stream
	stream.Deposit = remainingDepositValue
	stream.LastOutflowTime = nowTime
	err = k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimStreamAction,
			sdk.NewAttribute(types.AttributeKeyStreamId, strconv.FormatUint(stream.StreamId, 10)),
			sdk.NewAttribute(types.AttributeKeyStreamSender, stream.Sender),
			sdk.NewAttribute(types.AttributeKeyStreamReceiver, stream.Receiver),
			sdk.NewAttribute(types.AttributeKeyStreamClaimTotal, amountToClaim.String()),
			sdk.NewAttribute(types.AttributeKeyStreamClaimAmountReceived, finalClaimCoin.String()),
			sdk.NewAttribute(types.AttributeKeyStreamClaimValidatorFee, valFeeCoin.String()),
		),
	)

	return finalClaimCoin, valFeeCoin, nil
}

func (k Keeper) AddDeposit(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress, deposit sdk.Coin) (bool, error) {

	stream, ok := k.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return false, sdkerrors.Wrapf(types.ErrStreamDoesNotExist, "sender: %s, receiver %s", senderAddr.String(), receiverAddr.String())
	}

	duration := types.CalculateDuration(deposit, stream.FlowRate)

	// calculate if stream has "expired"
	var depositZeroTime time.Time
	nowTime := ctx.BlockTime()

	if stream.DepositZeroTime.Before(nowTime) {
		// stream has "expired" - the deposit end time is in the past.
		// Reset
		depositZeroTime = nowTime.Add(time.Second * time.Duration(duration))
	} else {
		// stream still "valid". Just extend the current deposit zero time
		depositZeroTime = stream.DepositZeroTime.Add(time.Second * time.Duration(duration))
	}

	// Send from user acc to module acc
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, sdk.NewCoins(stream.Deposit))

	if err != nil {
		return false, err
	}

	// set and save new stream data
	newDeposit := stream.Deposit
	stream.Deposit = newDeposit.Add(deposit)
	stream.DepositZeroTime = depositZeroTime
	stream.LastUpdatedTime = nowTime

	err = k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return false, err
	}

	// calculate total deposited in module
	totalDeposits, has := k.GetTotalDeposits(ctx)

	if !has {
		totalDeposits.Total = sdk.NewCoins(stream.Deposit)
	} else {
		totalDeposits.Total = totalDeposits.Total.Add(stream.Deposit)
	}

	k.SetTotalDeposits(ctx, totalDeposits)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositToStream,
			sdk.NewAttribute(types.AttributeKeyStreamId, strconv.FormatUint(stream.StreamId, 10)),
			sdk.NewAttribute(types.AttributeKeyStreamDepositAmount, deposit.String()),
			sdk.NewAttribute(types.AttributeKeyStreamDepositDuration, strconv.FormatInt(duration, 10)),
			sdk.NewAttribute(types.AttributeKeyStreamDepositZeroTime, strconv.FormatInt(depositZeroTime.Unix(), 10)),
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

	nowTime := ctx.BlockTime()
	depositZeroTime := nowTime
	duration := int64(0)

	if stream.Deposit.Amount.GT(sdk.NewIntFromUint64(0)) {
		// still has deposit. Calculate new deposit zero time based on new flow rate
		duration = types.CalculateDuration(stream.Deposit, newFlowRate)
		depositZeroTime = nowTime.Add(time.Second * time.Duration(duration))
	}

	// save new stream data
	stream.FlowRate = newFlowRate
	stream.DepositZeroTime = depositZeroTime
	stream.LastUpdatedTime = nowTime

	err := k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateFlowRate,
			sdk.NewAttribute(types.AttributeKeyStreamId, strconv.FormatUint(stream.StreamId, 10)),
			sdk.NewAttribute(types.AttributeKeyOldFlowRate, strconv.FormatInt(oldFlowRate, 10)),
			sdk.NewAttribute(types.AttributeKeyNewFlowRate, strconv.FormatInt(newFlowRate, 10)),
			sdk.NewAttribute(types.AttributeKeyStreamDepositDuration, strconv.FormatInt(duration, 10)),
			sdk.NewAttribute(types.AttributeKeyStreamDepositZeroTime, strconv.FormatInt(depositZeroTime.Unix(), 10)),
		),
	)

	return nil
}

func (k Keeper) CancelStreamBySenderReceiver(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress) error {

	stream, ok := k.GetStream(ctx, receiverAddr, senderAddr)

	if !ok {
		return sdkerrors.Wrapf(types.ErrStreamDoesNotExist, "sender: %s, receiver %s", senderAddr.String(), receiverAddr.String())
	}

	refundCoin := stream.Deposit
	// return any existing deposit to the sender
	if refundCoin.Amount.GT(sdk.NewIntFromUint64(0)) {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, sdk.NewCoins(refundCoin))
		if err != nil {
			return err
		}

		totalDeposits, _ := k.GetTotalDeposits(ctx)
		totalDeposits.Total = totalDeposits.Total.Sub(refundCoin)

		k.SetTotalDeposits(ctx, totalDeposits)
	}

	// set all to zero etc.
	nowTime := ctx.BlockTime()
	stream.Deposit = sdk.NewCoin(refundCoin.Denom, sdk.NewInt(0))
	stream.FlowRate = 0
	stream.DepositZeroTime = nowTime
	stream.LastUpdatedTime = nowTime

	err := k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStreamCancelled,
			sdk.NewAttribute(types.AttributeKeyStreamId, strconv.FormatUint(stream.StreamId, 10)),
		),
	)

	return nil
}

func (k Keeper) CreateIdLookup(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress, streamId uint64) error {

	idLookup := types.StreamIdLookup{
		Sender:   senderAddr.String(),
		Receiver: receiverAddr.String(),
	}

	err := k.SetUuidLookup(ctx, streamId, idLookup)

	if err != nil {
		return err
	}

	return nil
}

// CreateNewStream creates a new "empty" stream for a sender/receiver pair
// Deposit and Deposit Zero Time are handled by the AddDeposit function
func (k Keeper) CreateNewStream(ctx sdk.Context, receiverAddr, senderAddr sdk.AccAddress, deposit sdk.Coin, flowRate int64) (types.Stream, error) {

	streamId, err := k.GetHighestStreamId(ctx)

	if err != nil {
		return types.Stream{}, err
	}

	nowTime := ctx.BlockTime()

	stream := types.Stream{
		StreamId:        streamId,
		Sender:          senderAddr.String(),
		Receiver:        receiverAddr.String(),
		Deposit:         sdk.NewCoin(deposit.Denom, sdk.NewInt(0)), // set to zero for correct calculation in AddDeposit
		FlowRate:        flowRate,
		CreateTime:      nowTime,
		LastUpdatedTime: nowTime,
		LastOutflowTime: nowTime,
		DepositZeroTime: time.Unix(0, 0), // set to past, so deposit zero time correctly calculated in AddDeposit
		TotalStreamed:   sdk.NewCoin(deposit.Denom, sdk.NewInt(0)),
		Cancellable:     true, // default to true for now. Eventually, using eFUND will set to false
	}

	err = k.SetStream(ctx, receiverAddr, senderAddr, stream)

	if err != nil {
		return types.Stream{}, err
	}

	k.SetHighestStreamId(ctx, streamId+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateStreamAction,
			sdk.NewAttribute(types.AttributeKeyStreamId, strconv.FormatUint(streamId, 10)),
			sdk.NewAttribute(types.AttributeKeyStreamSender, senderAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStreamReceiver, receiverAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStreamFlowRate, strconv.FormatInt(flowRate, 10)),
		),
	)

	return stream, nil
}
