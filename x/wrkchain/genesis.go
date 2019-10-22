package wrkchain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	WrkChains []WrkChain `json:"registered_wrkchains"`
}

func NewGenesisState(wrkChains []WrkChain) GenesisState {
	return GenesisState{WrkChains: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.WrkChains {
		if record.Owner == nil {
			return fmt.Errorf("Invalid WrkChain: Owner: %s. Error: Missing Owner", record.Owner)
		}
		if record.WrkChainID == "" {
			return fmt.Errorf("Invalid WrkChain: WrkChainID: %s. Error: Missing ID", record.WrkChainID)
		}
		if record.GenesisHash == "" {
			return fmt.Errorf("Invalid WrkChain: GenesisHash: %s. Error: Missing Genesis Hash", record.GenesisHash)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		WrkChains: []WrkChain{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.WrkChains {
		keeper.SetWrkChain(ctx, record.WrkChainID, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []WrkChain
	iterator := k.GetWrkChainsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		wrkchainId := string(iterator.Key())
		wrkChain := k.GetWrkChain(ctx, wrkchainId)
		records = append(records, wrkChain)
	}
	return GenesisState{WrkChains: records}
}
