GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	@go build -ldflags="-s -w" -o '$(GOBIN)/gateway' ./cmd/gateway/main.go || exit
	@go build -ldflags="-s -w" -o '$(GOBIN)/bid' ./cmd/bid/main.go || exit
	@go build -ldflags="-s -w" -o '$(GOBIN)/pending' ./cmd/pending/main.go || exit

run:
	@go build -o '$(GOBIN)/gateway' ./cmd/gateway/main.go
	@'$(GOBIN)/gateway' --config='$(GOBASE)/configs/gateway.yml'

up:
	@docker-compose up -d --build

down:
	@docker-compose down

test:
	@go test -v -count=1 -race -timeout=60s ./...

install-deps:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0 && go mod vendor && go mod verify

lint: install-deps
	@golangci-lint run ./...

deps:
	@go mod tidy && go mod vendor && go mod verify

install:
	@go mod download

generate:
	@go generate ./...
	@find "$(GOBASE)/schema/" -name "*.proto" -type f -exec protoc --proto_path="$(GOBASE)" --micro_out="$(GOBASE)" --go_out="$(GOBASE)" {} \;

.PHONY: build
