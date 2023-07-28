.PHONY: all
all: lint test

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -coverprofile=coverage.out -covermode=atomic
