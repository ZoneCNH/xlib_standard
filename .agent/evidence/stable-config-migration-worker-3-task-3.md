# Stable Config Migration — Worker 3 Task 3 Review Evidence

Task: `3` — review/documentation slice only.
Scope: reconcile `.worktree/stable.md` omissions against current repository facts; no implementation or gate changes.

## Inputs and command evidence

- `.worktree/stable.md` is not available in this checkout. Command: `find /home/xlib-standard/.worktree/workspaces/stable-config-migration-gates -maxdepth 8 -type f -name stable.md -print` returned no paths.
- `.config/` is not available in this checkout. Command: `[ -e .config ] && find .config -maxdepth 3 -type f -print || echo '.config absent'` returned `.config absent`.
- The closest task source is `/home/xlib-standard/.worktree/workspaces/stable-config-migration-gates/.omx/context/stable-config-migration-gates-20260607T030733Z.md`; lines 15-20 state that the current stable gate coverage is too narrow and misses Go, shell, Makefile, templates, ignore files, GitHub, githooks, devcontainer, and tracked `.omx` context surfaces.
- This artifact documents omissions and verification commands only. It is not release-ready evidence.

## A. Legacy reference and gate coverage omissions

Current facts:

- `README.md:5` states old names are allowed only in migration docs and that the default downstream identity is `kernel`.
- `README.md:23` repeats that old names are not the main identity.
- `docs/migration/baselib-template-to-xlib-standard.md:21-28` requires release legacy scans to include source, scripts, `.agent/`, `.xlib/`, `.github/`, `.githooks/`, `.devcontainer/`, ignore files, downstream/templates, and runtime context evidence.
- `docs/migration/baselib-template-to-xlib-standard.md:30-37` records the broader static `rg --hidden` migration scan.
- `docs/migration/baselib-template-to-xlib-standard.md:39-44` requires classifying old names, `.agent/.xlib`, `.config`, GitHub, hooks, devcontainer, ignore files, and generated artifacts.
- `docs/migration/baselib-template-to-xlib-standard.md:46-51` says validation blocks release-ready before `.config` migration.
- `scripts/check_standard_impact.sh:193-253` already models wider impact buckets, including repository rules, context runtime, evidence, generators, harness, docs, downstream sync, and repository-rules audit conditions.
- `scripts/check_rendered_template.sh:109-120` and `scripts/check_rendered_template.sh:216-233` check rendered template requirements and stale legacy/module/package references.

Risk:

- A narrow stable scan that only checks `.agent/.xlib` and markdown/YAML/JSON files can miss executable or emitted legacy references in Go, shell, Makefile, workflow, hook, devcontainer, template, ignore, and runtime context surfaces.
- Some old-name references are legitimate migration context. A release gate must classify references instead of deleting all matches.

Closure item:

- Promote the documented static migration scan and classification rules into an active stable gate before any direct `v1.0.0` or `v1.0.0-rc.1` claim.

## B. Migration map and source-of-truth inventory omissions

Current facts:

- `.agent/index.yaml:3-13` defines the `.agent` logical classification index, authority order, and `physical_migration: true`.
- `.agent/index.yaml:36-42` classifies `.agent/harness/harness.yaml` as a `source_of_truth` machine contract.
- `cmd/goalcli/governance.go:1728-1760` validates `.agent/index.yaml` paths, required fields, duplicate paths, enum values, and unclassified `.agent` files.
- `.agent/registries/physical-migration-manifest.yaml:1-14` describes migration status, compatibility strategy, and the root exception model.
- `.agent/registries/generated-artifacts.yaml:1-62` is the generated-artifact inventory for downstream baselines, release manifests, generated rules, and evidence artifacts.
- `.agent/policies/runtime-file-ownership.yaml:11-31` assigns ownership for `.agent`, iron-rules, registry, and harness paths.
- `.agent/policies/runtime-file-ownership.yaml:52-91` classifies `.githooks`, `.github`, workflows, `.devcontainer`, `.gitignore`, and `.dockerignore`; `.github` is explicitly platform-native and not the config source of truth.
- `internal/xlibfacts/facts.go:15` still hard-codes `.xlib` facts.
- `internal/validation/validation.go:40-42` and `internal/validation/validation.go:89-91` require runtime ownership only for `.agent/`, `cmd/goalcli/`, and `contracts/`.
- `docs/v0.6.0/strict-config-root-omission-audit.md:1-4` concludes strict `.config/` root still has 15 gaps, with 6 P0 gaps.
- `docs/v0.6.0/strict-config-root-omission-audit.md:10` says `.config/` must become the unique machine source and old paths/flags should become violations.
- `docs/v0.6.0/strict-config-root-omission-audit.md:300-352` requires release manifest v2 fields for strict config root, config fingerprints, platform adapter, downstream cutover, and pathguard/strict-check proof.

