GIN_MODE=$(shell printenv GIN_MODE)

integration-test:
	@GIN_MODE=test ginkgo --randomizeAllSpecs -v ./tests/...

migrate:
	@tern migrate -m ./migrations -c ./tern$(GIN_MODE).conf

restart-dbs:
	@sudo service postgresql restart
	@sudo service redis-server restart

unit-test:
	go test

