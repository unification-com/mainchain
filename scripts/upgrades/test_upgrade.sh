#!/bin/bash

##########################################################################
# A script for quickly testing the in-place upgrade for the 1-ibc        #
# upgrade. The script will set up a single node network using und v1.5.1 #
# and the respective cosmovisor directory structure, with the current    #
# checked out repo as the version to upgrade to.                         #
#                                                                        #
# The script will run cosmovidor, then send an upgrade gov proposal to   #
# run the specified upgrade at block 10. Cosmovisor will auto-upgrade    #
# when the height is reached.                                            #
##########################################################################

TEST_PATH="/tmp/und_upgrade_test"
UND_HOME="${TEST_PATH}/.und_mainchain"
COSMOVISOR_HOME="${UND_HOME}/cosmovisor"
COSMOVISOR_BIN="${TEST_PATH}/cosmovisor"
UND_GEN_BIN="${COSMOVISOR_HOME}/genesis/bin/und"
UPGRADE_HEIGHT=10
CHAIN_ID="test-$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 10 | head -n 1)"
UPGRADE_PLAN_NAME="3-keyleth"
UND_GENESIS_VERSION="v1.7.0"

# cosmovisor will run as a background process.
# Catch and kill when ctrl-c is hit
trap "kill 0" EXIT

rm -rf "${TEST_PATH}"

mkdir -p "${COSMOVISOR_HOME}/genesis/bin"
mkdir -p "${COSMOVISOR_HOME}/upgrades/${UPGRADE_PLAN_NAME}/bin"

make build

cp "./build/und" "${COSMOVISOR_HOME}/upgrades/${UPGRADE_PLAN_NAME}/bin"

cd "${TEST_PATH}"

wget https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.2.0/cosmovisor-v1.2.0-linux-amd64.tar.gz
tar -zxvf cosmovisor-v1.2.0-linux-amd64.tar.gz

wget "https://github.com/unification-com/mainchain/releases/download/${UND_GENESIS_VERSION}/und_${UND_GENESIS_VERSION}_linux_x86_64.tar.gz"
tar -zxvf "und_${UND_GENESIS_VERSION}_linux_x86_64.tar.gz"
mv und "${UND_GEN_BIN}"

"${UND_GEN_BIN}" init test --home "${UND_HOME}"
"${UND_GEN_BIN}" unsafe-reset-all --home "${UND_HOME}"
"${UND_GEN_BIN}" config chain-id "${CHAIN_ID}" --home "${UND_HOME}"
"${UND_GEN_BIN}" config keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" config broadcast-mode block --home "${UND_HOME}"

"${UND_GEN_BIN}" init test --chain-id "${CHAIN_ID}" --overwrite --home "${UND_HOME}"

"${UND_GEN_BIN}" keys add validator --home "${UND_HOME}"

sed -i -e 's/"voting_period": "172800s"/"voting_period": "20s"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"stake"/"nund"/gi' "${UND_HOME}/config/genesis.json"

"${UND_GEN_BIN}" add-genesis-account validator 5000000000nund --keyring-backend test --home "${UND_HOME}"

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
