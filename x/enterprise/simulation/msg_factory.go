package simulation

import (
	"context"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/testutil/simsx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// MsgUndPurchaseOrderFactory raises a PO from a whitelisted account, then schedules
// future-block factories — one per ent signer — to process the PO via accept/reject.
// Decision scheduling is time-based via simsx.FutureOpsRegistry.
func MsgUndPurchaseOrderFactory(k keeper.Keeper) *simsx.LazyStateSimMsgFactory[*types.MsgUndPurchaseOrder] {
	return simsx.NewSimMsgFactoryWithFutureOps[*types.MsgUndPurchaseOrder](
		func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter, fOpsReg simsx.FutureOpsRegistry) ([]simsx.SimAccount, *types.MsgUndPurchaseOrder) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			r := testData.Rand()
			purchaser := testData.AnyAccount(reporter, simsx.WithSpendableBalance())
			if reporter.IsSkipped() {
				return nil, nil
			}
			if !k.AddressIsWhitelisted(sdkCtx, purchaser.Address) {
				reporter.Skip("address not whitelisted")
				return nil, nil
			}
			poId, err := k.GetHighestPurchaseOrderID(sdkCtx)
			if err != nil {
				reporter.Skip("cannot read highest PO ID")
				return nil, nil
			}
			amt := sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(r.IntInRange(1000, 1000000)))

			// Schedule a process op for each ent signer. Per legacy semantics:
			// most decisions arrive within a few blocks; some are delayed to test
			// the stale-PO path. Block time advances ~6s per block.
			entSignerList := k.GetParamEntSigners(sdkCtx)
			for _, signerStr := range strings.Split(entSignerList, ",") {
				signer, sErr := sdk.AccAddressFromBech32(signerStr)
				if sErr != nil {
					continue
				}
				delaySeconds := 6 * int64(r.IntInRange(1, 2))
				if r.Intn(4) == 3 {
					delaySeconds = 6 * 5
				}
				blockTime := sdkCtx.BlockTime().Add(time.Duration(delaySeconds) * time.Second)
				fOpsReg.Add(blockTime, processPurchaseOrderFactory(k, poId, signer))
			}

			return []simsx.SimAccount{purchaser}, types.NewMsgUndPurchaseOrder(purchaser.Address, amt)
		},
	)
}

// processPurchaseOrderFactory returns a factory that processes a specific PO via
// a specific ent signer. Used by MsgUndPurchaseOrderFactory's future-op schedule.
// Skips gracefully if the PO no longer exists or has already moved past STATUS_RAISED.
func processPurchaseOrderFactory(k keeper.Keeper, poId uint64, signer sdk.AccAddress) simsx.SimMsgFactoryFn[*types.MsgProcessUndPurchaseOrder] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgProcessUndPurchaseOrder) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		po, found := k.GetPurchaseOrder(sdkCtx, poId)
		if !found {
			reporter.Skip("PO not found")
			return nil, nil
		}
		if po.Status != types.StatusRaised {
			reporter.Skip("PO no longer raised")
			return nil, nil
		}
		signerAcc := testData.GetAccountbyAccAddr(reporter, signer)
		if reporter.IsSkipped() {
			return nil, nil
		}
		decision := types.StatusAccepted
		if testData.Rand().Intn(2) == 0 {
			decision = types.StatusRejected
		}
		return []simsx.SimAccount{signerAcc}, types.NewMsgProcessUndPurchaseOrder(po.Id, decision, signer)
	}
}

// MsgWhitelistAddressFactory adds or removes a random sim account from the
// enterprise whitelist using a random ent signer. Toggle semantics match legacy
// behaviour: if the target is whitelisted, action=remove; else action=add.
func MsgWhitelistAddressFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgWhitelistAddress] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgWhitelistAddress) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		target := testData.AnyAccount(reporter)
		if reporter.IsSkipped() {
			return nil, nil
		}
		action := types.WhitelistActionAdd
		if k.AddressIsWhitelisted(sdkCtx, target.Address) {
			action = types.WhitelistActionRemove
		}
		entSignerList := strings.Split(k.GetParamEntSigners(sdkCtx), ",")
		if len(entSignerList) == 0 {
			reporter.Skip("no ent signers configured")
			return nil, nil
		}
		r := testData.Rand()
		signerStr := entSignerList[r.Intn(len(entSignerList))]
		signer, err := sdk.AccAddressFromBech32(signerStr)
		if err != nil {
			reporter.Skip("invalid ent signer address")
			return nil, nil
		}
		signerAcc := testData.GetAccountbyAccAddr(reporter, signer)
		if reporter.IsSkipped() {
			return nil, nil
		}
		return []simsx.SimAccount{signerAcc}, types.NewMsgWhitelistAddress(target.Address, action, signer)
	}
}
