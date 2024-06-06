package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Stream{}, "stream/Stream", nil)
	cdc.RegisterConcrete(&StreamExport{}, "stream/StreamExport", nil)
	cdc.RegisterConcrete(&Params{}, "stream/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "stream/GenesisState", nil)

	legacy.RegisterAminoMsg(cdc, &MsgCreateStream{}, "mainchain/x/stream/MsgCreateStream")
	legacy.RegisterAminoMsg(cdc, &MsgClaimStream{}, "mainchain/x/stream/MsgClaimStream")
	legacy.RegisterAminoMsg(cdc, &MsgTopUpDeposit{}, "mainchain/x/stream/MsgTopUpDeposit")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateFlowRate{}, "mainchain/x/stream/MsgUpdateFlowRate")
	legacy.RegisterAminoMsg(cdc, &MsgCancelStream{}, "mainchain/x/stream/MsgCancelStream")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/stream/MsgUpdateParams")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateStream{},
		&MsgClaimStream{},
		&MsgTopUpDeposit{},
		&MsgUpdateFlowRate{},
		&MsgCancelStream{},
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)

	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec so that this can later be
	// used to properly serialize MsgSubmitProposal instances
	RegisterLegacyAminoCodec(govcodec.Amino)
}
