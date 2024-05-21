package keeper_test

import (
	"github.com/unification-com/mainchain/x/wrkchain/types"
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
				Authority: s.app.WrkchainKeeper.GetAuthority(),
				Params: types.Params{
					FeeRegister:         0,
					FeeRecord:           0,
					FeePurchaseStorage:  0,
					Denom:               "",
					DefaultStorageLimit: 0,
					MaxStorageLimit:     0,
				},
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			request: &types.MsgUpdateParams{
				Authority: s.app.WrkchainKeeper.GetAuthority(),
				Params: types.Params{
					FeeRegister:         24,
					FeeRecord:           2,
					FeePurchaseStorage:  24,
					Denom:               "test",
					DefaultStorageLimit: 99,
					MaxStorageLimit:     999,
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
