//go:build sims

package app

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/simsx"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
)

// BenchmarkFullAppSimulation runs the simulator under the benchmark framework
// to capture wall time, allocations, and memory for a single configured seed.
// b.N controls how many simulations are run per benchmark call; with realistic
// NumBlocks values the testing framework usually keeps b.N at 1.
//
// Profile with:
//
//	go test -benchmem -tags=sims -run=^$ -bench=^BenchmarkFullAppSimulation$ \
//	    ./app/... -NumBlocks=500 -BlockSize=200 -Commit=true \
//	    -cpuprofile cpu.out -memprofile mem.out -timeout 24h
func BenchmarkFullAppSimulation(b *testing.B) {
	b.ReportAllocs()

	cfg := simcli.NewConfigFromFlags()
	cfg.ChainID = simsx.SimAppChainID

	for i := 0; i < b.N; i++ {
		simsx.RunWithSeed(b, cfg, NewApp, setupStateFactory, cfg.Seed, nil)
	}
}
