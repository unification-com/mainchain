package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

// RegisterInvariants registers all enterprise invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(keeper))
}

// AllInvariants runs all invariants of the enterprise module
func AllInvariants(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return ModuleAccountInvariant(keeper)(ctx)
	}
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// locked UND held on store
func ModuleAccountInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		//var lockedByAccount sdk.Coins
		totalLocked := sdk.NewCoins(keeper.GetTotalLockedUnd(ctx))

		//keeper.GetAllLockedUndAccountsIterator(ctx, func(lockedUnd types.Deposit) bool {
		//	expectedDeposits = expectedDeposits.Add(deposit.Amount)
		//	return false
		//})

		macc := keeper.GetEnterpriseAccount(ctx)
		broken := !macc.GetCoins().IsEqual(totalLocked)

		return sdk.FormatInvariant(types.ModuleName, "locked",
			fmt.Sprintf("\tenterprise ModuleAccount coins: %s\n\tsum of locked amounts:  %s\n",
				macc.GetCoins(), totalLocked)), broken
	}
}

