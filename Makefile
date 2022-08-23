.PHONY: build

build:
	go build -v ./cmd/shortener

.DEFAULT_GOAL := build

run:
	go run -v -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" ./cmd/shortener/main.go
