package feecheck_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	protov2 "google.golang.org/protobuf/proto"

	"github.com/unification-com/mainchain/ante/feecheck"
)

const feeDenom = "nund"

var payer = sdk.AccAddress("feepayer____________")

type stubAccountKeeper struct{ acc sdk.AccountI }

func (s stubAccountKeeper) GetParams(context.Context) authtypes.Params              { return authtypes.Params{} }
func (s stubAccountKeeper) GetAccount(context.Context, sdk.AccAddress) sdk.AccountI { return s.acc }
func (s stubAccountKeeper) SetAccount(context.Context, sdk.AccountI)                {}
func (s stubAccountKeeper) GetModuleAddress(string) sdk.AccAddress                  { return nil }
func (s stubAccountKeeper) AddressCodec() address.Codec                             { return nil }

type stubBankKeeper struct{ balance, spendable sdk.Coins }

func (s stubBankKeeper) GetAllBalances(context.Context, sdk.AccAddress) sdk.Coins { return s.balance }
func (s stubBankKeeper) SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins { return s.spendable }

type stubEnterpriseKeeper struct{ locked sdk.Coin }

func (s stubEnterpriseKeeper) GetLockedUndAmountForAccount(sdk.Context, sdk.AccAddress) sdk.Coin {
	return s.locked
}

type stubFeeTx struct{ fee sdk.Coins }

func (s stubFeeTx) GetMsgs() []sdk.Msg                    { return nil }
func (s stubFeeTx) GetMsgsV2() ([]protov2.Message, error) { return nil, nil }
func (s stubFeeTx) GetGas() uint64                        { return 0 }
func (s stubFeeTx) GetFee() sdk.Coins                     { return s.fee }
func (s stubFeeTx) FeePayer() []byte                      { return payer }
func (s stubFeeTx) FeeGranter() []byte                    { return nil }

func coins(amt int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(feeDenom, amt)) }

// TestCheckFeePayerHasFunds_AbsentFeeDenom is the regression test for the nil
// pointer dereference that recovered as `runtime error: invalid memory address`
// inside the BEACON/WRKChain fee decorators. Coins.Find returns a zero-value Coin
// when the denomination is absent, whose nil math.Int panics as soon as SafeSub
// sanitises it; AmountOf returns a genuine zero instead. A Tx that carries no fee
// (as `--gas=auto` produces during simulation) must therefore pass the solvency
// check quietly and let the module's own fee check report the real problem.
func TestCheckFeePayerHasFunds_AbsentFeeDenom(t *testing.T) {
	for _, tc := range []struct {
		name string
		fee  sdk.Coins
	}{
		{"no fee at all", sdk.Coins{}},
		{"nil fee", nil},
		{"fee in another denomination", sdk.NewCoins(sdk.NewInt64Coin("foo", 1))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				err := feecheck.CheckFeePayerHasFunds(
					sdk.Context{},
					stubBankKeeper{balance: coins(10), spendable: coins(10)},
					stubAccountKeeper{acc: authtypes.NewBaseAccountWithAddress(payer)},
					stubEnterpriseKeeper{locked: sdk.NewInt64Coin(feeDenom, 0)},
					feeDenom,
					stubFeeTx{fee: tc.fee},
				)
				require.NoError(t, err)
			})
		})
	}
}

func TestCheckFeePayerHasFunds(t *testing.T) {
	acc := authtypes.NewBaseAccountWithAddress(payer)

	for _, tc := range []struct {
		name      string
		balance   sdk.Coins
		spendable sdk.Coins
		locked    sdk.Coin
		fee       sdk.Coins
		expErr    string
	}{
		{
			name:      "liquid balance covers the fee",
			balance:   coins(100),
			spendable: coins(100),
			locked:    sdk.NewInt64Coin(feeDenom, 0),
			fee:       coins(100),
		},
		{
			name:      "locked eFUND makes up the shortfall",
			balance:   coins(1),
			spendable: coins(1),
			locked:    sdk.NewInt64Coin(feeDenom, 999),
			fee:       coins(1000),
		},
		{
			name:      "short even with locked eFUND",
			balance:   coins(1),
			spendable: coins(1),
			locked:    sdk.NewInt64Coin(feeDenom, 1),
			fee:       coins(1000),
			expErr:    "insufficient und to pay for fees",
		},
		{
			name:      "balance covers it but the coins are not spendable",
			balance:   coins(1000),
			spendable: coins(1),
			locked:    sdk.NewInt64Coin(feeDenom, 0),
			fee:       coins(1000),
			expErr:    "insufficient spendable und to pay for fees",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := feecheck.CheckFeePayerHasFunds(
				sdk.Context{},
				stubBankKeeper{balance: tc.balance, spendable: tc.spendable},
				stubAccountKeeper{acc: acc},
				stubEnterpriseKeeper{locked: tc.locked},
				feeDenom,
				stubFeeTx{fee: tc.fee},
			)
			if tc.expErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.expErr)
		})
	}
}

func TestCheckFeePayerHasFunds_UnknownAccount(t *testing.T) {
	err := feecheck.CheckFeePayerHasFunds(
		sdk.Context{},
		stubBankKeeper{},
		stubAccountKeeper{acc: nil},
		stubEnterpriseKeeper{locked: sdk.NewInt64Coin(feeDenom, 0)},
		feeDenom,
		stubFeeTx{fee: coins(1)},
	)
	require.ErrorContains(t, err, "does not exist")
}
