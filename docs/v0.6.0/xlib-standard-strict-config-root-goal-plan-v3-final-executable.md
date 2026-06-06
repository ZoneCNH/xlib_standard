# xlib-standard strict-config-root v3：最彻底重构最终可执行 GOAL 方案

版本：v3 final executable
目标：禁止向后兼容，统一在 `.config/`，一次性从根源解决 xlib-standard 标准事实源分散、旧路径回潮、下游采纳不可验证、release evidence 无法证明 strict clean 的问题。

---

## 0. 总结判断

v2 方案方向正确，但必须补齐 6 个 P0 漏洞，否则 strict `.config/` 会在平台入口、下游校验、模板目录、release evidence、fingerprint、依赖自动化中重新长出第二套事实源。

最终 v3 方案为：

```text
strict .config root
+ platform adapter taxonomy
+ Renovate-only dependency automation
+ downstream standard-snapshot
+ root templates/mk 全部禁用
+ release manifest v2 strict fields
+ canonical JSON fingerprint
+ supply-chain/toolchain/security/dockerignore/governance rules
+ negative tests
+ flag-day cutover
```

核心原则：

```text
.config/ 是唯一机器事实源。
平台文件只能是 adapter 或 native interface。
旧路径不是兼容项，而是违规项。
下游不复制完整标准源，但必须携带可离线验证的 standard-snapshot。
所有 sha256 只基于 canonical JSON。
release evidence 必须证明 strict clean。
```

---

## 1. 最终目录

### 1.1 xlib-standard 标准源仓库

```text
.config/
├── xlib/
│   ├── xlib.json
│   ├── profiles/
│   ├── capabilities/
│   ├── downstream/
│   ├── rules/
│   ├── schemas/
│   ├── contracts/
│   ├── harness/
│   ├── agent/
│   ├── templates/
│   └── mk/
├── git/hooks/
├── golangci/golangci.yml
└── renovate/renovate.json

.github/workflows/        # GitHub 平台 adapter，只调用 .config 驱动命令
Dockerfile                # Docker 平台 native interface
.dockerignore             # Docker 平台 native interface，由 .config 规则校验
go.mod/go.sum             # Go 平台 native interface
Makefile                  # adapter，只 include .config/xlib/mk/*.mk
cmd/internal/pkg/docs/... # 代码、展示、工具入口
release/                  # Evidence 产物，不是事实源
```

### 1.2 下游仓库

下游不复制完整标准控制面，但必须有标准快照：

```text
.config/xlib/
├── xlib-standard.lock.json
├── standard-snapshot.json
├── adoption-proof.json
├── boundary-report.json
├── contract-fingerprint.json
└── profile-plan.json
```

下游禁止：

```text
xlib-standard.lock
.agent/
.xlib/
contracts/
.githooks/
templates/
mk/
```

---

## 2. Platform Adapter Taxonomy

新增：

```text
.config/xlib/rules/platform-adapters.json
```

目的：避免把平台必须存在的入口误判为兼容层，同时禁止它们承载 xlib 标准事实。

```json
{
  "schema_version": "2.0",
  "categories": {
    "thin_adapter": [
      ".github/workflows/ci.yml",
      ".github/workflows/release.yml",
      ".github/workflows/security.yml",
      "Makefile"
    ],
    "platform_native": [
      "go.mod",
      "go.sum",
      "Dockerfile",
      ".dockerignore",
      ".devcontainer/devcontainer.json",
      ".github/CODEOWNERS",
      ".github/rulesets/protect-main.json",
      ".github/pull_request_template.md",
      ".github/ISSUE_TEMPLATE"
    ],
    "forbidden_legacy": [
      ".agent",
      ".xlib",
      "contracts",
      "xlib-standard.lock",
      ".githooks",
      "templates",
      "mk",
      "renovate.json",
      ".golangci.yml"
    ]
  },
  "thin_adapter_policy": {
    "may_call": [
      "make",
      "go run ./cmd/goalcli"
    ],
    "must_pass_config": "--config .config/xlib/xlib.json",
    "must_not_define": [
      "profiles",
      "downstream_targets",
      "score_threshold",
      "contract_paths",
      "harness_gates",
      "release_policy"
    ]
  },
  "platform_native_policy": {
    "must_be_checked_by_pathguard": true,
    "must_not_override_xlib_facts": true
  }
}
```

Makefile 和 GitHub workflow 不是事实源，只是 adapter。`.dockerignore`、`go.mod`、`go.sum` 是平台 native interface，但它们的 policy 必须由 `.config/xlib/rules/*.json` 校验。

---

## 3. Dependency Automation：禁用 Dependabot，Renovate-only

### 3.1 决策

最彻底方案选择：

```text
禁止 .github/dependabot.yml
禁止 root renovate.json
只允许 .config/renovate/renovate.json
GitHub workflow 只作为 adapter 调用 Renovate CLI
```

