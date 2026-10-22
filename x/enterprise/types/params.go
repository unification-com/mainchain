package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxEntSigners is the upper bound on the number of authorised enterprise
// signers. Limits per-PO decision-iteration cost and per-tally authorisation
// scan.
const MaxEntSigners = 32

// MaxPOsPerBlock caps the number of raised + accepted purchase orders the
// BeginBlocker may process in a single block. Bounds the ABCI BeginBlock work
// per block — without it, a backlog of queued POs would block production
// proportionally to the queue size (which has no other bound). Leftovers
// carry forward to the next block.
const MaxPOsPerBlock = 100

// MaxOpenPOsPerPurchaser caps the number of simultaneously raised (not yet
// processed) purchase orders any single whitelisted address can have open.
// Defends against a whitelisted address spamming the raised queue and
// inflating BeginBlocker iteration cost. Once the cap is reached, the
// purchaser must wait for their existing POs to be processed before raising
// new ones.
const MaxOpenPOsPerPurchaser = 50

// MaxDecisionTimeLimit caps the gov-tunable DecisionTimeLimit (seconds).
// 30 days. Beyond this, stale-PO auto-reject ceases to be useful state
// hygiene and the raised queue can grow without natural cleanup.
const MaxDecisionTimeLimit uint64 = 30 * 24 * 60 * 60

func NewParams(denom string, minAccepts uint64, decisionLimit uint64, entSigners string) Params {
	return Params{
		EntSigners:        entSigners,
		Denom:             denom,
		MinAccepts:        minAccepts,
		DecisionTimeLimit: decisionLimit,
	}
}

// default enterprise FUND module parameters.
// Note: DefaultParams returns a non-functional configuration — EntSigners is
// set to the zero address, so no purchase order can be processed under the
// default. Production chains must set EntSigners via genesis or
// MsgUpdateParams. The legacy/devnet defaults are kept here only so that
// fresh chains boot without InitGenesis failing on param validation.
func DefaultParams() Params {
	return Params{
		// Zero-address placeholder. A production chain MUST override this via
		// genesis or MsgUpdateParams before purchase orders can be processed.
		EntSigners:        "und1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq5x8kpm",
		Denom:             sdk.DefaultBondDenom,
		MinAccepts:        1,
		// 23h 30min default. Slightly under one day. Production chains
		// typically set this via genesis to a value that matches their
		// gov-cycle timing. Kept verbatim from pre-vaxildan defaults rather
		// than "corrected" to 86400 (1 day) to avoid silently changing
		// behaviour on chains that boot with defaults.
		DecisionTimeLimit: 84600,
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateDenom(p.Denom); err != nil {
		return err
	}
	if err := validateMinAccepts(p.MinAccepts); err != nil {
		return err
	}
	if err := validateDecisionLimit(p.DecisionTimeLimit); err != nil {
		return err
	}
	signers, err := parseAndValidateEntSigners(p.EntSigners)
	if err != nil {
		return err
	}
	if len(signers) < int(p.MinAccepts) {
		return fmt.Errorf("number of authorised signers (%d) must be >= MinAccepts (%d)",
			len(signers), p.MinAccepts)
	}
	return nil
}

func validateDenom(v string) error {
	if strings.TrimSpace(v) == "" {
		return errors.New("denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}
	return nil
}

func validateMinAccepts(v uint64) error {
	if v == 0 {
		return fmt.Errorf("min accepts must be positive: %d", v)
	}
	return nil
}

func validateDecisionLimit(v uint64) error {
	if v == 0 {
		return fmt.Errorf("decision limit must be positive: %d", v)
	}
	if v > MaxDecisionTimeLimit {
		return fmt.Errorf("decision limit %d exceeds hard cap %d (seconds)", v, MaxDecisionTimeLimit)
	}
	return nil
}

// parseAndValidateEntSigners trims, splits, dedupes and validates the
// comma-separated EntSigners string. Caps the count at MaxEntSigners to bound
// per-PO authorisation-scan cost. Returns the cleaned slice for use by callers
// that want addresses without re-splitting.
func parseAndValidateEntSigners(s string) ([]string, error) {
	if s == "" {
		return nil, errors.New("must have at least one ent signer")
	}
	parts := strings.Split(s, ",")
	seen := make(map[string]struct{}, len(parts))
	out := make([]string, 0, len(parts))
	for i, raw := range parts {
		addr := strings.TrimSpace(raw)
		if addr == "" {
			return nil, fmt.Errorf("ent signer entry %d is empty", i)
		}
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return nil, fmt.Errorf("ent signer entry %d (%q) is not a valid address: %w", i, addr, err)
		}
		if _, dup := seen[addr]; dup {
			return nil, fmt.Errorf("ent signer entry %d (%q) is a duplicate", i, addr)
		}
		seen[addr] = struct{}{}
		out = append(out, addr)
	}
	if len(out) > MaxEntSigners {
		return nil, fmt.Errorf("ent signers count %d exceeds cap %d", len(out), MaxEntSigners)
	}
	return out, nil
}
