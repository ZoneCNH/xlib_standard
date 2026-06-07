# strict-config-root-v2 遗漏审计与补强清单

> 目标：检查《xlib-standard 严格 `.config/` 根治重构方案》是否还有遗漏、矛盾或执行风险。
> 结论：主方向正确，但需要补齐 15 个关键缺口；其中 P0 有 6 个，否则 strict-config-root 会在下游、平台 adapter、供应链和 release evidence 上出现新漂移。

---

## 0. 总体判断

原方案正确地抓住了根问题：`.config/` 必须成为唯一机器事实源，旧路径和旧参数必须变成违规状态，而不是 deprecated 或 mirror。

但方案仍有三个结构性缺口：

1. **平台强制入口没有完整建模**：`.github/workflows` 被处理了，但 `.github/dependabot.yml`、`.github/rulesets`、CODEOWNERS、`.dockerignore`、`.devcontainer`、Dockerfile、go.mod/go.sum 等没有被统一分类。
2. **下游执行模型前后不一致**：方案说下游只保留 lock/proof，不复制完整标准控制面，但后续验收又要求下游运行 `config-check`。没有 `.config/xlib/xlib.json` 的下游无法执行 config-driven check。
3. **模板与治理数据面的边界不够绝对**：方案禁止 `templates/l2` 和 `mk/governance.mk`，但没有明确禁止根目录 `templates/` 和整个 `mk/` 继续作为机器事实源。

---

## 1. P0 缺口一：平台强制入口必须单独建模

### 问题

“所有机器事实源都进 `.config/`”和“平台自动发现文件”之间存在天然冲突。

以下文件不能简单按 legacy path 删除，否则会失去平台能力：

```text
.github/workflows/*.yml
.github/dependabot.yml
.github/rulesets/*.json
.github/CODEOWNERS 或 CODEOWNERS
.github/pull_request_template.md
.github/ISSUE_TEMPLATE/*
.dockerignore
Dockerfile
.devcontainer/devcontainer.json
go.mod
go.sum
README.md
LICENSE
```

### 修正

新增：

```text
.config/xlib/rules/platform-adapters.json
```

定义三类路径：

```json
{
  "schema_version": "2.0",
  "platform_adapters": [
    {
      "path": ".github/workflows/*.yml",
      "platform": "github_actions",
      "mode": "thin_adapter",
      "allowed": true,
      "standard_facts_allowed": false
    },
    {
      "path": ".github/dependabot.yml",
      "platform": "github_dependabot",
      "mode": "platform_native_config",
      "decision_required": true,
      "allowed_only_if": "dependency_automation_policy.dependabot_enabled"
    },
    {
      "path": ".dockerignore",
      "platform": "docker",
      "mode": "platform_native_config",
      "allowed": true,
      "source_policy": ".config/xlib/rules/docker.json"
    },
    {
      "path": "go.mod",
      "platform": "go_toolchain",
      "mode": "platform_native_config",
      "allowed": true,
      "source_policy": ".config/xlib/rules/toolchain.json"
    }
  ]
}
```

### 原则

```text
平台 adapter 不是向后兼容。
平台 adapter 是外部平台 API。
它允许存在，但必须由 pathguard 限制为 adapter 或平台原生配置，不能成为 xlib 标准事实源。
```

---

## 2. P0 缺口二：Dependabot / Renovate 必须二选一或明确双轨

### 问题

原方案把 `renovate.json` 移到 `.config/renovate/renovate.json`，但没有处理 `.github/dependabot.yml`。如果继续使用 GitHub Dependabot，它只能从 `.github/dependabot.yml` 读取配置；如果严格禁止 `.github/dependabot.yml`，就必须关闭 Dependabot，改用 Renovate CLI / workflow 显式读取 `.config/renovate/renovate.json`。

### 修正选项 A：禁用 Dependabot，统一 Renovate

