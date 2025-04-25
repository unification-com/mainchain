package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/beacon interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {

	cdc.RegisterConcrete(&Beacon{}, "beacon/v1/Beacon", nil)
	cdc.RegisterConcrete(&BeaconTimestamp{}, "beacon/v1/BeaconTimestamp", nil)
	cdc.RegisterConcrete(&Params{}, "beacon/v1/Params", nil)
	cdc.RegisterConcrete(&BeaconStorageLimit{}, "beacon/v1/BeaconStorageLimit", nil)

	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/beacon/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgRegisterBeacon{}, "beacon/MsgRegisterBeacon")
	legacy.RegisterAminoMsg(cdc, &MsgRecordBeaconTimestamp{}, "beacon/MsgRecordBeaconTimestamp")
	legacy.RegisterAminoMsg(cdc, &MsgPurchaseBeaconStateStorage{}, "beacon/MsgPurchaseBeaconStateStorage")

}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterBeacon{},
		&MsgRecordBeaconTimestamp{},
		&MsgUpdateParams{},
		&MsgPurchaseBeaconStateStorage{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
