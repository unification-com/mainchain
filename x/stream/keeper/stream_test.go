package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/stream/types"
	"strconv"
	"time"
)

func (s *KeeperTestSuite) TestSetGetHighestStreamId() {
	// should be 1 by default
	highestId, err := s.app.StreamKeeper.GetHighestStreamId(s.ctx)
	s.Require().NoError(err)
	s.Equal(uint64(1), highestId)

	// should set and get
	expHighestId := uint64(24)
	s.app.StreamKeeper.SetHighestStreamId(s.ctx, expHighestId)
	highestId, err = s.app.StreamKeeper.GetHighestStreamId(s.ctx)
	s.Require().NoError(err)
	s.Equal(expHighestId, highestId)
}

func (s *KeeperTestSuite) TestIsStream() {
	ok := s.app.StreamKeeper.IsStream(s.ctx, s.addrs[1], s.addrs[0])
	s.False(ok)

	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		StreamId:        1,
		Sender:          s.addrs[0].String(),
		Receiver:        s.addrs[1].String(),
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
	s.True(ok)
}

func (s *KeeperTestSuite) TestSetGetStream() {

	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.False(ok)
	s.Require().Equal(types.Stream{}, stream)

	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		StreamId:        1,
		Sender:          s.addrs[0].String(),
		Receiver:        s.addrs[1].String(),
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

	s.True(ok)
	s.Require().Equal(expStream.StreamId, stream.StreamId)
	s.Require().Equal(expStream.Sender, stream.Sender)
	s.Require().Equal(expStream.Receiver, stream.Receiver)
	s.Require().Equal(expStream.Deposit, stream.Deposit)
	s.Require().Equal(expStream.FlowRate, stream.FlowRate)
	s.Require().Equal(expStream.CreateTime, stream.CreateTime)
	s.Require().Equal(expStream.LastUpdatedTime, stream.LastUpdatedTime)
	s.Require().Equal(expStream.LastOutflowTime, stream.LastOutflowTime)
	s.Require().Equal(expStream.DepositZeroTime, stream.DepositZeroTime)
	s.Require().Equal(expStream.TotalStreamed, stream.TotalStreamed)
	s.Require().Equal(expStream.Cancellable, stream.Cancellable)
}

func (s *KeeperTestSuite) TestSetGetIdLookup() {
	// default empty, doesn't exist
	lookup, ok := s.app.StreamKeeper.GetIdLookup(s.ctx, 24)
	s.False(ok)
	s.Equal(types.StreamIdLookup{}, lookup)

	// set
	expIdLookup := types.StreamIdLookup{
		Sender:   s.addrs[0].String(),
		Receiver: s.addrs[1].String(),
	}

	err := s.app.StreamKeeper.SetIdLookup(s.ctx, 24, expIdLookup)
	s.NoError(err)

	// get
	lookup, ok = s.app.StreamKeeper.GetIdLookup(s.ctx, 24)
	s.True(ok)
	s.Equal(expIdLookup, lookup)
}

