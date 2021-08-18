#!/bin/bash

#########################################################################
# A script for syncing to DevNet and exporting data to genesis and      #
# data files.                                                           #
#                                                                       #
# Run after DevNet has been populated using devnet_txs.sh               #
#                                                                       #
# Note: the script assumes jq is installed & "make build" has been run  #
#########################################################################

UND_BIN="./build/und"
UND_HOME_DIR="/tmp/.und_mainchain_DevNet"
DEVNET_GENESIS="./Docker/assets/node1/config/genesis.json"
CHAIN_ID="FUND-Mainchain-DevNet"
DEVNET_NODE1_HOST="localhost"
DEVNET_NODE1_RPC_PORT="26661"
DEVNET_RPC_HTTP="http://${DEVNET_NODE1_HOST}:${DEVNET_NODE1_RPC_PORT}"

rm -rf "${UND_HOME_DIR}"

${UND_BIN} init devnet_node --home="${UND_HOME_DIR}"

cp "${DEVNET_GENESIS}" "${UND_HOME_DIR}/config/genesis.json"
sed -i -e "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"0.25nund\"/g" "${UND_HOME_DIR}/config/app.toml"

sed -i -e "s/persistent_peers = \"\"/persistent_peers = \"53e857acc2df7127d5ef33b0dd98c55e7068ae06@localhost:26652,33a49c1eae31ce82ffab25ed821e8cec7f8bbd00@localhost:26653\"/g" "${UND_HOME_DIR}/config/config.toml"
sed -i -e "s/addr_book_strict = true/addr_book_strict = false/g" "${UND_HOME_DIR}/config/config.toml"

mkdir -p "${UND_HOME_DIR}/export"
HALT_HEIGHT=$(curl -s ${DEVNET_RPC_HTTP}/status | jq --raw-output '.result.sync_info.latest_block_height')

${UND_BIN} start --home="${UND_HOME_DIR}" --halt-height=${HALT_HEIGHT} --grpc.address="127.0.0.1:9999"

echo "Export state at height ${HALT_HEIGHT} to new genesis."
${UND_BIN} export --home="${UND_HOME_DIR}" --for-zero-height | jq > "${UND_HOME_DIR}/export/genesis.json"
echo "exported to ${UND_HOME_DIR}/export/genesis.json"
