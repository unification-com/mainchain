package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/stream/types"
)

// RegisterInvariants registers all stream invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(keeper))
}

// AllInvariants runs all invariants of the stream module
func AllInvariants(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return ModuleAccountInvariant(keeper)(ctx)
	}
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// total stream deposits held in the store
func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		totalDeposits := sdk.NewCoins()
		mAccBalance := k.GetStreamModuleAccountBalances(ctx)

		k.IterateAllStreams(ctx, func(receiverAddr, senderAddr sdk.AccAddress, stream types.Stream) bool {
			totalDeposits = totalDeposits.Add(stream.Deposit)
			return false
		})

		broken := !mAccBalance.Equal(totalDeposits)

		return sdk.FormatInvariant(types.ModuleName, "deposits",
			fmt.Sprintf("\tstream ModuleAccount coins: %s\n\ttotal deposits: %s\n",
				mAccBalance, totalDeposits)), broken
	}
}
