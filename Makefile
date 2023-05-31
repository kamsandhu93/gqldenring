SHELL := $(shell which bash)
DOC_COM := docker compose --project-name gqldenring

generate:
	go generate ./...

run:
	go run server.go

run-db:
	SQL_CONN="root:qwerty@tcp(0.0.0.0:3306)/db" go run server.go

fmt:
	goimports -w .

install-goimports:
	go install golang.org/x/tools/cmd/goimports@latest

tidy:
	go mod tidy

up:
	$(DOC_COM) up -d --build --wait

up-db:
	$(DOC_COM) up -d --wait db phpmyadmin

down:
	$(DOC_COM) down --remove-orphans

logs:
	$(DOC_COM) logs -f

lint:
	PWD=$(pwd) docker run -t --rm -v ${PWD}:/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run -v

release-snapshot:
	goreleaser release --snapshot --clean

vet:
	go vet .

build:
	go build -v .

test:
	go test ./...

check-git-diff:
	git diff --compact-summary --exit-code

# Requires all edits to be staged e.g. git add .
ci-checks: fmt tidy lint build test vet
	make check-git-diff || \
        (echo; echo "Unexpected difference in directories after code goimports and go mod tidy. Run the 'make fmt' and 'make tidy' commands then commit."; exit 1)
	@echo "All checks passed \U0001F44D"

