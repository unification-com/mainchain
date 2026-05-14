###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.18.1
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
#protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user $(shell id -u):$(shell id -g) $(protoImageName)

# protoReclaimOwnership: under rootless docker, files written by the protoc/buf
# container come back owned by an unprivileged sub-UID on the host. Run a root
# container with the same bind-mount and chown everything back to whatever uid
# the host sees as the workspace root. Under rootful docker this is a no-op.
# `stat -c "%u:%g" /workspace` reads the bind-mount's apparent ownership which
# matches the invoking host user under both rootless and rootful docker.
define protoReclaimOwnership
@$(DOCKER) run --rm -v $(CURDIR):/workspace --user 0:0 $(protoImageName) sh -c \
	'OWN=$$(stat -c "%u:%g" /workspace); chown -R "$$OWN" /workspace/proto /workspace/api /workspace/github.com 2>/dev/null; true'
endef

proto-all: proto-format proto-lint proto-gen proto-pulsar-gen

proto-gen:
	@echo "Generating Protobuf files"
	@chmod 777 proto && chmod 666 proto/buf.lock
	@mkdir -p github.com && chmod 777 github.com
	@$(protoImage) sh ./scripts/protocgen.sh
	$(protoReclaimOwnership)
	@cp -r github.com/unification-com/mainchain/* ./
	@rm -rf github.com
	@chmod 755 proto && chmod 644 proto/buf.lock

proto-pulsar-gen:
	@echo "Generating Protobuf Pulsar files"
	@chmod 777 proto && chmod 666 proto/buf.lock
	@chmod -R 777 api
	@$(protoImage) sh ./scripts/protocgen-pulsar.sh
	$(protoReclaimOwnership)
	@chmod 755 proto && chmod 644 proto/buf.lock
	@find api -type f -exec chmod 644 {} \; && find api -type d -exec chmod 755 {} \;

proto-format:
	@chmod 777 -R proto
	@$(protoImage) find ./proto -name "*.proto" -exec clang-format -i {} \;
	$(protoReclaimOwnership)
	@find proto -type f \( -name "*.proto" -o -name "*.yaml" -o -name "*.lock" \) -exec chmod 644 {} \;
	@find proto -type d -exec chmod 755 {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

proto-update-deps:
	@echo "Updating Protobuf dependencies"
	$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update


SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

swagger-proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git clone --depth 1 --branch $(COSMOS_SDK_SEM_VERSION) "https://github.com/cosmos/cosmos-sdk.git" && \
	rm -f ./cosmos-sdk/proto/buf.* && \
	mv ./cosmos-sdk/proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git clone --depth 1 --branch $(IBC_GO_SEM_VERSION) "https://github.com/cosmos/ibc-go.git" && \
	rm -f ./ibc-go/proto/buf.* && \
	mv ./ibc-go/proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin main && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL https://raw.githubusercontent.com/cosmos/gogoproto/main/gogoproto/gogo.proto > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL https://raw.githubusercontent.com/cosmos/ics23/master/proto/cosmos/ics23/v1/proofs.proto > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"


proto-swagger-gen:
	@echo
	@echo "=========== Generate Message ============"
	@echo
	@make swagger-proto-download-deps
	./scripts/protoc-swagger-gen.sh

	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
        echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
	@echo
	@echo "=========== Generate Complete ============"
	@echo

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps swagger-proto-download-deps
