package simulation

import (
	"math/rand"

	mathmod "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/unification-com/mainchain/x/stream/keeper"
	"github.com/unification-com/mainchain/x/stream/types"
)

// Simulation operation weights constants
//
//nolint:gosec // These aren't harcoded credentials.
const (
	OpWeightMsgCreateStream            = "op_weight_msg_create_stream"
	OpWeightMsgClaimStream             = "op_weight_msg_claim_stream"
	OpWeightMsgTopUpDeposit            = "op_weight_msg_top_up_deposit"
	OpWeightMsgUpdateFlowRate          = "op_weight_msg_update_flow_rate"
	OpWeightMsgCancelStream            = "op_weight_msg_cancel_stream"
	DefaultWeightMsgCreateStream   int = 100
	DefaultWeightMsgClaimStream    int = 100
	DefaultWeightMsgTopUpDeposit   int = 100
	DefaultWeightMsgUpdateFlowRate int = 100
	DefaultWeightMsgCancelStream   int = 50
)

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec,
	k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper,
) simulation.WeightedOperations {

	var (
		weightMsgCreateStream   int
		weightMsgClaimStream    int
		weightMsgTopUpDeposit   int
		weightMsgUpdateFlowRate int
		weightMsgCancelStream   int
	)

	appParams.GetOrGenerate(OpWeightMsgCreateStream, &weightMsgCreateStream, nil,
		func(_ *rand.Rand) {
			weightMsgCreateStream = DefaultWeightMsgCreateStream
		},
	)

	appParams.GetOrGenerate(OpWeightMsgClaimStream, &weightMsgClaimStream, nil,
		func(_ *rand.Rand) {
			weightMsgClaimStream = DefaultWeightMsgClaimStream
		},
	)

	appParams.GetOrGenerate(OpWeightMsgTopUpDeposit, &weightMsgTopUpDeposit, nil,
		func(_ *rand.Rand) {
			weightMsgTopUpDeposit = DefaultWeightMsgTopUpDeposit
		},
	)

	appParams.GetOrGenerate(OpWeightMsgUpdateFlowRate, &weightMsgUpdateFlowRate, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateFlowRate = DefaultWeightMsgUpdateFlowRate
		},
	)

	appParams.GetOrGenerate(OpWeightMsgCancelStream, &weightMsgCancelStream, nil,
		func(_ *rand.Rand) {
			weightMsgCancelStream = DefaultWeightMsgCancelStream
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateStream,
			SimulateMsgCreateStream(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateStream,
			SimulateMsgClaimStream(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateStream,
			SimulateMsgTopUpDeposit(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateStream,
			SimulateMsgUpdateFlowRate(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateStream,
			SimulateMsgCancelStream(ak, bk, k),
		),
	}
}

// SimulateMsgCreateStream generates MsgCreateStream with random values.
func SimulateMsgCreateStream(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		sender, _ := simtypes.RandomAcc(r, accs)
		receiver, _ := simtypes.RandomAcc(r, accs)
		if sender.Address.String() == receiver.Address.String() {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, "sender and receiver cannot be same"), nil, nil
		}

		if k.IsStream(ctx, receiver.Address, sender.Address) {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, "stream exists"), nil, nil
		}

		simAccount, _ := simtypes.FindAccount(accs, sender.Address)
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, "account private key is nil"), nil, nil // skip
		}

		senderAccount := ak.GetAccount(ctx, sender.Address)

		// ToDO - Spendable... no longer implemented in x/bank
		spendable := bk.SpendableCoin(ctx, senderAccount.GetAddress(), sdk.DefaultBondDenom)

		depositAmnt, err := simtypes.RandPositiveInt(r, spendable.Amount)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, err.Error()), nil, nil
		}

		if depositAmnt.LT(mathmod.NewIntFromUint64(60)) {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, "depositAmnt too small"), nil, nil
		}

		deposit := sdk.NewCoin(sdk.DefaultBondDenom, depositAmnt)

		maxFlowRate := deposit.Amount.Quo(mathmod.NewIntFromUint64(60))

		if maxFlowRate.LTE(mathmod.NewIntFromUint64(1)) {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, "maxFlowRate too low"), nil, nil
		}

		randFowRate := int64(simtypes.RandIntBetween(r, 1, int(maxFlowRate.Uint64())))

		msg := types.NewMsgCreateStream(deposit, randFowRate, receiver.Address, sender.Address)

		// fees need to be calculated from the remaining spendable coins after deposit is subtracted, so
		// GenSignedMockTx is used instead of GenAndDeliverTxWithRandFees
		var fees sdk.Coins
		var feeErr error
		remainingCoins, err := spendable.SafeSub(msg.Deposit)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, err.Error()), nil, nil
		}

		fees, feeErr = simtypes.RandomFees(r, ctx, sdk.NewCoins(remainingCoins))
		if feeErr != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.CreateStreamAction, "not enough for fees"), nil, nil
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig

		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
			[]sdk.Msg{msg},
			fees,
			simtestutil.DefaultGenTxGas,
			chainID,
			[]uint64{senderAccount.GetAccountNumber()},
			[]uint64{senderAccount.GetSequence()},
			simAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgClaimStream generates MsgClaimStream with random values.
