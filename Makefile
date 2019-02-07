SHELL := /bin/bash

export GO111MODULE := on

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build:
	go build ./...

.PHONY: test
test: format unit-test

.PHONY: format
format:
	go mod tidy

.PHONY: unit-test
unit-test: lint
	go test -v -race ./...

.PHONY: lint
lint:
	golangci-lint run --disable-all \
		--build-tags e2e \
		-e 'should have name of the form ErrFoo' -E 'deadcode,govet,golint' \
		./...

.PHONY: coverage
coverage:
	DEBUG=true go test -v -coverpkg ./... -covermode=atomic -tags e2e -coverprofile=coverage.txt -race $(shell go list -tags e2e ./...)

.PHONY: coverage-web
coverage-web: coverage
	go tool cover -html=coverage.txt
