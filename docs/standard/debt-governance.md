# Debt Governance Runtime

The debt governance runtime makes technical-debt controls a P0 release gate for xlib-standard. It is intentionally local, deterministic, and dependency-free: policy is declared under `.agent/debt/`, scanners reuse existing repository scripts, and generated evidence is written under ignored `release/debt/` paths.

## Policy surface

- `.agent/debt/rules.yaml` is the release-blocking rule set.
- `.agent/debt/rule-registry.yaml` maps rule identifiers to scanner scopes.
- `.agent/debt/profile.yaml` defines the default release profile and minimum score.
- `.agent/debt/register.md` documents ownership and operational expectations.

P0 debt cannot be excepted. The scanner treats P0 exception markers in `.agent/debt/rules.yaml` as a failure.

## Command surface

Use `GOWORK=off make debt` or `GOWORK=off go run ./cmd/xlibgate debt` for the full debt gate. Focused commands are also available:

- `architecture`
- `domain`
- `docs-drift`
- `dependency-debt`
- `security-debt`
- `testing-debt`
- `implementation-debt`
- `downstream-debt`

`GOWORK=off make debt-evidence` generates:

- `release/debt/latest.json`
- `release/debt/latest.md`
- `release/debt/latest.json.sha256`

These files are generated release evidence and must not be committed.

## Scanner reuse

The debt scanner delegates to existing gates instead of introducing a parallel toolchain:

- Architecture/domain debt reuse `scripts/check_boundary.sh`.
- Docs drift debt reuses `scripts/check_docs.sh`.
- Dependency debt reuses `scripts/check_dependency_diff.sh`.
- Security debt reuses `scripts/check_secrets.sh`.

Any delegated scanner failure fails the debt gate.

## Release integration

Release manifest generation records debt status, score, checksum, and check count. `release-final-check` requires a 9.8 minimum score and verifies debt evidence freshness as part of the release evidence bundle.
