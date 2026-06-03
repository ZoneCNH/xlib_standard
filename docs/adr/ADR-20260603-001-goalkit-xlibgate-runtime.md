# ADR-20260603-001: goalkit v0.1.0 runs through xlibgate and Harness Runtime

Status: Accepted

## Decision

goalkit v0.1.0 will not ship as a standalone CLI. The executable surface is Makefile targets that invoke `xlibgate`, and Harness Runtime remains the policy/control plane.

The source evidence ledger is `.agent/evidence/ledger.jsonl`. Generated packs under `release/evidence/goalkit/` are derived artifacts only.

## Rationale

This keeps the MVA aligned with the existing governance gate architecture, avoids a second command authority, and makes G12-G16 evidence auditable from the same runtime as the rest of xlib-standard.

## Consequences

- G12-G16 equivalents may be command-backed in PR-4, but they are non-blocking until later Harness activation.
- Reports must expose `mva_status: not-complete` until fresh evidence proves the full MVA.
- Future work must not introduce a mandatory external `goalkit` CLI for v0.1.0.

## Rejected

A standalone `goalkit` CLI was rejected because it would bypass the xlibgate executor and Harness control-plane contract for v0.1.0.
