package stream

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/stream/keeper"
	"github.com/unification-com/mainchain/x/stream/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, genState types.GenesisState) {

	moduleAcc := k.GetStreamModuleAccount(ctx)
	moduleHoldings := sdk.Coins{}
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	k.SetParams(ctx, genState.Params)

	for _, stream := range genState.Streams {

		senderAddr, err := sdk.AccAddressFromBech32(stream.Sender)
		if err != nil {
			panic(err)
		}

		receiverAddr, err := sdk.AccAddressFromBech32(stream.Receiver)
		if err != nil {
			panic(err)
		}

		err = k.SetStream(ctx, receiverAddr, senderAddr, stream.Stream)

		if err != nil {
			panic(err)
		}

		moduleHoldings = moduleHoldings.Add(stream.Stream.Deposit)
	}

	balances := bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		accountKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	if !balances.IsEqual(moduleHoldings) {
		panic(fmt.Sprintf("stream module acc balance does not match the module holdings: %s != %s", balances, moduleHoldings))
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// ToDo

	return genesis
}
