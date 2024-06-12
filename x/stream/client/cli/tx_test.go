package cli_test

import (
	"context"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/stream/client/cli"
)

func (s *CLITestSuite) TestCreateStreamTxCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)
	cmd := cli.GetCmdCreateStream()
	cmd.SetOutput(io.Discard)

	extraArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("photon", sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=test-chain", flags.FlagChainID),
	}

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		receiverAddr sdk.AccAddress
		from         string
		deposit      sdk.Coin
		flowRate     string
		extraArgs    []string
		expectErr    bool
	}{
		{
			"valid create",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(1000)),
			"10",
			extraArgs,
			false,
		},
		{
			"invalid receiver address",
			func() client.Context {
				return s.baseCtx
			},
			sdk.AccAddress{},
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(1000)),
			"10",
			extraArgs,
			true,
		},
		{
			"invalid deposit",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.Coin{},
			"10",
			extraArgs,
			true,
		},
		{
			"invalid flow rate",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(1000)),
			"ten",
			extraArgs,
			true,
		},
		{
			"fail msg validate basic - duration too short",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(10)),
			"10",
			extraArgs,
			true,
		},
		{
			"fail msg validate basic - sender receiver same",
			func() client.Context {
				return s.baseCtx
			},
			accounts[0].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(10000)),
			"10",
			extraArgs,
			true,
		},
		{
			"fail msg validate basic - zero deposit",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(0)),
			"10",
			extraArgs,
			true,
		},
		{
			"fail msg validate basic - zero flow rate",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(1000)),
			"0",
			extraArgs,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(append([]string{tc.receiverAddr.String(), tc.deposit.String(), tc.flowRate, tc.from}, tc.extraArgs...))

			s.Require().NoError(client.SetCmdClientContextHandler(tc.ctxGen(), cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestClaimStreamTxCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)
	cmd := cli.GetCmdClaimStream()
	cmd.SetOutput(io.Discard)

	extraArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("photon", sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=test-chain", flags.FlagChainID),
	}

	testCases := []struct {
		name       string
		ctxGen     func() client.Context
		senderAddr sdk.AccAddress
		from       string
		extraArgs  []string
		expectErr  bool
	}{
		{
			"valid claim",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			extraArgs,
			false,
		},
		{
			"invalid sender address",
			func() client.Context {
				return s.baseCtx
			},
			sdk.AccAddress{},
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			extraArgs,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(append([]string{tc.senderAddr.String(), tc.from}, tc.extraArgs...))

			s.Require().NoError(client.SetCmdClientContextHandler(tc.ctxGen(), cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestTopUpDepositTxCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)
	cmd := cli.GetCmdTopUpDeposit()
	cmd.SetOutput(io.Discard)

	extraArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("photon", sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=test-chain", flags.FlagChainID),
	}

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		receiverAddr sdk.AccAddress
		from         string
		deposit      sdk.Coin
		extraArgs    []string
		expectErr    bool
	}{
		{
			"valid top up",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(1000)),
			extraArgs,
			false,
		},
		{
			"invalid receiver address",
			func() client.Context {
				return s.baseCtx
			},
			sdk.AccAddress{},
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(1000)),
			extraArgs,
			true,
		},
		{
			"invalid deposit",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.Coin{},
			extraArgs,
			true,
		},
		{
			"fail msg validate basic - zero deposit",
			func() client.Context {
				return s.baseCtx
			},
			accounts[0].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			sdk.NewCoin("stake", sdk.NewInt(0)),
			extraArgs,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(append([]string{tc.receiverAddr.String(), tc.deposit.String(), tc.from}, tc.extraArgs...))

			s.Require().NoError(client.SetCmdClientContextHandler(tc.ctxGen(), cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestUpdateFlowRateTxCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)
	cmd := cli.GetCmdUpdateFlowRate()
	cmd.SetOutput(io.Discard)

	extraArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("photon", sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=test-chain", flags.FlagChainID),
	}

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		receiverAddr sdk.AccAddress
		from         string
		flowRate     string
		extraArgs    []string
		expectErr    bool
	}{
		{
			"valid update",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			"10",
			extraArgs,
			false,
		},
		{
			"invalid receiver address",
			func() client.Context {
				return s.baseCtx
			},
			sdk.AccAddress{},
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			"10",
			extraArgs,
			true,
		},
		{
			"invalid flow rate",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			"ten",
			extraArgs,
			true,
		},
		{
			"fail msg validate basic - zero flow rate",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			"0",
			extraArgs,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(append([]string{tc.receiverAddr.String(), tc.flowRate, tc.from}, tc.extraArgs...))

			s.Require().NoError(client.SetCmdClientContextHandler(tc.ctxGen(), cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *CLITestSuite) TestCancelStreamTxCmd() {
	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)
	cmd := cli.GetCmdCancelStream()
	cmd.SetOutput(io.Discard)

	extraArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("photon", sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=test-chain", flags.FlagChainID),
	}

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		receiverAddr sdk.AccAddress
		from         string
		extraArgs    []string
		expectErr    bool
	}{
		{
			"valid cancel",
			func() client.Context {
				return s.baseCtx
			},
			accounts[1].Address,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			extraArgs,
			false,
		},
		{
			"invalid receiver address",
			func() client.Context {
				return s.baseCtx
			},
			sdk.AccAddress{},
			fmt.Sprintf("--%s=%s", flags.FlagFrom, "key-0"),
			extraArgs,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(append([]string{tc.receiverAddr.String(), tc.from}, tc.extraArgs...))

			s.Require().NoError(client.SetCmdClientContextHandler(tc.ctxGen(), cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
