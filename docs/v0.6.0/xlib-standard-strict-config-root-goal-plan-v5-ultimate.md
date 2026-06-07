# xlib-standard strict-config-root v5：最终闭环执行版

版本：v5 ultimate closure
目标：在 v3/v4 的基础上，继续消除“伪严格”风险：平台原生文件、远程仓库策略、CODEOWNERS、GitHub rulesets、Docker context、下游工具链、full-copy render、projection 校验、环境变量逃逸、子模块和 symlink 绕过。

---

> **Current-status note (2026-06-07):** This strict `.config/` plan is historical planning material, not current authority to delete or migrate the live `.agent/**` runtime/adoption control plane. Current stable hardening keeps `.agent/**` protected under `CONSTITUTION.md`, `.agent/index.yaml`, `goalcli adoption-check`, Codex agent prompts, L2 templates, and release/adoption gates until an approved dual-read migration and compatibility test suite lands.

## 0. 最终结论

最彻底方案不是简单把目录移动到 `.config/`，而是建立三套互斥的事实域：

```text
A. xlib standard facts
   只能存在于 .config/xlib/。

B. platform-native executable facts
   平台强制读取的入口，例如 .github/workflows、.github/CODEOWNERS、.dockerignore、go.mod。
   它们不是向后兼容目录，而是平台 API。
   必须由 .config 声明、生成或审计，不能承载 xlib 标准事实。

C. forbidden legacy facts
   .agent、.xlib、contracts、root xlib-standard.lock、root templates、root mk、.githooks 等。
   出现即失败。
```

最终原则：

```text
.config/xlib 是唯一 xlib 标准事实源。
平台原生文件允许存在，但必须被 .config 监管。
legacy mirror / fallback / dual read 全部禁止。
下游只携带 effective snapshot，不复制完整标准源。
远程 GitHub 仓库设置必须进入 release evidence。
```

---

## 1. v5 新增的 P0 闭环点

v4 已经补齐 repo_kind、render prune、standard-snapshot、symlink/path escape、go:embed 和 pinned goalcli。但仍遗漏一个关键层：**平台远程状态和平台原生投影**。

新增 P0：

```text
P0-1  platform projection taxonomy：区分 platform-native、thin-adapter、generated-projection、forbidden-legacy。
P0-2  remote GitHub policy audit：branch protection / rulesets / Actions permissions / CODEOWNERS enforcement 不能只靠仓库文件证明。
P0-3  CODEOWNERS 不能移入 .config 后消失；GitHub 只识别 .github、root 或 docs 中的 CODEOWNERS。
P0-4  .github/rulesets/*.json 是 declarative artifact，不等于远程 ruleset 已启用；release 必须审计远程状态。
P0-5  render 必须使用 allowlist copy，而不是 copy-all 后 prune。
P0-6  platform projections 必须有 generated checksum，但 legacy mirrors 仍然禁止。
P0-7  .gitmodules、submodule、worktree、case-only path collision 必须纳入 pathguard。
P0-8  .config 也必须被 secret scan，不能因为是配置事实源就豁免。
```

---

## 2. 最终分类模型

新增：

```text
.config/xlib/rules/platform-projections.json
```

示例：

