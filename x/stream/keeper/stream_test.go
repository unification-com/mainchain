package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/stream/types"
)

func (s *KeeperTestSuite) TestIsStream() {
	ok := s.app.StreamKeeper.IsStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().False(ok)

	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
		FlowRate:        100,
		CreateTime:      nowTime,
		LastUpdatedTime: nowTime,
		LastOutflowTime: nowTime,
		DepositZeroTime: time.Unix(0, 0).UTC(),
		TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		Cancellable:     true,
	}

	err := s.app.StreamKeeper.SetStream(s.ctx, s.addrs[1], s.addrs[0], expStream)
	s.Require().NoError(err)

	ok = s.app.StreamKeeper.IsStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().True(ok)
}

func (s *KeeperTestSuite) TestSetGetStream() {

	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().False(ok)
	s.Require().Equal(types.Stream{}, stream)

	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
		FlowRate:        100,
		CreateTime:      nowTime,
		LastUpdatedTime: nowTime,
		LastOutflowTime: nowTime,
		DepositZeroTime: time.Unix(0, 0).UTC(),
		TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		Cancellable:     true,
	}

	err := s.app.StreamKeeper.SetStream(s.ctx, s.addrs[1], s.addrs[0], expStream)
	s.Require().NoError(err)

	stream, ok = s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])

	s.Require().True(ok)
	s.Require().Equal(expStream.Deposit, stream.Deposit)
	s.Require().Equal(expStream.FlowRate, stream.FlowRate)
	s.Require().Equal(expStream.CreateTime, stream.CreateTime)
	s.Require().Equal(expStream.LastUpdatedTime, stream.LastUpdatedTime)
	s.Require().Equal(expStream.LastOutflowTime, stream.LastOutflowTime)
	s.Require().Equal(expStream.DepositZeroTime, stream.DepositZeroTime)
	s.Require().Equal(expStream.TotalStreamed, stream.TotalStreamed)
	s.Require().Equal(expStream.Cancellable, stream.Cancellable)
}

func (s *KeeperTestSuite) TestCreateNewStream_BasicSuccess() {
	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)), // set to 0 when created. AddDeposit handles setting deposit in stream
		FlowRate:        100,
		CreateTime:      nowTime,
		LastUpdatedTime: nowTime,
		LastOutflowTime: nowTime,
		DepositZeroTime: time.Unix(0, 0).UTC(),
		TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		Cancellable:     true,
	}

	stream, err := s.app.StreamKeeper.CreateNewStream(s.ctx, s.addrs[1], s.addrs[0], sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)), 100)

	s.Require().NoError(err)

	s.Require().Equal(expStream.Deposit, stream.Deposit)
	s.Require().Equal(expStream.FlowRate, stream.FlowRate)
	s.Require().Equal(expStream.CreateTime, stream.CreateTime)
	s.Require().Equal(expStream.LastUpdatedTime, stream.LastUpdatedTime)
	s.Require().Equal(expStream.LastOutflowTime, stream.LastOutflowTime)
	s.Require().Equal(expStream.DepositZeroTime, stream.DepositZeroTime)
	s.Require().Equal(expStream.TotalStreamed, stream.TotalStreamed)
	s.Require().Equal(expStream.Cancellable, stream.Cancellable)

	events := s.ctx.EventManager().Events()

	hasCreateStreamEvent := false
	for _, ev := range events {
		if ev.Type == types.EventTypeCreateStreamAction {
			hasCreateStreamEvent = true

			attrSender, ok := ev.GetAttribute(types.AttributeKeyStreamSender)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamSender, attrSender.Key)
			s.Require().Equal(s.addrs[0].String(), attrSender.Value)

			attrReceiver, ok := ev.GetAttribute(types.AttributeKeyStreamReceiver)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
			s.Require().Equal(s.addrs[1].String(), attrReceiver.Value)

			attrFlowRate, ok := ev.GetAttribute(types.AttributeKeyFlowRate)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyFlowRate, attrFlowRate.Key)
			s.Require().Equal("100", attrFlowRate.Value)
		}
	}

	// should emit create_stream event
	s.Require().True(hasCreateStreamEvent)

	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])

	s.Require().True(ok)
	s.Require().Equal(expStream.Deposit, stream.Deposit)
	s.Require().Equal(expStream.FlowRate, stream.FlowRate)
	s.Require().Equal(expStream.CreateTime, stream.CreateTime)
	s.Require().Equal(expStream.LastUpdatedTime, stream.LastUpdatedTime)
	s.Require().Equal(expStream.LastOutflowTime, stream.LastOutflowTime)
	s.Require().Equal(expStream.DepositZeroTime, stream.DepositZeroTime)
	s.Require().Equal(expStream.TotalStreamed, stream.TotalStreamed)
	s.Require().Equal(expStream.Cancellable, stream.Cancellable)
}

func (s *KeeperTestSuite) TestCreateNewStream_Fail_Stream_Exists() {
	_, err := s.app.StreamKeeper.CreateNewStream(s.ctx, s.addrs[1], s.addrs[0], sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)), 100)
	s.Require().NoError(err)

	_, err = s.app.StreamKeeper.CreateNewStream(s.ctx, s.addrs[1], s.addrs[0], sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)), 100)
	s.Require().ErrorContains(err, "stream exists")
}

func (s *KeeperTestSuite) TestAddDeposit_Basic_Success() {
	// set stream
	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)), //default to zero when creating a new stream
		FlowRate:        1,
		CreateTime:      nowTime,
		LastUpdatedTime: nowTime,
		LastOutflowTime: nowTime,
		DepositZeroTime: nowTime,
		TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		Cancellable:     true,
	}

	// set stream
	err := s.app.StreamKeeper.SetStream(s.ctx, s.addrs[1], s.addrs[0], expStream)
	s.Require().NoError(err)

	// Add Deposit to stream
	ok, err := s.app.StreamKeeper.AddDeposit(s.ctx, s.addrs[1], s.addrs[0], sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)))
	s.Require().True(ok)
	s.Require().NoError(err)

	// check events ar emitted
	events := s.ctx.EventManager().Events()

	hasEvent := false
	for _, ev := range events {
		if ev.Type == types.EventTypeDepositToStream {
			hasEvent = true

			attrSender, ok := ev.GetAttribute(types.AttributeKeyStreamSender)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamSender, attrSender.Key)
			s.Require().Equal(s.addrs[0].String(), attrSender.Value)

			attrReceiver, ok := ev.GetAttribute(types.AttributeKeyStreamReceiver)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
			s.Require().Equal(s.addrs[1].String(), attrReceiver.Value)

			attrDepositAmount, evOk := ev.GetAttribute(types.AttributeKeyAmountDeposited)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyAmountDeposited, attrDepositAmount.Key)
			s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)).String(), attrDepositAmount.Value)

			attrDepositDuration, evOk := ev.GetAttribute(types.AttributeKeyDepositDuration)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyDepositDuration, attrDepositDuration.Key)
			s.Require().Equal("1000", attrDepositDuration.Value)

			attrDepositZeroTime, evOk := ev.GetAttribute(types.AttributeKeyDepositZeroTime)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyDepositZeroTime, attrDepositZeroTime.Key)
			s.Require().Equal(nowTime.Add(time.Second*1000).String(), attrDepositZeroTime.Value)

			attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
			s.Require().Equal("1000stake", attrRemainingDeposit.Value)
		}
	}

	// should emit stream_deposit event
	s.Require().True(hasEvent)

	// get stream from keeper
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	// should now be 1000stake
	s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)), stream.Deposit)
	// Deposit of 1000, flow rate of 100/s, should have deposit zero time of now + 10s
	s.Require().Equal(nowTime.Add(time.Second*1000), stream.DepositZeroTime)
	s.Require().Equal(nowTime, stream.LastUpdatedTime)
}

