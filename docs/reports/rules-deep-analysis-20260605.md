# `.agent/rules/` 深度分析报告

> 生成日期: 2026-06-05
> 分析范围: `.agent/rules/` 目录下全部 16 个规则文件 (4818 行)

---

## 1. 文件清单与规模

| 文件                     | 行数 | 类型     | 规则数   | 域                           |
| ------------------------ | ---- | -------- | -------- | ---------------------------- |
| README.md                | 77   | 索引     | —        | 导航/权威顺序                |
| iron-rules.md            | 49   | **SSOT** | 7 (压缩) | 铁律/退出码                  |
| core-rules.md            | 1008 | 机器渲染 | 49       | 核心控制层                   |
| schema-registry-rules.md | 1118 | 机器渲染 | 60       | Schema/Registry/Goal Pack    |
| agent-runtime-rules.md   | 1640 | 机器渲染 | 75       | Agent 运行时/工具链/度量     |
| goal-rules.md            | 452  | **手写** | ~40      | Goal 对象模型全生命周期      |
| worktree-rules.md        | 45   | 手写     | 5        | Worktree                     |
| commit-rules.md          | 43   | 手写     | 4        | Commit                       |
| pr-rules.md              | 43   | 手写     | 4        | PR                           |
| evidence-rules.md        | 49   | 手写     | 3        | Evidence                     |
| release-rules.md         | 60   | 手写     | 3        | Release                      |
| harness-rules.md         | 57   | 手写     | 4        | Harness                      |
| security-rules.md        | 39   | 手写     | 3        | Security                     |
| issue-rules.md           | 36   | 手写     | 3        | Issue                        |
| risk-decision-rules.md   | 51   | 手写     | 3        | Risk/Decision/Rollback       |
| self-improving-rules.md  | 51   | 手写     | 3        | Retrospective/Self-improving |

**规模分布**: 3 个机器渲染文件占 3766 行 (78%)，11 个手写文件占 935 行 (19%)，README 占 77 行 (2%)。

---

## 2. 发现的问题

### 2.1 [CRITICAL] 退出码定义冲突

两处定义不一致：

**iron-rules.md (权威)**:

| 退出码 | 含义                 |
| ------ | -------------------- |
| 5      | worktree / main 违规 |
| 6      | schema 校验失败      |
| 7      | secret / 凭据泄漏    |
| 8      | Evidence 缺失或伪造  |
| 9      | Traceability 断链    |
| 10     | Release 不完整       |

**agent-runtime-rules.md §84** (RULE-GOALCLI-EXIT-001):

| 退出码 | 含义                 |
| ------ | -------------------- |
| 2      | POLICY_VIOLATION     |
| 3      | SCHEMA_INVALID       |
| 4      | EVIDENCE_MISSING     |
| 5      | TRACEABILITY_BROKEN  |
| 6      | WORKTREE_INVALID     |
| 7      | SECRET_DETECTED      |
| 8      | RELEASE_BLOCKED      |
| 9      | NEEDS_HUMAN_APPROVAL |
| 10     | INCONSISTENT_STATE   |

**冲突**: 退出码 5-10 含义完全不同。iron-rules.md 是权威源，但 agent-runtime-rules.md 是机器渲染文件，会误导实现者。

**建议**: 以 iron-rules.md 为准，修复 agent-runtime-rules.md 的 §84 退出码映射（需调整 `scripts/render_domain_rules.py` 或源数据）。

---

### 2.2 [HIGH] 大量内容重复

#### 重复 1: goal-rules.md 与 iron-rules.md / core-rules.md

`goal-rules.md` 的"第一性原理铁律"章节完整复制了 `iron-rules.md` 的七律内容 (RULE-CORE-001 ~ RULE-CORE-006)，逐字重复约 70 行。

此外 `goal-rules.md` 还重复了：

- 状态机 (与 core-rules.md §5 RULE-STATE-\* 重复)
- ID 规则 (与 core-rules.md §4 RULE-ID-\* 重复)
- Context Recovery (与 core-rules.md §6 RULE-CONTEXT-\* 重复)
- 对象模型 (与 core-rules.md §3 RULE-OBJECT-\* 重复)

**影响**: 修改一处时需同步多处，违反 SSOT 原则（iron-rules.md 自身 RULE-SSOT-002）。

#### 重复 2: goal-rules.md 与其他手写文件