### 3.2 新规则

```text
.config/xlib/rules/dependency-automation.json
```

```json
{
  "schema_version": "2.0",
  "mode": "renovate_only",
  "forbidden_paths": [
    ".github/dependabot.yml",
    "renovate.json"
  ],
  "required_paths": [
    ".config/renovate/renovate.json"
  ],
  "github_workflow_policy": {
    "may_run_renovate": true,
    "required_config_file": ".config/renovate/renovate.json",
    "forbid_inline_renovate_config": true
  }
}
```

### 3.3 dependency-check 新规则

`goalcli dependency-check --config .config/xlib/xlib.json` 检查：

```text
1. .github/dependabot.yml 不存在
2. root renovate.json 不存在
3. .config/renovate/renovate.json 存在
4. workflow 不内联依赖策略
5. third-party Actions 40-char SHA pinned
6. 禁止 @latest
7. govulncheck 版本固定
8. go.mod/go.sum 与 toolchain rules 一致
```

---

## 4. Downstream standard-snapshot

### 4.1 为什么 lock-only 不够

如果下游只有：

```text
.config/xlib/xlib-standard.lock.json
```

就无法离线验证：

```text
forbidden paths
profile boundary
required gates
schema fingerprints
snapshot 是否与 lock 匹配
```

所以必须新增：

```text
.config/xlib/standard-snapshot.json
```

### 4.2 standard-snapshot.json

```json
{
  "schema_version": "2.0",
  "standard": {
    "name": "xlib-standard",
    "version": "v0.6.0",
    "commit": "<commit>",
    "strict_config_only": true
  },
  "profile": {
    "id": "l2.redis",
    "layer": "L2",
    "required_capabilities": [],
    "required_gates": [],
    "allowed_imports": [],
    "forbidden_imports": [],
    "forbidden_concepts": []
  },
  "rules": {
    "pathguard": {},
    "boundary": {},
    "security": {},
    "supply_chain": {},
    "toolchain": {},
    "fingerprint": {}
  },
  "contracts": {
    "config_schema_sha256": "<sha256>",
    "error_schema_sha256": "<sha256>",
    "health_schema_sha256": "<sha256>",
    "metrics_schema_sha256": "<sha256>",
    "release_manifest_schema_sha256": "<sha256>"
  },
  "fingerprint": {
    "canonical_json_sha256": "<sha256>"
  }
}
```

### 4.3 下游命令改造

下游不使用 `--config .config/xlib/xlib.json`，而使用 lock/snapshot：

```bash
goalcli pathguard --lock .config/xlib/xlib-standard.lock.json
goalcli boundary --lock .config/xlib/xlib-standard.lock.json
goalcli contracts --lock .config/xlib/xlib-standard.lock.json
goalcli adoption-check --lock .config/xlib/xlib-standard.lock.json --verify
goalcli strict-check --lock .config/xlib/xlib-standard.lock.json
```

`lock` 必须引用 snapshot：

```json
{
  "snapshot": {
    "path": ".config/xlib/standard-snapshot.json",
    "sha256": "<sha256>"
  }
}
```

---

## 5. 禁止 root templates/ 和 root mk/

如果目标是所有 xlib-controlled 机器事实源统一在 `.config/`，那么根目录 `templates/` 和 `mk/` 都不能保留。

迁移：

```text
templates/ -> .config/xlib/templates/
mk/        -> .config/xlib/mk/
```

`pathguard.json` 中直接禁止：

```json
{
  "forbidden_paths": [
    ".agent",
    ".xlib",
    "contracts",
    "xlib-standard.lock",
    ".githooks",
    "templates",
    "mk",
    "renovate.json",
    ".golangci.yml"
  ]
}
```

Makefile 只能：

```makefile
include .config/xlib/mk/governance.mk
include .config/xlib/mk/harness.mk
include .config/xlib/mk/release.mk
```

---

## 6. Release Manifest v2 strict fields

`release/manifest/latest.json` 仍然是生成 Evidence，不进入 `.config/`。但它必须证明 strict clean。

新增字段：

```json
{
  "strict_config_root": {
    "enabled": true,
    "config_root": ".config/xlib",
    "legacy_compatibility": false,
    "pathguard_status": "passed",
    "strict_check_status": "passed",
    "legacy_path_count": 0,
    "legacy_flag_count": 0,
    "root_lock_present": false,
    "contracts_root_present": false,
    "agent_root_present": false
  },
  "config_fingerprints": {
    "xlib_json_sha256": "<sha256>",
    "profiles_sha256": "<sha256>",
    "capabilities_sha256": "<sha256>",
    "downstream_sha256": "<sha256>",
    "rules_sha256": "<sha256>",
    "schemas_sha256": "<sha256>",
    "contracts_sha256": "<sha256>",
    "harness_sha256": "<sha256>",
    "agent_sha256": "<sha256>",
    "templates_sha256": "<sha256>"
  },
  "platform_adapters": {
    "github_workflows": "passed",
    "dockerignore": "passed",
    "dockerfile": "passed",
    "go_mod": "passed",
    "dependency_automation": "passed"
  },
  "downstream_strict_adoption": {
    "lock_path": ".config/xlib/xlib-standard.lock.json",
    "snapshot_required": true,
    "snapshot_sha256": "<sha256>"
  }
}
```