```json
{
  "schema_version": "2.0",
  "dependency_automation": {
    "mode": "renovate_cli_only",
    "renovate_config": ".config/renovate/renovate.json",
    "dependabot_enabled": false,
    "dependabot_file_forbidden": true
  }
}
```

Pathguard：

```text
.github/dependabot.yml exists -> fail
renovate.json at root exists -> fail
```

GitHub workflow 作为 adapter：

```yaml
- run: renovate --config-file .config/renovate/renovate.json
```

### 修正选项 B：保留 Dependabot 作为平台原生配置

```json
{
  "dependency_automation": {
    "mode": "github_dependabot_native",
    "dependabot_file": ".github/dependabot.yml",
    "dependabot_is_platform_adapter": true,
    "renovate_enabled": false
  }
}
```

但这会形成 `.github/dependabot.yml` 里的事实配置。若坚持“所有 xlib-controlled 事实进 `.config/`”，选项 A 更彻底。

### 推荐

选项 A：**禁用 Dependabot，使用 Renovate CLI，以 `.config/renovate/renovate.json` 为唯一事实源。**

---

## 3. P0 缺口三：下游仓库需要 `standard-snapshot`，否则 lock-only 不可离线验证

### 问题

方案中下游只保留：

```text
.config/xlib/xlib-standard.lock.json
.config/xlib/adoption-proof.json
.config/xlib/boundary-report.json
.config/xlib/contract-fingerprint.json
.config/xlib/profile-plan.json
```

但下游验收又要求：

```bash
GOWORK=off make pathguard
GOWORK=off make config-check
GOWORK=off make boundary
```

如果下游不复制完整 `.config/xlib/xlib.json`，这些命令无法知道：

```text
forbidden paths
forbidden flags
profile rules
boundary rules
schema paths
required gates
```

### 修正

下游不复制完整标准控制面，但必须保存一个**不可编辑标准快照**：

```text
.config/xlib/xlib-standard.lock.json
.config/xlib/standard-snapshot.json
.config/xlib/adoption-proof.json
.config/xlib/boundary-report.json
.config/xlib/contract-fingerprint.json
.config/xlib/profile-plan.json
```

`standard-snapshot.json` 只包含执行下游检查所需的最小规则：

```json
{
  "schema_version": "2.0",
  "standard": {
    "name": "xlib-standard",
    "version": "v0.6.0",
    "commit": "<commit>"
  },
  "profile": {
    "id": "l2.redis",
    "layer": "L2",
    "forbidden_imports": [],
    "forbidden_concepts": [],
    "required_gates": []
  },
  "pathguard": {
    "forbidden_paths": [],
    "forbidden_cli_flags": [],
    "forbidden_source_patterns": []
  },
  "boundary": {},
  "schemas": {
    "xlib_lock_sha256": "<sha256>",
    "adoption_proof_sha256": "<sha256>"
  },
  "snapshot_sha256": "<sha256>"
}
```

下游命令改成 lock/snapshot driven：

```bash
goalcli pathguard --lock .config/xlib/xlib-standard.lock.json
goalcli boundary --lock .config/xlib/xlib-standard.lock.json
goalcli adoption-check --lock .config/xlib/xlib-standard.lock.json --verify
```

不要在下游要求：

```bash
goalcli config check --config .config/xlib/xlib.json
```

除非下游确实复制完整标准控制面。

---

## 4. P0 缺口四：根目录 `templates/` 和 `mk/` 没有完全禁止

### 问题

方案只禁止：

```text
templates/l2
mk/governance.mk
```

但如果目标是 `.config/` 统一事实源，根目录 `templates/` 与 `mk/` 都是机器事实源风险点。

### 修正

pathguard 改为：

```json
{
  "forbidden_paths": [
    "templates",
    "mk",
    "contracts",
    ".agent",
    ".xlib",
    "xlib-standard.lock"
  ]
}
```

模板文件统一进入：

```text
.config/xlib/templates/fragments/**
.config/xlib/templates/render-policy.json
.config/xlib/templates/placeholders.json
.config/xlib/templates/golden-targets.json
```

