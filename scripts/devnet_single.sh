#!/bin/bash

#########################################################
# A script for quickly spinning up a single node        #
# local devnet using the specified released version     #
# of und                                                #
#                                                       #
# Usage:                                                #
#   ./scripts/devnet_single.sh <UND_VERS>               #
#                                                       #
# Example:                                              #
#   ./scripts/devnet_single.sh 1.6.2                    #
#                                                       #
# binary:   $HOME/und_devnets/<UND_VERS>/und            #
# und home: $HOME/und_devnets/<UND_VERS>/.und_mainchain #
# chain-id: FUND-DevNet                                 #
#                                                       #
# P2P:      tcp://localhost:16656                       #
# RPC:      (http || tcp)://localhost:16657             #
# REST:     http://localhost:1317/swagger/              #
#                                                       #
# gRPC:     localhost:9090                              #
# gRPC Web: http://localhost:9091                       #
#                                                       #
#########################################################

set -eu

# Testnet configuration
BASE_TEST_PATH="${HOME}/und_devnets"
CHAIN_ID="FUND-DevNet"
DEFAULT_UND_VERS="1.6.2"

UND_VERS="${1}"

if [ -z "$UND_VERS" ]; then
  echo "UND_VERS not set. Using und v${DEFAULT_UND_VERS}"
  UND_VERS="${DEFAULT_UND_VERS}"
  sleep 1s
fi

# Set the test path
TEST_PATH="${BASE_TEST_PATH}/${UND_VERS}"

# Internal VARS
DATA_DIR="${TEST_PATH}/.und_mainchain"
UND_BIN="${TEST_PATH}/und"
PREFIX="v"

function version_lt() { test "$(echo "$@" | tr " " "\n" | sort -rV | head -n 1)" != "$1"; }

if version_lt $UND_VERS "1.6.1"; then
  PREFIX=""
fi

if version_lt $UND_VERS "1.5.0"; then
  echo "versions < 1.5.0 not supported"
  exit 1
fi

# check for previous tests
if [ -d "$TEST_PATH" ]; then
  echo ""
  echo "Found previous test configuration in ${TEST_PATH}."
  echo ""
  echo "Either delete ${TEST_PATH} and rerun this script"
  echo "or start the chain again using:"
  echo ""
  echo "  ${UND_BIN} start --home ${DATA_DIR}"
  echo ""
  exit 0
fi

mkdir -p "${TEST_PATH}"

cd "${TEST_PATH}" || exit

# download & unpack und release
wget "https://github.com/unification-com/mainchain/releases/download/${PREFIX}${UND_VERS}/und_v${UND_VERS}_linux_x86_64.tar.gz"
tar -zxvf "und_v${UND_VERS}_linux_x86_64.tar.gz"

# init chain
"${UND_BIN}" init devnet-validator --home "${DATA_DIR}" --chain-id="${CHAIN_ID}"
"${UND_BIN}" config chain-id "${CHAIN_ID}" --home "${DATA_DIR}"
"${UND_BIN}" config keyring-backend test --home "${DATA_DIR}"

# change default denoms from stake to nund in genesis
sed -i "s/stake/nund/g" "${DATA_DIR}/config/genesis.json"

# add test keys to keychain
# add a validator
"${UND_BIN}" --home ${DATA_DIR} keys add validator --keyring-backend test
# add a test key for account 0. If v < 1.6.x, do not use account 0 in testing.
"${UND_BIN}" --home ${DATA_DIR} keys add t0 --keyring-backend test
# add a test key
"${UND_BIN}" --home ${DATA_DIR} keys add t1 --keyring-backend test
"${UND_BIN}" --home ${DATA_DIR} keys add t2 --keyring-backend test
E_ADDR_RES=$("${UND_BIN}" --home ${DATA_DIR} keys add ent1 --output json --keyring-backend test 2>&1)
ENT1=$(echo "${E_ADDR_RES}" | jq --raw-output '.address')

sed -i "s/\"ent_signers\": \"und1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq5x8kpm\"/\"ent_signers\": \"$ENT1\"/g" "${DATA_DIR}/config/genesis.json"
sed -i "s/\"voting_period\": \"172800s\"/\"voting_period\": \"30s\"/g" "${DATA_DIR}/config/genesis.json"

sed -i "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"25.0nund\"/g" "${DATA_DIR}/config/app.toml"
sed -i "s/enable = false/enable = true/g" "${DATA_DIR}/config/app.toml"
sed -i "s/swagger = false/swagger = true/g" "${DATA_DIR}/config/app.toml"

# add accounts to genesis
"${UND_BIN}" add-genesis-account t0 1000000000000000nund --home "${DATA_DIR}" --keyring-backend test
"${UND_BIN}" add-genesis-account validator 1000000000000000nund --home "${DATA_DIR}" --keyring-backend test
"${UND_BIN}" add-genesis-account t1 1000000000000000nund --home "${DATA_DIR}" --keyring-backend test
"${UND_BIN}" add-genesis-account t2 1000000000000000nund --home "${DATA_DIR}" --keyring-backend test
"${UND_BIN}" add-genesis-account ent1 1000000000000000nund --home "${DATA_DIR}" --keyring-backend test

# validator gentx
"${UND_BIN}" gentx validator 1000000nund --home ${DATA_DIR} --chain-id="${CHAIN_ID}"
"${UND_BIN}" collect-gentxs --home "${DATA_DIR}"

# start the daemon
"${UND_BIN}" start --home "${DATA_DIR}"
