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
UPGRADE_HEIGHT=8
CHAIN_ID="test-$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 10 | head -n 1)"
UPGRADE_PLAN_NAME="1-init_ibc"

# cosmovisor will run as a background process.
# Catch and kill when ctrl-c is hit
trap "kill 0" EXIT

rm -rf "${TEST_PATH}"

mkdir -p "${COSMOVISOR_HOME}/genesis/bin"
mkdir -p "${COSMOVISOR_HOME}/upgrades/${UPGRADE_PLAN_NAME}/bin"

make build

cp "./build/und" "${COSMOVISOR_HOME}/upgrades/${UPGRADE_PLAN_NAME}/bin"

cd "${TEST_PATH}"

wget https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz
tar -zxvf cosmovisor-v1.1.0-linux-amd64.tar.gz

wget https://github.com/unification-com/mainchain/releases/download/1.5.1/und_v1.5.1_linux_x86_64.tar.gz
tar -zxvf und_v1.5.1_linux_x86_64.tar.gz
mv und "${UND_GEN_BIN}"

"${UND_GEN_BIN}" init test --home "${UND_HOME}"
"${UND_GEN_BIN}" unsafe-reset-all --home "${UND_HOME}"
"${UND_GEN_BIN}" config chain-id "${CHAIN_ID}" --home "${UND_HOME}"
"${UND_GEN_BIN}" config keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" config broadcast-mode block --home "${UND_HOME}"

"${UND_GEN_BIN}" init test --chain-id "${CHAIN_ID}" --overwrite --home "${UND_HOME}"

"${UND_GEN_BIN}" keys add validator --home "${UND_HOME}"
"${UND_GEN_BIN}" keys add t1 --home ${UND_HOME}
"${UND_GEN_BIN}" keys add t2 --home ${UND_HOME}
E_ADDR_RES=$("${UND_GEN_BIN}" --home ${UND_HOME} keys add ent1 --output json 2>&1)
ENT1=$(echo "${E_ADDR_RES}" | jq --raw-output '.address')
E_ADDR_RES=$("${UND_GEN_BIN}" --home ${UND_HOME} keys add ent2 --output json 2>&1)
ENT2=$(echo "${E_ADDR_RES}" | jq --raw-output '.address')

sed -i "s/\"ent_signers\": \"und1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq5x8kpm\"/\"ent_signers\": \"$ENT1,$ENT2\"/g" "${UND_HOME}/config/genesis.json"
sed -i "s/\"min_accepts\": \"1\"/\"min_accepts\": \"2\"/g" "${UND_HOME}/config/genesis.json"
sed -i -e 's/"voting_period": "172800s"/"voting_period": "20s"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"stake"/"nund"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"historical_entries": 10000/"historical_entries": 3/gi' "${UND_HOME}/config/genesis.json"

sed -i "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"25.0nund\"/g" "${UND_HOME}/config/app.toml"
sed -i "s/enable = false/enable = true/g" "${UND_HOME}/config/app.toml"
sed -i "s/swagger = false/swagger = true/g" "${UND_HOME}/config/app.toml"

"${UND_GEN_BIN}" add-genesis-account validator 5000000000nund --keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" add-genesis-account t1 100000000000000nund --keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" add-genesis-account ent1 100000000000000nund --keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" add-genesis-account ent2 100000000000000nund --keyring-backend test --home "${UND_HOME}"

"${UND_GEN_BIN}" gentx validator 1000000nund --chain-id "${CHAIN_ID}" --home "${UND_HOME}"
"${UND_GEN_BIN}" collect-gentxs --home "${UND_HOME}"

export DAEMON_NAME=und
export DAEMON_HOME="${UND_HOME}"
export DAEMON_RESTART_AFTER_UPGRADE=true

"${UND_GEN_BIN}" unsafe-reset-all --home "${UND_HOME}"

echo "Start node & submit upgrade proposal "${UPGRADE_PLAN_NAME}" for height ${UPGRADE_HEIGHT}"

"${COSMOVISOR_BIN}" run start --home "${UND_HOME}" &

sleep 6s
"${UND_GEN_BIN}" tx gov submit-proposal software-upgrade "${UPGRADE_PLAN_NAME}" --title upgrade --description upgrade --upgrade-height ${UPGRADE_HEIGHT} --deposit 10000000nund --from validator --yes --home "${UND_HOME}" --gas auto --gas-adjustment 1.5 --gas-prices=25.0nund
"${UND_GEN_BIN}" tx gov vote 1 yes --from validator --yes --home "${UND_HOME}" --gas auto --gas-adjustment 1.5 --gas-prices=25.0nund

wait
