# Reusable governance fragment copied into rendered downstream repositories.
GOALCLI ?= go run ./cmd/goalcli

.PHONY: adoption-check
adoption-check:
	$(GOALCLI) adoption-check --verify