`goal-rules.md` 包含了 worktree、commit、PR、evidence、release、harness、security 等所有手写文件的规则摘要，形成 1:N 重复。

#### 重复 3: core-rules.md 与 goal-rules.md 的 Spec/Design/Task 规则

`goal-rules.md` 的 Spec/Design/Task/AutoResearch/Change Propagation 章节是独有的，但在 `core-rules.md` 中有部分重叠（通过 registry.yaml 映射）。

**量化**: 估算 `goal-rules.md` 中约 60% 内容 (270/452 行) 与其他文件重复。

---

### 2.3 [HIGH] registry.yaml 缺失

README.md 多次引用 `.agent/rules/registry.yaml`：

- "419 条规则机器化索引"
- "P0=119, P1=300"
- "make rules-verify 强制断言"

但 `.agent/rules/` 目录中不存在 `registry.yaml`。需要确认：

1. 是否已迁移到其他位置？
2. 是否尚未生成？
3. 是否需要重新运行 `scripts/extract_rules.py`？

**影响**: README 中的覆盖率数据 (363/419 active, 87%) 无法验证。

---

### 2.4 [MEDIUM] 文件命名不一致

存在两套命名体系：

**机器渲染文件** (按域命名):

- `core-rules.md` — 核心控制层
- `schema-registry-rules.md` — Schema/Registry 层
- `agent-runtime-rules.md` — 运行时层

**手写文件** (按对象命名):

- `goal-rules.md` — Goal 对象
- `worktree-rules.md` — Worktree 对象
- `evidence-rules.md` — Evidence 对象
- ...

README 中推荐的目录结构 (§113) 使用数字前缀: `00-index.md`, `01-core-rules.md`, `07-worktree-rules.md` 等，与实际命名不符。

---

### 2.5 [MEDIUM] 机器渲染文件过大

三个机器渲染文件合计 3766 行：

- `agent-runtime-rules.md`: 1640 行 (75 条规则)
- `schema-registry-rules.md`: 1118 行 (60 条规则)
- `core-rules.md`: 1008 行 (49 条规则)

`agent-runtime-rules.md` 覆盖了 Agent 协议、并发、命令契约、Bootstrap/Doctor、Dashboard、度量、Lease/Heartbeat 等多个子域，建议拆分。

---

### 2.6 [MEDIUM] 源引用可能失效

机器渲染文件标注 `source: §N Lxxxx`，指向 `.worktree/goal-patch.md` 的行号。如果源文件已变更或不存在，这些行号引用将失去意义。

README 自身也注明: `.worktree/goal-patch.md` — 历史推导，**仅供考古，不可作为依据**。

---

### 2.7 [LOW] 手写文件过于稀薄

11 个手写文件中，多数只有 3-5 条规则、36-60 行：

- `issue-rules.md`: 36 行, 3 条规则
- `security-rules.md`: 39 行, 3 条规则
- `commit-rules.md`: 43 行, 4 条规则

这些文件单独存在增加了导航成本，但合并又可能破坏"一个域一个文件"的组织原则。

---

### 2.8 [LOW] 缺少规则间交叉引用

手写文件之间缺乏交叉引用。例如：

- `evidence-rules.md` 未引用 `release-rules.md` (Release 需要 Evidence)
- `pr-rules.md` 未引用 `evidence-rules.md` (PR 需要 Evidence)
- `security-rules.md` 未引用 `harness-rules.md` (secret-check 是 Harness Gate)

---

## 3. 优化建议

### 3.1 [P0] 修复退出码冲突

**动作**: 统一到 `iron-rules.md` 定义的标准退出码。

**方法**:

1. 检查 `scripts/render_domain_rules.py` 中 §84 的映射逻辑
2. 修正源数据或渲染脚本
3. 重跑渲染生成

**预期**: 消除退出码歧义，确保 goalcli 实现者不被误导。

---

### 3.2 [P0] 确认 registry.yaml 状态

**动作**: 确认 `registry.yaml` 的存在性和准确性。

**方法**:

1. 检查是否在其他路径 (如 `.agent/registries/`)
2. 如缺失，运行 `python3 scripts/extract_rules.py` 重新生成
3. 更新 README 中的覆盖率数据

---

### 3.3 [P1] 消除 goal-rules.md 重复

**动作**: 重构 `goal-rules.md`，去除与其他文件的重复内容。

**方案 A (推荐)**: 保留 `goal-rules.md` 作为"Goal 全生命周期导航文档"，但：

