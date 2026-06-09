.PHONY: fmt vet lint test race contracts boundary render-smoke security ci release-check

## Gate pipeline: fmt → vet → lint → test → race → contracts → boundary → render-smoke → security

fmt:
	gofmt -s -w .

vet:
	GOWORK=off go vet ./...

lint:
	@which golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed, skipping"; exit 0; }
	golangci-lint run ./...

test:
	GOWORK=off go test -count=1 ./...

race:
	GOWORK=off go test -race -count=1 ./...

contracts:
	./scripts/check_contracts.sh

boundary:
	./scripts/check_boundary.sh

render-smoke:
	./scripts/check_rendered_template.sh

security:
	./scripts/check_security.sh

ci: fmt vet lint test race contracts boundary render-smoke security

release-check:
	@echo "Checking release readiness..."
	GOWORK=off go build ./...
	GOWORK=off go test -count=1 ./...
	./scripts/generate_manifest.sh
	@echo "Release check passed"
