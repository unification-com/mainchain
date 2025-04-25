package tests

import (
	"context"
	mathmod "cosmossdk.io/math"
	"fmt"
	enterpriseapi "github.com/unification-com/mainchain/api/mainchain/enterprise/v1"
	streamapi "github.com/unification-com/mainchain/api/mainchain/stream/v1"
	wrkchainapi "github.com/unification-com/mainchain/api/mainchain/wrkchain/v1"
	enterprisetypes "github.com/unification-com/mainchain/x/enterprise/types"
	streamtypes "github.com/unification-com/mainchain/x/stream/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"testing"
	"time"

	"github.com/cosmos/cosmos-proto/rapidproto"
	gogoproto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"pgregory.net/rapid"

	authapi "cosmossdk.io/api/cosmos/auth/v1beta1"
	v1beta1 "cosmossdk.io/api/cosmos/base/v1beta1"
	msgv1 "cosmossdk.io/api/cosmos/msg/v1"
	txv1beta1 "cosmossdk.io/api/cosmos/tx/v1beta1"
	"cosmossdk.io/x/tx/signing/aminojson"
	signing_testutil "cosmossdk.io/x/tx/signing/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	beaconapi "github.com/unification-com/mainchain/api/mainchain/beacon/v1"
	fundhelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/tests/rapidgen"
	"github.com/unification-com/mainchain/x/beacon"
	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	"github.com/unification-com/mainchain/x/enterprise"
	"github.com/unification-com/mainchain/x/stream"
	"github.com/unification-com/mainchain/x/wrkchain"
)

// TestAminoJSON_Equivalence tests that x/tx/Encoder encoding is equivalent to the legacy Encoder encoding.
// A custom generator is used to generate random messages that are then encoded using both encoders.  The custom
// generator only supports proto.Message (which implement the protoreflect API) so in order to test legacy gogo types
// we end up with a workflow as follows:
//
// 1. Generate a random protobuf proto.Message using the custom generator
// 2. Marshal the proto.Message to protobuf binary bytes
// 3. Unmarshal the protobuf bytes to a gogoproto.Message
// 4. Marshal the gogoproto.Message to amino JSON bytes
// 5. Marshal the proto.Message to amino JSON bytes
// 6. Compare the amino JSON bytes from steps 4 and 5
//
// In order for step 3 to work certain restrictions on the data generated in step 1 must be enforced and are described
// by the mutation of genOpts passed to the generator.
func TestAminoJSON_Equivalence(t *testing.T) {
	encCfg := testutil.MakeTestEncodingConfig(
		beacon.AppModuleBasic{}, enterprise.AppModuleBasic{}, stream.AppModuleBasic{}, wrkchain.AppModuleBasic{})
	legacytx.RegressionTestingAminoCodec = encCfg.Amino
	aj := aminojson.NewEncoder(aminojson.EncoderOptions{DoNotSortFields: true})

	for _, tt := range rapidgen.DefaultGeneratedTypes {
		desc := tt.Pulsar.ProtoReflect().Descriptor()
		name := string(desc.FullName())
		t.Run(name, func(t *testing.T) {
			gen := rapidproto.MessageGenerator(tt.Pulsar, tt.Opts)
			fmt.Printf("testing %s\n", tt.Pulsar.ProtoReflect().Descriptor().FullName())
			rapid.Check(t, func(t *rapid.T) {
				// uncomment to debug; catch a panic and inspect application state
				// defer func() {
				//	if r := recover(); r != nil {
				//		//fmt.Printf("Panic: %+v\n", r)
				//		t.FailNow()
				//	}
				// }()

				msg := gen.Draw(t, "msg")
				postFixPulsarMessage(msg)

				gogo := tt.Gogo
				sanity := tt.Pulsar

				protoBz, err := proto.Marshal(msg)
				require.NoError(t, err)

				err = proto.Unmarshal(protoBz, sanity)
				require.NoError(t, err)

				err = encCfg.Codec.Unmarshal(protoBz, gogo)
				require.NoError(t, err)

				legacyAminoJSON, err := encCfg.Amino.MarshalJSON(gogo)
				require.NoError(t, err)
				aminoJSON, err := aj.Marshal(msg)
				require.NoError(t, err)
				require.Equal(t, string(legacyAminoJSON), string(aminoJSON))

				// test amino json signer handler equivalence
				if !proto.HasExtension(desc.Options(), msgv1.E_Signer) {
					// not signable
					return
				}

				handlerOptions := signing_testutil.HandlerArgumentOptions{
					ChainID:       "test-chain",
					Memo:          "sometestmemo",
					Msg:           tt.Pulsar,
					AccNum:        1,
					AccSeq:        2,
					SignerAddress: "signerAddress",
					Fee: &txv1beta1.Fee{
						Amount: []*v1beta1.Coin{{Denom: "uatom", Amount: "1000"}},
					},
				}

				signerData, txData, err := signing_testutil.MakeHandlerArguments(handlerOptions)
				require.NoError(t, err)

				handler := aminojson.NewSignModeHandler(aminojson.SignModeHandlerOptions{})
				signBz, err := handler.GetSignBytes(context.Background(), signerData, txData)
				require.NoError(t, err)

				legacyHandler := tx.NewSignModeLegacyAminoJSONHandler()
				txBuilder := encCfg.TxConfig.NewTxBuilder()
				require.NoError(t, txBuilder.SetMsgs([]sdk.Msg{tt.Gogo}...))
				txBuilder.SetMemo(handlerOptions.Memo)
				txBuilder.SetFeeAmount(sdk.Coins{sdk.NewInt64Coin("uatom", 1000)})
				theTx := txBuilder.GetTx()

				legacySigningData := signing.SignerData{
					ChainID:       handlerOptions.ChainID,
					Address:       handlerOptions.SignerAddress,
					AccountNumber: handlerOptions.AccNum,
					Sequence:      handlerOptions.AccSeq,
				}
				legacySignBz, err := legacyHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
					legacySigningData, theTx)
				require.NoError(t, err)
				require.Equal(t, string(legacySignBz), string(signBz))
			})
		})
	}
}

