APP = caldera
RELEASE ?= v0.1.10
RELEASE_DATE = $(shell date +%FT%T%Z)
PROJECT = github.com/takama/caldera

LDFLAGS = "-s -w \
	-X $(PROJECT)/pkg/version.RELEASE=$(RELEASE) \
	-X $(PROJECT)/pkg/version.DATE=$(RELEASE_DATE)"

GO_PKG = $(shell go list $(PROJECT)/pkg/...)

all: run

vendor: bootstrap
	@echo "+ $@"
	@go mod tidy

run: clean build
	@echo "+ $@"
	./${APP}

build: vendor test lint
	@echo "+ $@"
	@go build -a -ldflags $(LDFLAGS) -o $(APP) $(PROJECT)/cmd/caldera

test:
	@echo "+ $@"
	@go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}"go test -race -cover {{.Dir}}"{{end}}' $(GO_PKG) | xargs -L 1 sh -c

fmt:
	@echo "+ $@"
	@go list -f '"gofmt -w -s -l {{.Dir}}"' $(GO_PKG) | xargs -L 1 sh -c

imports:
	@echo "+ $@"
	@go list -f '"goimports -w {{.Dir}}"' ${GO_PKG} | xargs -L 1 sh -c

lint: bootstrap
	@echo "+ $@"
	@golangci-lint run ./...

version:
	@./bumper.sh

clean:
	@rm -f ./${APP}

HAS_LINT := $(shell command -v golangci-lint;)
HAS_IMPORTS := $(shell command -v goimports;)

bootstrap:
ifndef HAS_LINT
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.32.2
endif
ifndef HAS_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif


.PHONY: all \
	vendor \
	run \
	build \
	test \
	fmt \
	imports \
	lint \
	version \
	clean \
	bootstrap
