LOCAL_BIN := $(CURDIR)/bin
PATH := $(PATH):$(PWD)/bin

.PHONY: bin-deps
bin-deps:
	$(info installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15

.PHONY: generate
generate:
	$(info generating code...)
	@$(LOCAL_BIN)/easyp generate

.PHONY: lint
lint:
	$(info linting proto...)
	@$(LOCAL_BIN)/easyp lint --path api

.PHONY: breaking
breaking:
	$(info backwarding compatibility...)
	@$(LOCAL_BIN)/easyp breaking --against main --path api

run:
	$(info running service-notes...)
	docker compose -f deploy/app/docker-compose.yml up --build -d
	$(info running service-client...)
	docker compose -f deploy/client/docker-compose.yml up --build -d

down:
	$(info stopping service-client...)
	docker compose -f deploy/client/docker-compose.yml down
	$(info stopping service-notes...)
	docker compose -f deploy/app/docker-compose.yml down
