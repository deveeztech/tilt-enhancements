export GOBIN ?= $(shell pwd)/target/go/bin

GO_FILES := $(shell \
	find . '(' -path '*/.*' -o -path './vendor' ')' -prune \
	-o -name '*.go' -print | cut -b3-)

.PHONY: build
build:
	go build ./...

.PHONY: install
install:
	go mod download

.PHONY: test
test:
	go test -race ./...

.PHONY: cover
cover:
	go test -coverprofile=cover.out -covermode=atomic -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint:
	gofmt -d -s $(GO_FILES)
	go vet ./...
	go install honnef.co/go/tools/cmd/staticcheck@2023.1.2
	$(GOBIN)/staticcheck ./...