package keeper_test

import (
	"github.com/unification-com/mainchain/x/beacon/types"
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
				Authority: s.app.BeaconKeeper.GetAuthority(),
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
				Authority: s.app.BeaconKeeper.GetAuthority(),
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

func (s *KeeperTestSuite) TestRegisterBeacon() {
	testCases := []struct {
		name        string
		request     *types.MsgRegisterBeacon
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid owner",
			request: &types.MsgRegisterBeacon{
				Owner: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "name too big",
			request: &types.MsgRegisterBeacon{
				Owner: s.addrs[0].String(),
				Name:  "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "name too big",
		},
		{
			name: "moniker too big",
			request: &types.MsgRegisterBeacon{
				Owner:   s.addrs[0].String(),
				Name:    "testname",
				Moniker: "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "moniker too big",
		},
		{
			name: "moniker empty",
			request: &types.MsgRegisterBeacon{
				Owner:   s.addrs[0].String(),
				Name:    "testname",
				Moniker: "",
			},
			expectErr:   true,
			expectedErr: "must have a moniker",
		},
		{
			name: "valid registration",
			request: &types.MsgRegisterBeacon{
				Owner:   s.addrs[0].String(),
				Name:    "testname",
				Moniker: "testmoniker",
			},
			expectErr:   false,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.RegisterBeacon(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRecordBeaconTimestamp() {

	// register test BEACON
	_, _ = s.msgServer.RegisterBeacon(s.ctx, &types.MsgRegisterBeacon{
		Owner:   s.addrs[0].String(),
		Name:    "testname",
		Moniker: "testmoniker",
	})

	testCases := []struct {
		name        string
		request     *types.MsgRecordBeaconTimestamp
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid owner",
			request: &types.MsgRecordBeaconTimestamp{
				Owner: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "beacon id zero",
			request: &types.MsgRecordBeaconTimestamp{
				Owner:    s.addrs[0].String(),
				BeaconId: 0,
			},
			expectErr:   true,
			expectedErr: "id must be greater than zero",
		},
		{
			name: "empty hash",
			request: &types.MsgRecordBeaconTimestamp{
				Owner:    s.addrs[0].String(),
				BeaconId: 1,
				Hash:     "",
			},
			expectErr:   true,
			expectedErr: "hash cannot be empty",
		},
		{
			name: "hash too big",
			request: &types.MsgRecordBeaconTimestamp{
				Owner:    s.addrs[0].String(),
				BeaconId: 1,
				Hash:     "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "hash too big",
		},
		{
			name: "beacon does not exist",
			request: &types.MsgRecordBeaconTimestamp{
				Owner:    s.addrs[0].String(),
				BeaconId: 99,
				Hash:     "testhash",
			},
			expectErr:   true,
			expectedErr: "beacon has not been registered yet",
		},
		{
			name: "not owner",
			request: &types.MsgRecordBeaconTimestamp{
				Owner:    s.addrs[1].String(),
				BeaconId: 1,
				Hash:     "testhash",
			},
			expectErr:   true,
			expectedErr: "you are not the owner of this beacon",
		},
		{
			name: "valid recording",
			request: &types.MsgRecordBeaconTimestamp{
				Owner:    s.addrs[0].String(),
				BeaconId: 1,
				Hash:     "testhash",
			},
			expectErr:   false,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.RecordBeaconTimestamp(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPurchaseBeaconStateStorage() {
	// register test BEACON
	_, _ = s.msgServer.RegisterBeacon(s.ctx, &types.MsgRegisterBeacon{
		Owner:   s.addrs[0].String(),
		Name:    "testname",
		Moniker: "testmoniker",
	})

	testCases := []struct {
		name        string
		request     *types.MsgPurchaseBeaconStateStorage
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid owner",
			request: &types.MsgPurchaseBeaconStateStorage{
				Owner: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "beacon id zero",
			request: &types.MsgPurchaseBeaconStateStorage{
				Owner:    s.addrs[0].String(),
				BeaconId: 0,
			},
			expectErr:   true,
			expectedErr: "id must be greater than zero",
		},
		{
			name: "number zero",
			request: &types.MsgPurchaseBeaconStateStorage{
				Owner:    s.addrs[0].String(),
				BeaconId: 1,
				Number:   0,
			},
			expectErr:   true,
			expectedErr: "cannot purchase zero",
		},
		{
			name: "beacon not registered",
			request: &types.MsgPurchaseBeaconStateStorage{
				Owner:    s.addrs[0].String(),
				BeaconId: 99,
				Number:   100,
			},
			expectErr:   true,
			expectedErr: "beacon has not been registered yet",
		},
		{
			name: "not owner",
			request: &types.MsgPurchaseBeaconStateStorage{
				Owner:    s.addrs[1].String(),
				BeaconId: 1,
				Number:   100,
			},
			expectErr:   true,
			expectedErr: "you are not the owner of this beacon",
		},
		{
			name: "purchase ok",
			request: &types.MsgPurchaseBeaconStateStorage{
				Owner:    s.addrs[0].String(),
				BeaconId: 1,
				Number:   100,
			},
			expectErr:   false,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.PurchaseBeaconStateStorage(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
