#!/bin/bash

########################################################################
# A script for exporting MainNet genesis and merging WrkChain & BEACON #
# data into DevNet genesis for import/export testing                   #
#                                                                      #
# Note: the script assumes MainNet is syncing to $HOME/.und_mainchain  #
# Note: the script assumes go and jq are installed                     #
# Outputs to ./build/genesis_merge                                     #
#                                                                      #
# Run from project root: ./scripts/genesis_devnet_merge.sh             #
########################################################################

echo "Build und binaries"
make build

mkdir -p ./build/genesis_merge
rm -f ./build/genesis_merge/*.json

echo "Copy DevNet genesis"
cp ./Docker/assets/node1/config/genesis.json ./build/genesis_merge/genesis_DevNet.json

echo "Export MainNet genesis"
./build/und export --for-zero-height --home="${HOME}/.und_mainchain" > ./build/genesis_merge/genesis_MainNet_export.json

echo "Merge WrkChain & BEACON data"

jq -M --argfile b ./build/genesis_merge/genesis_MainNet_export.json '.app_state.beacon *= $b.app_state.beacon' ./build/genesis_merge/genesis_DevNet.json > ./build/genesis_merge/genesis_DevNet_w_BEACON.json

jq -M --argfile b ./build/genesis_merge/genesis_MainNet_export.json '.app_state.wrkchain *= $b.app_state.wrkchain' ./build/genesis_merge/genesis_DevNet_w_BEACON.json > ./build/genesis_merge/genesis.json

rm ./build/genesis_merge/genesis_DevNet_w_BEACON.json

echo "Done"
echo "./build/genesis_merge/genesis.json can be copied to ./Docker/assets/node[n]/config/genesis.json for testing"