func (s *KeeperTestSuite) TestAddDeposit_Success_TopUpExistingNotExpired() {
	tCtx := s.ctx

	blockTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(1)
	nowTime := tCtx.BlockTime()

	testCases := []struct {
		name               string
		sender             sdk.AccAddress
		receiver           sdk.AccAddress
		stream             types.Stream
		deposit            sdk.Coin
		expDepositZeroTime time.Time
		expDeposit         sdk.Coin
		expDiff            int64
	}{
		{
			name:     "1",
			sender:   s.addrs[0],
			receiver: s.addrs[1],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(750)), // initial deposit was 1000, claim 250s ago
				FlowRate:        1,
				CreateTime:      time.Unix(nowTime.Unix()-500, 0).UTC(), // created 500s ago
				LastUpdatedTime: time.Unix(nowTime.Unix()-500, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-250, 0).UTC(), // last claim was 250s ago
				DepositZeroTime: nowTime.Add(time.Second * 500),         // have 500s left (created 500s ago, deposit 1000, flow rate 1/s)
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(250)),
				Cancellable:     true,
			},
			deposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expDepositZeroTime: nowTime.Add(time.Second * 1500),
			expDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1750)),
			expDiff:            1750,
		},
		{
			name:     "2",
			sender:   s.addrs[2],
			receiver: s.addrs[3],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(45343254343)),
				FlowRate:        142723,
				CreateTime:      time.Unix(nowTime.Unix()-637701, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix()-637701, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-317701, 0).UTC(),
				DepositZeroTime: nowTime.Add(time.Second * 227701),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				Cancellable:     true,
			},
			deposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(8359902543123)),
			expDepositZeroTime: nowTime.Add(time.Second * 58802020),
			expDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(8405245797466)),
			expDiff:            59119721,
		},
		{
			name: "3", // 4584/month stream, created 3 weeks ago. Last claim 1 week ago (approx half claimed).
			// 1 week until deposit zero. Top up with 1 month's worth 4584
			sender:   s.addrs[4],
			receiver: s.addrs[5],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(2474104396800)), // approx 2 weeks worth left
				FlowRate:        1744292,                                                   // approx 4584/month
				CreateTime:      time.Unix(nowTime.Unix()-1814400, 0).UTC(),                // 3 weeks ago
				LastUpdatedTime: time.Unix(nowTime.Unix()-604800, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-604800, 0).UTC(), // approx 1 week ago - 2 weeks claimed
				DepositZeroTime: nowTime.Add(time.Second * 604800),         // 1 week in future
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(2109895603200)),
				Cancellable:     true,
			},
			deposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(4584000000000)), // 4584
			expDepositZeroTime: nowTime.Add(time.Second * 3232800),                        // in approx 5 weeks. 1 week deposit remaining, plus 1 month more
			expDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(7058104396800)),
			expDiff:            3837600, // diff between last claim and deposit zero. 6 weeks
		},
	}

	for _, tc := range testCases {
		// deposit zero time is in the future, so just use SetStream instead of create & add deposit combo
		err := s.app.StreamKeeper.SetStream(tCtx, tc.receiver, tc.sender, tc.stream)
		s.Require().NoError(err, "SetStream NoError test name %s", tc.name)
		ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.deposit)
		s.Require().True(ok, "AddDeposit True test name %s", tc.name)
		s.Require().NoError(err, "AddDeposit NoError test name %s", tc.name)

		// events should NOT contain claim_stream
		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				hasEvent = true
			}
		}
		s.Require().False(hasEvent)

		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream True test name %s", tc.name)
		s.Require().Equal(tc.expDeposit, stream.Deposit, "GetStream Deposit Equal test name %s", tc.name)
		s.Require().Equal(tc.expDepositZeroTime, stream.DepositZeroTime, "GetStream DepositZeroTime Equal test name %s", tc.name)

		duration := stream.DepositZeroTime.Unix() - stream.LastOutflowTime.Unix()
		s.Require().Equal(tc.expDiff, duration, "duration test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestAddDeposit_Success_TopUpExistingExpired() {
	tCtx := s.ctx

	blockTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(1)
	nowTime := tCtx.BlockTime()

	testCases := []struct {
		name               string
		sender             sdk.AccAddress
		receiver           sdk.AccAddress
		stream             types.Stream
		initialDeposit     sdk.Coin
		newDeposit         sdk.Coin
		expDepositZeroTime time.Time
		expDeposit         sdk.Coin
		expTotalStreamed   sdk.Coin
	}{
		{
			name:     "1",
			sender:   s.addrs[0],
			receiver: s.addrs[1],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				FlowRate:        1,
				CreateTime:      time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				DepositZeroTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				Cancellable:     true,
			},
			initialDeposit:     sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			newDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expDepositZeroTime: nowTime.Add(time.Second * 1000),
			expDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expTotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:     "2",
			sender:   s.addrs[2],
			receiver: s.addrs[3],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				FlowRate:        142723,
				CreateTime:      time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				DepositZeroTime: nowTime,
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				Cancellable:     true,
			},
			initialDeposit:     sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			newDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(8359902543123)),
			expDepositZeroTime: nowTime.Add(time.Second * 58574319),
			expDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(8359902543123)),
			expTotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:     "3",
			sender:   s.addrs[4],
			receiver: s.addrs[5],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				FlowRate:        142723,
				CreateTime:      time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				DepositZeroTime: nowTime,
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				Cancellable:     true,
			},
			initialDeposit:     sdk.NewCoin("stake", sdk.NewIntFromUint64(234232455325)),
			newDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(8359902543123)),
			expDepositZeroTime: nowTime.Add(time.Second * 58574319),
			expDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(8359902543123)),
			expTotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(234232455325)),
		},
		{
			name:     "4",
			sender:   s.addrs[6],
			receiver: s.addrs[7],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				FlowRate:        142723,
				CreateTime:      time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				DepositZeroTime: time.Unix(nowTime.Unix()-1000, 0).UTC(),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				Cancellable:     true,
			},
			initialDeposit:     sdk.NewCoin("stake", sdk.NewIntFromUint64(234232455325)),
			newDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(8359902543123)),
			expDepositZeroTime: nowTime.Add(time.Second * 58574319),
			expDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(8359902543123)),
			expTotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(234232455325)),
		},
	}

	for _, tc := range testCases {
		// create
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.initialDeposit, tc.stream.FlowRate)
		s.Require().NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		if tc.initialDeposit.Amount.GT(sdk.NewIntFromUint64(0)) {
			ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.initialDeposit)
			s.Require().True(ok)
			s.Require().NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)
		}

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok)
		s.Require().Equal(tc.initialDeposit, stream.Deposit)

		// set times etc.
		stream.CreateTime = tc.stream.CreateTime
		stream.LastUpdatedTime = tc.stream.LastUpdatedTime
		stream.LastOutflowTime = tc.stream.LastOutflowTime
		stream.DepositZeroTime = tc.stream.DepositZeroTime
		err = s.app.StreamKeeper.SetStream(tCtx, tc.receiver, tc.sender, stream)
		s.Require().NoError(err, "SetStream NoError test name %s", tc.name)

		// top up with new deposit
		ok, err = s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.newDeposit)
		s.Require().True(ok, "AddDeposit True test name %s", tc.name)
		s.Require().NoError(err, "AddDeposit NoError test name %s", tc.name)

		// check events do contain claim_stream if expTotalStreamed > 0
		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				hasEvent = true
			}
		}
		if tc.expTotalStreamed.IsPositive() {
			s.Require().True(hasEvent)
		} else {
			s.Require().False(hasEvent)
		}

		// final check
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream True test name %s", tc.name)
		s.Require().Equal(tc.expDeposit, stream.Deposit, "GetStream Deposit Equal test name %s", tc.name)
		s.Require().Equal(tc.expDepositZeroTime, stream.DepositZeroTime, "GetStream DepositZeroTime Equal test name %s", tc.name)
		s.Require().Equal(tc.expTotalStreamed, stream.TotalStreamed, "GetStream TotalStreamed Equal test name %s", tc.name)
		s.Require().Equal(nowTime, stream.LastUpdatedTime, "GetStream TotalStreamed Equal test name %s", tc.name)
		if tc.expTotalStreamed.IsPositive() {
			s.Require().Equal(nowTime, stream.LastOutflowTime, "GetStream TotalStreamed Equal test name %s", tc.name)
		}

	}
}

