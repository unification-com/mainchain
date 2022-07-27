#!/bin/bash

##########################################################################
# A script for quickly testing the in-place upgrade for the 1-ibc        #
# upgrade. The script will set up a single node network using und v1.5.1 #
# and the respective cosmovisor directory structure, with the current    #
# checked out repo as the version to upgrade to.                         #
#                                                                        #
# The script will run cosmovidor, then send an upgrade gov proposal to   #
# run the 1-ibc upgrade at block 10. Cosmovisor will auto-upgrade when   #
# the height is reached.                                                 #
##########################################################################

TEST_PATH="/tmp/und_migrate_test"
UND_HOME="${TEST_PATH}/.und_mainchain"
UND_151_BIN="${TEST_PATH}/und_151"
UND_16X_BIN="${TEST_PATH}/und_16x"
CURR_UND_BIN="${UND_151_BIN}"
NUM_TO_SUB=100
CURRENT_HEIGHT=0
UPPER_CASE_HASH=0
CHAIN_ID_BASE="test-$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 10 | head -n 1)"
CHAIN_ID_V1="${CHAIN_ID_BASE}-1"
CHAIN_ID_V2="${CHAIN_ID_BASE}-2"
WC_SEQ=0
B_SEQ=0
WC_HEIGHT=1

# cosmovisor will run as a background process.
# Catch and kill when ctrl-c is hit
trap "kill 0" EXIT

function set_current_height() {
  CURRENT_HEIGHT=$(${CURR_UND_BIN} status --home "${UND_HOME}" | jq -r '.SyncInfo.latest_block_height')
}

function gen_hash() {
  local UPPER_CASE_HASH=${1:-$UPPER_CASE_HASH}
  local UUID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
  local HASH=$(echo "${UUID}" | openssl dgst -sha256)
  local SHA_HASH_ARR=($HASH)
  local SHA_HASH=${SHA_HASH_ARR[1]}
  if [ $UPPER_CASE_HASH -eq 1 ]
  then
    echo "${SHA_HASH^^}"
  else
    echo "${SHA_HASH}"
  fi
}

function submit_stuff() {
  for (( i=0; i<$NUM_TO_SUB; i++ ))
  do
    B_HASH=$(gen_hash)
    W_HASH=$(gen_hash)
    "${CURR_UND_BIN}" tx wrkchain record 1 --wc_height "${WC_HEIGHT}" --block_hash "${W_HASH}" --from wc --yes --home "${UND_HOME}" --sequence "${WC_SEQ}" --broadcast-mode sync --gas auto --gas-adjustment 1.5
    "${CURR_UND_BIN}" tx beacon record 1 --hash "${B_HASH}" --from bc --yes --home "${UND_HOME}" --sequence "${B_SEQ}" --broadcast-mode sync --gas auto --gas-adjustment 1.5

    WC_HEIGHT=$(awk "BEGIN {print $WC_HEIGHT+1}")
    WC_SEQ=$(awk "BEGIN {print $WC_SEQ+1}")
    B_SEQ=$(awk "BEGIN {print $B_SEQ+1}")
  done
}

function get_curr_acc_sequence() {
  local ADDR=$1
  local RES=$(${CURR_UND_BIN} query account $ADDR --home "${UND_HOME}" --output json)
  local CURR=$(echo "${RES}" | jq --raw-output '.sequence')
  local CURR_INT=$(awk "BEGIN {print $CURR}")
  echo "${CURR_INT}"
}

function query_wc_b() {
  echo "WC 1"
  ${CURR_UND_BIN} query wrkchain wrkchain 1 --home "${UND_HOME}"
  echo "B 1"
  ${CURR_UND_BIN} query beacon beacon 1 --home "${UND_HOME}"
  if [ "$CURR_UND_BIN" = "$UND_16X_BIN" ]; then
    echo "WC 1 STORAGE"
    ${CURR_UND_BIN} query wrkchain storage 1 --home "${UND_HOME}"
    echo "WC 1 SPENT eFUND"
    ${CURR_UND_BIN} query enterprise spent "${WC_ADDRESS}" --home "${UND_HOME}"
    echo "B 1 STORAGE"
    ${CURR_UND_BIN} query beacon storage 1 --home "${UND_HOME}"
    echo "B 1 SPENT eFUND"
    ${CURR_UND_BIN} query enterprise spent "${BC_ADDRESS}" --home "${UND_HOME}"
  fi
}