```json
{
  "schema_version": "2.0",
  "classes": {
    "xlib_standard_fact": {
      "root": ".config/xlib",
      "exclusive": true
    },
    "platform_native": {
      "description": "Files that external tools require at fixed paths.",
      "must_be_declared_in_config": true,
      "must_be_audited": true
    },
    "thin_adapter": {
      "description": "Entrypoint that only calls config-driven commands.",
      "must_not_define_standard_facts": true
    },
    "generated_projection": {
      "description": "Platform-required projection generated from .config.",
      "must_have_generated_marker": true,
      "must_match_checksum": true
    },
    "forbidden_legacy": {
      "description": "Old xlib-controlled paths. Existence fails pathguard.",
      "existence_allowed": false
    }
  },
  "paths": [
    {
      "path": ".config/xlib/**",
      "class": "xlib_standard_fact"
    },
    {
      "path": ".github/workflows/*.yml",
      "class": "thin_adapter",
      "required_calls": [
        "GOWORK=off make ci",
        "GOWORK=off make release-check"
      ]
    },
    {
      "path": ".github/CODEOWNERS",
      "class": "generated_projection",
      "source": ".config/github/codeowners.json"
    },
    {
      "path": ".github/rulesets/*.json",
      "class": "generated_projection",
      "source": ".config/github/rulesets/*.json",
      "remote_audit_required": true
    },
    {
      "path": "go.mod",
      "class": "platform_native",
      "reason": "Go toolchain module manifest"
    },
    {
      "path": "go.sum",
      "class": "platform_native",
      "reason": "Go checksum database artifact"
    },
    {
      "path": "Dockerfile",
      "class": "platform_native",
      "must_not_define_xlib_gates": true
    },
    {
      "path": ".dockerignore",
      "class": "platform_native",
      "must_match": ".config/docker/dockerignore.policy.json"
    },
    {
      "path": ".gitignore",
      "class": "platform_native",
      "must_match": ".config/git/gitignore.policy.json"
    },
    {
      "path": ".gitattributes",
      "class": "platform_native",
      "must_match": ".config/git/gitattributes.policy.json"
    },
    {
      "path": ".agent",
      "class": "forbidden_legacy"
    },
    {
      "path": "contracts",
      "class": "forbidden_legacy"
    },
    {
      "path": "xlib-standard.lock",
      "class": "forbidden_legacy"
    },
    {
      "path": "templates",
      "class": "forbidden_legacy"
    },
    {
      "path": "mk",
      "class": "forbidden_legacy"
    }
  ]
}
```

关键定义：

```text
legacy mirror = 禁止。
platform projection = 允许，但必须由 .config 生成/校验，且只为外部平台读取服务。
```

这避免把 `.github/CODEOWNERS`、`.github/workflows`、`.dockerignore` 误判成“向后兼容”。它们不是旧 xlib 事实源，而是外部平台的固定接口。

---

## 3. GitHub 远程策略审计

新增：

```text
.config/github/repository-policy.json
.config/github/codeowners.json
.config/github/rulesets/*.json
.config/xlib/rules/remote-policy.json
```

### 3.1 repository-policy.json

```json
{
  "schema_version": "2.0",
  "provider": "github",
  "repository": "ZoneCNH/xlib-standard",
  "default_branch": "main",
  "required": {
    "branch_protection_or_ruleset": true,
    "require_pull_request": true,
    "require_code_owner_review": true,
    "require_status_checks": true,
    "block_force_pushes": true,
    "block_deletions": true,
    "actions_permissions_restricted": true,
    "workflow_actions_pinned": true,
    "dependabot_config_forbidden": true,
    "renovate_config_required": true
  },
  "required_status_checks": [
    "pathguard",
    "strict-check",
    "config-check",
    "profile-check",
    "contracts",
    "boundary",
    "docs-check",
    "release-check"
  ]
}
```

### 3.2 新命令

```bash
goalcli platform github-audit --config .config/xlib/xlib.json
```

检查：

```text
1. 远程 default branch 是否为 main。
2. ruleset/branch protection 是否强制 PR。
3. 是否强制 CODEOWNERS review。
4. required status checks 是否包含 strict gate。
5. Actions 权限是否受限。
6. 是否禁止 force push / delete。
7. .github/rulesets/*.json 是否与远程 ruleset 等价。
8. .github/dependabot.yml 不存在。
9. .config/renovate/renovate.json 存在且 workflow 只调用 Renovate。
```

本地无 token 时：

```text
local ci: github-audit 可输出 unverified，但不能标记 passed。
release-check: github-audit unverified 必须 fail。
```

Release manifest 必须记录：

```json
{
  "remote_policy": {
    "provider": "github",
    "status": "passed",
    "default_branch": "main",
    "ruleset_audited": true,
    "codeowners_required": true,
    "required_status_checks": ["pathguard", "strict-check", "release-check"]
  }
}
```

---

## 4. CODEOWNERS 的处理

GitHub 只识别 `.github/CODEOWNERS`、root `CODEOWNERS` 或 `docs/CODEOWNERS`。因此最彻底方案不能把 CODEOWNERS 完全藏进 `.config` 后删除平台文件。

最终设计：

```text
.config/github/codeowners.json       canonical intent
.github/CODEOWNERS                   generated platform projection
```

这不是 legacy mirror，而是 platform projection。

### 4.1 codeowners.json

