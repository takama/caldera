PROTOC_IMAGE ?= takama/protoc:v0.1.12

PROTOC_RUN = docker run --rm -v $(shell pwd):/contracts -w /contracts $(PROTOC_IMAGE) -I.

all: generate

generate: info{{[- if .Example ]}} events{{[- end ]}}

info:
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. info/info.proto
	@$(PROTOC_RUN) --grpc-gateway_out=logtostderr=true,paths=source_relative:. info/info.proto
	@$(PROTOC_RUN) --swagger_out=logtostderr=true:. info/info.proto

{{[- if .Example ]}}

request:
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. request/request.proto
	@$(PROTOC_RUN) --grpc-gateway_out=logtostderr=true,paths=source_relative:. request/request.proto
	@$(PROTOC_RUN) --swagger_out=logtostderr=true:. request/request.proto

events: request
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. events/events.proto
	@$(PROTOC_RUN) --grpc-gateway_out=logtostderr=true,paths=source_relative:. events/events.proto
	@$(PROTOC_RUN) --swagger_out=logtostderr=true:. events/events.proto
{{[- end ]}}


.PHONY: all \
	info \
{{[- if .Example ]}}
	request \
	events \
{{[- end ]}}
	generate
