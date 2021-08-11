package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
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
// locked FUND held on store
func ModuleAccountInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		totalLocked := sdk.NewCoins(keeper.GetTotalLockedUnd(ctx))

		lockedByAccount := sdk.NewInt64Coin(keeper.GetParamDenom(ctx), 0)

		lockedIterator := keeper.GetAllLockedUndAccountsIterator(ctx)
		for ; lockedIterator.Valid(); lockedIterator.Next() {
			var l types.LockedUnd
			keeper.cdc.MustUnmarshalBinaryBare(lockedIterator.Value(), &l)
			lockedByAccount = lockedByAccount.Add(l.Amount)
		}

		macc := keeper.GetEnterpriseAccount(ctx)
		maccCoins := keeper.GetCoins(ctx, macc.GetAddress())

		broken := !maccCoins.IsEqual(totalLocked) || !maccCoins.IsEqual(sdk.NewCoins(lockedByAccount)) || !totalLocked.IsEqual(sdk.NewCoins(lockedByAccount))

		logger := keeper.Logger(ctx)
		logger.Debug("ModuleAccountInvariant - enterprise", "is_broken", broken, "maccCoins", maccCoins, "totalLocked", totalLocked, "lockedByAccount", lockedByAccount)

		return sdk.FormatInvariant(types.ModuleName, "locked",
			fmt.Sprintf("\tenterprise ModuleAccount coins: %s\n\ttotal locked: %s\n\t sum of locked: %s\n",
				maccCoins, totalLocked, lockedByAccount)), broken
	}
}