func (s *KeeperTestSuite) TestAddDeposit_Scenarios() {
	testCases := []struct {
		name                  string
		sender                sdk.AccAddress
		receiver              sdk.AccAddress
		flowRate              int64
		initialDeposit        sdk.Coin
		newDeposit            sdk.Coin
		createTimeOffset      int64    // seconds in past from "now"
		expInitialDepZeroTime int64    // seconds in future from create time
		expNewDeposit         sdk.Coin // after new deposit added
		expNewDepZeroTime     int64    // seconds in future from "now"
		expClaim              sdk.Coin
		expRemainDeposit      sdk.Coin // from claim event emission only
	}{
		{
			name:                  "simple 1 not expired",
			sender:                s.addrs[0],
			receiver:              s.addrs[1],
			flowRate:              1,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			newDeposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			createTimeOffset:      500,
			expInitialDepZeroTime: 1000,
			expNewDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(2000)),
			expNewDepZeroTime:     1500,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
		},
		{
			name:                  "simple 2 not expired",
			sender:                s.addrs[2],
			receiver:              s.addrs[3],
			flowRate:              4324532,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(10461907814400)), // 4 weeks worth
			newDeposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(10461907814400)), // another 4 weeks
			createTimeOffset:      1814400,                                                    // 3 weeks ago
			expInitialDepZeroTime: 2419200,                                                    // 4 weeks from creation date
			expNewDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(20923815628800)),
			expNewDepZeroTime:     3024000, // approx 5 weeks in the future from now
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(10461907814400)),
		},
		{
			name:                  "complex not expired",
			sender:                s.addrs[4],
			receiver:              s.addrs[5],
			flowRate:              54875,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(149772587375)),  // 54875 x 2729341
			newDeposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(1245790781440)), // 1245790781440 / 54875 = 22702337 seconds
			createTimeOffset:      1605494,                                                   // 1605494 seconds ago
			expInitialDepZeroTime: 2729341,                                                   // 2729341 seconds from creation date
			expNewDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1395563368815)), // 149772587375 + 1245790781440
			expNewDepZeroTime:     23826184,                                                  // (2729341-1605494) + 22702337
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(149772587375)),
		},
		{
			name:                  "simple expires now",
			sender:                s.addrs[6],
			receiver:              s.addrs[7],
			flowRate:              1,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			newDeposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			createTimeOffset:      1000,
			expInitialDepZeroTime: 1000,
			expNewDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expNewDepZeroTime:     1000,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:                  "simple expires in past",
			sender:                s.addrs[8],
			receiver:              s.addrs[9],
			flowRate:              1,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			newDeposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			createTimeOffset:      1500,
			expInitialDepZeroTime: 1000,
			expNewDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expNewDepZeroTime:     1000,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:                  "complex expires now",
			sender:                s.addrs[10],
			receiver:              s.addrs[11],
			flowRate:              87656,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(231352935168)),
			newDeposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(296752935417)),
			createTimeOffset:      2639328, // same as expInitialDepZeroTime
			expInitialDepZeroTime: 2639328, // 231352935168 / 87656
			expNewDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(296752935417)),
			expNewDepZeroTime:     3385426,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(231352935168)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:                  "complex expires in past",
			sender:                s.addrs[12],
			receiver:              s.addrs[13],
			flowRate:              782563,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(2535134750264)),
			newDeposit:            sdk.NewCoin("stake", sdk.NewIntFromUint64(3128529354197)),
			createTimeOffset:      3739341, // arbitrary - further in the past than expInitialDepZeroTime
			expInitialDepZeroTime: 3239528, // 2535134750264 / 782563
			expNewDeposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(3128529354197)),
			expNewDepZeroTime:     3997798, // 3128529354197 / 782563
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(2535134750264)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
	}

	for _, tc := range testCases {
		tCtx := s.ctx
		nowTime := time.Unix(time.Now().Unix(), 0).UTC()
		// set create time to past
		blockTimeCreate := time.Unix(nowTime.Unix()-tc.createTimeOffset, 0).UTC()
		tCtx = tCtx.WithBlockTime(blockTimeCreate).WithBlockHeight(1)

		// create
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.initialDeposit, tc.flowRate)

		s.Require().NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.initialDeposit)
		s.Require().True(ok, "initialDeposit ok NoError test name %s", tc.name)
		s.Require().NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		// should be in the future from the creation time
		expInitialDepZeroTime := time.Unix(blockTimeCreate.Unix()+tc.expInitialDepZeroTime, 0).UTC()
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Require().Equal(tc.initialDeposit, stream.Deposit, "tc.initialDeposit Equal stream.Deposit test name %s", tc.name)
		s.Require().Equal(expInitialDepZeroTime, stream.DepositZeroTime, "expInitialDepZeroTimeEqual stream.DepositZeroTime test name %s", tc.name)

		// set block time to now
		tCtx = tCtx.WithBlockTime(nowTime).WithBlockHeight(2)

		// add new deposit
		ok, err = s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.newDeposit)
		s.Require().True(ok, "newDeposit AddDeposit ok test name %s", tc.name)
		s.Require().NoError(err, "newDeposit AddDeposit NoError test name %s", tc.name)

		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
				attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
				// only for this stream
				if tc.sender.String() == attrSender.Value && tc.receiver.String() == attrReceiver.Value {
					hasEvent = true
				}
			}
		}

		if tc.expClaim.Amount.IsPositive() {
			s.Require().True(hasEvent)
			for _, ev := range events {
				if ev.Type == types.EventTypeClaimStreamAction {
					attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
					attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
					// only for this stream
					if tc.sender.String() != attrSender.Value || tc.receiver.String() != attrReceiver.Value {
						// skip events not for this stream
						continue
					}

					attrClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyClaimTotal)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimTotal, attrClaimTotal.Key)
					s.Require().Equal(tc.expClaim.String(), attrClaimTotal.Value, "AttributeKeyClaimTotal test name %s", tc.name)

					attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
					s.Require().Equal(tc.expRemainDeposit.String(), attrRemainingDeposit.Value, "AttributeKeyRemainingDeposit test name %s", tc.name)
				}
			}
		} else {
			s.Require().False(hasEvent)
		}

		// check results
		expDepZeroTime := time.Unix(nowTime.Unix()+tc.expNewDepZeroTime, 0).UTC()
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Require().Equal(tc.expNewDeposit, stream.Deposit, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
		s.Require().Equal(expDepZeroTime, stream.DepositZeroTime, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestAddDeposit_Fail_StreamNotExist() {
	ok, err := s.app.StreamKeeper.AddDeposit(s.ctx, s.addrs[1], s.addrs[0], sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)))
	s.Require().False(ok)
	s.Require().ErrorContains(err, "stream does not exist")

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().False(ok)
	s.Require().Equal(types.Stream{}, stream)
}

func (s *KeeperTestSuite) TestAddDeposit_Fail_InsufficientBalance() {
	// set stream
	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)), //default to zero when creating a new stream
		FlowRate:        100,
		CreateTime:      nowTime,
		LastUpdatedTime: time.Unix(0, 0).UTC(),
		LastOutflowTime: nowTime,
		DepositZeroTime: time.Unix(0, 0).UTC(),
		TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		Cancellable:     true,
	}

	// set stream
	err := s.app.StreamKeeper.SetStream(s.ctx, s.addrs[1], s.addrs[0], expStream)
	s.Require().NoError(err)

	// deposit more than sender's balance
	ok, err := s.app.StreamKeeper.AddDeposit(s.ctx, s.addrs[1], s.addrs[0], sdk.NewCoin("stake", sdk.NewIntFromUint64(1000000000000000001)))
	s.Require().False(ok)
	s.Require().ErrorContains(err, "insufficient funds")

}

