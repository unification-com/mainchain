package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/testutil/simsx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/beacon/keeper"
	"github.com/unification-com/mainchain/x/beacon/types"
)

func MsgRegisterBeaconFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgRegisterBeacon] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgRegisterBeacon) {
		fee := k.GetRegistrationFeeAsCoin(sdk.UnwrapSDKContext(ctx))
		owner := testData.AnyAccount(reporter, simsx.WithLiquidBalanceGTE(fee))
		if reporter.IsSkipped() {
			return nil, nil
		}
		r := testData.Rand()
		moniker := r.StringN(64)
		name := r.StringN(128)
		return []simsx.SimAccount{owner}, types.NewMsgRegisterBeacon(moniker, name, owner.Address)
	}
}

func MsgRecordBeaconTimestampFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgRecordBeaconTimestamp] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgRecordBeaconTimestamp) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		beacons := k.GetAllBeacons(sdkCtx)
		if len(beacons) == 0 {
			reporter.Skip("no beacons registered")
			return nil, nil
		}
		r := testData.Rand()
		beacon := beacons[r.Intn(len(beacons))]
		ownerAddr, err := sdk.AccAddressFromBech32(beacon.Owner)
		if err != nil {
			reporter.Skip("invalid beacon owner address")
			return nil, nil
		}
		fee := k.GetRecordFeeAsCoin(sdkCtx)
		owner := testData.GetAccountbyAccAddr(reporter, ownerAddr)
		if reporter.IsSkipped() {
			return nil, nil
		}
		if !owner.LiquidBalance().AmountOf(fee.Denom).GTE(fee.Amount) {
			reporter.Skip("beacon owner cannot pay record fee")
			return nil, nil
		}
		hash := r.StringN(64)
		return []simsx.SimAccount{owner}, types.NewMsgRecordBeaconTimestamp(
			beacon.BeaconId, hash, uint64(sdkCtx.BlockTime().Unix()), owner.Address,
		)
	}
}

func MsgPurchaseBeaconStateStorageFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgPurchaseBeaconStateStorage] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgPurchaseBeaconStateStorage) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		beacons := k.GetAllBeacons(sdkCtx)
		if len(beacons) == 0 {
			reporter.Skip("no beacons registered")
			return nil, nil
		}
		r := testData.Rand()
		beacon := beacons[r.Intn(len(beacons))]
		maxCanPurchase := k.GetMaxPurchasableSlots(sdkCtx, beacon.BeaconId)
		if maxCanPurchase == 0 {
			reporter.Skip("beacon max storage reached")
			return nil, nil
		}
		ownerAddr, err := sdk.AccAddressFromBech32(beacon.Owner)
		if err != nil {
			reporter.Skip("invalid beacon owner address")
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
			reporter.Skip("beacon owner cannot pay purchase storage fee")
			return nil, nil
		}
		return []simsx.SimAccount{owner}, types.NewMsgPurchaseBeaconStateStorage(beacon.BeaconId, randNumToPurchase, owner.Address)
	}
}
