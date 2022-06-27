#!/bin/bash

#########################################################################
# A script for generating and broadcasting some test transactions for   #
# populating DevNet                                                     #
#                                                                       #
# Note: the script assumes the accounts in Docker/README.md have been   #
#       imported into the keychain, jq is installed, and "make build"   #
#       has been run.                                                   #
#########################################################################

UNDCLI_BIN="./build/und"
DEVNET_RPC_IP="localhost"
DEVNET_RPC_PORT="26661"
DEVNET_RPC_TCP="tcp://${DEVNET_RPC_IP}:${DEVNET_RPC_PORT}"
DEVNET_RPC_HTTP="http://${DEVNET_RPC_IP}:${DEVNET_RPC_PORT}"
CHAIN_ID="FUND-DevNet-2"
BROADCAST_MODE="sync"
GAS_PRICES="25.0nund"
UPPER_CASE_HASH=0

NUM_TO_SUB=$1
if [ -z "$NUM_TO_SUB" ]; then
  NUM_TO_SUB=200
fi

# Account names as imported into undcli keys
ENT_ACC="ent"
ENT_ACC_SEQ=0
USER_ACCS=( "t1" "t2" "t3" "t4" )
TYPES=( "wrkchain" "beacon" "wrkchain" "beacon" )
ACC_SEQUENCESS=( 0 0 0 0 )

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

function get_addr() {
  local ADDR=$(${UNDCLI_BIN} keys show $1 -a)
  echo "${ADDR}"
}

function get_base_flags() {
  local BROADCAST=${1:-$BROADCAST_MODE}
  local FLAGS="--broadcast-mode=${BROADCAST} --output=json --chain-id=${CHAIN_ID} --node=${DEVNET_RPC_TCP} --gas=auto --gas-adjustment=1.5 -y"
  echo "${FLAGS}"
}

function get_gas_flags() {
  local FLAGS="--gas-prices=${GAS_PRICES}"
  echo "${FLAGS}"
}

function check_accounts_exist() {
  if { ${UNDCLI_BIN} keys show $1 2>&1 >&3 3>&- | grep '^' >&2; } 3>&1; then
    echo "${1} does not seem to exist in keyring. Exiting"
    exit 1
  else
    echo "Found ${1} in keyring"
  fi
}

function get_curr_acc_sequence() {
  local ACC=$1
  local CURR=$(${UNDCLI_BIN} query account $(get_addr ${ACC})  --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output=json | jq --raw-output '.sequence')
  local CURR_INT=$(awk "BEGIN {print $CURR}")
  echo "${CURR_INT}"
}

function update_account_sequences() {
  ENT_ACC_SEQ=$(get_curr_acc_sequence "${ENT_ACC}")

  for i in ${!USER_ACCS[@]}
  do
    ACC=${USER_ACCS[$i]}
    ACC_SEQUENCESS[$i]=$(get_curr_acc_sequence "${ACC}")
  done
}

check_accounts_exist ${ENT_ACC}

for i in ${!USER_ACCS[@]}
do
  check_accounts_exist "${USER_ACCS[$i]}"
done

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

update_account_sequences

for i in ${!USER_ACCS[@]}
do
  ACC=${USER_ACCS[$i]}
  IS_WHITELISTED=$(${UNDCLI_BIN} query enterprise whitelisted $(get_addr ${ACC}) --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json)
  if [ "$IS_WHITELISTED" = "true" ]; then
    echo "${ACC} already whitelisted"
  else
    echo "Whitelist  ${ACC} for Enterprise POs"
    ${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${ACC}) --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=${ENT_ACC_SEQ}
    ENT_ACC_SEQ=$(awk "BEGIN {print $ENT_ACC_SEQ+1}")
  fi
done

echo "Done. Wait for approx. 1 block"
sleep 6s

update_account_sequences

for i in ${!USER_ACCS[@]}
do
  ACC=${USER_ACCS[$i]}
  ACC_SEQ=${ACC_SEQUENCESS[$i]}
  echo "${ACC} raise Enterprise POs"
  ${UNDCLI_BIN} tx enterprise purchase 1000000000000000nund --from=${ACC} $(get_base_flags) $(get_gas_flags) --sequence="${ACC_SEQ}"
  ACC_SEQUENCESS[$i]=$(awk "BEGIN {print $ACC_SEQ+1}")
done

echo "Done. Wait for approx. 1 block"
sleep 6s

update_account_sequences

RAISED_POS=$(${UNDCLI_BIN} query enterprise orders --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json)

