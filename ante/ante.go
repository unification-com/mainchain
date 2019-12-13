package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/unification-com/mainchain/x/beacon"
	"github.com/unification-com/mainchain/x/enterprise"
	"github.com/unification-com/mainchain/x/wrkchain"
)

func NewAnteHandler(ak auth.AccountKeeper, supplyKeeper supply.Keeper, wrkchainKeeper wrkchain.Keeper, beaconKeeper beacon.Keeper, enterpriseKeeper enterprise.Keeper, sigGasConsumer auth.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		authante.NewMempoolFeeDecorator(),
		authante.NewValidateBasicDecorator(),
		authante.NewValidateMemoDecorator(ak),
		authante.NewConsumeGasForTxSizeDecorator(ak),
		authante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewValidateSigCountDecorator(ak),
		wrkchain.NewCorrectWrkChainFeeDecorator(ak, wrkchainKeeper, enterpriseKeeper), // WRKChain check Tx fees. Specifically check after MemPool, but before consuming fees/gas and undelegating locked UND
		beacon.NewCorrectBeaconFeeDecorator(ak, beaconKeeper, enterpriseKeeper),       // BEACON check Tx fees. Specifically check after MemPool, but before consuming fees/gas and undelegating locked UND
		enterprise.NewCheckLockedUndDecorator(enterpriseKeeper),                       // check for and unlock any locked UND for valid WRKChain/BEACON Txs
		authante.NewDeductFeeDecorator(ak, supplyKeeper),
		authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		authante.NewSigVerificationDecorator(ak),
		authante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}
