# Include variables from the .envrc file
include .envrc

# ========================================= #
# HELPERS
# ========================================= #

## help: print this help message
.PHONY: help
help:	
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ========================================== #
# DEVELOPMENT
# ========================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

# =========================================== #
# QUALITY CONTROL
# =========================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## test: run suite of End-to-End tests with base environment variables
.PHONY: api-test
api-test:
	@echo 'Running full test suite...'
	go test -v ./cmd/api -db-dsn=${TEST_GREENLIGHT_DB}

## test: run suite of integration tests with base environment variables
.PHONY: int-test
int-test:
	@echo 'Running full test suite...'
	go test -v ./internal/data/ -db-dsn=${TEST_GREENLIGHT_DB}