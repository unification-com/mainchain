########################################
### Simulations
#
# All targets run against the v0.54 testutil/simsx runner via `go test -tags=sims`.
# The `runsim` binary was retired upstream in v0.53; nothing here should shell out
# to it. Multi-seed coverage is provided by simsx.RunWithSeeds internally — there
# is no need for an external job runner.

SIMAPP = ./app/...

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true

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

test-sim-multi-seed-long:
	@echo "Running long multi-seed application simulation. This may take awhile!"
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestFullAppSimulation -Enabled=true \
		-NumBlocks=500 -Period=50 -v -timeout 24h

test-sim-multi-seed-short:
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@go test -mod=readonly -tags=sims $(SIMAPP) -run TestFullAppSimulation -Enabled=true \
		-NumBlocks=50 -Period=10 -v -timeout 24h

test-sim-benchmark:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -tags=sims -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

test-sim-profile:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -tags=sims -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h \
		-cpuprofile cpu.out -memprofile mem.out

test-sim-simple:
	@echo "Running simple test..."
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
	test-sim-multi-seed-long \
	test-sim-multi-seed-short \
	test-sim-benchmark \
	test-sim-profile \
	test-sim-simple \
	test-sim-fuzz
