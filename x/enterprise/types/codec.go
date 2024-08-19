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
	cdc.RegisterConcrete(&MsgUndPurchaseOrder{}, "enterprise/PurchaseUnd", nil)
	cdc.RegisterConcrete(&MsgProcessUndPurchaseOrder{}, "enterprise/ProcessUndPurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgWhitelistAddress{}, "enterprise/WhitelistAddress", nil)

	cdc.RegisterConcrete(&PurchaseOrderDecision{}, "enterprise/PurchaseOrderDecision", nil)
	cdc.RegisterConcrete(&EnterpriseUndPurchaseOrder{}, "enterprise/EnterpriseUndPurchaseOrder", nil)
	cdc.RegisterConcrete(&PurchaseOrders{}, "enterprise/PurchaseOrders", nil)
	cdc.RegisterConcrete(&LockedUnd{}, "enterprise/LockedUnd", nil)
	cdc.RegisterConcrete(&LockedUnds{}, "enterprise/LockedUnds", nil)
	cdc.RegisterConcrete(&UndSupply{}, "enterprise/UndSupply", nil)
	cdc.RegisterConcrete(&WhitelistAddresses{}, "enterprise/WhitelistAddresses", nil)
	cdc.RegisterConcrete(&Params{}, "enterprise/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "enterprise/GenesisState", nil)
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/enterprise/MsgUpdateParams")
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