```json
{
  "schema_version": "2.0",
  "default_owners": ["@ZoneCNH/release-owners"],
  "rules": [
    {
      "pattern": "*",
      "owners": ["@ZoneCNH/release-owners"]
    },
    {
      "pattern": "/.config/xlib/**",
      "owners": ["@ZoneCNH/architecture-owners", "@ZoneCNH/release-owners"]
    },
    {
      "pattern": "/cmd/**",
      "owners": ["@ZoneCNH/tooling-owners"]
    },
    {
      "pattern": "/internal/**",
      "owners": ["@ZoneCNH/tooling-owners"]
    },
    {
      "pattern": "/.github/**",
      "owners": ["@ZoneCNH/release-owners"]
    }
  ]
}
```

### 4.2 生成投影

```bash
goalcli platform codeowners render --config .config/xlib/xlib.json --out .github/CODEOWNERS
goalcli platform codeowners check --config .config/xlib/xlib.json
```

`pathguard` 允许 `.github/CODEOWNERS`，但要求：

```text
1. 有 generated marker。
2. sha256 与 .config/github/codeowners.json 匹配。
3. 文件大小低于 GitHub 限制。
4. 语法可解析。
5. .github/CODEOWNERS 自身归 release owners。
```

---

## 5. Render：从 copy-all 改为 allowlist materialization

v4 说 render 后 prune 完整 `.config/xlib`。v5 进一步收紧：**不要 copy-all 再 prune，必须 allowlist materialize**。

禁止流程：

```text
copy repo -> remove forbidden files -> hope clean
```

允许流程：

```text
resolve profile -> resolve fragment plan -> materialize allowlisted files -> write target .config/xlib artifacts -> strict-render-check
```

新增：

```text
.config/xlib/templates/materialization-policy.json
```

```json
{
  "schema_version": "2.0",
  "mode": "allowlist_only",
  "standard_source_never_copied": true,
  "forbidden_source_roots": [
    ".git",
    ".config/xlib/profiles",
    ".config/xlib/downstream",
    ".config/xlib/agent",
    ".config/xlib/harness",
    ".config/xlib/rules",
    ".config/xlib/schemas",
    "release",
    ".cache"
  ],
  "target_required_artifacts": [
    ".config/xlib/xlib-target.json",
    ".config/xlib/xlib-standard.lock.json",
    ".config/xlib/standard-snapshot.json",
    ".config/xlib/adoption-proof.json",
    ".config/xlib/boundary-report.json",
    ".config/xlib/contract-fingerprint.json",
    ".config/xlib/profile-plan.json"
  ]
}
```

`render_template.sh` 应该退化为 thin shell，真正渲染由：

```bash
goalcli render target \
  --config .config/xlib/xlib.json \
  --profile l2.redis \
  --module-name redisx \
  --module-path github.com/ZoneCNH/redisx \
  --package-name redisx \
  --out ../redisx
```

---

## 6. 下游标准快照：effective subset only

下游不能携带完整标准源。

下游允许：

```text
.config/xlib/xlib-target.json
.config/xlib/xlib-standard.lock.json
.config/xlib/standard-snapshot.json
.config/xlib/adoption-proof.json
.config/xlib/boundary-report.json
.config/xlib/contract-fingerprint.json
.config/xlib/profile-plan.json
```

下游禁止：

```text
.config/xlib/xlib.json
.config/xlib/profiles/
.config/xlib/downstream/
.config/xlib/agent/
.config/xlib/harness/
.config/xlib/templates/
.config/xlib/schemas/
```

### 6.1 standard-snapshot.json

```json
{
  "schema_version": "2.0",
  "snapshot_kind": "effective_target_subset",
  "standard": {
    "name": "xlib-standard",
    "version": "v0.6.0",
    "commit": "<commit>"
  },
  "target": {
    "module_name": "redisx",
    "module_path": "github.com/ZoneCNH/redisx",
    "profile_id": "l2.redis",
    "layer": "L2"
  },
  "effective_profile": {
    "id": "l2.redis",
    "required_capabilities": ["explicit_config", "repository_governance", "release_evidence"],
    "forbidden_imports": ["github.com/ZoneCNH/x.go", "github.com/bytechainx/x.go"],
    "forbidden_concepts": ["business_key_semantics", "hidden_global_client"]
  },
  "effective_gates": [
    {
      "id": "pathguard",
      "command": "xlibctl pathguard --lock .config/xlib/xlib-standard.lock.json"
    },
    {
      "id": "boundary",
      "command": "xlibctl boundary --lock .config/xlib/xlib-standard.lock.json"
    }
  ],
  "effective_schemas": {
    "xlib_lock_schema_sha256": "<sha256>",
    "adoption_proof_schema_sha256": "<sha256>"
  },
  "tooling": {
    "goalcli_distribution": "xlibctl",
    "version": "v0.6.0",
    "commit": "<commit>",
    "sha256": "<sha256>"
  }
}
```

