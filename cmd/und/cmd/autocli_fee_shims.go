package cmd

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
)

// BEACON and WRKChain Txs are charged flat, governance-set fees rather than
// gas-priced ones, and the ante chain rejects a Tx whose fee is not exactly the
// expected amount. The pre-autocli cobra commands therefore queried the module
// params and set --fees themselves; the autocli migration dropped that, which
// silently broke every operator script that never passed --fees. These shims put
// it back.
//
// Unlike the pre-autocli commands, an explicitly supplied --fees is left alone
// rather than overwritten. That keeps the fee-less scripts working exactly as
// they did before, and additionally allows offline / --generate-only signing,
// which cannot reach a node to read the params.

// feeKind selects which of the three flat fees a command has to pay.
type feeKind int

const (
	feeRegister feeKind = iota
	feeRecord
	feePurchaseStorage
)

// moduleFeeParams is the slice of a module's params that fee calculation needs.
// BEACON and WRKChain declare the same four fields on separate generated types,
// so they are normalised here and the arithmetic is written once.
type moduleFeeParams struct {
	denom              string
	feeRegister        uint64
	feeRecord          uint64
	feePurchaseStorage uint64
}

func beaconFeeParams(clientCtx client.Context) (moduleFeeParams, error) {
	res, err := beacontypes.NewQueryClient(clientCtx).Params(context.Background(), &beacontypes.QueryParamsRequest{})
	if err != nil {
		return moduleFeeParams{}, err
	}
	return moduleFeeParams{
		denom:              res.Params.Denom,
		feeRegister:        res.Params.FeeRegister,
		feeRecord:          res.Params.FeeRecord,
		feePurchaseStorage: res.Params.FeePurchaseStorage,
	}, nil
}

func wrkchainFeeParams(clientCtx client.Context) (moduleFeeParams, error) {
	res, err := wrkchaintypes.NewQueryClient(clientCtx).Params(context.Background(), &wrkchaintypes.QueryParamsRequest{})
	if err != nil {
		return moduleFeeParams{}, err
	}
	return moduleFeeParams{
		denom:              res.Params.Denom,
		feeRegister:        res.Params.FeeRegister,
		feeRecord:          res.Params.FeeRecord,
		feePurchaseStorage: res.Params.FeePurchaseStorage,
	}, nil
}

// flatFee resolves the exact fee a command owes. purchase-storage is the only
// kind that scales: it charges per slot, and the slot count is the second
// positional argument.
func flatFee(params moduleFeeParams, kind feeKind, args []string) (sdk.Coin, error) {
	var amount math.Int

	switch kind {
	case feeRegister:
		amount = math.NewIntFromUint64(params.feeRegister)
	case feeRecord:
		amount = math.NewIntFromUint64(params.feeRecord)
	case feePurchaseStorage:
		if len(args) < 2 {
			return sdk.Coin{}, fmt.Errorf("purchase-storage needs a slot count to price the fee")
		}
		numSlots, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return sdk.Coin{}, fmt.Errorf("invalid slot count %q: %w", args[1], err)
		}
		amount = math.NewIntFromUint64(params.feePurchaseStorage).Mul(math.NewIntFromUint64(numSlots))
	default:
		return sdk.Coin{}, fmt.Errorf("unknown fee kind %d", kind)
	}

	if params.denom == "" {
		return sdk.Coin{}, fmt.Errorf("module params returned an empty fee denomination")
	}

	return sdk.NewCoin(params.denom, amount), nil
}

// moduleFeeShim builds the PreRunE hook that reads the module params off the
// chain and pins --fees to the exact amount the ante chain will demand.
func moduleFeeShim(fetch func(client.Context) (moduleFeeParams, error), kind feeKind) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Leave an explicitly supplied fee alone, whichever form it takes. --gas-prices matters as
		// much as --fees here: the SDK refuses a Tx carrying both ("cannot provide both fees and
		// gas prices"), so injecting a fee on top of an operator's --gas-prices would fail the Tx
		// client-side instead of letting the chain judge it.
		if cmd.Flags().Changed(flags.FlagFees) || cmd.Flags().Changed(flags.FlagGasPrices) {
			return nil
		}

		// The root cmd's PersistentPreRunE installs the client context before any
		// PreRunE runs, so this only bites a cmd invoked outside Execute() — which
		// is worth an error rather than the nil-context panic cobra would give.
		if cmd.Context() == nil {
			return fmt.Errorf("no client context available to look up the %s fee; pass --fees explicitly", cmd.CommandPath())
		}

		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		params, err := fetch(clientCtx)
		if err != nil {
			return fmt.Errorf("could not read module fee params to set --fees (pass --fees explicitly to skip this lookup): %w", err)
		}

		fee, err := flatFee(params, kind, args)
		if err != nil {
			return err
		}

		return cmd.Flags().Set(flags.FlagFees, fee.String())
	}
}

// autoFeeCommands lists every Tx command that the chain charges a flat fee for.
func autoFeeCommands() map[string]func(*cobra.Command, []string) error {
	return map[string]func(*cobra.Command, []string) error{
		"tx beacon register":           moduleFeeShim(beaconFeeParams, feeRegister),
		"tx beacon record":             moduleFeeShim(beaconFeeParams, feeRecord),
		"tx beacon purchase-storage":   moduleFeeShim(beaconFeeParams, feePurchaseStorage),
		"tx wrkchain register":         moduleFeeShim(wrkchainFeeParams, feeRegister),
		"tx wrkchain record":           moduleFeeShim(wrkchainFeeParams, feeRecord),
		"tx wrkchain purchase-storage": moduleFeeShim(wrkchainFeeParams, feePurchaseStorage),
	}
}

// installAutoFeeShims attaches the fee shims to the autocli-generated cmd tree.
func installAutoFeeShims(rootCmd *cobra.Command) {
	for path, shim := range autoFeeCommands() {
		if cmd := findCommand(rootCmd, path); cmd != nil {
			chainPreRunE(cmd, shim)
		}
	}
}
