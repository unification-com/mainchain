#!/bin/bash

set -e

RPC="http://localhost:1317"

# Result formatting
R_DIV="=============================="
R_DIV="${R_DIV}${R_DIV}${R_DIV}"

# Vars
LAST_HEIGHT=0
TOTAL_SUPPLY=0
TOTAL_LOCKED_EFUND=0
TOTAL_SPENT_EFUND=0

function get_total_supply() {
  local TOT_S

  TOT_S=$(curl -s "${RPC}"/cosmos/bank/v1beta1/supply/by_denom?denom=nund | jq -r '.amount.amount')
  if [ "$TOT_S" -gt "0" ]; then
    TOTAL_SUPPLY=$(echo "${TOT_S}" | awk '{printf("%.2f", $1/(1000000000))}')
  fi
}

function get_locked_efund() {
  local TOTAL_L
  TOTAL_L=$(curl -s "${RPC}"/mainchain/enterprise/v1/locked | jq -r '.amount.amount')
  if [ "$TOTAL_L" -gt "0" ]; then
    TOTAL_LOCKED_EFUND=$(echo "${TOTAL_L}" | awk '{printf("%.2f", $1/(1000000000))}')
  fi
}

function get_spent_efund() {
  local TOTAL_SP
  TOTAL_SP=$(curl -s "${RPC}"/mainchain/enterprise/v1/total_spent | jq -r '.amount.amount')
  if [ "$TOTAL_SP" -gt "0" ]; then
    TOTAL_SPENT_EFUND=$(echo "${TOTAL_SP}" | awk '{printf("%.2f", $1/(1000000000))}')
  fi
}

function print_results() {

  printf "Total Locked eFUND    : %'.2f\n" "${TOTAL_LOCKED_EFUND}"
  printf "Total Spent eFUND     : %'.2f\n" "${TOTAL_SPENT_EFUND}"
  printf "Total Supply          : %'.2f\n\n" "${TOTAL_SUPPLY}"

}

function monitor() {
  while true; do
    get_total_supply
    get_locked_efund
    get_spent_efund
    clear
    print_results
    sleep 1
  done
}

# Run
monitor
