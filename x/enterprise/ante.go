package enterprise

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	beacon "github.com/unification-com/mainchain-cosmos/x/beacon/exported"
	wrkchain "github.com/unification-com/mainchain-cosmos/x/wrkchain/exported"
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
	entk Keeper
}

func NewCheckLockedUndDecorator(entk Keeper) CheckLockedUndDecorator {
	return CheckLockedUndDecorator{
		entk: entk,
	}
}

func (ld CheckLockedUndDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(FeeTx)

	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()

	if (wrkchain.CheckIsWrkChainTx(feeTx) || beacon.CheckIsBeaconTx(feeTx)) && ld.entk.IsLocked(ctx, feePayer) {
		// WRKChain/BEACON Tx and has locked Enterprise UND.
		// check for and Undelegate any Locked UND to pay for fees
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
