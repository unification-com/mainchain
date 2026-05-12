package keeper_test

import (
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	undapp "github.com/unification-com/mainchain/app"
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// setupBlockerTest configures the enterprise module with a realistic
// 2-of-3 multisig (the production configuration). Returns the test app,
// a context whose block time is fixed, the purchaser address, the three
// signer addresses, and the decision-time-limit param value.
func setupBlockerTest(t *testing.T) (
	app *undapp.App,
	ctx sdk.Context,
	purchaser sdk.AccAddress,
	signers []sdk.AccAddress,
	decisionTimeLimit uint64,
) {
	t.Helper()

	a := simapphelpers.Setup(t)
	c := a.BaseApp.NewContext(false)

	// Fix the block time so stale-PO maths are deterministic.
	c = c.WithBlockTime(time.Unix(1700000000, 0).UTC())

	purchaserAddrs := simapphelpers.GenerateRandomTestAccounts(1)
	signerAddrs := simapphelpers.GenerateRandomTestAccounts(3)

	signerStrs := make([]string, len(signerAddrs))
	for i, s := range signerAddrs {
		signerStrs[i] = s.String()
	}

	decisionLimit := uint64(84600) // default ~24h window

	err := a.EnterpriseKeeper.SetParams(c, types.Params{
		EntSigners:        strings.Join(signerStrs, ","),
		Denom:             sdk.DefaultBondDenom,
		MinAccepts:        2, // 2-of-3
		DecisionTimeLimit: decisionLimit,
	})
	require.NoError(t, err)

	// Whitelist the purchaser so RaiseNewPurchaseOrder succeeds in tests
	// that exercise it. The blocker doesn't depend on whitelist state, but
	// it's the path used to populate the raised queue.
	err = a.EnterpriseKeeper.AddAddressToWhitelist(c, purchaserAddrs[0])
	require.NoError(t, err)

	return a, c, purchaserAddrs[0], signerAddrs, decisionLimit
}

// raiseTestPO is a small helper that creates a raised purchase order
// for the given amount and returns its id.
func raiseTestPO(t *testing.T, app *undapp.App, ctx sdk.Context, purchaser sdk.AccAddress, amount int64) uint64 {
	t.Helper()
	poId, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, types.EnterpriseUndPurchaseOrder{
		Purchaser: purchaser.String(),
		Amount:    sdk.NewInt64Coin(sdk.DefaultBondDenom, amount),
	})
	require.NoError(t, err)
	return poId
}

// recordDecision appends a decision to a purchase order. Mirrors what
// the msg_server does, but skips the authorisation check (the blocker
// doesn't care who signed — only the count and Decision value matter).
func recordDecision(t *testing.T, app *undapp.App, ctx sdk.Context, poId uint64, signer sdk.AccAddress, decision types.PurchaseOrderStatus) {
	t.Helper()
	err := app.EnterpriseKeeper.ProcessPurchaseOrderDecision(ctx, poId, decision, signer)
	require.NoError(t, err)
}

// TestTallyPurchaseOrderDecisions_EmptyQueue: with no purchase orders raised,
// Tally must complete cleanly with no error and no state change.
func TestTallyPurchaseOrderDecisions_EmptyQueue(t *testing.T) {
	app, ctx, _, _, _ := setupBlockerTest(t)

	raisedBefore := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(ctx)
	require.Empty(t, raisedBefore)

	err := app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx)
	require.NoError(t, err)

	raisedAfter := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(ctx)
	require.Empty(t, raisedAfter)
	acceptedAfter := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx)
	require.Empty(t, acceptedAfter)
}

