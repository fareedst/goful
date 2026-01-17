.PHONY: help build run install test fmt vet lint tidy clean clean-release release app clean-app

GO        ?= go
BIN_DIR   ?= bin
BINARY    ?= goful
PKG       ?= ./...
MAIN      ?= .
DIST_DIR  ?= dist
RELEASE_PLATFORMS ?= linux/amd64 linux/arm64 darwin/arm64
RELEASE_LDFLAGS   ?= -s -w
SHASUM    ?= shasum -a 256

# macOS .app bundle variables
APP_NAME      ?= PanelDemo
APP_ID        ?= com.goful.paneldemo
APP_DIR       ?= $(DIST_DIR)/$(APP_NAME).app
APP_CONTENTS  ?= $(APP_DIR)/Contents

help:
	@echo "Available targets:"
	@echo "  build   - Compile the goful binary into $(BIN_DIR)/$(BINARY)"
	@echo "  lint    - Run gofmt and go vet"
	@echo "  test    - Run the full Go test suite"
	@echo "  release - Build CGO-disabled binaries for $(RELEASE_PLATFORMS) with checksums (set PLATFORM=os/arch to limit)"
	@echo "  run     - Run goful from sources"
	@echo "  install - Install goful into \$$GOBIN (or GOPATH/bin)"
	@echo "  fmt     - Format all Go sources with gofmt"
	@echo "  vet     - Run go vet on all packages"
	@echo "  tidy    - Sync go.mod/go.sum with imports"
	@echo "  clean   - Remove build artifacts"
	@echo "  clean-release - Remove $(DIST_DIR)/ release artifacts"

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY) $(MAIN)

fmt:
	$(GO) fmt $(PKG)

vet:
	$(GO) vet $(PKG)

lint: fmt vet

run:
	$(GO) run $(MAIN)

install:
	$(GO) install $(MAIN)

test:
	$(GO) test $(PKG)

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BIN_DIR)

clean-release:
	rm -rf $(DIST_DIR)

release:
	@set -euo pipefail; \
	platforms="$(PLATFORM)"; \
	if [ -z "$$platforms" ]; then platforms="$(RELEASE_PLATFORMS)"; fi; \
	mkdir -p $(DIST_DIR); \
	for platform in $$platforms; do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output="$(DIST_DIR)/$(BINARY)_$${os}_$${arch}"; \
		echo "DIAGNOSTIC: [IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX] building $$output"; \
		rm -f $$output $$output.sha256; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 $(GO) build -trimpath -ldflags "$(RELEASE_LDFLAGS)" -o $$output $(MAIN); \
		chmod +x $$output; \
		$(SHASUM) $$output > $$output.sha256; \
		echo "DIAGNOSTIC: [IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX] sha256 $$(cat $$output.sha256)"; \
	done

app:
	@echo "Building $(APP_NAME).app..."
	@mkdir -p $(APP_CONTENTS)/MacOS
	@printf '%s\n' '#!/bin/bash' \
		'logger -t $(APP_NAME) "$$(date +'"'"'%Y-%m-%d %H:%M:%S'"'"') : $$*"' \
		> $(APP_CONTENTS)/MacOS/$(APP_NAME)
	@chmod +x $(APP_CONTENTS)/MacOS/$(APP_NAME)
	@printf '%s\n' \
		'<?xml version="1.0" encoding="UTF-8"?>' \
		'<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' \
		'<plist version="1.0">' \
		'<dict>' \
		'  <key>CFBundleExecutable</key>' \
		'  <string>$(APP_NAME)</string>' \
		'  <key>CFBundleIdentifier</key>' \
		'  <string>$(APP_ID)</string>' \
		'  <key>CFBundleName</key>' \
		'  <string>$(APP_NAME)</string>' \
		'  <key>CFBundlePackageType</key>' \
		'  <string>APPL</string>' \
		'  <key>CFBundleVersion</key>' \
		'  <string>1.0</string>' \
		'  <key>CFBundleDocumentTypes</key>' \
		'  <array>' \
		'    <dict>' \
		'      <key>CFBundleTypeExtensions</key>' \
		'      <array><string>logged</string></array>' \
		'      <key>CFBundleTypeName</key>' \
		'      <string>Logged File</string>' \
		'      <key>CFBundleTypeRole</key>' \
		'      <string>Viewer</string>' \
		'    </dict>' \
		'  </array>' \
		'</dict>' \
		'</plist>' \
		> $(APP_CONTENTS)/Info.plist
	@echo "Created: $(APP_DIR)"
	@echo "Run 'open $(APP_DIR)' to register with Launch Services"

clean-app:
	rm -rf $(APP_DIR)
