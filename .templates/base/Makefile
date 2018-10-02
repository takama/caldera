PROJECT = github.com/{{[ .Github ]}}/{{[ .Name ]}}
APP = {{[ .Name ]}}

# Use the v0.0.0 tag for testing, it shouldn't clobber any release builds
RELEASE ?= v0.0.0
GOOS ?= linux
GOARCH ?= amd64
CA_DIR ?= certs

{{[- if .GKE.Enabled ]}}

# Configs for GKE
GKE_CLUSTER_NAME ?= {{[ .GKE.Project ]}}
GKE_PROJECT_ZONE ?= {{[ .GKE.Zone ]}}
GKE_PROJECT_ID ?= {{[ .GKE.Cluster ]}}

KUBE_CONTEXT ?= gke_$(GKE_PROJECT_ID)_$(GKE_PROJECT_ZONE)_$(GKE_CLUSTER_NAME)

REGISTRY ?= gcr.io/$(GKE_PROJECT_ID)
{{[- else ]}}

REGISTRY ?= {{[ .Github ]}}
{{[- end ]}}

# Common configuration
{{[- if .API.Enabled ]}}
{{[ toENV .Name ]}}_SERVER_PORT ?= {{[ .API.Config.Port ]}}
{{[ toENV .Name ]}}_EXTERNAL_PORT ?= {{[ .API.Config.Port ]}}
{{[- if .API.Gateway ]}}
{{[ toENV .Name ]}}_SERVER_GATEWAY_PORT ?= {{[ .API.Config.Gateway.Port ]}}
{{[ toENV .Name ]}}_GATEWAY_EXTERNAL_PORT ?= {{[ .API.Config.Gateway.Port ]}}
{{[- end ]}}
{{[- end ]}}
{{[ toENV .Name ]}}_INFO_PORT ?= 8080
{{[ toENV .Name ]}}_INFO_EXTERNAL_PORT ?= 8080
{{[ toENV .Name ]}}_LOG_LEVEL ?= 0
{{[ toENV .Name ]}}_CONFIG_PATH ?= /etc/$(shell echo "$(APP)" | tr - _)/default.conf

{{[- if .Storage.Enabled ]}}

# Database configuration
{{[ toENV .Name ]}}_DB_DRIVER ?= {{[ .Storage.Config.Driver ]}}
{{[- if .Storage.Postgres ]}}
{{[ toENV .Name ]}}_DB_DRIVER_VERSION ?= 10.5
{{[- end ]}}
{{[- if .Storage.MySQL ]}}
{{[ toENV .Name ]}}_DB_DRIVER_VERSION ?= 8.0
{{[- end ]}}
{{[ toENV .Name ]}}_DB_NAME ?= {{[ .Storage.Config.Name ]}}
{{[ toENV .Name ]}}_DB_USERNAME ?= {{[ .Storage.Config.Username ]}}
{{[ toENV .Name ]}}_DB_PASSWORD ?= {{[ .Storage.Config.Password ]}}
{{[ toENV .Name ]}}_DB_LOCATION ?= ~/docker/{{[ .Storage.Config.Driver ]}}
{{[- if .Storage.Postgres ]}}
{{[ toENV .Name ]}}_DB_PATH ?= /var/lib/postgresql/data
{{[- else ]}}
{{[ toENV .Name ]}}_DB_PATH ?= /var/lib/{{[ .Storage.Config.Driver ]}}
{{[- end ]}}
{{[ toENV .Name ]}}_DB_HOST ?= {{[ .Storage.Config.Host ]}}
{{[ toENV .Name ]}}_DB_PORT ?= {{[ .Storage.Config.Port ]}}
{{[ toENV .Name ]}}_DB_FIXTURES_DIR ?= fixtures

ifeq ("$({{[ toENV .Name ]}}_DB_DRIVER)","{{[ .Storage.Config.Driver ]}}")
{{[ toENV .Name ]}}_DB_LINK ?= --link $({{[ toENV .Name ]}}_DB_HOST):$({{[ toENV .Name ]}}_DB_HOST)
endif
{{[- end ]}}

# Namespace: dev, prod, username ...
NAMESPACE ?= dev
VALUES ?= values-$(NAMESPACE)

CONTAINER_IMAGE ?= $(REGISTRY)/$(APP)
CONTAINER_NAME ?= $(APP)

REPO_INFO=$(shell git config --get remote.origin.url)
REPO_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
RELEASE_DATE=$(shell date +%FT%T%Z)

ifndef REPO_COMMIT
REPO_COMMIT = git-$(shell git rev-parse --short HEAD)
endif

