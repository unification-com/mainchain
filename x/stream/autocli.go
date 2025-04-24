package stream

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/ed25519" // register so that it shows up in protoregistry.GlobalTypes
	"github.com/cosmos/cosmos-sdk/version"

	streamv1 "github.com/unification-com/mainchain/api/mainchain/stream/v1"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		// This is in place of func (AppModuleBasic) GetQueryCmd() *cobra.Command in module.go
		// and replaces client/cli/query.go
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: streamv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CalculateFlowRate",
					Use:       "calculate",
					Short:     "Calculate the Flow Rate for given parameters",
					Long:      "Calculate the Flow Rate for given parameters coin, duration (frequency of the payment, e.g. every 2 months) and duration type (e.g. month)",
					Example:   fmt.Sprintf("$ %s query stream calculate --coin 1000000000nund --period month --duration 1", version.AppName),
				},
				{
					RpcMethod: "Streams",
					Use:       "streams",
					Short:     "Query all streams",
					Long:      "Query all streams, with pagination flags",
					Example:   fmt.Sprintf("$ %s query stream streams", version.AppName),
				},
				{
					RpcMethod: "AllStreamsForReceiver",
					Use:       "streams-receiver [receiver_addr]",
					Short:     "Query all streams for given receiver address",
					Long:      "Query all streams being sent to the given receiver address",
					Example:   fmt.Sprintf("$ %s query stream streams-receiver und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "receiver_addr"},
					},
				},
				{
					RpcMethod: "StreamByReceiverSender",
					Use:       "stream [receiver_addr] [sender_addr]",
					Short:     "Query a stream for receiver/sender pair",
					Example:   fmt.Sprintf("$ %s query stream stream und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "receiver_addr"},
						{ProtoField: "sender_addr"},
					},
				},
				{
					RpcMethod: "StreamReceiverSenderCurrentFlow",
					Use:       "stream-flow [receiver_addr] [sender_addr]",
					Short:     "Query a stream's current flow for receiver/sender pair",
					Long:      "Query a stream's current flow rate for receiver/sender pair. This will be zero if the stream has expired",
					Example:   fmt.Sprintf("$ %s query stream stream-flow und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "receiver_addr"},
						{ProtoField: "sender_addr"},
					},
				},
				{
					RpcMethod: "AllStreamsForSender",
					Use:       "streams-sender [sender_addr]",
					Short:     "Query all streams for given sender address",
					Long:      "Query all streams being sent from the given sender address",
					Example:   fmt.Sprintf("$ %s query stream streams-sender und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "sender_addr"},
					},
				},
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current enterprise parameters information",
					Long:      "Query values set as enterprise parameters.",
				},
			},
		},
		// Note - we're still using func (AppModuleBasic) GetTxCmd() *cobra.Command in module.go for Tx commands
		// this is here just for an example for future use
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: streamv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateStream",
					Use:       "create [receiver] [deposit] [flow_rate] --from [sender]",
					Short:     "create a new payment stream",
					Long:      "create a new payment stream using the provided options",
					Example:   fmt.Sprintf("$ %s tx stream create und173qnkw458p646fahmd53xa45vqqvga7kyu6ryy 777000000000nund 299768 --from mykey", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "receiver"},
						{ProtoField: "deposit"},
						{ProtoField: "flow_rate"},
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
