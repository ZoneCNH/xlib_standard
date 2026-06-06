GOALCLI ?= go run ./cmd/goalcli
L2_MANIFEST ?= .agent/l2-capabilities.yaml
L2_EVIDENCE_DIR ?= .agent/evidence/l2
L2_RELEASE_READINESS ?= $(L2_EVIDENCE_DIR)/release-readiness.json

.PHONY: \
	l2-capability-check \
	l2-contract \
	l2-integration \
	l2-chaos \
	l2-benchmark \
	l2-adoption \
	l2-evidence \
	l2-release-readiness \
	l2-manifest-check \
	l2-contract-placeholder \
	l2-evidence-check \
	l2-release-readiness-check \
	downstream-baseline \
	downstream-adoption \
	p2-runtime-check

l2-capability-check:
	@test -f $(L2_MANIFEST)
	@echo "L2 capability manifest present: $(L2_MANIFEST)"

l2-contract: l2-capability-check
	@go test ./test/contract

l2-integration: l2-capability-check
	@echo "Standards-only integration stage: provider-neutral placeholder, downstream wires services"

l2-chaos: l2-capability-check
	@echo "Standards-only chaos stage: declare fault evidence expectations downstream"

l2-benchmark: l2-capability-check
	@echo "Standards-only benchmark stage: declare performance evidence expectations downstream"

l2-adoption: l2-capability-check
	@echo "Standards-only adoption stage: confirm rollout evidence and compatibility matrix downstream"

l2-evidence:
	@test -d $(L2_EVIDENCE_DIR)
	@echo "L2 evidence directory present: $(L2_EVIDENCE_DIR)"

l2-release-readiness: l2-capability-check l2-evidence
	@test -f $(L2_RELEASE_READINESS) || { echo "ERROR: missing L2 release-readiness evidence: $(L2_RELEASE_READINESS)" >&2; exit 1; }
	@python3 -m json.tool "$(L2_RELEASE_READINESS)" >/dev/null
	@L2_RELEASE_READINESS="$(L2_RELEASE_READINESS)" python3 -c 'import json, os, pathlib, sys; data = json.loads(pathlib.Path(os.environ["L2_RELEASE_READINESS"]).read_text()); required = {"schema_version", "adapter", "target_level", "score", "profiles", "evidence"}; evidence = data.get("evidence"); allowed_levels = {"L2-T0", "L2-T1", "L2-T2", "L2-T3", "L2-T4"}; errors = []; missing = sorted(required - set(data)); missing and errors.append("missing required fields: " + repr(missing)); data.get("target_level") not in allowed_levels and errors.append("target_level must be one of " + repr(sorted(allowed_levels))); (not isinstance(evidence, list) or not evidence) and errors.append("evidence must be a non-empty list"); blocking = [] if not isinstance(evidence, list) else [item for item in evidence if item.get("status") in {"fail", "missing"}]; blocking and errors.append("release evidence contains failing or missing statuses: " + repr(blocking)); errors and (sys.stderr.write("ERROR: invalid L2 release-readiness evidence: " + "; ".join(errors) + "\n"), sys.exit(1))'
	@echo "L2 release-readiness evidence accepted: $(L2_RELEASE_READINESS)"

l2-manifest-check: l2-capability-check

l2-contract-placeholder: l2-contract

l2-evidence-check: l2-evidence

l2-release-readiness-check: l2-release-readiness

downstream-baseline:
	$(GOALCLI) downstream-baseline --dry-run --verify

downstream-adoption:
	$(GOALCLI) downstream-adoption --dry-run --verify

p2-runtime-check: downstream-baseline downstream-adoption
