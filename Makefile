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
		--skip-files 'helper_test.go' \
		-e 'should have name of the form ErrFoo' -E 'deadcode,govet,golint' \
		./...

.PHONY: coverage
coverage:
	DEBUG=true go test -v -coverpkg ./... -covermode=atomic -coverprofile=coverage.txt -race $(shell go list ./...)

.PHONY: coverage-web
coverage-web: coverage
	go tool cover -html=coverage.txt