func (s *KeeperTestSuite) TestSetGetIdLookupFromStream() {
	lookup, ok := s.app.StreamKeeper.GetIdLookup(s.ctx, 24)
	s.False(ok)
	s.Equal(types.StreamIdLookup{}, lookup)

	// set stream
	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		StreamId:        1,
		Sender:          s.addrs[0].String(),
		Receiver:        s.addrs[1].String(),
		Deposit:         sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)),
		FlowRate:        100,
		CreateTime:      nowTime,
		LastUpdatedTime: nowTime,
		LastOutflowTime: nowTime,
		DepositZeroTime: nowTime,
		TotalStreamed:   sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
		Cancellable:     true,
	}

	err := s.app.StreamKeeper.SetStream(s.ctx, s.addrs[1], s.addrs[0], expStream)
	s.Require().NoError(err)

	// set lookup from stream info
	expIdLookup := types.StreamIdLookup{
		Sender:   expStream.Sender,
		Receiver: expStream.Receiver,
	}

	err = s.app.StreamKeeper.SetIdLookup(s.ctx, 1, expIdLookup)
	s.NoError(err)

	// get lookup
	lookup, ok = s.app.StreamKeeper.GetIdLookup(s.ctx, 1)
	s.True(ok)
	s.Equal(expIdLookup, lookup)

	sender, _ := sdk.AccAddressFromBech32(lookup.Sender)
	receiver, _ := sdk.AccAddressFromBech32(lookup.Receiver)

	// get stream
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, receiver, sender)

	s.True(ok)
	s.Require().Equal(expStream.StreamId, stream.StreamId)
	s.Require().Equal(expStream.Sender, stream.Sender)
	s.Require().Equal(expStream.Receiver, stream.Receiver)
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
		StreamId:        1,
		Sender:          s.addrs[0].String(),
		Receiver:        s.addrs[1].String(),
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

	s.NoError(err)

	s.Require().Equal(expStream.StreamId, stream.StreamId)
	s.Require().Equal(expStream.Sender, stream.Sender)
	s.Require().Equal(expStream.Receiver, stream.Receiver)
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
			attrStreamId, ok := ev.GetAttribute(types.AttributeKeyStreamId)
			s.True(ok)
			s.Equal(types.AttributeKeyStreamId, attrStreamId.Key)
			s.Equal("1", attrStreamId.Value)

			attrSender, ok := ev.GetAttribute(types.AttributeKeyStreamSender)
			s.True(ok)
			s.Equal(types.AttributeKeyStreamSender, attrSender.Key)
			s.Equal(expStream.Sender, attrSender.Value)

			attrReceiver, ok := ev.GetAttribute(types.AttributeKeyStreamReceiver)
			s.True(ok)
			s.Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
			s.Equal(expStream.Receiver, attrReceiver.Value)

			attrFlowRate, ok := ev.GetAttribute(types.AttributeKeyStreamFlowRate)
			s.True(ok)
			s.Equal(types.AttributeKeyStreamFlowRate, attrFlowRate.Key)
			s.Equal("100", attrFlowRate.Value)
		}
	}

	// should emit create_stream event
	s.True(hasCreateStreamEvent)

	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])

	s.True(ok)
	s.Require().Equal(expStream.StreamId, stream.StreamId)
	s.Require().Equal(expStream.Sender, stream.Sender)
	s.Require().Equal(expStream.Receiver, stream.Receiver)
	s.Require().Equal(expStream.Deposit, stream.Deposit)
	s.Require().Equal(expStream.FlowRate, stream.FlowRate)
	s.Require().Equal(expStream.CreateTime, stream.CreateTime)
	s.Require().Equal(expStream.LastUpdatedTime, stream.LastUpdatedTime)
	s.Require().Equal(expStream.LastOutflowTime, stream.LastOutflowTime)
	s.Require().Equal(expStream.DepositZeroTime, stream.DepositZeroTime)
	s.Require().Equal(expStream.TotalStreamed, stream.TotalStreamed)
	s.Require().Equal(expStream.Cancellable, stream.Cancellable)
}

func (s *KeeperTestSuite) TestCreateIdLookup() {
	// set stream
	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		StreamId:        1,
		Sender:          s.addrs[0].String(),
		Receiver:        s.addrs[1].String(),
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

	// set lookup from stream info
	expIdLookup := types.StreamIdLookup{
		Sender:   expStream.Sender,
		Receiver: expStream.Receiver,
	}

	err = s.app.StreamKeeper.CreateIdLookup(s.ctx, s.addrs[1], s.addrs[0], 1)
	s.NoError(err)

	// get lookup
	lookup, ok := s.app.StreamKeeper.GetIdLookup(s.ctx, 1)
	s.True(ok)
	s.Equal(expIdLookup, lookup)

	sender, _ := sdk.AccAddressFromBech32(lookup.Sender)
	receiver, _ := sdk.AccAddressFromBech32(lookup.Receiver)

	// get stream
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, receiver, sender)

	s.True(ok)
	s.Require().Equal(expStream.StreamId, stream.StreamId)
	s.Require().Equal(expStream.Sender, stream.Sender)
	s.Require().Equal(expStream.Receiver, stream.Receiver)
	s.Require().Equal(expStream.Deposit, stream.Deposit)
	s.Require().Equal(expStream.FlowRate, stream.FlowRate)
	s.Require().Equal(expStream.CreateTime, stream.CreateTime)
	s.Require().Equal(expStream.LastUpdatedTime, stream.LastUpdatedTime)
	s.Require().Equal(expStream.LastOutflowTime, stream.LastOutflowTime)
	s.Require().Equal(expStream.DepositZeroTime, stream.DepositZeroTime)
	s.Require().Equal(expStream.TotalStreamed, stream.TotalStreamed)
	s.Require().Equal(expStream.Cancellable, stream.Cancellable)
}