func newAny(t *testing.T, msg proto.Message) *anypb.Any {
	bz, err := proto.Marshal(msg)
	require.NoError(t, err)
	typeName := fmt.Sprintf("/%s", msg.ProtoReflect().Descriptor().FullName())
	return &anypb.Any{
		TypeUrl: typeName,
		Value:   bz,
	}
}

// TestAminoJSON_LegacyParity tests that the Encoder encoder produces the same output as the Encoder encoder.
func TestAminoJSON_LegacyParity(t *testing.T) {
	encCfg := testutil.MakeTestEncodingConfig(beacon.AppModuleBasic{}, enterprise.AppModuleBasic{},
		stream.AppModuleBasic{}, wrkchain.AppModuleBasic{})
	legacytx.RegressionTestingAminoCodec = encCfg.Amino

	aj := aminojson.NewEncoder(aminojson.EncoderOptions{DoNotSortFields: true})
	addr1 := sdk.AccAddress("addr1")
	now := time.Now()
	nowUint64 := uint64(now.Unix())

	randomHash := fundhelpers.GenerateRandomString(32)
	testSdkCoin := sdk.NewInt64Coin("nund", 1000)
	testV1beta1Coin := v1beta1.Coin{
		Denom:  "nund",
		Amount: "1000",
	}
	//
	//genericAuth, _ := codectypes.NewAnyWithValue(&authztypes.GenericAuthorization{Msg: "foo"})
	//genericAuthPulsar := newAny(t, &authzapi.GenericAuthorization{Msg: "foo"})
	//pubkeyAny, _ := codectypes.NewAnyWithValue(&secp256k1types.PubKey{Key: []byte("foo")})
	//pubkeyAnyPulsar := newAny(t, &secp256k1.PubKey{Key: []byte("foo")})
	//dec10bz, _ := math.LegacyNewDec(10).Marshal()

	cases := map[string]struct {
		gogo               gogoproto.Message
		pulsar             proto.Message
		pulsarMarshalFails bool

		// this will fail in cases where a lossy encoding of an empty array to protobuf occurs. the unmarshalled bytes
		// represent the array as nil, and a subsequent marshal to JSON represent the array as null instead of empty.
		roundTripUnequal bool

		// pulsar does not support marshaling a math.Dec as anything except a string.  Therefore, we cannot unmarshal
		// a pulsar encoded Math.dec (the string representation of a Decimal) into a gogo Math.dec (expecting an int64).
		protoUnmarshalFails bool
	}{
		"beacon/v1/Params": {
			gogo: &beacontypes.Params{
				FeeRegister:         10,
				FeeRecord:           1,
				FeePurchaseStorage:  10,
				Denom:               "nund",
				DefaultStorageLimit: 100,
				MaxStorageLimit:     200,
			},
			pulsar: &beaconapi.Params{
				FeeRegister:         10,
				FeeRecord:           1,
				FeePurchaseStorage:  10,
				Denom:               "nund",
				DefaultStorageLimit: 100,
				MaxStorageLimit:     200,
			},
		},
		"beacon/v1/Beacon": {
			gogo: &beacontypes.Beacon{
				BeaconId:        1,
				Moniker:         "testmoniker",
				Name:            "testname",
				LastTimestampId: 124,
				FirstIdInState:  24,
				NumInState:      100,
				RegTime:         nowUint64,
				Owner:           addr1.String(),
			},
			pulsar: &beaconapi.Beacon{
				BeaconId:        1,
				Moniker:         "testmoniker",
				Name:            "testname",
				LastTimestampId: 124,
				FirstIdInState:  24,
				NumInState:      100,
				RegTime:         nowUint64,
				Owner:           addr1.String(),
			},
		},
		"beacon/v1/BeaconTimestamp": {
			gogo: &beacontypes.BeaconTimestamp{
				TimestampId: 1,
				SubmitTime:  nowUint64,
				Hash:        randomHash,
			},
			pulsar: &beaconapi.BeaconTimestamp{
				TimestampId: 1,
				SubmitTime:  nowUint64,
				Hash:        randomHash,
			},
		},
		"beacon/v1/BeaconStorageLimit": {
			gogo: &beacontypes.BeaconStorageLimit{
				BeaconId:     1,
				InStateLimit: 1234,
			},
			pulsar: &beaconapi.BeaconStorageLimit{
				BeaconId:     1,
				InStateLimit: 1234,
			},
		},
		"enterprise/v1/PurchaseOrderDecision": {
			gogo: &enterprisetypes.PurchaseOrderDecision{
				Signer:       addr1.String(),
				Decision:     2,
				DecisionTime: nowUint64,
			},
			pulsar: &enterpriseapi.PurchaseOrderDecision{
				Signer:       addr1.String(),
				Decision:     2,
				DecisionTime: nowUint64,
			},
		},
		"enterprise/v1/EnterpriseUndPurchaseOrder": {
			gogo: &enterprisetypes.EnterpriseUndPurchaseOrder{
				Id:             1,
				Purchaser:      addr1.String(),
				Amount:         testSdkCoin,
				Status:         1,
				RaiseTime:      nowUint64,
				CompletionTime: 0,
				Decisions: enterprisetypes.PurchaseOrderDecisions{
					enterprisetypes.PurchaseOrderDecision{
						Signer:       addr1.String(),
						Decision:     2,
						DecisionTime: nowUint64,
					},
				},
			},
			pulsar: &enterpriseapi.EnterpriseUndPurchaseOrder{
				Id:             1,
				Purchaser:      addr1.String(),
				Amount:         &testV1beta1Coin,
				Status:         1,
				RaiseTime:      nowUint64,
				CompletionTime: 0,
				Decisions: []*enterpriseapi.PurchaseOrderDecision{
					{
						Signer:       addr1.String(),
						Decision:     2,
						DecisionTime: nowUint64,
					},
				},
			},
		},
		"enterprise/v1/PurchaseOrders": {
			gogo: &enterprisetypes.PurchaseOrders{
				PurchaseOrders: []*enterprisetypes.EnterpriseUndPurchaseOrder{
					{
						Id:             1,
						Purchaser:      addr1.String(),
						Amount:         testSdkCoin,
						Status:         1,
						RaiseTime:      nowUint64,
						CompletionTime: 0,
						Decisions: enterprisetypes.PurchaseOrderDecisions{
							enterprisetypes.PurchaseOrderDecision{
								Signer:       addr1.String(),
								Decision:     2,
								DecisionTime: nowUint64,
							},
						},
					},
				},
			},
			pulsar: &enterpriseapi.PurchaseOrders{
				PurchaseOrders: []*enterpriseapi.EnterpriseUndPurchaseOrder{
					{
						Id:             1,
						Purchaser:      addr1.String(),
						Amount:         &testV1beta1Coin,
						Status:         1,
						RaiseTime:      nowUint64,
						CompletionTime: 0,
						Decisions: []*enterpriseapi.PurchaseOrderDecision{
							{
								Signer:       addr1.String(),
								Decision:     2,
								DecisionTime: nowUint64,
							},
						},
					},
				},
			},
		},
		"enterprise/v1/LockedUnd": {
			gogo: &enterprisetypes.LockedUnd{
				Owner:  addr1.String(),
				Amount: testSdkCoin,
			},
			pulsar: &enterpriseapi.LockedUnd{
				Owner:  addr1.String(),
				Amount: &testV1beta1Coin,
			},
		},
		"enterprise/v1/SpentEFUND": {
			gogo: &enterprisetypes.SpentEFUND{
				Owner:  addr1.String(),
				Amount: testSdkCoin,
			},
			pulsar: &enterpriseapi.SpentEFUND{
				Owner:  addr1.String(),
				Amount: &testV1beta1Coin,
			},
		},
		"enterprise/v1/EnterpriseUserAccount": {
			gogo: &enterprisetypes.EnterpriseUserAccount{
				Owner:         addr1.String(),
				LockedEfund:   testSdkCoin,
				GeneralSupply: testSdkCoin,
				SpentEfund:    testSdkCoin,
				Spendable:     testSdkCoin,
			},
			pulsar: &enterpriseapi.EnterpriseUserAccount{
				Owner:         addr1.String(),
				LockedEfund:   &testV1beta1Coin,
				GeneralSupply: &testV1beta1Coin,
				SpentEfund:    &testV1beta1Coin,
				Spendable:     &testV1beta1Coin,
			},
		},
		"enterprise/v1/WhitelistAddresses": {
			gogo:   &enterprisetypes.WhitelistAddresses{Addresses: []string{addr1.String()}},
			pulsar: &enterpriseapi.WhitelistAddresses{Addresses: []string{addr1.String()}},
		},
		"enterprise/v1/Params": {
			gogo: &enterprisetypes.Params{
				EntSigners:        addr1.String(),
				Denom:             "nund",
				MinAccepts:        1,
				DecisionTimeLimit: 1234,
			},
			pulsar: &enterpriseapi.Params{
				EntSigners:        addr1.String(),
				Denom:             "nund",
				MinAccepts:        1,
				DecisionTimeLimit: 1234,
			},
		},
		"stream/v1/Stream": {
			gogo: &streamtypes.Stream{
				Deposit:         testSdkCoin,
				FlowRate:        1234,
				LastOutflowTime: now,
				DepositZeroTime: now,
				Cancellable:     true,
			},
			pulsar: &streamapi.Stream{
				Deposit:         &testV1beta1Coin,
				FlowRate:        1234,
				LastOutflowTime: timestamppb.New(now),
				DepositZeroTime: timestamppb.New(now),
				Cancellable:     true,
			},
		},
		"stream/v1/Params": {
			gogo:               &streamtypes.Params{ValidatorFee: mathmod.LegacyNewDecWithPrec(1, 2)},
			pulsar:             &streamapi.Params{ValidatorFee: "0.010000000000000000"},
			pulsarMarshalFails: true,
		},
		"wrkchain/v1/WrkChain": {
			gogo: &wrkchaintypes.WrkChain{
				WrkchainId:   1,
				Moniker:      "testmoniker",
				Name:         "testname",
				Genesis:      randomHash,
				BaseType:     "cosmos",
				Lastblock:    1,
				NumBlocks:    2,
				LowestHeight: 1,
				RegTime:      nowUint64,
				Owner:        addr1.String(),
			},
			pulsar: &wrkchainapi.WrkChain{
				WrkchainId:   1,
				Moniker:      "testmoniker",
				Name:         "testname",
				Genesis:      randomHash,
				BaseType:     "cosmos",
				Lastblock:    1,
				NumBlocks:    2,
				LowestHeight: 1,
				RegTime:      nowUint64,
				Owner:        addr1.String(),
			},
		},
		"wrkchain/v1/WrkChainStorageLimit": {
			gogo: &wrkchaintypes.WrkChainStorageLimit{
				WrkchainId:   1,
				InStateLimit: 1000,
			},
			pulsar: &wrkchainapi.WrkChainStorageLimit{
				WrkchainId:   1,
				InStateLimit: 1000,
			},
		},
		"wrkchain/v1/WrkChainBlock": {
			gogo: &wrkchaintypes.WrkChainBlock{
				Height:     24,
				Blockhash:  randomHash,
				Parenthash: randomHash,
				Hash1:      randomHash,
				Hash2:      randomHash,
				Hash3:      randomHash,
				SubTime:    nowUint64,
			},
			pulsar: &wrkchainapi.WrkChainBlock{
				Height:     24,
				Blockhash:  randomHash,
				Parenthash: randomHash,
				Hash1:      randomHash,
				Hash2:      randomHash,
				Hash3:      randomHash,
				SubTime:    nowUint64,
			},
		},
		"wrkchain/v1/Params": {
			gogo: &wrkchaintypes.Params{
				FeeRegister:         1,
				FeeRecord:           2,
				FeePurchaseStorage:  3,
				Denom:               "nund",
				DefaultStorageLimit: 4,
				MaxStorageLimit:     5,
			},
			pulsar: &wrkchainapi.Params{
				FeeRegister:         1,
				FeeRecord:           2,
				FeePurchaseStorage:  3,
				Denom:               "nund",
				DefaultStorageLimit: 4,
				MaxStorageLimit:     5,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gogoBytes, err := encCfg.Amino.MarshalJSON(tc.gogo)
			require.NoError(t, err)

			pulsarBytes, err := aj.Marshal(tc.pulsar)
			if tc.pulsarMarshalFails {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			fmt.Printf("pulsar: %s\n", string(pulsarBytes))
			fmt.Printf("  gogo: %s\n", string(gogoBytes))
			require.Equal(t, string(gogoBytes), string(pulsarBytes))

			pulsarProtoBytes, err := proto.Marshal(tc.pulsar)
			require.NoError(t, err)

			gogoType := reflect.TypeOf(tc.gogo).Elem()
			newGogo := reflect.New(gogoType).Interface().(gogoproto.Message)

			err = encCfg.Codec.Unmarshal(pulsarProtoBytes, newGogo)
			if tc.protoUnmarshalFails {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			newGogoBytes, err := encCfg.Amino.MarshalJSON(newGogo)
			require.NoError(t, err)
			if tc.roundTripUnequal {
				require.NotEqual(t, string(gogoBytes), string(newGogoBytes))
				return
			}
			require.Equal(t, string(gogoBytes), string(newGogoBytes))

			// test amino json signer handler equivalence
			msg, ok := tc.gogo.(legacytx.LegacyMsg)
			if !ok {
				// not signable
				return
			}

			handlerOptions := signing_testutil.HandlerArgumentOptions{
				ChainID:       "test-chain",
				Memo:          "sometestmemo",
				Msg:           tc.pulsar,
				AccNum:        1,
				AccSeq:        2,
				SignerAddress: "signerAddress",
				Fee: &txv1beta1.Fee{
					Amount: []*v1beta1.Coin{{Denom: "uatom", Amount: "1000"}},
				},
			}

			signerData, txData, err := signing_testutil.MakeHandlerArguments(handlerOptions)
			require.NoError(t, err)

			handler := aminojson.NewSignModeHandler(aminojson.SignModeHandlerOptions{})
			signBz, err := handler.GetSignBytes(context.Background(), signerData, txData)
			require.NoError(t, err)

			legacyHandler := tx.NewSignModeLegacyAminoJSONHandler()
			txBuilder := encCfg.TxConfig.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs([]sdk.Msg{msg}...))
			txBuilder.SetMemo(handlerOptions.Memo)
			txBuilder.SetFeeAmount(sdk.Coins{sdk.NewInt64Coin("uatom", 1000)})
			theTx := txBuilder.GetTx()

			legacySigningData := signing.SignerData{
				ChainID:       handlerOptions.ChainID,
				Address:       handlerOptions.SignerAddress,
				AccountNumber: handlerOptions.AccNum,
				Sequence:      handlerOptions.AccSeq,
			}
			legacySignBz, err := legacyHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
				legacySigningData, theTx)
			require.NoError(t, err)
			require.Equal(t, string(legacySignBz), string(signBz))
		})
	}
}

func postFixPulsarMessage(msg proto.Message) {
	if m, ok := msg.(*authapi.ModuleAccount); ok {
		if m.BaseAccount == nil {
			m.BaseAccount = &authapi.BaseAccount{}
		}
		_, _, bz := testdata.KeyTestPubAddr()
		// always set address to a valid bech32 address
		text, _ := bech32.ConvertAndEncode("cosmos", bz)
		m.BaseAccount.Address = text

		// see negative test
		if len(m.Permissions) == 0 {
			m.Permissions = nil
		}
	}
}
