package ante

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/beacon/exported"
	"github.com/unification-com/mainchain/x/beacon/types"
)

// CorrectBeaconFeeDecorator checks if the correct fees have been sent to pay for a
// BEACON register/record hash Tx, and if the fee paying account has sufficient funds to pay.
// It first checks if the Tx contains any BEACON Msgs, and if not, continues on to the next
// AnteHandler in the chain. If a BEACON Msg is detected, it then:
//
// 1. Checks sufficient fees have been included in the Tx, via the --fees flag
// 2. Checks the fee payer is the BEACON owner
// 3. Checks if the fee payer has sufficient funds in their account to pay for it, including any locked enterprise und
//
// If any of the checks fail, a suitable error is returned.
type CorrectBeaconFeeDecorator struct {
	bankKeeper   BankKeeper
	accKeeper    AccountKeeper
	beaconKeeper BeaconKeeper
	entKeeper    EnterpriseKeeper
}

func NewCorrectBeaconFeeDecorator(bankKeeper BankKeeper, accKeeper AccountKeeper, beaconKeeper BeaconKeeper, enterpriseKeeper EnterpriseKeeper) CorrectBeaconFeeDecorator {
	return CorrectBeaconFeeDecorator{
		bankKeeper:   bankKeeper,
		accKeeper:    accKeeper,
		beaconKeeper: beaconKeeper,
		entKeeper:    enterpriseKeeper,
	}
}

func (wfd CorrectBeaconFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)

	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "BEACON Tx must be a FeeTx")
	}

	// check if it's a BEACON Tx
	if !exported.CheckIsBeaconTx(feeTx) {
		// ignore and move on to the next decorator in the chain
		return next(ctx, tx, simulate)
	}

	// Check fees amount sent in Tx. Check during CheckTx
	if ctx.IsCheckTx() && !simulate {
		err := checkBeaconFees(ctx, feeTx, wfd.beaconKeeper)
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

func checkBeaconFees(ctx sdk.Context, tx sdk.FeeTx, bk BeaconKeeper) error {
	msgs := tx.GetMsgs()
	numMsgs := 0
	expectedFees := bk.GetZeroFeeAsCoin(ctx)
	expectedFeeDenom := bk.GetParamDenom(ctx)
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

	// go through Msgs wrapped in the Tx, and check for BEACON messages
	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgRegisterBeacon:
			expectedFees = expectedFees.Add(bk.GetRegistrationFeeAsCoin(ctx))
			numMsgs = numMsgs + 1
		case *types.MsgRecordBeaconTimestamp:
			expectedFees = expectedFees.Add(bk.GetRecordFeeAsCoin(ctx))
			numMsgs = numMsgs + 1
		}
	}

	totalFees := sdk.Coins{expectedFees}
	if tx.GetFee().IsAllLT(totalFees) {
		errMsg := fmt.Sprintf("insufficient fee to pay for beacon tx. numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(exported.ErrInsufficientBeaconFee, errMsg)
	}

	if tx.GetFee().IsAllGT(totalFees) {
		errMsg := fmt.Sprintf("too much fee sent to pay for beacon tx. numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(exported.ErrTooMuchBeaconFee, errMsg)
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
