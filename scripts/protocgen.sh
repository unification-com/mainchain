#!/usr/bin/env bash

set -eo pipefail

cd proto
buf mod update

echo "Generating gogo proto code"

proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
    echo "buf generate --template buf.gen.gogo.yaml "${file}""
    buf generate --template buf.gen.gogo.yaml "${file}"
    fi
  done
done

cd ..

# move proto files to the right places
#cp -r github.com/unification-com/mainchain/* ./
#rm -rf github.com

echo ""
echo "Now run:"
echo ""
echo "   sudo chown -R $(id -u):$(id -g) github.com"
echo "   cp -r github.com/unification-com/mainchain/* ./"
echo "   rm -rf github.com"
echo ""
echo "   make proto-pulsar-gen"
echo "   make proto-swagger-gen"
echo ""
