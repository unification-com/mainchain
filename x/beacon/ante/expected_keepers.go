package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/ante/feecheck"
)

// The account, bank and enterprise keeper contracts are identical for the BEACON
// and WRKChain fee decorators, so they live in ante/feecheck alongside the shared
// solvency check. Aliased rather than re-declared so existing references — app
// wiring, the app-level ante HandlerOptions, tests — keep compiling unchanged.
type (
	AccountKeeper    = feecheck.AccountKeeper
	BankKeeper       = feecheck.BankKeeper
	EnterpriseKeeper = feecheck.EnterpriseKeeper
)

type BeaconKeeper interface {
	GetZeroFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetParamDenom(ctx sdk.Context) string
	GetRegistrationFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetRecordFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetPurchaseStorageFeeAsCoin(ctx sdk.Context) sdk.Coin
	GetMaxPurchasableSlots(ctx sdk.Context, beaconId uint64) uint64
}
