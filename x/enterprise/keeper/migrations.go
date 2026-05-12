package keeper

import (
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper        Keeper
	accountKeeper types.AccountKeeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, accountKeeper types.AccountKeeper) Migrator {
	return Migrator{
		keeper:        keeper,
		accountKeeper: accountKeeper,
	}
}
