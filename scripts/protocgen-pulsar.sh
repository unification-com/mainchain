#!/usr/bin/env bash

# this script is for generating protobuf files for the new google.golang.org/protobuf API
set -eo pipefail

echo "Cleaning API directory"
(
	cd api
	find ./ -type f \( -iname \*.pulsar.go -o -iname \*.pb.go -o -iname \*.pb.gw.go \) -delete
	find . -empty -type d -delete
	cd ..
)

echo "Generating API module"
(
	cd proto
	buf generate --template buf.gen.pulsar.yaml
)

echo ""
echo "Now run:"
echo ""
echo "   # reclaim ownership via a root container; portable across rootless and rootful docker:"
echo "   docker run --rm -v \$(pwd):/workspace --user 0:0 alpine sh -c \\"
echo "     'chown -R \$(stat -c %u:%g /workspace) /workspace/api'"
echo ""
