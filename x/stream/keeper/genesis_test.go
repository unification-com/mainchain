package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestImportExportGenesis() {
	tCtx := s.ctx
	nowTime := time.Unix(time.Now().Unix(), 0).UTC()

	for i := int64(1); i < 100; i++ {
		deposit := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000*i)
		createTime := time.Unix(nowTime.Unix()-(i*9), 0).UTC()
		tCtx = tCtx.WithBlockTime(createTime)
		_, err := s.app.StreamKeeper.CreateNewStream(tCtx, s.addrs[i-1], s.addrs[i], deposit, i)
		s.Require().NoError(err)
		_, err = s.app.StreamKeeper.AddDeposit(tCtx, s.addrs[i-1], s.addrs[i], deposit)
		s.Require().NoError(err)

		// simulate some claims etc.
		tCtx = tCtx.WithBlockTime(nowTime)
		_, _, _, _, err = s.app.StreamKeeper.ClaimFromStream(tCtx, s.addrs[i-1], s.addrs[i])
		s.Require().NoError(err)
	}

	tCtx = tCtx.WithBlockTime(nowTime)
	genesis := s.app.StreamKeeper.ExportGenesis(tCtx)
	s.app.StreamKeeper.InitGenesis(s.ctx, genesis)
	newGenesis := s.app.StreamKeeper.ExportGenesis(tCtx)
	s.Require().Equal(genesis, newGenesis)
}
