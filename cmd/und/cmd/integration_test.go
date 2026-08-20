package cmd_test

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/cmd/und/cmd"
)

// TestLegacyFlagAliases_OnLiveRootCmd boots the real NewRootCmd (which
// runs the full autocli pipeline) and asserts that every (cmdPath,
// legacy → canonical) pair declared by the Stage 7c alias map lands on
// a real autocli-generated cobra command, with both spellings resolving
// to the same *pflag.Flag. This is the test that would have caught a
// silent autocli-skip regression — e.g. if proto regen for some reason
// dropped the canonical flag, the alias map would still parse but the
// resulting cmd would have no submit-time / height / etc. flag at all.
func TestLegacyFlagAliases_OnLiveRootCmd(t *testing.T) {
	rootCmd := cmd.NewRootCmd()

	type lookup struct {
		path      string
		legacy    string
		canonical string
	}
	cases := []lookup{
		{"tx beacon record", "subtime", "submit-time"},
		{"tx wrkchain register", "base", "base-type"},
		{"tx wrkchain register", "genesis", "genesis-hash"},
		{"tx wrkchain record", "wc_height", "height"},
		{"tx wrkchain record", "block_hash", "block-hash"},
		{"tx wrkchain record", "parent_hash", "parent-hash"},
	}
	for _, c := range cases {
		t.Run(c.path+":"+c.legacy, func(t *testing.T) {
			leaf := findCmdByPath(t, rootCmd, c.path)
			modern := leaf.Flags().Lookup(c.canonical)
			require.NotNil(t, modern,
				"%s: autocli did not register the canonical flag --%s — autocli/proto wiring drift",
				c.path, c.canonical)
			legacy := leaf.Flags().Lookup(c.legacy)
			require.NotNil(t, legacy,
				"%s: legacy spelling --%s did not normalise — install walker did not reach this cmd",
				c.path, c.legacy)
			require.Same(t, modern, legacy,
				"%s: legacy --%s and canonical --%s resolved to different flags",
				c.path, c.legacy, c.canonical)
		})
	}
}

// TestAutocliPurchaseStorageAlias_OnLiveRootCmd guards the cmd-name
// alias path: 'purchase-storage' is canonical, 'purchase_storage' is
// the legacy spelling cobra resolves via Aliases.
func TestAutocliPurchaseStorageAlias_OnLiveRootCmd(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	for _, module := range []string{"beacon", "wrkchain"} {
		t.Run(module, func(t *testing.T) {
			canonical := findCmdByPath(t, rootCmd, "tx "+module+" purchase-storage")
			require.NotNil(t, canonical)
			legacy := findCmdByPath(t, rootCmd, "tx "+module+" purchase_storage")
			require.Same(t, canonical, legacy,
				"legacy 'purchase_storage' alias must resolve to the same cmd as 'purchase-storage'")
		})
	}
}