rm -rf "${TEST_PATH}"
mkdir -p "${TEST_PATH}"

make build

cp "./build/und" "${UND_16X_BIN}"

cd "${TEST_PATH}"

wget https://github.com/unification-com/mainchain/releases/download/1.5.1/und_v1.5.1_linux_x86_64.tar.gz
tar -zxvf und_v1.5.1_linux_x86_64.tar.gz
mv und "${UND_151_BIN}"

"${CURR_UND_BIN}" init test --home "${UND_HOME}"
"${CURR_UND_BIN}" unsafe-reset-all --home "${UND_HOME}"
"${CURR_UND_BIN}" config chain-id "${CHAIN_ID_V1}" --home "${UND_HOME}"
"${CURR_UND_BIN}" config keyring-backend test --home "${UND_HOME}"
"${CURR_UND_BIN}" config broadcast-mode block --home "${UND_HOME}"

"${CURR_UND_BIN}" init test --chain-id "${CHAIN_ID_V1}" --overwrite --home "${UND_HOME}"

sed -i -e 's/"stake"/"nund"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"historical_entries": 10000/"historical_entries": 3/gi' "${UND_HOME}/config/genesis.json"

sed -i -e 's/pruning = "default"/pruning = "nothing"/gi' "${UND_HOME}/config/app.toml"
sed -i -e 's/enable = false/enable = true/gi' "${UND_HOME}/config/app.toml"
sed -i -e 's/swagger = false/swagger = true/gi' "${UND_HOME}/config/app.toml"
sed -i -e 's/minimum-gas-prices = ""/minimum-gas-prices = "1.0nund"/gi' "${UND_HOME}/config/app.toml"

# accounts
"${CURR_UND_BIN}" keys add validator --home "${UND_HOME}" --keyring-backend test
"${CURR_UND_BIN}" add-genesis-account validator 5000000000nund --keyring-backend test --home "${UND_HOME}"

ENT_ADDRESS=$("${CURR_UND_BIN}" keys add ent --home "${UND_HOME}" --keyring-backend test --output json | jq -r ".address")
"${CURR_UND_BIN}" add-genesis-account ent 5000000000nund --keyring-backend test --home "${UND_HOME}"

WC_ADDRESS=$("${CURR_UND_BIN}" keys add wc --home "${UND_HOME}" --keyring-backend test --output json | jq -r ".address")
"${CURR_UND_BIN}" add-genesis-account wc 5000000000nund --keyring-backend test --home "${UND_HOME}"
BC_ADDRESS=$("${CURR_UND_BIN}" keys add bc --home "${UND_HOME}" --keyring-backend test --output json | jq -r ".address")
"${CURR_UND_BIN}" add-genesis-account bc 5000000000nund --keyring-backend test --home "${UND_HOME}"

sed -i -e "s/\"ent_signers\": \"und1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq5x8kpm\"/\"ent_signers\": \"${ENT_ADDRESS}\"/gi" "${UND_HOME}/config/genesis.json"
sed -i -e "s/\"whitelist\": \[\]/\"whitelist\": \[\"${WC_ADDRESS}\",\"${BC_ADDRESS}\"\]/gi" "${UND_HOME}/config/genesis.json"

# gentx
"${CURR_UND_BIN}" gentx validator 1000000nund --chain-id "${CHAIN_ID_V1}" --home "${UND_HOME}"
"${CURR_UND_BIN}" collect-gentxs --home "${UND_HOME}"

"${CURR_UND_BIN}" unsafe-reset-all --home "${UND_HOME}"

