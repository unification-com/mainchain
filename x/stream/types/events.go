package types

const (
	EventTypeCreateStreamAction = "create_stream"
	EventTypeDepositToStream    = "stream_deposit"
	EventTypeClaimStreamAction  = "claim_stream"
	EventTypeUpdateFlowRate     = "update_flow_rate"
	EventTypeStreamCancelled    = "cancel_stream"

	AttributeKeyStreamSender        = "sender"
	AttributeKeyStreamReceiver      = "receiver"
	AttributeKeyFlowRate            = "flow_rate"
	AttributeKeyAmountDeposited     = "amount_deposited"
	AttributeKeyDepositDuration     = "deposit_duration"
	AttributeKeyDepositZeroTime     = "deposit_zero_time"
	AttributeKeyClaimAmountReceived = "amount_received"
	AttributeKeyClaimValidatorFee   = "validator_fee"
	AttributeKeyClaimTotal          = "claim_total"
	AttributeKeyOldFlowRate         = "old_flow_rate"
	AttributeKeyNewFlowRate         = "new_flow_rate"
	AttributeKeyRemainingDeposit    = "remaining_deposit"
	AttributeKeyRefundAmount        = "refund_amount"
)
