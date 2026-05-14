package cmd

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// Companion to legacy_flag_aliases.go: post-autocli behaviour shims that
// preserve legacy CLI ergonomics autocli can't express through descriptor
// metadata alone. Like the flag-alias walker, every shim chains onto the
// autocli-generated cmd via PreRunE so the mutation is visible to the
// downstream RunE that builds and broadcasts the tx.

// nowUnixSec is the clock used by the submit_time default-to-now shim.
// Tests override this to assert deterministic output.
var nowUnixSec = func() uint64 { return uint64(time.Now().Unix()) }

// legacyDecisionAliases maps the legacy `enterprise process` decision
// tokens to the autocli-stripped proto-enum forms. Both spellings remain
// accepted; the legacy form prints a one-shot deprecation notice via the
// existing warner.
var legacyDecisionAliases = map[string]string{
	"accept": "accepted",
	"reject": "rejected",
}

// installAutocliCompatShims wires the behaviour shims listed below onto
// the autocli-generated cmd tree. Must run after EnhanceRootCommand and
// after installLegacyFlagAliases, since both shims operate on flag /
// positional state that the alias walker has already normalised.
func installAutocliCompatShims(rootCmd *cobra.Command) {
	if cmd := findCommand(rootCmd, "tx beacon record"); cmd != nil {
		chainPreRunE(cmd, beaconRecordFillSubmitTime)
	}
	if cmd := findCommand(rootCmd, "tx enterprise process"); cmd != nil {
		chainPreRunE(cmd, enterpriseProcessAliasDecisionToken)
	}
}

// chainPreRunE prepends shim onto cmd.PreRunE, preserving any existing
// hook autocli or a parent cmd may already have installed.
func chainPreRunE(cmd *cobra.Command, shim func(*cobra.Command, []string) error) {
	existing := cmd.PreRunE
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		if err := shim(c, args); err != nil {
			return err
		}
		if existing != nil {
			return existing(c, args)
		}
		return nil
	}
}

// beaconRecordFillSubmitTime restores the pre-autocli convenience of
// defaulting submit_time to the current unix epoch when the operator
// omits --submit-time entirely. MsgRecordBeaconTimestamp.ValidateBasic()
// rejects submit_time=0, so without this shim the autocli surface would
// silently break production scripts that relied on the default. If the
// operator passes --submit-time (or its legacy --subtime alias) the
// supplied value flows through unchanged — including an explicit 0,
// which then surfaces the server-side rejection unchanged.
func beaconRecordFillSubmitTime(cmd *cobra.Command, _ []string) error {
	if cmd.Flags().Changed("submit-time") {
		return nil
	}
	return cmd.Flags().Set("submit-time", strconv.FormatUint(nowUnixSec(), 10))
}

// enterpriseProcessAliasDecisionToken rewrites the legacy decision
// tokens (`accept`, `reject`) onto the autocli-stripped proto-enum
// forms (`accepted`, `rejected`) before autocli's RunE binds the
// positional to msg.Decision. The proto enum names use the STATUS_
// prefix and autocli strips it mechanically — that's why the legacy
// shorter tokens stopped parsing once the legacy cobra cmd went away.
func enterpriseProcessAliasDecisionToken(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return nil
	}
	canonical, ok := legacyDecisionAliases[args[1]]
	if !ok {
		return nil
	}
	warnLegacyTokenOnce("tx enterprise process", args[1], canonical)
	args[1] = canonical
	return nil
}