"${CURR_UND_BIN}" start --home "${UND_HOME}" &
UND_PID="$!"

sleep 6s

# ent POs
"${CURR_UND_BIN}" tx enterprise purchase 1000000000000000nund --from wc --yes --home "${UND_HOME}" --gas-prices 1.0nund --gas auto --gas-adjustment 1.5
"${CURR_UND_BIN}" tx enterprise purchase 1000000000000000nund --from bc --yes --home "${UND_HOME}" --gas-prices 1.0nund --gas auto --gas-adjustment 1.5

sleep 6s
"${CURR_UND_BIN}" tx enterprise process 1 accept --from ent --yes --home "${UND_HOME}" --sequence 0 --gas-prices 1.0nund --gas auto --gas-adjustment 1.5
sleep 1s
"${CURR_UND_BIN}" tx enterprise process 2 accept --from ent --yes --home "${UND_HOME}" --sequence 1 --gas-prices 1.0nund --gas auto --gas-adjustment 1.5

sleep 15s

# register WC/BEACON
"${CURR_UND_BIN}" tx wrkchain register --moniker="wc1" --genesis="genhash" --name="WC 1" --base="geth" --from wc --yes --home "${UND_HOME}" --gas auto --gas-adjustment 1.5
"${CURR_UND_BIN}" tx beacon register --moniker="bc1" --name="BC 1" --from bc --yes --home "${UND_HOME}" --gas auto --gas-adjustment 1.5

# arbitrary gov proposal
"${CURR_UND_BIN}" tx gov submit-proposal --type text --title some_prop --description some_prop_desc --from validator --yes --home "${UND_HOME}" --gas-prices 1.0nund --gas auto --gas-adjustment 1.5
"${CURR_UND_BIN}" tx gov deposit 1 10000000nund --from validator --yes --home "${UND_HOME}" --gas-prices 1.0nund --gas auto --gas-adjustment 1.5
"${CURR_UND_BIN}" tx gov vote 1 yes --from validator --yes --home "${UND_HOME}" --gas-prices 1.0nund --gas auto --gas-adjustment 1.5

sleep 10s

WC_SEQ=2
B_SEQ=2

# submit stuff
submit_stuff

# wait for last Txs to process
sleep 7s

query_wc_b

# stop und 1.5.1
kill "${UND_PID}"

sleep 7s

echo "Exporting und v1.5.1 state"

"${CURR_UND_BIN}" export --for-zero-height --home "${UND_HOME}" > "${UND_HOME}"/v042_exported_state.json

cat "${UND_HOME}"/v042_exported_state.json | jq > "${UND_HOME}"/v042_exported_state_pretty.json

CURR_UND_BIN="${UND_16X_BIN}"

echo "Migrating state to v1.6.x format"

"${CURR_UND_BIN}" migrate "${UND_HOME}"/v042_exported_state.json --chain-id ${CHAIN_ID_V2} --genesis-time "2022-07-27T07:00:00Z" --log_level "" > "${UND_HOME}"/new_v045_genesis.json

cat "${UND_HOME}"/new_v045_genesis.json | jq > "${UND_HOME}"/new_v045_genesis_pretty.json

cp "${UND_HOME}"/new_v045_genesis.json "${UND_HOME}"/config/genesis.json

"${CURR_UND_BIN}" tendermint unsafe-reset-all --home "${UND_HOME}"

echo "Start chain"

"${CURR_UND_BIN}" start --home "${UND_HOME}" &

sleep 12s

# update chain id in client config
"${CURR_UND_BIN}" config chain-id "${CHAIN_ID_V2}" --home "${UND_HOME}"

sleep 1s

WC_SEQ=$(get_curr_acc_sequence "${WC_ADDRESS}")
B_SEQ=$(get_curr_acc_sequence "${BC_ADDRESS}")

echo "WC_SEQ=${WC_SEQ}"
echo "B_SEQ=${B_SEQ}"

query_wc_b

# submit more stuff
submit_stuff

sleep 7

query_wc_b

wait
