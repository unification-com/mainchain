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

TEST_PATH="/tmp/und_upgrade_test"
UND_HOME="${TEST_PATH}/.und_mainchain"
COSMOVISOR_HOME="${UND_HOME}/cosmovisor"
COSMOVISOR_BIN="${TEST_PATH}/cosmovisor"
UND_GEN_BIN="${COSMOVISOR_HOME}/genesis/bin/und"
UPGRADE_AFTER=10
NUM_TO_SUB=50
CURRENT_HEIGHT=0
UPPER_CASE_HASH=0
CHAIN_ID="test-$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 10 | head -n 1)"

# cosmovisor will run as a background process.
# Catch and kill when ctrl-c is hit
trap "kill 0" EXIT

function set_current_height() {
  CURRENT_HEIGHT=$(${UND_GEN_BIN} status --home "${UND_HOME}" | jq -r '.SyncInfo.latest_block_height')
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

rm -rf "${TEST_PATH}"

mkdir -p "${COSMOVISOR_HOME}/genesis/bin"
mkdir -p "${COSMOVISOR_HOME}/upgrades/1-ibc/bin"

make build

cp "./build/und" "${COSMOVISOR_HOME}/upgrades/1-ibc/bin"

cd "${TEST_PATH}"

wget https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.2.0/cosmovisor-v1.2.0-linux-amd64.tar.gz
tar -zxvf cosmovisor-v1.2.0-linux-amd64.tar.gz

wget https://github.com/unification-com/mainchain/releases/download/1.5.1/und_v1.5.1_linux_x86_64.tar.gz
tar -zxvf und_v1.5.1_linux_x86_64.tar.gz
mv und "${UND_GEN_BIN}"

"${UND_GEN_BIN}" init test --home "${UND_HOME}"
"${UND_GEN_BIN}" unsafe-reset-all --home "${UND_HOME}"
"${UND_GEN_BIN}" config chain-id "${CHAIN_ID}" --home "${UND_HOME}"
"${UND_GEN_BIN}" config keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" config broadcast-mode block --home "${UND_HOME}"

"${UND_GEN_BIN}" init test --chain-id "${CHAIN_ID}" --overwrite --home "${UND_HOME}"

sed -i -e 's/"voting_period": "172800s"/"voting_period": "20s"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"stake"/"nund"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"historical_entries": 10000/"historical_entries": 3/gi' "${UND_HOME}/config/genesis.json"

sed -i -e 's/pruning = "default"/pruning = "nothing"/gi' "${UND_HOME}/config/app.toml"
sed -i -e 's/enable = false/enable = true/gi' "${UND_HOME}/config/app.toml"
sed -i -e 's/swagger = false/swagger = true/gi' "${UND_HOME}/config/app.toml"

# accounts
"${UND_GEN_BIN}" keys add validator --home "${UND_HOME}" --keyring-backend test
"${UND_GEN_BIN}" add-genesis-account validator 5000000000nund --keyring-backend test --home "${UND_HOME}"

ENT_ADDRESS=$("${UND_GEN_BIN}" keys add ent --home "${UND_HOME}" --keyring-backend test --output json | jq -r ".address")
"${UND_GEN_BIN}" add-genesis-account ent 5000000000nund --keyring-backend test --home "${UND_HOME}"

WC_ADDRESS=$("${UND_GEN_BIN}" keys add wc --home "${UND_HOME}" --keyring-backend test --output json | jq -r ".address")
"${UND_GEN_BIN}" add-genesis-account wc 5000000000nund --keyring-backend test --home "${UND_HOME}"
BC_ADDRESS=$("${UND_GEN_BIN}" keys add bc --home "${UND_HOME}" --keyring-backend test --output json | jq -r ".address")
"${UND_GEN_BIN}" add-genesis-account bc 5000000000nund --keyring-backend test --home "${UND_HOME}"

sed -i -e "s/\"ent_signers\": \"und1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq5x8kpm\"/\"ent_signers\": \"${ENT_ADDRESS}\"/gi" "${UND_HOME}/config/genesis.json"
sed -i -e "s/\"whitelist\": \[\]/\"whitelist\": \[\"${WC_ADDRESS}\",\"${BC_ADDRESS}\"\]/gi" "${UND_HOME}/config/genesis.json"

# gentx
"${UND_GEN_BIN}" gentx validator 1000000nund --chain-id "${CHAIN_ID}" --home "${UND_HOME}"
"${UND_GEN_BIN}" collect-gentxs --home "${UND_HOME}"

export DAEMON_NAME=und
export DAEMON_HOME="${UND_HOME}"
export DAEMON_RESTART_AFTER_UPGRADE=true

"${UND_GEN_BIN}" unsafe-reset-all --home "${UND_HOME}"

"${COSMOVISOR_BIN}" run start --home "${UND_HOME}" &

sleep 6s

# ent POs
"${UND_GEN_BIN}" tx enterprise purchase 1000000000000000nund --from wc --yes --home "${UND_HOME}"
"${UND_GEN_BIN}" tx enterprise purchase 1000000000000000nund --from bc --yes --home "${UND_HOME}"

sleep 6s
"${UND_GEN_BIN}" tx enterprise process 1 accept --from ent --yes --home "${UND_HOME}" --sequence 0
sleep 1s
"${UND_GEN_BIN}" tx enterprise process 2 accept --from ent --yes --home "${UND_HOME}" --sequence 1

sleep 15s

# register WC/BEACON
"${UND_GEN_BIN}" tx wrkchain register --moniker="wc1" --genesis="genhash" --name="WC 1" --base="geth" --from wc --yes --home "${UND_HOME}"
"${UND_GEN_BIN}" tx beacon register --moniker="bc1" --name="BC 1" --from bc --yes --home "${UND_HOME}"

sleep 10s

WC_SEQ=2
B_SEQ=2

# submit stuff
for (( i=0; i<$NUM_TO_SUB; i++ ))
do
  HEIGHT=$(awk "BEGIN {print $i+1}")
  B_HASH=$(gen_hash)
  W_HASH=$(gen_hash)
  "${UND_GEN_BIN}" tx wrkchain record 1 --wc_height "${HEIGHT}" --block_hash "${W_HASH}" --from wc --yes --home "${UND_HOME}" --sequence "${WC_SEQ}" --broadcast-mode sync
  "${UND_GEN_BIN}" tx beacon record 1 --hash "${B_HASH}" --from bc --yes --home "${UND_HOME}" --sequence "${B_SEQ}" --broadcast-mode sync

  WC_SEQ=$(awk "BEGIN {print $WC_SEQ+1}")
  B_SEQ=$(awk "BEGIN {print $B_SEQ+1}")
done

# upgrade gov proposal
set_current_height
UPGRADE_HEIGHT=$(awk "BEGIN {print $CURRENT_HEIGHT+$UPGRADE_AFTER}")
"${UND_GEN_BIN}" tx gov submit-proposal software-upgrade 1-ibc --deposit 10000000nund --title upgrade --description upgrade --upgrade-height ${UPGRADE_HEIGHT} --from validator --yes --home "${UND_HOME}"
"${UND_GEN_BIN}" tx gov vote 1 yes --from validator --yes --home "${UND_HOME}"

echo "UPGRADE_HEIGHT=${UPGRADE_HEIGHT}"

wait