func (s *KeeperTestSuite) TestSetNewFlowRate_Success() {
	// set stream
	tCtx := s.ctx

	blockTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(1)
	nowTime := tCtx.BlockTime()
	deposit := sdk.NewCoin("stake", sdk.NewIntFromUint64(2400))

	// create stream
	_, err := s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[1], s.addrs[0], deposit, 1)
	s.Require().NoError(err)

	// add deposit
	ok, err := s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[1], s.addrs[0], deposit)
	s.Require().True(ok)
	s.Require().NoError(err)

	// Set new flow rate
	err = s.app.StreamKeeper.SetNewFlowRate(tCtx, s.addrs[1], s.addrs[0], 24)
	s.Require().NoError(err)

	// check events ar emitted
	events := tCtx.EventManager().Events()

	hasEvent := false
	for _, ev := range events {
		if ev.Type == types.EventTypeUpdateFlowRate {
			hasEvent = true
			attrSender, ok := ev.GetAttribute(types.AttributeKeyStreamSender)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamSender, attrSender.Key)
			s.Require().Equal(s.addrs[0].String(), attrSender.Value)

			attrReceiver, ok := ev.GetAttribute(types.AttributeKeyStreamReceiver)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
			s.Require().Equal(s.addrs[1].String(), attrReceiver.Value)

			attrOldFlowRate, evOk := ev.GetAttribute(types.AttributeKeyOldFlowRate)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyOldFlowRate, attrOldFlowRate.Key)
			s.Require().Equal("1", attrOldFlowRate.Value)

			attrNewFlowRate, evOk := ev.GetAttribute(types.AttributeKeyNewFlowRate)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyNewFlowRate, attrNewFlowRate.Key)
			s.Require().Equal("24", attrNewFlowRate.Value)

			attrStreamDepositDuration, evOk := ev.GetAttribute(types.AttributeKeyDepositDuration)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyDepositDuration, attrStreamDepositDuration.Key)
			s.Require().Equal("100", attrStreamDepositDuration.Value)

			attrStreamDepositZeroTime, evOk := ev.GetAttribute(types.AttributeKeyDepositZeroTime)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyDepositZeroTime, attrStreamDepositZeroTime.Key)
			s.Require().Equal(nowTime.Add(time.Second*100).String(), attrStreamDepositZeroTime.Value)

			attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
			s.Require().Equal("2400stake", attrRemainingDeposit.Value)
		}
	}

	// should emit stream_deposit event
	s.Require().True(hasEvent)

	// get stream from keeper
	stream, ok := s.app.StreamKeeper.GetStream(tCtx, s.addrs[1], s.addrs[0])
	s.Require().True(ok)
	// should now be 1000stake
	s.Require().Equal(int64(24), stream.FlowRate)
	s.Require().Equal(nowTime.Add(time.Second*100), stream.DepositZeroTime)
	s.Require().Equal(nowTime, stream.LastUpdatedTime)
}