Snapshot 不是 registry 副本，而是 target-specific compiled standard。

---

## 7. 工具链分发闭环

下游不能假设存在：

```text
./cmd/goalcli
```

也不能使用：

```text
go install ...@latest
```

最终选择：

```text
xlibctl = xlib-standard 发布出来的固定 CLI 二进制或固定源码 commit 构建产物。
```

lock 必须记录：

```json
{
  "tooling": {
    "name": "xlibctl",
    "version": "v0.6.0",
    "commit": "<40-char-commit>",
    "source": "github.com/ZoneCNH/xlib-standard/cmd/goalcli",
    "sha256": "<binary-or-source-checksum>",
    "install": {
      "method": "pinned_source_or_release_asset",
      "latest_forbidden": true
    }
  }
}
```

下游 Makefile：

```makefile
XLIB_LOCK ?= .config/xlib/xlib-standard.lock.json
XLIBCTL ?= ./bin/xlibctl

.PHONY: verify-xlibctl
verify-xlibctl:
	$(XLIBCTL) tooling verify --lock $(XLIB_LOCK)

.PHONY: pathguard
pathguard: verify-xlibctl
	$(XLIBCTL) pathguard --lock $(XLIB_LOCK)

.PHONY: boundary
boundary: verify-xlibctl
	$(XLIBCTL) boundary --lock $(XLIB_LOCK)
```

---

## 8. Pathguard v5 防绕过规则

`pathguard` 必须新增：

```text
1. EvalSymlinks：任何 .config/xlib 内 symlink fail。
2. Repo root containment：所有 configured paths 必须在 repo root 内。
3. Casefold scan：Contracts/、CONTRACTS/、.Agent/ 等 case variant fail。
4. Unicode normalization：NFC normalized path 与原始 path 不一致时 fail。
5. .gitmodules：默认禁止；如存在必须 platform-native declared。
6. Submodule boundary：submodule 不得作为 standard source 或 generated target 的一部分。
7. Hardlink / inode alias：在可检测平台上禁止关键配置 hardlink 到 legacy path。
8. go:embed scan：禁止 embed legacy path。
9. Docker COPY/ADD scan：禁止 COPY contracts/.agent/templates/mk。
10. Makefile/shell heredoc scan：禁止 legacy path 藏在多行字符串。
11. Binary fixture policy：二进制 fixture 默认不扫描内容，但必须在 fixture registry 中声明。
12. .config secret scan：.config 不是 secret-scan allowlist。
```

新增：

```text
.config/xlib/rules/fixture-policy.json
```

```json
{
  "schema_version": "2.0",
  "allowed_legacy_mentions": [
    {
      "path": "docs/migration/strict-config-root-v2.md",
      "reason": "migration history"
    },
    {
      "path": "internal/pathguard/testdata/negative/**",
      "reason": "negative tests only",
      "runtime_input_allowed": false
    }
  ]
}
```

---

## 9. Dependency automation 最终决策

彻底方案采用：

```text
Renovate-only。
.github/dependabot.yml forbidden。
.config/renovate/renovate.json canonical。
GitHub workflow 只作为 Renovate adapter。
```

新增：

```text
.config/xlib/rules/dependency-automation.json
```

```json
{
  "schema_version": "2.0",
  "mode": "renovate_only",
  "forbidden": [
    ".github/dependabot.yml",
    "dependabot.yml"
  ],
  "required": [
    ".config/renovate/renovate.json"
  ],
  "workflow_adapter": ".github/workflows/renovate.yml",
  "latest_forbidden": true,
  "pinning_required": true
}
```

`dependency-check` 不再读取 `.github/dependabot.yml`。它检查：

