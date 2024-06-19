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
	@rm -rf lint.log
	@echo "Checking gofmt"
	@gofmt -d -s $(GO_FILES) 2>&1 | tee lint.log
	@echo "Checking go vet"
	@go vet ./... 2>&1 | tee -a lint.log
	@echo "Checking staticcheck"
	@go install honnef.co/go/tools/cmd/staticcheck@2023.1.2
	@$(GOBIN)/staticcheck ./... 2>&1 |  tee -a lint.log
	@cat lint.log
	@[ ! -s lint.log ]