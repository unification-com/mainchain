package beacon

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/ed25519" // register so that it shows up in protoregistry.GlobalTypes
	"github.com/cosmos/cosmos-sdk/version"

	beaconv1 "github.com/unification-com/mainchain/api/mainchain/beacon/v1"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		// This is in place of func (AppModuleBasic) GetQueryCmd() *cobra.Command in module.go
		// and replaces client/cli/query.go
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: beaconv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Beacon",
					Use:       "beacon [beacon_id]",
					Short:     "Query a BEACON for given ID",
					Long:      "Query details about an individual BEACON.",
					Example:   fmt.Sprintf("$ %s query beacon beacon 1", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "beacon_id"},
					},
				},
				{
					RpcMethod: "BeaconsFiltered",
					Short:     "Query for all BEACONS",
					Long:      "Query details about all BEACONS on a network, with optional filters for owner and moniker.",
					Example:   fmt.Sprintf("$ %s query beacon beacons-filtered --owner und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy", version.AppName),
				},
				{
					RpcMethod: "BeaconTimestamp",
					Use:       "beacon-timestamp [beacon_id] [timestamp_id]",
					Short:     "Query a BEACON for given ID and timestamp ID to retrieve recorded timestamp",
					Example:   fmt.Sprintf("$ %s query beacon beacon-timestamp 1 24", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "beacon_id"},
						{ProtoField: "timestamp_id"},
					},
				},
				{
					RpcMethod: "BeaconStorage",
					Use:       "storage [beacon_id]",
					Short:     "Query a BEACON's storage capacity for given ID",
					Long:      "Query storage details about an individual BEACON.",
					Example:   fmt.Sprintf("$ %s query beacon storage 1", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "beacon_id"},
					},
				},
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current beacon parameters information",
					Long:      "Query values set as beacon parameters.",
				},
			},
		},
		// Note - we're still using func (AppModuleBasic) GetTxCmd() *cobra.Command in module.go for Tx commands
		// this is here just for an example for future use
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: beaconv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "RecordBeaconTimestamp",
					Use:       "record [beacon_id] --from [owner]",
					Short:     "record a BEACON's timestamp hash",
					Long:      "record a BEACON's timestamp hash along with a time submitted",
					Example:   fmt.Sprintf("$ %s tx beacon record 1 --hash d04b98f48e8 --subtime 1234356 --from mykey", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "beacon_id"},
					},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
			},
			EnhanceCustomCommand: false, // use custom commands only until v0.51
		},
	}
}
