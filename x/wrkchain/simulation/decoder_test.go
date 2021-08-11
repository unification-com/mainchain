package simulation

//import (
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/tendermint/tendermint/crypto/ed25519"
//	tmkv "github.com/tendermint/tendermint/libs/kv"
//
//	"github.com/cosmos/cosmos-sdk/codec"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/unification-com/mainchain/x/wrkchain/types"
//)
//
//var (
//	wcPk1   = ed25519.GenPrivKey().PubKey()
//	wcAddr1 = sdk.AccAddress(wcPk1.Address())
//)
//
//func makeTestCodec() (cdc *codec.Codec) {
//	cdc = codec.New()
//	sdk.RegisterCodec(cdc)
//	codec.RegisterCrypto(cdc)
//	types.RegisterCodec(cdc)
//	return
//}
//
//func TestDecodeStore(t *testing.T) {
//	cdc := makeTestCodec()
//
//	wrkChain := types.NewWrkChain()
//	wrkChain.WrkChainID = 1
//	wrkChain.Moniker = "WrkChain1"
//	wrkChain.Name = "Test WRKChain 1"
//	wrkChain.GenesisHash = "testgenesissha256hash"
//	wrkChain.LastBlock = 1
//	wrkChain.Owner = wcAddr1
//	wrkChain.RegisterTime = time.Now().Unix()
//
//	wrkChainBlock := types.NewWrkChainBlock()
//	wrkChainBlock.WrkChainID = 1
//	wrkChainBlock.Height = 1
//	wrkChainBlock.Owner = wcAddr1
//	wrkChainBlock.BlockHash = "arbitraryblockhashvalue"
//	wrkChainBlock.ParentHash = "arbitraryparenthashvalue"
//	wrkChainBlock.Hash1 = "hash1"
//	wrkChainBlock.Hash2 = "hash2"
//	wrkChainBlock.Hash3 = "hash3"
//	wrkChainBlock.SubmitTime = time.Now().Unix()
//
//	kvPairs := tmkv.Pairs{
//		tmkv.Pair{Key: types.WrkChainKey(1), Value: cdc.MustMarshalBinaryLengthPrefixed(wrkChain)},
//		tmkv.Pair{Key: types.WrkChainBlockKey(1, 1), Value: cdc.MustMarshalBinaryLengthPrefixed(wrkChainBlock)},
//		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
//	}
//
//	tests := []struct {
//		name        string
//		expectedLog string
//	}{
//		{"wrkchain", fmt.Sprintf("%v\n%v", wrkChain, wrkChain)},
//		{"wrkchain block", fmt.Sprintf("%v\n%v", wrkChainBlock, wrkChainBlock)},
//		{"other", ""},
//	}
//
//	for i, tt := range tests {
//		i, tt := i, tt
//		t.Run(tt.name, func(t *testing.T) {
//			switch i {
//			case len(tests) - 1:
//				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
//			default:
//				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
//			}
//		})
//	}
//}
