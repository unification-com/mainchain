package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v4 "github.com/unification-com/mainchain/x/wrkchain/migrations/v4"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(k Keeper) Migrator {
	return Migrator{
		keeper: k,
	}
}

func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.Migrate(ctx, ctx.KVStore(m.keeper.storeKey), m.keeper.cdc)
}