func (s *KeeperTestSuite) TestSetNewFlowRate_Success_ExistingNotExpired() {
	tCtx := s.ctx

	blockTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(1)
	nowTime := tCtx.BlockTime()

	testCases := []struct {
		name               string
		sender             sdk.AccAddress
		receiver           sdk.AccAddress
		stream             types.Stream
		newFlowRate        int64
		expDepositZeroTime time.Time
	}{
		{
			name:     "1",
			sender:   s.addrs[0],
			receiver: s.addrs[1],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
				FlowRate:        1,
				CreateTime:      time.Unix(nowTime.Unix()-500, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix(), 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix(), 0).UTC(),
				DepositZeroTime: nowTime.Add(time.Second * 500),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
				Cancellable:     true,
			},
			newFlowRate:        2,
			expDepositZeroTime: nowTime.Add(time.Second * 250),
		},
		{
			name:     "2", // effectively expired
			sender:   s.addrs[2],
			receiver: s.addrs[3],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				FlowRate:        1,
				CreateTime:      time.Unix(nowTime.Unix()-5000, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix()-500, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix()-500, 0).UTC(),
				DepositZeroTime: time.Unix(nowTime.Unix()-500, 0).UTC(),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
				Cancellable:     true,
			},
			newFlowRate:        200,
			expDepositZeroTime: nowTime,
		},
		{
			name:     "3",
			sender:   s.addrs[4],
			receiver: s.addrs[5],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				FlowRate:        1,
				CreateTime:      time.Unix(nowTime.Unix()-500, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix(), 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix(), 0).UTC(),
				DepositZeroTime: time.Unix(nowTime.Unix()-500, 0).UTC(),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
				Cancellable:     true,
			},
			newFlowRate:        99,
			expDepositZeroTime: nowTime,
		},
		{
			name:     "4",
			sender:   s.addrs[6],
			receiver: s.addrs[7],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(45343254343)),
				FlowRate:        142723,
				CreateTime:      time.Unix(nowTime.Unix()-637701, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix(), 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix(), 0).UTC(),
				DepositZeroTime: nowTime.Add(time.Second * 227701),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
				Cancellable:     true,
			},
			newFlowRate:        150000,
			expDepositZeroTime: nowTime.Add(time.Second * 302288),
		},
		{
			name:     "5",
			sender:   s.addrs[8],
			receiver: s.addrs[9],
			stream: types.Stream{
				Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(2474104396800)),
				FlowRate:        1744292,
				CreateTime:      time.Unix(nowTime.Unix()-1814400, 0).UTC(),
				LastUpdatedTime: time.Unix(nowTime.Unix()-604800, 0).UTC(),
				LastOutflowTime: time.Unix(nowTime.Unix(), 0).UTC(),
				DepositZeroTime: nowTime.Add(time.Second * 604800),
				TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(2109895603200)),
				Cancellable:     true,
			},
			newFlowRate:        1444444,
			expDepositZeroTime: nowTime.Add(time.Second * 1712842),
		},
	}

	for _, tc := range testCases {
		err := s.app.StreamKeeper.SetStream(tCtx, tc.receiver, tc.sender, tc.stream)
		s.Require().NoError(err, "SetStream NoError test name %s", tc.name)
		err = s.app.StreamKeeper.SetNewFlowRate(tCtx, tc.receiver, tc.sender, tc.newFlowRate)
		s.Require().NoError(err, "AddDeposit NoError test name %s", tc.name)

		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream True test name %s", tc.name)
		s.Require().Equal(tc.expDepositZeroTime, stream.DepositZeroTime, "GetStream DepositZeroTime Equal test name %s", tc.name)
		s.Require().Equal(tc.newFlowRate, stream.FlowRate, "GetStream FlowRate Equal test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestSetNewFlowRate_Scenarios() {
	testCases := []struct {
		name                  string
		sender                sdk.AccAddress
		receiver              sdk.AccAddress
		startFlowRate         int64
		initialDeposit        sdk.Coin
		newFlowRate           int64
		createTimeOffset      int64 // seconds in past from "now"
		expInitialDepZeroTime int64 // seconds in future from create time
		expNewDepZeroTime     int64 // seconds in future from "now"
		expClaim              sdk.Coin
		expRemainDeposit      sdk.Coin // from claim event emission only
		expDuration           int64    // from update_flow event
	}{
		{
			name:                  "not expired, has deposit, increase flow rate",
			sender:                s.addrs[0],
			receiver:              s.addrs[1],
			startFlowRate:         1,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			newFlowRate:           2,
			createTimeOffset:      500,
			expInitialDepZeroTime: 1000,
			expNewDepZeroTime:     250,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expDuration:           250,
		},
		{
			name:                  "expired, has deposit, increase flow rate",
			sender:                s.addrs[2],
			receiver:              s.addrs[3],
			startFlowRate:         1,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			newFlowRate:           2,
			createTimeOffset:      1500,
			expInitialDepZeroTime: 1000,
			expNewDepZeroTime:     0,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expDuration:           0,
		},
		{
			name:                  "not expired, has deposit, decrease flow rate",
			sender:                s.addrs[4],
			receiver:              s.addrs[5],
			startFlowRate:         2,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			newFlowRate:           1,
			createTimeOffset:      100,
			expInitialDepZeroTime: 500,
			expNewDepZeroTime:     800,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(200)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(800)),
			expDuration:           800,
		},
		{
			name:                  "expired, has deposit, decrease flow rate",
			sender:                s.addrs[6],
			receiver:              s.addrs[7],
			startFlowRate:         2,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			newFlowRate:           1,
			createTimeOffset:      1000,
			expInitialDepZeroTime: 500,
			expNewDepZeroTime:     0,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expDuration:           0,
		},
		{
			name:                  "expired, no deposit, increase flow rate",
			sender:                s.addrs[8],
			receiver:              s.addrs[9],
			startFlowRate:         1,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			newFlowRate:           2,
			createTimeOffset:      1500,
			expInitialDepZeroTime: 0,
			expNewDepZeroTime:     0,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expDuration:           0,
		},
		{
			name:                  "expired, no deposit, decrease flow rate",
			sender:                s.addrs[10],
			receiver:              s.addrs[11],
			startFlowRate:         2,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			newFlowRate:           3,
			createTimeOffset:      1500,
			expInitialDepZeroTime: 0,
			expNewDepZeroTime:     0,
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expDuration:           0,
		},
		{
			name:                  "complex not expired, has deposit, increase flow rate",
			sender:                s.addrs[12],
			receiver:              s.addrs[13],
			startFlowRate:         54875,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(149772587375)),
			newFlowRate:           69875,
			createTimeOffset:      1605494,
			expInitialDepZeroTime: 2729341,                                                 // 149772587375 / 54875
			expNewDepZeroTime:     882591,                                                  // 61671104125 / 69875
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(88101483250)), // 54875 * 1605494
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(61671104125)), // 149772587375 - 88101483250
			expDuration:           882591,
		},
		{
			name:                  "complex not expired, has deposit, decrease flow rate",
			sender:                s.addrs[14],
			receiver:              s.addrs[15],
			startFlowRate:         69416,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(349274587314)),
			newFlowRate:           51249,
			createTimeOffset:      2545494,
			expInitialDepZeroTime: 5031615,                                                  // 349274587314 / 69416
			expNewDepZeroTime:     3367413,                                                  // 172576575810 / 51249
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(176698011504)), // 69416 * 2545494
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(172576575810)), // 349274587314 - 176698011504
			expDuration:           882591,
		},
		{
			name:                  "complex expired, has deposit, increase flow rate",
			sender:                s.addrs[16],
			receiver:              s.addrs[17],
			startFlowRate:         54875,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(149772587375)),
			newFlowRate:           69875,
			createTimeOffset:      2729341,
			expInitialDepZeroTime: 2729341,                                                  // 149772587375 / 54875
			expNewDepZeroTime:     0,                                                        // 61671104125 / 69875
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(149772587375)), // 54875 * 1605494
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),            // 149772587375 - 88101483250
			expDuration:           0,
		},
		{
			name:                  "complex expired, has deposit, decrease flow rate",
			sender:                s.addrs[18],
			receiver:              s.addrs[19],
			startFlowRate:         69875,
			initialDeposit:        sdk.NewCoin("stake", sdk.NewIntFromUint64(249761587512)),
			newFlowRate:           32563,
			createTimeOffset:      3574405,
			expInitialDepZeroTime: 3574405,                                                  // 149772587375 / 54875
			expNewDepZeroTime:     0,                                                        // 61671104125 / 69875
			expClaim:              sdk.NewCoin("stake", sdk.NewIntFromUint64(249761587512)), // 54875 * 1605494
			expRemainDeposit:      sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),            // 149772587375 - 88101483250
			expDuration:           0,
		},
	}

	for _, tc := range testCases {
		tCtx := s.ctx
		nowTime := time.Unix(time.Now().Unix(), 0).UTC()
		// set create time to past
		blockTimeCreate := time.Unix(nowTime.Unix()-tc.createTimeOffset, 0).UTC()
		tCtx = tCtx.WithBlockTime(blockTimeCreate).WithBlockHeight(1)

		// create
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.initialDeposit, tc.startFlowRate)

		s.Require().NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		if tc.initialDeposit.IsPositive() {
			ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.initialDeposit)
			s.Require().True(ok, "initialDeposit ok NoError test name %s", tc.name)
			s.Require().NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)
		}

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		// should be in the future from the creation time
		expInitialDepZeroTime := time.Unix(blockTimeCreate.Unix()+tc.expInitialDepZeroTime, 0).UTC()
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Require().Equal(tc.initialDeposit, stream.Deposit, "tc.initialDeposit Equal stream.Deposit test name %s", tc.name)
		if tc.initialDeposit.IsPositive() {
			s.Require().Equal(expInitialDepZeroTime, stream.DepositZeroTime, "expInitialDepZeroTimeEqual stream.DepositZeroTime test name %s", tc.name)
		}
		// set block time to now
		tCtx = tCtx.WithBlockTime(nowTime).WithBlockHeight(2)

		// set new flow rate
		err = s.app.StreamKeeper.SetNewFlowRate(tCtx, tc.receiver, tc.sender, tc.newFlowRate)
		s.Require().NoError(err, "newFlowRate SetNewFlowRate NoError test name %s", tc.name)

		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
				attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
				// only for this stream
				if tc.sender.String() == attrSender.Value && tc.receiver.String() == attrReceiver.Value {
					hasEvent = true
				}
			}
		}

		if tc.expClaim.Amount.IsPositive() {
			s.Require().True(hasEvent)
			for _, ev := range events {
				if ev.Type == types.EventTypeClaimStreamAction {
					attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
					attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
					// only for this stream
					if tc.sender.String() != attrSender.Value || tc.receiver.String() != attrReceiver.Value {
						// skip events not for this stream
						continue
					}

					attrClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyClaimTotal)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimTotal, attrClaimTotal.Key)
					s.Require().Equal(tc.expClaim.String(), attrClaimTotal.Value, "AttributeKeyClaimTotal test name %s", tc.name)

					attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
					s.Require().Equal(tc.expRemainDeposit.String(), attrRemainingDeposit.Value, "AttributeKeyRemainingDeposit test name %s", tc.name)
				}
			}
		} else {
			s.Require().False(hasEvent)
		}

		// check results
		expDepZeroTime := time.Unix(nowTime.Unix()+tc.expNewDepZeroTime, 0).UTC()
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Require().Equal(tc.newFlowRate, stream.FlowRate, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
		s.Require().Equal(expDepZeroTime, stream.DepositZeroTime, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestSetNewFlowRate_Fail() {
	err := s.app.StreamKeeper.SetNewFlowRate(s.ctx, s.addrs[1], s.addrs[0], 24)
	s.Require().ErrorContains(err, "stream does not exist")

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().False(ok)
	s.Require().Equal(types.Stream{}, stream)
}

func (s *KeeperTestSuite) TestClaimFromStream_Success() {
	// set stream
	tCtx := s.ctx

	blockTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(1)
	nowTime := tCtx.BlockTime()

	// set validator fee
	valFee, err := sdk.NewDecFromStr("0.01")
	s.app.StreamKeeper.SetParams(tCtx, types.Params{ValidatorFee: valFee})

	deposit := sdk.NewCoin("stake", sdk.NewIntFromUint64(1000))

	// create stream
	_, err = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[1], s.addrs[0], deposit, 1)
	s.Require().NoError(err)
	// add deposit
	ok, err := s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[1], s.addrs[0], deposit)
	s.Require().True(ok)
	s.Require().NoError(err)

	// time travel
	future := time.Unix(nowTime.Unix()+500, 0).UTC()
	tCtx = tCtx.WithBlockTime(future).WithBlockHeight(2)

	// claim
	amntClaimed, valFeeSent, totalClaim, remainingDeposit, err := s.app.StreamKeeper.ClaimFromStream(tCtx, s.addrs[1], s.addrs[0])
	s.Require().NoError(err)
	s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(495)), amntClaimed, "amntClaimed")
	s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(5)), valFeeSent, "valFeeSent")
	s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(500)), totalClaim, "totalClaim")
	s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(500)), remainingDeposit, "remainingDeposit")

	// check event emission
	events := tCtx.EventManager().Events()

	hasEvent := false
	for _, ev := range events {
		if ev.Type == types.EventTypeClaimStreamAction {
			hasEvent = true

			attrSender, ok := ev.GetAttribute(types.AttributeKeyStreamSender)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamSender, attrSender.Key)
			s.Require().Equal(s.addrs[0].String(), attrSender.Value)

			attrReceiver, ok := ev.GetAttribute(types.AttributeKeyStreamReceiver)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
			s.Require().Equal(s.addrs[1].String(), attrReceiver.Value)

			attrStreamClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyClaimTotal)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyClaimTotal, attrStreamClaimTotal.Key)
			s.Require().Equal("500stake", attrStreamClaimTotal.Value)

			attrStreamClaimAmountReceived, evOk := ev.GetAttribute(types.AttributeKeyClaimAmountReceived)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyClaimAmountReceived, attrStreamClaimAmountReceived.Key)
			s.Require().Equal("495stake", attrStreamClaimAmountReceived.Value)

			attrStreamClaimValidatorFee, evOk := ev.GetAttribute(types.AttributeKeyClaimValidatorFee)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyClaimValidatorFee, attrStreamClaimValidatorFee.Key)
			s.Require().Equal("5stake", attrStreamClaimValidatorFee.Value)

			attrStreamRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyRemainingDeposit, attrStreamRemainingDeposit.Key)
			s.Require().Equal("500stake", attrStreamRemainingDeposit.Value)
		}
	}

	// should emit stream_deposit event
	s.Require().True(hasEvent)

	// check stream in keeper
	stream, ok := s.app.StreamKeeper.GetStream(tCtx, s.addrs[1], s.addrs[0])
	s.Require().True(ok)
	s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(500)), stream.Deposit, "stream.Deposit")
}

