# Traceability Matrix

| REQ | Requirement summary | Primary artifacts | Verification/Evidence |
| --- | --- | --- | --- |
| REQ-001 | xlib-standard 身份决策与 README/standard docs 对齐 | README.md; docs/adr/ADR-20260602-001-xlib-standard-role.md; docs/standard/module-boundary.md; docs/standard/repository-roles.md | docs-check; stale-name scan with migration exceptions |
| REQ-002 | 旧名迁移：baselib-template/foundationx 只在迁移上下文出现 | docs/migration/baselib-template-to-xlib-standard.md; docs/adr/ADR-20260602-002-kernel-rename.md | rg scan; module/name migration worker evidence |
| REQ-003 | Core gate 定义 release-final/check/score/preflight | docs/adr/ADR-20260602-003-core-gate.md; .agent/harness.yaml | release-final-check; score gate worker evidence |
| REQ-004 | module path/package/render 迁移 | docs/standard/template-generation-contract.md; docs/downstream-matrix.md | render/integration worker evidence |
| REQ-005 | Full Goal Runtime v3.1 .agent 文件完整 | .agent/goal-runtime.md; object-model; state-machine; traceability; harness; evidence/review/release/retro/patch files | file inventory; docs-check extension |
| REQ-006 | xlibgate/docs-check/score executable gates | .agent/harness.yaml; docs/standard/harness-gates.md | gate implementation worker evidence |
| REQ-007 | release manifest/hash/version Evidence | .agent/release-template.md; docs/release.md | release manifest worker evidence |
| REQ-008 | security and secret policy | docs/standard/security-and-secret-policy.md; docs/xgo-integration-boundary.md | security gate; secret leak scan |
| REQ-009 | downstream matrix and kernel integration | docs/downstream-matrix.md; docs/adr/ADR-20260602-002-kernel-rename.md | downstream integration worker evidence |
| REQ-010 | x.go integration boundary | docs/xgo-integration-boundary.md; docs/standard/module-boundary.md | boundary gate; import scan |
