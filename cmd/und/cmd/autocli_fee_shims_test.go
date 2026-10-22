package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFlatFee covers the fee arithmetic in isolation: register and record are
// flat lookups, purchase-storage multiplies the per-slot fee by the slot count
// taken from the second positional argument.
func TestFlatFee(t *testing.T) {
	// TestNet values at the time of writing.
	params := moduleFeeParams{
		denom:              "nund",
		feeRegister:        10000000000000,
		feeRecord:          1000000000,
		feePurchaseStorage: 5000000000,
	}

	for _, tc := range []struct {
		name   string
		kind   feeKind
		args   []string
		exp    string
		expErr string
	}{
		{name: "register", kind: feeRegister, args: []string{}, exp: "10000000000000nund"},
		{name: "record", kind: feeRecord, args: []string{"1"}, exp: "1000000000nund"},
		{name: "purchase one slot", kind: feePurchaseStorage, args: []string{"1", "1"}, exp: "5000000000nund"},
		{name: "purchase many slots", kind: feePurchaseStorage, args: []string{"1", "600000"}, exp: "3000000000000000nund"},
		{name: "purchase with no slot count", kind: feePurchaseStorage, args: []string{"1"}, expErr: "needs a slot count"},
		{name: "purchase with a junk slot count", kind: feePurchaseStorage, args: []string{"1", "lots"}, expErr: "invalid slot count"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := flatFee(params, tc.kind, tc.args)
			if tc.expErr != "" {
				require.ErrorContains(t, err, tc.expErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.exp, fee.String())
		})
	}
}

// TestFlatFee_NoOverflow checks the per-slot multiplication against a slot count
// that would overflow the int arithmetic the pre-autocli CLI used.
func TestFlatFee_NoOverflow(t *testing.T) {
	params := moduleFeeParams{denom: "nund", feePurchaseStorage: 5000000000}
	fee, err := flatFee(params, feePurchaseStorage, []string{"1", "18446744073709551615"})
	require.NoError(t, err)
	require.Equal(t, "92233720368547758075000000000nund", fee.String())
}

func TestFlatFee_EmptyDenom(t *testing.T) {
	_, err := flatFee(moduleFeeParams{feeRecord: 1}, feeRecord, nil)
	require.ErrorContains(t, err, "empty fee denomination")
}