```text
1. .github/dependabot.yml 不存在。
2. .config/renovate/renovate.json 存在。
3. Renovate workflow 是 thin adapter。
4. 所有 actions 固定到 40 位 commit SHA。
5. govulncheck 版本固定。
6. @latest 禁止。
```

---

## 10. Release manifest v5 字段

Release manifest 必须证明四类 clean：

```text
strict_config_clean
platform_projection_clean
remote_policy_clean
generated_target_clean
```

新增字段：

```json
{
  "strict_config_root": {
    "enabled": true,
    "repo_kind": "standard_source",
    "config_root": ".config/xlib",
    "legacy_compatibility": false,
    "pathguard_status": "passed",
    "legacy_path_count": 0,
    "legacy_flag_count": 0,
    "casefold_collision_count": 0,
    "symlink_escape_count": 0
  },
  "platform_projections": {
    "status": "passed",
    "github_workflows": "thin_adapter",
    "codeowners": "generated_projection",
    "rulesets": "generated_projection_with_remote_audit",
    "dockerignore": "platform_native_checked",
    "gitignore": "platform_native_checked"
  },
  "remote_policy": {
    "status": "passed",
    "provider": "github",
    "default_branch": "main",
    "codeowners_required": true,
    "required_status_checks_verified": true,
    "force_push_blocked": true,
    "rulesets_verified": true
  },
  "dependency_automation": {
    "mode": "renovate_only",
    "dependabot_config_absent": true,
    "renovate_config_sha256": "<sha256>"
  },
  "render_materialization": {
    "mode": "allowlist_only",
    "full_copy_forbidden": true,
    "standard_source_pruned": true,
    "generated_targets_verified": ["kernel", "configx", "redisx", "domain-market"]
  }
}
```

---

## 11. v5 新增命令

标准源仓库：

```bash
goalcli platform check --config .config/xlib/xlib.json
goalcli platform codeowners render --config .config/xlib/xlib.json
goalcli platform codeowners check --config .config/xlib/xlib.json
goalcli platform github-audit --config .config/xlib/xlib.json
goalcli render target --config .config/xlib/xlib.json --profile l2.redis ...
goalcli render verify --repo ../redisx
goalcli snapshot write --config .config/xlib/xlib.json --profile l2.redis --out ../redisx/.config/xlib/standard-snapshot.json
goalcli tooling package --config .config/xlib/xlib.json
goalcli tooling verify --lock .config/xlib/xlib-standard.lock.json
```

下游仓库：

```bash
xlibctl tooling verify --lock .config/xlib/xlib-standard.lock.json
xlibctl pathguard --lock .config/xlib/xlib-standard.lock.json
xlibctl boundary --lock .config/xlib/xlib-standard.lock.json
xlibctl contracts --lock .config/xlib/xlib-standard.lock.json
xlibctl adoption-check --lock .config/xlib/xlib-standard.lock.json --verify
```

---

## 12. v5 最终验收矩阵

### 12.1 标准源仓库

```bash
test ! -e .agent
test ! -e .xlib
test ! -e contracts
test ! -e xlib-standard.lock
test ! -e templates
test ! -e mk
test ! -e .githooks
test ! -e .github/dependabot.yml

test -f .config/xlib/xlib.json
test -f .config/xlib/rules/platform-projections.json
test -f .config/xlib/rules/pathguard.json
test -f .config/xlib/rules/dependency-automation.json
test -f .config/github/codeowners.json
test -f .github/CODEOWNERS

GOWORK=off make pathguard
GOWORK=off make strict-check
GOWORK=off make platform-check
GOWORK=off make github-policy-audit
GOWORK=off make config-check
GOWORK=off make profile-check
GOWORK=off make downstream-check
GOWORK=off make dependency-check
GOWORK=off make contracts
GOWORK=off make boundary
GOWORK=off make docs-check
GOWORK=off make integration
GOWORK=off make release-check
```

### 12.2 渲染下游

