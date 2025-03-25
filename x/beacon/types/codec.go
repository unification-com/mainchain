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

	cdc.RegisterConcrete(&Beacon{}, "beacon/Beacon", nil)
	cdc.RegisterConcrete(&BeaconTimestamp{}, "beacon/BeaconTimestamp", nil)
	cdc.RegisterConcrete(&BeaconTimestampGenesisExport{}, "beacon/BeaconTimestampGenesisExport", nil)
	cdc.RegisterConcrete(&Params{}, "beacon/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "beacon/GenesisState", nil)
	cdc.RegisterConcrete(&BeaconExport{}, "beacon/BeaconExport", nil)
	cdc.RegisterConcrete(&BeaconStorageLimit{}, "beacon/BeaconStorageLimit", nil)
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/beacon/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgRegisterBeacon{}, "beacon/RegisterBeacon")
	legacy.RegisterAminoMsg(cdc, &MsgRecordBeaconTimestamp{}, "beacon/RecordBeaconTimestamp")
	legacy.RegisterAminoMsg(cdc, &MsgPurchaseBeaconStateStorage{}, "beacon/PurchaseBeaconStateStorage")

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