func (s *KeeperTestSuite) TestAddDeposit_Basic_Success() {
	// set stream
	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		StreamId:        1,
		Sender:          s.addrs[0].String(),
		Receiver:        s.addrs[1].String(),
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
	s.True(ok)
	s.NoError(err)

	// check events ar emitted
	events := s.ctx.EventManager().Events()

	hasEvent := false
	for _, ev := range events {
		if ev.Type == types.EventTypeDepositToStream {
			hasEvent = true
			attrStreamId, evOk := ev.GetAttribute(types.AttributeKeyStreamId)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamId, attrStreamId.Key)
			s.Equal("1", attrStreamId.Value)

			attrDepositAmount, evOk := ev.GetAttribute(types.AttributeKeyStreamDepositAmount)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamDepositAmount, attrDepositAmount.Key)
			s.Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)).String(), attrDepositAmount.Value)

			attrDepositDuration, evOk := ev.GetAttribute(types.AttributeKeyStreamDepositDuration)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamDepositDuration, attrDepositDuration.Key)
			s.Equal("1000", attrDepositDuration.Value)

			attrDepositZeroTime, evOk := ev.GetAttribute(types.AttributeKeyStreamDepositZeroTime)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamDepositZeroTime, attrDepositZeroTime.Key)
			s.Equal(nowTime.Add(time.Second*1000).String(), attrDepositZeroTime.Value)
		}
	}

	// should emit stream_deposit event
	s.True(hasEvent)

	// get stream from keeper
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	// should now be 1000stake
	s.Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)), stream.Deposit)
	// Deposit of 1000, flow rate of 100/s, should have deposit zero time of now + 10s
	s.Equal(nowTime.Add(time.Second*1000), stream.DepositZeroTime)
	s.Equal(nowTime, stream.LastUpdatedTime)
}

