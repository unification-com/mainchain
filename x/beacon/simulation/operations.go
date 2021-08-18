package simulation

import (
	"errors"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simparams "github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/x/beacon/keeper"
	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	OpWeightMsgRegisterBeacon        = "op_weight_msg_register_beacon"
	OpWeightMsgRecordBeaconTimestamp = "op_weight_msg_record_beacon_timestamp"

	DefaultMsgRegisterBeacon        = 10
	DefaultMsgRecordBeaconTimestamp = 30
)

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler,
	k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper,
) simulation.WeightedOperations {

	var (
		weightMsgRegisterBeacon        int
		weightMsgRecordBeaconTimestamp int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRegisterBeacon, &weightMsgRegisterBeacon, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterBeacon = DefaultMsgRegisterBeacon
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRecordBeaconTimestamp, &weightMsgRecordBeaconTimestamp, nil,
		func(_ *rand.Rand) {
			weightMsgRecordBeaconTimestamp = DefaultMsgRecordBeaconTimestamp
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
	}

	return wEntOps

}

func SimulateMsgRegisterBeacon(k keeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		beaconOwnerAddr, _ := simtypes.RandomAcc(r, accs)

		simAccount, found := simtypes.FindAccount(accs, beaconOwnerAddr.Address)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.RegisterAction, "unable to find account"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, beaconOwnerAddr.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees := k.GetRegistrationFeeAsCoins(ctx)

		_, hasNeg := spendable.SafeSub(fees)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.RegisterAction, "not enough to pay beacon registration fee"), nil, nil // skip
		}

		moniker := simtypes.RandStringOfLength(r, 64)
		name := simtypes.RandStringOfLength(r, 128)


		msg := types.NewMsgRegisterBeacon(moniker, name, account.GetAddress())

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

		_, hasNeg := spendable.SafeSub(fees)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.RecordAction, "not enough to pay beacon record timestamp fee"), nil, nil // skip
		}

		hash := simtypes.RandStringOfLength(r, 64)

		msg := types.NewMsgRecordBeaconTimestamp(beacon.BeaconId, hash, uint64(ctx.BlockTime().Unix()), account.GetAddress())

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
		opMsg := simtypes.NewOperationMsg(msg, true, "")

		return opMsg, nil, nil
	}
}

func getRandomBeacon(r *rand.Rand, k keeper.Keeper,
	ctx sdk.Context) (beacon types.Beacon, err error){
	beacons := k.GetAllBeacons(ctx)
	if len(beacons) == 0 {
		return types.Beacon{}, errors.New("no beacons")
	}
	rndIdx := r.Intn(len(beacons))
	return beacons[rndIdx], nil
}
