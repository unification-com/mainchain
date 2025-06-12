package enterprise

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/ed25519" // register so that it shows up in protoregistry.GlobalTypes
	"github.com/cosmos/cosmos-sdk/version"

	enterprisev1 "github.com/unification-com/mainchain/api/mainchain/enterprise/v1"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		// This is in place of func (AppModuleBasic) GetQueryCmd() *cobra.Command in module.go
		// and replaces client/cli/query.go
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: enterprisev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "EnterpriseUndPurchaseOrders",
					Use:       "orders",
					Short:     "Query Enterprise eFUND purchase orders with optional filters and pagination",
					Long:      "Query for a all paginated Enterprise eFUND purchase orders that match optional filters.",
					Example:   fmt.Sprintf("$ %s query enterprise orders --status status-completed --purchaser und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
				},
				{
					RpcMethod: "EnterpriseUndPurchaseOrder",
					Use:       "order [purchase_order_id]",
					Short:     "Query an eFUND Purchase Order for given ID",
					Example:   fmt.Sprintf("$ %s query enterpise order 24", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "purchase_order_id"},
					},
				},
				{
					RpcMethod: "LockedUndByAddress",
					Use:       "locked [owner]",
					Short:     "Query a given address's locked, usable eFUND",
					Long:      "Query a given address's locked, usable eFUND",
					Example:   fmt.Sprintf("$ %s query enterprise locked und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
					},
				},
				{
					RpcMethod: "TotalLocked",
					Use:       "total-locked",
					Short:     "Query total eFUND currently locked",
					Example:   fmt.Sprintf("$ %s query enterprise total-locked", version.AppName),
				},
				{
					RpcMethod: "Whitelist",
					Use:       "whitelist",
					Short:     "Query whitelisted addresses",
					Long:      "Query all addresses currently authorised to raise new purchase orders. Paginated",
					Example:   fmt.Sprintf("$ %s query enterprise whitelist", version.AppName),
				},
				{
					RpcMethod: "Whitelisted",
					Use:       "is-whitelisted [address]",
					Short:     "Query a given address's whitelist status",
					Long:      "Query a given address's whitelist status. If whitelisted, they are authorised to raise new eFUND purchase orders",
					Example:   fmt.Sprintf("$ %s query enterprise is-whitelisted und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "EnterpriseAccount",
					Use:       "account [address]",
					Short:     "Query a given address's Enterprise account",
					Long:      "Query a given address's detailed Enterprise account information",
					Example:   fmt.Sprintf("$ %s query enterprise account und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "SpentEFUNDByAddress",
					Use:       "spent-efund [address]",
					Short:     "Query a given address's total spent eFUND to date",
					Long:      "Query a given address's total spent eFUND to date. This is how much eFUND the account has used to pay for fees etc. and is now in general supply as FUND",
					Example:   fmt.Sprintf("$ %s query enterprise spent-efund und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "TotalSpentEFUND",
					Use:       "total-spent-efund",
					Short:     "Query the network total spent eFUND to date",
					Long:      "Query the network's total spent eFUND to date. This is how much eFUND has used to pay for fees etc. and is now in general supply as FUND",
					Example:   fmt.Sprintf("$ %s query enterprise total-spent-efund", version.AppName),
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
			Service: enterprisev1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UndPurchaseOrder",
					Use:       "purchase [amount] --from [purchaser]",
					Short:     "raise a new eFUND purchase order",
					Long:      "raise a new eFUND purchase order for the given amount",
					Example:   fmt.Sprintf("$ %s tx enterprise purchase 1000000000000nund --from mykey", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
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
