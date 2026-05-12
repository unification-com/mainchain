package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/stream/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {

	moduleAcc := k.GetStreamModuleAccount(ctx)
	moduleHoldings := sdk.Coins{}
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	k.SetParams(ctx, genState.Params)

	// detect duplicate (sender, receiver, denom) triples in the genesis input.
	seen := make(map[string]struct{})

	for _, stream := range genState.Streams {

		senderAddr, err := sdk.AccAddressFromBech32(stream.Sender)
		if err != nil {
			panic(err)
		}

		receiverAddr, err := sdk.AccAddressFromBech32(stream.Receiver)
		if err != nil {
			panic(err)
		}

		denom := stream.Stream.Deposit.Denom

		key := stream.Sender + "|" + stream.Receiver + "|" + denom
		if _, dup := seen[key]; dup {
			panic(fmt.Sprintf("duplicate stream in genesis for (sender=%s, receiver=%s, denom=%s)", stream.Sender, stream.Receiver, denom))
		}
		seen[key] = struct{}{}

		err = k.SetStream(ctx, receiverAddr, senderAddr, denom, stream.Stream)

		if err != nil {
			panic(err)
		}

		moduleHoldings = moduleHoldings.Add(stream.Stream.Deposit)
	}

	balances := k.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		k.accKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	if !balances.Equal(moduleHoldings) {
		panic(fmt.Sprintf("stream module acc balance does not match the module holdings: %s != %s", balances, moduleHoldings))
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	var streams []types.StreamExport

	k.IterateAllStreams(ctx, func(receiverAddr, senderAddr sdk.AccAddress, _ string, stream types.Stream) bool {
		streams = append(streams,
			types.StreamExport{
				Receiver: receiverAddr.String(),
				Sender:   senderAddr.String(),
				Stream:   stream,
			},
		)
		return false
	})

	return types.NewGenesisState(streams, params)
}
