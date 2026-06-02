# Full Goal Runtime v3.1

Runtime: `xlib-standard` Full Goal Runtime v3.1.

## Purpose

把 Goal、Spec、Requirement、Acceptance Criteria、Design、ADR、Plan、Task、Test、Evidence、Risk、Decision、Review、Release、Retrospective 和 Patch 连接为同一条可审计链。完成不是 MVA，而是 `DONE with evidence:`：REQ-001..REQ-010 全部关闭、release-final-check 通过、score >= 9.8、kernel downstream integration 通过、manifest 完整。

## Required artifacts

- `.agent/object-model.md`
- `.agent/state-machine.md`
- `.agent/traceability-matrix.md`
- `.agent/harness.yaml`
- `.agent/evidence-protocol.md`
- `.agent/review-template.md`
- `.agent/release-template.md`
- `.agent/retrospective-template.md`
- `.agent/risk-register.md`
- `.agent/decision-log.md`
- `.agent/rollback-protocol.md`
- `.agent/prompt-patches.md`
- `.agent/harness-patches.md`
- `.agent/rule-patches.md`

## Stop condition

Only stop when Evidence proves the current slice is complete or a blocker is recorded with owner, scope and next action. A final goal claim must include commands, outputs, manifest/checksum state, score state and known gaps.
