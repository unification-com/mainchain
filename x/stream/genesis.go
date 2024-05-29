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
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	
	k.SetParams(ctx, genState.Params)
	k.SetHighestStreamId(ctx, genState.StartingStreamId)
	k.SetTotalDeposits(ctx, genState.TotalDeposits)

	for _, stream := range genState.Streams {
		s := types.Stream{
			StreamId:        stream.StreamId,
			Sender:          stream.Sender,
			Receiver:        stream.Receiver,
			Deposit:         stream.Deposit,
			FlowRate:        stream.FlowRate,
			CreateTime:      stream.CreateTime,
			LastUpdatedTime: stream.LastUpdatedTime,
			LastOutflowTime: stream.LastOutflowTime,
			DepositZeroTime: stream.DepositZeroTime,
			TotalStreamed:   stream.TotalStreamed,
			Cancellable:     stream.Cancellable,
		}

		senderAddr, err := sdk.AccAddressFromBech32(stream.Sender)
		if err != nil {
			panic(err)
		}

		receiverAddr, err := sdk.AccAddressFromBech32(stream.Receiver)
		if err != nil {
			panic(err)
		}

		err = k.SetStream(ctx, receiverAddr, senderAddr, s)

		if err != nil {
			panic(err)
		}

		idl := types.StreamIdLookup{
			Sender:   stream.Sender,
			Receiver: stream.Receiver,
		}

		err = k.SetUuidLookup(ctx, stream.StreamId, idl)

		if err != nil {
			panic(err)
		}
	}

	moduleHoldings := genState.TotalDeposits.Total

	balances := bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		accountKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	if !balances.IsEqual(moduleHoldings) {
		panic(fmt.Sprintf("enterprise module balance does not match the module holdings: %s <-> %s", balances, moduleHoldings))
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// ToDo

	return genesis
}
