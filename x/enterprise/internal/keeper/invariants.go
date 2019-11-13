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

		totalLocked := sdk.NewCoins(keeper.GetTotalLockedUnd(ctx))
		lockedByAccount := sdk.NewInt64Coin(types.DefaultDenomination, 0)

		lockedIterator := keeper.GetAllLockedUndAccountsIterator(ctx)
		for ; lockedIterator.Valid(); lockedIterator.Next() {
			var l types.LockedUnd
			keeper.cdc.MustUnmarshalBinaryBare(lockedIterator.Value(), &l)
			lockedByAccount = lockedByAccount.Add(l.Amount)
		}

		macc := keeper.GetEnterpriseAccount(ctx)
		broken := !macc.GetCoins().IsEqual(totalLocked) || !macc.GetCoins().IsEqual(sdk.NewCoins(lockedByAccount)) || !totalLocked.IsEqual(sdk.NewCoins(lockedByAccount))

		ctx.Logger().Info("ModuleAccountInvariant - enterprise", "broken", broken, "macc", macc.GetCoins(), "totalLocked", totalLocked, "lockedByAccount", lockedByAccount )

		return sdk.FormatInvariant(types.ModuleName, "locked",
			fmt.Sprintf("\tenterprise ModuleAccount coins: %s\n\ttotal locked: %s\n\t sum of locked: %s\n",
				macc.GetCoins(), totalLocked, lockedByAccount)), broken
	}
}

