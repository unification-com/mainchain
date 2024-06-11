package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/stream/types"
	"time"
)

func (s *KeeperTestSuite) TestMsgServerUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr bool
		expErrMsg string
	}{
		{
			name: "set valid params",
			request: &types.MsgUpdateParams{
				Authority: s.app.StreamKeeper.GetAuthority(),
				Params: types.Params{
					ValidatorFee: sdk.NewDecWithPrec(24, 2),
				},
			},
			expectErr: false,
			expErrMsg: "",
		},
		{
			name: "set invalid authority",
			request: &types.MsgUpdateParams{
				Authority: "foo",
			},
			expectErr: true,
			expErrMsg: "invalid authority",
		},
		{
			name: "set invalid params > 100%",
			request: &types.MsgUpdateParams{
				Authority: s.app.StreamKeeper.GetAuthority(),
				Params: types.Params{
					ValidatorFee: sdk.NewDecWithPrec(101, 2),
				},
			},
			expectErr: true,
			expErrMsg: "validator fee cannot be greater than 100",
		},
		{
			name: "set invalid params negative value",
			request: &types.MsgUpdateParams{
				Authority: s.app.StreamKeeper.GetAuthority(),
				Params: types.Params{
					ValidatorFee: sdk.NewDecWithPrec(-1, 2),
				},
			},
			expectErr: true,
			expErrMsg: "validator fee cannot be negative",
		},
		{
			name: "set invalid params nil value",
			request: &types.MsgUpdateParams{
				Authority: s.app.StreamKeeper.GetAuthority(),
				Params: types.Params{
					ValidatorFee: sdk.Dec{},
				},
			},
			expectErr: true,
			expErrMsg: "validator fee cannot be nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateParams(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgServerCreateStream() {
	testCases := []struct {
		name      string
		request   *types.MsgCreateStream
		expResult *types.MsgCreateStreamResponse
		expectErr bool
		expErrMsg string
	}{
		{
			name: "create valid stream",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: &types.MsgCreateStreamResponse{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expectErr: false,
			expErrMsg: "",
		},
		{
			name: "invalid - stream exists",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "use update stream msg to modify existing stream",
		},
		{
			name: "invalid - empty sender address",
			request: &types.MsgCreateStream{
				Sender:   "",
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name: "invalid - bad sender address",
			request: &types.MsgCreateStream{
				Sender:   "rubbish",
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},

		{
			name: "invalid - empty receiver address",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: "",
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name: "invalid - bad receiver address",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: "rubbish",
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name: "invalid - receiver address blocked",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: "und17xpfvakm2amg962yls6f84z3kell8c5lhuyfdm", // module account with no banking permissions
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "not allowed to receive funds",
		},
		{
			name: "invalid - receiver address same as sender address",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[0].String(), // module account with no banking permissions
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "sender and receiver cannot be same address",
		},
		{
			name: "invalid - nil deposit",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.Coin{},
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "deposit must be > zero",
		},
		{
			name: "invalid - zero deposit",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.NewInt64Coin("stake", 0),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "deposit must be > zero",
		},
		{
			name: "invalid - zero flow rate",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 0,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "flow rate must be > zero",
		},
		{
			name: "invalid - negative flow rate",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: -1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "flow rate must be > zero",
		},
		{
			name: "invalid - duration < 60 seconds",
			request: &types.MsgCreateStream{
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.NewInt64Coin("stake", 10),
				FlowRate: 1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "calculated duration too short. Must be > 1 minute",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res, err := s.msgServer.CreateStream(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResult, res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgServerClaimStream() {
	testCases := []struct {
		name      string
		create    *types.MsgCreateStream
		claim     *types.MsgClaimStream
		expResult *types.MsgClaimStreamResponse
		expectErr bool
		expErrMsg string
	}{
		{
			name: "valid claim",
			create: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			claim: &types.MsgClaimStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
			},
			expResult: &types.MsgClaimStreamResponse{
				TotalClaimed:     sdk.NewInt64Coin("stake", 1000),
				StreamPayment:    sdk.NewInt64Coin("stake", 990),
				ValidatorFee:     sdk.NewInt64Coin("stake", 10), // default fee is 1%
				RemainingDeposit: sdk.NewInt64Coin("stake", 0),
			},
			expectErr: false,
			expErrMsg: "",
		},
		{
			name:   "invalid claim - bad receiver address",
			create: nil,
			claim: &types.MsgClaimStream{
				Sender:   s.addrs[0].String(),
				Receiver: "rubbish",
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid claim - empty receiver address",
			create: nil,
			claim: &types.MsgClaimStream{
				Sender:   s.addrs[0].String(),
				Receiver: "",
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid claim - bad sender address",
			create: nil,
			claim: &types.MsgClaimStream{
				Sender:   "rubbish",
				Receiver: s.addrs[1].String(),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid claim - empty sender address",
			create: nil,
			claim: &types.MsgClaimStream{
				Sender:   "",
				Receiver: s.addrs[1].String(),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid claim - stream does not exist",
			create: nil,
			claim: &types.MsgClaimStream{
				Sender:   s.addrs[8].String(), // created earlier
				Receiver: s.addrs[9].String(),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "stream not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tCtx := s.ctx
			nowTime := time.Unix(time.Now().Unix(), 0).UTC()

			if tc.create != nil {
				createTime := time.Unix(nowTime.Unix()-1000, 0).UTC()
				tCtx = tCtx.WithBlockTime(createTime)
				_, _ = s.msgServer.CreateStream(tCtx, tc.create)
			}

			tCtx = tCtx.WithBlockTime(nowTime)

			res, err := s.msgServer.ClaimStream(tCtx, tc.claim)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResult, res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgServerTopUpDeposit() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()

	testCases := []struct {
		name      string
		create    *types.MsgCreateStream
		topup     *types.MsgTopUpDeposit
		expResult *types.MsgTopUpDepositResponse
		expectErr bool
		expErrMsg string
	}{
		{
			name: "valid topup",
			create: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			topup: &types.MsgTopUpDeposit{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
			},
			expResult: &types.MsgTopUpDepositResponse{
				DepositAmount:   sdk.NewInt64Coin("stake", 1000),
				CurrentDeposit:  sdk.NewInt64Coin("stake", 2000),
				DepositZeroTime: time.Unix(nowTime.Unix()+1500, 0).UTC(),
			},
			expectErr: false,
			expErrMsg: "",
		},
		{
			name: "invalid topup - denom mistmatch",
			create: &types.MsgCreateStream{
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			topup: &types.MsgTopUpDeposit{
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.NewInt64Coin("notstake", 1000),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "top up denom does not match stream denom",
		},
		{
			name:   "invalid topup - bad sender address",
			create: nil,
			topup: &types.MsgTopUpDeposit{
				Sender:   "rubbish",
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid topup - empty sender address",
			create: nil,
			topup: &types.MsgTopUpDeposit{
				Sender:   "",
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid topup - bad receiver address",
			create: nil,
			topup: &types.MsgTopUpDeposit{
				Sender:   s.addrs[0].String(),
				Receiver: "rubbish",
				Deposit:  sdk.NewInt64Coin("stake", 1000),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid topup - empty receiver address",
			create: nil,
			topup: &types.MsgTopUpDeposit{
				Sender:   s.addrs[1].String(),
				Receiver: "",
				Deposit:  sdk.NewInt64Coin("stake", 1000),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid topup - zero deposit",
			create: nil,
			topup: &types.MsgTopUpDeposit{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 0),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "deposit must be > zero",
		},
		{
			name:   "invalid topup - nil deposit",
			create: nil,
			topup: &types.MsgTopUpDeposit{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.Coin{},
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "deposit must be > zero",
		},
		{
			name:   "invalid topup - stream not exist",
			create: nil,
			topup: &types.MsgTopUpDeposit{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[9].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "stream not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.create != nil {
				// created 500 seconds ago
				createTime := time.Unix(nowTime.Unix()-500, 0).UTC()
				tCtx = tCtx.WithBlockTime(createTime)
				_, _ = s.msgServer.CreateStream(tCtx, tc.create)
			}

			tCtx = tCtx.WithBlockTime(nowTime)

			res, err := s.msgServer.TopUpDeposit(tCtx, tc.topup)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResult, res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgServerUpdateFlowRate() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()

	testCases := []struct {
		name      string
		create    *types.MsgCreateStream
		flowRate  *types.MsgUpdateFlowRate
		expResult *types.MsgUpdateFlowRateResponse
		expectErr bool
		expErrMsg string
	}{
		{
			name: "valid flow rate",
			create: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				FlowRate: 2,
			},
			expResult: &types.MsgUpdateFlowRateResponse{
				FlowRate: 2,
			},
			expectErr: false,
			expErrMsg: "",
		},
		{
			name:   "invalid flow rate - bad sender address",
			create: nil,
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   "rubbish",
				Receiver: s.addrs[1].String(),
				FlowRate: 2,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid flow rate - empty sender address",
			create: nil,
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   "",
				Receiver: s.addrs[1].String(),
				FlowRate: 2,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid flow rate - bad receiver address",
			create: nil,
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   s.addrs[1].String(),
				Receiver: "rubbish",
				FlowRate: 2,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid flow rate - empty receiver address",
			create: nil,
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   s.addrs[1].String(),
				Receiver: "",
				FlowRate: 2,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid flow rate - zero flow rate",
			create: nil,
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				FlowRate: 0,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "flow rate must be > zero",
		},
		{
			name:   "invalid flow rate - negative flow rate",
			create: nil,
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				FlowRate: -1,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "flow rate must be > zero",
		},
		{
			name:   "invalid flow rate - stream not exist",
			create: nil,
			flowRate: &types.MsgUpdateFlowRate{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[9].String(),
				FlowRate: 2,
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "stream not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.create != nil {
				// created 500 seconds ago
				createTime := time.Unix(nowTime.Unix()-500, 0).UTC()
				tCtx = tCtx.WithBlockTime(createTime)
				_, _ = s.msgServer.CreateStream(tCtx, tc.create)
			}

			tCtx = tCtx.WithBlockTime(nowTime)

			res, err := s.msgServer.UpdateFlowRate(tCtx, tc.flowRate)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResult, res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgServerCancelStream() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()

	testCases := []struct {
		name      string
		create    *types.MsgCreateStream
		cancel    *types.MsgCancelStream
		expResult *types.MsgCancelStreamResponse
		expectErr bool
		expErrMsg string
	}{
		{
			name: "valid cancel",
			create: &types.MsgCreateStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewInt64Coin("stake", 1000),
				FlowRate: 1,
			},
			cancel: &types.MsgCancelStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
			},
			expResult: &types.MsgCancelStreamResponse{},
			expectErr: false,
			expErrMsg: "",
		},
		{
			name:   "invalid cancel - bad sender address",
			create: nil,
			cancel: &types.MsgCancelStream{
				Sender:   "rubbish",
				Receiver: s.addrs[1].String(),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid cancel - empty sender address",
			create: nil,
			cancel: &types.MsgCancelStream{
				Sender:   "",
				Receiver: s.addrs[1].String(),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid cancel - bad receiver address",
			create: nil,
			cancel: &types.MsgCancelStream{
				Sender:   s.addrs[0].String(),
				Receiver: "rubbish",
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid cancel - empty receiver address",
			create: nil,
			cancel: &types.MsgCancelStream{
				Sender:   s.addrs[0].String(),
				Receiver: "",
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:   "invalid cancel - stream not exist",
			create: nil,
			cancel: &types.MsgCancelStream{
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[9].String(),
			},
			expResult: nil,
			expectErr: true,
			expErrMsg: "stream not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.create != nil {
				// created 500 seconds ago
				createTime := time.Unix(nowTime.Unix()-500, 0).UTC()
				tCtx = tCtx.WithBlockTime(createTime)
				_, _ = s.msgServer.CreateStream(tCtx, tc.create)
			}

			tCtx = tCtx.WithBlockTime(nowTime)

			res, err := s.msgServer.CancelStream(tCtx, tc.cancel)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResult, res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgServerCancelStream_Fail_NotCancellable() {

	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	expStream := types.Stream{
		Deposit:         sdk.NewInt64Coin("stake", 1000),
		FlowRate:        1,
		LastOutflowTime: nowTime,
		DepositZeroTime: time.Unix(0, 0).UTC(),
		Cancellable:     false,
	}

	err := s.app.StreamKeeper.SetStream(s.ctx, s.addrs[1], s.addrs[0], expStream)
	s.Require().NoError(err)

	cancelMsg := &types.MsgCancelStream{
		Sender:   s.addrs[0].String(),
		Receiver: s.addrs[1].String(),
	}

	resp, err := s.msgServer.CancelStream(s.ctx, cancelMsg)
	s.Require().Nil(resp)
	s.Require().ErrorContains(err, "cannot be cancelled")

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().True(ok)
	s.Require().Equal(expStream, stream)
}
