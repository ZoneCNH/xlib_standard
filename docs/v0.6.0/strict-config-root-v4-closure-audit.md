# xlib-standard strict-config-root v4：闭环审计补丁

目标：在 v3 final executable 的基础上，继续寻找“严格 `.config/` 根治方案”仍可能被绕过、误伤平台、或让下游验证失真的第三层漏洞。本文件不是替代 v3，而是 v3 的闭环补丁。

---

## 0. 总判断

v3 已经补齐了 P0 大洞：platform adapter taxonomy、Renovate-only、downstream standard-snapshot、root templates/mk 禁止、release manifest v2、canonical JSON fingerprint、supply-chain/toolchain/docker/security rules、negative tests、flag-day cutover。

但继续反向审计后，还要补 13 个闭环点：

```text
P0-1  repo_kind 必须拆分 standard_source 与 generated_target
P0-2  下游不能复制完整 .config/xlib 标准源；render 后必须 prune
P0-3  平台 native 文件不等于兼容层，必须定义 observed-native-fact 规则
P0-4  symlink / path escape / case-insensitive bypass 必须被 pathguard 禁止
P0-5  go:embed / testdata / fixtures / generated docs 必须有 legacy allowlist 模型
P0-6  下游 goalcli 工具分发必须固定，否则 adoption-check 不可重复
P1-1  .gitignore / .gitattributes / .editorconfig 需要 platform-native 分类
P1-2  所有默认路径常量必须集中到 .config/xlib/rules/path-constants.json
P1-3  环境变量覆盖必须受限，不能把 XLIB_CONFIG 指到 legacy path
P1-4  config-check 必须区分 source repo 与 generated repo
P1-5  snapshot 必须声明 effective profile，不是原始 registry 副本
P1-6  migration 文档是唯一允许提旧路径的文档域
P1-7  negative tests 必须覆盖 shell、Makefile、Go code、go:embed、Dockerfile、workflow 五类入口
```

最终闭环原则：

```text
标准源仓库有 .config/xlib/xlib.json。
生成下游仓库没有完整 .config/xlib/xlib.json，只有 lock + standard-snapshot + reports。
平台 native 文件可以存在，但只能作为被 .config 规则验证的执行接口。
所有 path 规则必须抵抗 symlink、大小写、路径穿越、go:embed 和 test fixture 绕过。
```

---

## 1. P0-1：必须加入 repo_kind

### 问题

v3 中 `.config/xlib/xlib.json` 被设计为标准源入口；但下游仓库只应携带 lock + snapshot。如果 `config-check` 默认总是寻找 `.config/xlib/xlib.json`，就会误伤下游；如果下游也复制完整 xlib.json 和 registries，又会让下游变成伪标准源。

### 修正

标准源仓库：

```json
{
  "schema_version": "2.0",
  "repo_kind": "standard_source",
  "standard": {
    "name": "xlib-standard",
    "config_root": ".config/xlib",
    "mode": "strict_config_only",
    "legacy_compatibility": false
  }
}
```

生成下游仓库：

```json
{
  "schema_version": "2.0",
  "repo_kind": "generated_target",
  "lock": ".config/xlib/xlib-standard.lock.json",
  "snapshot": ".config/xlib/standard-snapshot.json"
}
```

推荐文件：

```text
standard source repo:
  .config/xlib/xlib.json

generated target repo:
  .config/xlib/xlib-target.json
  .config/xlib/xlib-standard.lock.json
  .config/xlib/standard-snapshot.json
```

命令分流：

```bash
# 标准源仓库
goalcli strict-check --config .config/xlib/xlib.json

# 下游仓库
goalcli strict-check --target .config/xlib/xlib-target.json
# 或
goalcli strict-check --lock .config/xlib/xlib-standard.lock.json
```

强规则：

```text
standard_source 允许 profiles/rules/schemas/harness/agent/templates/downstream registry。
generated_target 禁止完整 profiles/rules/schemas/harness/agent/templates/downstream registry。
generated_target 必须有 lock + standard-snapshot。
```

---

## 2. P0-2：render 后必须 prune 标准源

### 问题

当前 render 脚本语义是复制整个仓库，然后移动 `pkg/templatex` 并替换标识。strict v3 如果只把标准源移进 `.config/xlib`，full-copy render 会把完整 `.config/xlib` 也复制到下游，导致下游携带标准源副本。

### 修正

render 流程必须改为：

```text
copy repo
  -> remove standard-source-only config
  -> generate target .config/xlib subset
  -> write lock
  -> write standard-snapshot
  -> write reports/proof placeholders
  -> strict-render-check
```

必须删除的下游路径：

