package simulation

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simparams "github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgUndPurchaseOrder        = "op_weight_msg_ent_raise_po"
	OpWeightMsgProcessUndPurchaseOrder = "op_weight_msg_process_po"
	OpWeightMsgWhitelistAddress        = "op_weight_msg_process_po"

	DefaultMsgUndPurchaseOrder        = 20
	DefaultMsgProcessUndPurchaseOrder = 20
	DefaultMsgWhitelistAddress        = 20
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler,
	k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper,
) simulation.WeightedOperations {

	var (
		weightMsgUndPurchaseOrder        int
		weightMsgProcessUndPurchaseOrder int
		weightMsgWhitelistAddress        int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgUndPurchaseOrder, &weightMsgUndPurchaseOrder, nil,
		func(_ *rand.Rand) {
			weightMsgUndPurchaseOrder = DefaultMsgUndPurchaseOrder
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgProcessUndPurchaseOrder, &weightMsgProcessUndPurchaseOrder, nil,
		func(_ *rand.Rand) {
			weightMsgProcessUndPurchaseOrder = DefaultMsgProcessUndPurchaseOrder
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgWhitelistAddress, &weightMsgWhitelistAddress, nil,
		func(_ *rand.Rand) {
			weightMsgWhitelistAddress = DefaultMsgWhitelistAddress
		},
	)

	wEntOps := simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUndPurchaseOrder,
			SimulateMsgUndPurchaseOrder(k, bk, ak),
		),
		simulation.NewWeightedOperation(
			weightMsgWhitelistAddress,
			SimulateMsgWhitelistAddress(k, bk, ak),
		),
		//simulation.NewWeightedOperation(
		//	weightMsgProcessUndPurchaseOrder,
		//	SimulateMsgProcessUndPurchaseOrder(k, bk, ak),
		//),
	}

	return wEntOps
}

func SimulateMsgUndPurchaseOrder(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		purchaserAddr, _ := simtypes.RandomAcc(r, accs)

		isWhitelisted := k.AddressIsWhitelisted(ctx, purchaserAddr.Address)

		if !isWhitelisted {
			return simtypes.NoOpMsg(types.ModuleName, types.WhitelistAddressAction, "address not whitelisted"), nil, nil
		}

		simAccount, found := simtypes.FindAccount(accs, purchaserAddr.Address)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseAction, "unable to find account"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, purchaserAddr.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseAction, "unable to generate fees"), nil, err
		}

		randAmt := int64(rand.Intn(1000000) + 1)

		msg := types.NewMsgUndPurchaseOrder(account.GetAddress(), sdk.NewInt64Coin(sdk.DefaultBondDenom, randAmt))

		txGen := simparams.MakeTestEncodingConfig().TxConfig

		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, res, err := app.Deliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, errors.New(res.Log)
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil

	}
}

func SimulateMsgWhitelistAddress(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		accToWhitelist, _ := simtypes.RandomAcc(r, accs)
		wlAction := types.WhitelistActionAdd

		isWhitelisted := k.AddressIsWhitelisted(ctx, accToWhitelist.Address)

		if isWhitelisted {
			wlAction = types.WhitelistActionRemove
		}

		enSignerAddr, err := sdk.AccAddressFromBech32(k.GetParamEntSigners(ctx))

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.WhitelistAddressAction, "unable to get addr from bech32"), nil, err
		}

		enSignerAccount, found := simtypes.FindAccount(accs, enSignerAddr)

		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.WhitelistAddressAction, fmt.Sprintf("unable to find ent signer account %s", enSignerAddr)), nil, errors.New("unable to find ent signer account")
		}

		account := ak.GetAccount(ctx, enSignerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.WhitelistAddressAction, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgWhitelistAddress(accToWhitelist.Address, wlAction, enSignerAccount.Address)

		txGen := simparams.MakeTestEncodingConfig().TxConfig

		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			enSignerAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, res, err := app.Deliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, errors.New(res.Log)
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil

	}
}

func SimulateMsgProcessUndPurchaseOrder(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		poId, found := randomPurchaseOrderId(r, k, ctx, types.StatusRaised)

		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "no raised purchase orders"), nil, nil
		}

		enSignerAddr, err := sdk.AccAddressFromBech32(k.GetParamEntSigners(ctx))

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "unable to get addr from bech32"), nil, err
		}

		enSignerAccount, found := simtypes.FindAccount(accs, enSignerAddr)

		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, fmt.Sprintf("unable to find ent signer account %s", enSignerAddr)), nil, errors.New("unable to find ent signer account")
		}

		po, found := k.GetPurchaseOrder(ctx, poId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "purchase order not found"), nil, nil
		}

		if len(po.Decisions) > 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "decision already made"), nil, nil
		}

		account := ak.GetAccount(ctx, enSignerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.WhitelistAddressAction, "unable to generate fees"), nil, err
		}

		rnd := rand.Intn(100)
		decision := types.StatusAccepted
		if rnd >= 50 {
			decision = types.StatusRejected
		}

		msg := types.NewMsgProcessUndPurchaseOrder(po.Id, decision, enSignerAccount.Address)

		txGen := simparams.MakeTestEncodingConfig().TxConfig

		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			enSignerAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, res, err := app.Deliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, errors.New(res.Log)
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil

	}
}

func randomPurchaseOrderId(r *rand.Rand, k keeper.Keeper,
	ctx sdk.Context, status types.PurchaseOrderStatus) (poId uint64, found bool) {

	poId, _ = k.GetHighestPurchaseOrderID(ctx)

	if poId > 1 {
		poId = uint64(simtypes.RandIntBetween(r, int(1), int(poId)))
	}

	po, found := k.GetPurchaseOrder(ctx, poId)
	if !found || po.Status != status {
		return poId, false
	}

	return poId, true
}
