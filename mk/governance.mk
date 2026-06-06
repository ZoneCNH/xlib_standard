# Reusable governance fragment copied into rendered downstream repositories.
GOALCLI ?= go run ./cmd/goalcli

.PHONY: require-gowork-off
require-gowork-off:
	@if [ "$${GOWORK:-}" != "off" ]; then \
		echo "GOWORK=off is required for governance checks"; \
		exit 1; \
	fi

.PHONY: adoption-check
adoption-check: require-gowork-off
	$(GOALCLI) adoption-check --verify

.PHONY: downstream-baseline
downstream-baseline: require-gowork-off
	$(GOALCLI) downstream-baseline --dry-run --verify

.PHONY: downstream-adoption
downstream-adoption: require-gowork-off
	$(GOALCLI) downstream-adoption --dry-run --verify

.PHONY: p2-runtime-check
p2-runtime-check: require-gowork-off downstream-baseline downstream-adoption adoption-check
