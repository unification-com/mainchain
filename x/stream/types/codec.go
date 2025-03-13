package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
