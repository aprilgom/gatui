.PHONY: test lint fix precommit

test:
	go test ./...

lint:
	golangci-lint run ./...

fix:
	go fix ./...

precommit:
	.githooks/pre-commit
