
GIN_MODE=$(shell printenv GIN_MODE)

migrate:
	@tern migrate -m ./migrations -c ./tern$(GIN_MODE).conf

unit-test:
	go test

integration-test:
	@GIN_MODE=test ginkgo watch --randomizeAllSpecs -v ./tests/...