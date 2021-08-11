#!/usr/bin/env bash

# see pre-requests:
# - https://grpc.io/docs/languages/go/quickstart/
# - gocosmos plugin is automatically installed during scaffolding.

set -eo pipefail

go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null
go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc 2>/dev/null

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  buf protoc \
  -I "proto" \
  -I "third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')

  # command to generate gRPC gateway (*.pb.gw.go in respective modules) files
  buf protoc \
  -I "proto" \
  -I "third_party/proto" \
  --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done


# move proto files to the right places
cp -r github.com/unification-com/mainchain/* ./
rm -rf github.com
