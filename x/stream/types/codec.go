package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Stream{}, "stream/v1/Stream", nil)
	cdc.RegisterConcrete(&Params{}, "stream/v1/Params", nil)

	legacy.RegisterAminoMsg(cdc, &MsgCreateStream{}, "stream/MsgCreateStream")
	legacy.RegisterAminoMsg(cdc, &MsgClaimStream{}, "stream/MsgClaimStream")
	legacy.RegisterAminoMsg(cdc, &MsgTopUpDeposit{}, "stream/MsgTopUpDeposit")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateFlowRate{}, "stream/MsgUpdateFlowRate")
	legacy.RegisterAminoMsg(cdc, &MsgCancelStream{}, "stream/MsgCancelStream")
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