// TestTallyPurchaseOrderDecisions_Accept: PO with 2-of-3 accepts is
// moved from Raised → Accepted, removed from raised queue, added to
// accepted queue. Status and CompletionTime are set on the PO record.
func TestTallyPurchaseOrderDecisions_Accept(t *testing.T) {
	app, ctx, purchaser, signers, _ := setupBlockerTest(t)

	poId := raiseTestPO(t, app, ctx, purchaser, 100)
	recordDecision(t, app, ctx, poId, signers[0], types.StatusAccepted)
	recordDecision(t, app, ctx, poId, signers[1], types.StatusAccepted)
	// signer 2 abstains

	err := app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx)
	require.NoError(t, err)

	po, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
	require.True(t, found)
	require.Equal(t, types.StatusAccepted, po.Status)
	require.Equal(t, uint64(ctx.BlockHeader().Time.Unix()), po.CompletionTime)

	// Removed from raised queue
	raised := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(ctx)
	require.NotContains(t, raised, poId, "accepted PO must be removed from raised queue")

	// Added to accepted queue
	accepted := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx)
	require.Contains(t, accepted, poId, "accepted PO must be in accepted queue")
}

// TestTallyPurchaseOrderDecisions_Reject: PO with 2-of-3 rejects
// (numRejects=2, rejectThreshold = 3-2 = 1) is moved to Rejected and
// removed from the raised queue without being added to the accepted
// queue.
func TestTallyPurchaseOrderDecisions_Reject(t *testing.T) {
	app, ctx, purchaser, signers, _ := setupBlockerTest(t)

	poId := raiseTestPO(t, app, ctx, purchaser, 100)
	recordDecision(t, app, ctx, poId, signers[0], types.StatusRejected)
	recordDecision(t, app, ctx, poId, signers[1], types.StatusRejected)

	err := app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx)
	require.NoError(t, err)

	po, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
	require.True(t, found)
	require.Equal(t, types.StatusRejected, po.Status)
	require.Equal(t, uint64(ctx.BlockHeader().Time.Unix()), po.CompletionTime)

	raised := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(ctx)
	require.NotContains(t, raised, poId)
	accepted := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx)
	require.NotContains(t, accepted, poId, "rejected PO must NOT be in accepted queue")
}

// TestTallyPurchaseOrderDecisions_AutoRejectStale: PO with no decisions
// (or insufficient accepts) past the DecisionTimeLimit gets auto-rejected.
func TestTallyPurchaseOrderDecisions_AutoRejectStale(t *testing.T) {
	app, ctx, purchaser, _, decisionLimit := setupBlockerTest(t)

	poId := raiseTestPO(t, app, ctx, purchaser, 100)
	// no decisions recorded — PO will be stale once we advance time

	// Advance ctx block time past the decision limit
	staleCtx := ctx.WithBlockTime(ctx.BlockHeader().Time.Add(time.Duration(decisionLimit+1) * time.Second))

	err := app.EnterpriseKeeper.TallyPurchaseOrderDecisions(staleCtx)
	require.NoError(t, err)

	po, found := app.EnterpriseKeeper.GetPurchaseOrder(staleCtx, poId)
	require.True(t, found)
	require.Equal(t, types.StatusRejected, po.Status, "stale PO must be auto-rejected")
	require.Equal(t, uint64(staleCtx.BlockHeader().Time.Unix()), po.CompletionTime)

	raised := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(staleCtx)
	require.NotContains(t, raised, poId)
	accepted := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(staleCtx)
	require.NotContains(t, accepted, poId)
}

// TestTallyPurchaseOrderDecisions_Pending: PO with only one accept (below
// MinAccepts=2), no rejects, and within the decision window must remain in
// the raised queue with Status=Raised.
func TestTallyPurchaseOrderDecisions_Pending(t *testing.T) {
	app, ctx, purchaser, signers, _ := setupBlockerTest(t)

	poId := raiseTestPO(t, app, ctx, purchaser, 100)
	recordDecision(t, app, ctx, poId, signers[0], types.StatusAccepted) // only 1 of 2 needed

	err := app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx)
	require.NoError(t, err)

	po, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
	require.True(t, found)
	require.Equal(t, types.StatusRaised, po.Status, "pending PO must stay Raised")
	require.Equal(t, uint64(0), po.CompletionTime, "pending PO must not have a CompletionTime")

	raised := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(ctx)
	require.Contains(t, raised, poId, "pending PO must remain in raised queue")
	accepted := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx)
	require.NotContains(t, accepted, poId)
}

