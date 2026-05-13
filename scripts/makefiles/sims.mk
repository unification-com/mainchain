########################################
### Simulations
#
# All targets run against the v0.54 testutil/simsx runner via `go test -tags=sims`.
# The `runsim` binary was retired upstream in v0.53; nothing here should shell out
# to it. Multi-seed coverage is provided by simsx.RunWithSeeds internally — there
# is no need for an external job runner.

SIMAPP = ./app/...

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 50
SIM_COMMIT ?= true

# Hand-picked seeds known to produce non-pathological sim trajectories at deep
# NumBlocks (carried over from the pre-simsx runsim invocation). Default seeds
# from simsx.Run / rand.Int63 hit "all validators jailed" too often at 500
# blocks, which is why test-sim-simple uses single-seed and test-sim-multi-seed-long
# iterates this curated list explicitly.
SIM_SEEDS = 1 2 4 7 32 123 124 582 1893 2989 3012 4728 37827 981928 87821 \
            891823782 989182 89182392 11 22 44 77 99 2020 3232 123123 124124 \
            582582 18931893 29892989 30123012 47284728 7601778 8090485 \
            977367484 491163361 424254581 673398983

# Parallelism for the seed-iterating targets. 4 matches the historical runsim -Jobs=4.
SIM_JOBS ?= 4

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.und_mainchain/config/genesis.json will be used."
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestFullAppSimulation \
		-Genesis=${HOME}/.und_mainchain/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

test-sim-import-export:
	@echo "Running application import/export simulation. This may take several minutes..."
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestAppImportExport -Enabled=true \
		-NumBlocks=50 -Period=5 -v -timeout 24h

test-sim-after-import:
	@echo "Running application simulation-after-import. This may take several minutes..."
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestAppSimulationAfterImport -Enabled=true \
		-NumBlocks=50 -Period=5 -v -timeout 24h

test-sim-custom-genesis-multi-seed:
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.und_mainchain/config/genesis.json will be used."
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestFullAppSimulation \
		-Genesis=${HOME}/.und_mainchain/config/genesis.json \
		-Enabled=true -NumBlocks=400 -Period=5 -v -timeout 24h

# Multi-seed determinism (TestAppStateDeterminism's 3 seeds × 3 runs per seed).
# Quick smoke at 50 blocks — sufficient because the random seeds it generates
# rarely produce empty-validator-set states at this depth.
test-sim-multi-seed-short:
	@echo "Running short multi-seed determinism check."
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=50 -BlockSize=10 -Commit=true -Period=0 -v -timeout 24h

# Deep multi-seed simulation: each curated seed in SIM_SEEDS gets a full
# TestFullAppSimulation run at 500 blocks, $(SIM_JOBS) at a time. The seed
# list is hand-picked to avoid the empty-validator-set pathology that random
# seeds hit at this depth. Matches the historical runsim -Jobs=4 shape.
test-sim-multi-seed-long:
	@echo "Running long multi-seed simulation: $(words $(SIM_SEEDS)) seeds × 500 blocks, $(SIM_JOBS) at a time..."
	@printf '%s\n' $(SIM_SEEDS) | xargs -I{} -P $(SIM_JOBS) -n 1 sh -c '\
		echo "=== seed={} ==="; \
		go test -mod=readonly -tags=sims $(SIMAPP) -run TestFullAppSimulation -Enabled=true \
			-Seed={} -NumBlocks=500 -BlockSize=50 -Commit=true -Period=50 -timeout 24h'

# Benchmark pins seed 42 — empirically completes 500 blocks × 50 BlockSize
# without hitting "empty validator set" in the benchmark harness. Note that
# seed-stability differs between TestFullAppSimulation (all 38 SIM_SEEDS pass at
# 500×50) and BenchmarkFullAppSimulation (some seeds skip on the bench harness
# even though they pass on the test harness — same simsx internals, different
# testing.TB type). Override via `make SIM_BENCH_SEED=N test-sim-benchmark`.
SIM_BENCH_SEED ?= 42

test-sim-benchmark:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE), seed=$(SIM_BENCH_SEED). This may take awhile!"
	@go test -mod=readonly -tags=sims -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ -benchtime=1x \
		-Enabled=true -Seed=$(SIM_BENCH_SEED) -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

test-sim-profile:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE), seed=$(SIM_BENCH_SEED). This may take awhile!"
	@go test -mod=readonly -tags=sims -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ -benchtime=1x \
		-Enabled=true -Seed=$(SIM_BENCH_SEED) -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h \
		-cpuprofile cpu.out -memprofile mem.out

# Single-seed deep simulation. TestFullAppSimulation is single-seed by design;
# vary the seed with -Seed=N or use test-sim-multi-seed-* for multi-seed signal.
test-sim-simple:
	@echo "Running simple single-seed simulation (NumBlocks=$(SIM_NUM_BLOCKS), BlockSize=$(SIM_BLOCK_SIZE))..."
	@mkdir -p $(HOME)/.und_simapp
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestFullAppSimulation -Enabled=true \
		-NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -Period=100 -v -timeout 24h -Verbose=true \
		-ExportParamsPath=$(HOME)/.und_simapp/params.json \
		-ExportStatePath=$(HOME)/.und_simapp/state.json \
		-ExportStatsPath=$(HOME)/.und_simapp/statistics.json

SIM_FUZZ_TIME ?= 5m

test-sim-fuzz:
	@echo "Running application fuzz for numBlocks=2, blockSize=20 (fuzztime=$(SIM_FUZZ_TIME))."
	@go test -mod=readonly -tags=sims -timeout=60m -fuzztime=$(SIM_FUZZ_TIME) \
		-run=^$$ -fuzz=^FuzzFullAppSimulation$$ $(SIMAPP) \
		-NumBlocks=2 -BlockSize=20

.PHONY: \
	test-sim-nondeterminism \
	test-sim-custom-genesis-fast \
	test-sim-import-export \
	test-sim-after-import \
	test-sim-custom-genesis-multi-seed \
	test-sim-multi-seed-short \
	test-sim-multi-seed-long \
	test-sim-benchmark \
	test-sim-profile \
	test-sim-simple \
	test-sim-fuzz