func (s *KeeperTestSuite) TestAddDeposit_Success_TopUpExistingNotExpired() {
	tCtx := s.ctx

	blockTime := time.Now()
	tCtx = tCtx.WithBlockTime(blockTime)
	nowTime := tCtx.BlockTime()

	testCases := []struct {
		name               string
		stream             types.Stream
		deposit            sdk.Coin
		expDepositZeroTime time.Time
		expDeposit         sdk.Coin
		expDiff            int64
	}{
		{
			name: "1",
			stream: types.Stream{
				StreamId:        1,
				Sender:          s.addrs[0].String(),
				Receiver:        s.addrs[1].String(),
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
			name: "2",
			stream: types.Stream{
				StreamId:        2,
				Sender:          s.addrs[2].String(),
				Receiver:        s.addrs[3].String(),
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
			stream: types.Stream{
				StreamId:        3,
				Sender:          s.addrs[4].String(),
				Receiver:        s.addrs[5].String(),
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
		sendAddr, _ := sdk.AccAddressFromBech32(tc.stream.Sender)
		recAddr, _ := sdk.AccAddressFromBech32(tc.stream.Receiver)
		// deposit zero time is in the future, so just use SetStream instead of create & add deposit combo
		err := s.app.StreamKeeper.SetStream(tCtx, recAddr, sendAddr, tc.stream)
		s.NoError(err, "SetStream NoError test name %s", tc.name)
		ok, err := s.app.StreamKeeper.AddDeposit(tCtx, recAddr, sendAddr, tc.deposit)
		s.True(ok, "AddDeposit True test name %s", tc.name)
		s.NoError(err, "AddDeposit NoError test name %s", tc.name)

		// events should NOT contain claim_stream
		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				hasEvent = true
			}
		}
		s.False(hasEvent)

		stream, ok := s.app.StreamKeeper.GetStream(tCtx, recAddr, sendAddr)
		s.True(ok, "GetStream True test name %s", tc.name)
		s.Equal(tc.expDeposit, stream.Deposit, "GetStream Deposit Equal test name %s", tc.name)
		s.Equal(tc.expDepositZeroTime, stream.DepositZeroTime, "GetStream DepositZeroTime Equal test name %s", tc.name)

		duration := stream.DepositZeroTime.Unix() - stream.LastOutflowTime.Unix()
		s.Equal(tc.expDiff, duration, "duration test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestAddDeposit_Success_TopUpExistingExpired() {
	tCtx := s.ctx

	blockTime := time.Now()
	tCtx = tCtx.WithBlockTime(blockTime)
	nowTime := tCtx.BlockTime()

	testCases := []struct {
		name               string
		stream             types.Stream
		initialDeposit     sdk.Coin
		newDeposit         sdk.Coin
		expDepositZeroTime time.Time
		expDeposit         sdk.Coin
		expTotalStreamed   sdk.Coin
	}{
		{
			name: "1",
			stream: types.Stream{
				StreamId:        1,
				Sender:          s.addrs[0].String(),
				Receiver:        s.addrs[1].String(),
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
			name: "2",
			stream: types.Stream{
				StreamId:        2,
				Sender:          s.addrs[2].String(),
				Receiver:        s.addrs[3].String(),
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
			name: "3",
			stream: types.Stream{
				StreamId:        3,
				Sender:          s.addrs[4].String(),
				Receiver:        s.addrs[5].String(),
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
			name: "4",
			stream: types.Stream{
				StreamId:        4,
				Sender:          s.addrs[6].String(),
				Receiver:        s.addrs[7].String(),
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
		sendAddr, _ := sdk.AccAddressFromBech32(tc.stream.Sender)
		recAddr, _ := sdk.AccAddressFromBech32(tc.stream.Receiver)

		// create
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, recAddr, sendAddr, tc.initialDeposit, tc.stream.FlowRate)
		s.NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		if tc.initialDeposit.Amount.GT(sdk.NewIntFromUint64(0)) {
			ok, err := s.app.StreamKeeper.AddDeposit(tCtx, recAddr, sendAddr, tc.initialDeposit)
			s.True(ok)
			s.NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)
		}

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, recAddr, sendAddr)
		s.True(ok)
		s.Equal(tc.initialDeposit, stream.Deposit)

		// set times etc.
		stream.CreateTime = tc.stream.CreateTime
		stream.LastUpdatedTime = tc.stream.LastUpdatedTime
		stream.LastOutflowTime = tc.stream.LastOutflowTime
		stream.DepositZeroTime = tc.stream.DepositZeroTime
		err = s.app.StreamKeeper.SetStream(tCtx, recAddr, sendAddr, stream)
		s.NoError(err, "SetStream NoError test name %s", tc.name)

		// top up with new deposit
		ok, err = s.app.StreamKeeper.AddDeposit(tCtx, recAddr, sendAddr, tc.newDeposit)
		s.True(ok, "AddDeposit True test name %s", tc.name)
		s.NoError(err, "AddDeposit NoError test name %s", tc.name)

		// check events do contain claim_stream if expTotalStreamed > 0
		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				hasEvent = true
			}
		}
		if tc.expTotalStreamed.IsPositive() {
			s.True(hasEvent)
		} else {
			s.False(hasEvent)
		}

		// final check
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, recAddr, sendAddr)
		s.True(ok, "GetStream True test name %s", tc.name)
		s.Equal(tc.expDeposit, stream.Deposit, "GetStream Deposit Equal test name %s", tc.name)
		s.Equal(tc.expDepositZeroTime, stream.DepositZeroTime, "GetStream DepositZeroTime Equal test name %s", tc.name)
		s.Equal(tc.expTotalStreamed, stream.TotalStreamed, "GetStream TotalStreamed Equal test name %s", tc.name)
		s.Equal(nowTime, stream.LastUpdatedTime, "GetStream TotalStreamed Equal test name %s", tc.name)
		if tc.expTotalStreamed.IsPositive() {
			s.Equal(nowTime, stream.LastOutflowTime, "GetStream TotalStreamed Equal test name %s", tc.name)
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
		tCtx = tCtx.WithBlockTime(blockTimeCreate)

		// create
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.initialDeposit, tc.flowRate)

		s.NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.initialDeposit)
		s.True(ok, "initialDeposit ok NoError test name %s", tc.name)
		s.NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		// should be in the future from the creation time
		expInitialDepZeroTime := time.Unix(blockTimeCreate.Unix()+tc.expInitialDepZeroTime, 0).UTC()
		s.True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Equal(tc.initialDeposit, stream.Deposit, "tc.initialDeposit Equal stream.Deposit test name %s", tc.name)
		s.Equal(expInitialDepZeroTime, stream.DepositZeroTime, "expInitialDepZeroTimeEqual stream.DepositZeroTime test name %s", tc.name)

		// set block time to now
		tCtx = tCtx.WithBlockTime(nowTime)

		// add new deposit
		ok, err = s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.newDeposit)
		s.True(ok, "newDeposit AddDeposit ok test name %s", tc.name)
		s.NoError(err, "newDeposit AddDeposit NoError test name %s", tc.name)

		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				attrStreamId, _ := ev.GetAttribute(types.AttributeKeyStreamId)
				// only for this stream
				if strconv.Itoa(int(stream.StreamId)) == attrStreamId.Value {
					hasEvent = true
				}
			}
		}

		if tc.expClaim.Amount.IsPositive() {
			s.True(hasEvent)
			for _, ev := range events {
				if ev.Type == types.EventTypeClaimStreamAction {
					attrStreamId, _ := ev.GetAttribute(types.AttributeKeyStreamId)
					if strconv.Itoa(int(stream.StreamId)) != attrStreamId.Value {
						// skip events not for this stream
						continue
					}
					attrClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyStreamClaimTotal)
					s.True(evOk)
					s.Equal(types.AttributeKeyStreamClaimTotal, attrClaimTotal.Key)
					s.Equal(tc.expClaim.String(), attrClaimTotal.Value, "AttributeKeyStreamClaimTotal test name %s", tc.name)

					attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
					s.True(evOk)
					s.Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
					s.Equal(tc.expRemainDeposit.String(), attrRemainingDeposit.Value, "AttributeKeyRemainingDeposit test name %s", tc.name)
				}
			}
		} else {
			s.False(hasEvent)
		}

		// check results
		expDepZeroTime := time.Unix(nowTime.Unix()+tc.expNewDepZeroTime, 0).UTC()
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Equal(tc.expNewDeposit, stream.Deposit, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
		s.Equal(expDepZeroTime, stream.DepositZeroTime, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestAddDeposit_Fail_StreamNotExist() {
	ok, err := s.app.StreamKeeper.AddDeposit(s.ctx, s.addrs[1], s.addrs[0], sdk.NewCoin("stake", sdk.NewIntFromUint64(1000)))
	s.False(ok)
	s.ErrorContains(err, "stream does not exist")

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.False(ok)
	s.Equal(types.Stream{}, stream)
}

func (s *KeeperTestSuite) TestAddDeposit_Fail_InsufficientBalance() {
	// set stream
	nowTime := s.ctx.BlockTime()

	expStream := types.Stream{
		StreamId:        1,
		Sender:          s.addrs[0].String(),
		Receiver:        s.addrs[1].String(),
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
	s.False(ok)
	s.ErrorContains(err, "insufficient funds")

}

func (s *KeeperTestSuite) TestSetNewFlowRate_Success() {
	// set stream
	tCtx := s.ctx

	blockTime := time.Now()
	tCtx = tCtx.WithBlockTime(blockTime)
	nowTime := tCtx.BlockTime()
	deposit := sdk.NewCoin("stake", sdk.NewIntFromUint64(2400))

	// create stream
	_, err := s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[1], s.addrs[0], deposit, 1)
	s.Require().NoError(err)

	// add deposit
	ok, err := s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[1], s.addrs[0], deposit)
	s.True(ok)
	s.Require().NoError(err)

	// Set new flow rate
	err = s.app.StreamKeeper.SetNewFlowRate(tCtx, s.addrs[1], s.addrs[0], 24)
	s.NoError(err)

	// check events ar emitted
	events := tCtx.EventManager().Events()

	hasEvent := false
	for _, ev := range events {
		if ev.Type == types.EventTypeUpdateFlowRate {
			hasEvent = true
			attrStreamId, evOk := ev.GetAttribute(types.AttributeKeyStreamId)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamId, attrStreamId.Key)
			s.Equal("1", attrStreamId.Value)

			attrOldFlowRate, evOk := ev.GetAttribute(types.AttributeKeyOldFlowRate)
			s.True(evOk)
			s.Equal(types.AttributeKeyOldFlowRate, attrOldFlowRate.Key)
			s.Equal("1", attrOldFlowRate.Value)

			attrNewFlowRate, evOk := ev.GetAttribute(types.AttributeKeyNewFlowRate)
			s.True(evOk)
			s.Equal(types.AttributeKeyNewFlowRate, attrNewFlowRate.Key)
			s.Equal("24", attrNewFlowRate.Value)

			attrStreamDepositDuration, evOk := ev.GetAttribute(types.AttributeKeyStreamDepositDuration)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamDepositDuration, attrStreamDepositDuration.Key)
			s.Equal("100", attrStreamDepositDuration.Value)

			attrStreamDepositZeroTime, evOk := ev.GetAttribute(types.AttributeKeyStreamDepositZeroTime)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamDepositZeroTime, attrStreamDepositZeroTime.Key)
			s.Equal(nowTime.Add(time.Second*100).String(), attrStreamDepositZeroTime.Value)
		}
	}

	// should emit stream_deposit event
	s.True(hasEvent)

	// get stream from keeper
	stream, ok := s.app.StreamKeeper.GetStream(tCtx, s.addrs[1], s.addrs[0])
	s.True(ok)
	// should now be 1000stake
	s.Equal(int64(24), stream.FlowRate)
	s.Equal(nowTime.Add(time.Second*100), stream.DepositZeroTime)
	s.Equal(nowTime, stream.LastUpdatedTime)
}

