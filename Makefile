.DEFAULT_GOAL := test

export CGO_ENABLED := 0

.PHONY: test
test:
	@set -e; \
	trap '$(MAKE) test-deps-down' EXIT; \
	$(MAKE) test-deps-up; \
	$(MAKE) test-run

.PHONY: test-deps-up
test-deps-up:
	docker compose up --detach --wait

.PHONY: test-deps-down
test-deps-down:
	docker compose down -v

.PHONY:test-run
test-run:
	@postgres_address="$$(docker compose port postgres 5432)" && \
	PGHOST=localhost PGPORT="$${postgres_address##*:}" \
	PGUSER=test PGPASSWORD=test PGSSLMODE=disable \
	go test -timeout 2m -cover -coverprofile=coverage.txt -covermode=atomic ./...
