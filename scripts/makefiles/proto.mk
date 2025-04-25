###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
#protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user $(shell id -u):$(shell id -g) $(protoImageName)

proto-all: proto-format proto-lint proto-gen

# NOTE: when using rootless docker, this will fail. Before running, run:
#   chmod 777 proto/buf.lock
#   mkdir github.com && chmod 777 github.com
# After running, run:
#   sudo chown -R $(id -u):$(id -g) github.com
#   cp -r github.com/unification-com/mainchain/* ./
#   rm -rf github.com
proto-gen:
	@echo "Generating Protobuf files"
	@chmod 777 proto/buf.lock
	@mkdir github.com && chmod 777 github.com
	@$(protoImage) sh ./scripts/protocgen.sh

proto-pulsar-gen:
	@echo "Generating Protobuf Pulsar files"
	@chmod 777 proto/buf.lock
	@chmod 777 -R api
	@$(protoImage) sh ./scripts/protocgen-pulsar.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

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
