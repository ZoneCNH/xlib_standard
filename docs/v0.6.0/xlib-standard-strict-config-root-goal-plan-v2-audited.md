# xlib-standard 严格 `.config/` 根治重构方案

> 版本：strict-config-root-v2
> 目标：禁止向后兼容，统一在 `.config/`，一次性从根源解决标准事实源分散、模板治理漂移、下游采纳不可证明的问题。
> 适用对象：`ZoneCNH/xlib-standard`、由它生成或治理的 Foundation / L2.5 / business governance-only 下游仓库。
> 生成日期：2026-06-06

---

## 0. 总结结论

`xlib-standard` 不应继续作为“模板仓库 + 若干脚本 + 若干文档”的组合存在，而应被重构为：

> **严格 `.config/` 标准控制面 + profile 编译器 + pathguard 架构防回潮系统 + lock/evidence/adoption 证明工厂。**

最终原则：

```text
.config/ 是唯一机器事实源。
旧路径不是兼容入口，而是违规状态。
旧参数不是 deprecated，而是 forbidden。
旧 lock 不是 mirror，而是 forbidden artifact。
docs/ 只做人读投影。
release/ 只做 Evidence 产物。
Makefile / .github 只做平台 adapter。
```

这次重构应作为一次 **breaking release** 执行，例如：

```text
v0.6.0 strict-config-root breaking release
```

不要做：

```text
dual-read
legacy mirror
fallback mode
root lock
contracts/ mirror
.agent/ mirror
--enable-governance
--layer
```

必须做：

```text
strict_config_only
.config/xlib/xlib.json 统一入口
profile-driven render
lock only in .config/xlib/xlib-standard.lock.json
pathguard as P0 gate
strict-render-check
downstream full regeneration
```

---

## 1. 底层本质

当前问题不是目录不好看，而是 **标准事实源分散**：

```text
.agent/          goal runtime / harness / evidence protocol
contracts/       schema / contract
xlib-standard.lock
.githooks/
mk/governance.mk
docs/            matrix / sync policy / generation rules
scripts/         render behavior
Makefile         gate behavior
```

只要这些路径继续承载机器事实，系统就会持续出现：

```text
标准漂移
文档漂移
下游采纳伪证明
模板渲染残留
旧参数回潮
lock 路径分叉
review 依赖人脑记忆
```

根治方案是把所有机器事实源压缩成一个控制面：

```text
.config/xlib/
```

并通过 `pathguard` 禁止旧路径和旧参数再次出现。

---

## 2. 不可违背的基本真理

### 真理 1：兼容层会保留旧系统的权力

只要旧路径还能被读取，它就会重新成为事实源。
只要旧参数还能运行，团队就会绕过 profile。
只要根目录 `xlib-standard.lock` 还能存在，下游就不会真正迁移到 `.config/`。

因此：

```text
legacy path = error
legacy CLI flag = error
legacy generated output = error
legacy docs current usage = error
```

### 真理 2：配置必须集中，入口可以不集中

`.github/workflows/*.yml` 和根目录 `Makefile` 可以保留，因为它们是 GitHub / 本地开发平台入口。
但它们只能做 adapter，不能承载标准事实。

正确边界：

```text
.config/xlib/     标准事实源
Makefile          adapter
.github/workflows adapter
docs/             projection
release/          evidence artifact
```

### 真理 3：schema、contract、harness、agent runtime 必须同源

不允许：

```text
contracts/
.agent/
.xlib/
mk/governance.mk
```

作为机器事实源继续存在。

它们的内容必须迁移到：

```text
.config/xlib/schemas/
.config/xlib/contracts/
.config/xlib/harness/
.config/xlib/agent/
.config/xlib/mk/
```

### 真理 4：ZoneCNH 架构规则必须机器化

需要被机器固化的架构规则：

```text
L2.5 只放共享领域值对象、枚举、语义模型。
L2.5 不放 Provider 实现、策略逻辑、执行策略。
contracts 只放跨域稳定端口、事件协议、DTO。
x.go 只做配置加载、依赖创建、模块 wiring、生命周期。
x.go 不做因子计算、信号判断、风控规则、订单路由。
signal-factory / optimizer 必须先经过 risk-engine。
禁止绕过 risk-engine 直接调用 order-engine。
执行反馈通过 fills / positions / PnL / exposure events。
```

这些规则进入：

```text
.config/xlib/rules/architecture.json
.config/xlib/rules/boundary.json
.config/xlib/profiles/*.json
```

---

## 3. 最终目标架构

### 3.1 根目录最终形态

