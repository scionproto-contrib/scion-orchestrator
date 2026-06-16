SCION_VERSION := v0.12.0
SCION_REPO    := https://github.com/scionproto/scion.git
SCION_DIR     := dev/scion
RELEASE_DIR   := release

# Relative paths inside SCION_DIR to build
SCION_CMDS := \
	scion/cmd/scion \
	scion-pki/cmd/scion-pki \
	router/cmd/router \
	control/cmd/control \
	daemon/cmd/daemon \
	dispatcher/cmd/dispatcher

RELEASE_TARGETS := \
	release-linux-amd64 \
#	release-linux-arm64 \
	release-darwin-amd64 \
	release-darwin-arm64 \
	release-windows-amd64

# ---------------------------------------------------------------------------
# Internal helper — always invoked via $(MAKE) _build-release GOOS=… GOARCH=…
# Using a sub-make avoids the double-$$ escaping that $(define)+$(eval) requires.
# ---------------------------------------------------------------------------
RELEASE_SUBDIR := leaf-as

.PHONY: _build-release
_build-release:
	@echo "==> $(GOOS)/$(GOARCH)"
	@mkdir -p $(RELEASE_DIR)/$(GOOS)_$(GOARCH)/$(RELEASE_SUBDIR)/bin
	@cd $(SCION_DIR) && for cmd in $(SCION_CMDS); do \
		name=$$(basename $$cmd); \
		if [ "$(GOOS)" = "windows" ]; then name="$${name}.exe"; fi; \
		out="$(CURDIR)/$(RELEASE_DIR)/$(GOOS)_$(GOARCH)/$(RELEASE_SUBDIR)/bin/$$name"; \
		echo "  $$cmd -> $$out"; \
		( cd "$$cmd" && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o "$$out" . ); \
	done
	@if [ "$(GOOS)" = "windows" ]; then \
		CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
			go build -o $(RELEASE_DIR)/$(GOOS)_$(GOARCH)/$(RELEASE_SUBDIR)/scion-orchestrator.exe .; \
	else \
		CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
			go build -o $(RELEASE_DIR)/$(GOOS)_$(GOARCH)/$(RELEASE_SUBDIR)/scion-orchestrator .; \
	fi
	@cp -r $(RELEASE_DIR)/$(RELEASE_SUBDIR)/config \
		$(RELEASE_DIR)/$(GOOS)_$(GOARCH)/$(RELEASE_SUBDIR)/
	@cd $(RELEASE_DIR)/$(GOOS)_$(GOARCH)/$(RELEASE_SUBDIR) && \
		tar -czf ../../scion-orchestrator_$(SCION_VERSION)_$(GOOS)_$(GOARCH).tar.gz \
			scion-orchestrator* bin config && \
		zip -r   ../../scion-orchestrator_$(SCION_VERSION)_$(GOOS)_$(GOARCH).zip \
			scion-orchestrator* bin config
	@echo "  -> $(RELEASE_DIR)/scion-orchestrator_$(SCION_VERSION)_$(GOOS)_$(GOARCH).{tar.gz,zip}"

# ---------------------------------------------------------------------------
# Per-platform targets
# ---------------------------------------------------------------------------
.PHONY: release-linux-amd64
release-linux-amd64: scion-source
	$(MAKE) _build-release GOOS=linux GOARCH=amd64

.PHONY: release-linux-arm64
release-linux-arm64: scion-source
	$(MAKE) _build-release GOOS=linux GOARCH=arm64

.PHONY: release-darwin-amd64
release-darwin-amd64: scion-source
	$(MAKE) _build-release GOOS=darwin GOARCH=amd64

.PHONY: release-darwin-arm64
release-darwin-arm64: scion-source
	$(MAKE) _build-release GOOS=darwin GOARCH=arm64

.PHONY: release-windows-amd64
release-windows-amd64: scion-source
	$(MAKE) _build-release GOOS=windows GOARCH=amd64

# ---------------------------------------------------------------------------
# Top-level targets
# ---------------------------------------------------------------------------
.PHONY: release
release: $(RELEASE_TARGETS)
	@echo "==> Done. Archives in $(RELEASE_DIR)/"

.PHONY: scion-source
scion-source: $(SCION_DIR)/.git

$(SCION_DIR)/.git:
	mkdir -p dev
	git clone $(SCION_REPO) $(SCION_DIR)
	cd $(SCION_DIR) && git checkout $(SCION_VERSION)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o scion-orchestrator .

.PHONY: clean
clean:
	rm -f scion-orchestrator scion-orchestrator.exe

.PHONY: clean-release
clean-release:
	rm -rf $(RELEASE_DIR)

.PHONY: clean-scion
clean-scion:
	rm -rf dev

.PHONY: help
help:
	@echo "Targets:"
	@echo "  release              Build all platforms + create archives"
	@echo "  release-<os>-<arch>  Build a single platform, e.g. release-linux-amd64"
	@echo "  scion-source         Clone SCION $(SCION_VERSION) into $(SCION_DIR)/"
	@echo "  build                Build scion-orchestrator for the host"
	@echo "  clean                Remove local build artefacts"
	@echo "  clean-release        Remove $(RELEASE_DIR)/"
	@echo "  clean-scion          Remove dev/ (cloned SCION source)"
	@echo ""
	@echo "Platforms: $(RELEASE_TARGETS)"
	@echo "SCION version: $(SCION_VERSION)  (override with SCION_VERSION=...)"
