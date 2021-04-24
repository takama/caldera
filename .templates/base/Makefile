PROJECT = {{[ .Project ]}}
APP = {{[ .Name ]}}
BIN = {{[ .Bin ]}}
SERVICE_NAME ?= $(shell echo "$(APP)" | tr - _)

# Use the v0.0.0 tag for testing, it shouldn't clobber any release builds
RELEASE ?= v0.0.0
{{[- if .API.Enabled ]}}
API_VERSION ?= {{[ .API.Version ]}}
{{[- end ]}}
GOOS ?= linux
GOARCH ?= amd64
CA_DIR ?= certs

# Include special environments
# Uncomment this to use your specific envs and secrets
# ifneq ("$(wildcard .env.mk)","")
# include .env.mk
# endif

# Namespace: public, commerce, billing,
NAMESPACE ?= {{[ .Namespace ]}}
# Cluster: dev, prod, ...
CLUSTER ?= dev
VALUES ?= values-$(CLUSTER)

# Certificates and hosts
{{[ toENV .Name ]}}_SSL_CERT_NAME ?= $(APP)-$(CLUSTER)-certs
{{[ toENV .Name ]}}_GRPC_HOST ?= $(APP)-grpc.$(CLUSTER)
{{[ toENV .Name ]}}_REST_HOST ?= $(APP)-rest.$(CLUSTER)
{{[ toENV .Name ]}}_CORS_ALLOWED_HOSTS ?= localhost,127.0.0.1,$({{[ toENV .Name ]}}_REST_HOST)

{{[- if .GKE.Enabled ]}}

# Configs for GKE
GKE_PROJECT_ID ?= {{[ .GKE.Project ]}}
GKE_PROJECT_REGION ?= {{[ .GKE.Region ]}}
GKE_CLUSTER_NAME ?= {{[ .GKE.Cluster ]}}

KUBE_CONTEXT ?= gke_$(GKE_PROJECT_ID)_$(GKE_PROJECT_REGION)_$(GKE_CLUSTER_NAME)

REGISTRY ?= gcr.io/$(GKE_PROJECT_ID)/$(NAMESPACE)
{{[- else ]}}

REGISTRY ?= {{[ .Github ]}}
{{[- end ]}}

CONTAINER_IMAGE ?= $(REGISTRY)/$(APP)
CONTAINER_NAME ?= $(NAMESPACE).$(APP)

{{[- if and .API.Gateway .API.UI ]}}

OPENAPI_GENERATOR_VERSION ?= v0.1.7
{{[- end ]}}
LINTER_VERSION ?= v{{[ .Linter.Version ]}}

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
{{[ toENV .Name ]}}_LOGGER_LEVEL ?= 0
{{[ toENV .Name ]}}_CONFIG_PATH ?= /etc/$(SERVICE_NAME)/default.conf

{{[- if .Storage.Enabled ]}}

# Default database configuration
{{[ toENV .Name ]}}_DB_DRIVER ?= {{[ .Storage.Config.Driver ]}}
{{[ toENV .Name ]}}_DB_DRIVER_VERSION ?= {{[ .Storage.Config.Version ]}}
{{[ toENV .Name ]}}_DB_NAME ?= $(NAMESPACE)-$(APP)
{{[ toENV .Name ]}}_DB_USERNAME ?= {{[ .Storage.Config.Username ]}}
{{[ toENV .Name ]}}_DB_PASSWORD ?= {{[ .Storage.Config.Password ]}}
{{[ toENV .Name ]}}_DB_LOCATION ?= ~/.caldera/$(NAMESPACE)/$(APP)/{{[ .Storage.Config.Driver ]}}
{{[- if .Storage.Postgres ]}}
{{[ toENV .Name ]}}_DB_PATH ?= /var/lib/postgresql/data
{{[- else ]}}
{{[ toENV .Name ]}}_DB_PATH ?= /var/lib/{{[ .Storage.Config.Driver ]}}
{{[- end ]}}
{{[ toENV .Name ]}}_DB_HOST ?= $(CONTAINER_NAME).{{[ .Storage.Config.Host ]}}
{{[ toENV .Name ]}}_DB_PORT ?= {{[ .Storage.Config.Port ]}}
{{[ toENV .Name ]}}_DB_FIXTURES_DIR ?= fixtures

