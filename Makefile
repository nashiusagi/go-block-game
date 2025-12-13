GO ?= go
PKG := ./...
BIN := block-game
BIN_DIR := bin

.PHONY: run lint fmt test build

run:
	$(GO) run .

lint:
	$(GO) vet $(PKG)

fmt:
	$(GO) fmt $(PKG)

test:
	$(GO) test $(PKG)

build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BIN) ./...

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