func (s *KeeperTestSuite) TestClaimFromStream_Scenarios() {
	testCases := []struct {
		name              string
		sender            sdk.AccAddress
		receiver          sdk.AccAddress
		flowRate          int64
		deposit           sdk.Coin
		valFee            string
		createTimeOffset  int64 // seconds in past from "now"
		expTotalClaim     sdk.Coin
		expReceiverAmount sdk.Coin
		expRemainDeposit  sdk.Coin
		expValFee         sdk.Coin
	}{
		{
			name:              "simple, not expired, no val fee",
			sender:            s.addrs[0],
			receiver:          s.addrs[1],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.0",
			createTimeOffset:  500,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:              "simple, expired, no val fee",
			sender:            s.addrs[2],
			receiver:          s.addrs[3],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.0",
			createTimeOffset:  9999,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:              "simple, not expired, 5% val fee",
			sender:            s.addrs[4],
			receiver:          s.addrs[5],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.05",
			createTimeOffset:  500,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(475)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(25)),
		},
		{
			name:              "simple, expired, 5% val fee",
			sender:            s.addrs[6],
			receiver:          s.addrs[7],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.05",
			createTimeOffset:  9999,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(950)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(50)),
		},
		{
			name:              "not expired, no val fee",
			sender:            s.addrs[8],
			receiver:          s.addrs[9],
			flowRate:          54769,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(198770009867)),
			valFee:            "0.0",
			createTimeOffset:  2628000,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(143932932000)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(143932932000)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(54837077867)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:              "not expired, 3% val fee",
			sender:            s.addrs[10],
			receiver:          s.addrs[11],
			flowRate:          69249,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(275585190123)),
			valFee:            "0.03",
			createTimeOffset:  1739523,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(120460228227)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(116846421381)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(155124961896)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(3613806846)),
		},
		{
			name:              "expired, 9% val fee",
			sender:            s.addrs[12],
			receiver:          s.addrs[13],
			flowRate:          73269,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(313563770856)),
			valFee:            "0.09",
			createTimeOffset:  5000000,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(313563770856)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(285343031479)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(28220739377)),
		},
	}

	for _, tc := range testCases {
		tCtx := s.ctx
		nowTime := time.Unix(time.Now().Unix(), 0).UTC()
		// set create time to past
		blockTimeCreate := time.Unix(nowTime.Unix()-tc.createTimeOffset, 0).UTC()
		tCtx = tCtx.WithBlockTime(blockTimeCreate).WithBlockHeight(1)

		// set params
		newParams := types.Params{
			ValidatorFee: sdk.MustNewDecFromStr(tc.valFee),
		}

		err := s.app.StreamKeeper.SetParams(tCtx, newParams)
		s.Require().NoError(err, "SetParams NoError test name %s", tc.name)

		// create
		_, err = s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.deposit, tc.flowRate)

		s.Require().NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		if tc.deposit.IsPositive() {
			ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.deposit)
			s.Require().True(ok, "initialDeposit ok NoError test name %s", tc.name)
			s.Require().NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)
		}

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)

		// set block time to now
		tCtx = tCtx.WithBlockTime(nowTime).WithBlockHeight(2)

		// claim
		receiverAmount, valFee, claimTotal, remainingDeposit, err := s.app.StreamKeeper.ClaimFromStream(tCtx, tc.receiver, tc.sender)
		s.Require().NoError(err, "ClaimFromStream NoError test name %s", tc.name)
		s.Require().Equal(tc.expReceiverAmount, receiverAmount, "ClaimFromStream receiverAmount test name %s", tc.name)
		s.Require().Equal(tc.expValFee, valFee, "ClaimFromStream valFee test name %s", tc.name)
		s.Require().Equal(tc.expTotalClaim, claimTotal, "ClaimFromStream totalClaim test name %s", tc.name)
		s.Require().Equal(tc.expRemainDeposit, remainingDeposit, "ClaimFromStream remainingDeposit test name %s", tc.name)

		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
				attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
				// only for this stream
				if tc.sender.String() == attrSender.Value && tc.receiver.String() == attrReceiver.Value {
					hasEvent = true
				}
			}
		}

		if tc.expReceiverAmount.Amount.IsPositive() {
			s.Require().True(hasEvent)
			for _, ev := range events {
				if ev.Type == types.EventTypeClaimStreamAction {
					attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
					attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
					// only for this stream
					if tc.sender.String() != attrSender.Value || tc.receiver.String() != attrReceiver.Value {
						// skip events not for this stream
						continue
					}

					attrClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyClaimTotal)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimTotal, attrClaimTotal.Key)
					s.Require().Equal(tc.expTotalClaim.String(), attrClaimTotal.Value, "AttributeKeyClaimTotal test name %s", tc.name)

					attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
					s.Require().Equal(tc.expRemainDeposit.String(), attrRemainingDeposit.Value, "AttributeKeyRemainingDeposit test name %s", tc.name)

					attrClaimAmountReceived, evOk := ev.GetAttribute(types.AttributeKeyClaimAmountReceived)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimAmountReceived, attrClaimAmountReceived.Key)
					s.Require().Equal(tc.expReceiverAmount.String(), attrClaimAmountReceived.Value, "AttributeKeyClaimAmountReceived test name %s", tc.name)

					attrClaimValidatorFee, evOk := ev.GetAttribute(types.AttributeKeyClaimValidatorFee)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimValidatorFee, attrClaimValidatorFee.Key)
					s.Require().Equal(tc.expValFee.String(), attrClaimValidatorFee.Value, "AttributeKeyClaimValidatorFee test name %s", tc.name)
				}
			}
		} else {
			s.Require().False(hasEvent)
		}

		// check results
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Require().Equal(tc.expRemainDeposit, stream.Deposit, "tc.expRemainDeposit Equal stream.Deposit test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestClaimFromStream_Fail_NotExist() {
	c, v, t, d, err := s.app.StreamKeeper.ClaimFromStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().ErrorContains(err, "stream does not exist")
	s.Require().Equal(sdk.Coin{}, c)
	s.Require().Equal(sdk.Coin{}, v)
	s.Require().Equal(sdk.Coin{}, t)
	s.Require().Equal(sdk.Coin{}, d)

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().False(ok)
	s.Require().Equal(types.Stream{}, stream)
}

