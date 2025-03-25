package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"io"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/gogoproto/proto"

	"github.com/unification-com/mainchain/x/stream/client/cli"
	"github.com/unification-com/mainchain/x/stream/types"
)

func (s *CLITestSuite) TestGetAllStreamsCmd() {

	cmd := cli.GetCmdGetAllStreams()
	cmd.SetOutput(io.Discard)

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"valid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryStreamsResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryStreamsResponse{},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			s.Require().NoError(client.SetCmdClientContextHandler(clientCtx, cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(outBuf.Bytes(), tc.expectResult))
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestGetAllStreamsByReceiverCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 1)

	cmd := cli.GetCmdGetAllStreamsByReceiver()
	cmd.SetOutput(io.Discard)

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"valid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryAllStreamsForReceiverResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				accounts[0].Address.String(),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryAllStreamsForReceiverResponse{},
			false,
		},
		{
			"invalid receiver address",
			func() client.Context {
				return s.baseCtx
			},
			[]string{
				"foo",
			},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			s.Require().NoError(client.SetCmdClientContextHandler(clientCtx, cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(outBuf.Bytes(), tc.expectResult))
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestGetCmdGetStreamCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)

	cmd := cli.GetCmdGetStream()
	cmd.SetOutput(io.Discard)

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"valid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryStreamByReceiverSenderResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				accounts[0].Address.String(),
				accounts[1].Address.String(),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryStreamByReceiverSenderResponse{},
			false,
		},
		{
			"invalid receiver address",
			func() client.Context {
				return s.baseCtx
			},
			[]string{
				"foo",
				accounts[1].Address.String(),
			},
			nil,
			true,
		},
		{
			"invalid sender address",
			func() client.Context {
				return s.baseCtx
			},
			[]string{
				accounts[1].Address.String(),
				"foo",
			},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			s.Require().NoError(client.SetCmdClientContextHandler(clientCtx, cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(outBuf.Bytes(), tc.expectResult))
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestGetCmdGetStreamByReceiverSenderCurrentFlowCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)

	cmd := cli.GetCmdGetStreamByReceiverSenderCurrentFlow()
	cmd.SetOutput(io.Discard)

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"valid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryStreamReceiverSenderCurrentFlowResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				accounts[0].Address.String(),
				accounts[1].Address.String(),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryStreamReceiverSenderCurrentFlowResponse{},
			false,
		},
		{
			"invalid receiver address",
			func() client.Context {
				return s.baseCtx
			},
			[]string{
				"foo",
				accounts[1].Address.String(),
			},
			nil,
			true,
		},
		{
			"invalid sender address",
			func() client.Context {
				return s.baseCtx
			},
			[]string{
				accounts[1].Address.String(),
				"foo",
			},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			s.Require().NoError(client.SetCmdClientContextHandler(clientCtx, cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(outBuf.Bytes(), tc.expectResult))
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestGetCmdGetAllStreamsBySenderCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 1)

	cmd := cli.GetCmdGetAllStreamsBySender()
	cmd.SetOutput(io.Discard)

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"valid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryAllStreamsForSenderResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				accounts[0].Address.String(),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryAllStreamsForSenderResponse{},
			false,
		},
		{
			"invalid sender address",
			func() client.Context {
				return s.baseCtx
			},
			[]string{
				"foo",
			},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			s.Require().NoError(client.SetCmdClientContextHandler(clientCtx, cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(outBuf.Bytes(), tc.expectResult))
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestGetCmdCalculateFlowRateCmd() {
	cmd := cli.GetCmdCalculateFlowRate()
	cmd.SetOutput(io.Discard)

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"valid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryCalculateFlowRateResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=month", cli.FlagStreamPeriod),
				fmt.Sprintf("--%s=1000000000stake", cli.FlagStreamCoin),
				fmt.Sprintf("--%s=1", cli.FlagStreamDuration),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryCalculateFlowRateResponse{},
			false,
		},
		{
			"invalid period",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryCalculateFlowRateResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=foo", cli.FlagStreamPeriod),
				fmt.Sprintf("--%s=1000000000stake", cli.FlagStreamCoin),
				fmt.Sprintf("--%s=1", cli.FlagStreamDuration),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			nil,
			true,
		},
		{
			"invalid coin",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryCalculateFlowRateResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=month", cli.FlagStreamPeriod),
				fmt.Sprintf("--%s=foo", cli.FlagStreamCoin),
				fmt.Sprintf("--%s=1", cli.FlagStreamDuration),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			nil,
			true,
		},
		{
			"invalid duration",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryCalculateFlowRateResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=month", cli.FlagStreamPeriod),
				fmt.Sprintf("--%s=1000000000stake", cli.FlagStreamCoin),
				fmt.Sprintf("--%s=0", cli.FlagStreamDuration),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			s.Require().NoError(client.SetCmdClientContextHandler(clientCtx, cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(outBuf.Bytes(), tc.expectResult))
				s.Require().NoError(err)
			}
		})
	}
}
