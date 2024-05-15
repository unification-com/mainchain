package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/beacon interfaces and concrete types
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
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "mainchain/x/beacon/MsgUpdateParams")

}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterBeacon{},
		&MsgRecordBeaconTimestamp{},
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

	// Register all Amino interfaces and concrete types on the authz Amino codec
	// so that this can later be used to properly serialize MsgGrant and MsgExec
	// instances.
	RegisterLegacyAminoCodec(authzcodec.Amino)
}
