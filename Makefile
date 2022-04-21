.PHONY: build

build:
	go build -v ./cmd/shortener

.DEFAULT_GOAL := build

run:
	go run ./cmd/shortener/main.go
