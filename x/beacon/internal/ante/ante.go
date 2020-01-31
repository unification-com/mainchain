package ante

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/unification-com/mainchain/x/beacon/exported"
	"github.com/unification-com/mainchain/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
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
	ak auth.AccountKeeper
	bk keeper.Keeper
	ek types.EnterpriseKeeper
}

func NewCorrectBeaconFeeDecorator(ak auth.AccountKeeper, beaconKeeper keeper.Keeper, enterpriseKeeper types.EnterpriseKeeper) CorrectBeaconFeeDecorator {
	return CorrectBeaconFeeDecorator{
		ak: ak,
		bk: beaconKeeper,
		ek: enterpriseKeeper,
	}
}

func (wfd CorrectBeaconFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(FeeTx)

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
		err := checkBeaconFees(ctx, feeTx, wfd.bk)
		if err != nil {
			return ctx, err
		}
	}

	// check sender has sufficient funds
	err := checkFeePayerHasFunds(ctx, wfd.ak, wfd.ek, feeTx)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func checkBeaconFees(ctx sdk.Context, tx FeeTx, bk keeper.Keeper) error {
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
		return sdkerrors.Wrap(types.ErrIncorrectFeeDenomination, errMsg)
	}

	// go through Msgs wrapped in the Tx, and check for BEACON messages
	for _, msg := range msgs {
		switch msg.(type) {
		case types.MsgRegisterBeacon:
			expectedFees = expectedFees.Add(bk.GetRegistrationFeeAsCoin(ctx))
			numMsgs = numMsgs + 1
		case types.MsgRecordBeaconTimestamp:
			expectedFees = expectedFees.Add(bk.GetRecordFeeAsCoin(ctx))
			numMsgs = numMsgs + 1
		}
	}

	totalFees := sdk.Coins{expectedFees}
	if tx.GetFee().IsAllLT(totalFees) {
		errMsg := fmt.Sprintf("insufficient fee to pay for beacon tx. numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(types.ErrInsufficientBeaconFee, errMsg)
	}

	if tx.GetFee().IsAllGT(totalFees) {
		errMsg := fmt.Sprintf("too much fee sent to pay for beacon tx. numMsgs in tx: %v, expected fees: %v, sent fees: %v", numMsgs, totalFees.String(), tx.GetFee())
		return sdkerrors.Wrap(types.ErrTooMuchBeaconFee, errMsg)
	}

	return nil
}

func checkFeePayerHasFunds(ctx sdk.Context, ak auth.AccountKeeper, ek types.EnterpriseKeeper, tx FeeTx) error {
	feePayer := tx.FeePayer()
	feePayerAcc := ak.GetAccount(ctx, feePayer)
	blockTime := ctx.BlockHeader().Time
	fees := tx.GetFee()

	if feePayerAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid fee: %s", fees)
	}

	coins := feePayerAcc.GetCoins()

	potentialCoins := coins

	//get any locked enterprise UND
	lockedUnd := ek.GetLockedUndAmountForAccount(ctx, feePayer)

	lockedUndCoins := sdk.NewCoins(lockedUnd)
	// include any locked UND in potential coins. We need to do this because if these checks pass,
	// the locked UND will be unlocked in the next decorator
	potentialCoins = potentialCoins.Add(lockedUndCoins...)

	// verify the account has enough funds to pay for fees, including any locked enterprise UND
	_, hasNeg := potentialCoins.SafeSub(fees)
	if hasNeg {
		err := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %s", coins, potentialCoins, fees)
		return err
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := feePayerAcc.SpendableCoins(blockTime)
	potentialSpendableCoins := spendableCoins

	// include any locked UND in potential coins. We need to do this because if these checks pass,
	// the locked UND will be unlocked in the next decorator
	potentialSpendableCoins = potentialSpendableCoins.Add(lockedUndCoins...)

	if _, hasNeg := potentialSpendableCoins.SafeSub(fees); hasNeg {
		err := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient spendable und to pay for fees. unlocked und: %s, including locked und: %s, fee: %s", spendableCoins, potentialSpendableCoins, fees)
		return err
	}

	return nil
}
