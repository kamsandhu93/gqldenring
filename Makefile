generate:
	go generate ./...

run:
	go run server.go

fmt:
	go fmt ./...

up:
	docker compose --project-name gqldenring up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

lint:
	PWD=$(pwd) docker run -t --rm -v ${PWD}:/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run -v

release-snapshot:
	goreleaser release --snapshot --clean