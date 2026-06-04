# Downstream Sync Plan

- generated_by: `goalcli downstream-sync-plan`
- impact_report: `/home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/impact.md`
- downstream_sync_required: `true`
- downstream_release_decision: `required`
- repository_rules_release_decision: `audit_required`
- primary_downstream: `github.com/ZoneCNH/kernel`
- changed_file_count: `3`
- adoption_claim: `not_claimed`

## Impact Categories

| Category | Files |
| --- | ---: |
| `contracts` | 1 |
| `context_runtime` | 1 |
| `governance_registry` | 1 |
| `harness` | 0 |
| `repository_rules` | 0 |
| `generator` | 0 |
| `downstream_context` | 0 |
| `evidence` | 0 |
| `docs` | 0 |
| `other` | 0 |

## Target Plan

| Target | Layer | Priority | Action | Status |
| --- | --- | --- | --- | --- |
| `kernel` | `L0` | `P0` | `primary_sync_required` | `blocked_pending_downstream_workspace` |
| `configx` | `L1` | `P1` | `sync_required` | `blocked_pending_downstream_workspace` |
| `observex` | `L1` | `P1` | `sync_required` | `blocked_pending_downstream_workspace` |
| `testkitx` | `L1` | `P1` | `sync_required` | `blocked_pending_downstream_workspace` |
| `postgresx` | `L2` | `P2` | `sync_required` | `blocked_pending_downstream_workspace` |
| `redisx` | `L2` | `P2` | `sync_required` | `blocked_pending_downstream_workspace` |
| `kafkax` | `L2` | `P2` | `sync_required` | `blocked_pending_downstream_workspace` |
| `natsx` | `L2` | `P2` | `sync_required` | `blocked_pending_downstream_workspace` |
| `taosx` | `L2` | `P2` | `sync_required` | `blocked_pending_downstream_workspace` |
| `ossx` | `L2` | `P2` | `sync_required` | `blocked_pending_downstream_workspace` |
| `clickhousex` | `L2` | `P2` | `sync_required` | `blocked_pending_downstream_workspace` |
| `x.go` | `consumer_only` | `review` | `consumer_only_review_required` | `review_pending_no_standard_write` |

## Sync Commands

### kernel

```bash
scripts/render_template.sh --module-name kernel --module-path github.com/ZoneCNH/kernel --package-name kernel --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kernel
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kernel && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kernel && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kernel && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kernel && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kernel && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kernel && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### configx

```bash
scripts/render_template.sh --module-name configx --module-path github.com/ZoneCNH/configx --package-name configx --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/configx
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/configx && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/configx && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/configx && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/configx && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/configx && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/configx && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### observex

```bash
scripts/render_template.sh --module-name observex --module-path github.com/ZoneCNH/observex --package-name observex --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/observex
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/observex && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/observex && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/observex && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/observex && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/observex && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/observex && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### testkitx

```bash
scripts/render_template.sh --module-name testkitx --module-path github.com/ZoneCNH/testkitx --package-name testkitx --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/testkitx
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/testkitx && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/testkitx && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/testkitx && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/testkitx && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/testkitx && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/testkitx && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### postgresx

```bash
scripts/render_template.sh --module-name postgresx --module-path github.com/ZoneCNH/postgresx --package-name postgresx --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/postgresx
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/postgresx && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/postgresx && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/postgresx && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/postgresx && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/postgresx && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/postgresx && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### redisx

```bash
scripts/render_template.sh --module-name redisx --module-path github.com/ZoneCNH/redisx --package-name redisx --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/redisx
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/redisx && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/redisx && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/redisx && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/redisx && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/redisx && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/redisx && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### kafkax

```bash
scripts/render_template.sh --module-name kafkax --module-path github.com/ZoneCNH/kafkax --package-name kafkax --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kafkax
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kafkax && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kafkax && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kafkax && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kafkax && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kafkax && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/kafkax && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### natsx

```bash
scripts/render_template.sh --module-name natsx --module-path github.com/ZoneCNH/natsx --package-name natsx --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/natsx
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/natsx && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/natsx && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/natsx && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/natsx && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/natsx && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/natsx && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### taosx

```bash
scripts/render_template.sh --module-name taosx --module-path github.com/ZoneCNH/taosx --package-name taosx --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/taosx
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/taosx && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/taosx && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/taosx && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/taosx && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/taosx && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/taosx && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### ossx

```bash
scripts/render_template.sh --module-name ossx --module-path github.com/ZoneCNH/ossx --package-name ossx --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/ossx
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/ossx && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/ossx && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/ossx && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/ossx && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/ossx && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/ossx && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### clickhousex

```bash
scripts/render_template.sh --module-name clickhousex --module-path github.com/ZoneCNH/clickhousex --package-name clickhousex --out /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/clickhousex
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/clickhousex && GOWORK=off go mod tidy
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/clickhousex && GOWORK=off go test ./...
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/clickhousex && GOWORK=off make contracts
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/clickhousex && GOWORK=off make boundary
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/clickhousex && CHECK_STATUS=passed GOWORK=off make evidence
cd /home/xlib-standard/.downstream-sync-plan-test/TestRunDownstreamSyncPlanWritesRequiredPlan/workspace/clickhousex && RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

## Consumer Review

- `x.go`: `consumer_only_review_required` / `review_pending_no_standard_write`.

## Evidence Rules

- This plan is not proof-based adoption.
- Adoption truth remains `.agent/registries/downstream-adoption-status.yaml` and `.agent/evidence/truth-state.yaml`.
- This command must not modify downstream repositories or adoption truth files.
- Generated output is local evidence and is ignored by git.
