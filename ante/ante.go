package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	beaconante "github.com/unification-com/mainchain/x/beacon/ante"
	entante "github.com/unification-com/mainchain/x/enterprise/ante"
	wrkante "github.com/unification-com/mainchain/x/wrkchain/ante"
)

func NewAnteHandler(
	ak authante.AccountKeeper,
	bankKeeper BankKeeper,
	wrkchainKeeper wrkante.WrkchainKeeper,
	beaconKeeper beaconante.BeaconKeeper,
	enterpriseKeeper entante.EnterpriseKeeper,
	sigGasConsumer authante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		authante.NewRejectExtensionOptionsDecorator(),
		authante.NewMempoolFeeDecorator(),
		authante.NewValidateBasicDecorator(),
		authante.TxTimeoutHeightDecorator{},
		authante.NewValidateMemoDecorator(ak),
		authante.NewConsumeGasForTxSizeDecorator(ak),
		authante.NewRejectFeeGranterDecorator(),
		authante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewValidateSigCountDecorator(ak),
		wrkante.NewCorrectWrkChainFeeDecorator(bankKeeper, ak, wrkchainKeeper, enterpriseKeeper), // WRKChain check Tx fees. Specifically check after MemPool, but before consuming fees/gas and undelegating locked FUND
		beaconante.NewCorrectBeaconFeeDecorator(bankKeeper, ak, beaconKeeper, enterpriseKeeper),       // BEACON check Tx fees. Specifically check after MemPool, but before consuming fees/gas and undelegating locked FUND
		entante.NewCheckLockedUndDecorator(enterpriseKeeper),                       // check for and unlock any locked FUND for valid WRKChain/BEACON Txs
		authante.NewDeductFeeDecorator(ak, bankKeeper),
		authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		authante.NewSigVerificationDecorator(ak, signModeHandler),
		authante.NewIncrementSequenceDecorator(ak),
	)
}
