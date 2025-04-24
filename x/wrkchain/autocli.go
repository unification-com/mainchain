package wrkchain

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/ed25519" // register so that it shows up in protoregistry.GlobalTypes
	"github.com/cosmos/cosmos-sdk/version"

	wrkchainv1 "github.com/unification-com/mainchain/api/mainchain/wrkchain/v1"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		// This is in place of func (AppModuleBasic) GetQueryCmd() *cobra.Command in module.go
		// and replaces client/cli/query.go
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: wrkchainv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "WrkChain",
					Use:       "wrkchain [wrkchain_id]",
					Short:     "Query a WrkChain for given ID",
					Long:      "Query details about an individual WrkChain.",
					Example:   fmt.Sprintf("$ %s query wrkchain wrkchain 1", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wrkchain_id"},
					},
				},
				{
					RpcMethod: "WrkChainsFiltered",
					Use:       "wrkchains-filtered",
					Short:     "Query for all WrkChains",
					Long:      "Query details about all WrkChains on a network, with optional filters for owner and moniker.",
					Example:   fmt.Sprintf("$ %s query wrkchain wrkchains-filtered --owner und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy", version.AppName),
				},
				{
					RpcMethod: "WrkChainBlock",
					Use:       "wrkchain-block [wrkchain_id] [wc_height]",
					Short:     "Query a WrkChain for given ID and WrkChain height to retrieve recorded block hash data",
					Example:   fmt.Sprintf("$ %s query wrkchain wrkchain-block 1 24", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wrkchain_id"},
						{ProtoField: "wc_height"},
					},
				},
				{
					RpcMethod: "WrkChainStorage",
					Use:       "storage [wrkchain_id]",
					Short:     "Query a WrkChain's storage capacity for given ID",
					Long:      "Query storage details about an individual WrkChain.",
					Example:   fmt.Sprintf("$ %s query wrkchain storage 1", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wrkchain_id"},
					},
				},
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current wrkchain parameters information",
					Long:      "Query values set as wrkchain parameters.",
				},
			},
		},
		// Note - we're still using func (AppModuleBasic) GetTxCmd() *cobra.Command in module.go for Tx commands
		// this is here just for an example for future use
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: wrkchainv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "RecordWrkChainBlock",
					Use:       "record [wrkchain_id] --from [owner]",
					Short:     "record a WrkChain's block hashes",
					Long:      "record a WrkChain's block hash along with optional additional hashes such as parent block hash",
					Example:   fmt.Sprintf("$ %s tx wrkchain record 1 --wc_height 24 --block_hash d04b98f48e8 --parent_hash f8bcc15c6ae --from mykey", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wrkchain_id"},
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
