//go:build sims

package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/testutil/simsx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// Blank import: runs SetAddressPrefixes/RegisterDenoms in its init() so that
	// NewApp captures the und HRP in its keeper codecs. Without this, default
	// enterprise params (und1qqqq...) fail to decode against the captured cosmos
	// HRP and the simulator aborts on the first stream/enterprise op.
	_ "github.com/unification-com/mainchain/app/params"
)

func init() {
	simcli.GetSimulatorFlags()
}

func setupStateFactory(app *App) simsx.SimStateFactory {
	return simsx.SimStateFactory{
		Codec:         app.AppCodec(),
		AppStateFn:    simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		BlockedAddr:   BlockedAddresses(),
		AccountSource: app.AccountKeeper,
		BalanceSource: app.BankKeeper,
	}
}

// isEmptyValidatorSetErr lets us skip-not-fail when the simulator ends with zero
// validators (the v0.50.7-fixed bug class).
func isEmptyValidatorSetErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "validator set is empty after InitGenesis")
}

func assertEqualStores(t *testing.T, a, b *App, ctxA, ctxB sdk.Context) {
	t.Helper()

	skipPrefixes := map[string][][]byte{
		stakingtypes.StoreKey: {
			stakingtypes.UnbondingQueueKey,
			stakingtypes.RedelegationQueueKey,
			stakingtypes.ValidatorQueueKey,
			stakingtypes.HistoricalInfoKey,
			stakingtypes.UnbondingIDKey,
			stakingtypes.UnbondingIndexKey,
			stakingtypes.UnbondingTypeKey,
			stakingtypes.ValidatorUpdatesKey,
		},
		authzkeeper.StoreKey:   {authzkeeper.GrantQueuePrefix},
		feegrant.StoreKey:      {feegrant.FeeAllowanceQueueKeyPrefix.Bytes()},
		slashingtypes.StoreKey: {slashingtypes.ValidatorMissedBlockBitmapKeyPrefix},
	}

	storeKeys := a.GetStoreKeys()
	require.NotEmpty(t, storeKeys)

	for _, keyA := range storeKeys {
		if _, ok := keyA.(*storetypes.KVStoreKey); !ok {
			continue
		}
		name := keyA.Name()
		keyB := b.GetKey(name)

		storeA := ctxA.KVStore(keyA)
		storeB := ctxB.KVStore(keyB)

		failedA, failedB := simtestutil.DiffKVStores(storeA, storeB, skipPrefixes[name])
		require.Equal(t, len(failedA), len(failedB), "unequal sets of key-values to compare %s", name)
		fmt.Printf("compared %d different key/value pairs between %s and %s\n", len(failedA), keyA, keyB)
		require.Equal(t, 0, len(failedA),
			simtestutil.GetSimulationLog(name, a.SimulationManager().StoreDecoders, failedA, failedB))
	}
}

// Hold references so the helpers above keep compiling while the remaining three
// test bodies below are stubs. Replaced by real call sites when the bodies land.
var (
	_ = isEmptyValidatorSetErr
	_ = assertEqualStores
)

// TestFullAppSimulation runs the simulator end-to-end. With no `-Seed=N` flag,
// simsx iterates its 38 default seeds in parallel; pass `-Seed=N` to run a single
// seed against a fixed value (e.g. for reproducing a failure).
func TestFullAppSimulation(t *testing.T) {
	cfg := simcli.NewConfigFromFlags()
	if cfg.Seed != simcli.DefaultSeedValue {
		cfg.ChainID = simsx.SimAppChainID
		simsx.RunWithSeed(t, cfg, NewApp, setupStateFactory, cfg.Seed, nil)
		return
	}
	simsx.Run(t, NewApp, setupStateFactory)
}

func TestAppImportExport(t *testing.T) {
	t.Skip("WIP: body lands in stage 7 follow-up commit; see docs/planning/vaxildan/07-stage7-tests.md")
}

func TestAppSimulationAfterImport(t *testing.T) {
	t.Skip("WIP: body lands in stage 7 follow-up commit; see docs/planning/vaxildan/07-stage7-tests.md")
}

func TestAppStateDeterminism(t *testing.T) {
	t.Skip("WIP: body lands in stage 7 follow-up commit; see docs/planning/vaxildan/07-stage7-tests.md")
}
