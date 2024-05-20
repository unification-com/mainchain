#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./tmp-swagger-gen

cd proto
proto_dirs=$(find ./mainchain -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package $file &>/dev/null; then
      buf generate --template buf.gen.swagger.yaml $file
    fi
  done
done

cd ..

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./client/docs/config.json -o ./tmp-swagger-gen/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true


echo ""
echo "Now run:"
echo ""
echo "   sudo chown -R $(id -u):$(id -g) tmp-swagger-gen"
echo "   cp tmp-swagger-gen/swagger.yaml client/docs/swagger-ui/swagger.yaml"
echo "   make update-swagger-docs"
echo "   rm -rf tmp-swagger-gen"
echo ""

# clean swagger files
#rm -rf ./tmp-swagger-gen
