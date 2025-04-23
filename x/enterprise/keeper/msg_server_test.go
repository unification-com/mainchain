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
