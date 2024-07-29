package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FreezerDecorator struct{}

func NewFreezerDecorator() FreezerDecorator {
	return FreezerDecorator{}
}

func (fd FreezerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	//sigTx, ok := tx.(authsigning.SigVerifiableTx)
	//
	//if !ok {
	//	return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	//}
	//
	//msgs := tx.GetMsgs()
	//
	//// check all signers in all messages in the Tx
	//for _, msg := range msgs {
	//	for i, signer := range msg.GetSigners() {
	//		if isFrozen(signer) {
	//			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "wallet %s in msg %d is frozen", signer, i)
	//		}
	//	}
	//}
	//
	//// check actual Tx signers
	//for _, signer := range sigTx.GetSigners() {
	//	if isFrozen(signer) {
	//		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "wallet %s frozen", signer)
	//	}
	//}

	return next(ctx, tx, simulate)
}

func isFrozen(signer sdk.AccAddress) bool {
	//frozen, _ := sdk.AccAddressFromBech32("und18mcmhkq6fmhu9hpy3sx5cugqwv6z0wrz7nn5d7")
	//return signer.Equals(frozen)
	return false
}
