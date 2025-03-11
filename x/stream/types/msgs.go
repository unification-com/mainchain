package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CreateStreamAction   = "create_stream"
	ClaimStreamAction    = "claim_stream"
	TopUpDepositAction   = "top_up_deposit"
	UpdateFlowRateAction = "update_flow_rate"
	CancelStreamAction   = "cancel_stream"
)

var (
	_ sdk.Msg = &MsgCreateStream{}
	_ sdk.Msg = &MsgClaimStream{}
	_ sdk.Msg = &MsgTopUpDeposit{}
	_ sdk.Msg = &MsgUpdateFlowRate{}
	_ sdk.Msg = &MsgCancelStream{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// --- Create New Stream Msg ---

// NewMsgCreateStream is a constructor function for MsgCreateStream
func NewMsgCreateStream(
	deposit sdk.Coin,
	flowRate int64,
	receiver,
	sender sdk.AccAddress) *MsgCreateStream {
	return &MsgCreateStream{
		Receiver: receiver.String(),
		Sender:   sender.String(),
		Deposit:  deposit,
		FlowRate: flowRate,
	}
}

// Route should return the name of the module
func (msg MsgCreateStream) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateStream) Type() string { return CreateStreamAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateStream) ValidateBasic() error {
	_, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return accErr
	}
	_, accErr = sdk.AccAddressFromBech32(msg.Receiver)
	if accErr != nil {
		return accErr
	}

	if msg.Deposit.IsNil() || msg.Deposit.IsNegative() || msg.Deposit.IsZero() {
		return errorsmod.Wrap(ErrInvalidData, "deposit must be > zero")
	}

	if msg.FlowRate < 1 {
		return errorsmod.Wrap(ErrInvalidData, "flow rate must be > zero")
	}

	if msg.Sender == msg.Receiver {
		return errorsmod.Wrap(ErrInvalidData, "receiver cannot be same as sender")
	}

	duration := CalculateDuration(msg.Deposit, msg.FlowRate)

	if duration < 60 {
		return errorsmod.Wrap(ErrInvalidData, "calculated duration too short. Must be > 1 minute")
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgCreateStream) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// --- Claim Stream By sender & receiver Msg ---

// NewMsgClaimStream is a constructor function for MsgClaimStream
func NewMsgClaimStream(
	receiver sdk.AccAddress,
	sender sdk.AccAddress) *MsgClaimStream {
	return &MsgClaimStream{
		Receiver: receiver.String(),
		Sender:   sender.String(),
	}
}

// Route should return the name of the module
func (msg MsgClaimStream) Route() string { return RouterKey }

// Type should return the action
func (msg MsgClaimStream) Type() string { return ClaimStreamAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgClaimStream) ValidateBasic() error {
	_, accErr := sdk.AccAddressFromBech32(msg.Receiver)
	if accErr != nil {
		return accErr
	}

	_, accErr = sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return accErr
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgClaimStream) GetSigners() []sdk.AccAddress {
	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{receiver}
}

// --- Top up Deposit Msg ---

// NewMsgTopUpDeposit is a constructor function for MsgTopUpDeposit
func NewMsgTopUpDeposit(
	receiver, sender sdk.AccAddress,
	deposit sdk.Coin,
) *MsgTopUpDeposit {
	return &MsgTopUpDeposit{
		Receiver: receiver.String(),
		Sender:   sender.String(),
		Deposit:  deposit,
	}
}

// Route should return the name of the module
func (msg MsgTopUpDeposit) Route() string { return RouterKey }

// Type should return the action
func (msg MsgTopUpDeposit) Type() string { return TopUpDepositAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgTopUpDeposit) ValidateBasic() error {
	_, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return accErr
	}

	_, accErr = sdk.AccAddressFromBech32(msg.Receiver)
	if accErr != nil {
		return accErr
	}

	if msg.Deposit.IsNil() || msg.Deposit.IsNegative() || msg.Deposit.IsZero() {
		return errorsmod.Wrap(ErrInvalidData, "deposit must be > zero")
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgTopUpDeposit) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// --- Update Flow Rate Msg ---

// NewMsgUpdateFlowRate is a constructor function for MsgUpdateFlowRate
func NewMsgUpdateFlowRate(
	receiver, sender sdk.AccAddress,
	flowRate int64,
) *MsgUpdateFlowRate {
	return &MsgUpdateFlowRate{
		Receiver: receiver.String(),
		Sender:   sender.String(),
		FlowRate: flowRate,
	}
}

// Route should return the name of the module
func (msg MsgUpdateFlowRate) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateFlowRate) Type() string { return UpdateFlowRateAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateFlowRate) ValidateBasic() error {
	_, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return accErr
	}

	_, accErr = sdk.AccAddressFromBech32(msg.Receiver)
	if accErr != nil {
		return accErr
	}

	if msg.FlowRate < 1 {
		return errorsmod.Wrap(ErrInvalidData, "flow rate must be > zero")
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgUpdateFlowRate) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// --- Cancel Stream Msg ---

// NewMsgCancelStream is a constructor function for MsgCancelStream
func NewMsgCancelStream(
	reciever,
	sender sdk.AccAddress) *MsgCancelStream {
	return &MsgCancelStream{
		Receiver: reciever.String(),
		Sender:   sender.String(),
	}
}

// Route should return the name of the module
func (msg MsgCancelStream) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCancelStream) Type() string { return CancelStreamAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelStream) ValidateBasic() error {
	_, accErr := sdk.AccAddressFromBech32(msg.Sender)
	if accErr != nil {
		return accErr
	}

	_, accErr = sdk.AccAddressFromBech32(msg.Receiver)
	if accErr != nil {
		return accErr
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgCancelStream) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// --- Modify Params Msg Type ---

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := m.Params.Validate(); err != nil {
		return err
	}

	return nil
}
