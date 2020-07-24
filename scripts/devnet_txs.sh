#!/bin/bash

#########################################################################
# A script for generating and broadcasting some test transactions for   #
# populating DevNet                                                     #
#                                                                       #
# Note: the script assumes the accounts in Docker/README.md have been   #
#       imported into the keychain.                                     #
#########################################################################

UNDCLI_BIN="./build/undcli"
DEVNET_RPC_IP="172.25.0.3"
DEVNET_RPC_PORT="26661"
DEVNET_RPC_TCP="tcp://${DEVNET_RPC_IP}:${DEVNET_RPC_PORT}"
DEVNET_RPC_HTTP="http://${DEVNET_RPC_IP}:${DEVNET_RPC_PORT}"
CHAIN_ID="FUND-Mainchain-DevNet"
BROADCAST_MODE="block"
GAS_PRICES="0.25nund"

# Account names as imported into undcli keys
NODE1_ACC="node1"
NODE2_ACC="node2"
NODE3_ACC="node3"
ENT_ACC="ent"
T1_ACC="t1"
T2_ACC="t2"
T3_ACC="t3"

gen_hash() {
  UUID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
  HASH=$(echo "${UUID}" | openssl dgst -sha256)
  SHA_HASH_ARR=($HASH)
  SHA_HASH=${SHA_HASH_ARR[1]}
  echo "${SHA_HASH}"
}

get_addr() {
  ADDR=$(${UNDCLI_BIN} keys show $1 -a)
  echo "${ADDR}"
}

get_base_flags() {
  BROADCAST=${1:-$BROADCAST_MODE}
  FLAGS="--broadcast-mode=${BROADCAST} --chain-id=${CHAIN_ID} --node=${DEVNET_RPC_TCP} --gas=auto --gas-adjustment=1.5 -y"
  echo "${FLAGS}"
}

get_gas_flags() {
  FLAGS="--gas-prices=${GAS_PRICES}"
  echo "${FLAGS}"
}

check_accounts_exist() {
  if { ${UNDCLI_BIN} keys show $1 2>&1 >&3 3>&- | grep '^' >&2; } 3>&1; then
    echo "${1} does not seem to exist in keyring. Exiting"
    exit 1
  else
    echo "Found ${1} in keyring"
  fi
}

check_accounts_exist ${NODE1_ACC}
check_accounts_exist ${NODE2_ACC}
check_accounts_exist ${NODE3_ACC}
check_accounts_exist ${ENT_ACC}
check_accounts_exist ${T1_ACC}
check_accounts_exist ${T2_ACC}
check_accounts_exist ${T3_ACC}

# Wait for Node1 to come online
echo "Waiting for DevNet Node1 to come online"
until nc -z 172.25.0.3 26661;
do
  echo "Waiting for DevNet Node1 to come online"
  sleep 1
done

echo "Node1 is online"

# Wait for first block
until [ $(curl -s ${DEVNET_RPC_HTTP}/status | jq --raw-output '.result.sync_info.latest_block_height') -gt 1 ]
do
  echo "Waiting for DevNet Block #1"
  sleep 1
done

echo "Block >= 1 committed"
echo "Running transactions"

START_TIME=$(date +%s)

${UNDCLI_BIN} tx send $(get_addr ${NODE1_ACC}) $(get_addr ${T1_ACC}) 100000000000nund $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${T1_ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${T2_ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${T3_ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${T1_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${T2_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${T3_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise process 1 accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise process 2 accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx enterprise process 3 accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags)
sleep 1s

${UNDCLI_BIN} tx wrkchain register --moniker="wrkchain1" --genesis="$(gen_hash)" --name="Wrkchain 1" --base="geth" --from=${T1_ACC} $(get_base_flags)
sleep 1s

${UNDCLI_BIN} tx beacon register --moniker="beacon1" --name="Beacon 1" --from=${T2_ACC} $(get_base_flags)
sleep 1s

${UNDCLI_BIN} tx beacon register --moniker="beacon2" --name="Beacon 2" --from=${T3_ACC} $(get_base_flags)
sleep 1s

${UNDCLI_BIN} tx wrkchain record 1 --wc_height=1 --block_hash="$(gen_hash)" --parent_hash="$(gen_hash)" --hash1="$(gen_hash)" --hash2="$(gen_hash)" --hash3="$(gen_hash)" --from=${T1_ACC} $(get_base_flags)
sleep 1s

${UNDCLI_BIN} tx beacon record 1 --hash="$(gen_hash)" --subtime=$(date +%s) --from=${T2_ACC} $(get_base_flags sync)
sleep 1s

${UNDCLI_BIN} tx send $(get_addr ${NODE1_ACC}) $(get_addr ${T1_ACC}) 100000000000nund $(get_base_flags sync) $(get_gas_flags)
sleep 5s

${UNDCLI_BIN} tx beacon record 2 --hash="$(gen_hash)" --subtime=$(date +%s) --from=${T3_ACC} $(get_base_flags sync)
sleep 1s

${UNDCLI_BIN} tx send $(get_addr ${NODE2_ACC}) $(get_addr ${T2_ACC}) 10000000000nund $(get_base_flags sync) $(get_gas_flags)
sleep 5s

${UNDCLI_BIN} tx wrkchain record 1 --wc_height=2 --block_hash="$(gen_hash)" --parent_hash="$(gen_hash)" --hash1="$(gen_hash)" --hash2="$(gen_hash)" --hash3="$(gen_hash)" --from=${T1_ACC} $(get_base_flags sync)
sleep 1s

${UNDCLI_BIN} tx beacon record 2 --hash="$(gen_hash)" --subtime=$(date +%s) --from=${T3_ACC} $(get_base_flags sync)
sleep 1s

${UNDCLI_BIN} tx send $(get_addr ${NODE3_ACC}) $(get_addr ${T3_ACC}) 24000000000nund $(get_base_flags sync) $(get_gas_flags)

${UNDCLI_BIN} query wrkchain wrkchain 1 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 1 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 2 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}

echo "Finished in $(($(date +%s) - START_TIME)) seconds."
echo ""