```bash
tmp="$(mktemp -d)"

goalcli render target \
  --config .config/xlib/xlib.json \
  --profile l2.redis \
  --module-name redisx \
  --module-path github.com/ZoneCNH/redisx \
  --package-name redisx \
  --standard-version v0.6.0 \
  --standard-commit "$(git rev-parse HEAD)" \
  --out "$tmp/redisx"

test -f "$tmp/redisx/.config/xlib/xlib-target.json"
test -f "$tmp/redisx/.config/xlib/xlib-standard.lock.json"
test -f "$tmp/redisx/.config/xlib/standard-snapshot.json"
test ! -f "$tmp/redisx/.config/xlib/xlib.json"
test ! -d "$tmp/redisx/.config/xlib/profiles"
test ! -d "$tmp/redisx/.config/xlib/downstream"
test ! -d "$tmp/redisx/.config/xlib/agent"
test ! -e "$tmp/redisx/.agent"
test ! -e "$tmp/redisx/contracts"
test ! -e "$tmp/redisx/xlib-standard.lock"

(
  cd "$tmp/redisx"
  ./bin/xlibctl tooling verify --lock .config/xlib/xlib-standard.lock.json
  ./bin/xlibctl pathguard --lock .config/xlib/xlib-standard.lock.json
  ./bin/xlibctl boundary --lock .config/xlib/xlib-standard.lock.json
  ./bin/xlibctl adoption-check --lock .config/xlib/xlib-standard.lock.json --verify
)
```

### 12.3 负向测试

必须失败：

```bash
mkdir .agent && ! make pathguard
mkdir contracts && ! make pathguard
touch xlib-standard.lock && ! make pathguard
mkdir templates && ! make pathguard
mkdir mk && ! make pathguard
mkdir .githooks && ! make pathguard
touch .github/dependabot.yml && ! make dependency-check
scripts/render_template.sh --enable-governance && exit 1 || true
scripts/render_template.sh --layer L2 && exit 1 || true
```

必须失败的高级绕过：

```text
1. symlink .config/xlib/profiles -> ../contracts
2. create Contracts/ on macOS-like case-insensitive simulation
3. add //go:embed contracts/*.json
4. Dockerfile COPY contracts /app/contracts
5. workflow uses actions/checkout@v4 instead of 40-char SHA
6. workflow defines SCORE_MIN=9.8 instead of reading .config/xlib/harness/score.json
7. downstream contains full .config/xlib/profiles/
8. downstream xlibctl checksum mismatches lock
9. release-check runs without github remote policy audit in release context
```

---

## 13. Commit order for final breaking PR

仍然是一个 breaking PR，但 commit 顺序调整为：

```text
1. add .config/xlib strict source model and repo_kind
2. add platform projection taxonomy and remote policy schema
3. add configroot/profile/downstream/pathguard/lockfile/snapshot/tooling packages
4. add CODEOWNERS/rulesets/Docker/gitignore platform projection checks
5. add Renovate-only dependency automation and forbid Dependabot
6. add allowlist render materialization and remove copy-all render
7. add xlibctl toolchain packaging and downstream lock verification
8. move schemas/contracts/harness/agent/templates/mk into .config
9. remove legacy paths entirely
10. rewrite Makefile and GitHub workflows as adapter-only
11. rewrite README/docs/generation/standard docs
12. add negative tests for paths, flags, symlinks, go:embed, Dockerfile, workflow, downstream full config
13. update release manifest v5 evidence
14. run golden renders and release gate
```

---

## 14. 最终系统定律

```text
1. .config/xlib 是 xlib 标准事实源。
2. 平台原生文件不是兼容层，但必须由 .config 声明、生成或审计。
3. 旧 xlib 路径永远是 forbidden legacy。
4. 下游只能携带 effective snapshot，不能携带标准源 registry。
5. render 必须 allowlist materialize，不能 copy-all prune。
6. release evidence 必须证明本地路径、平台投影、远程仓库策略、下游生成物全部 strict clean。
7. 任何无法被 pathguard / strict-check / release manifest 证明的合规状态，都视为不存在。
```

---

## 15. 最终推荐

以 v3 作为主实施方案，v4 和 v5 作为必须合并的闭环补丁。真正的最终版应该命名为：

```text
strict-config-root-v2-breaking-release
```

其中“v2”表示标准配置 schema 版本，不表示兼容升级。

最终落地标准：

```text
No legacy path.
No legacy flag.
No full-copy render.
No downstream pseudo-standard-source.
No unpinned tooling.
No unaudited platform native facts.
No remote policy gap.
No release without strict evidence.
```