Release gate 不允许只说 `release-check passed`，必须证明 strict-config-root clean。

---

## 7. Canonical JSON fingerprint

新增：

```text
.config/xlib/rules/fingerprint.json
```

```json
{
  "schema_version": "2.0",
  "algorithm": "sha256",
  "encoding": "canonical-json-v1",
  "rules": {
    "object_keys_sorted": true,
    "line_endings": "LF",
    "trailing_newline": true,
    "unicode_normalization": "NFC",
    "exclude_fields": [
      "generated_at",
      "last_checked_at",
      "duration_ms",
      "hostname",
      "pid",
      "tmpdir"
    ]
  }
}
```

新增命令：

```bash
goalcli fingerprint --config .config/xlib/xlib.json --path .config/xlib/profiles/l2.redis.json
goalcli fingerprint tree --config .config/xlib/xlib.json --path .config/xlib/rules
```

所有 lock、snapshot、release manifest 指纹必须使用 canonical JSON。禁止直接对 map marshal 后 hash。

---

## 8. Supply Chain / Toolchain / Docker / Security Rules

### 8.1 supply-chain.json

```json
{
  "schema_version": "2.0",
  "github_actions": {
    "third_party_actions_must_be_40_char_sha": true,
    "allow_official_actions_by_sha_only": true,
    "forbid_latest": true
  },
  "go_tools": {
    "govulncheck_version": "golang.org/x/vuln/cmd/govulncheck@v1.1.4",
    "forbid_latest": true
  },
  "dependencies": {
    "renovate_config": ".config/renovate/renovate.json",
    "dependabot_forbidden": true
  }
}
```

### 8.2 toolchain.json

```json
{
  "schema_version": "2.0",
  "go": {
    "go_mod_required": true,
    "go_sum_required": true,
    "go_work_forbidden_in_repo": true,
    "release_gates_require_gowork_off": true
  },
  "lint": {
    "config": ".config/golangci/golangci.yml"
  }
}
```

### 8.3 docker.json

```json
{
  "schema_version": "2.0",
  "dockerignore": {
    "path": ".dockerignore",
    "must_exclude": [
      ".git",
      ".cache",
      "release/manifest/latest.json",
      "release/manifest/latest.json.sha256",
      ".env",
      "*.pem",
      "*.key"
    ],
    "must_not_exclude": [
      ".config/xlib"
    ]
  },
  "build_context": {
    "config_required": true,
    "release_artifacts_not_source": true
  }
}
```

### 8.4 security.json

```json
{
  "schema_version": "2.0",
  "secrets": {
    "forbidden_paths": [
      "/home/k8s/secrets/env/*"
    ],
    "forbidden_outputs": [
      "README.md",
      "docs/",
      "release/manifest/latest.json",
      "PR description",
      "test logs"
    ]
  },
  "vulnerability_scan": {
    "weekly_window_required": true,
    "force_flag": "XLIB_FORCE_VULNCHECK",
    "enable_flag": "XLIB_ENABLE_VULNCHECK"
  }
}
```

---

## 9. Negative Test Matrix

必须新增负向测试，证明旧系统无法回潮。

```text
TestLegacyAgentPathFails
TestLegacyContractsPathFails
TestLegacyRootLockFails
TestLegacyGitHooksFails
TestLegacyTemplatesRootFails
TestLegacyMkRootFails
TestRenderEnableGovernanceFlagFails
TestRenderLayerFlagFails
TestWorkflowInlineScoreThresholdFails
TestWorkflowActionLatestFails
TestWorkflowActionTagFails
TestDependabotConfigFails
TestRootRenovateConfigFails
TestDownstreamWithoutSnapshotFails
TestDownstreamRootLockFails
TestNonCanonicalFingerprintFails
TestReleaseManifestMissingStrictConfigRootFails
TestDocsCurrentUsageLegacyPathFails
TestDockerignoreExcludesConfigFails
TestGoWorkCommittedFails
```

每个负向测试都必须有 fixture：

```text
testdata/strict/legacy-agent/
testdata/strict/root-lock/
testdata/strict/dependabot/
testdata/strict/workflow-latest/
testdata/strict/downstream-no-snapshot/
```

---

## 10. Flag-day Cutover Protocol

禁止兼容意味着不能逐步读旧路径。必须 flag-day：

