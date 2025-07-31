.DEFAULT_GOAL := build

fmt:
	go fmt .\...
.PHONY:fmt

lint: fmt
	golint .\...
.PHONY:lint

vet: fmt
	go vet .\...
.PHONY:vet

build: vet
	go build -o a.out ./cmd/server/main.go
.PHONY:build
