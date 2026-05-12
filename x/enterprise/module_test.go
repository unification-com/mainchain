package enterprise_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/enterprise"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}

// TestImportExportGenesisRoundTrip seeds rich enterprise state into a running
// app, exports the genesis, re-imports it into the same app, then re-exports.
// The two exports must be byte-equivalent. This catches:
//   - schema-format regressions during the upgrade (any change that breaks
//     serialise/deserialise symmetry surfaces here)
//   - state that's written by SetXxx but not read back by GetXxx (or vice versa)
//   - non-idempotent ExportGenesis or InitGenesis
//
// Mirrors the pattern in x/stream/keeper/genesis_test.go:TestImportExportGenesis.
// TA-28 — enterprise ExportGenesis was 0% covered before this test.
func TestImportExportGenesisRoundTrip(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	// Seed addresses (funded so the bank-related view is realistic).
	purchasers := simapphelpers.GenerateRandomTestAccounts(3)
	signers := simapphelpers.GenerateRandomTestAccounts(3)
	whitelisted := simapphelpers.GenerateRandomTestAccounts(2)

	denom := sdk.DefaultBondDenom

	// Configure params with a 2-of-3 signer set.
	signerStrs := make([]string, len(signers))
	for i, s := range signers {
		signerStrs[i] = s.String()
	}
	err := app.EnterpriseKeeper.SetParams(ctx, types.Params{
		EntSigners:        strings.Join(signerStrs, ","),
		Denom:             denom,
		MinAccepts:        2,
		DecisionTimeLimit: 84600,
	})
	require.NoError(t, err)

	// Whitelist a couple of addresses.
	for _, w := range whitelisted {
		require.NoError(t, app.EnterpriseKeeper.AddAddressToWhitelist(ctx, w))
	}
	// Also whitelist the purchasers so RaiseNewPurchaseOrder succeeds.
	for _, p := range purchasers {
		require.NoError(t, app.EnterpriseKeeper.AddAddressToWhitelist(ctx, p))
	}

	// Raise three purchase orders with distinct statuses to exercise all
	// queue paths in InitGenesis (raised → raised queue, accepted → accepted
	// queue, completed → no queue).
	// PO that stays in Raised status (will land in the raised queue on InitGenesis).
	_, err = app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, types.EnterpriseUndPurchaseOrder{
		Purchaser: purchasers[0].String(),
		Amount:    sdk.NewInt64Coin(denom, 100),
	})
	require.NoError(t, err)

	// PO that gets tallied to Accepted (will land in the accepted queue on InitGenesis).
	poAccepted, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, types.EnterpriseUndPurchaseOrder{
		Purchaser: purchasers[1].String(),
		Amount:    sdk.NewInt64Coin(denom, 200),
	})
	require.NoError(t, err)
	require.NoError(t, app.EnterpriseKeeper.ProcessPurchaseOrderDecision(ctx, poAccepted, types.StatusAccepted, signers[0]))
	require.NoError(t, app.EnterpriseKeeper.ProcessPurchaseOrderDecision(ctx, poAccepted, types.StatusAccepted, signers[1]))
	require.NoError(t, app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx))

	// PO that runs end-to-end to Completed (no queue placement on re-import).
	poCompleted, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, types.EnterpriseUndPurchaseOrder{
		Purchaser: purchasers[2].String(),
		Amount:    sdk.NewInt64Coin(denom, 300),
	})
	require.NoError(t, err)
	require.NoError(t, app.EnterpriseKeeper.ProcessPurchaseOrderDecision(ctx, poCompleted, types.StatusAccepted, signers[0]))
	require.NoError(t, app.EnterpriseKeeper.ProcessPurchaseOrderDecision(ctx, poCompleted, types.StatusAccepted, signers[1]))
	require.NoError(t, app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx))
	require.NoError(t, app.EnterpriseKeeper.ProcessAcceptedPurchaseOrders(ctx))

	// Independently seed locked + spent for one purchaser so the LockedUnd
	// and SpentEFUND fields in the genesis are non-empty.
	require.NoError(t, app.EnterpriseKeeper.CreateAndLockEFUND(ctx, purchasers[0], sdk.NewInt64Coin(denom, 1500)))
	require.NoError(t, app.EnterpriseKeeper.SetSpentEFUNDForAccount(ctx, types.SpentEFUND{
		Owner:  purchasers[0].String(),
		Amount: sdk.NewInt64Coin(denom, 250),
	}))

	// Move the total-spent counter so it's non-zero (matches what
	// incrementSpentEFUND would do internally in production).
	require.NoError(t, app.EnterpriseKeeper.SetTotalSpentEFUND(ctx, sdk.NewInt64Coin(denom, 250)))

	// Sanity: confirm the seeded state actually has the breadth we need.
	require.NotEmpty(t, app.EnterpriseKeeper.GetAllPurchaseOrders(ctx), "must have POs")
	require.NotEmpty(t, app.EnterpriseKeeper.GetAllLockedUnds(ctx), "must have locked entries")
	require.NotEmpty(t, app.EnterpriseKeeper.GetAllSpentEFUNDs(ctx), "must have spent entries")
	require.NotEmpty(t, app.EnterpriseKeeper.GetAllWhitelistedAddresses(ctx), "must have whitelisted addrs")
	require.True(t, app.EnterpriseKeeper.GetTotalLockedUnd(ctx).IsPositive(), "totalLocked must be positive")
	require.True(t, app.EnterpriseKeeper.GetTotalSpentEFUND(ctx).IsPositive(), "totalSpent must be positive")

	// Capture the first export.
	firstExport := enterprise.ExportGenesis(ctx, app.EnterpriseKeeper)
	require.NotNil(t, firstExport)

	// Re-import the same state. InitGenesis must be idempotent for state that
	// was just exported.
	enterprise.InitGenesis(ctx, app.EnterpriseKeeper, app.BankKeeper, app.AccountKeeper, *firstExport)

	// Re-export. Should match the first export exactly.
	secondExport := enterprise.ExportGenesis(ctx, app.EnterpriseKeeper)
	require.NotNil(t, secondExport)

	require.Equal(t, firstExport, secondExport, "ExportGenesis must be idempotent across InitGenesis re-imports")
}
