package ante_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"

	undapp "github.com/unification-com/mainchain/app"
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	beaconante "github.com/unification-com/mainchain/x/beacon/ante"
	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	entante "github.com/unification-com/mainchain/x/enterprise/ante"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	wrkante "github.com/unification-com/mainchain/x/wrkchain/ante"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
)

// GAS PARITY — the guard against accidentally state-breaking a patch release.
//
// The BEACON and WRKChain decorators, and the enterprise decorator that unlocks eFUND to pay their
// fees, all read module params from the store. Store reads in the ante chain are gas metered, that
// gas lands in the Tx's GasUsed, and GasUsed is hashed into the block's LastResultsHash. So ANY
// change to how many store reads these decorators perform during delivery changes the results hash,
// and nodes on either side of the change compute different ones — a chain split.
//
// It does not take a visible refactor to do this. Deleting a line whose argument happens to be a
// param getter is enough, because the read is inside the call:
//
//	_, feeToPay := feesToPay.Find(k.GetParamDenom(ctx))   // <- that is a store read
//
// An idle chain will not show it. Empty blocks execute identically either way, so the split only
// appears once a block carries a BEACON or WRKChain Tx — by which time the release has shipped.
//
// This test pins the exact gas the custom ante segment consumes during DELIVERY for a canonical Tx
// of each shape. If a change moves any of these numbers, the change is state-machine breaking:
//
//   - it does NOT belong in a patch release,
//   - it needs a coordinated upgrade with an upgrade handler and a height,
//   - and the constant below is updated as part of that upgrade, never to make the test pass.
//
// The numbers were verified equal against the v1.13.0 tree, decorator for decorator.
const (
	gasBeaconRecordLiquid   = 8581  // BEACON record, fee paid from liquid FUND
	gasBeaconRecordLocked   = 59239 // BEACON record, fee unlocked+minted from locked eFUND
	gasWrkChainRecordLiquid = 8581  // WRKChain record, fee paid from liquid FUND
	gasWrkChainRecordLocked = 59239 // WRKChain record, fee unlocked+minted from locked eFUND
)

const gasParityChainID = "und-gas-parity-test"

// customAnteSegment chains the three custom decorators in the order NewAnteHandler wires them, so
// what is measured is what a validator actually runs.
func customAnteSegment(app *undapp.App) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		wrkante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper),
		beaconante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper),
		entante.NewCheckLockedUndDecorator(app.EnterpriseKeeper),
	)
}

func gasParityFund(ctx sdk.Context, bk bankkeeper.Keeper, addr sdk.AccAddress, amt sdk.Coins) error {
	if err := bk.MintCoins(ctx, enttypes.ModuleName, amt); err != nil {
		return err
	}
	return bk.SendCoinsFromModuleToAccount(ctx, enttypes.ModuleName, addr, amt)
}

// measureDeliveryGas runs one Tx through the custom ante segment in DELIVERY mode (IsCheckTx
// false, simulate false) and returns the gas consumed.
func measureDeliveryGas(t *testing.T, msgFor func(sdk.AccAddress) sdk.Msg, locked bool) uint64 {
	t.Helper()

	r := rand.New(rand.NewSource(1))
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false) // delivery, NOT CheckTx
	txGen := app.GetTxConfig()

	const recordFee = 2

	// The BEACON, WRKChain and enterprise modules each carry their own fee denomination, and
	// UnlockAndMintCoinsForFees looks the Tx fee up by the enterprise one. In production all three
	// are the chain's bond denom, so take that as given and align the other two to it — the
	// measurement then reflects the path a validator actually walks.
	denom := app.EnterpriseKeeper.GetParams(ctx).Denom
	require.NotEmpty(t, denom)

	require.NoError(t, app.BeaconKeeper.SetParams(ctx, beacontypes.NewParams(24, recordFee, 2, denom, 200, 300)))
	require.NoError(t, app.WrkchainKeeper.SetParams(ctx, wrkchaintypes.NewParams(24, recordFee, 2, denom, 200, 300)))

	privK := ed25519.GenPrivKey()
	addr := sdk.AccAddress(privK.PubKey().Address())
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)

	if locked {
		// Only a dust liquid balance, so the fee must come from unlocked eFUND — this is what
		// exercises UnlockAndMintCoinsForFees, and therefore x/enterprise/keeper/locked.go.
		require.NoError(t, gasParityFund(ctx, app.BankKeeper, addr, sdk.NewCoins(sdk.NewInt64Coin(denom, 1))))
		require.NoError(t, app.EnterpriseKeeper.SetLockedUndForAccount(ctx, enttypes.LockedUnd{
			Owner:  addr.String(),
			Amount: sdk.NewInt64Coin(denom, recordFee*100),
		}))
	} else {
		require.NoError(t, gasParityFund(ctx, app.BankKeeper, addr, sdk.NewCoins(sdk.NewInt64Coin(denom, 1000000))))
	}

	fee := sdk.NewCoins(sdk.NewInt64Coin(denom, recordFee))
	tx, err := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msgFor(addr)}, fee, uint64(0), gasParityChainID, []uint64{0}, []uint64{0}, privK)
	require.NoError(t, err)

	before := ctx.GasMeter().GasConsumed()
	_, err = customAnteSegment(app)(ctx, tx, false)
	require.NoError(t, err)

	return ctx.GasMeter().GasConsumed() - before
}

func TestDeliveryGasParity(t *testing.T) {
	beaconMsg := func(addr sdk.AccAddress) sdk.Msg {
		return beacontypes.NewMsgRecordBeaconTimestamp(1, "somehash", 1, addr)
	}
	wrkchainMsg := func(addr sdk.AccAddress) sdk.Msg {
		return wrkchaintypes.NewMsgRecordWrkChainBlock(1, 1, "blockhash", "", "", "", "", addr)
	}

	for _, tc := range []struct {
		name   string
		msgFor func(sdk.AccAddress) sdk.Msg
		locked bool
		want   uint64
	}{
		{"beacon record, liquid FUND", beaconMsg, false, gasBeaconRecordLiquid},
		{"beacon record, locked eFUND", beaconMsg, true, gasBeaconRecordLocked},
		{"wrkchain record, liquid FUND", wrkchainMsg, false, gasWrkChainRecordLiquid},
		{"wrkchain record, locked eFUND", wrkchainMsg, true, gasWrkChainRecordLocked},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := measureDeliveryGas(t, tc.msgFor, tc.locked)
			require.Equal(t, tc.want, got, deliveryGasFailureMessage(tc.name, tc.want, got))
		})
	}
}

func deliveryGasFailureMessage(name string, want, got uint64) string {
	return fmt.Sprintf(
		"delivery-path gas for %q changed: %d -> %d.\n\n"+
			"GasUsed is hashed into the block's LastResultsHash, so a node running this code computes a\n"+
			"different results hash from one that does not. Nodes on either side of the change will split\n"+
			"the moment a block carries a BEACON or WRKChain Tx.\n\n"+
			"If this change is intended it is a STATE-MACHINE change: it needs a coordinated upgrade with an\n"+
			"upgrade handler and a height, and this constant is updated as part of that upgrade. Do NOT update\n"+
			"the constant to make the test pass in a patch release.\n\n"+
			"Most likely cause: a store read added to or removed from the ante path. Watch for param getters\n"+
			"used as call arguments — `feesToPay.Find(k.GetParamDenom(ctx))` contains a store read.",
		name, want, got)
}
