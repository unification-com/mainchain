package enterprise_test

import (
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"testing"

	"github.com/stretchr/testify/require"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}