```text
.config/xlib/xlib.json
.config/xlib/profiles/
.config/xlib/capabilities/
.config/xlib/downstream/
.config/xlib/rules/
.config/xlib/schemas/
.config/xlib/contracts/
.config/xlib/harness/
.config/xlib/agent/
.config/xlib/templates/
.config/xlib/mk/
```

下游最终只允许：

```text
.config/xlib/xlib-target.json
.config/xlib/xlib-standard.lock.json
.config/xlib/standard-snapshot.json
.config/xlib/adoption-proof.json
.config/xlib/boundary-report.json
.config/xlib/contract-fingerprint.json
.config/xlib/profile-plan.json
```

新增命令：

```bash
goalcli render-prune \
  --repo "$out_dir" \
  --profile l2.redis \
  --lock .config/xlib/xlib-standard.lock.json
```

strict-render-check 必须失败：

```text
if generated_target contains .config/xlib/profiles/ -> fail
if generated_target contains .config/xlib/rules/ -> fail
if generated_target contains .config/xlib/xlib.json -> fail
```

---

## 3. P0-3：平台 native 文件是 observed fact，不是 xlib fact

### 问题

如果说“所有机器事实都必须在 .config”，会与 Go/GitHub/Docker 的平台事实冲突：

```text
go.mod      Go 工具链必须读取 module path / go version
.github/*   GitHub 平台必须读取 workflow / CODEOWNERS / ruleset exports
.dockerignore Docker 必须读取 build context 边界
.gitignore  Git 必须读取忽略规则
```

这些不能迁入 `.config` 后仍被平台自动识别。

### 修正

新增概念：

```text
xlib-owned fact      必须在 .config
platform-native fact 必须留在平台原生位置，但由 .config 规则验证
observed fact        xlib 读取平台事实并验证一致性，但不拥有其源位置
```

新增：

```text
.config/xlib/rules/platform-native-facts.json
```

示例：

```json
{
  "schema_version": "2.0",
  "facts": [
    {
      "path": "go.mod",
      "kind": "go_module",
      "owner": "go_toolchain",
      "xlib_policy": "observe_and_validate",
      "must_not_define": ["profile", "harness_gate", "downstream_target"]
    },
    {
      "path": ".dockerignore",
      "kind": "docker_context_policy",
      "owner": "docker",
      "xlib_policy": "observe_and_validate",
      "config_rule": ".config/xlib/rules/docker.json"
    },
    {
      "path": ".github/workflows/ci.yml",
      "kind": "github_workflow",
      "owner": "github_actions",
      "xlib_policy": "thin_adapter_only"
    }
  ]
}
```

验收：

```text
platform-native 文件可以存在。
它们不能定义 profile、schema path、downstream matrix、score threshold、release policy。
它们必须被 pathguard 和 platform-adapter-check 验证。
```

---

## 4. P0-4：pathguard 必须抵抗 symlink、大小写、路径穿越

### 问题

只做 `os.Stat("contracts")` 不够。绕过方式包括：

```text
.config/xlib/schemas -> ../../contracts    symlink escape
Contracts/                              case-insensitive filesystem
.config/xlib/../contracts               path traversal
cοntracts                               Unicode homoglyph
```

### 修正

pathguard 必须新增：

```text
1. filepath.Clean
2. filepath.EvalSymlinks
3. repoRoot containment check
4. lowercase normalized forbidden path check
5. ASCII-only policy for xlib-controlled paths
6. forbid symlink inside .config/xlib unless explicitly allowed
```

规则：

```json
{
  "schema_version": "2.0",
  "path_safety": {
    "forbid_symlink_under_config_xlib": true,
    "config_paths_must_resolve_inside_repo": true,
    "xlib_controlled_paths_ascii_only": true,
    "case_insensitive_forbidden_match": true,
    "forbid_path_traversal_segments": true
  }
}
```

负向测试：

```text
TestConfigSymlinkEscapesRepoFails
TestConfigSymlinkToContractsFails
TestContractsCaseInsensitiveFails
TestPathTraversalToContractsFails
TestUnicodeHomoglyphContractsFails
```

---

## 5. P0-5：go:embed / testdata / fixtures 需要显式 allowlist

### 问题

旧路径可能藏在：

```text
//go:embed contracts/*.json
//go:embed .agent/*
testdata/legacy-root-lock/
fixtures/contracts/
Dockerfile COPY contracts /
Makefile cat xlib-standard.lock
shell scripts grep .agent
```

如果完全禁止字符串，会误杀负向测试；如果放开 testdata，又会让 legacy fixture 变成实际输入。

### 修正

新增：

```text
.config/xlib/rules/legacy-reference-allowlist.json
```

