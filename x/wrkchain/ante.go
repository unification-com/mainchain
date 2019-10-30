package wrkchain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
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

// CorrectWrkChainFeeDecorator checks if the correct fees have been sent to pay for a
// WRKChain register/record hash Tx, and if the fee paying account has sufficient funds to pay.
// It first checks if the Tx contains any WRKChain Msgs, and if not, continues on to the next
// AnteHAndler in the chain. If a WRKChain Msg is detected, it then:
//
// 1. Checks sufficient fees have been included in the Tx, via the --fees flag
// 2. Checks the fee payer is the WRKChain owner
// 3. Checks if the fee payer has sufficient funds in their account to pay for it
//
// If any of the checks fail, a suitable error is returned.
type CorrectWrkChainFeeDecorator struct{
	ak           auth.AccountKeeper
	wck          Keeper
}

func NewWrkChainFeeDecorator(ak auth.AccountKeeper, wrkchainKeeper Keeper) CorrectWrkChainFeeDecorator {
	return CorrectWrkChainFeeDecorator{
		ak: ak,
		wck: wrkchainKeeper,
	}
}

func (wfd CorrectWrkChainFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(FeeTx)

	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "WRKChain Tx must be a FeeTx")
	}

	// check if it's a WRKChain Tx
	if !checkIsWrkChainTx(feeTx) {
		// ignore and move on to the next decorator in the chain
		return next(ctx, tx, simulate)
	}

	// Check fees sent in Tx
	err := checkWrkchainFees(feeTx)
	if err != nil {
		return ctx, err
	}

	// check if the WRKChain exists. We don't want to charge fees unnecessarily for re-registration
	err = checkWrkChainExists(ctx, wfd.wck, feeTx)
	if err != nil {
		return ctx, err
	}

	// check fee payer is WRKChain Owner
	err  = checkWrkChainOwnerFeePayer(feeTx)
	if err != nil {
		return ctx, err
	}

	// check sender has sufficient funds
	err = checkFeePayerHasFunds(ctx, wfd.ak, feeTx)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func checkIsWrkChainTx(tx FeeTx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch msg.(type) {
		case MsgRegisterWrkChain:
			return true
		case MsgRecordWrkChainBlock:
			return true
		}
	}
	return false
}

func checkWrkChainExists(ctx sdk.Context, wck Keeper, tx FeeTx) error {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch m := msg.(type) {
		case MsgRegisterWrkChain:
			if wck.IsWrkChainRegistered(ctx, m.WrkChainID) {
				return sdkerrors.Wrapf(ErrWrkChainAlreadyRegistered, "already registered: %s", m.WrkChainID)
			}
		case MsgRecordWrkChainBlock:
			if !wck.IsWrkChainRegistered(ctx, m.WrkChainID) {
				return sdkerrors.Wrapf(ErrWrkChainDoesNotExist, "does not exist: %s", m.WrkChainID)
			}
		}
	}

	return nil
}

func checkWrkchainFees(tx FeeTx) error {
	msgs := tx.GetMsgs()
	numMsgs := 0
	expectedFees := FeesBaseDenomination

	// go through Msgs wrapped in the Tx, and check for WRKChain messages
	for _, msg := range msgs {
		switch msg.(type) {
		case MsgRegisterWrkChain:
			expectedFees = expectedFees.Add(FeesWrkChainRegistrationCoin)
			numMsgs = numMsgs + 1
		case MsgRecordWrkChainBlock:
			expectedFees = expectedFees.Add(FeesWrkChainRecordHashCoin)
			numMsgs = numMsgs + 1
		}
	}

	totalFees := sdk.Coins{expectedFees}
	if tx.GetFee().IsAllLT(totalFees) {
		errMsg := fmt.Sprintf("numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(ErrInsufficientWrkChainFee, errMsg)
	}

	if tx.GetFee().IsAllGT(totalFees) {
		errMsg := fmt.Sprintf("numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(ErrTooMuchWrkChainFee, errMsg)
	}

	return nil
}

func checkWrkChainOwnerFeePayer(tx FeeTx) error {
	msgs := tx.GetMsgs()
	feePayer := tx.FeePayer()
	for _, msg := range msgs {
		switch m := msg.(type) {
		case MsgRegisterWrkChain:
			if !feePayer.Equals(m.Owner) {
				errMsg := fmt.Sprintf("Owner: %s, Fee Payer: %s", m.Owner, feePayer)
				return sdkerrors.Wrap(ErrFeePayerNotOwner, errMsg)
			}
		case MsgRecordWrkChainBlock:
			if !feePayer.Equals(m.Owner) {
				errMsg := fmt.Sprintf("Owner: %s, Fee Payer: %s", m.Owner, feePayer)
				return sdkerrors.Wrap(ErrFeePayerNotOwner, errMsg)
			}
		}
	}
	return nil
}

func checkFeePayerHasFunds(ctx sdk.Context, ak auth.AccountKeeper, tx FeeTx) error {
	feePayer := tx.FeePayer()
	feePayerAcc := ak.GetAccount(ctx, feePayer)
	blockTime := ctx.BlockHeader().Time
	coins := feePayerAcc.GetCoins()
	fees := tx.GetFee()

	if feePayerAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid fee: %s", fees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := feePayerAcc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	return nil
}
