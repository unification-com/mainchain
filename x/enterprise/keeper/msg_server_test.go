package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/enterprise/types"
)

func (s *KeeperTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr bool
	}{
		{
			name: "set invalid authority",
			request: &types.MsgUpdateParams{
				Authority: "foo",
			},
			expectErr: true,
		},
		{
			name: "set invalid params",
			request: &types.MsgUpdateParams{
				Authority: s.app.EnterpriseKeeper.GetAuthority(),
				Params: types.Params{
					EntSigners:        "",
					Denom:             "",
					MinAccepts:        0,
					DecisionTimeLimit: 0,
				},
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			request: &types.MsgUpdateParams{
				Authority: s.app.EnterpriseKeeper.GetAuthority(),
				Params: types.Params{
					EntSigners:        "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6,und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed",
					Denom:             sdk.DefaultBondDenom,
					MinAccepts:        2,
					DecisionTimeLimit: 2,
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateParams(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUndPurchaseOrder() {

	entSigner := s.addrs[0]
	whitelisted := s.addrs[1]
	notWhitelisted := s.addrs[2]

	_ = s.app.EnterpriseKeeper.SetParams(s.ctx, types.Params{
		EntSigners:        entSigner.String(),
		Denom:             sdk.DefaultBondDenom,
		MinAccepts:        1,
		DecisionTimeLimit: 999999999,
	})

	_ = s.app.EnterpriseKeeper.AddAddressToWhitelist(s.ctx, whitelisted)

	testCases := []struct {
		name        string
		request     *types.MsgUndPurchaseOrder
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid purchaser",
			request: &types.MsgUndPurchaseOrder{
				Purchaser: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "incorrect denomination",
			request: &types.MsgUndPurchaseOrder{
				Purchaser: whitelisted.String(),
				Amount:    sdk.NewInt64Coin("rubbish", 1),
			},
			expectErr:   true,
			expectedErr: "denomination must be nund",
		},
		{
			name: "amount zero",
			request: &types.MsgUndPurchaseOrder{
				Purchaser: whitelisted.String(),
				Amount:    sdk.NewInt64Coin("nund", 0),
			},
			expectErr:   true,
			expectedErr: "amount must be greater than zero",
		},
		{
			name: "not whitelisted",
			request: &types.MsgUndPurchaseOrder{
				Purchaser: notWhitelisted.String(),
				Amount:    sdk.NewInt64Coin("nund", 100),
			},
			expectErr:   true,
			expectedErr: "not whitelisted to raise purchase orders",
		},
		{
			name: "valid",
			request: &types.MsgUndPurchaseOrder{
				Purchaser: whitelisted.String(),
				Amount:    sdk.NewInt64Coin("nund", 100),
			},
			expectErr:   false,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.UndPurchaseOrder(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestProcessUndPurchaseOrder() {
	entSigner := s.addrs[0]

	_ = s.app.EnterpriseKeeper.SetParams(s.ctx, types.Params{
		EntSigners:        entSigner.String(),
		Denom:             sdk.DefaultBondDenom,
		MinAccepts:        1,
		DecisionTimeLimit: 999999999,
	})

	_ = s.app.EnterpriseKeeper.SetPurchaseOrder(s.ctx, types.EnterpriseUndPurchaseOrder{
		Id:             1,
		Purchaser:      s.addrs[1].String(),
		Amount:         sdk.NewInt64Coin("nund", 100),
		Status:         types.StatusRaised,
		RaiseTime:      0,
		CompletionTime: 0,
		Decisions:      nil,
	})

	_ = s.app.EnterpriseKeeper.SetPurchaseOrder(s.ctx, types.EnterpriseUndPurchaseOrder{
		Id:             2,
		Purchaser:      s.addrs[1].String(),
		Amount:         sdk.NewInt64Coin("nund", 100),
		Status:         types.StatusCompleted,
		RaiseTime:      0,
		CompletionTime: 0,
		Decisions:      nil,
	})

	testCases := []struct {
		name        string
		request     *types.MsgProcessUndPurchaseOrder
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid signer",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "unauthorised signer",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer: s.addrs[1].String(),
			},
			expectErr:   true,
			expectedErr: "unauthorised signer processing purchase order",
		},
		{
			name: "po id zero",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSigner.String(),
				PurchaseOrderId: 0,
			},
			expectErr:   true,
			expectedErr: "purchase order id must be greater than zero",
		},
		{
			name: "po not exist",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSigner.String(),
				PurchaseOrderId: 99,
			},
			expectErr:   true,
			expectedErr: "purchase order does not exist",
		},
		{
			name: "invalid decision",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSigner.String(),
				PurchaseOrderId: 1,
				Decision:        types.StatusCompleted,
			},
			expectErr:   true,
			expectedErr: "decision should be accept or reject",
		},
		{
			name: "already processed",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSigner.String(),
				PurchaseOrderId: 2,
				Decision:        types.StatusAccepted,
			},
			expectErr:   true,
			expectedErr: "already processed",
		},
		{
			name: "valid",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSigner.String(),
				PurchaseOrderId: 1,
				Decision:        types.StatusAccepted,
			},
			expectErr:   false,
			expectedErr: "",
		},
		{
			name: "already decided",
			request: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSigner.String(),
				PurchaseOrderId: 1,
				Decision:        types.StatusAccepted,
			},
			expectErr:   true,
			expectedErr: "already decided",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.ProcessUndPurchaseOrder(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestWhitelistAddress() {
	entSigner := s.addrs[0]

	_ = s.app.EnterpriseKeeper.SetParams(s.ctx, types.Params{
		EntSigners:        entSigner.String(),
		Denom:             sdk.DefaultBondDenom,
		MinAccepts:        1,
		DecisionTimeLimit: 999999999,
	})

	testCases := []struct {
		name        string
		request     *types.MsgWhitelistAddress
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid signer",
			request: &types.MsgWhitelistAddress{
				Signer: "bored_with_unit_tests",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "set invalid address",
			request: &types.MsgWhitelistAddress{
				Signer:  entSigner.String(),
				Address: "yarbles",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "unauthorised",
			request: &types.MsgWhitelistAddress{
				Signer:  s.addrs[2].String(),
				Address: s.addrs[1].String(),
			},
			expectErr:   true,
			expectedErr: "unauthorised signer modifying whitelist",
		},
		{
			name: "invalid action",
			request: &types.MsgWhitelistAddress{
				Signer:  entSigner.String(),
				Address: s.addrs[1].String(),
				Action:  types.WhitelistActionNil,
			},
			expectErr:   true,
			expectedErr: "action should be add or remove",
		},
		{
			name: "ok",
			request: &types.MsgWhitelistAddress{
				Signer:  entSigner.String(),
				Address: s.addrs[1].String(),
				Action:  types.WhitelistActionAdd,
			},
			expectErr:   false,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.WhitelistAddress(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
