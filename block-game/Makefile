GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
BIN := block-game
BIN_DIR := bin

.PHONY: run lint fmt test build

run:
	$(GO) run cmd/main.go

lint:
	$(GO) vet ./... 

fmt:
	$(GOFMT) -w $(GOFILES)

fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

test:
	$(GO) test ./... 

build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BIN) cmd/main.go

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

