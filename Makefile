TOOLS=docker-compose.tools.yaml

.PHONY: mocks
mocks:
	rm -rf ./internal/generated/mocks
	mkdir -p ./internal/generated/mocks

	docker compose -f $(TOOLS) run --rm mockery

.PHONY: lint
lint:
	docker compose -f $(TOOLS) run --rm lint

.PHONY: test
test:
	./scripts/test.sh
