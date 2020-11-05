#!/bin/bash

#########################################################################
# A script for generating and broadcasting some test transactions for   #
# populating DevNet                                                     #
#                                                                       #
# Note: the script assumes the accounts in Docker/README.md have been   #
#       imported into the keychain, jq is installed, and "make build"   #
#       has been run.                                                   #
#########################################################################

UNDCLI_BIN="./build/undcli"
DEVNET_RPC_IP="localhost"
DEVNET_RPC_PORT="26661"
DEVNET_RPC_TCP="tcp://${DEVNET_RPC_IP}:${DEVNET_RPC_PORT}"
DEVNET_RPC_HTTP="http://${DEVNET_RPC_IP}:${DEVNET_RPC_PORT}"
CHAIN_ID="FUND-Mainchain-DevNet"
BROADCAST_MODE="sync"
GAS_PRICES="0.25nund"
NUM_TO_SUB=200
UPPER_CASE_HASH=0

# Account names as imported into undcli keys
NODE1_ACC="node1"
NODE2_ACC="node2"
NODE3_ACC="node3"
ENT_ACC="ent"
T1_ACC="t1"
T2_ACC="t2"
T3_ACC="t3"
T4_ACC="t4"

gen_hash() {
  UPPER_CASE_HASH=${1:-$UPPER_CASE_HASH}
  UUID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
  HASH=$(echo "${UUID}" | openssl dgst -sha256)
  SHA_HASH_ARR=($HASH)
  SHA_HASH=${SHA_HASH_ARR[1]}
  if [ $UPPER_CASE_HASH -eq 1 ]
  then
    echo "${SHA_HASH^^}"
  else
    echo "${SHA_HASH}"
  fi
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
check_accounts_exist ${T4_ACC}

# Wait for Node1 to come online
echo "Waiting for DevNet Node1 to come online"
until nc -z localhost 26661;
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

${UNDCLI_BIN} tx send $(get_addr ${NODE1_ACC}) $(get_addr ${T1_ACC}) 100000000000nund $(get_base_flags) $(get_gas_flags) --sequence=0

echo "Whitelist  ${T1_ACC}, ${T2_ACC}, ${T3_ACC}, ${T4_ACC} for Enterprise POs"
${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${T1_ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=0
${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${T2_ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=1
${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${T3_ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=2
${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${T4_ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=3
echo "Done. Wait for approx. 1 block"
sleep 6s

echo "${T1_ACC}, ${T2_ACC}, ${T3_ACC}, ${T4_ACC} raise Enterprise POs"
${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${T1_ACC} $(get_base_flags) $(get_gas_flags) --sequence=0
${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${T2_ACC} $(get_base_flags) $(get_gas_flags) --sequence=0
${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${T3_ACC} $(get_base_flags) $(get_gas_flags) --sequence=0
${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${T4_ACC} $(get_base_flags) $(get_gas_flags) --sequence=0
echo "Done. Wait for approx. 1 block"
sleep 6s

echo "Process Enterprise POs"
${UNDCLI_BIN} tx enterprise process 1 accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=4
${UNDCLI_BIN} tx enterprise process 2 accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=5
${UNDCLI_BIN} tx enterprise process 3 accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=6
${UNDCLI_BIN} tx enterprise process 4 accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=7
echo "Done. Wait for approx. 1 block"
sleep 6s

WC1_GEN_HASH="0x$(gen_hash)"
WC2_GEN_HASH="$(gen_hash 1)"

echo "${T1_ACC}, ${T2_ACC}, ${T3_ACC}, ${T4_ACC} register WRKChain and BEACONs"
${UNDCLI_BIN} tx wrkchain register --moniker="wrkchain1" --genesis="${WC1_GEN_HASH}" --name="Wrkchain 1: geth" --base="geth" --from=${T1_ACC} $(get_base_flags) --sequence=1
${UNDCLI_BIN} tx beacon register --moniker="beacon1" --name="Beacon 1" --from=${T2_ACC} $(get_base_flags) --sequence=1
${UNDCLI_BIN} tx beacon register --moniker="beacon2" --name="Beacon 2" --from=${T3_ACC} $(get_base_flags) --sequence=1
${UNDCLI_BIN} tx wrkchain register --moniker="wrkchain2" --genesis="${WC2_GEN_HASH}" --name="Wrkchain 2: tendermint" --base="tendermint" --from=${T4_ACC} $(get_base_flags) --sequence=1
echo "Done. Wait for approx. 1 block"
sleep 6s

CURRENT_SEQUENCE_ACC_1=$(${UNDCLI_BIN} query account $(get_addr ${T1_ACC})  --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} | jq --raw-output '.account.value.sequence')
CURRENT_SEQUENCE_ACC_2=$(${UNDCLI_BIN} query account $(get_addr ${T2_ACC})  --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} | jq --raw-output '.account.value.sequence')
CURRENT_SEQUENCE_ACC_3=$(${UNDCLI_BIN} query account $(get_addr ${T3_ACC})  --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} | jq --raw-output '.account.value.sequence')
CURRENT_SEQUENCE_ACC_4=$(${UNDCLI_BIN} query account $(get_addr ${T4_ACC})  --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} | jq --raw-output '.account.value.sequence')

CURRENT_SEQUENCE_INT_ACC_1=$(awk "BEGIN {print $CURRENT_SEQUENCE_ACC_1}")
CURRENT_SEQUENCE_INT_ACC_2=$(awk "BEGIN {print $CURRENT_SEQUENCE_ACC_2}")
CURRENT_SEQUENCE_INT_ACC_3=$(awk "BEGIN {print $CURRENT_SEQUENCE_ACC_3}")
CURRENT_SEQUENCE_INT_ACC_4=$(awk "BEGIN {print $CURRENT_SEQUENCE_ACC_4}")

WC1_P_HASH=${WC1_GEN_HASH}
WC2_P_HASH=${WC2_GEN_HASH}

for (( i=0; i<$NUM_TO_SUB; i++ ))
do
  SUB_NUM_OF=$(awk "BEGIN {print $i+1}")
  NEXT_SEQ_ACC_1=$(awk "BEGIN {print $CURRENT_SEQUENCE_INT_ACC_1+$i}")
  NEXT_SEQ_ACC_2=$(awk "BEGIN {print $CURRENT_SEQUENCE_INT_ACC_2+$i}")
  NEXT_SEQ_ACC_3=$(awk "BEGIN {print $CURRENT_SEQUENCE_INT_ACC_3+$i}")
  NEXT_SEQ_ACC_4=$(awk "BEGIN {print $CURRENT_SEQUENCE_INT_ACC_4+$i}")

  WC_HEIGHT=$(awk "BEGIN {print $i+1}")

  WC1_HASH="0x$(gen_hash)"
  echo "${T1_ACC} submit WC ${WC_HEIGHT} - ${WC1_HASH}"
  ${UNDCLI_BIN} tx wrkchain record 1 --wc_height=${WC_HEIGHT} --block_hash="${WC1_HASH}" --parent_hash="${WC1_P_HASH}" --hash1="0x$(gen_hash)" --hash2="0x$(gen_hash)" --hash3="0x$(gen_hash)" --from=${T1_ACC} $(get_base_flags) --sequence=${NEXT_SEQ_ACC_1}
  WC1_P_HASH=${WC1_HASH}
  echo "${T2_ACC} submit BEACON hash ${SUB_NUM_OF} / ${NUM_TO_SUB}"
  ${UNDCLI_BIN} tx beacon record 1 --hash="$(gen_hash)" --subtime=$(date +%s) --from=${T2_ACC} $(get_base_flags sync) --sequence=${NEXT_SEQ_ACC_2}
  echo "${T3_ACC} submit BEACON hash ${SUB_NUM_OF} / ${NUM_TO_SUB}"
  ${UNDCLI_BIN} tx beacon record 2 --hash="$(gen_hash)" --subtime=$(date +%s) --from=${T3_ACC} $(get_base_flags sync) --sequence=${NEXT_SEQ_ACC_3}
  WC2_HASH="$(gen_hash 1)"
  echo "${T4_ACC} submit WC ${WC_HEIGHT} - ${WC2_HASH}"
  ${UNDCLI_BIN} tx wrkchain record 2 --wc_height=${WC_HEIGHT} --block_hash="${WC2_HASH}" --parent_hash="${WC2_P_HASH}" --from=${T4_ACC} $(get_base_flags) --sequence=${NEXT_SEQ_ACC_4}
  WC2_P_HASH=${WC2_HASH}
done

echo "Done. Wait for approx. 1 block and query"

sleep 6s
${UNDCLI_BIN} query wrkchain wrkchain 1 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query wrkchain wrkchain 2 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 1 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 2 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}

echo "Finished in $(($(date +%s) - START_TIME)) seconds."
echo ""
