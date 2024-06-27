.PHONY: tools lint test

tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0

lint:
	@golangci-lint run ./...

test:
	@go test ./pkg/analyzer/...