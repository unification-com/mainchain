package simulation

import (
	"context"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/simsx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/stream/keeper"
	"github.com/unification-com/mainchain/x/stream/types"
)

// streamRef is a stream identified by its (receiver, sender, denom) triple plus
// the loaded value. Used by the factories that need to find an existing stream.
type streamRef struct {
	receiver sdk.AccAddress
	sender   sdk.AccAddress
	denom    string
	stream   types.Stream
}

// pickStream returns a random existing stream from the keeper or skips the
// reporter when none exist. Returns false if skipped.
func pickStream(ctx sdk.Context, k keeper.Keeper, r *simsx.XRand, reporter simsx.SimulationReporter) (streamRef, bool) {
	var all []streamRef
	k.IterateAllStreams(ctx, func(receiver, sender sdk.AccAddress, denom string, s types.Stream) bool {
		all = append(all, streamRef{receiver: receiver, sender: sender, denom: denom, stream: s})
		return false
	})
	if len(all) == 0 {
		reporter.Skip("no streams")
		return streamRef{}, false
	}
	return all[r.Intn(len(all))], true
}

// MsgCreateStreamFactory picks two distinct accounts; sender funds a stream to
// receiver in nund. Deposit is drawn via the framework's RandSubsetCoin which
// tracks per-account allocations, avoiding over-commitment when the same account
// is used by other ops in the same block. Flow rate is randomised up to deposit/60.
func MsgCreateStreamFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgCreateStream] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgCreateStream) {
		r := testData.Rand()
		sender := testData.AnyAccount(reporter, simsx.WithSpendableBalance())
		if reporter.IsSkipped() {
			return nil, nil
		}
		receiver := testData.AnyAccount(reporter, simsx.ExcludeAccounts(sender))
		if reporter.IsSkipped() {
			return nil, nil
		}
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if k.IsStream(sdkCtx, receiver.Address, sender.Address, sdk.DefaultBondDenom) {
			reporter.Skip("stream already exists for pair")
			return nil, nil
		}
		deposit := sender.LiquidBalance().RandSubsetCoin(reporter, sdk.DefaultBondDenom)
		if reporter.IsSkipped() {
			return nil, nil
		}
		if deposit.Amount.LT(math.NewIntFromUint64(120)) {
			reporter.Skip("deposit too small")
			return nil, nil
		}
		maxFlow := deposit.Amount.Quo(math.NewIntFromUint64(60))
		if maxFlow.LTE(math.OneInt()) {
			reporter.Skip("maxFlow too low")
			return nil, nil
		}
		flowRate := int64(r.IntInRange(1, int(maxFlow.Uint64())))
		return []simsx.SimAccount{sender}, types.NewMsgCreateStream(deposit, flowRate, receiver.Address, sender.Address)
	}
}

func MsgClaimStreamFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgClaimStream] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgClaimStream) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		ref, ok := pickStream(sdkCtx, k, testData.Rand(), reporter)
		if !ok {
			return nil, nil
		}
		nowTime := sdkCtx.BlockTime()
		claimTotal, _ := types.CalculateAmountToClaim(nowTime, ref.stream.DepositZeroTime, ref.stream.LastOutflowTime, ref.stream.Deposit, ref.stream.FlowRate)
		if claimTotal.IsNil() || claimTotal.IsZero() || claimTotal.IsNegative() {
			reporter.Skip("nothing to claim")
			return nil, nil
		}
		receiver := testData.GetAccountbyAccAddr(reporter, ref.receiver)
		if reporter.IsSkipped() {
			return nil, nil
		}
		return []simsx.SimAccount{receiver}, types.NewMsgClaimStream(ref.receiver, ref.sender, ref.denom)
	}
}

func MsgTopUpDepositFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgTopUpDeposit] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgTopUpDeposit) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		ref, ok := pickStream(sdkCtx, k, testData.Rand(), reporter)
		if !ok {
			return nil, nil
		}
		sender := testData.GetAccountbyAccAddr(reporter, ref.sender)
		if reporter.IsSkipped() {
			return nil, nil
		}
		deposit := sender.LiquidBalance().RandSubsetCoin(reporter, ref.denom)
		if reporter.IsSkipped() {
			return nil, nil
		}
		if deposit.Amount.LT(math.NewIntFromUint64(60)) {
			reporter.Skip("top-up amount too small")
			return nil, nil
		}
		return []simsx.SimAccount{sender}, types.NewMsgTopUpDeposit(ref.receiver, ref.sender, deposit)
	}
}

func MsgUpdateFlowRateFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgUpdateFlowRate] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgUpdateFlowRate) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		r := testData.Rand()
		ref, ok := pickStream(sdkCtx, k, r, reporter)
		if !ok {
			return nil, nil
		}
		sender := testData.GetAccountbyAccAddr(reporter, ref.sender)
		if reporter.IsSkipped() {
			return nil, nil
		}
		halfFlow := ref.stream.FlowRate / 2
		newFlow := ref.stream.FlowRate
		if r.Intn(2) == 0 {
			newFlow -= halfFlow
		} else {
			newFlow += halfFlow
		}
		if newFlow <= 0 {
			reporter.Skip("new flow rate must be positive")
			return nil, nil
		}
		return []simsx.SimAccount{sender}, types.NewMsgUpdateFlowRate(ref.receiver, ref.sender, newFlow, ref.denom)
	}
}

func MsgCancelStreamFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgCancelStream] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgCancelStream) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		ref, ok := pickStream(sdkCtx, k, testData.Rand(), reporter)
		if !ok {
			return nil, nil
		}
		sender := testData.GetAccountbyAccAddr(reporter, ref.sender)
		if reporter.IsSkipped() {
			return nil, nil
		}
		return []simsx.SimAccount{sender}, types.NewMsgCancelStream(ref.receiver, ref.sender, ref.denom)
	}
}