ifneq ("$({{[ toENV .Name ]}}_DB_DRIVER)","stub")
{{[ toENV .Name ]}}_DB_LINK ?= --link $({{[ toENV .Name ]}}_DB_HOST):$({{[ toENV .Name ]}}_DB_HOST)
endif

{{[ toENV .Name ]}}_DB_ENV ?= 

{{[- if .Storage.Postgres ]}}

# PostgreSQL database configuration
ifeq ("$({{[ toENV .Name ]}}_DB_DRIVER)","postgres")
{{[ toENV .Name ]}}_DB_DRIVER_VERSION = 12.6
{{[ toENV .Name ]}}_DB_LOCATION ?= ~/.caldera/$(NAMESPACE)/$(APP)/postgres
{{[ toENV .Name ]}}_DB_PATH = /var/lib/postgresql/data
{{[ toENV .Name ]}}_DB_HOST ?= $(CONTAINER_NAME).postgres
{{[ toENV .Name ]}}_DB_PORT = 5432
{{[ toENV .Name ]}}_DB_ENV = \
	-e "POSTGRES_DB=$({{[ toENV .Name ]}}_DB_NAME)" \
	-e "POSTGRES_USER=$({{[ toENV .Name ]}}_DB_USERNAME)" \
	-e "POSTGRES_PASSWORD=$({{[ toENV .Name ]}}_DB_PASSWORD)"
endif
{{[- end ]}}
{{[- if .Storage.MySQL ]}}

# MySQL database configuration
ifeq ("$({{[ toENV .Name ]}}_DB_DRIVER)","mysql")
{{[ toENV .Name ]}}_DB_DRIVER_VERSION = 8.0
{{[ toENV .Name ]}}_DB_LOCATION ?= ~/.caldera/$(NAMESPACE)/$(APP)/mysql
{{[ toENV .Name ]}}_DB_PATH = /var/lib/mysql
{{[ toENV .Name ]}}_DB_HOST ?= $(CONTAINER_NAME).mysql
{{[ toENV .Name ]}}_DB_PORT = 3306
{{[ toENV .Name ]}}_DB_ROOT_PASSWORD ?= mysql
{{[ toENV .Name ]}}_DB_ENV = \
	-e "MYSQL_ROOT_PASSWORD=$({{[ toENV .Name ]}}_DB_ROOT_PASSWORD)" \
	-e "MYSQL_DATABASE=$({{[ toENV .Name ]}}_DB_NAME)" \
	-e "MYSQL_USER=$({{[ toENV .Name ]}}_DB_USERNAME)" \
	-e "MYSQL_PASSWORD=$({{[ toENV .Name ]}}_DB_PASSWORD)"
endif
{{[- end ]}}
{{[- end ]}}

REPO_INFO = $(shell git config --get remote.origin.url)
REPO_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
RELEASE_DATE = $(shell date +%FT%T%Z)

ifndef REPO_COMMIT
REPO_COMMIT = git-$(shell git rev-parse --short HEAD)
endif

BUILD ?= $(RELEASE)
DEPLOY_PARAMS ?= --atomic
ifneq ("$(CLUSTER)","prod")
BUILD = $(RELEASE)-$(REPO_COMMIT)-$(CLUSTER)
DEPLOY_PARAMS = --wait --force
endif

LDFLAGS = "-s -w \
	-X $(PROJECT)/pkg/version.RELEASE=$(RELEASE) \
{{[- if .API.Enabled ]}}
	-X $(PROJECT)/pkg/version.API=$(API_VERSION) \
{{[- end ]}}
	-X $(PROJECT)/pkg/version.DATE=$(RELEASE_DATE) \
	-X $(PROJECT)/pkg/version.REPO=$(REPO_INFO) \
	-X $(PROJECT)/pkg/version.COMMIT=$(REPO_COMMIT) \
	-X $(PROJECT)/pkg/version.BRANCH=$(REPO_BRANCH)"

