# xlib-standard 项目深度分析报告

> 生成日期：2026-06-02
> 分析范围：全仓库代码、测试、CI、治理体系
> 对标文档：docs/goal.md v2.9.3 Complete

---

## 一、项目定位与核心身份

`xlib-standard` 是一个**基础库标准与交付运行时仓库**，承担五类职责：

| 职责 | 说明 |
|------|------|
| **Standard Source** | 维护基础库 P0 标准、仓库角色、分层、模块边界、DoD、安全、release 和 Evidence 协议 |
| **Go Reference Template** | 提供可编译参考包 `pkg/templatex`，证明标准可落地 |
| **Generator** | 通过 `scripts/render_template.sh` 渲染具体基础库 |
| **Harness** | 通过 Makefile、scripts、CI 固化 required、extended、score 和 final gate |
| **Evidence Runtime** | 记录可追溯完成状态，生成 release manifest |

**旧名 `baselib-template` / `foundationx` 已迁移**，新默认下游集成目标是 `kernel`。

---

## 二、技术栈与依赖

| 维度 | 现状 |
|------|------|
| 语言 | Go 1.23 |
| 模块路径 | `github.com/ZoneCNH/xlib-standard` |
| 外部依赖 | **零**（`go.sum` 仅含标准库间接依赖） |
| Linter | golangci-lint v2.1.6，启用 17 个 linter |
| 安全扫描 | govulncheck v1.3.0 + 自研 secret scanner |
| CI | GitHub Actions，所有 Action 固定 40 位 SHA |

**亮点：零外部依赖 = 极低供应链风险。**

---

## 三、代码结构分析

```
总 Go 代码行数：~68,900 行（含 .worktree 历史副本）
实际源码（不含 worktree）：约 3,000-4,000 行
```

### 核心包

| 包 | 职责 | 覆盖率 |
|---|------|--------|
| `pkg/templatex` | 公共 API 参考实现（Client、Config、Health、Metrics、Errors） | **100%** |
| `internal/sanitize` | 密钥脱敏（`Secret()` → `***`） | **100%** |
| `internal/validation` | 输入校验（`RequireNonEmpty`） | **100%** |
| `internal/releasequality` | Release 质量评分（10 维度打分） | **100%** |
| `internal/tools/releasemanifest` | Release manifest 生成工具 | **100%** |
| `cmd/xlibgate` | Gate 路由入口（14 个子命令） | **100%** |
| `testkit` | 可复用测试夹具、断言、golden 文件工具 | **100%** |

### 测试质量

- **全部 12 个包测试通过**
- **所有可测包覆盖率 100%**
- 包含：单元测试、属性测试、fuzz smoke、golden 测试、fixture 测试

---

## 四、Harness 体系（Makefile + CI）

### Gate 分层

| 层级 | Target | 包含检查 |
|------|--------|----------|
| **基础** | `make ci` | fmt → vet → lint → test → race → boundary → security → contracts → score |
| **扩展** | `make ci-extended` | ci + property + golden + fuzz-smoke |
| **发布** | `make release-check` | ci + integration + dependency-check + standard-impact-check + docs-check + score-check + evidence + hash + verify |
| **最终** | `make release-final-check` | release-check + min score 9.5 + clean workspace |
| **预发布** | `make release-preflight VERSION=vX.Y.Z` | 版本检查 + main 同步 + tag 检查 + CHANGELOG + 工具 + final check |

### CI Workflow（4 个）

| Workflow | 触发 | 职责 |
|----------|------|------|
| `ci.yml` | PR + push main | 完整 release-check + score + evidence upload |
| `integration.yml` | PR + push main | 模板渲染集成测试 |
| `release.yml` | tag push | 发布验证 |
| `security.yml` | 定时 + push | 安全扫描 |

---

## 五、Release Scorecard 评分体系

`internal/releasequality` 实现了 **10 维度评分**，每维度权重 1.0：

| 维度 | 检查内容 |
|------|----------|
| scorecard_doc | `docs/scorecard.md` 存在 |
| manifest_score_schema | manifest 包含 `score`/`workflow_run_id`/`artifact_url` |
| score_cli | xlibgate score 命令可运行 |
| score_gate | Makefile 中有 score-check |
| manifest_min_score_verify | release evidence 验证脚本检查分数阈值 |
| security_gate | secret scanner 覆盖 token 和私钥 |
| release_docs | 发布文档绑定 score 和 CI artifact |
| supply_chain_docs | 供应链文档包含 score/workflow evidence |
| retrospective_template | 复盘模板包含 Score/Gate/Patch |
| release_template | 发布模板要求 score 和 artifact evidence |

**当前评分：10.0/10.0（满分）**

---

## 六、治理与标准体系

### `.agent/` 目录（Goal Runtime v3.1）

