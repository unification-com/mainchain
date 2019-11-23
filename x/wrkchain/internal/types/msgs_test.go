package types

import (

	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestMsgRegisterWrkChain(t *testing.T) {

	tests := []struct {
		moniker  string
		name     string
		genesis  string
		owner    sdk.AccAddress
		expectPass bool
	}{
		{"wc1", "WC 1", "genhash1", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"", "WC 1", "genhash", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"wc2", "", "genhash2", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"wc3", "WC 3", "", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"wc4", "WC 4", "genhash4", sdk.AccAddress{}, false},
	}

	for i, tc := range tests {
		msg := NewMsgRegisterWrkChain(
			tc.moniker,
			tc.genesis,
			tc.name,
			tc.owner,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgRecordWrkChainBlock(t *testing.T) {

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	emptyAddr := sdk.AccAddress{}

	tests := []struct {
		WrkChainID uint64
		Height     uint64
		BlockHash  string
		ParentHash string
		Hash1      string
		Hash2      string
		Hash3      string
		Owner      sdk.AccAddress
		expectPass bool
	}{
		{1, 1, "blockhash", "parenthash", "hash1", "hash2", "hash3", addr, true},
		{1, 1, "blockhash", "parenthash", "hash1", "hash2", "", addr, true},
		{1, 1, "blockhash", "parenthash", "hash1", "", "", addr, true},
		{1, 1, "blockhash", "parenthash", "", "", "", addr, true},
		{1, 1, "blockhash", "", "", "", "", addr, true},
		{1, 1, "blockhash", "parenthash", "", "", "", emptyAddr, false},
		{1, 1, "", "", "", "", "", addr, false},
		{1, 0, "blockhash", "", "", "", "", addr, false},
		{0, 1, "blockhash", "", "", "", "", addr, false},
	}

	for i, tc := range tests {
		msg := NewMsgRecordWrkChainBlock(
			tc.WrkChainID,
			tc.Height,
			tc.BlockHash,
			tc.ParentHash,
			tc.Hash1,
			tc.Hash2,
			tc.Hash3,
			tc.Owner,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