func (s *KeeperTestSuite) TestClaimFromStream_Fail_NoDeposit() {
	testCases := []struct {
		name     string
		sender   sdk.AccAddress
		receiver sdk.AccAddress
		stream   types.Stream
	}{
		{
			name:     "zero deposit",
			sender:   s.addrs[0],
			receiver: s.addrs[1],
			stream: types.Stream{
				Deposit: sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			},
		},
		{
			name:     "nil deposit",
			sender:   s.addrs[2],
			receiver: s.addrs[3],
			stream: types.Stream{
				Deposit: sdk.Coin{},
			},
		},
	}

	for _, tc := range testCases {
		err := s.app.StreamKeeper.SetStream(s.ctx, tc.receiver, tc.sender, tc.stream)
		s.Require().NoError(err, "SetStream NoError test name %s", tc.name)
		c, v, t, d, err := s.app.StreamKeeper.ClaimFromStream(s.ctx, tc.receiver, tc.sender)
		s.Require().ErrorContains(err, "stream deposit is zero")
		s.Require().Equal(sdk.Coin{}, c)
		s.Require().Equal(sdk.Coin{}, v)
		s.Require().Equal(sdk.Coin{}, t)
		s.Require().Equal(sdk.Coin{}, d)
	}
}

func (s *KeeperTestSuite) TestCancelStreamBySenderReceiver_Success() {
	tCtx := s.ctx

	blockTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(1)
	nowTime := tCtx.BlockTime()

	// set validator fee
	valFee, err := sdk.NewDecFromStr("0.01")
	s.app.StreamKeeper.SetParams(tCtx, types.Params{ValidatorFee: valFee})

	deposit := sdk.NewCoin("stake", sdk.NewIntFromUint64(1000))

	// create stream
	_, err = s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[1], s.addrs[0], deposit, 1)
	s.Require().NoError(err)
	// add deposit
	ok, err := s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[1], s.addrs[0], deposit)
	s.Require().True(ok)
	s.Require().NoError(err)

	// time travel
	future := time.Unix(nowTime.Unix()+500, 0).UTC()
	tCtx = tCtx.WithBlockTime(future).WithBlockHeight(2)

	// claim
	err = s.app.StreamKeeper.CancelStreamBySenderReceiver(tCtx, s.addrs[1], s.addrs[0])
	s.Require().NoError(err)

	// check event emission
	events := tCtx.EventManager().Events()

	hasCancelEvent := false
	hasClaimEvent := false

	for _, ev := range events {
		if ev.Type == types.EventTypeClaimStreamAction {
			hasClaimEvent = true
			attrClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyClaimTotal)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyClaimTotal, attrClaimTotal.Key)
			s.Require().Equal("500stake", attrClaimTotal.Value)

			attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
			// claim event occurs before cancel, so still has 500
			s.Require().Equal("500stake", attrRemainingDeposit.Value)

			attrClaimAmountReceived, evOk := ev.GetAttribute(types.AttributeKeyClaimAmountReceived)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyClaimAmountReceived, attrClaimAmountReceived.Key)
			s.Require().Equal("495stake", attrClaimAmountReceived.Value)

			attrClaimValidatorFee, evOk := ev.GetAttribute(types.AttributeKeyClaimValidatorFee)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyClaimValidatorFee, attrClaimValidatorFee.Key)
			s.Require().Equal("5stake", attrClaimValidatorFee.Value)
		}

		if ev.Type == types.EventTypeStreamCancelled {
			hasCancelEvent = true

			attrSender, ok := ev.GetAttribute(types.AttributeKeyStreamSender)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamSender, attrSender.Key)
			s.Require().Equal(s.addrs[0].String(), attrSender.Value)

			attrReceiver, ok := ev.GetAttribute(types.AttributeKeyStreamReceiver)
			s.Require().True(ok)
			s.Require().Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
			s.Require().Equal(s.addrs[1].String(), attrReceiver.Value)

			attrRefundAmount, evOk := ev.GetAttribute(types.AttributeKeyRefundAmount)
			s.Require().True(evOk)
			s.Require().Equal(types.AttributeKeyRefundAmount, attrRefundAmount.Key)
			s.Require().Equal("500stake", attrRefundAmount.Value)
		}
	}

	// should emit stream_deposit and claim events
	s.Require().True(hasCancelEvent)
	s.Require().True(hasClaimEvent)

	// check stream in keeper
	stream, ok := s.app.StreamKeeper.GetStream(tCtx, s.addrs[1], s.addrs[0])
	s.Require().True(ok)
	s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(0)), stream.Deposit, "stream.Deposit")
}

