package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	undapp "github.com/unification-com/mainchain/app"
	apphelpers "github.com/unification-com/mainchain/app/helpers"
	streamtypes "github.com/unification-com/x-stream/x/stream/types"
)

// TestUpgradeNameIs8Vaxildan asserts that the registered upgrade name matches the
// on-chain plan name. Lets a stray "7-taryon" leftover fail loudly.
func TestUpgradeNameIs8Vaxildan(t *testing.T) {
	require.Equal(t, "8-vaxildan", undapp.UpgradeName)
}

// TestFreshBootVersionMap asserts that a fresh-boot App reports a non-empty consensus
// version map and that x/stream is at ConsensusVer = 2 (the v1→v2 bump that ships in
// the 8-vaxildan upgrade). The other modules' versions float with SDK upstreams; checking
// the stream module specifically is the integration anchor for stage 5b.
func TestFreshBootVersionMap(t *testing.T) {
	app := apphelpers.Setup(t)
	ctx := app.NewContext(false)

	versionMap := app.ModuleManager.GetVersionMap()
	require.NotEmpty(t, versionMap)

	streamVersion, ok := versionMap[streamtypes.ModuleName]
	require.True(t, ok, "stream module must appear in fresh-boot VersionMap")
	require.Equal(t, uint64(2), streamVersion,
		"stream module ConsensusVersion must be 2 after the 8-vaxildan widening")

	// Quietly use ctx so the import isn't pruned; ctx is the canonical handle a
	// runtime upgrade callback receives, so keeping a reference here is the right shape.
	_ = ctx
}

// TestRunMigrationsFromTaryonVersionMap simulates the canonical taryon→vaxildan in-place
// upgrade by handing a fromVM that matches the consensus versions a taryon-era binary
// would have committed, then asserts that ModuleManager.RunMigrations advances every
// module to its current version without error. This is the most direct stand-in for the
// upgrade-simulator rehearsal that stage 7 runs against the actual taryon binary.
func TestRunMigrationsFromTaryonVersionMap(t *testing.T) {
	app := apphelpers.Setup(t)
	ctx := app.NewContext(false)

	// Synthesise the "from" version map by starting from the current (post-vaxildan)
	// map and rolling x/stream back to 1. Modules whose consensus versions did not
	// change between taryon and vaxildan are left at their current value — a
	// no-op for RunMigrations, which is precisely what we want to assert.
	fromVM := app.ModuleManager.GetVersionMap()
	fromVM[streamtypes.ModuleName] = 1

	toVM, err := app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
	require.NoError(t, err, "RunMigrations must succeed for the 8-vaxildan upgrade plan")
	require.Equal(t, uint64(2), toVM[streamtypes.ModuleName],
		"x/stream must be at version 2 after RunMigrations")

	// Every entry from the input must appear in the output, and no entry should
	// regress. RunMigrations is supposed to be monotonic; this guards against a
	// future SDK change that would silently drop modules from the map.
	for name, fromVer := range fromVM {
		toVer, ok := toVM[name]
		require.True(t, ok, "module %s missing from post-migration VersionMap", name)
		require.GreaterOrEqual(t, toVer, fromVer, "module %s regressed: %d → %d", name, fromVer, toVer)
	}
}