func (s *KeeperTestSuite) TestSetNewFlowRate_Success_ExistingNotExpired() {
	tCtx := s.ctx

	blockTime := time.Now()
	tCtx = tCtx.WithBlockTime(blockTime)
	nowTime := tCtx.BlockTime()

	testCases := []struct {
		name               string
		stream             types.Stream
		newFlowRate        int64
		expDepositZeroTime time.Time
	}{
		{
			name: "1",
			stream: types.Stream{
				StreamId:        1,
				Sender:          s.addrs[0].String(),
				Receiver:        s.addrs[1].String(),
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
			name: "2", // effectively expired
			stream: types.Stream{
				StreamId:        2,
				Sender:          s.addrs[2].String(),
				Receiver:        s.addrs[3].String(),
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
			name: "3",
			stream: types.Stream{
				StreamId:        3,
				Sender:          s.addrs[4].String(),
				Receiver:        s.addrs[5].String(),
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
			name: "4",
			stream: types.Stream{
				StreamId:        4,
				Sender:          s.addrs[6].String(),
				Receiver:        s.addrs[7].String(),
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
			name: "5",
			stream: types.Stream{
				StreamId:        5,
				Sender:          s.addrs[8].String(),
				Receiver:        s.addrs[9].String(),
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
		sendAddr, _ := sdk.AccAddressFromBech32(tc.stream.Sender)
		recAddr, _ := sdk.AccAddressFromBech32(tc.stream.Receiver)
		err := s.app.StreamKeeper.SetStream(tCtx, recAddr, sendAddr, tc.stream)
		s.NoError(err, "SetStream NoError test name %s", tc.name)
		err = s.app.StreamKeeper.SetNewFlowRate(tCtx, recAddr, sendAddr, tc.newFlowRate)
		s.NoError(err, "AddDeposit NoError test name %s", tc.name)

		stream, ok := s.app.StreamKeeper.GetStream(tCtx, recAddr, sendAddr)
		s.True(ok, "GetStream True test name %s", tc.name)
		s.Equal(tc.expDepositZeroTime, stream.DepositZeroTime, "GetStream DepositZeroTime Equal test name %s", tc.name)
		s.Equal(tc.newFlowRate, stream.FlowRate, "GetStream FlowRate Equal test name %s", tc.name)
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
			sender:                s.addrs[16],
			receiver:              s.addrs[17],
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
		tCtx = tCtx.WithBlockTime(blockTimeCreate)

		// create
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, tc.receiver, tc.sender, tc.initialDeposit, tc.startFlowRate)

		s.NoError(err, "CreateNewStream NoError test name %s", tc.name)

		// add initial deposit
		if tc.initialDeposit.IsPositive() {
			ok, err := s.app.StreamKeeper.AddDeposit(tCtx, tc.receiver, tc.sender, tc.initialDeposit)
			s.True(ok, "initialDeposit ok NoError test name %s", tc.name)
			s.NoError(err, "initialDeposit AddDeposit NoError test name %s", tc.name)
		}

		// check stream
		stream, ok := s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		// should be in the future from the creation time
		expInitialDepZeroTime := time.Unix(blockTimeCreate.Unix()+tc.expInitialDepZeroTime, 0).UTC()
		s.True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Equal(tc.initialDeposit, stream.Deposit, "tc.initialDeposit Equal stream.Deposit test name %s", tc.name)
		if tc.initialDeposit.IsPositive() {
			s.Equal(expInitialDepZeroTime, stream.DepositZeroTime, "expInitialDepZeroTimeEqual stream.DepositZeroTime test name %s", tc.name)
		}
		// set block time to now
		tCtx = tCtx.WithBlockTime(nowTime)

		// set new flow rate
		err = s.app.StreamKeeper.SetNewFlowRate(tCtx, tc.receiver, tc.sender, tc.newFlowRate)
		s.NoError(err, "newDeposit AddDeposit NoError test name %s", tc.name)

		events := tCtx.EventManager().Events()
		hasEvent := false
		for _, ev := range events {
			if ev.Type == types.EventTypeClaimStreamAction {
				attrStreamId, _ := ev.GetAttribute(types.AttributeKeyStreamId)
				// only for this stream
				if strconv.Itoa(int(stream.StreamId)) == attrStreamId.Value {
					hasEvent = true
				}
			}
		}

		if tc.expClaim.Amount.IsPositive() {
			s.True(hasEvent)
			for _, ev := range events {
				if ev.Type == types.EventTypeClaimStreamAction {
					attrStreamId, _ := ev.GetAttribute(types.AttributeKeyStreamId)
					if strconv.Itoa(int(stream.StreamId)) != attrStreamId.Value {
						// skip events not for this stream
						continue
					}
					attrClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyStreamClaimTotal)
					s.True(evOk)
					s.Equal(types.AttributeKeyStreamClaimTotal, attrClaimTotal.Key)
					s.Equal(tc.expClaim.String(), attrClaimTotal.Value, "AttributeKeyStreamClaimTotal test name %s", tc.name)

					attrRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
					s.True(evOk)
					s.Equal(types.AttributeKeyRemainingDeposit, attrRemainingDeposit.Key)
					s.Equal(tc.expRemainDeposit.String(), attrRemainingDeposit.Value, "AttributeKeyRemainingDeposit test name %s", tc.name)
				}
			}
		} else {
			s.False(hasEvent)
		}

		// check results
		expDepZeroTime := time.Unix(nowTime.Unix()+tc.expNewDepZeroTime, 0).UTC()
		stream, ok = s.app.StreamKeeper.GetStream(tCtx, tc.receiver, tc.sender)
		s.True(ok, "GetStream ok NoError test name %s", tc.name)
		s.Equal(tc.newFlowRate, stream.FlowRate, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
		s.Equal(expDepZeroTime, stream.DepositZeroTime, "tc.expNewDeposit Equal stream.Deposit test name %s", tc.name)
	}
}

func (s *KeeperTestSuite) TestSetNewFlowRate_Fail() {
	err := s.app.StreamKeeper.SetNewFlowRate(s.ctx, s.addrs[1], s.addrs[0], 24)
	s.ErrorContains(err, "stream does not exist")

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.False(ok)
	s.Equal(types.Stream{}, stream)
}

func (s *KeeperTestSuite) TestClaimFromStream_Success() {
	// set stream
	tCtx := s.ctx

	blockTime := time.Now()
	tCtx = tCtx.WithBlockTime(blockTime)
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
	s.True(ok)
	s.Require().NoError(err)

	// time travel
	future := time.Unix(nowTime.Unix()+501, 0).UTC()
	tCtx = tCtx.WithBlockTime(future)

	// claim
	amntClaimed, valFeeSent, totalClaim, remainingDeposit, err := s.app.StreamKeeper.ClaimFromStream(tCtx, s.addrs[1], s.addrs[0])
	s.Require().NoError(err)
	s.Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(495)), amntClaimed, "amntClaimed")
	s.Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(5)), valFeeSent, "valFeeSent")
	s.Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(500)), totalClaim, "totalClaim")
	s.Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(500)), remainingDeposit, "remainingDeposit")

	// check event emission
	events := tCtx.EventManager().Events()

	hasEvent := false
	for _, ev := range events {
		if ev.Type == types.EventTypeClaimStreamAction {
			hasEvent = true
			attrStreamId, evOk := ev.GetAttribute(types.AttributeKeyStreamId)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamId, attrStreamId.Key)
			s.Equal("1", attrStreamId.Value)

			attrSender, ok := ev.GetAttribute(types.AttributeKeyStreamSender)
			s.True(ok)
			s.Equal(types.AttributeKeyStreamSender, attrSender.Key)
			s.Equal(s.addrs[0].String(), attrSender.Value)

			attrReceiver, ok := ev.GetAttribute(types.AttributeKeyStreamReceiver)
			s.True(ok)
			s.Equal(types.AttributeKeyStreamReceiver, attrReceiver.Key)
			s.Equal(s.addrs[1].String(), attrReceiver.Value)

			attrStreamClaimTotal, evOk := ev.GetAttribute(types.AttributeKeyStreamClaimTotal)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamClaimTotal, attrStreamClaimTotal.Key)
			s.Equal("500stake", attrStreamClaimTotal.Value)

			attrStreamClaimAmountReceived, evOk := ev.GetAttribute(types.AttributeKeyStreamClaimAmountReceived)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamClaimAmountReceived, attrStreamClaimAmountReceived.Key)
			s.Equal("495stake", attrStreamClaimAmountReceived.Value)

			attrStreamClaimValidatorFee, evOk := ev.GetAttribute(types.AttributeKeyStreamClaimValidatorFee)
			s.True(evOk)
			s.Equal(types.AttributeKeyStreamClaimValidatorFee, attrStreamClaimValidatorFee.Key)
			s.Equal("5stake", attrStreamClaimValidatorFee.Value)

			attrStreamRemainingDeposit, evOk := ev.GetAttribute(types.AttributeKeyRemainingDeposit)
			s.True(evOk)
			s.Equal(types.AttributeKeyRemainingDeposit, attrStreamRemainingDeposit.Key)
			s.Equal("500stake", attrStreamRemainingDeposit.Value)
		}
	}

	// should emit stream_deposit event
	s.True(hasEvent)

	// check stream in keeper
	stream, ok := s.app.StreamKeeper.GetStream(tCtx, s.addrs[1], s.addrs[0])
	s.True(ok)
	s.Equal(sdk.NewCoin("stake", sdk.NewIntFromUint64(500)), stream.Deposit, "stream.Deposit")
}

