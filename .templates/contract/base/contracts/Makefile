INCLUDE_GRPC_PROTO ?= -I. -I $(GOPATH)/src -I /usr/local/include -I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway -I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

all: generate

generate: bootstrap events

request:
	@echo "+ $@"
	@protoc $(INCLUDE_GRPC_PROTO) --go_out=plugins=grpc:$(GOPATH)/src request/request.proto
	@protoc $(INCLUDE_GRPC_PROTO) --grpc-gateway_out=logtostderr=true:$(GOPATH)/src request/request.proto
	@protoc $(INCLUDE_GRPC_PROTO) --swagger_out=logtostderr=true:. request/request.proto

events: request
	@echo "+ $@"
	@protoc $(INCLUDE_GRPC_PROTO) --go_out=plugins=grpc:$(GOPATH)/src events/events.proto
	@protoc $(INCLUDE_GRPC_PROTO) --grpc-gateway_out=logtostderr=true:$(GOPATH)/src events/events.proto
	@protoc $(INCLUDE_GRPC_PROTO) --swagger_out=logtostderr=true:. events/events.proto

HAS_PROTOC := $(shell command -v protoc;)
HAS_PROTOC_GO := $(shell command -v protoc-gen-go;)
HAS_PROTOC_GRPC_GATEWAY := $(shell command -v protoc-gen-grpc-gateway;)
HAS_PROTOC_SWAGGER := $(shell command -v protoc-gen-swagger;)

bootstrap:
ifndef HAS_PROTOC
	@echo "Protobuf compiler should be installed (v >= 3)"
	@echo "Pre-compiled binaries for your platform: https://github.com/google/protobuf/releases"
	@echo "Do not forget to copy content of 'include' directory into '/usr/local/include'"
	@exit 1
endif
ifndef HAS_PROTOC_GO
	go get -u github.com/golang/protobuf/protoc-gen-go
endif
ifndef HAS_PROTOC_GRPC_GATEWAY
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
endif
ifndef HAS_PROTOC_SWAGGER
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
endif

.PHONY: all \
	generate \
	request \
	events \
	bootstrap
