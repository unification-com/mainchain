package ante

import (
	errorsmod "cosmossdk.io/errors"
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
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()

	if (wrkchain.CheckIsWrkChainTx(feeTx) || beacon.CheckIsBeaconTx(feeTx)) && ld.entk.IsLocked(ctx, feePayer) {
		// WRKChain/BEACON Tx and has locked Enterprise FUND.
		// check for and mint any Locked FUND to pay for fees
		// We unlock and mint (instead of msg_server) because
		// fees are paid during the Ante process, further in the chain
		// WRKChain/BEACON Txs have been checked before this decorator is called

		err := ld.entk.UnlockAndMintCoinsForFees(ctx, feePayer, feeTx.GetFee())

		if err != nil {
			return ctx, errorsmod.Wrap(err, "failed to unlock enterprise und")
		}

	}

	return next(ctx, tx, simulate)
}
