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

func (s *KeeperTestSuite) TestRegisterWrkChain() {
	testCases := []struct {
		name        string
		request     *types.MsgRegisterWrkChain
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid owner",
			request: &types.MsgRegisterWrkChain{
				Owner: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "name too big",
			request: &types.MsgRegisterWrkChain{
				Owner: s.addrs[0].String(),
				Name:  "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "name too big",
		},
		{
			name: "moniker too big",
			request: &types.MsgRegisterWrkChain{
				Owner:   s.addrs[0].String(),
				Name:    "testname",
				Moniker: "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "moniker too big",
		},
		{
			name: "moniker empty",
			request: &types.MsgRegisterWrkChain{
				Owner:   s.addrs[0].String(),
				Name:    "testname",
				Moniker: "",
			},
			expectErr:   true,
			expectedErr: "must have a moniker",
		},
		{
			name: "valid registration",
			request: &types.MsgRegisterWrkChain{
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
			_, err := s.msgServer.RegisterWrkChain(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRecordWrkChainBlock() {

	_, _ = s.msgServer.RegisterWrkChain(s.ctx, &types.MsgRegisterWrkChain{
		Owner:   s.addrs[0].String(),
		Name:    "testname",
		Moniker: "testmoniker",
	})

	testCases := []struct {
		name        string
		request     *types.MsgRecordWrkChainBlock
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid owner",
			request: &types.MsgRecordWrkChainBlock{
				Owner: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "zero height",
			request: &types.MsgRecordWrkChainBlock{
				Owner:  s.addrs[0].String(),
				Height: 0,
			},
			expectErr:   true,
			expectedErr: "height must be > 0",
		},
		{
			name: "block hash too big",
			request: &types.MsgRecordWrkChainBlock{
				Owner:     s.addrs[0].String(),
				Height:    1,
				BlockHash: "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "block hash too big",
		},
		{
			name: "parent hash too big",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[0].String(),
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "parent hash too big",
		},
		{
			name: "hash1 too big",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[0].String(),
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "parenthash",
				Hash1:      "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "hash1 too big",
		},
		{
			name: "hash2 too big",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[0].String(),
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "parenthash",
				Hash1:      "hash1",
				Hash2:      "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "hash2 too big",
		},
		{
			name: "hash3 too big",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[0].String(),
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "parenthash",
				Hash1:      "hash1",
				Hash2:      "hash2",
				Hash3:      "lkwnefokwenfkowenfkowenfkowenokeokfnklenwfklwnlkewnlenflkwflkwnlknwelkfnwelfnlkweflkwflkwnflkwnflknelkfwlkflkwnflkwnflkwnlkwnwefwes",
			},
			expectErr:   true,
			expectedErr: "hash3 too big",
		},
		{
			name: "wrkchain not registered",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[0].String(),
				WrkchainId: 99,
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "parenthash",
				Hash1:      "hash1",
				Hash2:      "hash2",
				Hash3:      "hash3",
			},
			expectErr:   true,
			expectedErr: "wrkchain has not been registered yet",
		},
		{
			name: "not owner",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[1].String(),
				WrkchainId: 1,
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "parenthash",
				Hash1:      "hash1",
				Hash2:      "hash2",
				Hash3:      "hash3",
			},
			expectErr:   true,
			expectedErr: "you are not the owner of this wrkchain",
		},
		{
			name: "valid",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[0].String(),
				WrkchainId: 1,
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "parenthash",
				Hash1:      "hash1",
				Hash2:      "hash2",
				Hash3:      "hash3",
			},
			expectErr:   false,
			expectedErr: "",
		},
		{
			name: "height must be higher than previous",
			request: &types.MsgRecordWrkChainBlock{
				Owner:      s.addrs[0].String(),
				WrkchainId: 1,
				Height:     1,
				BlockHash:  "blockhash",
				ParentHash: "parenthash",
				Hash1:      "hash1",
				Hash2:      "hash2",
				Hash3:      "hash3",
			},
			expectErr:   true,
			expectedErr: "wrkchain block hashes height must be > last height recorded",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.RecordWrkChainBlock(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPurchaseWrkChainStateStorage() {

	_, _ = s.msgServer.RegisterWrkChain(s.ctx, &types.MsgRegisterWrkChain{
		Owner:   s.addrs[0].String(),
		Name:    "testname",
		Moniker: "testmoniker",
	})

	testCases := []struct {
		name        string
		request     *types.MsgPurchaseWrkChainStateStorage
		expectErr   bool
		expectedErr string
	}{
		{
			name: "set invalid owner",
			request: &types.MsgPurchaseWrkChainStateStorage{
				Owner: "invalidaddr",
			},
			expectErr:   true,
			expectedErr: "decoding bech32 failed",
		},
		{
			name: "wrkchain id zero",
			request: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      s.addrs[0].String(),
				WrkchainId: 0,
			},
			expectErr:   true,
			expectedErr: "id must be greater than zero",
		},
		{
			name: "number zero",
			request: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      s.addrs[0].String(),
				WrkchainId: 1,
				Number:     0,
			},
			expectErr:   true,
			expectedErr: "cannot purchase zero",
		},
		{
			name: "wrkchain not registered",
			request: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      s.addrs[0].String(),
				WrkchainId: 99,
				Number:     100,
			},
			expectErr:   true,
			expectedErr: "wrkchain has not been registered yet",
		},
		{
			name: "not owner",
			request: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      s.addrs[1].String(),
				WrkchainId: 1,
				Number:     100,
			},
			expectErr:   true,
			expectedErr: "you are not the owner of this wrkchain",
		},
		{
			name: "purchase ok",
			request: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      s.addrs[0].String(),
				WrkchainId: 1,
				Number:     100,
			},
			expectErr:   false,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.PurchaseWrkChainStateStorage(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
