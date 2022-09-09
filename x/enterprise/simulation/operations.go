package simulation

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simparams "github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"math/rand"
	"strings"
)

// Simulation operation weights constants
const (
	OpWeightMsgUndPurchaseOrder        = "op_weight_msg_raise_ent_po"
	OpWeightMsgProcessUndPurchaseOrder = "op_weight_msg_proc_ent_po"
	OpWeightMsgWhitelistAddress        = "op_weight_msg_ent_whitelist"

	DefaultMsgUndPurchaseOrder        = 20
	DefaultMsgProcessUndPurchaseOrder = 20
	DefaultMsgWhitelistAddress        = 20
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec,
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
			weightMsgWhitelistAddress,
			SimulateMsgWhitelistAddress(k, bk, ak),
		),
		simulation.NewWeightedOperation(
			weightMsgUndPurchaseOrder,
			SimulateMsgUndPurchaseOrder(k, bk, ak),
		),
		simulation.NewWeightedOperation(
			weightMsgProcessUndPurchaseOrder,
			SimulateMsgProcessUndPurchaseOrder(k, bk, ak),
		),
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

		randAmt := int64(simtypes.RandIntBetween(r, 1000, 1000000))

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

		// submit the PO
		opMsg := simtypes.NewOperationMsg(msg, true, "", nil)

		poId, err := k.GetHighestPurchaseOrderID(ctx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate poId"), nil, err
		}

		// schedule POs for decisions
		// first, get ent signers
		entSignerList := k.GetParamEntSigners(ctx)
		entSignerArray := strings.Split(entSignerList, ",")

		// allow some POs to go stale
		blocksInFuture := int64(1)
		switch r.Intn(4) {
		case 0:
		case 1:
		case 2:
		default:
			break
		case 3:
			blocksInFuture = int64(5)
		}

		// generate future operations for decisions
		fops := make([]simtypes.FutureOperation, len(entSignerArray))
		for i := 0; i < len(entSignerArray); i++ {
			whenDecide := ctx.BlockHeader().Height + blocksInFuture
			signer, _ := sdk.AccAddressFromBech32(entSignerArray[i])
			fops[i] = simtypes.FutureOperation{
				BlockHeight: int(whenDecide),
				Op:          operationSimulateMsgProcessUndPurchaseOrder(k, bk, ak, signer, int64(poId-1)),
			}
		}

		return opMsg, fops, nil
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

		enSignerAddr, err := getRandomEntSignerAcc(r, k, ctx)

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

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil

	}
}

func operationSimulateMsgProcessUndPurchaseOrder(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper,
	signer sdk.AccAddress, poId int64) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		if poId < 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "not processing POID -1"), nil, nil
		}

		po, found := k.GetPurchaseOrder(ctx, uint64(poId))
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "purchase order not found"), nil, nil
		}

		enSignerAccount, found := simtypes.FindAccount(accs, signer)

		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, fmt.Sprintf("unable to find ent signer account %s", signer)), nil, errors.New("unable to find ent signer account")
		}

		if po.Status != types.StatusRaised {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "purchase order does not have raised status"), nil, nil
		}

		//if len(po.Decisions) > 0 {
		//	for _, d := range po.Decisions {
		//		if d.Signer == signer.String() {
		//			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, fmt.Sprintf("signer %s already decided %s for poID %d", signer, d.Decision, poId)), nil, nil
		//		}
		//	}
		//}

		account := ak.GetAccount(ctx, signer)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.ProcessAction, "unable to generate fees"), nil, err
		}

		decision := randomDecisionOption(r)

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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(),
					fmt.Sprintf("unable to deliver tx PoId %d, signer %s, decision %s. Err: %s", po.Id, signer, decision, err)),
				nil,
				nil
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func SimulateMsgProcessUndPurchaseOrder(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return operationSimulateMsgProcessUndPurchaseOrder(k, bk, ak, sdk.AccAddress{}, -1)
}

func getRandomEntSignerAcc(r *rand.Rand, k keeper.Keeper,
	ctx sdk.Context) (entSigner sdk.AccAddress, err error) {
	entSignerList := k.GetParamEntSigners(ctx)
	entSignerArray := strings.Split(entSignerList, ",")
	rndIdx := r.Intn(len(entSignerArray))
	return sdk.AccAddressFromBech32(entSignerArray[rndIdx])
}

// Pick a random voting option
func randomDecisionOption(r *rand.Rand) types.PurchaseOrderStatus {
	switch r.Intn(2) {
	case 0:
		return types.StatusRejected
	case 1:
		return types.StatusAccepted
	default:
		panic("invalid vote option")
	}
}