```text
.
├── .config/
│   ├── xlib/
│   ├── git/
│   ├── golangci/
│   └── renovate/
├── .github/
├── cmd/
├── docs/
├── examples/
├── internal/
├── pkg/
├── release/
├── scripts/
├── testkit/
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### 3.2 `.config/` 最终形态

```text
.config/
├── xlib/
│   ├── xlib.json
│   ├── profiles/
│   │   ├── registry.json
│   │   ├── l0.kernel.json
│   │   ├── l1.config.json
│   │   ├── l1.observability.json
│   │   ├── l1.testkit.json
│   │   ├── l2.redis.json
│   │   ├── l2.postgres.json
│   │   ├── l2.kafka.json
│   │   ├── l2.nats.json
│   │   ├── l2.oss.json
│   │   ├── l2.clickhouse.json
│   │   ├── l25.domain-model.json
│   │   ├── contracts.cross-domain.json
│   │   ├── business.governance-only.json
│   │   └── xgo.consumer.json
│   ├── capabilities/
│   │   └── registry.json
│   ├── downstream/
│   │   ├── targets.json
│   │   ├── adoption-status.json
│   │   └── release-sync-policy.json
│   ├── rules/
│   │   ├── pathguard.json
│   │   ├── architecture.json
│   │   ├── boundary.json
│   │   ├── dependency.json
│   │   ├── security.json
│   │   ├── release.json
│   │   ├── naming.json
│   │   └── evidence.json
│   ├── schemas/
│   │   ├── xlib.schema.json
│   │   ├── profile.schema.json
│   │   ├── capability.schema.json
│   │   ├── downstream-target.schema.json
│   │   ├── adoption-status.schema.json
│   │   ├── xlib-lock.schema.json
│   │   ├── release-manifest.schema.json
│   │   ├── config.schema.json
│   │   ├── error.schema.json
│   │   ├── health.schema.json
│   │   └── metrics.schema.json
│   ├── contracts/
│   │   ├── config.contract.json
│   │   ├── error.contract.json
│   │   ├── health.contract.json
│   │   ├── metrics.contract.json
│   │   ├── release-manifest.contract.json
│   │   └── adoption-proof.contract.json
│   ├── harness/
│   │   ├── gates.json
│   │   ├── required.json
│   │   ├── extended.json
│   │   ├── release.json
│   │   ├── score.json
│   │   └── docker.json
│   ├── agent/
│   │   ├── runtime.json
│   │   ├── object-model.json
│   │   ├── state-machine.json
│   │   ├── evidence-protocol.json
│   │   ├── traceability-matrix.json
│   │   ├── risk-register.json
│   │   ├── decision-log.json
│   │   └── rollback-protocol.json
│   ├── templates/
│   │   ├── render-policy.json
│   │   ├── placeholders.json
│   │   ├── fragments.json
│   │   └── golden-targets.json
│   └── mk/
│       ├── governance.mk
│       ├── release.mk
│       └── harness.mk
├── git/
│   └── hooks/
│       ├── pre-commit
│       └── pre-push
├── golangci/
│   └── golangci.yml
└── renovate/
    └── renovate.json
