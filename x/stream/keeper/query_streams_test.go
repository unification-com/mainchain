package keeper_test

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/unification-com/mainchain/x/stream/types"
)

func (s *KeeperTestSuite) TestQueryStreams() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(nowTime)

	for i := int64(1); i <= 10; i++ {
		deposit := sdk.NewInt64Coin("stake", 1000)
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[i-1], s.addrs[i], deposit, 1)
		s.Require().NoError(err)
		_, err = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[i-1], s.addrs[i], deposit)
		s.Require().NoError(err)
	}

	expResp := &types.QueryStreamsResponse{
		Streams: []*types.StreamResult{
			{
				Receiver: s.addrs[0].String(),
				Sender:   s.addrs[1].String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
			{
				Receiver: s.addrs[1].String(),
				Sender:   s.addrs[2].String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
		},
	}

	req := &types.QueryStreamsRequest{
		Pagination: &query.PageRequest{
			Limit: 2,
		},
	}

	resp, err := s.app.StreamKeeper.Streams(tCtx, req)
	s.Require().NoError(err)
	s.Require().Equal(expResp.Streams, resp.Streams)
}

func (s *KeeperTestSuite) TestQueryStreamByReceiverSender() {

	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(nowTime)

	successSender := s.addrs[0]
	successReceiver := s.addrs[1]

	testCases := []struct {
		name      string
		sender    sdk.AccAddress
		receiver  sdk.AccAddress
		query     *types.QueryStreamByReceiverSenderRequest
		flowRate  int64
		deposit   sdk.Coin
		expResp   *types.QueryStreamByReceiverSenderResponse
		expErr    bool
		expErrMsg string
	}{
		{
			name:     "success",
			sender:   successSender,
			receiver: successReceiver,
			flowRate: 1,
			deposit:  sdk.NewInt64Coin("stake", 1000),
			query: &types.QueryStreamByReceiverSenderRequest{
				SenderAddr:   successSender.String(),
				ReceiverAddr: successReceiver.String(),
			},
			expResp: &types.QueryStreamByReceiverSenderResponse{
				Stream: types.StreamResult{
					Receiver: successReceiver.String(),
					Sender:   successSender.String(),
					Stream: &types.Stream{
						Deposit:         sdk.NewInt64Coin("stake", 1000),
						FlowRate:        1,
						LastOutflowTime: nowTime,
						DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
						Cancellable:     true,
					},
				},
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name: "fail - bad sender address",
			query: &types.QueryStreamByReceiverSenderRequest{
				SenderAddr:   "rubbish",
				ReceiverAddr: s.addrs[0].String(),
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name: "fail - empty sender address",
			query: &types.QueryStreamByReceiverSenderRequest{
				SenderAddr:   "",
				ReceiverAddr: s.addrs[0].String(),
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name: "fail - bad receiver address",
			query: &types.QueryStreamByReceiverSenderRequest{
				SenderAddr:   s.addrs[0].String(),
				ReceiverAddr: "rubbish",
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name: "fail - empty receiver address",
			query: &types.QueryStreamByReceiverSenderRequest{
				SenderAddr:   s.addrs[0].String(),
				ReceiverAddr: "",
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name: "fail - stream does not exist",
			query: &types.QueryStreamByReceiverSenderRequest{
				SenderAddr:   s.addrs[7].String(),
				ReceiverAddr: s.addrs[9].String(),
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "stream not found",
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			if tc.sender != nil {
				// create & deposit
				_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.deposit, tc.flowRate)
				_, _ = s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.deposit)
			}

			resp, err := s.app.StreamKeeper.StreamByReceiverSender(tCtx, tc.query)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResp, resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQueryStreamReceiverSenderCurrentFlow() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(nowTime)

	successSender := s.addrs[0]
	successReceiver := s.addrs[1]

	testCases := []struct {
		name        string
		sender      sdk.AccAddress
		receiver    sdk.AccAddress
		query       *types.QueryStreamReceiverSenderCurrentFlowRequest
		flowRate    int64
		deposit     sdk.Coin
		queryFuture int64
		expResp     *types.QueryStreamReceiverSenderCurrentFlowResponse
		expErr      bool
		expErrMsg   string
	}{
		{
			name:        "success - not expired",
			sender:      successSender,
			receiver:    successReceiver,
			flowRate:    1,
			deposit:     sdk.NewInt64Coin("stake", 1000),
			queryFuture: 0,
			query: &types.QueryStreamReceiverSenderCurrentFlowRequest{
				SenderAddr:   successSender.String(),
				ReceiverAddr: successReceiver.String(),
			},
			expResp: &types.QueryStreamReceiverSenderCurrentFlowResponse{
				ConfiguredFlowRate: 1,
				CurrentFlowRate:    1,
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name:        "success - expired",
			queryFuture: 1001,
			query: &types.QueryStreamReceiverSenderCurrentFlowRequest{
				SenderAddr:   successSender.String(),
				ReceiverAddr: successReceiver.String(),
			},
			expResp: &types.QueryStreamReceiverSenderCurrentFlowResponse{
				ConfiguredFlowRate: 1,
				CurrentFlowRate:    0,
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name:        "fail - bad sender address",
			queryFuture: 0,
			query: &types.QueryStreamReceiverSenderCurrentFlowRequest{
				SenderAddr:   "rubbish",
				ReceiverAddr: s.addrs[0].String(),
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:        "fail - empty sender address",
			queryFuture: 0,
			query: &types.QueryStreamReceiverSenderCurrentFlowRequest{
				SenderAddr:   "",
				ReceiverAddr: s.addrs[0].String(),
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:        "fail - bad receiver address",
			queryFuture: 0,
			query: &types.QueryStreamReceiverSenderCurrentFlowRequest{
				SenderAddr:   s.addrs[0].String(),
				ReceiverAddr: "rubbish",
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "decoding bech32 failed",
		},
		{
			name:        "fail - empty receiver address",
			queryFuture: 0,
			query: &types.QueryStreamReceiverSenderCurrentFlowRequest{
				SenderAddr:   s.addrs[0].String(),
				ReceiverAddr: "",
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "empty address string is not allowed",
		},
		{
			name:        "fail - stream does not exist",
			queryFuture: 0,
			query: &types.QueryStreamReceiverSenderCurrentFlowRequest{
				SenderAddr:   s.addrs[7].String(),
				ReceiverAddr: s.addrs[9].String(),
			},
			expResp:   nil,
			expErr:    true,
			expErrMsg: "stream not found",
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			if tc.sender != nil {
				// create & deposit
				_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.deposit, tc.flowRate)
				_, _ = s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.deposit)
			}

			queryTime := time.Unix(nowTime.Unix()+tc.queryFuture, 0).UTC()
			tCtx = tCtx.WithBlockTime(queryTime)

			resp, err := s.app.StreamKeeper.StreamReceiverSenderCurrentFlow(tCtx, tc.query)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResp, resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQueryAllStreamsForReceiver_Success() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(nowTime)

	for i := int64(1); i <= 10; i++ {
		deposit := sdk.NewInt64Coin("stake", 1000)
		receiverAddr := s.addrs[i-1]
		if i < 4 {
			receiverAddr = s.addrs[0]
		}
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, receiverAddr, s.addrs[i], deposit, 1)
		s.Require().NoError(err)
		_, err = s.app.StreamKeeper.AddDeposit(tCtx, receiverAddr, s.addrs[i], deposit)
		s.Require().NoError(err)
	}

	expResp := &types.QueryAllStreamsForReceiverResponse{
		Streams: []*types.StreamResult{
			{
				Receiver: s.addrs[0].String(),
				Sender:   s.addrs[1].String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
			{
				Receiver: s.addrs[0].String(),
				Sender:   s.addrs[2].String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
			{
				Receiver: s.addrs[0].String(),
				Sender:   s.addrs[3].String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
		},
	}

	req := &types.QueryAllStreamsForReceiverRequest{
		ReceiverAddr: s.addrs[0].String(),
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	resp, err := s.app.StreamKeeper.AllStreamsForReceiver(tCtx, req)
	s.Require().NoError(err)
	s.Require().Equal(expResp.Streams, resp.Streams)
}

func (s *KeeperTestSuite) TestQueryAllStreamsForReceiver_Success_No_Results() {
	req := &types.QueryAllStreamsForReceiverRequest{
		ReceiverAddr: s.addrs[9].String(),
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	expResp := &types.QueryAllStreamsForReceiverResponse{
		Streams: []*types.StreamResult{},
	}

	resp, err := s.app.StreamKeeper.AllStreamsForReceiver(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(expResp.Streams, resp.Streams)
}

func (s *KeeperTestSuite) TestQueryAllStreamsForReceiver_Fail_Bad_Or_Empty_Receiver() {
	req1 := &types.QueryAllStreamsForReceiverRequest{
		ReceiverAddr: "rubbish",
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	resp1, err1 := s.app.StreamKeeper.AllStreamsForReceiver(s.ctx, req1)
	s.Require().Error(err1)
	s.Require().ErrorContains(err1, "decoding bech32 failed")
	s.Require().Nil(resp1)

	req2 := &types.QueryAllStreamsForReceiverRequest{
		ReceiverAddr: "",
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	resp2, err2 := s.app.StreamKeeper.AllStreamsForReceiver(s.ctx, req2)
	s.Require().Error(err2)
	s.Require().ErrorContains(err2, "empty address string is not allowed")
	s.Require().Nil(resp2)
}

// ToDo
func (s *KeeperTestSuite) TestQueryAllStreamsForSender_Success() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(nowTime)

	deposit := sdk.NewInt64Coin("stake", 1000)
	qSender := s.addrs[0]
	// first three for results, sender = s.addrs[0]
	_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[1], qSender, deposit, 1)
	_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[2], qSender, deposit, 1)
	_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[3], qSender, deposit, 1)

	_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[1], s.addrs[3], deposit, 1)
	_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[2], s.addrs[4], deposit, 1)
	_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[3], s.addrs[5], deposit, 1)
	_, _ = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[4], s.addrs[6], deposit, 1)

	_, _ = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[1], qSender, deposit)
	_, _ = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[2], qSender, deposit)
	_, _ = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[3], qSender, deposit)
	_, _ = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[1], s.addrs[3], deposit)
	_, _ = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[2], s.addrs[4], deposit)
	_, _ = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[3], s.addrs[5], deposit)
	_, _ = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[4], s.addrs[6], deposit)

	expResp := &types.QueryAllStreamsForSenderResponse{
		Streams: []*types.StreamResult{
			{
				Receiver: s.addrs[1].String(),
				Sender:   qSender.String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
			{
				Receiver: s.addrs[2].String(),
				Sender:   qSender.String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
			{
				Receiver: s.addrs[3].String(),
				Sender:   qSender.String(),
				Stream: &types.Stream{
					Deposit:         sdk.NewInt64Coin("stake", 1000),
					FlowRate:        1,
					LastOutflowTime: nowTime,
					DepositZeroTime: time.Unix(nowTime.Unix()+1000, 0).UTC(),
					Cancellable:     true,
				},
			},
		},
	}

	req := &types.QueryAllStreamsForSenderRequest{
		SenderAddr: s.addrs[0].String(),
	}

	resp, err := s.app.StreamKeeper.AllStreamsForSender(tCtx, req)
	s.Require().NoError(err)
	s.Require().Equal(expResp.Streams, resp.Streams)
}

func (s *KeeperTestSuite) TestQueryAllStreamsForSender_Success_No_Results() {
	req := &types.QueryAllStreamsForSenderRequest{
		SenderAddr: s.addrs[9].String(),
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	expResp := &types.QueryAllStreamsForSenderResponse{
		Streams: []*types.StreamResult{},
	}

	resp, err := s.app.StreamKeeper.AllStreamsForSender(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(expResp.Streams, resp.Streams)
}

func (s *KeeperTestSuite) TestQueryAllStreamsForSender_Fail_Bad_Or_Empty_Receiver() {
	req1 := &types.QueryAllStreamsForSenderRequest{
		SenderAddr: "rubbish",
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	resp1, err1 := s.app.StreamKeeper.AllStreamsForSender(s.ctx, req1)
	s.Require().Error(err1)
	s.Require().ErrorContains(err1, "decoding bech32 failed")
	s.Require().Nil(resp1)

	req2 := &types.QueryAllStreamsForSenderRequest{
		SenderAddr: "",
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	resp2, err2 := s.app.StreamKeeper.AllStreamsForSender(s.ctx, req2)
	s.Require().Error(err2)
	s.Require().ErrorContains(err2, "empty address string is not allowed")
	s.Require().Nil(resp2)
}
