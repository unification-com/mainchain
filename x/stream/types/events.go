package types

const (
	EventTypeCreateStreamAction = "create_stream"
	EventTypeDepositToStream    = "stream_deposit"
	EventTypeClaimStreamAction  = "claim_stream"
	EventTypeUpdateFlowRate     = "update_flow_rate"
	EventTypeStreamCancelled    = "cancel_stream"

	AttributeKeyStreamId                  = "stream_id"
	AttributeKeyStreamSender              = "sender"
	AttributeKeyStreamReceiver            = "receiver"
	AttributeKeyStreamFlowRate            = "flow_rate"
	AttributeKeyStreamDepositAmount       = "deposit"
	AttributeKeyStreamDepositDuration     = "deposit_duration"
	AttributeKeyStreamDepositZeroTime     = "deposit_zero_time"
	AttributeKeyStreamClaimAmountReceived = "amount_received"
	AttributeKeyStreamClaimValidatorFee   = "validator_fee"
	AttributeKeyStreamClaimTotal          = "claim_total"
	AttributeKeyOldFlowRate               = "old_flow_rate"
	AttributeKeyNewFlowRate               = "new_flow_rate"
)
