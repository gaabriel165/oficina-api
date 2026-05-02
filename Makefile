.PHONY: run build test test/unit test/integration test/coverage lint swagger mocks up down logs

run:
	go run ./cmd/api/main.go

build:
	go build -o bin/api ./cmd/api/main.go

test:
	go test ./...

test/unit:
	go test ./internal/application/usecase/... ./internal/domain/entity/...

test/integration:
	go test ./internal/integration/... -timeout 120s

test/coverage:
	go test ./internal/application/usecase/... ./internal/domain/entity/... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

mocks:
	~/go/bin/darwin_amd64/mockery

swagger:
	~/go/bin/swag init -g cmd/api/main.go

scan:
	trivy fs --scanners vuln,secret,misconfig .

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f api
