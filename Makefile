include .env
export $(shell sed 's/=.*//' .env | grep -v '^#')
GIN_MODE=$(shell printenv GIN_MODE)
export PATH := /usr/local/go/bin:$(PATH)

build:
	@go build -o bin/api .

dev: migrate-up
	@if command -v air >/dev/null 2>&1; then air; else go run .; fi

integration-test:
	@GIN_MODE=test ginkgo --randomizeAllSpecs -v ./tests/...

migrate-down:
	@if command -v tern>/dev/null 2>&1; then echo ""; else go install github.com/jackc/tern/v2@latest; fi 
	@tern migrate -m ./migrations --conn-string $(DATABASE_URL) -$(ls -1 ./migrations | wc -l)

migrate-up:
	@if command -v tern>/dev/null 2>&1; then echo ""; else go install github.com/jackc/tern/v2@latest; fi 
	@tern migrate -m ./migrations --conn-string $(DATABASE_URL)

unit-test:
	@go test

