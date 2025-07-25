GOLANGCI_LINT_IMAGE=golangci/golangci-lint:v2.2-alpine
MOCKERY_PACKAGE=github.com/vektra/mockery/v3@v3.5.1

.PHONY: deps
deps:
	docker pull $(GOLANGCI_LINT_IMAGE)
	go install $(MOCKERY_PACKAGE)

.PHONY: mocks
mocks:
	rm -rf ./internal/generated/mocks
	mkdir -p ./internal/generated/mocks

	mockery

.PHONY: lint
lint:
	./scripts/lint.sh $(GOLANGCI_LINT_IMAGE)

.PHONY: test
test:
	./scripts/test.sh
