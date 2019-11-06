package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/unification-com/mainchain-cosmos/x/enterprise"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain"
)

func NewAnteHandler(ak auth.AccountKeeper, supplyKeeper supply.Keeper, wrkchainKeeper wrkchain.Keeper, enterpriseKeeper enterprise.Keeper, sigGasConsumer auth.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		authante.NewMempoolFeeDecorator(),
		authante.NewValidateBasicDecorator(),
		authante.NewValidateMemoDecorator(ak),
		enterprise.NewCheckLockedUndDecorator(ak, enterpriseKeeper), // check if account has enough unlocked UND
		wrkchain.NewWrkChainFeeDecorator(ak, wrkchainKeeper), // WRKChain check Tx fees. Specifically check after MemPool, but before consuming fees/gas
		authante.NewConsumeGasForTxSizeDecorator(ak),
		authante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewValidateSigCountDecorator(ak),
		authante.NewDeductFeeDecorator(ak, supplyKeeper),
		authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		authante.NewSigVerificationDecorator(ak),
		authante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}