Risk:

- The repository still treats `.agent` and `.xlib` as active machine fact surfaces while `.config/` is absent, so a stable or release-ready statement would be misleading.
- Runtime ownership policy covers broader surfaces than the current validation hard requirements.

Closure item:

- Add a machine-readable `.config` migration inventory and strict gate that distinguishes canonical config, platform adapters, generated artifacts, and temporary compatibility paths.

## C. GitHub, hooks, devcontainer, ignore, and generated-artifact omissions

Current facts:

- `.github/workflows/ci.yml:37-45` runs `GOWORK=off XLIB_CONTEXT=ci_pull_request make release-check` in CI.
- `.github/workflows/release.yml:37-45` runs release final checks before manifest generation and upload.
- `.github/workflows/release.yml:43-68` handles release manifest hash/upload.
- `.github/workflows/security.yml:31-40` installs and runs the pinned `govulncheck` path.
- `.github/CODEOWNERS:1` is a platform-native ownership file.
- `.github/dependabot.yml:14`, `.github/dependabot.yml:34`, and `.github/dependabot.yml:43-54` define dependency-review labels and Docker ecosystem coverage.
- `.github/pull_request_template.md:14-20` and `.github/pull_request_template.md:28` document release evidence expectations and the prohibition on committed `release/manifest/latest.json`.
- `.githooks/pre-commit:4-15` blocks main/master commits and scans staged secrets.
- `.githooks/pre-push:4-8` blocks direct main/master pushes.
- `.devcontainer/devcontainer.json:3` points to compose configuration, and `.devcontainer/devcontainer.json:8` sets `XLIB_CONTEXT=docker_toolchain`.
- `.gitignore:27-56`, `.gitignore:64-68`, `.gitignore:86-90`, and `.gitignore:107-110` classify generated evidence/runtime artifacts, observed GitHub rules, `.omx/.worktree`, caches, and duplicate generated ignores.
- `.dockerignore:5` excludes `.agent/inbox`; `.dockerignore:11` excludes `release/downstream-sync/latest.md`.
- `docs/v0.6.0/strict-config-root-omission-audit.md:20-42` lists platform-forced entries such as `.github`, Dependabot, CODEOWNERS, `.dockerignore`, Dockerfile, `.devcontainer`, Go module files, README, and LICENSE.
- `docs/v0.6.0/strict-config-root-omission-audit.md:90-95` says platform adapters are external platform APIs and may be allowed only when pathguard-limited and not treated as xlib standard fact sources.

Risk:

- Platform-native files cannot be blindly migrated into `.config`, but they must be inventoried and guarded as adapters.
- Generated artifacts are intentionally ignored and should not be committed as a side effect of migration proof.

Closure item:

- Add a platform-adapter inventory and generated-artifact policy check to the stable gate so GitHub, hook, devcontainer, ignore, and generated outputs are explicitly classified.

## D. Downstream and template reference omissions

Current facts:

