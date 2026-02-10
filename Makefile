APP_NAME := termfolio
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)
HOST_KEY ?= .ssh/host_ed25519
CONFIG ?= config.yaml
GO ?= go
GORELEASER ?= goreleaser

.PHONY: help build run test fmt vet tidy keys release-check clean

help:
	@echo "usage: make <target>"
	@echo ""
	@echo "targets:"
	@echo "  help           show available make targets"
	@echo "  build          build local binary at $(BIN)"
	@echo "  run            run server using CONFIG=$(CONFIG)"
	@echo "  test           run go test for all packages"
	@echo "  fmt            format go source files"
	@echo "  vet            run go vet checks"
	@echo "  tidy           sync go modules"
	@echo "  keys           generate ssh host key at HOST_KEY=$(HOST_KEY)"
	@echo "  release-check  validate goreleaser config"
	@echo "  clean          remove build artifacts"

build: | $(BIN_DIR)
	$(GO) build -o $(BIN) .

run:
	$(GO) run . -c $(CONFIG)

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

keys:
	mkdir -p "$(dir $(HOST_KEY))"
	ssh-keygen -t ed25519 -f "$(HOST_KEY)" -N "" -q
	@echo "generated host key at $(HOST_KEY)"

release-check:
	$(GORELEASER) check --config .goreleaser.yaml

clean:
	rm -rf $(BIN_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)
