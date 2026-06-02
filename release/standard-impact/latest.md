# Standard Impact Report

- generated_at: `2026-06-02T02:54:42Z`
- downstream_sync_required: `true`
- primary_downstream: `github.com/ZoneCNH/kernel`
- changed_file_count: `21`

## Downstream

- `github.com/ZoneCNH/kernel`
- `github.com/ZoneCNH/configx`
- `github.com/ZoneCNH/observex`
- `github.com/ZoneCNH/testkitx`
- `github.com/ZoneCNH/postgresx`
- `github.com/ZoneCNH/redisx`
- `github.com/ZoneCNH/kafkax`
- `github.com/ZoneCNH/taosx`
- `github.com/ZoneCNH/ossx`
- `github.com/ZoneCNH/clickhousex`
- `github.com/ZoneCNH/x.go`

## docs

- `README.md`
- `docs/adr/3.md`
- `docs/downstream-sync-policy.md`
- `docs/standard/README.md`
- `docs/standard/evidence-protocol.md`
- `docs/standard/harness-gates.md`
- `docs/supply-chain.md`

## contracts

- 无变化

## harness

- `Makefile`
- `cmd/xlibgate/main.go`
- `scripts/check_dependency_diff.sh`
- `scripts/check_docs.sh`
- `scripts/check_standard_impact.sh`

## generator

- 无变化

## evidence

- `internal/tools/releasemanifest/main.go`
- `internal/tools/releasemanifest/main_test.go`
- `release/manifest/template.json`
- `release/standard-impact/latest.md`

## other

- `.github/dependabot.yml`
- `.gitignore`
- `renovate.json`
- `scripts/check_dependency_diff_test.go`
- `scripts/check_standard_impact_test.go`

## Sync Decision

- `downstream-sync-required`
- 原因：contracts、harness、generator 或 evidence 影响面发生变化。