```

---

## 4. 必须删除的旧路径

一次性删除，不做 mirror，不做 fallback：

```bash
git rm -r .agent
git rm -r .xlib
git rm -r contracts
git rm -r .githooks
git rm -r templates/l2
git rm -f xlib-standard.lock
git rm -f releasemanifest
git rm -f .golangci.yml
git rm -f renovate.json
git rm -f mk/governance.mk
```

迁移映射：

| 旧路径 | 新路径 |
|---|---|
| `.agent/harness/harness.yaml` | `.config/xlib/harness/gates.json` |
| `.agent/traceability/*` | `.config/xlib/agent/traceability-matrix.json` |
| `.agent/*evidence*` | `.config/xlib/agent/evidence-protocol.json` |
| `.xlib/*` | `.config/xlib/*` |
| `contracts/*.schema.json` | `.config/xlib/schemas/*.schema.json` |
| `contracts/*metrics*` | `.config/xlib/contracts/metrics.contract.json` |
| `xlib-standard.lock` | `.config/xlib/xlib-standard.lock.json` |
| `templates/l2/*` | `.config/xlib/templates/fragments/l2/*` 或 `templates/fragments/l2/*` |
| `.githooks/*` | `.config/git/hooks/*` |
| `.golangci.yml` | `.config/golangci/golangci.yml` |
| `renovate.json` | `.config/renovate/renovate.json` |
| `mk/governance.mk` | `.config/xlib/mk/governance.mk` |

---

## 5. `.config/xlib/xlib.json`

```json
{
  "schema_version": "2.0",
  "standard": {
    "name": "xlib-standard",
    "config_root": ".config/xlib",
    "mode": "strict_config_only",
    "legacy_compatibility": false
  },
  "paths": {
    "profiles": ".config/xlib/profiles/registry.json",
    "capabilities": ".config/xlib/capabilities/registry.json",
    "downstream_targets": ".config/xlib/downstream/targets.json",
    "adoption_status": ".config/xlib/downstream/adoption-status.json",
    "rules": ".config/xlib/rules",
    "schemas": ".config/xlib/schemas",
    "contracts": ".config/xlib/contracts",
    "harness": ".config/xlib/harness",
    "agent": ".config/xlib/agent",
    "templates": ".config/xlib/templates",
    "mk": ".config/xlib/mk",
    "git_hooks": ".config/git/hooks",
    "golangci": ".config/golangci/golangci.yml",
    "renovate": ".config/renovate/renovate.json"
  },
  "policy": {
    "config_is_only_source_of_truth": true,
    "legacy_paths_forbidden": true,
    "legacy_cli_flags_forbidden": true,
    "generated_docs_must_match_config": true,
    "release_artifacts_are_not_source": true,
    "docs_are_not_machine_source": true,
    "root_lock_forbidden": true,
    "contracts_root_forbidden": true,
    "agent_root_forbidden": true,
    "xlib_root_forbidden": true
  },
  "platform_adapters": {
    "github_workflows": ".github/workflows",
    "makefile": "Makefile",
    "adapter_only": true,
    "standard_facts_allowed": false
  }
}
```

---

## 6. `pathguard.json`

`pathguard` 是根治方案的 P0。没有它，旧路径必然回潮。

```json
{
  "schema_version": "2.0",
  "mode": "strict",
  "forbidden_paths": [
    ".agent",
    ".xlib",
    "contracts",
    "xlib-standard.lock",
    "releasemanifest",
    "templates/l2",
    "mk/governance.mk",
    ".githooks",
    ".golangci.yml",
    "renovate.json"
  ],
  "forbidden_cli_flags": [
    "--enable-governance",
    "--layer"
  ],
  "forbidden_source_patterns": [
    ".agent/",
    ".xlib/",
    "contracts/",
    "xlib-standard.lock",
    "core.hooksPath=.githooks",
    "core.hooksPath .githooks",
    "--enable-governance",
    "--layer L0",
    "--layer L1",
    "--layer L2"
  ],
  "allowed_platform_adapters": [
    ".github/workflows/ci.yml",
    ".github/workflows/release.yml",
    ".github/workflows/security.yml",
    "Makefile"
  ],
  "platform_adapter_policy": {
    "must_not_define_standard_facts": true,
    "must_call_config_driven_commands": true,
    "required_config_arg": "--config .config/xlib/xlib.json"
  },
  "allowed_mentions": [
    "docs/migration/strict-config-root-v2.md",
    "CHANGELOG.md"
  ]
}
```

`pathguard` 检查：

```text
1. forbidden_paths 不存在。
2. forbidden_cli_flags 不在当前用法中出现。
3. forbidden_source_patterns 不在 scripts、Makefile、docs、workflow 当前用法中出现。
4. `.github/workflows` 只做 adapter。
5. Makefile 只做 adapter。
6. docs 不能把 legacy path 描述成当前用法。
7. release/ 不能反向成为 source。
```

---

## 7. Profile 体系

### 7.1 `profiles/registry.json`

```json
{
  "schema_version": "2.0",
  "profiles": [
    {
      "id": "l0.kernel",
      "file": ".config/xlib/profiles/l0.kernel.json"
    },
    {
      "id": "l1.config",
      "file": ".config/xlib/profiles/l1.config.json"
    },
    {
      "id": "l1.observability",
      "file": ".config/xlib/profiles/l1.observability.json"
    },
    {
      "id": "l2.redis",
      "file": ".config/xlib/profiles/l2.redis.json"
    },
    {
      "id": "l25.domain-model",
      "file": ".config/xlib/profiles/l25.domain-model.json"
    },
    {
      "id": "contracts.cross-domain",
      "file": ".config/xlib/profiles/contracts.cross-domain.json"
    },
    {
      "id": "business.governance-only",
      "file": ".config/xlib/profiles/business.governance-only.json"
    },
    {
      "id": "xgo.consumer",
      "file": ".config/xlib/profiles/xgo.consumer.json"
    }
  ]
}
```

### 7.2 `l2.redis.json`

```json
{
  "schema_version": "2.0",
  "id": "l2.redis",
  "layer": "L2",
  "kind": "infrastructure_adapter",
  "description": "Redis foundation adapter profile.",
  "required_capabilities": [
    "explicit_config",
    "config_redaction",
    "typed_error",
    "health_check",
    "metrics_contract",
    "context_aware_operations",
    "lifecycle_close_idempotency",
    "docker_toolchain",
    "repository_governance",
    "release_evidence"
  ],
  "allowed_imports": [
    "std",
    "github.com/ZoneCNH/kernel",
    "github.com/ZoneCNH/configx",
    "github.com/ZoneCNH/observex",
    "github.com/redis/go-redis"
  ],
  "forbidden_imports": [
    "github.com/ZoneCNH/x.go",
    "github.com/bytechainx/x.go",
    "github.com/ZoneCNH/market-data",
    "github.com/ZoneCNH/factor-engine",
    "github.com/ZoneCNH/signal-factory",
    "github.com/ZoneCNH/risk-engine",
    "github.com/ZoneCNH/order-engine"
  ],
  "forbidden_concepts": [
    "business_key_semantics",
    "application_cache_policy",
    "strategy_state",
    "order_routing",
    "hidden_global_client",
    "implicit_production_secret"
  ],
  "required_gates": [
    "pathguard",
    "config-check",
    "fmt",
    "vet",
    "test",
    "race",
    "docs-check",
    "contracts",
    "boundary",
    "docker-toolchain-check",
    "release-evidence-check",
    "adoption-check"
  ],
  "lock_policy": {
    "path": ".config/xlib/xlib-standard.lock.json",
    "root_lock_forbidden": true
  }
}
```

### 7.3 `l25.domain-model.json`

```json
{
  "schema_version": "2.0",
  "id": "l25.domain-model",
  "layer": "L2.5",
  "kind": "shared_domain_semantics",
  "description": "Shared domain value object and semantic model profile.",
  "required_capabilities": [
    "value_object",
    "enum_contract",
    "validation",
    "canonical_json",
    "golden_samples",
    "compatibility_tests",
    "fuzz_smoke",
    "repository_governance",
    "release_evidence"
  ],
  "allowed_imports": [
    "std",
    "github.com/ZoneCNH/decimalx"
  ],
  "forbidden_imports": [
    "github.com/ZoneCNH/x.go",
    "github.com/bytechainx/x.go",
    "github.com/ZoneCNH/binance",
    "github.com/ZoneCNH/okx",
    "github.com/ZoneCNH/bybit",
    "github.com/ZoneCNH/market-data",
    "github.com/ZoneCNH/factor-engine",
    "github.com/ZoneCNH/risk-engine",
    "github.com/ZoneCNH/order-engine"
  ],
  "forbidden_concepts": [
    "provider_implementation",
    "exchange_sdk_runtime",
    "strategy_logic",
    "execution_policy",
    "database_access",
    "message_bus_runtime",
    "application_wiring",
    "production_secret"
  ],
  "required_gates": [
    "pathguard",
    "config-check",
    "fmt",
    "vet",
    "test",
    "property",
    "golden",
    "fuzz-smoke",
    "contracts",
    "boundary",
    "release-evidence-check"
  ],
  "lock_policy": {
    "path": ".config/xlib/xlib-standard.lock.json",
    "root_lock_forbidden": true
  }
}
```

### 7.4 `business.governance-only.json`

```json
{
  "schema_version": "2.0",
  "id": "business.governance-only",
  "layer": "Business",
  "kind": "governance_shell",
  "description": "Business repository governance shell without generated business implementation.",
  "required_capabilities": [
    "repository_governance",
    "release_evidence"
  ],
  "allowed_imports": [
    "std",
    "github.com/ZoneCNH/contracts",
    "github.com/ZoneCNH/decimalx",
    "github.com/ZoneCNH/domain-market",
    "github.com/ZoneCNH/domain-exchange",
    "github.com/ZoneCNH/domain-macro",
    "github.com/ZoneCNH/kernel",
    "github.com/ZoneCNH/configx",
    "github.com/ZoneCNH/observex"
  ],
  "forbidden_imports": [
    "github.com/ZoneCNH/x.go"
  ],
  "forbidden_concepts": [
    "generated_business_logic",
    "generated_trading_strategy",
    "generated_risk_rule",
    "generated_order_routing",
    "generated_provider_runtime"
  ],
  "required_gates": [
    "pathguard",
    "config-check",
    "contracts",
    "boundary",
    "docs-check",
    "release-evidence-check",
    "adoption-check"
  ],
  "lock_policy": {
    "path": ".config/xlib/xlib-standard.lock.json",
    "root_lock_forbidden": true
  }
}
```

### 7.5 `xgo.consumer.json`

```json
{
  "schema_version": "2.0",
  "id": "xgo.consumer",
  "layer": "Composition",
  "kind": "composition_root_consumer",
  "description": "x.go composition root checks.",
  "required_capabilities": [
    "repository_governance",
    "release_evidence"
  ],
  "allowed_concepts": [
    "config_load",
    "dependency_creation",
    "module_wiring",
    "lifecycle_management",
    "graceful_shutdown"
  ],
  "forbidden_concepts": [
    "factor_calculation",
    "signal_judgement",
    "risk_rule_body",
    "order_routing_body",
    "provider_business_logic"
  ],
  "required_gates": [
    "pathguard",
    "config-check",
    "architecture-check",
    "boundary",
    "docs-check",
    "release-evidence-check"
  ],
  "lock_policy": {
    "path": ".config/xlib/xlib-standard.lock.json",
    "root_lock_forbidden": true
  }
}
```

---

## 8. Capabilities Registry

```json
{
  "schema_version": "2.0",
  "capabilities": [
    {
      "id": "explicit_config",
      "description": "Public API accepts explicit config and never reads production secrets implicitly.",
      "required_contracts": [
        ".config/xlib/schemas/config.schema.json"
      ],
      "required_tests": [
        "config_validation",
        "secret_redaction"
      ]
    },
    {
      "id": "config_redaction",
      "description": "Sensitive fields are redacted from errors, logs, examples and evidence.",
      "required_tests": [
        "redaction_golden"
      ]
    },
    {
      "id": "typed_error",
      "description": "Errors expose stable machine-readable kinds.",
      "required_contracts": [
        ".config/xlib/schemas/error.schema.json"
      ]
    },
    {
      "id": "health_check",
      "description": "Health checks expose a stable contract.",
      "required_contracts": [
        ".config/xlib/schemas/health.schema.json"
      ]
    },
    {
      "id": "metrics_contract",
      "description": "Metrics names and labels follow stable low-cardinality contract.",
      "required_contracts": [
        ".config/xlib/schemas/metrics.schema.json"
      ]
    },
    {
      "id": "repository_governance",
      "description": "Repository emits strict .config lock and adoption proof.",
      "required_files": [
        ".config/xlib/xlib-standard.lock.json",
        ".config/xlib/adoption-proof.json",
        ".config/xlib/profile-plan.json"
      ],
      "forbidden_files": [
        "xlib-standard.lock",
        ".agent",
        "contracts"
      ]
    },
    {
      "id": "release_evidence",
      "description": "Release evidence is generated under release/ and never used as config source."
    }
  ]
}
```

---

## 9. Architecture Rules

`.config/xlib/rules/architecture.json`：

```json
{
  "schema_version": "2.0",
  "rules": [
    {
      "id": "xgo-composition-root-only",
      "type": "concept_ban",
      "target": "x.go",
      "allowed": [
        "config_load",
        "dependency_creation",
        "module_wiring",
        "lifecycle_management",
        "graceful_shutdown"
      ],
      "forbidden": [
        "factor_calculation",
        "signal_judgement",
        "risk_rule_body",
        "order_routing_body",
        "provider_business_logic"
      ],
      "required_evidence": [
        "wiring_lifecycle_test"
      ]
    },
    {
      "id": "l25-domain-shared-only",
      "type": "layer_concept_ban",
      "target_layer": "L2.5",
      "forbidden": [
        "provider_implementation",
        "strategy_logic",
        "execution_policy",
        "exchange_sdk_runtime",
        "database_access",
        "message_bus_runtime"
      ]
    },
    {
      "id": "risk-engine-required-before-order-engine",
      "type": "flow_contract",
      "from": [
        "signal-factory",
        "optimizer"
      ],
      "must_pass": [
        "risk-engine"
      ],
      "before": [
        "order-engine"
      ],
      "forbidden": [
        "direct_order_engine_call",
        "direct_exchange_sdk_call"
      ],
      "required_evidence": [
        "paper_trade_path",
        "contract_trace"
      ]
    },
    {
      "id": "execution-feedback-via-events",
      "type": "feedback_contract",
      "from_domain": "execution",
      "to_domain": "decision",
      "allowed_medium": [
        "fills",
        "positions",
        "pnl",
        "exposure_events"
      ],
      "forbidden": [
        "sync_import_strategy_internal",
        "sync_import_backtest_internal"
      ]
    }
  ]
}
```

---

## 10. Harness Gates

`.config/xlib/harness/gates.json`：

```json
{
  "schema_version": "2.0",
  "gates": [
    {
      "id": "pathguard",
      "command": "GOWORK=off go run ./cmd/goalcli pathguard --config .config/xlib/xlib.json",
      "phase": "p0",
      "required": true
    },
    {
      "id": "config-check",
      "command": "GOWORK=off go run ./cmd/goalcli config check --config .config/xlib/xlib.json",
      "phase": "p0",
      "required": true
    },
    {
      "id": "profile-check",
      "command": "GOWORK=off go run ./cmd/goalcli profile check --config .config/xlib/xlib.json",
      "phase": "p0",
      "required": true
    },
    {
      "id": "downstream-check",
      "command": "GOWORK=off go run ./cmd/goalcli downstream check --config .config/xlib/xlib.json",
      "phase": "p0",
      "required": true
    },
    {
      "id": "contracts",
      "command": "GOWORK=off go run ./cmd/goalcli contracts --config .config/xlib/xlib.json",
      "phase": "required",
      "required": true
    },
    {
      "id": "boundary",
      "command": "GOWORK=off go run ./cmd/goalcli boundary --config .config/xlib/xlib.json",
      "phase": "required",
      "required": true
    },
    {
      "id": "docs-check",
      "command": "GOWORK=off go run ./cmd/goalcli docs-check --config .config/xlib/xlib.json",
      "phase": "required",
      "required": true
    },
    {
      "id": "release-evidence-check",
      "command": "GOWORK=off go run ./cmd/goalcli release-evidence-check --config .config/xlib/xlib.json",
      "phase": "release",
      "required": true
    }
  ]
}
```

---

## 11. `goalcli` 新命令

必须新增或改造：

```bash
goalcli pathguard --config .config/xlib/xlib.json
goalcli strict-check --config .config/xlib/xlib.json

goalcli config check --config .config/xlib/xlib.json
goalcli config resolve --config .config/xlib/xlib.json --profile l2.redis --field layer
goalcli config resolve --config .config/xlib/xlib.json --profile l2.redis --field has_capability:repository_governance

goalcli profile check --config .config/xlib/xlib.json
goalcli downstream check --config .config/xlib/xlib.json

goalcli lock write \
  --config .config/xlib/xlib.json \
  --profile l2.redis \
  --module-name redisx \
  --module-path github.com/ZoneCNH/redisx \
  --package-name redisx \
  --standard-version v0.6.0 \
  --standard-commit "$(git rev-parse HEAD)" \
  --out .config/xlib/xlib-standard.lock.json

goalcli strict-render-check \
  --repo /tmp/redisx \
  --lock /tmp/redisx/.config/xlib/xlib-standard.lock.json

goalcli contracts --config .config/xlib/xlib.json
goalcli boundary --config .config/xlib/xlib.json
goalcli adoption-check --config .config/xlib/xlib.json --verify
goalcli docs-check --config .config/xlib/xlib.json
```

---

## 12. 新增内部包

```text
internal/configroot/
├── model.go
├── load.go
├── validate.go
├── resolve.go
├── checksum.go
└── errors.go

internal/pathguard/
├── model.go
├── load.go
├── check_paths.go
├── check_content.go
├── check_adapters.go
└── report.go

internal/profile/
├── model.go
├── registry.go
├── load.go
├── validate.go
└── resolve.go

internal/downstream/
├── model.go
├── load.go
├── validate.go
└── report.go

internal/lockfile/
├── model.go
├── write.go
├── read.go
├── fingerprint.go
└── validate.go

internal/strictcheck/
├── check.go
├── report.go
└── rules.go

internal/configdocs/
├── render.go
├── check.go
└── markers.go
```

职责：

| 包 | 职责 |
|---|---|
| `configroot` | 读取和校验 `.config/xlib/xlib.json` |
| `pathguard` | 禁止旧路径、旧参数、旧事实源回潮 |
| `profile` | 读取 profile registry 和 profile 文件 |
| `downstream` | 校验 targets 与 adoption status |
| `lockfile` | 写入和校验 `.config/xlib/xlib-standard.lock.json` |
| `strictcheck` | 聚合严格模式检查 |
| `configdocs` | 从 `.config` 生成 / 校验 docs 投影 |

---

## 13. `render_template.sh` 最终用法

### 13.1 新命令

```bash
scripts/render_template.sh \
  --config .config/xlib/xlib.json \
  --profile l2.redis \
  --module-name redisx \
  --module-path github.com/ZoneCNH/redisx \
  --package-name redisx \
  --standard-version v0.6.0 \
  --standard-commit "$(git rev-parse HEAD)" \
  --out ../redisx
```

### 13.2 禁止参数

```bash
--enable-governance
--layer
```

脚本开头必须直接失败：

```bash
for arg in "$@"; do
  case "$arg" in
    --enable-governance|--layer)
      echo "ERROR: legacy flag $arg is forbidden. Use --config and --profile." >&2
      exit 2
      ;;
  esac
done
```

### 13.3 渲染流程

```text
1. pathguard forbids legacy flags.
2. parse --config / --profile / target fields.
3. goalcli config check.
4. goalcli profile check.
5. goalcli config resolve --field layer.
6. copy/render template.
7. goalcli lock write.
8. generate profile-plan/adoption-proof/boundary-report/contract-fingerprint placeholders.
9. goalcli strict-render-check.
```

### 13.4 lock 只写这里

```text
.config/xlib/xlib-standard.lock.json
```

禁止写：

```text
xlib-standard.lock
```

---

## 14. 下游仓库最终形态

下游只保留 strict proof，不复制完整标准控制面：

```text
.config/
└── xlib/
    ├── xlib-standard.lock.json
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
mk/governance.mk
```

下游 `adoption-check` 必须证明：

```text
1. .config/xlib/xlib-standard.lock.json 存在。
2. root xlib-standard.lock 不存在。
3. .agent 不存在。
4. contracts 不存在。
5. profile_id 可解析。
6. profile fingerprint 与标准源一致。
7. boundary-report 当前生成。
8. contract-fingerprint 当前生成。
9. adoption-proof 当前生成。
10. release-evidence-check 通过。
```

---

## 15. Makefile 严格 adapter

Makefile 顶部：

```makefile
XLIB_CONFIG ?= .config/xlib/xlib.json
GOALCLI ?= GOWORK=off go run ./cmd/goalcli

include .config/xlib/mk/governance.mk
include .config/xlib/mk/harness.mk
include .config/xlib/mk/release.mk
```

核心 target：

```makefile
.PHONY: pathguard
pathguard:
	$(GOALCLI) pathguard --config $(XLIB_CONFIG)

.PHONY: strict-check
strict-check:
	$(GOALCLI) strict-check --config $(XLIB_CONFIG)

.PHONY: config-check
config-check:
	$(GOALCLI) config check --config $(XLIB_CONFIG)

.PHONY: profile-check
profile-check:
	$(GOALCLI) profile check --config $(XLIB_CONFIG)

.PHONY: downstream-check
downstream-check:
	$(GOALCLI) downstream check --config $(XLIB_CONFIG)

.PHONY: boundary
boundary:
	$(GOALCLI) boundary --config $(XLIB_CONFIG)

.PHONY: contracts
contracts:
	$(GOALCLI) contracts --config $(XLIB_CONFIG)

.PHONY: adoption-check
adoption-check:
	$(GOALCLI) adoption-check --config $(XLIB_CONFIG) --verify

.PHONY: docs-check
docs-check:
	$(GOALCLI) docs-check --config $(XLIB_CONFIG)

.PHONY: ci
ci: pathguard strict-check config-check profile-check downstream-check fmt vet lint test race boundary architecture domain secret-check security contracts governance-check debt score rules-verify

.PHONY: release-check
release-check: pathguard strict-check config-check profile-check downstream-check require-gowork-off ci integration dependency-check standard-impact-check docs-check docs-drift score-check governance-check p1-governance-check p2-runtime-check debt-evidence fact-audit
	CHECK_STATUS=passed $(MAKE) evidence
	$(MAKE) release-evidence-hash
	$(MAKE) release-evidence-check
	$(MAKE) release-evidence-checksum-check
```

---

## 16. Git hooks 迁移

原 `.githooks` 删除。新路径：

```text
.config/git/hooks/
```

Makefile：

```makefile
.PHONY: install-hooks
install-hooks:
	@git config core.hooksPath .config/git/hooks
	@echo "✅ git hooks enabled: core.hooksPath=.config/git/hooks"

.PHONY: doctor-hooks
doctor-hooks:
	@[ "$$(git config --get core.hooksPath)" = ".config/git/hooks" ] || { \
		echo "ERROR: core.hooksPath must be .config/git/hooks"; \
		echo "run: make install-hooks"; \
		exit 1; \
	}
	@echo "✅ hooks config ok"
```

`pathguard` 禁止：

```text
.githooks
core.hooksPath=.githooks
core.hooksPath .githooks
```

---

## 17. GitHub Actions 规则

`.github/workflows` 可以保留，但只做平台 adapter。

允许：

```yaml
- run: GOWORK=off make ci
- run: GOWORK=off make release-check
```

禁止：

```yaml
env:
  SCORE_MIN: "9.8"
  PROFILE: l2.redis
  CONTRACTS_DIR: contracts
```

禁止：

```yaml
- run: go run ./cmd/goalcli contracts
```

必须是：

```yaml
- run: GOWORK=off go run ./cmd/goalcli contracts --config .config/xlib/xlib.json
```

或：

```yaml
- run: GOWORK=off make contracts
```

---

## 18. Docs 规则

`docs/` 不是事实源，只是人读投影。

必须改：

```text
README.md
docs/generation.md
docs/standard/xlib-standard.md
docs/standard/module-boundary.md
docs/standard/evidence-protocol.md
docs/downstream-matrix.md
docs/downstream-sync-policy.md
docs/release.md
```

删除当前用法：

```text
--enable-governance
--layer
xlib-standard.lock
.agent/harness/harness.yaml
contracts/
.githooks
mk/governance.mk
```

替换成：

```text
--config .config/xlib/xlib.json
--profile <profile-id>
.config/xlib/xlib-standard.lock.json
.config/xlib/harness/gates.json
.config/xlib/schemas/
.config/git/hooks
.config/xlib/mk/governance.mk
```

生成块格式：

```markdown
<!-- xlib-generated:start source=.config/xlib/downstream/targets.json sha256=<sha256> -->
...
<!-- xlib-generated:end -->
```

---

## 19. JSON 作为唯一机器格式

第一阶段 canonical config 全部使用 JSON：

```text
.config/xlib/**/*.json
```

原因：

```text
1. Go 标准库可直接解析。
2. 不需要 YAML parser。
3. 不引入新的供应链变量。
4. schemas / contracts / release manifest 本身更接近 JSON 生态。
```

策略：

```text
- 不引入 YAML parser。
- 不引入 JSON Schema runtime validator。
- schema 文件保留为 contract/documentation。
- 真实校验用 Go 结构校验。
```

---

## 20. GOAL 总目标树

```text
G0: xlib-standard strict-config-root-v2 根治
  G1: 所有机器事实源统一进入 .config/
    G1.1: xlib.json 总入口
    G1.2: profiles/capabilities/downstream/rules/schemas/contracts/harness/agent/templates/mk 入 .config
    G1.3: hooks/lint/renovate 入 .config
  G2: 旧路径彻底非法化
    G2.1: 删除 .agent/.xlib/contracts/.githooks/root lock
    G2.2: pathguard 检测 forbidden path
    G2.3: docs/scripts/Makefile/workflows 检测 forbidden pattern
  G3: 渲染改为 profile-driven
    G3.1: render_template.sh 只接受 --config + --profile
    G3.2: layer 从 profile resolve
    G3.3: governance 从 profile capabilities resolve
    G3.4: lock 只写 .config/xlib/xlib-standard.lock.json
  G4: 下游采纳 proof 化
    G4.1: 下游只保留 .config/xlib/*.json proof
    G4.2: adoption-check 验证 no legacy path
    G4.3: drift-check 比较 profile/rules/schema/template fingerprint
  G5: ZoneCNH 架构规则机器化
    G5.1: L2.5 不含 provider/strategy/execution
    G5.2: x.go composition-root-only
    G5.3: risk-engine 必经
    G5.4: execution feedback via events
  G6: release gate 收口
    G6.1: pathguard as P0
    G6.2: strict-check as aggregate
    G6.3: golden downstream strict render
```

---

## 21. Breaking PR 计划

PR 名称：

```text
strict-config-root-v2
```

建议一个 PR，内部拆 12 个 commit，不拆多 PR。
理由：拆多 PR 会产生中间兼容态，而本方案要求禁止兼容。

### Commit 1：add `.config/xlib` strict control plane

新增：

```text
.config/xlib/xlib.json
.config/xlib/profiles/*
.config/xlib/capabilities/registry.json
.config/xlib/downstream/*
.config/xlib/rules/*
.config/xlib/schemas/*
.config/xlib/contracts/*
.config/xlib/harness/*
.config/xlib/agent/*
.config/xlib/templates/*
.config/xlib/mk/*
.config/git/hooks/*
.config/golangci/golangci.yml
.config/renovate/renovate.json
```

验收：

```bash
test -f .config/xlib/xlib.json
test -f .config/xlib/rules/pathguard.json
```

### Commit 2：add `configroot/profile/downstream/lockfile` packages

新增：

```text
internal/configroot/*
internal/profile/*
internal/downstream/*
internal/lockfile/*
```

验收：

```bash
GOWORK=off go test ./internal/configroot ./internal/profile ./internal/downstream ./internal/lockfile
```

### Commit 3：add `pathguard` and `strict-check`

新增：

```text
internal/pathguard/*
internal/strictcheck/*
cmd/goalcli/pathguard.go
cmd/goalcli/strict_check.go
```

验收：

```bash
GOWORK=off go run ./cmd/goalcli pathguard --config .config/xlib/xlib.json
GOWORK=off go run ./cmd/goalcli strict-check --config .config/xlib/xlib.json
```

### Commit 4：delete legacy paths

执行：

```bash
git rm -r .agent
git rm -r .xlib
git rm -r contracts
git rm -r .githooks
git rm -r templates/l2
git rm -f xlib-standard.lock
git rm -f releasemanifest
git rm -f .golangci.yml
git rm -f renovate.json
git rm -f mk/governance.mk
```

验收：

```bash
test ! -e .agent
test ! -e .xlib
test ! -e contracts
test ! -e xlib-standard.lock
test ! -e .githooks
test ! -e mk/governance.mk
```

### Commit 5：Makefile adapter-only

修改：

```text
Makefile
```

验收：

```bash
GOWORK=off make pathguard
GOWORK=off make config-check
```

### Commit 6：goalcli config-driven migration

修改所有相关命令：

```text
boundary
contracts
adoption-check
docs-check
evidence
release-evidence-check
standard-impact-check
score
```

统一加：

```bash
--config .config/xlib/xlib.json
```

验收：

```bash
GOWORK=off make contracts
GOWORK=off make boundary
GOWORK=off make docs-check
```

### Commit 7：render_template strict profile-driven

修改：

```text
scripts/render_template.sh
scripts/check_rendered_template.sh
```

验收：

```bash
! scripts/render_template.sh --enable-governance
! scripts/render_template.sh --layer L2
```

### Commit 8：lock write and strict-render-check

新增：

```text
goalcli lock write
goalcli strict-render-check
```

验收：

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
test ! -e "$tmp/redisx/xlib-standard.lock"
test ! -e "$tmp/redisx/.agent"
test ! -e "$tmp/redisx/contracts"
```

### Commit 9：docs rewrite

修改：

```text
README.md
docs/generation.md
docs/standard/*.md
docs/downstream-matrix.md
docs/downstream-sync-policy.md
docs/release.md
```

验收：

```bash
GOWORK=off make docs-check
GOWORK=off make pathguard
```

### Commit 10：GitHub workflow thin adapter

修改：

```text
.github/workflows/*.yml
```

验收：

```bash
GOWORK=off make pathguard
```

### Commit 11：integration golden strict render

渲染：

```text
kernel        -> l0.kernel
configx       -> l1.config
redisx        -> l2.redis
domain-market -> l25.domain-model
```

验收：

```bash
GOWORK=off make integration
```

### Commit 12：release gate

最终验收：

```bash
GOWORK=off make ci
GOWORK=off make release-check
```

---

## 22. 1 天行动计划

目标：完成破坏性骨架和 pathguard。

```text
1. 创建 .config/xlib/xlib.json
2. 创建 .config/xlib/rules/pathguard.json
3. 创建 .config/xlib/profiles/registry.json
4. 创建 l0.kernel / l1.config / l2.redis / l25.domain-model profile
5. 创建 .config/xlib/harness/gates.json
6. 创建 .config/xlib/capabilities/registry.json
7. 创建 .config/xlib/schemas/*.schema.json
8. 实现 goalcli config check
9. 实现 goalcli pathguard
10. 修改 Makefile 增加 pathguard/config-check
11. 禁止 --enable-governance / --layer
12. 删除 legacy paths
```

当天验收：

```bash
test ! -e .agent
test ! -e .xlib
test ! -e contracts
test ! -e xlib-standard.lock
test ! -e .githooks

GOWORK=off make pathguard
GOWORK=off make config-check
```

---

## 23. 7 天行动计划

### Day 1：strict config root skeleton

交付：

```text
.config/xlib skeleton
pathguard
config-check
legacy delete
```

### Day 2：profile / downstream / lock packages

交付：

```text
profile-check
downstream-check
lock write
```

### Day 3：render_template breaking rewrite

交付：

```text
--config / --profile only
lock-in-.config-only
strict-render-check
```

### Day 4：Makefile / hooks / workflow adapter

交付：

```text
Makefile adapter-only
.config/git/hooks
.github thin workflow
```

### Day 5：contracts / boundary / docs config-driven

交付：

```text
contracts --config
boundary --config
docs-check --config
```

### Day 6：integration golden

交付：

```text
kernel
configx
redisx
domain-market
```

### Day 7：release gate 收口

交付：

```text
ci
integration
release-check
strict clean report
```

第 7 天验收：

```bash
GOWORK=off make pathguard
GOWORK=off make strict-check
GOWORK=off make config-check
GOWORK=off make profile-check
GOWORK=off make downstream-check
GOWORK=off make contracts
GOWORK=off make boundary
GOWORK=off make docs-check
GOWORK=off make ci
GOWORK=off make integration
GOWORK=off make release-check
```

---

## 24. 30 天行动计划

### Week 1：strict `.config` root landing

结果：

```text
xlib-standard repo 无 legacy path
render_template 只接受 --config + --profile
lock 只在 .config/xlib
pathguard 进入 ci/release-check 第一关
```

### Week 2：下游强制重生成

目标下游：

```text
kernel
configx
observex
testkitx
redisx
postgresx
kafkax
natsx
taosx
ossx
clickhousex
decimalx
domain-market
domain-exchange
domain-macro
```

每个下游必须满足：

```text
.config/xlib/xlib-standard.lock.json exists
no root xlib-standard.lock
no .agent
no contracts
adoption-check passed
release-check passed
```

### Week 3：drift-check / upgrade-plan

新增：

```bash
goalcli drift check --lock .config/xlib/xlib-standard.lock.json
goalcli upgrade plan --repo ../redisx
```

输出：

```text
profile drift
rules drift
schema drift
harness drift
template drift
adoption proof stale
```

### Week 4：architecture harness + fragment renderer

交付：

```text
architecture check
profile fragment plan
strict fragment renderer
generated docs projection
```

---

## 25. 验收矩阵

| 维度 | 必须通过 |
|---|---|
| 旧路径 | `.agent`、`.xlib`、`contracts`、`xlib-standard.lock`、`.githooks` 不存在 |
| 旧参数 | `--enable-governance`、`--layer` 在当前用法中不存在 |
| 配置入口 | 所有 xlib 命令显式读 `.config/xlib/xlib.json` |
| lock | 只生成 `.config/xlib/xlib-standard.lock.json` |
| hooks | `core.hooksPath` 只允许 `.config/git/hooks` |
| contracts | 只读 `.config/xlib/schemas` 和 `.config/xlib/contracts` |
| harness | gate 定义只在 `.config/xlib/harness` |
| agent runtime | 只在 `.config/xlib/agent` |
| downstream | registry 只在 `.config/xlib/downstream` |
| profile | generator 只由 `--profile` 驱动 |
| docs | docs 只做人读投影，不做事实源 |
| release | release 只做 Evidence，不做 config source |
| CI | workflow 只做 adapter |
| Makefile | Makefile 只做 adapter |
| pathguard | `ci` 和 `release-check` 第一关 |
| integration | golden downstream 无 legacy path |

---

## 26. 最终执行命令

本仓库：

```bash
GOWORK=off make pathguard
GOWORK=off make strict-check
GOWORK=off make config-check
GOWORK=off make profile-check
GOWORK=off make downstream-check
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

渲染验证：

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

## 27. 衡量指标

| 指标 | 目标 |
|---|---|
| Legacy path count | 0 |
| Legacy CLI flag count | 0 |
| Root lock count | 0 |
| `.config` coverage | 100% xlib-controlled machine facts |
| Profile coverage | Foundation + L2.5 + governance-only |
| Downstream regenerated coverage | 30 天内 100% |
| Adoption proof freshness | 不超过 7 天 |
| Drift age | 不超过 7 天 |
| Boundary P0 violation | 0 |
| x.go reverse dependency | 0 |
| L2.5 provider/strategy/execution violation | 0 |
| Docs generated-block drift | 0 |
| Release evidence completeness | 100% |
| Golden render success rate | 100% |
| CI false positive rate | < 5% |

---

## 28. 迭代优化机制

每次失败都必须进入分类，并决定是否沉淀为规则：

```text
template_defect
profile_defect
downstream_misuse
contract_breaking_change
documentation_drift
evidence_gap
architecture_violation
pathguard_gap
```

处理规则：

| 失败类型 | 处理 |
|---|---|
| `template_defect` | 修 fragment / render，并加 golden case |
| `profile_defect` | 修改 profile schema 与 profile-check |
| `downstream_misuse` | 加 adoption-check / boundary rule |
| `contract_breaking_change` | 要求 major 或 compatibility adapter |
| `documentation_drift` | docs generated block check |
| `evidence_gap` | proof schema + checksum |
| `architecture_violation` | architecture rule + pathguard |
| `pathguard_gap` | 增加 forbidden path / pattern |

---

## 29. AI / 自动化介入位置

AI 不应该优先生成业务代码，而应该用在这些高杠杆点：

```text
1. legacy path 自动扫描和修复建议
2. profile 推导和差异分析
3. 下游 lock drift 归因
4. architecture rule violation 解释
5. docs projection 生成
6. release note / breaking migration guide 生成
7. adoption proof 审查
8. CI failure 分类
9. upgrade plan 自动生成
10. PR review checklist 自动生成
```

---

## 30. 最终推荐路径

最终执行路径：

```text
P0: 发布 v0.6.0 strict-config-root breaking release
P1: 删除所有 legacy path
P2: .config/xlib/xlib.json 成为唯一入口
P3: profiles / schemas / contracts / rules / harness / agent / downstream 全部进入 .config
P4: render_template.sh 只接受 --config + --profile
P5: lock 只写 .config/xlib/xlib-standard.lock.json
P6: pathguard 阻止旧路径和旧参数回归
P7: Makefile / .github 只做 adapter
P8: 下游全部重新生成或强制迁移
P9: release gate 必须证明 strict clean
```

最终判断：

> **`.config/` 不是新目录，而是新的权力中心。旧路径不是历史遗留，而是违规状态。**

这是最彻底、最干净、最不容易反复的根治方案。


---

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
