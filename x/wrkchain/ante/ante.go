package ante

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/wrkchain/exported"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// CorrectWrkChainFeeDecorator checks if the correct fees have been sent to pay for a
// WRKChain register/record hash Tx, and if the fee paying account has sufficient funds to pay.
// It first checks if the Tx contains any WRKChain Msgs, and if not, continues on to the next
// AnteHandler in the chain. If a WRKChain Msg is detected, it then:
//
// 1. Checks sufficient fees have been included in the Tx, via the --fees flag
// 2. Checks the fee payer is the WRKChain owner
// 3. Checks if the fee payer has sufficient funds in their account to pay for it, including any locked enterprise und
//
// If any of the checks fail, a suitable error is returned.
type CorrectWrkChainFeeDecorator struct {
	bankKeeper     BankKeeper
	accKeeper      AccountKeeper
	wrkchainKeeper WrkchainKeeper
	entKeeper      EnterpriseKeeper
}

func NewCorrectWrkChainFeeDecorator(bankKeeper BankKeeper, accKeeper AccountKeeper, wrkchainKeeper WrkchainKeeper, enterpriseKeeper EnterpriseKeeper) CorrectWrkChainFeeDecorator {
	return CorrectWrkChainFeeDecorator{
		bankKeeper:     bankKeeper,
		accKeeper:      accKeeper,
		wrkchainKeeper: wrkchainKeeper,
		entKeeper:      enterpriseKeeper,
	}
}

func (wfd CorrectWrkChainFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)

	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "WRKChain Tx must be a FeeTx")
	}

	// check if it's a WRKChain Tx
	if !exported.CheckIsWrkChainTx(feeTx) {
		// ignore and move on to the next decorator in the chain
		return next(ctx, tx, simulate)
	}

	// Check fees amount sent in Tx. Check during CheckTx
	if ctx.IsCheckTx() && !simulate {
		err := checkWrkchainFees(ctx, feeTx, wfd.wrkchainKeeper)
		if err != nil {
			return ctx, err
		}
	}

	// check sender has sufficient funds
	err := checkFeePayerHasFunds(ctx, wfd.bankKeeper, wfd.accKeeper, wfd.entKeeper, feeTx)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func checkWrkchainFees(ctx sdk.Context, tx sdk.FeeTx, wck WrkchainKeeper) error {
	msgs := tx.GetMsgs()
	numMsgs := 0
	expectedFees := wck.GetZeroFeeAsCoin(ctx)
	expectedFeeDenom := wck.GetParamDenom(ctx)
	hasFeeDenom := false

	for _, feeCoin := range tx.GetFee() {
		if feeCoin.Denom == expectedFeeDenom {
			hasFeeDenom = true
		}
	}

	if !hasFeeDenom {
		errMsg := fmt.Sprintf("incorrect fee denomination. expected %s", expectedFeeDenom)
		return sdkerrors.Wrap(exported.ErrIncorrectFeeDenomination, errMsg)
	}

	// go through Msgs wrapped in the Tx, and check for WRKChain messages
	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgRegisterWrkChain:
			expectedFees = expectedFees.Add(wck.GetRegistrationFeeAsCoin(ctx))
			numMsgs = numMsgs + 1
		case *types.MsgRecordWrkChainBlock:
			expectedFees = expectedFees.Add(wck.GetRecordFeeAsCoin(ctx))
			numMsgs = numMsgs + 1
		}
	}

	totalFees := sdk.Coins{expectedFees}
	if tx.GetFee().IsAllLT(totalFees) {
		errMsg := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(exported.ErrInsufficientWrkChainFee, errMsg)
	}

	if tx.GetFee().IsAllGT(totalFees) {
		errMsg := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(exported.ErrTooMuchWrkChainFee, errMsg)
	}

	return nil
}

func checkFeePayerHasFunds(ctx sdk.Context, bankKeeper BankKeeper, accKeeper AccountKeeper, ek EnterpriseKeeper, tx sdk.FeeTx) error {
	feePayer := tx.FeePayer()
	feePayerAcc := accKeeper.GetAccount(ctx, feePayer)
	//blockTime := ctx.BlockHeader().Time
	fees := tx.GetFee()

	if feePayerAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid fee: %s", fees)
	}

	coins := bankKeeper.GetAllBalances(ctx, feePayerAcc.GetAddress()) //feePayerAcc.GetCoins()

	potentialCoins := coins

	//get any locked enterprise FUND
	lockedUnd := ek.GetLockedUndAmountForAccount(ctx, feePayer)

	lockedUndCoins := sdk.NewCoins(lockedUnd)
	// include any locked FUND in potential coins. We need to do this because if these checks pass,
	// the locked FUND will be unlocked in the next decorator
	potentialCoins = potentialCoins.Add(lockedUndCoins...)

	// verify the account has enough funds to pay for fees, including any locked enterprise FUND
	_, hasNeg := potentialCoins.SafeSub(fees)
	if hasNeg {
		err := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %s", coins, potentialCoins, fees)
		return err
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := bankKeeper.SpendableCoins(ctx, feePayerAcc.GetAddress()) //feePayerAcc.SpendableCoins(blockTime)
	potentialSpendableCoins := spendableCoins

	// include any locked FUND in potential coins. We need to do this because if these checks pass,
	// the locked FUND will be unlocked in the next decorator
	potentialSpendableCoins = potentialSpendableCoins.Add(lockedUndCoins...)

	if _, hasNeg := potentialSpendableCoins.SafeSub(fees); hasNeg {
		err := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient spendable und to pay for fees. unlocked und: %s, including locked und: %s, fee: %s", spendableCoins, potentialSpendableCoins, fees)
		return err
	}

	return nil
}
