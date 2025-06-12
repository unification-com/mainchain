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
UND_HOME="/tmp/.und_mainchain_devnet"
TX_DIR="${UND_HOME}/txs"

NUM_TO_SUB=$1
if [ -z "$NUM_TO_SUB" ]; then
  NUM_TO_SUB=200
fi

# Account names as imported into undcli keys
ENT_ACC="ent"
ENT_ACC_SEQ=0
ENT_ACC_NUM=0
USER_ACCS=( "t1" "t2" "t3" "t4" )
TYPES=( "wrkchain" "beacon" "wrkchain" "beacon" )
WRK_BEAC_IDS=( 0 0 0 0 )
ACC_SEQUENCESS=( 0 0 0 0 )
ACC_NUMS=( 0 0 0 0 )

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
  local ADDR=$(${UNDCLI_BIN} keys show $1 -a --home "${UND_HOME}" --keyring-backend test)
  echo "${ADDR}"
}

function get_base_tx_flags() {
  local BROADCAST=${1:-$BROADCAST_MODE}
  local FLAGS="--output=json --chain-id=${CHAIN_ID} --node=${DEVNET_RPC_TCP} --home ${UND_HOME} --keyring-backend test -y"
  echo "${FLAGS}"
}

function get_gas_flags() {
  local FLAGS="--gas-prices=${GAS_PRICES} --gas auto --gas-adjustment 1.5"
  echo "${FLAGS}"
  echo ""
}

function check_accounts_exist() {
  local KEY_FILE
  KEYFILE="${UND_HOME}/keyring-test/${1}.info"
  if { ${UNDCLI_BIN} keys show $1 --home "${UND_HOME}" --keyring-backend test 2>&1 >&3 3>&- | grep '^' >&2; } 3>&1; then
    echo "${1} does not seem to exist in keyring. Exiting"
    exit 1
  else
    echo "Found ${1} in keyring"
  fi
}

function get_curr_acc_sequence() {
  local ACC=$1
  local CURR=$(${UNDCLI_BIN} query auth account-info $(get_addr ${ACC}) --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output=json | jq --raw-output '.info.sequence')
  local CURR_INT
  if [ "$CURR" = "null" ]; then
    echo "1"
  else
    CURR_INT=$(awk "BEGIN {print $CURR}")
    echo "${CURR_INT}"
  fi
}

function get_curr_acc_number() {
  local ACC=$1
  local CURR=$(${UNDCLI_BIN} query auth account-info $(get_addr ${ACC}) --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output=json | jq --raw-output '.info.account_number')
  local CURR_INT=$(awk "BEGIN {print $CURR}")
  echo "${CURR_INT}"
}

#function update_account_sequences() {
#  ENT_ACC_SEQ=$(get_curr_acc_sequence "${ENT_ACC}")
#
#  for i in ${!USER_ACCS[@]}
#  do
#    ACC=${USER_ACCS[$i]}
#    ACC_SEQUENCESS[$i]=$(get_curr_acc_sequence "${ACC}")
#  done
#}

function update_account_numbers() {
  ENT_ACC_NUM=$(get_curr_acc_number "${ENT_ACC}")

  for i in ${!USER_ACCS[@]}
  do
    ACC=${USER_ACCS[$i]}
    ACC_NUMS[$i]=$(get_curr_acc_number "${ACC}")
  done
}

function sign_tx() {
  local TX_UNSIGNED=${1}
  local TX_SIGNED=${2}
  local FROM=${3}
  local ACC_N=${4}
  local ACC_S=${5}
  ${UNDCLI_BIN} tx sign ${TX_UNSIGNED} --home ${UND_HOME} --chain-id ${CHAIN_ID} --keyring-backend test --from ${FROM} --offline --account-number ${ACC_N} --sequence ${ACC_S} > "${TX_SIGNED}"
}

function broadcast_tx() {
  local TX_SIGNED=${1}
  ${UNDCLI_BIN} tx broadcast ${TX_SIGNED} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output=json --broadcast-mode sync
}

