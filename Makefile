#!/usr/bin/make -f

PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(HOME)/go/bin
SIMAPP = ./utilities/simulation/make

# Docker variables
DOCKER := $(shell which docker)

DOCKER_IMAGE_NAME = persistenceone/persistencesdk
DOCKER_TAG_NAME = latest
DOCKER_CONTAINER_NAME = persistencesdk-container
DOCKER_CMD ?= "/bin/sh"

export GO111MODULE = on

all: build test

# The below include contains the tools and runsim targets.
include utilities/simulation/make/Makefile

########################################
### Build

build: go.sum
	@go build -mod=readonly -o bin/ ./...
.PHONY: build

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

distclean:
	rm -rf \
    gitian-build-darwin/ \
    gitian-build-linux/ \
    gitian-build-windows/ \
    .gitian-builder-cache/
.PHONY: distclean

### Testing

test: test-unit
test-all: test-unit test-ledger-mock test-race test-cover

test-ledger-mock:
	@go test -mod=readonly `go list github.com/cosmos/cosmos-sdk/crypto` -tags='cgo ledger test_ledger_mock'

test-ledger: test-ledger-mock
	@go test -mod=readonly -v `go list github.com/cosmos/cosmos-sdk/crypto` -tags='cgo ledger'

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly $(PACKAGES_NOSIMULATION) -tags='ledger test_ledger_mock'

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race $(PACKAGES_NOSIMULATION)

.PHONY: test test-all test-ledger-mock test-ledger test-unit test-race

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.assetNode/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Genesis=${HOME}/.assetNode/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 50 5 TestAppImportExport

test-sim-after-import: runsim
	@echo "Running application simulation-after-import. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 50 5 TestAppSimulationAfterImport

test-sim-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.assetNode/config/genesis.json will be used."
	@$(BINDIR)/runsim -Genesis=${HOME}/.assetNode/config/genesis.json -SimAppPkg=$(SIMAPP) 400 5 TestFullAppSimulation

test-sim-multi-seed-long: runsim
	@echo "Running long multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 500 50 TestFullAppSimulation

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 50 10 TestFullAppSimulation

test-sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Period=1 -Commit=true -Seed=57 -v -timeout 24h

.PHONY: \
test-sim-nondeterminism \
test-sim-custom-genesis-fast \
test-sim-import-export \
test-sim-after-import \
test-sim-custom-genesis-multi-seed \
test-sim-multi-seed-short \
test-sim-multi-seed-long \
test-sim-benchmark-invariants

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true

test-sim-benchmark:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

test-sim-profile:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out

.PHONY: test-sim-profile test-sim-benchmark

test-cover:
	@export VERSION=$(VERSION); bash -x tests/test_cover.sh
.PHONY: test-cover

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES_NOSIMULATION)
.PHONY: benchmark


# Commands for running docker
#
# Run persistenceCore on docker
# Example Usage:
#   make docker-build   ## Builds persistenceCore binary in 2 stages, 1st builder 2nd Runner
#                          Final image only has the compiled persistenceCore binary
#   make docker-interactive   ## Will start an shell session into the docker container
#                                Access to persistenceCore binary here
#       NOTE: To be used for testing only, since the container will be removed after stopping
#   make docker-run DOCKER_CMD=sleep 10000000 DOCKER_OPTS=-d   ## Will run the container in the background
#       NOTE: Recommeded to use docker commands directly for long running processes
#   make docker-clean  # Will clean up the running container, as well as delete the image
#                        after one is done testing
docker-build:
	${DOCKER} build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} .

docker-build-no-cache:
	${DOCKER} build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} . --no-cache

docker-build-push: docker-build
	${DOCKER} push ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-run:
	${DOCKER} run ${DOCKER_OPTS} --name=${DOCKER_CONTAINER_NAME} ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} ${DOCKER_CMD}

docker-interactive:
	${MAKE} docker-run DOCKER_CMD=/bin/sh DOCKER_OPTS="--rm -it"

docker-clean-container:
	-${DOCKER} stop ${DOCKER_CONTAINER_NAME}
	-${DOCKER} rm ${DOCKER_CONTAINER_NAME}

docker-clean-image:
	-${DOCKER} rmi ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-clean: docker-clean-container docker-clean-image