for row in $(echo "${RAISED_POS}" | jq -r ".purchase_orders[] | @base64"); do
  _jq() {
    echo ${row} | base64 --decode | jq -r ${1}
  }
  POID=$(_jq '.id')
  PO_STATUS=$(_jq '.status')
  echo "Process Enterprise PO ${POID} - status=${PO_STATUS}"
  if [ "$PO_STATUS" = "STATUS_RAISED" ]; then
    ${UNDCLI_BIN} tx enterprise process ${POID} accept --from=${ENT_ACC} $(get_base_flags) $(get_gas_flags) --sequence=${ENT_ACC_SEQ}
    ENT_ACC_SEQ=$(awk "BEGIN {print $ENT_ACC_SEQ+1}")
  fi
done

echo "Done. Wait for approx. 2 blocks"
sleep 15s

update_account_sequences

POS=$(${UNDCLI_BIN} query enterprise orders --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json)

echo ${POS} | jq

for i in ${!USER_ACCS[@]}
do
  ACC=${USER_ACCS[$i]}
  ACC_SEQ=${ACC_SEQUENCESS[$i]}
  TYPE=${TYPES[$i]}
  MONIKER="${TYPE}_${ACC}"
  GEN_HASH="0x$(gen_hash)"
  THING_EXISTS=$(${UNDCLI_BIN} query ${TYPE} search --moniker="${MONIKER}" --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json | jq ".${TYPE}s[]")
  if [ "$THING_EXISTS" = "" ]; then
    echo "Register ${TYPE} for ${ACC}"
    if [ "$TYPE" = "wrkchain" ]; then
      ${UNDCLI_BIN} tx wrkchain register --moniker="${MONIKER}" --genesis="${GEN_HASH}" --name="${MONIKER}" --base="geth" --from=${ACC} $(get_base_flags) --sequence=${ACC_SEQ}
    else
      ${UNDCLI_BIN} tx beacon register --moniker="${MONIKER}" --name="${MONIKER}" --from=${ACC} $(get_base_flags) --sequence=${ACC_SEQ}
    fi
    ACC_SEQUENCESS[$i]=$(awk "BEGIN {print $ACC_SEQ+1}")
  else
    echo "${TYPE} ${MONIKER} already registered"
  fi
done

echo "Done. Wait for approx. 1 block"
sleep 7s

update_account_sequences

for (( i=0; i<$NUM_TO_SUB; i++ ))
do
  for j in ${!USER_ACCS[@]}
  do
    ACC=${USER_ACCS[$j]}
    ACC_SEQ=${ACC_SEQUENCESS[$j]}
    TYPE=${TYPES[$j]}
    MONIKER="${TYPE}_${ACC}"
    RES=""
    TX_HASH=""
    RAW_LOG=""
    ID=$(${UNDCLI_BIN} query ${TYPE} search --moniker="${MONIKER}" --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json | jq -r ".${TYPE}s[0].${TYPE}_id")
    if [ "$TYPE" = "wrkchain" ]; then
      WC_HASH="0x$(gen_hash)"
      WC_HEIGHT=$(awk "BEGIN {print $i+1}")
      echo "record wrkchain block ${WC_HEIGHT} / ${NUM_TO_SUB} for ${MONIKER}"
      RES=$(${UNDCLI_BIN} tx wrkchain record ${ID} --wc_height=${WC_HEIGHT} --block_hash="${WC_HASH}" --from=${ACC} $(get_base_flags) --sequence=${ACC_SEQ})
      RAW_LOG=$(echo ${RES} | jq -r ".raw_log")
      TX_HASH=$(echo ${RES} | jq -r ".txhash")
    else
      B_HASH="$(gen_hash)"
      TS=$(awk "BEGIN {print $i+1}")
      echo "record beacon timestamp block ${TS} / ${NUM_TO_SUB} for ${MONIKER}"
      RES=$(${UNDCLI_BIN} tx beacon record ${ID} --hash="$(gen_hash)" --subtime=$(date +%s) --from=${ACC} $(get_base_flags) --sequence=${ACC_SEQ})
      RAW_LOG=$(echo ${RES} | jq -r ".raw_log")
      TX_HASH=$(echo ${RES} | jq -r ".txhash")
    fi

    if [ "$RAW_LOG" = "unauthorized: signature verification failed; verify correct account sequence and chain-id" ]; then
      echo "ERROR:"
      echo "${RAW_LOG}"
      CURR_ACC_SEQ=$(get_curr_acc_sequence "${ACC}")
      echo "ACC_SEQ=${ACC_SEQ}"
      echo "CURR_ACC_SEQ=${CURR_ACC_SEQ}"
      ACC_SEQUENCESS[$j]=$(get_curr_acc_sequence "${ACC}")
    else
      echo "Submitted in tx ${TX_HASH}"
      ACC_SEQUENCESS[$j]=$(awk "BEGIN {print $ACC_SEQ+1}")
    fi
  done
#  sleep 1s
done

echo "Done. Wait for approx. 1 block and query"

sleep 6s
${UNDCLI_BIN} query wrkchain wrkchain 1 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query wrkchain wrkchain 2 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 1 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 2 --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}

echo "Finished in $(($(date +%s) - START_TIME)) seconds."
echo ""
