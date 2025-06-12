package rapidgen

import (
	"fmt"
	cosmos_proto "github.com/cosmos/cosmos-proto"
	"github.com/cosmos/cosmos-proto/rapidproto"
	gogoproto "github.com/cosmos/gogoproto/proto"
	enterpriseapi "github.com/unification-com/mainchain/api/mainchain/enterprise/v1"
	streamapi "github.com/unification-com/mainchain/api/mainchain/stream/v1"
	wrkchainapi "github.com/unification-com/mainchain/api/mainchain/wrkchain/v1"
	enterprisetypes "github.com/unification-com/mainchain/x/enterprise/types"
	streamtypes "github.com/unification-com/mainchain/x/stream/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"pgregory.net/rapid"

	beaconapi "github.com/unification-com/mainchain/api/mainchain/beacon/v1"
	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
)

type GeneratedType struct {
	Pulsar proto.Message
	Gogo   gogoproto.Message
	Opts   rapidproto.GeneratorOptions
}

func GenType(gogo gogoproto.Message, pulsar proto.Message, opts rapidproto.GeneratorOptions) GeneratedType {
	return GeneratedType{
		Pulsar: pulsar,
		Gogo:   gogo,
		Opts:   opts,
	}
}

func GeneratorFieldMapper(t *rapid.T, field protoreflect.FieldDescriptor, name string) (protoreflect.Value, bool) {
	opts := field.Options()
	if proto.HasExtension(opts, cosmos_proto.E_Scalar) {
		scalar := proto.GetExtension(opts, cosmos_proto.E_Scalar).(string)
		switch scalar {
		case "cosmos.Int":
			i32 := rapid.Int32().Draw(t, name)
			return protoreflect.ValueOfString(fmt.Sprintf("%d", i32)), true
		case "cosmos.Dec":
			if field.Kind() == protoreflect.BytesKind {
				return protoreflect.ValueOfBytes([]byte{}), true
			}

			return protoreflect.ValueOfString(""), true
		}
	}

	return protoreflect.Value{}, false
}

var (
	GenOpts = rapidproto.GeneratorOptions{
		Resolver:  protoregistry.GlobalTypes,
		FieldMaps: []rapidproto.FieldMapper{GeneratorFieldMapper},
	}
	SignableTypes = []GeneratedType{
		// beacon
		GenType(&beacontypes.MsgRegisterBeacon{}, &beaconapi.MsgRegisterBeacon{}, GenOpts),
		GenType(&beacontypes.MsgPurchaseBeaconStateStorage{}, &beaconapi.MsgPurchaseBeaconStateStorage{}, GenOpts),
		GenType(&beacontypes.MsgRecordBeaconTimestamp{}, &beaconapi.MsgRecordBeaconTimestamp{}, GenOpts),
		GenType(&beacontypes.MsgUpdateParams{}, &beaconapi.MsgUpdateParams{}, GenOpts.WithDisallowNil()),

		//// enterprise
		GenType(&enterprisetypes.MsgUndPurchaseOrder{}, &enterpriseapi.MsgUndPurchaseOrder{}, GenOpts.WithDisallowNil()),
		GenType(&enterprisetypes.MsgProcessUndPurchaseOrder{}, &enterpriseapi.MsgProcessUndPurchaseOrder{}, GenOpts),
		GenType(&enterprisetypes.MsgWhitelistAddress{}, &enterpriseapi.MsgWhitelistAddress{}, GenOpts),
		GenType(&enterprisetypes.MsgUpdateParams{}, &enterpriseapi.MsgUpdateParams{}, GenOpts.WithDisallowNil()),

		// stream
		GenType(&streamtypes.MsgCreateStream{}, &streamapi.MsgCreateStream{}, GenOpts.WithDisallowNil()),
		GenType(&streamtypes.MsgClaimStream{}, &streamapi.MsgClaimStream{}, GenOpts),
		GenType(&streamtypes.MsgTopUpDeposit{}, &streamapi.MsgTopUpDeposit{}, GenOpts.WithDisallowNil()),
		GenType(&streamtypes.MsgUpdateFlowRate{}, &streamapi.MsgUpdateFlowRate{}, GenOpts),
		GenType(&streamtypes.MsgCancelStream{}, &streamapi.MsgCancelStream{}, GenOpts),
		GenType(&streamtypes.MsgUpdateParams{}, &streamapi.MsgUpdateParams{}, GenOpts.WithDisallowNil()),

		// wrkchain
		GenType(&wrkchaintypes.MsgRegisterWrkChain{}, &wrkchainapi.MsgRegisterWrkChain{}, GenOpts),
		GenType(&wrkchaintypes.MsgPurchaseWrkChainStateStorage{}, &wrkchainapi.MsgPurchaseWrkChainStateStorage{}, GenOpts),
		GenType(&wrkchaintypes.MsgRecordWrkChainBlock{}, &wrkchainapi.MsgRecordWrkChainBlock{}, GenOpts),
		GenType(&wrkchaintypes.MsgUpdateParams{}, &wrkchainapi.MsgUpdateParams{}, GenOpts.WithDisallowNil()),
	}
	NonsignableTypes = []GeneratedType{
		GenType(&beacontypes.Beacon{}, &beaconapi.Beacon{}, GenOpts),
		GenType(&beacontypes.BeaconStorageLimit{}, &beaconapi.BeaconStorageLimit{}, GenOpts),
		GenType(&beacontypes.BeaconTimestamp{}, &beaconapi.BeaconTimestamp{}, GenOpts),
		GenType(&beacontypes.Params{}, &beaconapi.Params{}, GenOpts),

		GenType(&enterprisetypes.PurchaseOrderDecision{}, &enterpriseapi.PurchaseOrderDecision{}, GenOpts),
		GenType(&enterprisetypes.EnterpriseUndPurchaseOrder{}, &enterpriseapi.EnterpriseUndPurchaseOrder{}, GenOpts.WithDisallowNil()),
		GenType(&enterprisetypes.PurchaseOrders{}, &enterpriseapi.PurchaseOrders{}, GenOpts.WithDisallowNil()),
		GenType(&enterprisetypes.LockedUnd{}, &enterpriseapi.LockedUnd{}, GenOpts.WithDisallowNil()),
		GenType(&enterprisetypes.SpentEFUND{}, &enterpriseapi.SpentEFUND{}, GenOpts.WithDisallowNil()),
		GenType(&enterprisetypes.EnterpriseUserAccount{}, &enterpriseapi.EnterpriseUserAccount{}, GenOpts.WithDisallowNil()),
		GenType(&enterprisetypes.WhitelistAddresses{}, &enterpriseapi.WhitelistAddresses{}, GenOpts),
		GenType(&enterprisetypes.Params{}, &enterpriseapi.Params{}, GenOpts),

		GenType(&streamtypes.Stream{}, &streamapi.Stream{}, GenOpts.WithDisallowNil()),
		GenType(&streamtypes.Params{}, &streamapi.Params{}, GenOpts),

		GenType(&wrkchaintypes.WrkChain{}, &wrkchainapi.WrkChain{}, GenOpts),
		GenType(&wrkchaintypes.WrkChainStorageLimit{}, &wrkchainapi.WrkChainStorageLimit{}, GenOpts),
		GenType(&wrkchaintypes.WrkChainBlock{}, &wrkchainapi.WrkChainBlock{}, GenOpts),
		GenType(&wrkchaintypes.Params{}, &wrkchainapi.Params{}, GenOpts),
	}
	DefaultGeneratedTypes = append(SignableTypes, NonsignableTypes...)
)
