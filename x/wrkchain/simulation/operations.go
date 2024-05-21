package simulation

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/unification-com/mainchain/x/wrkchain/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/types"
	"math/rand"
)

const (
	OpWeightMsgRegisterWrkChain             = "op_weight_msg_register_wrkchain"
	OpWeightMsgRecordWrkChainBlock          = "op_weight_msg_record_wrkchain_hash"
	OpWeightMsgPurchaseWrkChainStateStorage = "op_weight_msg_wrkchain_purchase_storage"

	DefaultMsgRegisterWrkChain             = 10
	DefaultMsgRecordWrkChainBlock          = 30
	DefaultMsgPurchaseWrkChainStateStorage = 5
)

//func WeightedOperations(
//	appParams simtypes.AppParams, cdc codec.JSONCodec,
//	k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper,
//) simulation.WeightedOperations {
//	return nil
//}

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec,
	k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper,
) simulation.WeightedOperations {

	var (
		weightMsgRegisterWrkChain             int
		weightMsgRecordWrkChainBlock          int
		weightMsgPurchaseWrkChainStateStorage int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRegisterWrkChain, &weightMsgRegisterWrkChain, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterWrkChain = DefaultMsgRegisterWrkChain
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRecordWrkChainBlock, &weightMsgRecordWrkChainBlock, nil,
		func(_ *rand.Rand) {
			weightMsgRecordWrkChainBlock = DefaultMsgRecordWrkChainBlock
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgPurchaseWrkChainStateStorage, &weightMsgPurchaseWrkChainStateStorage, nil,
		func(_ *rand.Rand) {
			weightMsgPurchaseWrkChainStateStorage = DefaultMsgPurchaseWrkChainStateStorage
		},
	)

	wEntOps := simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgRegisterWrkChain,
			SimulateMsgRegisterWrkChain(k, bk, ak),
		),
		simulation.NewWeightedOperation(
			weightMsgRecordWrkChainBlock,
			SimulateMsgRecordWrkChainBlock(k, bk, ak),
		),
		simulation.NewWeightedOperation(
			weightMsgPurchaseWrkChainStateStorage,
			SimulateMsgPurchaseWrkChainStateStorage(k, bk, ak),
		),
	}

	return wEntOps

}

func SimulateMsgRegisterWrkChain(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees := k.GetRegistrationFeeAsCoins(ctx)
		feeAsCoin := k.GetRegistrationFeeAsCoin(ctx)

		_, hasNeg := spendable.SafeSub(feeAsCoin)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.RegisterAction, "not enough to pay wrkchain registration fee"), nil, nil // skip
		}

		moniker := simtypes.RandStringOfLength(r, 64)
		genesisHash := simtypes.RandStringOfLength(r, 64)
		name := simtypes.RandStringOfLength(r, 128)
		baseType := simtypes.RandStringOfLength(r, 5)

		msg := types.NewMsgRegisterWrkChain(moniker, genesisHash, name, baseType, account.GetAddress())

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig

		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
			[]sdk.Msg{msg},
			fees,
			simtestutil.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		// submit the PO
		opMsg := simtypes.NewOperationMsg(msg, true, "", nil)

		return opMsg, nil, nil
	}
}

func SimulateMsgRecordWrkChainBlock(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		wrkChain, err := getRandomWrkChain(r, k, ctx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "no wrkchains"), nil, nil // skip
		}

		wrkChainOwnerAddr, err := sdk.AccAddressFromBech32(wrkChain.Owner)

		simAccount, found := simtypes.FindAccount(accs, wrkChainOwnerAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "unable to find account"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, wrkChainOwnerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees := k.GetRecordFeeAsCoins(ctx)
		feeAsCoin := k.GetRecordFeeAsCoin(ctx)

		_, hasNeg := spendable.SafeSub(feeAsCoin)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "not enough to pay wrkchain record timestamp fee"), nil, nil // skip
		}

		hash := simtypes.RandStringOfLength(r, 64)
		ph := ""
		h1 := ""
		h2 := ""
		h3 := ""

		if wrkChain.WrkchainId%2 == 0 {
			ph = simtypes.RandStringOfLength(r, 64)
			h1 = simtypes.RandStringOfLength(r, 64)
			h2 = simtypes.RandStringOfLength(r, 64)
			h3 = simtypes.RandStringOfLength(r, 64)
		}

		height := wrkChain.Lastblock + 1

		msg := types.NewMsgRecordWrkChainBlock(wrkChain.WrkchainId, height, hash, ph, h1, h2, h3, account.GetAddress())

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig

		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
			[]sdk.Msg{msg},
			fees,
			simtestutil.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		// submit the PO
		opMsg := simtypes.NewOperationMsg(msg, true, "", nil)

		return opMsg, nil, nil
	}
}

func SimulateMsgPurchaseWrkChainStateStorage(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		wrkchain, err := getRandomWrkChain(r, k, ctx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseStorageAction, "no wrkchains"), nil, nil // skip
		}

		wrkchainOwnerAddr, err := sdk.AccAddressFromBech32(wrkchain.Owner)

		simAccount, found := simtypes.FindAccount(accs, wrkchainOwnerAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "unable to find account"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, wrkchainOwnerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		maxCanPurchase := k.GetMaxPurchasableSlots(ctx, wrkchain.WrkchainId)
		if maxCanPurchase == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseStorageAction, "max storage reached"), nil, nil // skip
		}

		randNumToPurchase := uint64(1)
		if maxCanPurchase > 1 {
			randNumToPurchase = uint64(simtypes.RandIntBetween(r, 1, int(maxCanPurchase)))
		}

		wcParams := k.GetParams(ctx)
		actualPurchaseAmt := wcParams.FeePurchaseStorage
		actualFeeDenom := wcParams.Denom

		feeInt := int64(actualPurchaseAmt * randNumToPurchase)
		feeAsCoin := sdk.NewInt64Coin(actualFeeDenom, feeInt)
		fees := sdk.NewCoins(feeAsCoin)

		_, hasNeg := spendable.SafeSub(feeAsCoin)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseStorageAction, "not enough to pay wrkchain purchase storage fee"), nil, nil // skip
		}

		msg := types.NewMsgPurchaseWrkChainStateStorage(wrkchain.WrkchainId, randNumToPurchase, account.GetAddress())

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig

		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
			[]sdk.Msg{msg},
			fees,
			simtestutil.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		// submit the PO
		opMsg := simtypes.NewOperationMsg(msg, true, "", nil)

		return opMsg, nil, nil
	}
}

func getRandomWrkChain(r *rand.Rand, k keeper.Keeper,
	ctx sdk.Context) (wrkchain types.WrkChain, err error) {
	wrkChains := k.GetAllWrkChains(ctx)
	if len(wrkChains) == 0 {
		return types.WrkChain{}, errors.New("no wrkChains")
	}
	rndIdx := r.Intn(len(wrkChains))
	return wrkChains[rndIdx], nil
}