- 删除与 iron-rules.md 重复的"第一性原理铁律"章节，改为引用
- 删除与 core-rules.md 重复的状态机/ID/Context 章节，改为引用
- 保留独有的 Spec/Design/Task/AutoResearch/Change Propagation 规则
- 保留独有的 xlib-standard/kernel/x.go 专用规则
- 保留独有的 goalcli 最小命令集和 Makefile Gate 最小集合

**方案 B**: 将 `goal-rules.md` 标记为 deprecated，独有内容迁移到新的 `spec-design-task-rules.md`。

**预期**: 减少约 270 行重复内容，消除 SSOT 冲突风险。

---

### 3.4 [P1] 拆分 agent-runtime-rules.md

**动作**: 将 1640 行的大文件拆分为 2-3 个子文件。

**建议拆分**:

- `agent-protocol-rules.md` — Agent 执行协议/权限边界/并发/Lease/Heartbeat (~§57-§64, §100-§101, §216-§217)
- `goalcli-arch-rules.md` — goalcli 命令契约/退出码/架构/Checker/事务/Dry-run (~§83-§85, §120-§121, §175-§176)
- `governance-metrics-rules.md` — Dashboard/度量/治理节奏/Bootstrap/Doctor/Repair (~§77, §189-§190, §206-§207, §228-§229)

**前提**: 需要调整 `scripts/render_domain_rules.py` 的分区逻辑。

---

### 3.5 [P2] 统一命名体系

**动作**: 采用 README §113 推荐的数字前缀命名。

**映射**:

| 当前名称                 | 建议名称                    |
| ------------------------ | --------------------------- |
| core-rules.md            | 01-core-rules.md            |
| schema-registry-rules.md | 02-schema-registry-rules.md |
| agent-runtime-rules.md   | 03-agent-runtime-rules.md   |
| goal-rules.md            | 04-goal-rules.md            |
| worktree-rules.md        | 05-worktree-rules.md        |
| commit-rules.md          | 06-commit-rules.md          |
| pr-rules.md              | 07-pr-rules.md              |
| evidence-rules.md        | 08-evidence-rules.md        |
| release-rules.md         | 09-release-rules.md         |
| harness-rules.md         | 10-harness-rules.md         |
| security-rules.md        | 11-security-rules.md        |
| issue-rules.md           | 12-issue-rules.md           |
| risk-decision-rules.md   | 13-risk-decision-rules.md   |
| self-improving-rules.md  | 14-self-improving-rules.md  |

**注意**: 重命名会影响所有引用这些文件的文档和脚本，需要批量更新。

---

### 3.6 [P2] 增加交叉引用

**动作**: 在手写文件间添加交叉引用。

**示例**:

- `evidence-rules.md` 末尾添加: `> 相关: [Release 规则](release-rules.md) · [PR 规则](pr-rules.md)`
- `pr-rules.md` 末尾添加: `> 相关: [Evidence 规则](evidence-rules.md) · [Harness 规则](harness-rules.md)`

---

## 4. 优先级排序

| 优先级 | 动作                        | 影响           | 工作量                   |
| ------ | --------------------------- | -------------- | ------------------------ |
| P0     | 修复退出码冲突              | 消除实现歧义   | 低 (修改渲染脚本)        |
| P0     | 确认 registry.yaml          | 验证覆盖率数据 | 低 (检查/运行脚本)       |
| P1     | 消除 goal-rules.md 重复     | 消除 SSOT 违规 | 中 (重构文档)            |
| P1     | 拆分 agent-runtime-rules.md | 提升可维护性   | 中 (调整渲染脚本)        |
| P2     | 统一命名体系                | 提升导航体验   | 高 (批量重命名+更新引用) |
| P2     | 增加交叉引用                | 提升可发现性   | 低 (添加链接)            |

---

## 5. 结论

`.agent/rules/` 体系设计完备，419 条规则覆盖了 Goal Runtime 的完整生命周期。当前主要问题集中在：

1. **一致性**: 退出码冲突、SSOT 违规 (goal-rules.md 重复)
2. **可维护性**: 机器渲染文件过大、registry.yaml 缺失
3. **可发现性**: 命名不一致、缺乏交叉引用

建议按 P0 → P1 → P2 顺序逐步修复。P0 项可在 1 小时内完成，P1 项需 2-4 小时，P2 项可作为后续迭代。

---

_报告结束_
