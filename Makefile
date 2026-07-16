APPNAME=$(shell basename $(shell go list))
VERSION?=snapshot
COMMIT=$(shell git rev-parse --verify HEAD)
DATE?=$(shell date +%FT%T%z)
RELEASE?=0

GOPATH?=$(shell go env GOPATH)
GO_LDFLAGS+=-X main.appName=$(APPNAME)
GO_LDFLAGS+=-X main.buildVersion=$(VERSION)
GO_LDFLAGS+=-X main.buildCommit=$(COMMIT)
GO_LDFLAGS+=-X main.buildDate=$(DATE)
ifeq ($(RELEASE), 1)
	# Strip debug information from the binary
	GO_LDFLAGS+=-s -w
endif
GO_LDFLAGS:=-ldflags="$(GO_LDFLAGS)"

SUITE=*.venom.yml

DIST=dist
# OS/arch pairs to cross-compile for release. `arm` maps to GOARM=7 (armv7).
PLATFORMS=linux/amd64 linux/arm64 linux/386 linux/arm \
	darwin/amd64 darwin/arm64 \
	windows/amd64 windows/arm64 windows/386

.PHONY: default
default: build

GOLANGCILINTVERSION:=1.64.8
GOLANGCILINT=$(GOPATH)/bin/golangci-lint
$(GOLANGCILINT):
	curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v$(GOLANGCILINTVERSION)

VENOMVERSION:=v1.0.1
VENOM=$(GOPATH)/bin/venom
$(VENOM):
	go install github.com/ovh/venom/cmd/venom@$(VENOMVERSION)

GOCOVMERGE=$(GOPATH)/bin/gocovmerge
$(GOCOVMERGE):
	go install github.com/wadey/gocovmerge@latest

.PHONY: build
build:
	go build $(GO_LDFLAGS) -o ./build/$(APPNAME)

# release cross-compiles the binary for every platform in PLATFORMS, packages each
# one as a tar.gz (zip on Windows) alongside README/LICENSE, and writes checksums.
# Usage: make release VERSION=v1.2.3 RELEASE=1
.PHONY: release
release: clean-dist
	@mkdir -p $(DIST)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		archname=$$arch; goarm=; bin=$(APPNAME); \
		if [ "$$arch" = "arm" ]; then archname=armv7; goarm=7; fi; \
		if [ "$$os" = "windows" ]; then bin=$(APPNAME).exe; fi; \
		stage=$(APPNAME)_$(VERSION)_$${os}_$${archname}; \
		echo ">> building $$os/$$arch -> $$stage"; \
		mkdir -p $(DIST)/$$stage; \
		cp README.md LICENSE $(DIST)/$$stage/; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch GOARM=$$goarm \
			go build -trimpath $(GO_LDFLAGS) -o $(DIST)/$$stage/$$bin . || exit 1; \
		if [ "$$os" = "windows" ]; then \
			(cd $(DIST)/$$stage && zip -q -r ../$$stage.zip .); \
		else \
			tar -czf $(DIST)/$$stage.tar.gz -C $(DIST)/$$stage .; \
		fi; \
		rm -rf $(DIST)/$$stage; \
	done
	@cd $(DIST) && (sha256sum * 2>/dev/null || shasum -a 256 *) > checksums.txt
	@echo ">> artifacts written to $(DIST)/"
	@ls -1 $(DIST)

.PHONY: lint
lint: $(GOLANGCILINT)
	$(GOLANGCILINT) run ./...

.PHONY: test
test:
	go test -v -race -coverprofile=build/cover.out ./...

.PHONY: integration
integration: $(VENOM) clean
	go test -race -tags=integration -coverprofile=venom.cover.out -coverpkg="./..." -c . -o ./build/$(APPNAME).test
	$(VENOM) run -v --output-dir=build --format=xml tests/$(SUITE)

.PHONY: coverage
coverage: $(GOCOVMERGE)
	$(GOCOVMERGE) build/*.cover.out > build/coverage.out
	go tool cover -html=build/coverage.out -o build/coverage.html

# Claude Code skill: symlink the versioned skill into the user's skills dir so it
# stays maintained in-repo while being available globally in Claude Code.
CLAUDE_SKILLS_DIR?=$(HOME)/.claude/skills
SKILL_NAME:=adr
SKILL_SRC:=$(CURDIR)/.claude/skills/$(SKILL_NAME)

.PHONY: install-skill
install-skill:
	@mkdir -p $(CLAUDE_SKILLS_DIR)
	@ln -sfn $(SKILL_SRC) $(CLAUDE_SKILLS_DIR)/$(SKILL_NAME)
	@echo "linked $(CLAUDE_SKILLS_DIR)/$(SKILL_NAME) -> $(SKILL_SRC)"

.PHONY: uninstall-skill
uninstall-skill:
	@rm -f $(CLAUDE_SKILLS_DIR)/$(SKILL_NAME)
	@echo "removed $(CLAUDE_SKILLS_DIR)/$(SKILL_NAME)"

.PHONY: clean
clean:
	find ./build -mindepth 1 ! -name .gitkeep -delete

.PHONY: clean-dist
clean-dist:
	rm -rf $(DIST)
