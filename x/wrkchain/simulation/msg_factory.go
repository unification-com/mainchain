package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/testutil/simsx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/wrkchain/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

func MsgRegisterWrkChainFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgRegisterWrkChain] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgRegisterWrkChain) {
		fee := k.GetRegistrationFeeAsCoin(sdk.UnwrapSDKContext(ctx))
		owner := testData.AnyAccount(reporter, simsx.WithLiquidBalanceGTE(fee))
		if reporter.IsSkipped() {
			return nil, nil
		}
		r := testData.Rand()
		moniker := r.StringN(64)
		genesisHash := r.StringN(64)
		name := r.StringN(128)
		baseType := r.StringN(5)
		return []simsx.SimAccount{owner}, types.NewMsgRegisterWrkChain(moniker, genesisHash, name, baseType, owner.Address)
	}
}

func MsgRecordWrkChainBlockFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgRecordWrkChainBlock] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgRecordWrkChainBlock) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		wrkchains := k.GetAllWrkChains(sdkCtx)
		if len(wrkchains) == 0 {
			reporter.Skip("no wrkchains registered")
			return nil, nil
		}
		r := testData.Rand()
		wc := wrkchains[r.Intn(len(wrkchains))]
		ownerAddr, err := sdk.AccAddressFromBech32(wc.Owner)
		if err != nil {
			reporter.Skip("invalid wrkchain owner address")
			return nil, nil
		}
		fee := k.GetRecordFeeAsCoin(sdkCtx)
		owner := testData.GetAccountbyAccAddr(reporter, ownerAddr)
		if reporter.IsSkipped() {
			return nil, nil
		}
		if !owner.LiquidBalance().AmountOf(fee.Denom).GTE(fee.Amount) {
			reporter.Skip("wrkchain owner cannot pay record fee")
			return nil, nil
		}
		hash := r.StringN(64)
		ph, h1, h2, h3 := "", "", "", ""
		if wc.WrkchainId%2 == 0 {
			ph = r.StringN(64)
			h1 = r.StringN(64)
			h2 = r.StringN(64)
			h3 = r.StringN(64)
		}
		return []simsx.SimAccount{owner}, types.NewMsgRecordWrkChainBlock(
			wc.WrkchainId, wc.Lastblock+1, hash, ph, h1, h2, h3, owner.Address,
		)
	}
}

func MsgPurchaseWrkChainStateStorageFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgPurchaseWrkChainStateStorage] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgPurchaseWrkChainStateStorage) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		wrkchains := k.GetAllWrkChains(sdkCtx)
		if len(wrkchains) == 0 {
			reporter.Skip("no wrkchains registered")
			return nil, nil
		}
		r := testData.Rand()
		wc := wrkchains[r.Intn(len(wrkchains))]
		maxCanPurchase := k.GetMaxPurchasableSlots(sdkCtx, wc.WrkchainId)
		if maxCanPurchase == 0 {
			reporter.Skip("wrkchain max storage reached")
			return nil, nil
		}
		ownerAddr, err := sdk.AccAddressFromBech32(wc.Owner)
		if err != nil {
			reporter.Skip("invalid wrkchain owner address")
			return nil, nil
		}
		owner := testData.GetAccountbyAccAddr(reporter, ownerAddr)
		if reporter.IsSkipped() {
			return nil, nil
		}
		randNumToPurchase := uint64(1)
		if maxCanPurchase > 1 {
			randNumToPurchase = uint64(r.Intn(int(maxCanPurchase))) + 1
		}
		params := k.GetParams(sdkCtx)
		feeAsCoin := sdk.NewInt64Coin(params.Denom, int64(params.FeePurchaseStorage*randNumToPurchase))
		if !owner.LiquidBalance().AmountOf(feeAsCoin.Denom).GTE(feeAsCoin.Amount) {
			reporter.Skip("wrkchain owner cannot pay purchase storage fee")
			return nil, nil
		}
		return []simsx.SimAccount{owner}, types.NewMsgPurchaseWrkChainStateStorage(wc.WrkchainId, randNumToPurchase, owner.Address)
	}
}
