// Package feecheck holds the fee-payer solvency check shared by the BEACON and
// WRKChain fee decorators. Both modules charge flat, parameter-driven fees rather
// than gas-priced ones, and both must account for locked enterprise FUND that a
// later decorator will unlock to pay with, so the check itself is identical in
// each — only the fee denomination differs.
//
// The package deliberately imports nothing from the app-level ante package or from
// x/beacon and x/wrkchain, so that those may import it without an import cycle.
package feecheck

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CheckFeePayerHasFunds verifies the fee payer can cover feeDenom's share of the
// Tx fee, counting any locked enterprise FUND as available because the lock-checking
// decorator further down the chain unlocks it to pay with. Both the total balance
// and the spendable balance are checked, so vesting accounts cannot pay from coins
// they do not yet control.
func CheckFeePayerHasFunds(
	ctx sdk.Context,
	bankKeeper BankKeeper,
	accKeeper AccountKeeper,
	ek EnterpriseKeeper,
	feeDenom string,
	tx sdk.FeeTx,
) error {
	feePayer := tx.FeePayer()
	feePayerAcc := accKeeper.GetAccount(ctx, feePayer)
	fees := tx.GetFee()

	if feePayerAcc == nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	if !fees.IsValid() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid fee: %s", fees)
	}

	// AmountOf, not Coins.Find: Find returns a zero-value Coin when the denomination
	// is absent, and a zero-value Coin carries a nil math.Int that panics the moment
	// SafeSub sanitises it. AmountOf yields a genuine zero instead, so a Tx with no
	// fee (or a fee in some other denomination) falls through to the module's own
	// fee check and gets a legible error rather than a recovered nil dereference.
	fee := sdk.NewCoin(feeDenom, fees.AmountOf(feeDenom))

	// include any locked FUND in potential coins. We need to do this because if these
	// checks pass, the locked FUND will be unlocked in the next decorator
	lockedUndCoins := sdk.NewCoins(ek.GetLockedUndAmountForAccount(ctx, feePayer))

	coins := bankKeeper.GetAllBalances(ctx, feePayerAcc.GetAddress())
	potentialCoins := coins.Add(lockedUndCoins...)

	if _, hasNeg := potentialCoins.SafeSub(fee); hasNeg {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %s", coins, potentialCoins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := bankKeeper.SpendableCoins(ctx, feePayerAcc.GetAddress())
	potentialSpendableCoins := spendableCoins.Add(lockedUndCoins...)

	if _, hasNeg := potentialSpendableCoins.SafeSub(fee); hasNeg {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient spendable und to pay for fees. unlocked und: %s, including locked und: %s, fee: %s", spendableCoins, potentialSpendableCoins, fees)
	}

	return nil
}
