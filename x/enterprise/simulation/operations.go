package simulation

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/unification-com/mainchain/simapp/helpers"
	"github.com/unification-com/mainchain/x/enterprise/internal/keeper"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// SimulateMsgRaisePurchaseOrder generates a MsgUndPurchaseOrder with random values
// nolint: funlen
func SimulateMsgRaisePurchaseOrder(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		coins := account.SpendableCoins(ctx.BlockTime())

		fees, err := simulation.RandomFees(r, ctx, coins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgUndPurchaseOrder(
			simAccount.Address,
			sdk.NewInt64Coin(k.GetParamDenom(ctx), int64(simulation.RandIntBetween(r, 100000000000, 1000000000000))),
		)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil

	}
}

func SimulateMsgProcessUndPurchaseOrder(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		params := types.NewQueryPurchaseOrdersParams(1, 100, types.StatusRaised, sdk.AccAddress{})
		raisedPos := k.GetPurchaseOrdersFiltered(ctx, params)

		// needs to be sent specifically by the designated Ent account
		entAcc := GenerateEntSourceSimAccount()
		account := ak.GetAccount(ctx, entAcc.Address)

		if len(raisedPos) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		rndPo := 0
		if len(raisedPos) > 1 {
			rndPo = rand.Intn(len(raisedPos) - 1)
		}

		po := raisedPos[rndPo]

		for _, d := range po.Decisions {
			if d.Signer.Equals(entAcc.Address) {
				// decision already made
				return simulation.NoOpMsg(types.ModuleName), nil, nil
			}
		}

		coins := account.SpendableCoins(ctx.BlockTime())

		fees, err := simulation.RandomFees(r, ctx, coins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgProcessUndPurchaseOrder(po.PurchaseOrderID, keeper.RandomDecision(), entAcc.Address)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			entAcc.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