func (s *KeeperTestSuite) TestCancelStreamBySenderReceiver_Scenarios() {
	testCases := []struct {
		name              string
		sender            sdk.AccAddress
		receiver          sdk.AccAddress
		flowRate          int64
		deposit           sdk.Coin
		valFee            string
		cancelTimeOffset  int64    // seconds from "now"
		expTotalClaim     sdk.Coin // after cancel
		expReceiverAmount sdk.Coin // after cancel
		expRemainDeposit  sdk.Coin // after cancel - claim event only
		expValFee         sdk.Coin // after cancel
		expRefundAmount   sdk.Coin // after cancel
	}{
		{
			name:              "has expired, has deposit, no val fee",
			sender:            s.addrs[0],
			receiver:          s.addrs[1],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.0",
			cancelTimeOffset:  1001,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRefundAmount:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:              "not expired, has deposit, no val fee",
			sender:            s.addrs[2],
			receiver:          s.addrs[3],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.0",
			cancelTimeOffset:  500,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(500)), // claim before cancel, so will have remaining deposit
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRefundAmount:   sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
		},
		{
			name:              "has expired, no deposit left, no val fee",
			sender:            s.addrs[4],
			receiver:          s.addrs[5],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			valFee:            "0.0",
			cancelTimeOffset:  1001,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expRefundAmount:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:              "has expired, has deposit, 3% val fee",
			sender:            s.addrs[6],
			receiver:          s.addrs[7],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.03",
			cancelTimeOffset:  1001,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(970)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(30)),
			expRefundAmount:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		},
		{
			name:              "not expired, has deposit, 3% val fee",
			sender:            s.addrs[8],
			receiver:          s.addrs[9],
			flowRate:          1,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
			valFee:            "0.03",
			cancelTimeOffset:  500,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(485)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(15)),
			expRefundAmount:   sdk.NewCoin("stake", sdk.NewIntFromUint64(500)),
		},
		{
			name:              "not expired, has deposit, 4% val fee",
			sender:            s.addrs[10],
			receiver:          s.addrs[11],
			flowRate:          2342343,
			deposit:           sdk.NewCoin("stake", sdk.NewIntFromUint64(6155677404000)),
			valFee:            "0.04",
			cancelTimeOffset:  1814421,
			expTotalClaim:     sdk.NewCoin("stake", sdk.NewIntFromUint64(4249996328403)),
			expReceiverAmount: sdk.NewCoin("stake", sdk.NewIntFromUint64(4079996475267)),
			expRemainDeposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(1905681075597)),
			expValFee:         sdk.NewCoin("stake", sdk.NewIntFromUint64(169999853136)),
			expRefundAmount:   sdk.NewCoin("stake", sdk.NewIntFromUint64(1905681075597)),
		},
	}

	for _, tc := range testCases {
		tCtx := s.ctx
		nowTime := time.Unix(time.Now().Unix(), 0).UTC()
		// set create time to past
		blockTimeCreate := nowTime
		blockNum := int64(1)
		tCtx = tCtx.WithBlockTime(blockTimeCreate).WithBlockHeight(blockNum)
		blockNum += 1

		// set params
		newParams := types.Params{
			ValidatorFee: sdk.MustNewDecFromStr(tc.valFee),
		}

		err := s.app.StreamKeeper.SetParams(tCtx, newParams)
		s.Require().NoError(err, "SetParams NoError test name %s", tc.name)

		// create
		_, err = s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.deposit, tc.flowRate)

		s.Require().NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		if tc.deposit.IsPositive() {
			ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.deposit)
			s.Require().True(ok, "initialDeposit ok NoError test name %s", tc.name)
			s.Require().NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)
		}

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)

		cancelTime := time.Unix(nowTime.Unix()+tc.cancelTimeOffset, 0).UTC()
		tCtx = tCtx.WithBlockTime(cancelTime).WithBlockHeight(blockNum)

		// cancel
		err = s.app.StreamKeeper.CancelStreamBySenderReceiver(tCtx, tc.receiver, tc.sender)
		s.Require().NoError(err)

		events := tCtx.EventManager().Events()
		hasClaimEvent := false
		hasCancelEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
				attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
				// only for this stream
				if tc.sender.String() == attrSender.Value && tc.receiver.String() == attrReceiver.Value {
					hasClaimEvent = true
				}
			}
			if ev.Type == types.EventTypeStreamCancelled {
				attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
				attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
				// only for this stream
				if tc.sender.String() == attrSender.Value && tc.receiver.String() == attrReceiver.Value {
					hasCancelEvent = true
				}
			}
		}

		if tc.expReceiverAmount.Amount.IsPositive() {
			s.Require().True(hasClaimEvent)
			for _, ev := range events {
				if ev.Type == types.EventTypeClaimStreamAction {
					attrSender, _ := ev.GetAttribute(types.AttributeKeyStreamSender)
					attrReceiver, _ := ev.GetAttribute(types.AttributeKeyStreamReceiver)
					// only for this stream
					if tc.sender.String() != attrSender.Value || tc.receiver.String() != attrReceiver.Value {
						// skip events not for this stream
						continue
					}

					attrClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyClaimTotal)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimTotal, attrClaimTotal.Key)
					s.Require().Equal(tc.expTotalClaim.String(), attrClaimTotal.Value, "AttributeKeyClaimTotal test name %s", tc.name)

					attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
					s.Require().Equal(tc.expRemainDeposit.String(), attrRemainingDeposit.Value, "AttributeKeyRemainingDeposit test name %s", tc.name)

					attrClaimAmountReceived, evOk := ev.GetAttribute(types.AttributeKeyClaimAmountReceived)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimAmountReceived, attrClaimAmountReceived.Key)
					s.Require().Equal(tc.expReceiverAmount.String(), attrClaimAmountReceived.Value, "AttributeKeyClaimAmountReceived test name %s", tc.name)

					attrClaimValidatorFee, evOk := ev.GetAttribute(types.AttributeKeyClaimValidatorFee)
					s.Require().True(evOk)
					s.Require().Equal(types.AttributeKeyClaimValidatorFee, attrClaimValidatorFee.Key)
					s.Require().Equal(tc.expValFee.String(), attrClaimValidatorFee.Value, "AttributeKeyClaimValidatorFee test name %s", tc.name)
				}
			}
		} else {
			s.Require().False(hasClaimEvent)
		}

		for _, ev := range events {
			if ev.Type == types.EventTypeStreamCancelled {
				attrSender, sOk := ev.GetAttribute(types.AttributeKeyStreamSender)
				attrReceiver, rOk := ev.GetAttribute(types.AttributeKeyStreamReceiver)
				// only for this stream
				if tc.sender.String() != attrSender.Value || tc.receiver.String() != attrReceiver.Value {
					continue
				}

				s.Require().True(sOk)
				s.Require().Equal(types.AttributeKeyStreamSender, attrSender.Key)
				s.Require().Equal(tc.sender.String(), attrSender.Value)

				s.Require().True(rOk)
				s.Require().Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
				s.Require().Equal(tc.receiver.String(), attrReceiver.Value)

				attrRefundAmount, evOk := ev.GetAttribute(types.AttributeKeyRefundAmount)
				s.Require().True(evOk)
				s.Require().Equal(types.AttributeKeyRefundAmount, attrRefundAmount.Key)
				s.Require().Equal(tc.expRefundAmount.String(), attrRefundAmount.Value)

				attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
				s.Require().True(evOk)
				s.Require().Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
				s.Require().Equal("0stake", attrRemainingDeposit.Value)
			}
		}

		s.Require().True(hasCancelEvent)

		// check results
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.Require().True(ok, "GetStream ok NoError test name %s", tc.name)
		// should be nothing left in the stream
		s.Require().Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(0)), stream.Deposit, "0stake Equal stream.Deposit test name %s", tc.name)
		s.Require().Equal(tCtx.BlockTime(), stream.DepositZeroTime, "stream.DepositZeroTime is now test name %s", tc.name)
		s.Require().Equal(int64(0), stream.FlowRate, "stream.FlowRate is zero test name %s", tc.name)
		s.Require().Equal(tCtx.BlockTime(), stream.LastUpdatedTime, "stream.LastUpdatedTime is now test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestCancelStreamBySenderReceiver_Fail_NotExist() {
	err := s.app.StreamKeeper.CancelStreamBySenderReceiver(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().ErrorContains(err, "stream does not exist")

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.Require().False(ok)
	s.Require().Equal(types.Stream{}, stream)
}

func (s *KeeperTestSuite) TestGetTotalDeposits() {

	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(nowTime).WithBlockHeight(1)

	// create and deposit
	for i := int64(1); i <= 10; i++ {
		deposit := sdk.NewInt64Coin("stake", 1000*i)
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[i-1], s.addrs[i], deposit, i)
		s.Require().NoError(err)
		ok, err := s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[i-1], s.addrs[i], deposit)
		s.Require().NoError(err)
		s.Require().True(ok)
	}

	expectedDeposits := sdk.NewCoins(sdk.NewInt64Coin("stake", 55000))
	currentDeposits := s.app.StreamKeeper.GetTotalDeposits(tCtx)
	s.Require().Equal(expectedDeposits, currentDeposits)

	// move 500 seconds
	blockTime := time.Unix(nowTime.Unix()+500, 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(2)

	// claim - not expired
	for i := int64(1); i <= 10; i++ {
		_, _, _, _, err := s.app.StreamKeeper.ClaimFromStream(tCtx, s.addrs[i-1], s.addrs[i])
		s.Require().NoError(err)
	}

	expectedDeposits = sdk.NewCoins(sdk.NewInt64Coin("stake", 27500))
	currentDeposits = s.app.StreamKeeper.GetTotalDeposits(tCtx)
	s.Require().Equal(expectedDeposits, currentDeposits)

	// move 1001 seconds
	blockTime = time.Unix(nowTime.Unix()+1001, 0).UTC()
	tCtx = tCtx.WithBlockTime(blockTime).WithBlockHeight(3)

	// claim - expired
	for i := int64(1); i <= 10; i++ {
		_, _, _, _, err := s.app.StreamKeeper.ClaimFromStream(tCtx, s.addrs[i-1], s.addrs[i])
		s.Require().NoError(err)
	}

	expectedDeposits = sdk.NewCoins(sdk.NewInt64Coin("stake", 0))
	currentDeposits = s.app.StreamKeeper.GetTotalDeposits(tCtx)
	s.Require().Equal(expectedDeposits, currentDeposits)
}

func (s *KeeperTestSuite) TestIterateAllStreams() {

	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()
	tCtx = tCtx.WithBlockTime(nowTime)

	streams := map[string]map[string]types.Stream{}

	for i := int64(1); i < 100; i++ {
		deposit := sdk.NewInt64Coin("stake", 1000*i)
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[i-1], s.addrs[i], deposit, i)
		s.Require().NoError(err)
		stream, _ := s.app.StreamKeeper.GetStream(tCtx, s.addrs[i-1], s.addrs[i])
		streams[s.addrs[i-1].String()] = map[string]types.Stream{}
		streams[s.addrs[i-1].String()][s.addrs[i].String()] = stream
	}

	s.app.StreamKeeper.IterateAllStreams(tCtx, func(receiverAddr, senderAddr sdk.AccAddress, stream types.Stream) bool {

		expectedStream, exists := streams[receiverAddr.String()][senderAddr.String()]
		s.Require().True(exists)
		s.Require().Equal(expectedStream, stream)

		return false
	})
}
