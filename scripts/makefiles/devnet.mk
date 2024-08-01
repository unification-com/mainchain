### DevNet Docker compositions

devnet:
	docker-compose -f Docker/docker-compose.local.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.local.yml up --build

devnet-down:
	docker-compose -f Docker/docker-compose.local.yml down --remove-orphans

devnet-latest-release:
	@echo "${LATEST_RELEASE}" > ./.vers_docker
	docker-compose -f Docker/docker-compose.upstream.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.upstream.yml up --build

devnet-latest-release-down:
	docker-compose -f Docker/docker-compose.upstream.yml down

devnet-master:
	@echo "master" > ./.vers_docker
	docker-compose -f Docker/docker-compose.upstream.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.upstream.yml up --build

devnet-master-down:
	docker-compose -f Docker/docker-compose.upstream.yml down

.PHONY: devnet devnet-down devnet-latest-release devnet-latest-release-down devnet-master devnet-master-down