```json
{
  "schema_version": "2.0",
  "allowed_legacy_mentions": [
    {
      "path": "docs/migration/strict-config-root-v3.md",
      "reason": "migration documentation only"
    },
    {
      "path": "CHANGELOG.md",
      "reason": "release history only"
    },
    {
      "path": "testdata/negative/legacy-paths/**",
      "reason": "negative pathguard fixtures only",
      "must_be_referenced_only_by": ["pathguard_negative_tests"]
    }
  ],
  "forbidden_in_all_other_files": [
    ".agent/",
    "contracts/",
    "xlib-standard.lock",
    "--enable-governance",
    "--layer"
  ]
}
```

pathguard 必须识别：

```text
go:embed directives
Dockerfile COPY/ADD
Makefile shell commands
GitHub workflow run/env blocks
shell scripts
Go string constants
```

---

## 6. P0-6：下游 goalcli 工具分发必须固定

### 问题

下游需要运行：

```bash
goalcli pathguard --lock ...
goalcli boundary --lock ...
goalcli adoption-check --lock ...
```

但下游到底从哪里得到 goalcli？如果仍用 `go run ./cmd/goalcli`，下游必须复制工具代码；如果用远程 `go run github.com/ZoneCNH/xlib-standard/cmd/goalcli@latest`，则不可重复且违反供应链固定规则。

### 修正

新增：

```text
.config/xlib/tools/goalcli-distribution.json
```

推荐策略：

```json
{
  "schema_version": "2.0",
  "tool": "goalcli",
  "distribution": "pinned_module_tool",
  "module": "github.com/ZoneCNH/xlib-standard/cmd/goalcli",
  "version": "v0.6.0",
  "commit": "<commit>",
  "checksum": "<sha256>",
  "forbid_latest": true,
  "downstream_invocation": "go run github.com/ZoneCNH/xlib-standard/cmd/goalcli@<pinned-version>",
  "offline_override": {
    "env": "XLIB_GOALCLI_BIN",
    "checksum_required": true
  }
}
```

下游 lock 必须记录：

```json
{
  "tooling": {
    "goalcli": {
      "distribution": "pinned_module_tool",
      "version": "v0.6.0",
      "commit": "<commit>",
      "checksum": "<sha256>",
      "forbid_latest": true
    }
  }
}
```

验收：

```text
下游不得使用 @latest。
下游不得使用未校验的本地 goalcli。
下游可以通过 XLIB_GOALCLI_BIN 离线运行，但必须校验 checksum。
```

---

## 7. P1：根目录平台文件补齐分类

必须加入 platform-native：

```text
.gitignore
.gitattributes
.editorconfig
Dockerfile
.dockerignore
go.mod
go.sum
.devcontainer/devcontainer.json
.github/workflows/*.yml
.github/CODEOWNERS
.github/pull_request_template.md
.github/ISSUE_TEMPLATE/*
LICENSE
README.md
```

其中：

```text
README.md / docs/ 是 human projection，不是 machine source。
.gitignore / .dockerignore 是 platform-native policy adapter，必须被 .config 规则校验。
go.mod 是 Go native source，xlib 只 observe_and_validate。
```

---

## 8. P1：路径常量集中治理

新增：

```text
.config/xlib/rules/path-constants.json
```

```json
{
  "schema_version": "2.0",
  "constants": {
    "standard_config": ".config/xlib/xlib.json",
    "target_descriptor": ".config/xlib/xlib-target.json",
    "target_lock": ".config/xlib/xlib-standard.lock.json",
    "target_snapshot": ".config/xlib/standard-snapshot.json",
    "target_adoption_proof": ".config/xlib/adoption-proof.json",
    "target_boundary_report": ".config/xlib/boundary-report.json",
    "target_contract_fingerprint": ".config/xlib/contract-fingerprint.json"
  },
  "forbidden_constants": [
    "xlib-standard.lock",
    ".agent",
    "contracts",
    ".githooks",
    "templates",
    "mk/governance.mk"
  ]
}
```

所有 Go code / shell / Makefile / docs 当前用法中的路径必须从这里 resolve 或被 pathguard 认可。

---

## 9. P1：环境变量覆盖必须受限

允许：

```text
XLIB_CONFIG=.config/xlib/xlib.json
XLIB_LOCK=.config/xlib/xlib-standard.lock.json
XLIB_TARGET=.config/xlib/xlib-target.json
XLIB_GOALCLI_BIN=/abs/path/to/verified/goalcli
```

禁止：

```text
XLIB_CONFIG=contracts/...
XLIB_CONFIG=.agent/...
XLIB_LOCK=xlib-standard.lock
XLIB_GOALCLI_BIN without checksum
```

新增规则：

