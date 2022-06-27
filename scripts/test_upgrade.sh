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
UPGRADE_HEIGHT=10

# cosmovisor will run as a background process.
# Catch and kill when ctrl-c is hit
trap "kill 0" EXIT

rm -rf "${TEST_PATH}"

mkdir -p "${COSMOVISOR_HOME}/genesis/bin"
mkdir -p "${COSMOVISOR_HOME}/upgrades/1-ibc/bin"

make build

cp "./build/und" "${COSMOVISOR_HOME}/upgrades/1-ibc/bin"

cd "${TEST_PATH}"

wget https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz
tar -zxvf cosmovisor-v1.1.0-linux-amd64.tar.gz

wget https://github.com/unification-com/mainchain/releases/download/1.5.1/und_v1.5.1_linux_x86_64.tar.gz
tar -zxvf und_v1.5.1_linux_x86_64.tar.gz
mv und "${UND_GEN_BIN}"

"${UND_GEN_BIN}" init test --home "${UND_HOME}"
"${UND_GEN_BIN}" unsafe-reset-all --home "${UND_HOME}"
"${UND_GEN_BIN}" config chain-id test --home "${UND_HOME}"
"${UND_GEN_BIN}" config keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" config broadcast-mode block --home "${UND_HOME}"

"${UND_GEN_BIN}" init test --chain-id test --overwrite --home "${UND_HOME}"

sed -i -e 's/"voting_period": "172800s"/"voting_period": "20s"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"stake"/"nund"/gi' "${UND_HOME}/config/genesis.json"
sed -i -e 's/"historical_entries": 10000/"historical_entries": 3/gi' "${UND_HOME}/config/genesis.json"

"${UND_GEN_BIN}" keys add validator --home "${UND_HOME}"
"${UND_GEN_BIN}" add-genesis-account validator 5000000000nund --keyring-backend test --home "${UND_HOME}"
"${UND_GEN_BIN}" gentx validator 1000000nund --chain-id test --home "${UND_HOME}"
"${UND_GEN_BIN}" collect-gentxs --home "${UND_HOME}"

export DAEMON_NAME=und
export DAEMON_HOME="${UND_HOME}"
export DAEMON_RESTART_AFTER_UPGRADE=true

"${UND_GEN_BIN}" unsafe-reset-all --home "${UND_HOME}"

"${COSMOVISOR_BIN}" run start --home "${UND_HOME}" &

sleep 6s
"${UND_GEN_BIN}" tx gov submit-proposal software-upgrade 1-ibc --title upgrade --description upgrade --upgrade-height ${UPGRADE_HEIGHT} --from validator --yes --home "${UND_HOME}"
"${UND_GEN_BIN}" tx gov deposit 1 10000000nund --from validator --yes --home "${UND_HOME}"
"${UND_GEN_BIN}" tx gov vote 1 yes --from validator --yes --home "${UND_HOME}"

wait