// TestCompatShims_OnLiveRootCmd asserts the autocli compat shims are
// actually attached to the live autocli-generated cmds: invoking PreRunE
// on `tx beacon record` (with no --submit-time) fills the flag from the
// process clock; invoking PreRunE on `tx enterprise process` rewrites
// args[1] from accept→accepted.
//
// --fees is supplied throughout so the auto-fee shim short-circuits: it
// shares the same PreRunE chain and would otherwise reach for a node.
func TestCompatShims_OnLiveRootCmd(t *testing.T) {
	rootCmd := cmd.NewRootCmd()

	t.Run("beacon record submit-time default fills", func(t *testing.T) {
		record := findCmdByPath(t, rootCmd, "tx beacon record")
		require.NotNil(t, record.PreRunE,
			"compat shim must be attached as PreRunE on tx beacon record")
		require.NoError(t, record.ParseFlags([]string{"--hash", "abc", "--fees", "1000000000nund"}))
		require.NoError(t, record.PreRunE(record, []string{"1"}))
		got, err := record.Flags().GetUint64("submit-time")
		require.NoError(t, err)
		require.NotZero(t, got,
			"submit-time must be populated from the process clock when omitted")
	})

	t.Run("beacon record submit-time explicit preserved", func(t *testing.T) {
		// Re-build to get a fresh flagset.
		root2 := cmd.NewRootCmd()
		record := findCmdByPath(t, root2, "tx beacon record")
		require.NoError(t, record.ParseFlags([]string{"--submit-time", "42", "--hash", "abc", "--fees", "1000000000nund"}))
		require.NoError(t, record.PreRunE(record, []string{"1"}))
		got, err := record.Flags().GetUint64("submit-time")
		require.NoError(t, err)
		require.Equal(t, uint64(42), got)
	})

	t.Run("enterprise process accept→status-accepted", func(t *testing.T) {
		process := findCmdByPath(t, rootCmd, "tx enterprise process")
		require.NotNil(t, process.PreRunE,
			"compat shim must be attached as PreRunE on tx enterprise process")
		args := []string{"42", "accept"}
		require.NoError(t, process.PreRunE(process, args))
		require.Equal(t, "status-accepted", args[1])
	})

	t.Run("enterprise process reject→status-rejected", func(t *testing.T) {
		process := findCmdByPath(t, rootCmd, "tx enterprise process")
		args := []string{"42", "reject"}
		require.NoError(t, process.PreRunE(process, args))
		require.Equal(t, "status-rejected", args[1])
	})

	t.Run("enterprise process accepted→status-accepted", func(t *testing.T) {
		process := findCmdByPath(t, rootCmd, "tx enterprise process")
		args := []string{"42", "accepted"}
		require.NoError(t, process.PreRunE(process, args))
		require.Equal(t, "status-accepted", args[1])
	})

	t.Run("enterprise process status-accepted passes through", func(t *testing.T) {
		process := findCmdByPath(t, rootCmd, "tx enterprise process")
		args := []string{"42", "status-accepted"}
		require.NoError(t, process.PreRunE(process, args))
		require.Equal(t, "status-accepted", args[1])
	})
}

// TestAutoFeeShims_OnLiveRootCmd guards the v1.13.0 regression: the pre-autocli
// cobra commands queried the module params and set --fees themselves, because
// BEACON and WRKChain Txs are charged flat governance-set fees that the ante
// chain matches exactly. The autocli migration dropped that, so any operator
// script that had never needed to pass --fees started failing. The shim must be
// attached to all six flat-fee commands, and must leave an explicitly supplied
// --fees alone (which is also what makes offline signing possible).
func TestAutoFeeShims_OnLiveRootCmd(t *testing.T) {
	for _, path := range []string{
		"tx beacon register",
		"tx beacon record",
		"tx beacon purchase-storage",
		"tx wrkchain register",
		"tx wrkchain record",
		"tx wrkchain purchase-storage",
	} {
		t.Run(path, func(t *testing.T) {
			leaf := findCmdByPath(t, cmd.NewRootCmd(), path)

			require.NotNil(t, leaf.Flags().Lookup("fees"),
				"%s: autocli did not register --fees", path)
			require.NotNil(t, leaf.PreRunE,
				"%s: auto-fee shim must be attached as PreRunE", path)

			require.NoError(t, leaf.ParseFlags([]string{"--fees", "1nund"}))
			require.NoError(t, leaf.PreRunE(leaf, []string{"1", "1"}),
				"%s: an explicit --fees must short-circuit the shim before it reaches for a node", path)

			got, err := leaf.Flags().GetString("fees")
			require.NoError(t, err)
			require.Equal(t, "1nund", got,
				"%s: explicitly supplied --fees must not be overwritten", path)
		})
	}
}

func findCmdByPath(t *testing.T, rootCmd *cobra.Command, path string) *cobra.Command {
	t.Helper()
	cur := rootCmd
	for _, part := range strings.Fields(path) {
		var next *cobra.Command
		for _, sub := range cur.Commands() {
			if sub.Name() == part {
				next = sub
				break
			}
			for _, alias := range sub.Aliases {
				if alias == part {
					next = sub
					break
				}
			}
			if next != nil {
				break
			}
		}
		require.NotNilf(t, next, "no subcommand %q under %q", part, cur.CommandPath())
		cur = next
	}
	return cur
}
