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