包含 28 个治理文件，覆盖：
- goal-runtime、state-machine、object-model
- evidence-protocol、evidence-template
- risk-register、decision-log、rollback-protocol
- harness、gates、agent-teams
- retrospective-template、release-template
- spec、plan、design、traceability-matrix

### `docs/standard/` 目录

12 个标准文档，覆盖：
- 仓库角色、分层、模块边界、DoD
- Harness gates、Evidence 协议、Release 标准
- 安全与密钥策略、下游兼容性、模板生成契约
- 复盘与补丁、x.go 集成边界

---

## 七、与 Goal v2.9.3 的差距分析

Goal v2.9.3 定义了 P0/P1/P2 共 **52 个 Issue**，当前实现状态：

### P0 Minimal Kernel（16 个 Issue）

| Issue | 标题 | 状态 |
|-------|------|------|
| P0-001 | Minimal Constitution | ⚠️ 有 CONSTITUTION 相关文档，但无独立 `CONSTITUTION.md` |
| P0-002 | Minimal Kernel Policy | ❌ 无 `.agent/minimal-kernel.yaml` |
| P0-003 | xlibgate CLI Skeleton | ✅ 已实现（14 个子命令） |
| P0-004 | main-guard | ❌ 无 `main_guard.go`，无上下文感知阻断 |
| P0-005 | worktree-guard | ❌ 无 `worktree_guard.go` |
| P0-006 | evidence-check | ⚠️ 有 evidence 生成，但无 DONE 协议解析 |
| P0-007 | boundary no-xgo-import | ✅ 已实现（`scripts/check_boundary.sh`） |
| P0-008 | no-secret-default | ✅ 已实现（`scripts/check_secrets.sh`） |
| P0-009 | Makefile governance-check | ⚠️ 有 `make ci` 但无独立 `governance-check` target |
| P0-010 | CI Required Checks Skeleton | ✅ 已实现（`ci.yml`） |
| P0-011 | Release Manifest Skeleton | ✅ 已实现 |
| P0-012 | DONE with evidence Protocol | ⚠️ 有文档但无 CLI 解析 |
| P0-013 | Execution Context Policy | ❌ 无 `.agent/execution-context.yaml` |
| P0-014 | xlibgate CLI Contract | ⚠️ 有 CLI 但无 report schema |
| P0-015 | Issue/Command Registry | ❌ 无 `.agent/issue-registry.yaml` |
| P0-016 | Makefile Baseline | ⚠️ 有 Makefile 但无 baseline YAML |

**P0 完成度：约 5/16（31%）**

### P1 Governance Hardening（21 个 Issue）

几乎全部未实现。仅部分有文档基础。

**P1 完成度：约 1/21（5%）**

### P2 Runtime & Conformance（15 个 Issue）

全部未实现。

**P2 完成度：0/15（0%）**

---

## 八、架构优势与风险

### 优势

1. **零外部依赖**：供应链风险极低
2. **100% 测试覆盖率**：所有可测包全覆盖
3. **多层 Gate 体系**：ci → ci-extended → release-check → release-final-check
4. **Evidence 可追溯**：manifest + SHA256 + CI artifact
5. **Action SHA pinning**：所有 GitHub Actions 固定 40 位 SHA
6. **Fuzz + Property + Golden 测试**：超越基础单元测试

### 风险

| 风险 | 等级 | 说明 |
|------|------|------|
| P0 治理 Gate 缺失 | **高** | 无 main-guard/worktree-guard/execution-context，Agent 可在 main 上直接开发 |
| `.agent/` 文档与 YAML 混用 | 中 | 当前 `.agent/` 全是 `.md`，Goal 要求大量 `.yaml` |
| xlibgate 命令与 Goal 不对齐 | 高 | 当前 14 个命令 vs Goal 要求 30+ 个命令 |
| `internal/` 包过少 | 中 | 仅有 sanitize/validation/releasequality，Goal 要求大量 guards/evidence/boundary/scope 包 |
| `.worktree/` 残留数据 | 低 | 有历史 team worktree 数据残留，已 gitignore |

---

## 九、总结

| 维度 | 评分 | 说明 |
|------|------|------|
| 代码质量 | ⭐⭐⭐⭐⭐ | 100% 覆盖率，零依赖，多层测试 |
| CI/Harness | ⭐⭐⭐⭐⭐ | 完整 gate 链，Evidence 可追溯 |
| 标准文档 | ⭐⭐⭐⭐ | 文档丰富，但部分与 Goal YAML 要求不一致 |
| P0 治理实现 | ⭐⭐ | 核心 guard/context/registry 未实现 |
| P1/P2 实现 | ⭐ | 几乎未开始 |

**一句话结论：代码质量和 Harness 体系已达生产级，但 Goal v2.9.3 定义的治理运行时（guards、context、registry、conformance）仍需从零构建。项目处于"基础设施就绪，治理内核待建"的状态。**
