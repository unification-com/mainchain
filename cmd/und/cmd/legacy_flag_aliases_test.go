package cmd

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// resetLegacyFlagWarnerForTest clears the one-shot tracking state. Tests
// that assert on warning output must call this in setup so cases don't
// poison each other.
func resetLegacyFlagWarnerForTest() {
	legacyFlagWarnedOnce = sync.Map{}
}

func TestLegacyFlagAliasMapEntries(t *testing.T) {
	// Sanity check: the registry must list every cmd/flag pair documented
	// in the Stage 7c plan. If any pair is removed or renamed, this test
	// fails loudly so the smoke checklist + release notes stay in sync.
	want := map[string]map[string]string{
		"tx beacon record": {
			"subtime": "submit-time",
		},
		"tx wrkchain register": {
			"base":    "base-type",
			"genesis": "genesis-hash",
		},
		"tx wrkchain record": {
			"wc_height":   "height",
			"block_hash":  "block-hash",
			"parent_hash": "parent-hash",
		},
	}
	got := legacyFlagAliasMap()
	require.Equal(t, len(want), len(got), "alias map cmd count drifted")
	for path, wantPairs := range want {
		aliases, ok := got[path]
		require.True(t, ok, "missing cmd path %q", path)
		require.Equal(t, len(wantPairs), len(aliases), "alias count drifted for %q", path)
		for _, a := range aliases {
			canonical, ok := wantPairs[a.legacy]
			require.True(t, ok, "unexpected legacy flag %q under %q", a.legacy, path)
			require.Equal(t, canonical, a.canonical, "canonical drifted for %q under %q", a.legacy, path)
		}
	}
}

func TestWarnLegacyFlagOnce_EmitsOnceThenSilent(t *testing.T) {
	resetLegacyFlagWarnerForTest()
	var buf bytes.Buffer
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &buf
	defer func() { legacyFlagWarnSink = prev }()

	warnLegacyFlagOnce("tx beacon record", "subtime", "submit-time")
	warnLegacyFlagOnce("tx beacon record", "subtime", "submit-time")
	warnLegacyFlagOnce("tx beacon record", "subtime", "submit-time")

	out := buf.String()
	require.Equal(t, 1, strings.Count(out, "deprecated"),
		"expected exactly one deprecation line for repeated legacy flag use, got: %q", out)
	require.Contains(t, out, "--subtime")
	require.Contains(t, out, "--submit-time")
	require.Contains(t, out, "tx beacon record")
}

func TestWarnLegacyFlagOnce_PerLegacyPerCmd(t *testing.T) {
	// Different (cmdPath, legacy) pairs each get their own one-shot slot.
	resetLegacyFlagWarnerForTest()
	var buf bytes.Buffer
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &buf
	defer func() { legacyFlagWarnSink = prev }()

	warnLegacyFlagOnce("tx wrkchain record", "wc_height", "height")
	warnLegacyFlagOnce("tx wrkchain record", "block_hash", "block-hash")
	warnLegacyFlagOnce("tx wrkchain register", "base", "base-type")
	// Repeats of the first two — should not emit again.
	warnLegacyFlagOnce("tx wrkchain record", "wc_height", "height")
	warnLegacyFlagOnce("tx wrkchain record", "block_hash", "block-hash")

	out := buf.String()
	require.Equal(t, 3, strings.Count(out, "deprecated"),
		"expected exactly three distinct deprecation lines, got: %q", out)
}

// stubTree constructs a minimal cobra tree mirroring the paths the
// installer walks, with one canonical flag per (cmd, alias) tuple drawn
// from the alias map. Each canonical flag is a no-op string so we can
// inspect Lookup() resolution.
func stubTree() *cobra.Command {
	root := &cobra.Command{Use: "und"}

	tx := &cobra.Command{Use: "tx"}
	root.AddCommand(tx)

	beacon := &cobra.Command{Use: "beacon"}
	tx.AddCommand(beacon)
	beaconRecord := &cobra.Command{Use: "record"}
	beaconRecord.Flags().String("submit-time", "", "")
	beacon.AddCommand(beaconRecord)

	wrkchain := &cobra.Command{Use: "wrkchain"}
	tx.AddCommand(wrkchain)

	wrkchainRegister := &cobra.Command{Use: "register"}
	wrkchainRegister.Flags().String("base-type", "", "")
	wrkchainRegister.Flags().String("genesis-hash", "", "")
	wrkchain.AddCommand(wrkchainRegister)

	wrkchainRecord := &cobra.Command{Use: "record"}
	wrkchainRecord.Flags().String("height", "", "")
	wrkchainRecord.Flags().String("block-hash", "", "")
	wrkchainRecord.Flags().String("parent-hash", "", "")
	wrkchain.AddCommand(wrkchainRecord)

	return root
}