// TestTallyPurchaseOrderDecisions_MultipleOrders: a single Tally call must
// correctly dispatch each PO to its outcome — Accept, Reject, AutoReject,
// or remain Pending — in one pass.
func TestTallyPurchaseOrderDecisions_MultipleOrders(t *testing.T) {
	app, ctx, purchaser, signers, decisionLimit := setupBlockerTest(t)

	// PO #1: will be accepted (2 accepts).
	poAccept := raiseTestPO(t, app, ctx, purchaser, 100)
	recordDecision(t, app, ctx, poAccept, signers[0], types.StatusAccepted)
	recordDecision(t, app, ctx, poAccept, signers[1], types.StatusAccepted)

	// PO #2: will be rejected (2 rejects, threshold=1).
	poReject := raiseTestPO(t, app, ctx, purchaser, 200)
	recordDecision(t, app, ctx, poReject, signers[0], types.StatusRejected)
	recordDecision(t, app, ctx, poReject, signers[1], types.StatusRejected)

	// PO #3: stale, no decisions → will be auto-rejected.
	poStale := raiseTestPO(t, app, ctx, purchaser, 300)

	// PO #4: only 1 accept → remains pending.
	poPending := raiseTestPO(t, app, ctx, purchaser, 400)
	recordDecision(t, app, ctx, poPending, signers[0], types.StatusAccepted)

	// Advance time past the limit so poStale qualifies as stale. POs 1, 2, 4
	// also see the advanced time; for poAccept and poReject the deciding
	// branch fires before the stale check (and even if it didn't, the stale
	// check requires numAccepts < MinAccepts, which fails for poAccept).
	// poPending has only 1 accept (< MinAccepts=2) and IS stale at this
	// point — it will therefore be auto-rejected here too.
	staleCtx := ctx.WithBlockTime(ctx.BlockHeader().Time.Add(time.Duration(decisionLimit+1) * time.Second))

	err := app.EnterpriseKeeper.TallyPurchaseOrderDecisions(staleCtx)
	require.NoError(t, err)

	poAcc, _ := app.EnterpriseKeeper.GetPurchaseOrder(staleCtx, poAccept)
	require.Equal(t, types.StatusAccepted, poAcc.Status, "poAccept")

	poRej, _ := app.EnterpriseKeeper.GetPurchaseOrder(staleCtx, poReject)
	require.Equal(t, types.StatusRejected, poRej.Status, "poReject")

	poStl, _ := app.EnterpriseKeeper.GetPurchaseOrder(staleCtx, poStale)
	require.Equal(t, types.StatusRejected, poStl.Status, "poStale")

	poPend, _ := app.EnterpriseKeeper.GetPurchaseOrder(staleCtx, poPending)
	require.Equal(t, types.StatusRejected, poPend.Status, "poPending (also stale, also auto-rejected)")

	// Only poAccept goes to the accepted queue.
	accepted := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(staleCtx)
	require.Contains(t, accepted, poAccept)
	require.NotContains(t, accepted, poReject)
	require.NotContains(t, accepted, poStale)
	require.NotContains(t, accepted, poPending)

	// Raised queue must be empty — every PO was decided this tally round.
	raised := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(staleCtx)
	require.Empty(t, raised)
}

// TestProcessAcceptedPurchaseOrders_EmptyQueue: no accepted POs → no-op.
func TestProcessAcceptedPurchaseOrders_EmptyQueue(t *testing.T) {
	app, ctx, _, _, _ := setupBlockerTest(t)

	accepted := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx)
	require.Empty(t, accepted)

	err := app.EnterpriseKeeper.ProcessAcceptedPurchaseOrders(ctx)
	require.NoError(t, err)

	accepted = app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx)
	require.Empty(t, accepted)
}

