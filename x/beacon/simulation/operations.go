package simulation

import (
	"errors"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/unification-com/mainchain-cosmos/simapp/helpers"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
)

// SimulateMsgRegisterBeacon generates a MsgRegisterBeacon with random values
// nolint: funlen
func SimulateMsgRegisterBeacon(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		moniker := simulation.RandStringOfLength(r, 16)
		name := simulation.RandStringOfLength(r, 16)

		fees := k.GetRegistrationFeeAsCoins(ctx)

		coins := account.SpendableCoins(ctx.BlockTime())

		_, hasNeg := coins.SafeSub(fees)

		if hasNeg {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgRegisterBeacon(
			moniker,
			name,
			simAccount.Address,
		)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil

	}
}

// SimulateMsgRecordBeaconTimestamp generates a MsgRecordBeaconTimestamp with random values
// nolint: funlen
func SimulateMsgRecordBeaconTimestamp(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		// randomly select a WRKChain
		beacons := k.GetAllBeacons(ctx)
		if len(beacons) == 0 {
			// nothing registered yet
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		rB := 0
		if len(beacons) > 1 {
			rB = rand.Intn(len(beacons) - 1)
		}
		beacon := beacons[rB]

		ownerAddr := beacon.Owner
		var simAccount simulation.Account

		for _, ac := range accs {
			if ac.Address.String() == ownerAddr.String() {
				simAccount = ac
			}
		}

		account := ak.GetAccount(ctx, ownerAddr)

		msg := types.NewMsgRecordBeaconTimestamp(
			beacon.BeaconID,
			simulation.RandStringOfLength(r, 32),
			uint64(time.Now().Unix()),
			ownerAddr,
		)

		fees := k.GetRecordFeeAsCoins(ctx)

		coins := account.SpendableCoins(ctx.BlockTime())
		_, hasNeg := coins.SafeSub(fees)

		if hasNeg {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil

	}
}