func TestFindCommand_ResolvesNestedPaths(t *testing.T) {
	root := stubTree()
	for _, path := range []string{
		"tx beacon record",
		"tx wrkchain register",
		"tx wrkchain record",
	} {
		got := findCommand(root, path)
		require.NotNil(t, got, "findCommand returned nil for %q", path)
		require.Equal(t, lastSegment(path), got.Name(), "wrong cmd resolved for %q", path)
	}
	require.Nil(t, findCommand(root, "tx beacon nope"))
	require.Nil(t, findCommand(root, "no-such tree"))
}

func TestFindCommand_HonoursAliases(t *testing.T) {
	root := stubTree()
	beacon := findCommand(root, "tx beacon")
	require.NotNil(t, beacon)
	purchase := &cobra.Command{Use: "purchase-storage", Aliases: []string{"purchase_storage"}}
	beacon.AddCommand(purchase)

	require.Equal(t, purchase, findCommand(root, "tx beacon purchase-storage"))
	require.Equal(t, purchase, findCommand(root, "tx beacon purchase_storage"))
}

func TestInstallLegacyFlagAliases_BothSpellingsResolveToSameFlag(t *testing.T) {
	resetLegacyFlagWarnerForTest()
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &bytes.Buffer{}
	defer func() { legacyFlagWarnSink = prev }()

	root := stubTree()
	installLegacyFlagAliases(root)

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
		cmd := findCommand(root, c.path)
		require.NotNil(t, cmd, "%s: cmd missing", c.path)
		modern := cmd.Flags().Lookup(c.canonical)
		require.NotNil(t, modern, "%s: canonical flag --%s missing", c.path, c.canonical)
		legacy := cmd.Flags().Lookup(c.legacy)
		require.NotNil(t, legacy, "%s: legacy flag --%s did not resolve", c.path, c.legacy)
		require.Same(t, modern, legacy,
			"%s: legacy --%s and canonical --%s resolved to different *Flag", c.path, c.legacy, c.canonical)
	}
}

func TestInstallLegacyFlagAliases_ParseAcceptsLegacySpelling(t *testing.T) {
	resetLegacyFlagWarnerForTest()
	var buf bytes.Buffer
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &buf
	defer func() { legacyFlagWarnSink = prev }()

	root := stubTree()
	installLegacyFlagAliases(root)

	cmd := findCommand(root, "tx wrkchain record")
	require.NotNil(t, cmd)

	require.NoError(t, cmd.ParseFlags([]string{
		"--wc_height", "42",
		"--block_hash", "deadbeef",
		"--parent-hash", "cafef00d",
	}))

	height, err := cmd.Flags().GetString("height")
	require.NoError(t, err)
	require.Equal(t, "42", height)

	blockHash, err := cmd.Flags().GetString("block-hash")
	require.NoError(t, err)
	require.Equal(t, "deadbeef", blockHash)

	parentHash, err := cmd.Flags().GetString("parent-hash")
	require.NoError(t, err)
	require.Equal(t, "cafef00d", parentHash)

	out := buf.String()
	require.Contains(t, out, "--wc_height")
	require.Contains(t, out, "--block_hash")
	require.NotContains(t, out, "--parent-hash is deprecated",
		"modern spelling must not trigger a deprecation notice")
}

// TestInstallLegacyFlagAliases_TopLevelRootCmd guards against accidental
// regressions where the walker breaks if rootCmd already has a flagset
// normalizer of its own (the cobra default is one that maps `.` → `-`).
func TestInstallLegacyFlagAliases_PreservesUnrelatedFlags(t *testing.T) {
	resetLegacyFlagWarnerForTest()
	prev := legacyFlagWarnSink
	legacyFlagWarnSink = &bytes.Buffer{}
	defer func() { legacyFlagWarnSink = prev }()

	root := stubTree()
	cmd := findCommand(root, "tx beacon record")
	cmd.Flags().String("unrelated", "", "")

	installLegacyFlagAliases(root)

	require.NotNil(t, cmd.Flags().Lookup("unrelated"),
		"unrelated flags must not be affected by the alias normalizer")
	require.NotNil(t, cmd.Flags().Lookup("submit-time"))
	require.NotNil(t, cmd.Flags().Lookup("subtime"))
}

func lastSegment(path string) string {
	parts := strings.Fields(path)
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

// guard against a future cobra version reordering positional parsing —
// pflag.NormalizedName must remain a string alias.
var _ pflag.NormalizedName = pflag.NormalizedName("")
