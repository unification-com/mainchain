package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	beacon "github.com/unification-com/mainchain/x/beacon/exported"
	wrkchain "github.com/unification-com/mainchain/x/wrkchain/exported"
)

type CheckLockedUndDecorator struct {
	entk EnterpriseKeeper
}

func NewCheckLockedUndDecorator(entk EnterpriseKeeper) CheckLockedUndDecorator {
	return CheckLockedUndDecorator{
		entk: entk,
	}
}

func (ld CheckLockedUndDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)

	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()

	if (wrkchain.CheckIsWrkChainTx(feeTx) || beacon.CheckIsBeaconTx(feeTx)) && ld.entk.IsLocked(ctx, feePayer) {
		// WRKChain/BEACON Tx and has locked Enterprise FUND.
		// check for and Undelegate any Locked FUND to pay for fees
		// We undelegate and unlock here (instead of handler) because
		// fees are paid during the Ante process, further in the chain
		// WRKChain/BEACON Txs have been checked before this decorator is called

		err := ld.entk.UnlockCoinsForFees(ctx, feePayer, feeTx.GetFee())

		if err != nil {
			return ctx, sdkerrors.Wrap(err, "failed to unlock enterprise und")
		}

	}

	return next(ctx, tx, simulate)
}
