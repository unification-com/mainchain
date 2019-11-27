package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/enterprise"
	"github.com/unification-com/mainchain-cosmos/x/mint/internal/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper, keeper enterprise.Keeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	var mintedCoins sdk.Coins
	var mintedCoin sdk.Coin
	bondedRatio := k.BondedRatio(ctx)

	if ctx.BlockHeight() < 1577880 {
		mintedCoin = sdk.NewCoin(params.MintDenom, sdk.NewInt(25781428245))
		mintedCoins = sdk.NewCoins(mintedCoin)
	} else if ctx.BlockHeight() < 3155760 {
		mintedCoin = sdk.NewCoin(params.MintDenom, sdk.NewInt(12890714123))
		mintedCoins = sdk.NewCoins(mintedCoin)
	} else if ctx.BlockHeight() < 6311520 {
		mintedCoin = sdk.NewCoin(params.MintDenom, sdk.NewInt(6445357061))
		mintedCoins = sdk.NewCoins(mintedCoin)
	} else {
		// recalculate inflation rate
		totalUNDSupply := keeper.GetTotalUndSupply(ctx)
		totalLockedUND := keeper.GetTotalLockedUnd(ctx)
		liquidUND := totalUNDSupply.Sub(totalLockedUND)

		minter.Inflation = minter.NextInflationRate(params, bondedRatio)
		minter.AnnualProvisions = minter.NextAnnualProvisions(params, liquidUND.Amount)
		k.SetMinter(ctx, minter)

		// mint coins, update supply
		mintedCoin = minter.BlockProvision(params)
		mintedCoins = sdk.NewCoins(mintedCoin)
	}

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