BUILD =? $(RELEASE)
ifneq ("$(NAMESPACE)","prod")
BUILD = $(RELEASE)-$(REPO_COMMIT)-$(NAMESPACE)
endif

LDFLAGS = "-s -w \
	-X $(PROJECT)/pkg/version.RELEASE=$(RELEASE) \
	-X $(PROJECT)/pkg/version.DATE=$(RELEASE_DATE) \
	-X $(PROJECT)/pkg/version.REPO=$(REPO_INFO) \
	-X $(PROJECT)/pkg/version.COMMIT=$(REPO_COMMIT) \
	-X $(PROJECT)/pkg/version.BRANCH=$(REPO_BRANCH)"

GO_PACKAGES=$(shell go list $(PROJECT)/...)

BUILDTAGS=

all: build

check-all: fmt imports test lint

{{[- if .GKE.Enabled ]}}

project:
	@echo "+ $@"
ifneq ("$(GKE_PROJECT_ID)", "$(shell gcloud config get-value project)")
	@gcloud config set project $(GKE_PROJECT_ID)
endif
ifneq ("$(GKE_PROJECT_ZONE)", "$(shell gcloud config get-value compute/zone)")
	@gcloud config set compute/zone $(GKE_PROJECT_ZONE)
endif
ifneq ("$(GKE_CLUSTER_NAME)", "$(shell gcloud config get-value container/cluster)")
	@gcloud config set container/cluster $(GKE_CLUSTER_NAME)
endif

cluster:
	@echo "+ $@"
ifneq ("$(GKE_CLUSTER_NAME)", "$(shell gcloud container clusters list | grep $(GKE_CLUSTER_NAME) | cut -d ' ' -f1)")
	@gcloud container clusters get-credentials $(GKE_CLUSTER_NAME) --zone $(GKE_PROJECT_ZONE) --project $(GKE_PROJECT_ID)
endif
ifneq ("$(KUBE_CONTEXT)", "$(shell kubectl config current-context)")
	@kubectl config use-context $(KUBE_CONTEXT)
endif
{{[- end ]}}

vendor: bootstrap
	@echo "+ $@"
	@dep ensure -vendor-only

{{[- if .Contract ]}}

contracts:
	@echo "+ $@"
	@$(MAKE) -C contracts generate
{{[- end ]}}

compile: {{[- if .Contract ]}} contracts {{[- end ]}} vendor test lint
	@echo "+ $@"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo \
		-ldflags $(LDFLAGS) -o bin/$(GOOS)-$(GOARCH)/service $(PROJECT)/cmd

certs:
ifeq ("$(wildcard $(CA_DIR)/ca-certificates.crt)","")
	@echo "+ $@"
	@docker run --name $(CONTAINER_NAME)-certs -d alpine:latest \
	sh -c "apk --update upgrade && apk add ca-certificates && update-ca-certificates"
	@docker wait $(CONTAINER_NAME)-certs
	@mkdir -p $(CA_DIR)
	@docker cp $(CONTAINER_NAME)-certs:/etc/ssl/certs/ca-certificates.crt $(CA_DIR)
	@docker rm -f $(CONTAINER_NAME)-certs
endif

build: compile certs
	@echo "+ $@"
	@docker build --pull -t $(CONTAINER_IMAGE):$(BUILD) .

push: build project
	@echo "+ $@"
{{[- if .GKE.Enabled ]}}
	@gcloud docker -- push $(CONTAINER_IMAGE):$(BUILD)
{{[- else ]}}
	@docker push $(CONTAINER_IMAGE):$(BUILD)
{{[- end ]}}

{{[- if .Storage.Enabled ]}}

HAS_DB_RUNNED := $(shell docker ps | grep $({{[ toENV .Name ]}}_DB_HOST))
HAS_DB_EXITED := $(shell docker ps -a | grep $({{[ toENV .Name ]}}_DB_HOST))

db:
	@echo "+ $@"