Make fragments 统一进入：

```text
.config/xlib/mk/governance.mk
.config/xlib/mk/harness.mk
.config/xlib/mk/release.mk
```

根目录 Makefile 只能 include `.config/xlib/mk/*.mk`。

---

## 5. P0 缺口五：release manifest v2 没有定义 strict 字段

### 问题

方案保留 release/ 作为 Evidence artifact，但没有定义 strict-config-root 后 release manifest 必须新增哪些字段。

### 修正

`.config/xlib/schemas/release-manifest.schema.json` 必须要求：

```json
{
  "strict_config_root": {
    "enabled": true,
    "config_root": ".config/xlib",
    "legacy_compatibility": false,
    "pathguard_status": "passed",
    "legacy_path_count": 0,
    "legacy_flag_count": 0
  },
  "config_fingerprints": {
    "xlib_json_sha256": "<sha256>",
    "profiles_sha256": "<sha256>",
    "rules_sha256": "<sha256>",
    "schemas_sha256": "<sha256>",
    "harness_sha256": "<sha256>",
    "templates_sha256": "<sha256>"
  },
  "platform_adapters": {
    "github_workflows": "passed",
    "dockerignore": "passed",
    "go_mod": "passed",
    "dependabot_or_renovate": "passed"
  },
  "downstream_cutover": {
    "required": true,
    "targets_total": 0,
    "targets_passed": 0,
    "targets_blocked": []
  }
}
```

release gate 必须证明：

```text
pathguard passed
strict-check passed
legacy path count = 0
legacy flag count = 0
platform adapter check passed
all golden renders strict clean
```

---

## 6. P0 缺口六：canonical JSON fingerprint 没有规范

### 问题

方案大量依赖 sha256，但没有定义 JSON 如何序列化。Go 的 map 遍历顺序不稳定；如果直接 marshal 任意 map，fingerprint 可能漂移。

### 修正

新增：

```text
.config/xlib/rules/fingerprint.json
```

规则：

```json
{
  "schema_version": "2.0",
  "canonical_json": {
    "object_keys_sorted": true,
    "indent": "",
    "trailing_newline": true,
    "unicode_normalization": "NFC",
    "line_endings": "LF",
    "ignore_fields": [
      "generated_at",
      "last_checked_at"
    ]
  }
}
```

新增命令：

```bash
goalcli fingerprint --config .config/xlib/xlib.json --path .config/xlib/profiles/l2.redis.json
```

所有 lock、snapshot、release manifest 都必须使用同一 canonical encoder。

---

## 7. P1 缺口：供应链安全策略没有完全迁入 `.config`

### 缺失内容

当前方案没有完整定义：

```text
GitHub Actions SHA pinning
Action 来源 tag 注释
禁止 @latest
govulncheck 固定版本
vulncheck schedule
XLIB_ENABLE_VULNCHECK 默认值
golangci-lint 版本
Renovate/Dependabot 策略
go version policy
```

### 修正

新增：

```text
.config/xlib/rules/supply-chain.json
.config/xlib/rules/toolchain.json
.config/xlib/rules/security.json
```

关键字段：

```json
{
  "github_actions": {
    "third_party_actions_must_pin_40_char_sha": true,
    "source_tag_comment_required": true,
    "forbidden_refs": ["@latest", "@main", "@master"]
  },
  "govulncheck": {
    "version": "golang.org/x/vuln/cmd/govulncheck@v1.1.4",
    "weekly_schedule_utc": "Monday 03:17",
    "default_enabled_in_ci": false,
    "force_env": "XLIB_FORCE_VULNCHECK"
  },
  "go": {
    "version_source": "go.mod",
    "committed_go_work_forbidden": true,
    "all_release_gates_require_gowork_off": true
  }
}
```

---

## 8. P1 缺口：`.gitignore` / `.dockerignore` / build context 没有纳入验收

### 问题

