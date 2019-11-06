package enterprise

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain"
)

var (
	_ FeeTx = (*auth.StdTx)(nil) // assert StdTx implements FeeTx
)

// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
type FeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() sdk.Coins
	FeePayer() sdk.AccAddress
}

type CheckLockedUndDecorator struct {
	ak   auth.AccountKeeper
	entk Keeper
}

func NewCheckLockedUndDecorator(ak auth.AccountKeeper, entk Keeper) CheckLockedUndDecorator {
	return CheckLockedUndDecorator{
		ak:   ak,
		entk: entk,
	}
}

func (ld CheckLockedUndDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(FeeTx)

	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	blockTime := ctx.BlockHeader().Time

	if wrkchain.CheckIsWrkChainTx(feeTx) && ld.entk.IsLocked(ctx, feePayer) {
		// WRKChain Tx and has locked Enterprise UND.
		// check for and Undelegate any Locked UND to pay for fees
		// We undelegate and unlock here (instead of handler) because
		// fees are paid during the Ante process, further in the chain

		lockedUnd := ld.entk.GetLockedUnd(ctx, feePayer).Amount
		lockedUndCoins := sdk.NewCoins(lockedUnd)
		fees := feeTx.GetFee()

		// calculate how much Locked UND would be left over after deducting Tx fees
		_, hasNeg := lockedUndCoins.SafeSub(fees)

		if !hasNeg {
			// locked UND >= total fees
			// undelegate the fee amount to allow for payment
			err := ld.entk.SupplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, fees)

			if err != nil {
				return ctx, sdkerrors.Wrap(err, "failed to undelegate enterprise und")
			}

			// decrement the tracked locked UND
			feeNund := fees.AmountOf(DefaultDenomination)
			feeNundCoin := sdk.NewCoin(DefaultDenomination, feeNund)
			err = ld.entk.DecrementLockedUnd(ctx, feeTx.FeePayer(), feeNundCoin)
			if err != nil {
				return ctx, sdkerrors.Wrap(err, "failed to deduct locked und")
			}
		} else {
			// calculate how much can be undelegated, and if, by undelegating, the account
			// would have enough to pay for the fees. If not, don't undelegate
			feePayerAcc := ld.ak.GetAccount(ctx, feePayer)

			// How many spendable UND does the account have
			spendableCoins := feePayerAcc.SpendableCoins(blockTime)

			// calculate how much would be available if UND were unlocked
			potentiallyAvailable := spendableCoins.Add(lockedUndCoins)

			// is this enough to pay for the fees
			_, hasNeg := potentiallyAvailable.SafeSub(fees)

			if !hasNeg {
				// undelegate the fee amount to allow for payment
				err := ld.entk.SupplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, lockedUndCoins)

				if err != nil {
					return ctx, sdkerrors.Wrap(err, "failed to undelegate enterprise und")
				}

				err = ld.entk.DecrementLockedUnd(ctx, feeTx.FeePayer(), lockedUnd)
				if err != nil {
					return ctx, sdkerrors.Wrap(err, "failed to deduct locked und")
				}
			}

		}
	}

	return next(ctx, tx, simulate)
}
