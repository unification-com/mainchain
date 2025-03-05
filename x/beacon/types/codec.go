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

// RegisterLegacyAminoCodec registers the necessary x/beacon interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	//cdc.RegisterConcrete(MsgRegisterBeacon{}, "beacon/RegisterBeacon", nil)
	//cdc.RegisterConcrete(MsgRecordBeaconTimestamp{}, "beacon/RecordBeaconTimestamp", nil)
	//cdc.RegisterConcrete(MsgPurchaseBeaconStateStorage{}, "beacon/PurchaseBeaconStateStorage", nil)

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