- `scripts/render_template.sh:7-12` defines template rendering and governance-lock behavior.
- `scripts/render_template.sh:205-207` and `scripts/render_template.sh:237-275` cover template package movement, legacy/name replacement, and governance pack emission.
- `scripts/check_rendered_template.sh:45-107` checks rendered Docker, devcontainer, and governance requirements.
- `scripts/check_rendered_template.sh:216-233` checks rendered stale legacy/module/package references.
- `scripts/render_template_test.go:11-44` asserts rendered templates exclude `.omc`, `.omx`, `.worktree`, release manifests, downstream sync, and debt artifacts.
- `scripts/render_template_test.go:66-84` checks control-plane paths.
- `scripts/render_template_test.go:124-150` checks governance pack rendering.
- `scripts/render_template_test.go:320-369` checks archive/runtime omissions and untracked file handling.
- `templates/l2/.agent/l2-capabilities.yaml:21-28` still references `.agent/evidence/l2` template evidence paths.
- `scripts/verify_l2_standard.py:29-33` requires template files under `templates/l2/.agent/...`.
- `docs/downstream-matrix.md:23-28` defines downstream evidence and release-proof expectations.
- `docs/standard/downstream-registry.md:3-26` says current validation is patch-only, proof-oriented, and not downstream adoption truth.
- `.agent/registries/downstream-registry.yaml:1-31` defines the local-only downstream registry, no downstream writes, required Docker contract, and gaps.
- `.agent/registries/downstream-adoption-status.yaml:1-40` and `.agent/registries/downstream-adoption-status.yaml:121-133` classify downstreams as registered-not-adopted and forbid false upgrade claims.
- `cmd/goalcli/adoption_check.go:39-63` requires `.agent/harness/harness.yaml`, `.githooks`, `.github`, `mk/governance.mk`, command registries, and Makefile adoption targets.
- `cmd/goalcli/adoption_check.go:177-185` keeps the adoption contract Makefile-driven.

Risk:

- Templates and downstream checks can continue to emit or require `.agent` paths after a stable config gate if migration is only checked in current-repo docs.
- Downstream evidence proves local patch/proof behavior; it does not prove external repository adoption.

Closure item:

- Add template/downstream migration coverage to the stable gate and maintain an allowed legacy-context list for migration docs and compatibility adapters.

## Review risks and stop conditions

- `.worktree/stable.md` is absent, so this task cannot update it inline.
- `.config/` is absent, so this task cannot truthfully claim config-root completion.
- Existing docs explicitly block release-ready until `.config` migration and classification gates exist.
- Old names are not all defects; migration context must be classified.
- Generated artifacts remain ignored and should not be committed.

## Verification commands for this slice

Documentation/index validation:

```sh
GOWORK=off make docs-check
GOWORK=off make evidence-check
```

Targeted contract/template/downstream validation:

```sh
GOWORK=off go test ./contracts -run 'TestExecutionEvidenceContractMatchesEvidenceManifest' -count=1
GOWORK=off go test ./scripts -run 'TestRenderTemplateIncludesGovernancePack|TestRenderTemplateIncludesDockerContract|TestRenderTemplatePrunesOmittedAgentInboxIndexEntries|TestRenderTemplateGitArchivePrunesRuntimeState|TestRenderTemplateGitArchiveSkipsUntrackedFiles' -count=1
GOWORK=off go test ./cmd/goalcli -run 'TestRunDownstreamSyncPlan|TestRunDispatchesDownstreamSyncPlan|TestRunAdoptionCheck' -count=1
```

Full local verification when promoting this review to implementation:

```sh
GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off make lint
GOWORK=off make p2-runtime-check
GOWORK=off make release-check
```

Static migration classification scan from `docs/migration/baselib-template-to-xlib-standard.md:32-36`:

```sh
rg -n "baselib-template|foundationx|\.agent|\.xlib|\.config" \
  README.md docs .agent .xlib .github .githooks .devcontainer templates contracts examples \
  cmd internal pkg scripts testkit Makefile Dockerfile docker-compose.yml .gitignore .dockerignore \
  --glob '*.go' --glob '*.sh' --glob '*.md' --glob '*.yaml' --glob '*.yml' --glob '*.json' --glob '*.tmpl'
```

## Delegation evidence

- Repository map probe `019ea01d-9644-7af1-a92c-34522d91c238` confirmed `.worktree/stable.md` absence and supplied A-D repository anchors.
- Review/risk probe `019ea01d-a939-76a1-8064-6c5bd73c1a70` confirmed `.config/` absence, `.agent/.xlib` source-of-truth risks, and current validation gaps.
- Test/verification probe `019ea01d-c450-7a83-a9af-69b6eea611c2` confirmed existing verification targets and independently ran `GOWORK=off make vet`, `make lint`, `make test`, `make docs-check`, and `make p2-runtime-check` successfully before this artifact was added.