GO_PACKAGES = $(shell go list $(PROJECT)/pkg/...)
PRIVATE_REPOS ?= {{[ .PrivateRepo ]}}

BUILDTAGS=

all: build

check-all: fmt imports test lint

{{[- if .GKE.Enabled ]}}

project:
	@echo "+ $@"
ifneq ("$(GKE_PROJECT_ID)", "$(shell gcloud config get-value project)")
	@gcloud config set project $(GKE_PROJECT_ID)
endif
ifneq ("$(GKE_PROJECT_REGION)", "$(shell gcloud config get-value compute/region)")
	@gcloud config set compute/region $(GKE_PROJECT_REGION)
endif
ifneq ("$(GKE_CLUSTER_NAME)", "$(shell gcloud config get-value container/cluster)")
	@gcloud config set container/cluster $(GKE_CLUSTER_NAME)
endif
ifneq ("$(KUBE_CONTEXT)", "$(shell kubectl config view -o jsonpath='{.contexts[?(@.name == "$(KUBE_CONTEXT)")].name}')")
	@gcloud container clusters get-credentials $(GKE_CLUSTER_NAME) --region $(GKE_PROJECT_REGION) --project $(GKE_PROJECT_ID)
endif
ifneq ("$(KUBE_CONTEXT)", "$(shell kubectl config current-context)")
	@kubectl config use-context $(KUBE_CONTEXT)
endif
{{[- end ]}}

{{[- if .API.GRPC ]}}

contracts:
	@echo "+ $@"
	@$(MAKE) -C contracts generate
{{[- if .API.UI ]}}

