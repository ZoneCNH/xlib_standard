.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: race
race:
	go test -race ./...

.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed"; \
		exit 1; \
	fi

.PHONY: integration
integration:
	./scripts/run_integration.sh

.PHONY: docs-check
docs-check:
	./scripts/check_docs.sh

.PHONY: security
security:
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not installed"; \
		exit 1; \
	fi
	./scripts/check_secrets.sh

.PHONY: boundary
boundary:
	./scripts/check_boundary.sh

.PHONY: contracts
contracts:
	./scripts/check_contracts.sh

.PHONY: property
property:
	go test ./... -run 'Test.*Property|Test.*Invariant'

.PHONY: fuzz-smoke
fuzz-smoke:
	./scripts/run_fuzz_smoke.sh

.PHONY: golden
golden:
	go test ./... -run 'Test.*Golden|Test.*Snapshot'

.PHONY: evidence
evidence:
	./scripts/generate_manifest.sh

.PHONY: release-evidence-hash
release-evidence-hash:
	./scripts/hash_release_evidence.sh >/dev/null

.PHONY: release-evidence-check
release-evidence-check:
	RELEASE_EVIDENCE_REQUIRE_PASSED=1 ./scripts/check_release_evidence.sh

.PHONY: release-evidence-checksum-check
release-evidence-checksum-check:
	./scripts/hash_release_evidence.sh --check

.PHONY: require-gowork-off
require-gowork-off:
	@if [ "$(GOWORK)" != "off" ]; then \
		echo "GOWORK=off is required for release targets"; \
		exit 1; \
	fi

.PHONY: ci
ci: fmt vet lint test race boundary security contracts

.PHONY: ci-extended
ci-extended: ci property golden fuzz-smoke

.PHONY: release-check
release-check: require-gowork-off ci integration docs-check
	CHECK_STATUS=passed $(MAKE) evidence
	$(MAKE) release-evidence-hash
	$(MAKE) release-evidence-check
	$(MAKE) release-evidence-checksum-check

.PHONY: release-check-extended
release-check-extended: require-gowork-off ci-extended integration docs-check
	CHECK_STATUS=passed $(MAKE) evidence
	$(MAKE) release-evidence-hash
	$(MAKE) release-evidence-check
	$(MAKE) release-evidence-checksum-check

.PHONY: release-final-check
release-final-check: release-check
	RELEASE_EVIDENCE_REQUIRE_PASSED=1 RELEASE_EVIDENCE_REQUIRE_CLEAN=1 ./scripts/check_release_evidence.sh
	$(MAKE) release-evidence-checksum-check

.PHONY: release-preflight
release-preflight:
	./scripts/check_release_preflight.sh "$(VERSION)"
	GOWORK=off VERSION="$(VERSION)" $(MAKE) release-final-check