```json
{
  "schema_version": "2.0",
  "env_overrides": {
    "XLIB_CONFIG": {
      "allowed_prefix": ".config/xlib/",
      "must_end_with": "xlib.json"
    },
    "XLIB_LOCK": {
      "allowed_value": ".config/xlib/xlib-standard.lock.json"
    },
    "XLIB_GOALCLI_BIN": {
      "absolute_path_required": true,
      "checksum_required": true
    }
  }
}
```

---

## 10. P1：standard-snapshot 只能是 effective subset

`snapshot` 不能复制完整 registry，否则下游仍然携带标准源副本。

允许：

```text
effective profile
resolved required gates
resolved forbidden imports/concepts
resolved pathguard rules
schema fingerprints
canonical fingerprint rules
tooling distribution
```

禁止：

```text
all profiles registry
all downstream targets
all release policies unrelated to this target
all templates
all agent runtime
```

新增 snapshot policy：

```json
{
  "schema_version": "2.0",
  "snapshot_policy": {
    "mode": "effective_subset_only",
    "forbid_full_registry_copy": true,
    "include_only_target_profile": true,
    "include_only_required_rules": true,
    "include_schema_fingerprints_not_full_schemas": true
  }
}
```

---

## 11. P1：render output pruning 验收命令

新增验收：

```bash
# 下游不允许携带标准源 registry
test ! -e "$tmp/redisx/.config/xlib/xlib.json"
test ! -e "$tmp/redisx/.config/xlib/profiles"
test ! -e "$tmp/redisx/.config/xlib/downstream"
test ! -e "$tmp/redisx/.config/xlib/templates"
test ! -e "$tmp/redisx/.config/xlib/harness"
test ! -e "$tmp/redisx/.config/xlib/agent"

# 下游只允许 target subset
test -f "$tmp/redisx/.config/xlib/xlib-target.json"
test -f "$tmp/redisx/.config/xlib/xlib-standard.lock.json"
test -f "$tmp/redisx/.config/xlib/standard-snapshot.json"
```

---

## 12. P1：final additional negative tests

```text
TestGeneratedTargetContainsXlibJSONFails
TestGeneratedTargetContainsProfilesRegistryFails
TestGeneratedTargetContainsRulesRegistryFails
TestGeneratedTargetSnapshotContainsAllProfilesFails
TestConfigSymlinkToContractsFails
TestConfigPathEscapesRepoFails
TestCaseInsensitiveContractsFails
TestGoEmbedContractsFails
TestDockerfileCopyContractsFails
TestMakefileGoalcliWithoutConfigFails
TestEnvOverrideToRootLockFails
TestDownstreamGoalcliLatestFails
TestDownstreamMissingGoalcliChecksumFails
TestPlatformNativeDefinesProfileFails
```

---

## 13. Final modified GOAL stack

```text
G0: strict-config-root 最终闭环

  G1: 标准源仓库唯一事实源
    - .config/xlib/xlib.json
    - repo_kind=standard_source
    - all registries in .config/xlib

  G2: 生成下游不是标准源
    - repo_kind=generated_target
    - xlib-target.json
    - lock + effective standard-snapshot
    - no full registry copy

  G3: 平台 native 文件被管控而非误删
    - platform-native-facts.json
    - observed_and_validate
    - no xlib facts in platform files

  G4: 路径安全不能被绕过
    - symlink guard
    - path escape guard
    - case-insensitive forbidden match
    - ASCII-only xlib-controlled paths

  G5: 下游工具链可重复
    - pinned goalcli distribution
    - no @latest
    - offline binary requires checksum

  G6: 旧路径只能出现在迁移文档和负向测试
    - legacy-reference-allowlist
    - go:embed/Dockerfile/Makefile/workflow scan

  G7: release evidence 证明 strict clean
    - release manifest v2
    - canonical fingerprints
    - platform adapter status
```

---

## 14. 最终结论

v3 已经是可执行的 strict 方案，但如果要真正“从根源上解决”，必须补上本文件的 v4 闭环：

```text
1. 区分 standard_source 与 generated_target。
2. render 后 prune，避免下游复制完整标准源。
3. 把 go.mod/.github/.dockerignore 等平台 native 文件纳入 observed fact 模型。
4. pathguard 抵抗 symlink、大小写、路径穿越、go:embed、fixture 误用。
5. 下游 goalcli 分发固定且可离线校验。
6. snapshot 只包含 effective subset，不复制 full registry。
```

最终原则：

> `.config/` 是 xlib 权力中心；平台 native 文件是受监管的执行接口；下游是标准消费者，不是标准副本；任何 legacy path、legacy flag、full registry leak、unverified toolchain 都是违规状态。