func SimulateMsgClaimStream(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		var sender simtypes.Account
		var receiver simtypes.Account
		maxTries := 2000
		haveStream := false

		for i := 0; i < maxTries; i += 1 {
			sen, _ := simtypes.RandomAcc(r, accs)
			rec, _ := simtypes.RandomAcc(r, accs)
			if k.IsStream(ctx, rec.Address, sen.Address) {
				haveStream = true
				sender = sen
				receiver = rec
				break
			}
		}

		if !haveStream {
			return simtypes.NoOpMsg(types.ModuleName, types.ClaimStreamAction, "suitable stream not found"), nil, nil
		}

		stream, _ := k.GetStream(ctx, receiver.Address, sender.Address)

		nowTime := ctx.BlockTime()
		claimTotal, _ := types.CalculateAmountToClaim(nowTime, stream.DepositZeroTime, stream.LastOutflowTime, stream.Deposit, stream.FlowRate)

		if claimTotal.IsZero() || claimTotal.IsNegative() || claimTotal.IsNil() {
			return simtypes.NoOpMsg(types.ModuleName, types.ClaimStreamAction, "nothing to claim"), nil, nil
		}

		simAccount, _ := simtypes.FindAccount(accs, receiver.Address)
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.ClaimStreamAction, "account private key is nil"), nil, nil // skip
		}

		msg := types.NewMsgClaimStream(receiver.Address, sender.Address)

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgTopUpDeposit generates MsgTopUpDeposit with random values.
func SimulateMsgTopUpDeposit(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		var sender simtypes.Account
		var receiver simtypes.Account
		maxTries := 2000
		haveStream := false

		for i := 0; i < maxTries; i += 1 {
			sen, _ := simtypes.RandomAcc(r, accs)
			rec, _ := simtypes.RandomAcc(r, accs)
			if k.IsStream(ctx, rec.Address, sen.Address) {
				haveStream = true
				sender = sen
				receiver = rec
				break
			}
		}

		if !haveStream {
			return simtypes.NoOpMsg(types.ModuleName, types.TopUpDepositAction, "suitable stream not found"), nil, nil
		}

		simAccount, _ := simtypes.FindAccount(accs, sender.Address)
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TopUpDepositAction, "account private key is nil"), nil, nil // skip
		}

		senderAccount := ak.GetAccount(ctx, sender.Address)

		spendable := bk.SpendableCoin(ctx, senderAccount.GetAddress(), sdk.DefaultBondDenom)

		depositAmnt, err := simtypes.RandPositiveInt(r, spendable.Amount)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TopUpDepositAction, err.Error()), nil, nil
		}

		if depositAmnt.LT(mathmod.NewIntFromUint64(60)) {
			return simtypes.NoOpMsg(types.ModuleName, types.TopUpDepositAction, "depositAmnt too small"), nil, nil
		}

		deposit := sdk.NewCoin(sdk.DefaultBondDenom, depositAmnt)

		msg := types.NewMsgTopUpDeposit(receiver.Address, sender.Address, deposit)

		// fees
		var fees sdk.Coins
		var feeErr error
		remainingCoins, err := spendable.SafeSub(msg.Deposit)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TopUpDepositAction, err.Error()), nil, nil
		}

		fees, feeErr = simtypes.RandomFees(r, ctx, sdk.NewCoins(remainingCoins))
		if feeErr != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TopUpDepositAction, "not enough for fees"), nil, nil
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig

		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
			[]sdk.Msg{msg},
			fees,
			simtestutil.DefaultGenTxGas,
			chainID,
			[]uint64{senderAccount.GetAccountNumber()},
			[]uint64{senderAccount.GetSequence()},
			simAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgUpdateFlowRate generates MsgUpdateFlowRate with random values.
func SimulateMsgUpdateFlowRate(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		var sender simtypes.Account
		var receiver simtypes.Account
		maxTries := 2000
		haveStream := false

		for i := 0; i < maxTries; i += 1 {
			sen, _ := simtypes.RandomAcc(r, accs)
			rec, _ := simtypes.RandomAcc(r, accs)
			if k.IsStream(ctx, rec.Address, sen.Address) {
				haveStream = true
				sender = sen
				receiver = rec
				break
			}
		}

		if !haveStream {
			return simtypes.NoOpMsg(types.ModuleName, types.UpdateFlowRateAction, "suitable stream not found"), nil, nil
		}

		stream, _ := k.GetStream(ctx, receiver.Address, sender.Address)

		simAccount, _ := simtypes.FindAccount(accs, sender.Address)
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.UpdateFlowRateAction, "account private key is nil"), nil, nil // skip
		}

		halfFlow := stream.FlowRate / 2
		newFlow := stream.FlowRate
		rnd := simtypes.RandIntBetween(r, 1, 2)

		if rnd == 1 {
			newFlow = stream.FlowRate - halfFlow
		} else {
			newFlow = stream.FlowRate + halfFlow
		}

		// unlikely but just in case
		if newFlow <= 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.UpdateFlowRateAction, "new flow must be greater than zero"), nil, nil
		}

		msg := types.NewMsgUpdateFlowRate(receiver.Address, sender.Address, newFlow)

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgCancelStream generates MsgCancelStream with random values.
func SimulateMsgCancelStream(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		var sender simtypes.Account
		var receiver simtypes.Account
		maxTries := 2000
		haveStream := false

		for i := 0; i < maxTries; i += 1 {
			sen, _ := simtypes.RandomAcc(r, accs)
			rec, _ := simtypes.RandomAcc(r, accs)
			if k.IsStream(ctx, rec.Address, sen.Address) {
				haveStream = true
				sender = sen
				receiver = rec
				break
			}
		}

		if !haveStream {
			return simtypes.NoOpMsg(types.ModuleName, types.CancelStreamAction, "suitable stream not found"), nil, nil
		}

		simAccount, _ := simtypes.FindAccount(accs, sender.Address)
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.CancelStreamAction, "account private key is nil"), nil, nil // skip
		}

		msg := types.NewMsgCancelStream(receiver.Address, sender.Address)

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
