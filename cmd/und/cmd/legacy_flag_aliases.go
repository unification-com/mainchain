package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// legacyFlagAlias maps a legacy flag spelling to the canonical (modern)
// flag name that autocli derives from the proto field. The two spellings
// resolve to the same *pflag.Flag — only the legacy spelling triggers a
// one-shot deprecation notice on stderr.
type legacyFlagAlias struct {
	legacy    string
	canonical string
}

// legacyFlagAliasMap is the registry of legacy → modern flag aliases for
// each Tx subcommand that drifted during the autocli migration. Path is
// the space-separated walk from rootCmd (e.g. "tx beacon record").
func legacyFlagAliasMap() map[string][]legacyFlagAlias {
	return map[string][]legacyFlagAlias{
		"tx beacon record": {
			{legacy: "subtime", canonical: "submit-time"},
		},
		"tx wrkchain register": {
			{legacy: "base", canonical: "base-type"},
			{legacy: "genesis", canonical: "genesis-hash"},
		},
		"tx wrkchain record": {
			{legacy: "wc_height", canonical: "height"},
			{legacy: "block_hash", canonical: "block-hash"},
			{legacy: "parent_hash", canonical: "parent-hash"},
		},
	}
}

// legacyFlagWarnSink is the destination for deprecation notices. Tests
// override it to capture output without polluting test stderr.
var legacyFlagWarnSink io.Writer = os.Stderr

// legacyFlagWarnedOnce tracks which (cmdPath, legacy) pairs have already
// emitted their one-shot deprecation notice during this process.
var legacyFlagWarnedOnce sync.Map // map[string]struct{}

func warnLegacyFlagOnce(cmdPath, legacy, canonical string) {
	warnLegacyOnce(cmdPath, "flag --"+legacy, "--"+canonical)
}

// warnLegacyTokenOnce is the positional-argument counterpart of
// warnLegacyFlagOnce: same one-shot semantics, but the message frames
// the value as a token rather than a flag spelling.
func warnLegacyTokenOnce(cmdPath, legacy, canonical string) {
	warnLegacyOnce(cmdPath, "value '"+legacy+"'", "'"+canonical+"'")
}

func warnLegacyOnce(cmdPath, legacyLabel, canonicalLabel string) {
	key := cmdPath + " " + legacyLabel
	if _, loaded := legacyFlagWarnedOnce.LoadOrStore(key, struct{}{}); loaded {
		return
	}
	fmt.Fprintf(legacyFlagWarnSink,
		"und: %s on `und %s` is deprecated; use %s instead.\n",
		legacyLabel, cmdPath, canonicalLabel)
}

// installLegacyFlagAliases walks rootCmd and installs a per-flagset
// normalize function on every Tx subcommand that drifted during the
// autocli migration. Each normalizer rewrites legacy flag spellings onto
// the canonical proto-derived name, so production scripts that pass
// --subtime / --wc_height / --block_hash / --parent_hash / --base /
// --genesis keep working unchanged. The legacy spelling is purely an
// alias — pflag never registers it as its own flag, so help output and
// tx-body serialization always see the canonical name.
//
// The walker is best-effort: if autocli fails to produce a target cmd,
// installation is skipped silently in production. The companion test
// asserts the walker reaches every expected target.
func installLegacyFlagAliases(rootCmd *cobra.Command) {
	for path, aliases := range legacyFlagAliasMap() {
		cmd := findCommand(rootCmd, path)
		if cmd == nil {
			continue
		}
		path := path
		aliases := append([]legacyFlagAlias(nil), aliases...)
		cmd.Flags().SetNormalizeFunc(func(_ *pflag.FlagSet, name string) pflag.NormalizedName {
			for _, a := range aliases {
				if name == a.legacy {
					warnLegacyFlagOnce(path, a.legacy, a.canonical)
					return pflag.NormalizedName(a.canonical)
				}
			}
			return pflag.NormalizedName(name)
		})
	}
}

// findCommand resolves a space-separated path under rootCmd, descending
// one segment at a time and honouring cobra's per-cmd Aliases.
func findCommand(rootCmd *cobra.Command, path string) *cobra.Command {
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
		if next == nil {
			return nil
		}
		cur = next
	}
	return cur
}
