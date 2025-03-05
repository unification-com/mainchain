package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/enterprise interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	//cdc.RegisterConcrete(MsgRegisterWrkChain{}, "wrkchain/RegisterWrkChain", nil)
	//cdc.RegisterConcrete(MsgRecordWrkChainBlock{}, "wrkchain/MsgRecordWrkChainBlock", nil)
	//cdc.RegisterConcrete(MsgPurchaseWrkChainStateStorage{}, "wrkchain/PurchaseWrkChainStateStorage", nil)

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