release manifest、checksum、cache、runtime state 必须忽略；Docker build context 必须包含 `.config`，同时排除 `.git` 和敏感本地状态。

### 修正

新增：

```text
.config/xlib/rules/gitignore.json
.config/xlib/rules/docker.json
```

必须检查：

```text
release/manifest/latest.json ignored
release/manifest/latest.json.sha256 ignored
.cache/ ignored
.omc/ 或旧 runtime state ignored/forbidden
.env forbidden
.config/xlib/** included in Docker context
.git excluded from Docker context
secrets excluded
```

如果保留根 `.dockerignore` 作为 Docker 平台 adapter，pathguard 必须检查它只符合 `.config/xlib/rules/docker.json`。

---

## 9. P1 缺口：repository rules / branch protection / CODEOWNERS 没有定义

### 问题

当前生成脚本曾要求 `.github/rulesets/protect-main.json`，但新方案只处理了 workflow，没有处理 repository rules、CODEOWNERS、review policy、branch governance。

### 修正

新增：

```text
.config/xlib/rules/repository-governance.json
.config/github/rulesets/protect-main.json
.config/github/CODEOWNERS
```

如果 GitHub 平台需要 `.github/CODEOWNERS` 或 `.github/rulesets`，它们必须被归类为 platform adapter，而不是 legacy mirror。

推荐策略：

```text
规则事实源：.config/xlib/rules/repository-governance.json
平台应用：goalcli repository-rules apply/check
Evidence：release manifest 记录 repository_rules_release_decision
```

---

## 10. P1 缺口：hard-coded path audit 不够强

### 问题

pathguard 只扫普通字符串还不够。Go、Shell、Docker、Makefile、YAML 都可能通过不同方式硬编码旧路径。

### 修正

新增扫描器：

```bash
goalcli legacy-ref-audit --config .config/xlib/xlib.json
```

检查：

```text
普通字符串
//go:embed contracts/*
//go:embed .agent/*
os.ReadFile("contracts/...")
filepath.Join(".agent", ...)
Makefile include mk/*.mk
Docker COPY contracts /
workflow env CONTRACTS_DIR
README 当前用法
模板文件中的旧路径
测试 fixture 中的旧路径
```

允许旧路径只出现在：

```text
docs/migration/strict-config-root-v2.md
CHANGELOG.md
```

并且必须标记为 forbidden legacy。

---

## 11. P1 缺口：cutover / rollback 操作手册不够硬

### 问题

“禁止兼容”意味着跨仓库会出现 flag day。需要明确合并顺序，否则 xlib-standard 合入后，下游全部短时不可用。

### 修正

新增：

```text
docs/migration/strict-config-root-v2.md
```

必须包含：

```text
1. freeze window
2. xlib-standard strict branch
3. tag v0.6.0
4. downstream regeneration branches
5. downstream adoption proof collection
6. ZoneCNH composition consumer update
7. merge order
8. rollback is full revert of strict PR before downstream cutover, not compatibility mode
9. after tag, forward-fix only
```

推荐合并顺序：

```text
xlib-standard strict branch green
-> tag v0.6.0-rc.1
-> regenerate kernel/configx/observex/testkitx
-> regenerate L2 infrastructure libs
-> regenerate L2.5 domain libs
-> apply business governance-only shells
-> update ZoneCNH x.go consumer checks
-> final release v0.6.0
```

---

## 12. P1 缺口：tool distribution model 没有定义

### 问题

下游运行 `goalcli` 的方式不明确。当前 render 是复制仓库式，生成库可能带 `cmd/goalcli`；如果未来改成 fragment render，下游可能没有 `cmd/goalcli`。

### 修正

必须选一个：

### 方案 A：下游复制 `cmd/goalcli` 和必要 internal 包

优点：离线、无外部工具安装。
缺点：下游工具代码会漂移，需要 lock fingerprint。

### 方案 B：独立 `xlibctl` 工具

```bash
go run github.com/ZoneCNH/xlib-standard/cmd/goalcli@v0.6.0 ...
```

