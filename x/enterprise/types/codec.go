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

	cdc.RegisterConcrete(&PurchaseOrderDecision{}, "enterprise/v1/PurchaseOrderDecision", nil)
	cdc.RegisterConcrete(&EnterpriseUndPurchaseOrder{}, "enterprise/v1/EnterpriseUndPurchaseOrder", nil)
	cdc.RegisterConcrete(&PurchaseOrders{}, "enterprise/v1/PurchaseOrders", nil)
	cdc.RegisterConcrete(&LockedUnd{}, "enterprise/v1/LockedUnd", nil)
	cdc.RegisterConcrete(&LockedUnds{}, "enterprise/v1/LockedUnds", nil)
	cdc.RegisterConcrete(&SpentEFUND{}, "enterprise/v1/SpentEFUND", nil)
	cdc.RegisterConcrete(&EnterpriseUserAccount{}, "enterprise/v1/EnterpriseUserAccount", nil)
	cdc.RegisterConcrete(&WhitelistAddresses{}, "enterprise/v1/WhitelistAddresses", nil)
	cdc.RegisterConcrete(&Params{}, "enterprise/v1/Params", nil)

	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/enterprise/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgUndPurchaseOrder{}, "enterprise/MsgUndPurchaseOrder")
	legacy.RegisterAminoMsg(cdc, &MsgProcessUndPurchaseOrder{}, "enterprise/MsgProcessUndPurchaseOrder")
	legacy.RegisterAminoMsg(cdc, &MsgWhitelistAddress{}, "enterprise/MsgWhitelistAddress")
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
