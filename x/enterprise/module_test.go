package enterprise_test

import (
	"github.com/unification-com/mainchain/app/test_helpers"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	//app.InitChain(
	//	abcitypes.RequestInitChain{
	//		AppStateBytes: []byte("{}"),
	//		ChainId:       "test-chain-id",
	//	},
	//)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}