# init the home dir
echo "init test environment"
${UNDCLI_BIN} keys list --home "${UND_HOME}" --keyring-backend test
cp ./Docker/assets/keychain/* "${UND_HOME}"/keyring-test/
mkdir -p ${TX_DIR}

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

update_account_numbers
#update_account_sequences

for i in ${!USER_ACCS[@]}
do
  ACC=${USER_ACCS[$i]}
  IS_WHITELISTED=$(${UNDCLI_BIN} query enterprise whitelisted $(get_addr ${ACC}) --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json | jq -r '.whitelisted')
  if [ "$IS_WHITELISTED" = "true" ]; then
    echo "${ACC} already whitelisted"
  else
    echo "Whitelist  ${ACC} for Enterprise POs"
    TX_JSON_UNSIGNED="${TX_DIR}/wl_ent_${ENT_ACC_SEQ}_unsigned.json"
    TX_JSON_SIGNED="${TX_DIR}/wl_ent_${ENT_ACC_SEQ}_signed.json"
    ${UNDCLI_BIN} tx enterprise whitelist add $(get_addr ${ACC}) --home ${UND_HOME}  --keyring-backend test --from ${ENT_ACC} --chain-id ${CHAIN_ID} --generate-only  --fees 5000000nund > "${TX_JSON_UNSIGNED}"

    sign_tx "${TX_JSON_UNSIGNED}" "${TX_JSON_SIGNED}" "${ENT_ACC}" "${ENT_ACC_NUM}" "${ENT_ACC_SEQ}"
    broadcast_tx ${TX_JSON_SIGNED}
    ENT_ACC_SEQ=$(awk "BEGIN {print $ENT_ACC_SEQ+1}")
    rm "${TX_JSON_UNSIGNED}"
    rm "${TX_JSON_SIGNED}"
  fi
done

echo "Done. Wait for approx. 1 block"
sleep 6s

#update_account_sequences

for i in ${!USER_ACCS[@]}
do
  ACC=${USER_ACCS[$i]}
  ACC_SEQ=${ACC_SEQUENCESS[$i]}
  echo "${ACC} raise Enterprise POs"
  # no need to generate & sign offline
  ${UNDCLI_BIN} tx enterprise purchase 1199500000000nund --from=${ACC} $(get_base_tx_flags) $(get_gas_flags)
  ACC_SEQUENCESS[$i]=$(awk "BEGIN {print $ACC_SEQ+1}")
done

echo "Done. Wait for approx. 1 block"
sleep 6s

#update_account_sequences

RAISED_POS=$(${UNDCLI_BIN} query enterprise orders --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json)

for row in $(echo "${RAISED_POS}" | jq -r ".purchase_orders[] | @base64"); do
  _jq() {
    echo ${row} | base64 --decode | jq -r ${1}
  }
  POID=$(_jq '.id')
  PO_STATUS=$(_jq '.status')
  echo "Process Enterprise PO ${POID} - status=${PO_STATUS}"
  if [ "$PO_STATUS" = "STATUS_RAISED" ]; then
    TX_JSON_UNSIGNED="${TX_DIR}/wl_ent_${ENT_ACC_SEQ}_unsigned.json"
    TX_JSON_SIGNED="${TX_DIR}/wl_ent_${ENT_ACC_SEQ}_signed.json"
    ${UNDCLI_BIN} tx enterprise process ${POID} accept --from=${ENT_ACC} --home ${UND_HOME}  --keyring-backend test --from ${ENT_ACC} --chain-id ${CHAIN_ID} --generate-only --fees 5000000nund > "${TX_JSON_UNSIGNED}"

    sign_tx "${TX_JSON_UNSIGNED}" "${TX_JSON_SIGNED}" "${ENT_ACC}" "${ENT_ACC_NUM}" "${ENT_ACC_SEQ}"
    broadcast_tx ${TX_JSON_SIGNED}

    ENT_ACC_SEQ=$(awk "BEGIN {print $ENT_ACC_SEQ+1}")
    rm "${TX_JSON_UNSIGNED}"
    rm "${TX_JSON_SIGNED}"
  fi
done

echo "Done. Wait for approx. 2 blocks"
sleep 15s

#update_account_sequences

POS=$(${UNDCLI_BIN} query enterprise orders --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json)

echo ${POS} | jq


for i in ${!USER_ACCS[@]}
do
  ACC=${USER_ACCS[$i]}
  ACC_SEQ=${ACC_SEQUENCESS[$i]}
  TYPE=${TYPES[$i]}
  MONIKER="${TYPE}_${ACC}"
  GEN_HASH="0x$(gen_hash)"
  THING_EXISTS=$(${UNDCLI_BIN} query ${TYPE} search --moniker="${MONIKER}" --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json | jq ".${TYPE}s[]")
  if [ "$THING_EXISTS" = "" ]; then
    echo "Register ${TYPE} for ${ACC}"
    if [ "$TYPE" = "wrkchain" ]; then
      ${UNDCLI_BIN} tx wrkchain register --moniker="${MONIKER}" --genesis="${GEN_HASH}" --name="${MONIKER}" --base="geth" --from=${ACC} $(get_base_tx_flags)
    else
      ${UNDCLI_BIN} tx beacon register --moniker="${MONIKER}" --name="${MONIKER}" --from=${ACC} $(get_base_tx_flags)
    fi
    ACC_SEQUENCESS[$i]=$(awk "BEGIN {print $ACC_SEQ+1}")
  else
    echo "${TYPE} ${MONIKER} already registered"
  fi
done

echo "Done. Wait for approx. 1 block"
sleep 7s

#update_account_sequences

for (( i=0; i<$NUM_TO_SUB; i++ ))
do
  for j in ${!USER_ACCS[@]}
  do
    ACC=${USER_ACCS[$j]}
    ACC_SEQ=${ACC_SEQUENCESS[$j]}
    ACC_NUM=${ACC_NUMS[$j]}
    TYPE=${TYPES[$j]}
    MONIKER="${TYPE}_${ACC}"
    ID=${WRK_BEAC_IDS[$j]}
    TX_JSON_UNSIGNED="${TX_DIR}/${ACC}_${ACC_SEQ}_unsigned.json"
    TX_JSON_SIGNED="${TX_DIR}/${ACC}_${ACC_SEQ}_signed.json"
    if [ "$ID" = "0" ]; then
      ID=$(${UNDCLI_BIN} query ${TYPE} search --moniker="${MONIKER}" --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID} --output json | jq -r ".${TYPE}s[0].${TYPE}_id")
      WRK_BEAC_IDS[$j]="${ID}"
    fi

    if [ "$TYPE" = "wrkchain" ]; then
      WC_HASH="0x$(gen_hash)"
      WC_HEIGHT=$(awk "BEGIN {print $i+1}")
      echo "record wrkchain block ${WC_HEIGHT} / ${NUM_TO_SUB} for ${MONIKER}"
      ${UNDCLI_BIN} tx wrkchain record ${ID} --wc_height=${WC_HEIGHT} --block_hash="${WC_HASH}" --home ${UND_HOME} --keyring-backend test --from ${ACC} --node=${DEVNET_RPC_TCP} --chain-id ${CHAIN_ID} --generate-only  --fees 1000000000nund > "${TX_JSON_UNSIGNED}"
    else
      B_HASH="$(gen_hash)"
      TS=$(awk "BEGIN {print $i+1}")
      echo "record beacon timestamp block ${TS} / ${NUM_TO_SUB} for ${MONIKER}"
      ${UNDCLI_BIN} tx beacon record ${ID} --hash="$(gen_hash)" --subtime=$(date +%s) --home ${UND_HOME} --keyring-backend test --from ${ACC} --node=${DEVNET_RPC_TCP} --chain-id ${CHAIN_ID} --generate-only  --fees 1000000000nund > "${TX_JSON_UNSIGNED}"
    fi

    sign_tx "${TX_JSON_UNSIGNED}" "${TX_JSON_SIGNED}" "${ACC}" "${ACC_NUM}" "${ACC_SEQ}"
    broadcast_tx "${TX_JSON_SIGNED}"
    ACC_SEQUENCESS[$j]=$(awk "BEGIN {print $ACC_SEQ+1}")
    rm "${TX_JSON_UNSIGNED}"
    rm "${TX_JSON_SIGNED}"
  done
#  sleep 1s
done

echo "Done. Wait for approx. 1 block and query"

sleep 6s
${UNDCLI_BIN} query wrkchain wrkchain 1 --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query wrkchain wrkchain 2 --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 1 --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}
${UNDCLI_BIN} query beacon beacon 2 --home ${UND_HOME} --node=${DEVNET_RPC_TCP} --chain-id=${CHAIN_ID}

echo "Finished in $(($(date +%s) - START_TIME)) seconds."
echo ""
