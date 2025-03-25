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

	cdc.RegisterConcrete(&WrkChain{}, "wrkchain/WrkChain", nil)
	cdc.RegisterConcrete(&WrkChainBlock{}, "wrkchain/WrkChainBlock", nil)
	cdc.RegisterConcrete(&WrkChainBlockGenesisExport{}, "wrkchain/WrkChainBlockGenesisExport", nil)
	cdc.RegisterConcrete(&Params{}, "wrkchain/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "wrkchain/GenesisState", nil)
	cdc.RegisterConcrete(&WrkChainExport{}, "wrkchain/WrkChainExport", nil)
	cdc.RegisterConcrete(&WrkChainStorageLimit{}, "wrkchain/WrkChainStorageLimit", nil)
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/wrkchain/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgRegisterWrkChain{}, "wrkchain/RegisterWrkChain")
	legacy.RegisterAminoMsg(cdc, &MsgRecordWrkChainBlock{}, "wrkchain/MsgRecordWrkChainBlock")
	legacy.RegisterAminoMsg(cdc, &MsgPurchaseWrkChainStateStorage{}, "wrkchain/PurchaseWrkChainStateStorage")

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