```text
T-2 days: freeze xlib-standard structural changes
T-1 day : generate downstream migration PR plan
T day   : merge strict-config-root-v3 breaking PR
T+0     : regenerate kernel/configx/redisx/domain-market
T+1     : regenerate all foundation downstreams
T+2     : migrate business governance-only repos
T+3     : enable drift-check org-wide
```

Rollback 策略：

```text
只能 full revert strict-config-root-v3 PR。
不允许引入 compatibility mode。
不允许恢复 root lock。
不允许恢复 .agent/ contracts/。
```

---

## 11. Final GOAL Stack

```text
G0: xlib-standard strict-config-root v3 根治

  G1: 唯一事实源
    G1.1: 所有 xlib config in .config/xlib
    G1.2: platform-adapters taxonomy
    G1.3: pathguard P0

  G2: 旧系统不可回潮
    G2.1: forbidden legacy paths
    G2.2: forbidden legacy flags
    G2.3: negative tests

  G3: 生成器 profile-driven
    G3.1: --config required
    G3.2: --profile required
    G3.3: layer/governance derived from profile
    G3.4: lock only in .config/xlib

  G4: 下游可离线验证
    G4.1: lock
    G4.2: standard-snapshot
    G4.3: adoption proof
    G4.4: boundary/contract reports

  G5: release evidence 可证明 strict clean
    G5.1: release manifest v2 strict fields
    G5.2: canonical fingerprints
    G5.3: platform adapter status

  G6: supply chain / toolchain / docker / security 收口
    G6.1: Renovate-only
    G6.2: Action SHA pinning
    G6.3: govulncheck fixed version
    G6.4: go.work forbidden
    G6.5: dockerignore checked
```

---

## 12. Final Acceptance Commands

标准源仓库：

```bash
test ! -e .agent
test ! -e .xlib
test ! -e contracts
test ! -e xlib-standard.lock
test ! -e .githooks
test ! -e templates
test ! -e mk
test ! -e renovate.json
test ! -e .golangci.yml
test ! -e .github/dependabot.yml

GOWORK=off make pathguard
GOWORK=off make strict-check
GOWORK=off make config-check
GOWORK=off make profile-check
GOWORK=off make downstream-check
GOWORK=off make dependency-check
GOWORK=off make contracts
GOWORK=off make boundary
GOWORK=off make docs-check
GOWORK=off make ci
GOWORK=off make integration
GOWORK=off make release-check
```

旧命令必须失败：

```bash
! scripts/render_template.sh --enable-governance
! scripts/render_template.sh --layer L2
```

渲染：

```bash
tmp="$(mktemp -d)"

scripts/render_template.sh \
  --config .config/xlib/xlib.json \
  --profile l2.redis \
  --module-name redisx \
  --module-path github.com/ZoneCNH/redisx \
  --package-name redisx \
  --standard-version v0.6.0 \
  --standard-commit "$(git rev-parse HEAD)" \
  --out "$tmp/redisx"

test -f "$tmp/redisx/.config/xlib/xlib-standard.lock.json"
test -f "$tmp/redisx/.config/xlib/standard-snapshot.json"
test ! -e "$tmp/redisx/xlib-standard.lock"
test ! -e "$tmp/redisx/.agent"
test ! -e "$tmp/redisx/contracts"

(
  cd "$tmp/redisx"
  GOWORK=off make pathguard
  GOWORK=off make adoption-check
  GOWORK=off make release-check
)
```

---

## 13. 最终推荐实施顺序

一个 breaking PR，内部 commit 顺序：

```text
1. add .config/xlib v3 control plane
2. add platform-adapters / fingerprint / supply-chain / docker / toolchain rules
3. add standard-snapshot model
4. add configroot/profile/downstream/pathguard/fingerprint/lockfile packages
5. add strict-check and negative test fixtures
6. move schemas/contracts/harness/agent/templates/mk into .config
7. remove legacy paths
8. rewrite render_template.sh profile-driven only
9. rewrite Makefile and GitHub workflow adapters
10. rewrite README/docs
11. update integration golden targets
12. release manifest v2 strict fields
13. full release gate
```

---

## 14. 最终结论

最彻底方案不是“把旧文件移动进 `.config`”。

最彻底方案是：

```text
1. 定义平台 adapter 与标准事实源的边界；
2. 删除所有 legacy source；
3. 让 pathguard 成为 P0 gate；
4. 让下游通过 lock + standard-snapshot 离线验证；
5. 让 release manifest 证明 strict clean；
6. 用 canonical JSON 让 fingerprint 稳定；
7. 用负向测试证明旧系统不能回潮。
```

一句话：

> `.config/` 不是目录选择，而是权力收敛。旧路径不是历史遗留，而是架构违规。平台入口不是兼容层，而是受 `.config` 管控的 adapter。