ifndef HAS_DB_RUNNED
ifndef HAS_DB_EXITED
	@mkdir -p $({{[ toENV .Name ]}}_DB_LOCATION)
	@docker run -d	--name $({{[ toENV .Name ]}}_DB_DRIVER) \
	-p $({{[ toENV .Name ]}}_DB_PORT):$({{[ toENV .Name ]}}_DB_PORT) \
{{[- if .Storage.MySQL ]}}
	-e "MYSQL_ROOT_PASSWORD={{[ randStr ]}}" \
	-e "MYSQL_DATABASE=$({{[ toENV .Name ]}}_DB_NAME)" \
	-e "MYSQL_USER=$({{[ toENV .Name ]}}_DB_USERNAME)" \
	-e "MYSQL_PASSWORD=$({{[ toENV .Name ]}}_DB_PASSWORD)" \
{{[- end ]}}
{{[- if .Storage.Postgres ]}}
	-e "POSTGRES_DB=$({{[ toENV .Name ]}}_DB_NAME)" \
	-e "POSTGRES_USER=$({{[ toENV .Name ]}}_DB_USERNAME)" \
	-e "POSTGRES_PASSWORD=$({{[ toENV .Name ]}}_DB_PASSWORD)" \
{{[- end ]}}
	-v $({{[ toENV .Name ]}}_DB_LOCATION):$({{[ toENV .Name ]}}_DB_PATH) \
	$({{[ toENV .Name ]}}_DB_DRIVER):$({{[ toENV .Name ]}}_DB_DRIVER_VERSION)
	@sleep {{[- if .Storage.Postgres ]}} 45 {{[ else ]}} 90 {{[- end ]}}
else
	@docker start $({{[ toENV .Name ]}}_DB_HOST)
endif
endif

migrate-up: build db
	@echo "+ $@"
	@docker run --rm $({{[ toENV .Name ]}}_DB_LINK) \
		-e "{{[ toENV .Name ]}}_LOG_LEVEL=$({{[ toENV .Name ]}}_LOG_LEVEL)" \
		-e "{{[ toENV .Name ]}}_CONFIG_PATH=$({{[ toENV .Name ]}}_CONFIG_PATH)" \
		-e "{{[ toENV .Name ]}}_DB_DRIVER=$({{[ toENV .Name ]}}_DB_DRIVER)" \
		-e "{{[ toENV .Name ]}}_DB_HOST=$({{[ toENV .Name ]}}_DB_HOST)" \
		-e "{{[ toENV .Name ]}}_DB_PORT=$({{[ toENV .Name ]}}_DB_PORT)" \
		-v $(shell pwd)/config/default.conf:$({{[ toENV .Name ]}}_CONFIG_PATH):ro \
		$(CONTAINER_IMAGE):$(RELEASE) /dcm migrate up

migrate-down: build db
	@echo "+ $@"
	@docker run --rm $({{[ toENV .Name ]}}_DB_LINK) \
		-e "{{[ toENV .Name ]}}_LOG_LEVEL=$({{[ toENV .Name ]}}_LOG_LEVEL)" \
		-e "{{[ toENV .Name ]}}_CONFIG_PATH=$({{[ toENV .Name ]}}_CONFIG_PATH)" \
		-e "{{[ toENV .Name ]}}_DB_DRIVER=$({{[ toENV .Name ]}}_DB_DRIVER)" \
		-e "{{[ toENV .Name ]}}_DB_HOST=$({{[ toENV .Name ]}}_DB_HOST)" \
		-e "{{[ toENV .Name ]}}_DB_PORT=$({{[ toENV .Name ]}}_DB_PORT)" \
		-v $(shell pwd)/config/default.conf:$({{[ toENV .Name ]}}_CONFIG_PATH):ro \
		$(CONTAINER_IMAGE):$(RELEASE) /dcm migrate down
{{[- end ]}}

