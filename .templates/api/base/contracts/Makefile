PROTOC_IMAGE ?= takama/protoc:v0.3.5

PROTOC_RUN = docker run --rm -v $(shell pwd):/contracts -w /contracts $(PROTOC_IMAGE) -I.

all: generate

generate: clean info{{[- if .Example ]}} events{{[- end ]}}

info:
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:. \
		--openapiv2_out=logtostderr=true:. info/info.proto

{{[- if .Example ]}}

events:
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. \
		--openapiv2_out=logtostderr=true:. events/events.proto
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:. \
		--openapiv2_out=logtostderr=true:. events/public.proto
{{[- end ]}}

clean:
	@rm -rf */*.pb.go
	@rm -rf */*.pb.gw.go
	@rm -rf */*.swagger.json

.PHONY: all \
	generate \
	info \
{{[- if .Example ]}}
	events \
{{[- end ]}}
	clean
