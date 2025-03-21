package simulation

import (
	"errors"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/unification-com/mainchain/x/beacon/keeper"
	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	OpWeightMsgRegisterBeacon             = "op_weight_msg_register_beacon"
	OpWeightMsgRecordBeaconTimestamp      = "op_weight_msg_record_beacon_timestamp"
	OpWeightMsgPurchaseBeaconStateStorage = "op_weight_msg_beacon_purchase_storage"

	DefaultMsgRegisterBeacon             = 10
	DefaultMsgRecordBeaconTimestamp      = 30
	DefaultMsgPurchaseBeaconStateStorage = 5
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
		weightMsgRegisterBeacon             int
		weightMsgRecordBeaconTimestamp      int
		weightMsgPurchaseBeaconStateStorage int
	)

	appParams.GetOrGenerate(OpWeightMsgRegisterBeacon, &weightMsgRegisterBeacon, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterBeacon = DefaultMsgRegisterBeacon
		},
	)

	appParams.GetOrGenerate(OpWeightMsgRecordBeaconTimestamp, &weightMsgRecordBeaconTimestamp, nil,
		func(_ *rand.Rand) {
			weightMsgRecordBeaconTimestamp = DefaultMsgRecordBeaconTimestamp
		},
	)

	appParams.GetOrGenerate(OpWeightMsgPurchaseBeaconStateStorage, &weightMsgPurchaseBeaconStateStorage, nil,
		func(_ *rand.Rand) {
			weightMsgPurchaseBeaconStateStorage = DefaultMsgPurchaseBeaconStateStorage
		},
	)

	wEntOps := simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgRegisterBeacon,
			SimulateMsgRegisterBeacon(k, bk, ak),
		),
		simulation.NewWeightedOperation(
			weightMsgRecordBeaconTimestamp,
			SimulateMsgRecordBeaconTimestamp(k, bk, ak),
		),
		simulation.NewWeightedOperation(
			weightMsgPurchaseBeaconStateStorage,
			SimulateMsgPurchaseBeaconStateStorage(k, bk, ak),
		),
	}

	return wEntOps

}

func SimulateMsgRegisterBeacon(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
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
			return simtypes.NoOpMsg(types.ModuleName, types.RegisterAction, "not enough to pay beacon registration fee"), nil, nil // skip
		}

		moniker := simtypes.RandStringOfLength(r, 64)
		name := simtypes.RandStringOfLength(r, 128)

		msg := types.NewMsgRegisterBeacon(moniker, name, account.GetAddress())

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
		opMsg := simtypes.NewOperationMsg(msg, true, "")

		return opMsg, nil, nil
	}
}

func SimulateMsgRecordBeaconTimestamp(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		beacon, err := getRandomBeacon(r, k, ctx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "no beacons"), nil, nil // skip
		}

		beaconOwnerAddr, err := sdk.AccAddressFromBech32(beacon.Owner)

		simAccount, found := simtypes.FindAccount(accs, beaconOwnerAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "unable to find account"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, beaconOwnerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees := k.GetRecordFeeAsCoins(ctx)
		feeAsCoin := k.GetRecordFeeAsCoin(ctx)

		_, hasNeg := spendable.SafeSub(feeAsCoin)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "not enough to pay beacon record timestamp fee"), nil, nil // skip
		}

		hash := simtypes.RandStringOfLength(r, 64)

		msg := types.NewMsgRecordBeaconTimestamp(beacon.BeaconId, hash, uint64(ctx.BlockTime().Unix()), account.GetAddress())

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
		opMsg := simtypes.NewOperationMsg(msg, true, "")

		return opMsg, nil, nil
	}
}

func SimulateMsgPurchaseBeaconStateStorage(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		beacon, err := getRandomBeacon(r, k, ctx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseStorageAction, "no beacons"), nil, nil // skip
		}

		beaconOwnerAddr, err := sdk.AccAddressFromBech32(beacon.Owner)

		simAccount, found := simtypes.FindAccount(accs, beaconOwnerAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "unable to find account"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, beaconOwnerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		maxCanPurchase := k.GetMaxPurchasableSlots(ctx, beacon.BeaconId)
		if maxCanPurchase == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseStorageAction, "max storage reached"), nil, nil // skip
		}

		randNumToPurchase := uint64(1)
		if maxCanPurchase > 1 {
			randNumToPurchase = uint64(simtypes.RandIntBetween(r, 1, int(maxCanPurchase)))
		}

		bParams := k.GetParams(ctx)
		actualPurchaseAmt := bParams.FeePurchaseStorage
		actualFeeDenom := bParams.Denom

		feeInt := int64(actualPurchaseAmt * randNumToPurchase)
		feeAsCoin := sdk.NewInt64Coin(actualFeeDenom, feeInt)
		fees := sdk.NewCoins(feeAsCoin)

		_, hasNeg := spendable.SafeSub(feeAsCoin)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.PurchaseStorageAction, "not enough to pay beacon purchase storage fee"), nil, nil // skip
		}

		msg := types.NewMsgPurchaseBeaconStateStorage(beacon.BeaconId, uint64(randNumToPurchase), account.GetAddress())

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
		opMsg := simtypes.NewOperationMsg(msg, true, "")

		return opMsg, nil, nil
	}
}

func getRandomBeacon(r *rand.Rand, k keeper.Keeper,
	ctx sdk.Context) (beacon types.Beacon, err error) {
	beacons := k.GetAllBeacons(ctx)
	if len(beacons) == 0 {
		return types.Beacon{}, errors.New("no beacons")
	}
	rndIdx := r.Intn(len(beacons))
	return beacons[rndIdx], nil
}
