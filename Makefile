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

.PHONY: default
default: build

GOLANGCILINTVERSION:=1.46.2
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

.PHONY: lint
lint: $(GOLANGCILINT)
	$(GOLANGCILINT) run ./...

.PHONY: test
test:
	go test -v -race -coverprofile=build/cover.out ./...

.PHONY: integration
integration: $(VENOM) clean
	go test -race -coverprofile=venom.cover.out -coverpkg="./..." -c . -o ./build/$(APPNAME).test
	$(VENOM) run -v --output-dir=build --format=xml tests/$(SUITE)

.PHONY: coverage
coverage: $(GOCOVMERGE)
	$(GOCOVMERGE) build/*.cover.out > build/coverage.out
	go tool cover -html=build/coverage.out -o build/coverage.html

.PHONY: clean
clean:
	mv ./build/.gitkeep /tmp/.gitkeep
	rm -rf ./build/*
	mv /tmp/.gitkeep ./build/.gitkeep
