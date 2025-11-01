-include app.env

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ''

## run: Start the Go server
.PHONY: run
run:
	@echo starting the Go server
	go run cmd/main.go

## migrations name=<name>: Create a new migration
.PHONY: new-migration
migrations:
	@echo creating migration file
	@cd internal/infra/database/migrations && tern new $(name)

## migrations-up: Apply up migrations
.PHONY: migrations-up
migrations-up:
	@echo "Running up migrations..."
	tern migrate --migrations internal/infra/database/migrations

## migrations-down: Roll back the last applied migration
.PHONY: migrations-down
migrations-down:
	@echo "Running down migrations..."
	tern migrate \
		--conn-string "postgres://admin:secret@localhost:5432/merchcore?sslmode=disable" \
		--migrations internal/infra/database/migrations \
		--destination 0

## migrations-force version=<version>: Force migrations to a version
.PHONY: migrations-force
migrations-force:
	@echo "Forcing migration version $(version)"
	migrate -database=$(DSN) -path=./internal/db/migrations force $(version)

## test: Run all unit tests
.PHONY: test
test:
	@echo running all unit tests
	go test -v -cover -count=1 ./...

## audit: Format, lint, test
.PHONY: audit
audit:
	@echo 'Formatting...'
	go fmt ./...

	@echo 'Linting...'
	golangci-lint run

	# @echo 'Running tests...'
	# go test -race -vet=off ./...

## vendor: Tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo Vendoring...
	go mod tidy
	go mod verify
	go mod vendor

## build: Build the Go binary
.PHONY: build
build:
	@echo Building the Go binary
	go build -o bin/app cmd/main.go

## clean: Remove build artifacts
.PHONY: clean
clean:
	@echo cleaning our binarys
	rm -rf ./bin/*

## mock: filename=<relative path to where you want mock genrated to be stored> interface-name=<iface> Generate mocks
.PHONY: mock
mock:
	mockgen -package mockdb -destination $(filename) build/internal/domain/users $(interface-name)
	mockgen -package mockdb -destination ./internal/domain/user/userdb/mock/user.go build/internal/domain/users UserRepository


## redis: run the redis client   
.PHONY: redis
redis:
	docker run --name redis-bankapp -p 6379:6379 -d  redis:8.2.0-alpine

## compose-build: run compose build
.PHONY: compose-build
compose-build:
	docker compose build --no-cache

## compose-up: run compose up
.PHONY: compose-up 
compose-up:
	docker compose -f docker-compose.yaml up

## compose-debug: run compose debug
.PHONY: compose-debug
compose-debug:
	docker compose -f docker-compose.yaml -f docker-compose-debug.yaml up

## compose-down: run compose down
.PHONY: compose-down
compose-down:
	docker compose down

## compose-test: run compose test
.PHONY: compose-test
compose-test:
	docker compose -f docker-compose.yaml -f docker-compose-test.yaml run --build simplebank


# docker run -d   \
# 	--name redis-merchcore \
# 	-p 6379:6379 \
# 	-v redisdata:/data \
# 	redis:8.2.0-alpine \
# 	redis-server --requirepass redis1234



# docker run -d \
#   --name postgres-merchcore \
#   -e POSTGRES_USER=admin \
#   -e POSTGRES_PASSWORD=secret \
#   -e POSTGRES_DB=merchcore \
#   -p 5432:5432 \
#   -v storehqdata:/var/lib/postgresql/data \
#   postgres:18.0