run: clean build {{[- if .Storage.Enabled ]}} db {{[- end ]}}
	@echo "+ $@"
	@docker run --name $(CONTAINER_NAME) \
{{[- if .Storage.Enabled ]}}
		$({{[ toENV .Name ]}}_DB_LINK) \
{{[- end ]}}
{{[- if .API.Enabled ]}}
		-p $({{[ toENV .Name ]}}_EXTERNAL_PORT):$({{[ toENV .Name ]}}_SERVER_PORT) \
{{[- if .API.Gateway ]}}
		-p $({{[ toENV .Name ]}}_GATEWAY_EXTERNAL_PORT):$({{[ toENV .Name ]}}_SERVER_GATEWAY_PORT) \
{{[- end ]}}
{{[- end ]}}
		-p $({{[ toENV .Name ]}}_INFO_EXTERNAL_PORT):$({{[ toENV .Name ]}}_INFO_PORT) \
{{[- if .API.Enabled ]}}
		-e "{{[ toENV .Name ]}}_SERVER_PORT=$({{[ toENV .Name ]}}_SERVER_PORT)" \
{{[- if .API.Gateway ]}}
		-e "{{[ toENV .Name ]}}_GATEWAY_PORT=$({{[ toENV .Name ]}}_SERVER_GATEWAY_PORT)" \
{{[- end ]}}
{{[- end ]}}
		-e "{{[ toENV .Name ]}}_INFO_PORT=$({{[ toENV .Name ]}}_INFO_PORT)" \
		-e "{{[ toENV .Name ]}}_LOG_LEVEL=$({{[ toENV .Name ]}}_LOG_LEVEL)" \
		-e "{{[ toENV .Name ]}}_CONFIG_PATH=$({{[ toENV .Name ]}}_CONFIG_PATH)" \
		-e "{{[ toENV .Name ]}}_DB_DRIVER=$({{[ toENV .Name ]}}_DB_DRIVER)" \
		-e "{{[ toENV .Name ]}}_DB_HOST=$({{[ toENV .Name ]}}_DB_HOST)" \
		-e "{{[ toENV .Name ]}}_DB_PORT=$({{[ toENV .Name ]}}_DB_PORT)" \
		-v $(shell pwd)/config/default.conf:$({{[ toENV .Name ]}}_CONFIG_PATH):ro \
		-v $(shell pwd)/fixtures:/$({{[ toENV .Name ]}}_DB_FIXTURES_DIR):ro \
		-d $(CONTAINER_IMAGE):$(BUILD)
	@sleep 1
	@docker logs $(CONTAINER_NAME)

logs:
	@echo "+ $@"
	@docker logs -f $(CONTAINER_NAME)

deploy: push cluster
	@echo "+ $@"
	@helm upgrade $(CONTAINER_NAME) -f .helm/$(VALUES).yaml .helm --kube-context $(KUBE_CONTEXT) \
		--namespace $(NAMESPACE) --version=$(RELEASE) --set image.tag=$(BUILD) -i --wait

charts:
	@echo "+ $@"
	@helm template .helm -n ${NAMESPACE} -f .helm/${VALUES}.yaml

test:
	@echo "+ $@"
	@go test -race -cover $(GO_PACKAGES)

cover:
	@echo "+ $@"
	@echo "mode: set" > coverage.txt
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}} && \
		cat {{.Dir}}/.coverprofile | sed 1d >> coverage.txt"{{end}}' $(GO_PACKAGES) | xargs -L 1 sh -c

fmt:
	@echo "+ $@"
	@go list -f '"gofmt -w -s -l {{.Dir}}"' $(GO_PACKAGES) | xargs -L 1 sh -c

imports:
	@echo "+ $@"
	@go list -f '"goimports -w {{.Dir}}"' ${GO_PACKAGES} | xargs -L 1 sh -c

lint: bootstrap
	@echo "+ $@"
	@golangci-lint run --enable-all --skip-dirs vendor ./...

HAS_RUNNED := $(shell docker ps | grep $(CONTAINER_NAME))
HAS_EXITED := $(shell docker ps -a | grep $(CONTAINER_NAME))

stop:
ifdef HAS_RUNNED
	@echo "+ $@"
	@docker stop $(CONTAINER_NAME)
endif

start: stop
	@echo "+ $@"
	@docker start $(CONTAINER_NAME)

rm:
ifdef HAS_EXITED
	@echo "+ $@"
	@docker rm $(CONTAINER_NAME)
endif

version:
	@./bumper.sh

clean: stop rm
	@rm -f bin/$(GOOS)-$(GOARCH)/$(APP)

HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golangci-lint;)
HAS_IMPORTS := $(shell command -v goimports;)
{{[- if .GKE.Enabled ]}}
HAS_GCLOUD := $(shell command -v gcloud;)
{{[- end ]}}

bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifndef HAS_LINT
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif
ifndef HAS_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif
{{[- if .GKE.Enabled ]}}
ifndef HAS_GCLOUD
	@echo "gcloud cli utility should be installed"
	@echo "Pre-compiled binaries for your platform:
	@echo https://console.cloud.google.com/storage/browser/cloud-sdk-release?authuser=0"
	@exit 1
endif
{{[- end ]}}

.PHONY: all \
{{[- if .GKE.Enabled ]}}
	project \
	cluster \
{{[- end ]}}
	vendor \
{{[- if .Contract ]}}
	contracts \
{{[- end ]}}
	compile \
	build \
	certs \
	push \
{{[- if .Storage.Enabled ]}}
	db \
	migrate-up \
	migrate-down \
{{[- end ]}}
	run \
	logs \
	deploy \
	charts \
	test \
	cover \
	fmt \
	lint \
	stop \
	start \
	rm \
	version \
	clean \
	bootstrap
