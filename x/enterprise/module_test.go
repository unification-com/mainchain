package enterprise_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}