func (s *KeeperTestSuite) TestClaimFromStream_Scenarios() {
	// Todo
}

func (s *KeeperTestSuite) TestClaimFromStream_Fail_NotExist() {
	c, v, t, d, err := s.app.StreamKeeper.ClaimFromStream(s.ctx, s.addrs[1], s.addrs[0])
	s.ErrorContains(err, "stream does not exist")
	s.Equal(sdk.Coin{}, c)
	s.Equal(sdk.Coin{}, v)
	s.Equal(sdk.Coin{}, t)
	s.Equal(sdk.Coin{}, d)

	// double check
	stream, ok := s.app.StreamKeeper.GetStream(s.ctx, s.addrs[1], s.addrs[0])
	s.False(ok)
	s.Equal(types.Stream{}, stream)
}

func (s *KeeperTestSuite) TestClaimFromStream_Fail_NoDeposit() {
	testCases := []struct {
		name   string
		stream types.Stream
	}{
		{
			name: "zero deposit",
			stream: types.Stream{
				StreamId: 1,
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
				Deposit:  sdk.NewCoin("stake", sdk.NewIntFromUint64(0)),
			},
		},
		{
			name: "nil deposit",
			stream: types.Stream{
				StreamId: 2,
				Sender:   s.addrs[2].String(),
				Receiver: s.addrs[3].String(),
				Deposit:  sdk.Coin{},
			},
		},
	}

	for _, tc := range testCases {
		sendAddr, _ := sdk.AccAddressFromBech32(tc.stream.Sender)
		recAddr, _ := sdk.AccAddressFromBech32(tc.stream.Receiver)
		err := s.app.StreamKeeper.SetStream(s.ctx, recAddr, sendAddr, tc.stream)
		s.NoError(err, "SetStream NoError test name %s", tc.name)
		c, v, t, d, err := s.app.StreamKeeper.ClaimFromStream(s.ctx, recAddr, sendAddr)
		s.ErrorContains(err, "stream deposit is zero")
		s.Equal(sdk.Coin{}, c)
		s.Equal(sdk.Coin{}, v)
		s.Equal(sdk.Coin{}, t)
		s.Equal(sdk.Coin{}, d)
	}
}

// Todo - CancelStreamBySenderReceiver
// Todo - GetTotalDeposits
