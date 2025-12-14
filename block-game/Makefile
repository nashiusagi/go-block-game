GO ?= go
PKG := ./...
BIN := block-game
BIN_DIR := bin

.PHONY: run lint fmt test build

run:
	$(GO) run cmd/main.go

lint:
	$(GO) vet $(PKG)

fmt:
	$(GO) fmt $(PKG)

test:
	$(GO) test $(PKG)

build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BIN) cmd/main.go

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