openapi: bootstrap contracts
	@echo "+ $@"
{{[- if .Example ]}}
	@cp contracts/events/*.swagger.json web/public
{{[- else ]}}
	@cp contracts/info/info.swagger.json web/public/public.swagger.json
{{[- end ]}}
	@statik -m -f -ns public -p openapi -src web/public -dest public
{{[- end ]}}
{{[- end ]}}

vendor: {{[- if .API.UI ]}} openapi {{[- else ]}}{{[- if .API.GRPC ]}} contracts {{[- else ]}} bootstrap {{[- end ]}}{{[- end ]}}
	@echo "+ $@"
ifeq ("$(wildcard go.mod)","")
	@GOPRIVATE=$(PRIVATE_REPOS) go mod init {{[ .Project ]}}
endif
	@GOPRIVATE=$(PRIVATE_REPOS) go mod tidy

compile: vendor test lint
	@echo "+ $@"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo \
		-ldflags $(LDFLAGS) -o bin/$(GOOS)-$(GOARCH)/$(BIN) $(PROJECT)/cmd

certs:
ifeq ("$(wildcard $(CA_DIR)/ca-certificates.crt)","")
	@echo "+ $@"
	@docker run --name $(CONTAINER_NAME).certs -d alpine:latest \
	sh -c "apk --update upgrade && apk add ca-certificates && update-ca-certificates"
	@docker wait $(CONTAINER_NAME).certs
	@mkdir -p $(CA_DIR)
	@docker cp $(CONTAINER_NAME).certs:/etc/ssl/certs/ca-certificates.crt $(CA_DIR)
	@docker rm -f $(CONTAINER_NAME).certs
endif

build: compile certs
	@echo "+ $@"
	@docker build --pull -t $(CONTAINER_IMAGE):$(BUILD) .

push: project build
	@echo "+ $@"
	@docker push $(CONTAINER_IMAGE):$(BUILD)

{{[- if .Storage.Enabled ]}}

HAS_DB_RUNNED := $(shell docker ps -q -f name=^/$({{[ toENV .Name ]}}_DB_HOST)$$)
HAS_DB_EXITED := $(shell docker ps -a -q -f name=^/$({{[ toENV .Name ]}}_DB_HOST)$$)

db:
	@echo "+ $@"
ifneq ("$({{[ toENV .Name ]}}_DB_DRIVER)","stub")
ifndef HAS_DB_RUNNED
ifndef HAS_DB_EXITED
	@mkdir -p $({{[ toENV .Name ]}}_DB_LOCATION)
	@docker run -d	--name $({{[ toENV .Name ]}}_DB_HOST) \
	-p $({{[ toENV .Name ]}}_DB_PORT):$({{[ toENV .Name ]}}_DB_PORT) \
	$({{[ toENV .Name ]}}_DB_ENV) \
	-v $({{[ toENV .Name ]}}_DB_LOCATION):$({{[ toENV .Name ]}}_DB_PATH) \
	$({{[ toENV .Name ]}}_DB_DRIVER):$({{[ toENV .Name ]}}_DB_DRIVER_VERSION)
	@sleep {{[- if .Storage.Postgres ]}} 45 {{[ else ]}} 90 {{[- end ]}}
else
	@docker start $({{[ toENV .Name ]}}_DB_HOST)
endif
endif
endif

migrate-up: build db
	@echo "+ $@"
	@docker run --rm $({{[ toENV .Name ]}}_DB_LINK) \
		-e "{{[ toENV .Name ]}}_LOGGER_LEVEL=$({{[ toENV .Name ]}}_LOGGER_LEVEL)" \
		-e "{{[ toENV .Name ]}}_CONFIG_PATH=$({{[ toENV .Name ]}}_CONFIG_PATH)" \
		-e "{{[ toENV .Name ]}}_DB_DRIVER=$({{[ toENV .Name ]}}_DB_DRIVER)" \
		-e "{{[ toENV .Name ]}}_DB_HOST=$({{[ toENV .Name ]}}_DB_HOST)" \
		-e "{{[ toENV .Name ]}}_DB_PORT=$({{[ toENV .Name ]}}_DB_PORT)" \
		-v $(shell pwd)/config/default.conf:$({{[ toENV .Name ]}}_CONFIG_PATH):ro \
		$(CONTAINER_IMAGE):$(BUILD) /$(BIN) migrate up

migrate-down: build db
	@echo "+ $@"
	@docker run --rm $({{[ toENV .Name ]}}_DB_LINK) \
		-e "{{[ toENV .Name ]}}_LOGGER_LEVEL=$({{[ toENV .Name ]}}_LOGGER_LEVEL)" \
		-e "{{[ toENV .Name ]}}_CONFIG_PATH=$({{[ toENV .Name ]}}_CONFIG_PATH)" \
		-e "{{[ toENV .Name ]}}_DB_DRIVER=$({{[ toENV .Name ]}}_DB_DRIVER)" \
		-e "{{[ toENV .Name ]}}_DB_HOST=$({{[ toENV .Name ]}}_DB_HOST)" \
		-e "{{[ toENV .Name ]}}_DB_PORT=$({{[ toENV .Name ]}}_DB_PORT)" \
		-v $(shell pwd)/config/default.conf:$({{[ toENV .Name ]}}_CONFIG_PATH):ro \
		$(CONTAINER_IMAGE):$(BUILD) /$(BIN) migrate down
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
		-e "{{[ toENV .Name ]}}_LOGGER_LEVEL=$({{[ toENV .Name ]}}_LOGGER_LEVEL)" \
		-e "{{[ toENV .Name ]}}_CONFIG_PATH=$({{[ toENV .Name ]}}_CONFIG_PATH)" \
{{[- if .Storage.Enabled ]}}
		-e "{{[ toENV .Name ]}}_DB_DRIVER=$({{[ toENV .Name ]}}_DB_DRIVER)" \
		-e "{{[ toENV .Name ]}}_DB_HOST=$({{[ toENV .Name ]}}_DB_HOST)" \
		-e "{{[ toENV .Name ]}}_DB_PORT=$({{[ toENV .Name ]}}_DB_PORT)" \
{{[- end ]}}
		-v $(shell pwd)/config/default.conf:$({{[ toENV .Name ]}}_CONFIG_PATH):ro \
		-d $(CONTAINER_IMAGE):$(BUILD)
	@sleep 1
	@docker logs $(CONTAINER_NAME)

logs:
	@echo "+ $@"
	@docker logs -f $(CONTAINER_NAME)

deploy: push
	@echo "+ $@"
	@helm upgrade $(APP) -f .helm/$(VALUES).yaml .helm --kube-context $(KUBE_CONTEXT) \
		--namespace $(NAMESPACE) --version=$(RELEASE) -i $(DEPLOY_PARAMS) \
		--set main.image.registry=$(REGISTRY) \
		--set main.image.name=$(APP) \
		--set main.image.tag=$(BUILD) \
		--set main.ingress.tls.secret.name=$({{[ toENV .Name ]}}_SSL_CERT_NAME) \
		--set main.ingress.hosts.grpc.host=$({{[ toENV .Name ]}}_GRPC_HOST) \
		--set main.ingress.hosts.rest.host=$({{[ toENV .Name ]}}_REST_HOST) \
		--set main.ingress.hosts.rest.cors.allow.origin="{$({{[ toENV .Name ]}}_CORS_ALLOWED_HOSTS)}"

charts:
	@echo "+ $@"
	@helm template .helm -n $(APP) --namespace $(NAMESPACE) -f .helm/$(VALUES).yaml \
		--set main.image.registry=$(REGISTRY) \
		--set main.image.name=$(APP) \
		--set main.image.tag=$(BUILD) \
		--set main.ingress.tls.secret.name=$({{[ toENV .Name ]}}_SSL_CERT_NAME) \
		--set main.ingress.hosts.grpc.host=$({{[ toENV .Name ]}}_GRPC_HOST) \
		--set main.ingress.hosts.rest.host=$({{[ toENV .Name ]}}_REST_HOST) \
		--set main.ingress.hosts.rest.cors.allow.origin="{$({{[ toENV .Name ]}}_CORS_ALLOWED_HOSTS)}"

test:
	@echo "+ $@"
	@go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}"go test -race -cover {{.Dir}}"{{end}}' \
		$(GO_PACKAGES) | xargs -L 1 sh -c

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
	@GOPRIVATE=$(PRIVATE_REPOS) golangci-lint run ./...

HAS_RUNNED := $(shell docker ps -q -f name=^/$(CONTAINER_NAME)$$)
HAS_EXITED := $(shell docker ps -a -q -f name=^/$(CONTAINER_NAME)$$)

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
	@rm -f bin/$(GOOS)-$(GOARCH)/$(BIN)

HAS_LINT := $(shell command -v golangci-lint;)
HAS_IMPORTS := $(shell command -v goimports;)
{{[- if .API.UI ]}}
HAS_OPENAPI_GENERATOR := $(shell command -v statik;)
{{[- end ]}}
{{[- if .GKE.Enabled ]}}
HAS_GCLOUD := $(shell command -v gcloud;)
HAS_DOCKER_GCR := $(shell command -v docker-credential-gcr;)
{{[- end ]}}

bootstrap:
ifndef HAS_LINT
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINTER_VERSION)
endif
ifndef HAS_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif
{{[- if .API.UI ]}}
ifndef HAS_OPENAPI_GENERATOR
	go get -u github.com/rakyll/statik@$(OPENAPI_GENERATOR_VERSION)
endif
{{[- end ]}}
{{[- if .GKE.Enabled ]}}
ifndef HAS_GCLOUD
	@echo "gcloud cli utility should be installed"
	@echo "Pre-compiled binaries for your platform:"
	@echo "https://console.cloud.google.com/storage/browser/cloud-sdk-release?authuser=0"
endif
ifdef HAS_GCLOUD
ifndef HAS_DOCKER_GCR
	@gcloud components install docker-credential-gcr -q
	@docker-credential-gcr configure-docker
endif
ifeq ($(findstring *, $(shell gcloud auth list 2>&1 )),)
	@echo "you should be authorized in gcloud with your work account"
	@echo "you can do it with command:"
	@echo "gcloud auth login"
	@exit 1
endif
endif
{{[- end ]}}

.PHONY: all \
{{[- if .GKE.Enabled ]}}
	project \
{{[- end ]}}
{{[- if .API.GRPC ]}}
	contracts \
{{[- if .API.UI ]}}
	openapi \
{{[- end ]}}
{{[- end ]}}
	vendor \
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
