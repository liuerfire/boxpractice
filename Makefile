.PHONY: all
all: build

build:
	go build -o target/boxpractice ./cmd/boxpractice

TEST_FLAGS = -race
.PHONY: unittest
unittest:
	./scripts/migrate drop -f
	./scripts/migrate up
	go test $(TEST_FLAGS) -coverprofile=cover.out -p 1 ./...
	go tool cover -func=cover.out

.PHONY: format
format: $(GOIMPORTS)
	go install golang.org/x/tools/cmd/goimports@latest
	go list -f '{{.Dir}}' ./... | xargs goimports -w -local github.com/liuerfire

.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
	golangci-lint run

mysql:
	docker run --rm -d -p 3306:3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=1 -e MYSQL_DATABASE=boxpractice mysql:8

run:
	./scripts/migrate up
	./target/boxpractice
