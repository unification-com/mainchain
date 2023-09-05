package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	ibcante "github.com/cosmos/ibc-go/v5/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	beaconante "github.com/unification-com/mainchain/x/beacon/ante"
	entante "github.com/unification-com/mainchain/x/enterprise/ante"
	wrkante "github.com/unification-com/mainchain/x/wrkchain/ante"
)

type HandlerOptions struct {
	authante.HandlerOptions

	BK               BankKeeper
	BeaconKeeper     beaconante.BeaconKeeper
	EnterpriseKeeper entante.EnterpriseKeeper
	IBCKeeper        *ibckeeper.Keeper
	WrkchainKeeper   wrkante.WrkchainKeeper
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {

	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.BK == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for AnteHandler")
	}
	if options.WrkchainKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "wrkchain keeper is required for AnteHandler")
	}
	if options.BeaconKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "beacon keeper is required for AnteHandler")
	}
	if options.EnterpriseKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "enterprise keeper is required for AnteHandler")
	}
	if options.EnterpriseKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "ibc keeper is required for AnteHandler")
	}

	anteDecorators := []sdk.AnteDecorator{
		authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		authante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		authante.NewValidateBasicDecorator(),
		authante.NewTxTimeoutHeightDecorator(),
		authante.NewValidateMemoDecorator(options.AccountKeeper),
		authante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		wrkante.NewCorrectWrkChainFeeDecorator(options.BK, options.AccountKeeper, options.WrkchainKeeper, options.EnterpriseKeeper), // WRKChain check Tx fees. Specifically check after MemPool, but before consuming fees/gas and undelegating locked FUND
		beaconante.NewCorrectBeaconFeeDecorator(options.BK, options.AccountKeeper, options.BeaconKeeper, options.EnterpriseKeeper),  // BEACON check Tx fees. Specifically check after MemPool, but before consuming fees/gas and undelegating locked FUND
		entante.NewCheckLockedUndDecorator(options.EnterpriseKeeper),                                                                // check for and unlock any locked FUND for valid WRKChain/BEACON Txs
		authante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		authante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewValidateSigCountDecorator(options.AccountKeeper),
		authante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		authante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		authante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
