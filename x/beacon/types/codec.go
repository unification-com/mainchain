package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/enterprise interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgRegisterBeacon{}, "beacon/RegisterBeacon", nil)
	cdc.RegisterConcrete(MsgRecordBeaconTimestamp{}, "beacon/RecordBeaconTimestamp", nil)

	cdc.RegisterConcrete(&Beacon{}, "beacon/Beacon", nil)
	cdc.RegisterConcrete(&BeaconTimestamp{}, "beacon/BeaconTimestamp", nil)
	cdc.RegisterConcrete(&BeaconTimestampGenesisExport{}, "beacon/BeaconTimestampGenesisExport", nil)
	cdc.RegisterConcrete(&Params{}, "beacon/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "beacon/GenesisState", nil)
	cdc.RegisterConcrete(&BeaconExport{}, "beacon/BeaconExport", nil)

}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterBeacon{},
		&MsgRecordBeaconTimestamp{},
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
	amino.Seal()
}