// TestProcessAcceptedPurchaseOrders_Mints: an accepted PO transitions to
// Completed, eFUND is locked for the purchaser in the requested amount,
// the PO is removed from the accepted queue, and total locked / module
// supply update accordingly.
func TestProcessAcceptedPurchaseOrders_Mints(t *testing.T) {
	app, ctx, purchaser, signers, _ := setupBlockerTest(t)

	amount := int64(100)
	poId := raiseTestPO(t, app, ctx, purchaser, amount)
	recordDecision(t, app, ctx, poId, signers[0], types.StatusAccepted)
	recordDecision(t, app, ctx, poId, signers[1], types.StatusAccepted)

	// First, Tally moves it into the accepted queue.
	require.NoError(t, app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx))
	require.Contains(t, app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx), poId)

	totalLockedBefore := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	purchaserLockedBefore := app.EnterpriseKeeper.GetLockedUndAmountForAccount(ctx, purchaser)

	err := app.EnterpriseKeeper.ProcessAcceptedPurchaseOrders(ctx)
	require.NoError(t, err)

	// Status flipped to Completed
	po, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
	require.True(t, found)
	require.Equal(t, types.StatusCompleted, po.Status)

	// Removed from accepted queue
	require.NotContains(t, app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx), poId)

	// Purchaser's locked eFUND grew by exactly the PO amount
	purchaserLockedAfter := app.EnterpriseKeeper.GetLockedUndAmountForAccount(ctx, purchaser)
	expDelta := sdk.NewInt64Coin(sdk.DefaultBondDenom, amount)
	require.True(t, purchaserLockedAfter.IsEqual(purchaserLockedBefore.Add(expDelta)),
		"purchaser locked: before %s + %s = %s, got %s",
		purchaserLockedBefore, expDelta, purchaserLockedBefore.Add(expDelta), purchaserLockedAfter)

	// Total locked grew by the same amount
	totalLockedAfter := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedAfter.IsEqual(totalLockedBefore.Add(expDelta)),
		"total locked: before %s + %s = %s, got %s",
		totalLockedBefore, expDelta, totalLockedBefore.Add(expDelta), totalLockedAfter)
}

// TestProcessAcceptedPurchaseOrders_MultipleOrders: process N accepted
// orders in one call. Each should be Completed and minted; total locked
// should reflect the sum.
func TestProcessAcceptedPurchaseOrders_MultipleOrders(t *testing.T) {
	app, ctx, purchaser, signers, _ := setupBlockerTest(t)

	amounts := []int64{100, 250, 500}
	poIds := make([]uint64, 0, len(amounts))
	expectedTotalDelta := int64(0)

	for _, amt := range amounts {
		poId := raiseTestPO(t, app, ctx, purchaser, amt)
		recordDecision(t, app, ctx, poId, signers[0], types.StatusAccepted)
		recordDecision(t, app, ctx, poId, signers[1], types.StatusAccepted)
		poIds = append(poIds, poId)
		expectedTotalDelta += amt
	}

	require.NoError(t, app.EnterpriseKeeper.TallyPurchaseOrderDecisions(ctx))
	require.Len(t, app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx), len(amounts))

	totalLockedBefore := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)

	err := app.EnterpriseKeeper.ProcessAcceptedPurchaseOrders(ctx)
	require.NoError(t, err)

	// Every PO is Completed
	for _, poId := range poIds {
		po, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
		require.True(t, found)
		require.Equal(t, types.StatusCompleted, po.Status, "poId %d", poId)
	}

	// Accepted queue drained
	require.Empty(t, app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx))

	// Total locked grew by the sum of amounts
	totalLockedAfter := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	expDelta := sdk.NewInt64Coin(sdk.DefaultBondDenom, expectedTotalDelta)
	require.True(t, totalLockedAfter.IsEqual(totalLockedBefore.Add(expDelta)))
}
