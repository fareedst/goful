.PHONY: help build run install test fmt tidy clean

GO       ?= go
BIN_DIR  ?= bin
BINARY   ?= goful
PKG      ?= ./...
MAIN     ?= .

help:
	@echo "Available targets:"
	@echo "  build   - Compile the goful binary into $(BIN_DIR)/$(BINARY)"
	@echo "  run     - Run goful from sources"
	@echo "  install - Install goful into \$$GOBIN (or GOPATH/bin)"
	@echo "  test    - Run the full Go test suite"
	@echo "  fmt     - Format all Go sources with gofmt"
	@echo "  tidy    - Sync go.mod/go.sum with imports"
	@echo "  clean   - Remove build artifacts"

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY) $(MAIN)

run:
	$(GO) run $(MAIN)

install:
	$(GO) install $(MAIN)

test:
	$(GO) test $(PKG)

fmt:
	$(GO) fmt $(PKG)

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BIN_DIR)

