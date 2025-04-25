package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/enterprise interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {

	cdc.RegisterConcrete(&WrkChain{}, "wrkchain/v1/WrkChain", nil)
	cdc.RegisterConcrete(&WrkChainBlock{}, "wrkchain/v1/WrkChainBlock", nil)
	cdc.RegisterConcrete(&Params{}, "wrkchain/v1/Params", nil)
	cdc.RegisterConcrete(&WrkChainStorageLimit{}, "wrkchain/v1/WrkChainStorageLimit", nil)
	legacy.RegisterAminoMsg(cdc, &MsgRegisterWrkChain{}, "wrkchain/MsgRegisterWrkChain")
	legacy.RegisterAminoMsg(cdc, &MsgRecordWrkChainBlock{}, "wrkchain/MsgRecordWrkChainBlock")
	legacy.RegisterAminoMsg(cdc, &MsgPurchaseWrkChainStateStorage{}, "wrkchain/MsgPurchaseWrkChainStorage")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/wrkchain/MsgUpdateParams")

}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterWrkChain{},
		&MsgRecordWrkChainBlock{},
		&MsgUpdateParams{},
		&MsgPurchaseWrkChainStateStorage{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
