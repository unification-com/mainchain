package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestMsgRegisterWrkChain_Route(t *testing.T) {
	msg := MsgRegisterWrkChain{}
	require.Equal(t, ModuleName, msg.Route())
}

func TestMsgRegisterWrkChain_Type(t *testing.T) {
	msg := MsgRegisterWrkChain{}
	require.Equal(t, "register_wrkchain", msg.Type())
}

func TestMsgRegisterWrkChain_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	ownerAddr := sdk.AccAddress(pubKey2.Address())
	msg := MsgRegisterWrkChain{Owner: ownerAddr}
	require.True(t, msg.GetSigners()[0].Equals(ownerAddr))
}

func TestMsgRecordWrkChainBlock_Route(t *testing.T) {
	msg := MsgRecordWrkChainBlock{}
	require.Equal(t, ModuleName, msg.Route())
}

func TestMsgRecordWrkChainBlock_Type(t *testing.T) {
	msg := MsgRecordWrkChainBlock{}
	require.Equal(t, "record_wrkchain_hash", msg.Type())
}

func TestMsgRecordWrkChainBlock_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	ownerAddr := sdk.AccAddress(pubKey2.Address())
	msg := MsgRecordWrkChainBlock{Owner: ownerAddr}
	require.True(t, msg.GetSigners()[0].Equals(ownerAddr))
}

func TestMsgRegisterWrkChain(t *testing.T) {

	tests := []struct {
		moniker    string
		name       string
		genesis    string
		baseType   string
		owner      sdk.AccAddress
		expectPass bool
	}{
		{"wc1", "WC 1", "genhash1", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"", "WC 1", "genhash", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"", "", "genhash2", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"", "WC 3", "", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"", "WC 4", "genhash4", "", sdk.AccAddress{}, false},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "WC 1", "genhash", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a51", "WC 1", "genhash", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a51", "genhash", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "genhash", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad72", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad721", "geth", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
	}

	for i, tc := range tests {
		msg := NewMsgRegisterWrkChain(
			tc.moniker,
			tc.genesis,
			tc.name,
			tc.baseType,
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
		{1,
			2,
			"0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad72",
			"0xd806cf49e61dba73e1d92d5a04fabb94b9318f9007a1c9ab74b6678b8d7bd405",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0xe89cf713dd1eb2928d6b007ffc43c10780fdb53bfd29ad59e76a0b0b97815101", addr, true},
		{1,
			3,
			"0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad721",
			"0xd806cf49e61dba73e1d92d5a04fabb94b9318f9007a1c9ab74b6678b8d7bd405",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0xe89cf713dd1eb2928d6b007ffc43c10780fdb53bfd29ad59e76a0b0b97815101", addr, false},
		{1,
			3,
			"0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad72",
			"0xd806cf49e61dba73e1d92d5a04fabb94b9318f9007a1c9ab74b6678b8d7bd4051",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0xe89cf713dd1eb2928d6b007ffc43c10780fdb53bfd29ad59e76a0b0b97815101", addr, false},
		{1,
			3,
			"0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad72",
			"0xd806cf49e61dba73e1d92d5a04fabb94b9318f9007a1c9ab74b6678b8d7bd405",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b4211",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0xe89cf713dd1eb2928d6b007ffc43c10780fdb53bfd29ad59e76a0b0b97815101", addr, false},
		{1,
			3,
			"0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad72",
			"0xd806cf49e61dba73e1d92d5a04fabb94b9318f9007a1c9ab74b6678b8d7bd405",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b4211",
			"0xe89cf713dd1eb2928d6b007ffc43c10780fdb53bfd29ad59e76a0b0b97815101", addr, false},
		{1,
			3,
			"0x191ee45e03daae0d6a90afccab1a49b37b4297dec658e9d7596fb7545d62ad72",
			"0xd806cf49e61dba73e1d92d5a04fabb94b9318f9007a1c9ab74b6678b8d7bd405",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"0xe89cf713dd1eb2928d6b007ffc43c10780fdb53bfd29ad59e76a0b0b978151011", addr, false},
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
