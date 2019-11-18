package simulation

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/unification-com/mainchain-cosmos/simapp/helpers"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// SimulateMsgRegisterWrkChain generates a MsgRegisterWrkChain with random values
// nolint: funlen
func SimulateMsgRegisterWrkChain(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		randSuffix := uint64(r.Intn(1000000))
		moniker := "wrkchain_" + fmt.Sprint(randSuffix)
		name := "WRKChain " + fmt.Sprint(randSuffix)

		fees := types.FeesWrkChainRegistration

		coins := account.SpendableCoins(ctx.BlockTime())

		coins, hasNeg := coins.SafeSub(fees)
		if hasNeg {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgRegisterWrkChain(
			moniker,
			"genesishash",
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


// SimulateMsgRecordWrkChainBlock generates a MsgRecordWrkChainBlock with random values
// nolint: funlen
func SimulateMsgRecordWrkChainBlock(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		params := types.NewQueryWrkChainParams(1, 100, "", simAccount.Address)
		registeredWrkChains := k.GetWrkChainsFiltered(ctx, params)

		if len(registeredWrkChains) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// get a random wrkchain
		rWc := 0
		if len(registeredWrkChains) > 1 {
			rWc = rand.Intn(len(registeredWrkChains) - 1)
		}

		wrkChain := registeredWrkChains[rWc]
		height := wrkChain.LastBlock + 1

		msg := types.NewMsgRecordWrkChainBlock(
			wrkChain.WrkChainID,
			height,
			"blockhash",
			"parenthash",
			"hash1",
			"hash2",
			"hash3",
			simAccount.Address,
		)

		fees := types.FeesWrkChainRecordHash

		coins := account.SpendableCoins(ctx.BlockTime())
		coins, hasNeg := coins.SafeSub(fees)
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
