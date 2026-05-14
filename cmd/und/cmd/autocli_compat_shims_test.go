package cmd

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func withFixedNow(t *testing.T, fixed uint64) {
	t.Helper()
	prev := nowUnixSec
	nowUnixSec = func() uint64 { return fixed }
	t.Cleanup(func() { nowUnixSec = prev })
}

func TestBeaconRecordFillSubmitTime_DefaultsWhenOmitted(t *testing.T) {
	withFixedNow(t, 1700000000)
	cmd := &cobra.Command{Use: "record"}
	cmd.Flags().Uint64("submit-time", 0, "")
	cmd.Flags().String("hash", "", "")

	require.NoError(t, cmd.ParseFlags([]string{"--hash", "abc"}))
	require.NoError(t, beaconRecordFillSubmitTime(cmd, nil))

	got, err := cmd.Flags().GetUint64("submit-time")
	require.NoError(t, err)
	require.Equal(t, uint64(1700000000), got)
}

func TestBeaconRecordFillSubmitTime_PreservesExplicitValue(t *testing.T) {
	withFixedNow(t, 1700000000)
	cmd := &cobra.Command{Use: "record"}
	cmd.Flags().Uint64("submit-time", 0, "")

	require.NoError(t, cmd.ParseFlags([]string{"--submit-time", "42"}))
	require.NoError(t, beaconRecordFillSubmitTime(cmd, nil))

	got, err := cmd.Flags().GetUint64("submit-time")
	require.NoError(t, err)
	require.Equal(t, uint64(42), got,
		"shim must never overwrite an explicit user-supplied --submit-time")
}

func TestBeaconRecordFillSubmitTime_PreservesExplicitZero(t *testing.T) {
	// If the operator explicitly passes --submit-time 0, the shim must
	// not paper over it. The server-side ValidateBasic rejection then
	// surfaces unchanged — which is the legacy behaviour too: legacy
	// cli only defaulted when the flag was *omitted*, not when set to 0.
	cmd := &cobra.Command{Use: "record"}
	cmd.Flags().Uint64("submit-time", 0, "")

	require.NoError(t, cmd.ParseFlags([]string{"--submit-time", "0"}))
	require.NoError(t, beaconRecordFillSubmitTime(cmd, nil))

	got, err := cmd.Flags().GetUint64("submit-time")
	require.NoError(t, err)
	require.Equal(t, uint64(0), got)
}

func TestEnterpriseProcessAliasDecisionToken_LegacyTokens(t *testing.T) {
	legacyFlagWarnedOnce = sync.Map{}
	var buf bytes.Buffer
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &buf
	defer func() { legacyFlagWarnSink = prev }()

	cases := []struct{ in, want string }{
		{"accept", "status-accepted"},
		{"accepted", "status-accepted"},
		{"reject", "status-rejected"},
		{"rejected", "status-rejected"},
	}
	for _, c := range cases {
		args := []string{"42", c.in}
		require.NoError(t, enterpriseProcessAliasDecisionToken(nil, args))
		require.Equal(t, c.want, args[1],
			"legacy token %q must be rewritten to %q in-place", c.in, c.want)
	}

	out := buf.String()
	require.Contains(t, out, "value 'accept'",
		"decision-token deprecation must be phrased as a value, not a flag spelling")
	require.Contains(t, out, "value 'reject'")
	require.Contains(t, out, "'status-accepted'")
	require.Contains(t, out, "'status-rejected'")
	require.Contains(t, out, "tx enterprise process")
	require.NotContains(t, out, "--accept",
		"decision tokens are positionals; warning must not look like a flag")
}

func TestEnterpriseProcessAliasDecisionToken_ModernTokensUnchanged(t *testing.T) {
	legacyFlagWarnedOnce = sync.Map{}
	var buf bytes.Buffer
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &buf
	defer func() { legacyFlagWarnSink = prev }()

	for _, tok := range []string{"status-accepted", "status-rejected", "anythingelse"} {
		args := []string{"42", tok}
		require.NoError(t, enterpriseProcessAliasDecisionToken(nil, args))
		require.Equal(t, tok, args[1],
			"non-legacy decision token %q must pass through unchanged", tok)
	}

	require.Empty(t, buf.String(),
		"non-legacy decision tokens must not trigger a deprecation notice")
}

func TestEnterpriseProcessAliasDecisionToken_ShortArgsNoOp(t *testing.T) {
	// Defensive: if the operator omits the decision positional entirely
	// the shim must not panic; autocli's own positional validation will
	// surface the usage error.
	require.NoError(t, enterpriseProcessAliasDecisionToken(nil, nil))
	require.NoError(t, enterpriseProcessAliasDecisionToken(nil, []string{"only-one-arg"}))
}

func TestEnterpriseProcessAliasDecisionToken_WarnOnceAcrossInvocations(t *testing.T) {
	legacyFlagWarnedOnce = sync.Map{}
	var buf bytes.Buffer
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &buf
	defer func() { legacyFlagWarnSink = prev }()

	for i := 0; i < 5; i++ {
		args := []string{"1", "accept"}
		require.NoError(t, enterpriseProcessAliasDecisionToken(nil, args))
		require.Equal(t, "status-accepted", args[1])
	}
	require.Equal(t, 1, strings.Count(buf.String(), "deprecated"),
		"repeated legacy-token use must emit at most one deprecation notice per process")
}

func TestChainPreRunE_AppliesShimThenExisting(t *testing.T) {
	order := []string{}
	existing := func(*cobra.Command, []string) error {
		order = append(order, "existing")
		return nil
	}
	shim := func(*cobra.Command, []string) error {
		order = append(order, "shim")
		return nil
	}

	cmd := &cobra.Command{Use: "x", PreRunE: existing}
	chainPreRunE(cmd, shim)
	require.NoError(t, cmd.PreRunE(cmd, nil))
	require.Equal(t, []string{"shim", "existing"}, order,
		"shim must run first so its mutations are visible to the existing hook")
}

func TestChainPreRunE_ShimErrorShortCircuits(t *testing.T) {
	existingCalled := false
	existing := func(*cobra.Command, []string) error {
		existingCalled = true
		return nil
	}
	boom := errors.New("boom")
	shim := func(*cobra.Command, []string) error { return boom }

	cmd := &cobra.Command{Use: "x", PreRunE: existing}
	chainPreRunE(cmd, shim)
	require.ErrorIs(t, cmd.PreRunE(cmd, nil), boom)
	require.False(t, existingCalled,
		"shim error must short-circuit the existing hook")
}

func TestChainPreRunE_NoExistingHook(t *testing.T) {
	called := false
	cmd := &cobra.Command{Use: "x"}
	chainPreRunE(cmd, func(*cobra.Command, []string) error {
		called = true
		return nil
	})
	require.NoError(t, cmd.PreRunE(cmd, nil))
	require.True(t, called)
}