优点：下游轻。
缺点：依赖网络或 module cache。

### 推荐

短期方案 A，长期抽出：

```text
github.com/ZoneCNH/xlibctl
```

并在 lock 中记录：

```json
{
  "tooling": {
    "mode": "embedded_goalcli",
    "goalcli_sha256": "<sha256>"
  }
}
```

---

## 13. P2 缺口：测试矩阵缺少 negative tests

### 必须新增负向测试

```text
render_template --enable-governance 必须失败
render_template --layer L2 必须失败
root xlib-standard.lock 存在时 pathguard 必须失败
.agent 存在时 pathguard 必须失败
contracts 存在时 pathguard 必须失败
README 当前用法出现 --layer 必须失败
workflow 使用 @latest action 必须失败
Makefile 调 goalcli contracts 但不传 --config 必须失败
下游没有 standard-snapshot 时 adoption-check 必须失败
lock 写到 root 时必须失败
```

---

## 14. P2 缺口：profile coverage 仍不完整

### 原方案已有

```text
l0.kernel
l1.config
l1.observability
l1.testkit
l2.redis
l2.postgres
l2.kafka
l2.nats
l2.oss
l2.clickhouse
l25.domain-model
contracts.cross-domain
business.governance-only
xgo.consumer
```

### 需要补齐

```text
l1.resilience       -> resiliencx
l1.scheduler        -> schedulex
l2.taos             -> taosx
l2.object-storage   -> ossx
l2.analytics-db     -> clickhousex
l2.message-broker   -> shared kafka/nats base
l25.decimal         -> decimalx 特化，而不是泛 domain-model
l25.exchange        -> domain-exchange 特化
l25.macro           -> domain-macro 特化
```

---

## 15. P2 缺口：语义扫描需要 allowlist schema

### 问题

禁止概念扫描会误伤文档、迁移说明、测试 fixture。

### 修正

新增：

```text
.config/xlib/rules/pathguard-allowlist.json
```

格式：

```json
{
  "schema_version": "2.0",
  "allowlist": [
    {
      "path": "docs/migration/strict-config-root-v2.md",
      "pattern": "--enable-governance",
      "reason": "documents forbidden legacy flag"
    }
  ],
  "allowlist_policy": {
    "reason_required": true,
    "expires_at_required_for_non_migration": true
  }
}
```

---

# 最终补强后的 P0 必改清单

在原方案基础上，必须追加：

```text
1. .config/xlib/rules/platform-adapters.json
2. dependency automation 二选一：禁 Dependabot + Renovate CLI，或明确 Dependabot platform adapter
3. 下游增加 .config/xlib/standard-snapshot.json
4. pathguard 禁止 root templates/ 和 root mk/，不是只禁 templates/l2 和 mk/governance.mk
5. release-manifest v2 增加 strict_config_root / config_fingerprints / platform_adapters / downstream_cutover
6. canonical JSON fingerprint 规则
7. supply-chain/toolchain/security rules
8. gitignore/dockerignore checks
9. repository governance / CODEOWNERS / rulesets policy
10. legacy-ref-audit 扫 go:embed、Docker、Makefile、workflow、fixture、templates
11. strict cutover + rollback guide
12. tooling distribution model
13. negative tests
14. expanded profiles
15. pathguard allowlist schema
```

---

# 修正后的最终判断

原方案方向正确，但还不能算“最彻底”。

真正完整的版本应该从：

```text
strict .config root + pathguard
```

升级为：

```text
strict .config root
+ platform adapter taxonomy
+ downstream standard snapshot
+ canonical fingerprint
+ supply-chain/toolchain rules
+ repository governance rules
+ release manifest v2
+ legacy-ref-audit
+ flag-day cutover protocol
```

这样才不会在删除 `.agent / contracts / root lock` 后，又从 `.github/dependabot.yml`、`.dockerignore`、downstream lock-only check、root templates、hard-coded go:embed 或 release manifest 字段里重新长出第二套事实源。
