L2_MANIFEST ?= .agent/l2-capabilities.yaml
L2_EVIDENCE_DIR ?= .agent/evidence/l2

.PHONY: l2-manifest-check l2-contract-placeholder l2-evidence-check l2-release-readiness-check

l2-manifest-check:
	@test -f $(L2_MANIFEST)
	@echo "L2 capability manifest present: $(L2_MANIFEST)"

l2-contract-placeholder:
	@echo "Wire testkitx contract packs in downstream L2 repository"

l2-evidence-check:
	@test -d $(L2_EVIDENCE_DIR)
	@echo "L2 evidence directory present: $(L2_EVIDENCE_DIR)"

l2-release-readiness-check: l2-manifest-check l2-evidence-check
	@echo "Run xlibgate release-readiness adjudication downstream"
