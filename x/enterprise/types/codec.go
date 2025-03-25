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

	cdc.RegisterConcrete(&PurchaseOrderDecision{}, "enterprise/PurchaseOrderDecision", nil)
	cdc.RegisterConcrete(&EnterpriseUndPurchaseOrder{}, "enterprise/EnterpriseUndPurchaseOrder", nil)
	cdc.RegisterConcrete(&PurchaseOrders{}, "enterprise/PurchaseOrders", nil)
	cdc.RegisterConcrete(&LockedUnd{}, "enterprise/LockedUnd", nil)
	cdc.RegisterConcrete(&LockedUnds{}, "enterprise/LockedUnds", nil)
	cdc.RegisterConcrete(&SpentEFUND{}, "enterprise/SpentEFUND", nil)
	cdc.RegisterConcrete(&EnterpriseUserAccount{}, "enterprise/EnterpriseUserAccount", nil)
	cdc.RegisterConcrete(&UndSupply{}, "enterprise/UndSupply", nil)
	cdc.RegisterConcrete(&WhitelistAddresses{}, "enterprise/WhitelistAddresses", nil)
	cdc.RegisterConcrete(&Params{}, "enterprise/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "enterprise/GenesisState", nil)

	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/enterprise/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgUndPurchaseOrder{}, "enterprise/PurchaseUnd")
	legacy.RegisterAminoMsg(cdc, &MsgProcessUndPurchaseOrder{}, "enterprise/ProcessUndPurchaseOrder")
	legacy.RegisterAminoMsg(cdc, &MsgWhitelistAddress{}, "enterprise/WhitelistAddress")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUndPurchaseOrder{},
		&MsgProcessUndPurchaseOrder{},
		&MsgWhitelistAddress{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
