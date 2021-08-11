package keeper_test

import (
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// WRKChainEqual checks if two WRKChains are equal
func WRKChainEqual(wcA types.WrkChain, wcB types.WrkChain) bool {
	return wcA == wcB
}
//
//// ParamsEqual checks params are equal
//func ParamsEqual(paramsA, paramsB types.Params) bool {
//	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(paramsA),
//		types.ModuleCdc.MustMarshalBinaryBare(paramsB))
//}
//
// WRKChainBlockEqual checks if two WRKChainBlocks are equal
func WRKChainBlockEqual(bA, bB types.WrkChainBlock) bool {
	return bA == bB
}
