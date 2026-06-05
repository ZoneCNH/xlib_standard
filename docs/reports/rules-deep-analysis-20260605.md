# .agent/rules/ 深度分析报告（事实核对修订版）

> 生成日期: 2026-06-05
> 复核日期: 2026-06-05
> 分析范围: `.agent/rules/` 目录下 18 个文件：16 个 Markdown 文件与 2 个 YAML 文件，共 8687 行

分析范围：`.agent/rules/` 当前版本、与规则生成/校验强相关的 `README.md`、`.agent/registries/generated-artifacts.yaml`、`scripts/verify_rules.py`。

本报告是在隔离 worktree `codex/rules-report-fix-20260605T062208Z` 中经 Agent Team 初修后做的 leader 终审版。目标不是重写规则，而是核对原报告结论，删除错误判断，保留能被当前仓库证据支撑的问题和修复顺序。

| 文件 | 行数 | 类型 | 规则数 / 用途 | 域 |
| --- | ---: | --- | --- | --- |
| README.md | 77 | 索引 | — | 导航/权威顺序 |
| iron-rules.md | 49 | **SSOT** | 7 (压缩) | 铁律/标准退出码 |
| registry.yaml | 3798 | 机器索引 | 419 | 规则机器化索引（P0=119, P1=300） |
| enforcement-normalization.yaml | 71 | 归一化配置 | — | `enforced_by` 命令归一化 |
| core-rules.md | 1008 | 机器渲染 | 49 | 核心控制层 |
| schema-registry-rules.md | 1118 | 机器渲染 | 60 | Schema/Registry/Goal Pack |
| agent-runtime-rules.md | 1640 | 机器渲染 | 75 | Agent 运行时/工具链/度量 |
| goal-rules.md | 452 | **手写** | ~40 | Goal 对象模型全生命周期 |
| worktree-rules.md | 45 | 手写 | 5 | Worktree |
| commit-rules.md | 43 | 手写 | 4 | Commit |
| pr-rules.md | 43 | 手写 | 4 | PR |
| evidence-rules.md | 49 | 手写 | 3 | Evidence |
| release-rules.md | 60 | 手写 | 3 | Release |
| harness-rules.md | 57 | 手写 | 4 | Harness |
| security-rules.md | 39 | 手写 | 3 | Security |
| issue-rules.md | 36 | 手写 | 3 | Issue |
| risk-decision-rules.md | 51 | 手写 | 3 | Risk/Decision/Rollback |
| self-improving-rules.md | 51 | 手写 | 3 | Retrospective/Self-improving |

**规模分布**: 3 个机器渲染 Markdown 文件占 3766 行；`registry.yaml` 占 3798 行；`enforcement-normalization.yaml` 占 71 行；README、SSOT 与手写域文件合计 1052 行。原报告“16 个规则文件 (4818 行)”只覆盖了 Markdown 视角且遗漏了两个当前存在的 YAML 文件。

- “`.agent/rules/registry.yaml` 缺失 / 治理索引不存在”是错误结论。当前目录存在 `.agent/rules/registry.yaml`，且 `python3 scripts/verify_rules.py` 可读取并通过校验。
- “需要恢复 `00-index.md`、`01-core-rules.md` 等数字前缀文件名”是错误方向。当前 README 已声明真实树为 `README.md`、`registry.yaml`、`core-rules.md`、`schema-registry-rules.md`、`agent-runtime-rules.md` 等非数字前缀文件。

当前真实问题按风险排序如下：

### 2.1 [HIGH] 退出码定义冲突

两处定义仍不一致，但优先级应从 CRITICAL/P0 调整为 HIGH/P1：`README.md` 已声明 `iron-rules.md` 与 `registry.yaml` 为优先依据，冲突主要会误导阅读 `agent-runtime-rules.md` 的实现者，而不是证明当前 gate 已失效。

**iron-rules.md（权威标准退出码）**:

| 退出码 | 含义 |
| --- | --- |
| 5 | worktree / main 违规 |
| 6 | schema 校验失败 |
| 7 | secret / 凭据泄漏 |
| 8 | Evidence 缺失或伪造 |
| 9 | Traceability 断链 |
| 10 | Release 不完整 |

**agent-runtime-rules.md §84** (`RULE-GOALCLI-EXIT-001`):

| 退出码 | 含义 |
| --- | --- |
| 2 | POLICY_VIOLATION |
| 3 | SCHEMA_INVALID |
| 4 | EVIDENCE_MISSING |
| 5 | TRACEABILITY_BROKEN |
| 6 | WORKTREE_INVALID |
| 7 | SECRET_DETECTED |
| 8 | RELEASE_BLOCKED |
| 9 | NEEDS_HUMAN_APPROVAL |
| 10 | INCONSISTENT_STATE |

**冲突**: 退出码 5-10 在两处含义不同。`iron-rules.md` 是标准退出码 SSOT，`agent-runtime-rules.md` 是机器渲染文件；修复路径应是调整源数据或 `scripts/render_domain_rules.py` 后重新渲染，而不是手改生成文件。

## 2. 真实问题

### 2.2 [HIGH] `goal-rules.md` 与其他文件存在概念重叠

原报告把问题描述为“逐字重复约 70 行”和“约 60% 内容重复”，当前证据不足以支撑该精确量化；但维护风险仍然真实存在。

#### 重叠 1: `goal-rules.md` 与 `iron-rules.md` / `core-rules.md`

`goal-rules.md` 的“第一性原理铁律”、状态机、ID 规则、Context Recovery 与对象模型叙述，与 `iron-rules.md` 和 `core-rules.md` 中的权威规则存在明显概念重叠。

**影响**: 修改一处时需要同步多处叙述，否则会形成 SSOT 漂移。此风险应保留为 P1 维护性问题，但不应在未重新量化前继续声称“270/452 行、60%”。

#### 重叠 2: `goal-rules.md` 与其他手写域文件

`goal-rules.md` 覆盖 worktree、commit、PR、evidence、release、harness、security 等域的摘要。摘要本身有导航价值，但应明确“导航摘要”与“权威规则”边界，避免读者把摘要当作各域 SSOT。

影响：

### 2.3 [已纠正 / INFO] `registry.yaml` 当前存在，原“缺失”结论失效

原报告称 `.agent/rules/registry.yaml` 缺失，这是错误结论。当前目录中存在 `.agent/rules/registry.yaml`，且 `wc -l` 显示 3798 行；`rg` 验证的关键字段为：

- `generated_at: 2026-06-03`
- `total_rules: 419`
- `p0_count: 119`
- `p1_count: 300`
- `active_count: 363`
- `indexed_count: 56`

**影响修正**: README 中“419 条规则机器化索引”“P0=119, P1=300”“363/419 active”等覆盖率数据可由 `registry.yaml` 验证，不再是 P0 缺失问题。

**剩余风险**: `registry.yaml` 是大型生成索引，应继续用 `python3 scripts/verify_rules.py`、`make rules-verify` 或 `go run ./cmd/goalcli rules-consistency-check` 验证新鲜度与一致性。

影响：

### 2.4 [MEDIUM] 目录命名建议与实际文件名存在生成规则层面的差异

原报告把数字前缀命名归因于 “README §113”，该引用不准确。当前 `README.md` 的目录树与实际文件名一致，并未要求 `00-index.md`、`01-core-rules.md` 这类数字前缀。

真实可保留的问题是：`agent-runtime-rules.md` 中的 `RULE-REPO-LAYOUT-001` / repo layout 建议仍提到数字前缀布局，与当前 `.agent/rules/` 实际命名不同。

**影响修正**: 不应把“统一命名体系”列为当前 P2 批量重命名动作；更安全的动作是记录该差异为已知偏离，或在源规则/渲染产物中更新推荐布局。批量重命名会影响大量文档、脚本与引用，不适合作为本报告的直接优化建议。

修复建议：

- 在 README 中把“人改哪里 / 机器生成哪里 / 如何重建 / 如何验证”写成明确流程。
- 对 generated Markdown 增加更醒目的头部约束：不要手改，改源数据后重新生成。

三个机器渲染 Markdown 文件合计 3766 行：

`commit-rules.md`、`pr-rules.md`、`evidence-rules.md`、`release-rules.md`、`harness-rules.md`、`security-rules.md`、`issue-rules.md`、`risk-decision-rules.md`、`self-improving-rules.md`、`worktree-rules.md` 等文件短小，部分与 registry/铁律存在重复表述。

`agent-runtime-rules.md` 覆盖 Agent 协议、并发、命令契约、Bootstrap/Doctor、Dashboard、度量、Lease/Heartbeat 等多个子域，维护和 review 成本偏高。由于这些是机器渲染文件，拆分应通过源数据和 `scripts/render_domain_rules.py` 完成，不能直接手改产物。

## 3. 对原报告建议的修正

### 2.6 [MEDIUM] 源引用是考古线索，不应作为当前权威依据

机器渲染文件标注 `source: §N Lxxxx`，指向 `.worktree/goal-patch.md` 的历史行号。README 已说明 `.worktree/goal-patch.md` 是历史推导材料，“仅供考古，不可作为依据”。

**影响修正**: 这不是“源引用必然失效”的直接证据，而是新鲜度与可追溯性风险：报告或实现不能只依赖这些历史行号，应优先引用 `iron-rules.md`、`registry.yaml`、当前域文件和验证命令输出。

- 保留 exit code 一致性修复，但优先级调整为 P1。
- 保留路径漂移修复，但定位为 P2。
- 保留维护入口清晰化，但应围绕 generated/source contract，而不是围绕文件名重排。

### 2.7 [LOW] 手写文件较薄，导航成本偏高

11 个手写域文件中，多数只有 3-5 条规则、36-60 行：

## 5. 验证证据

这些文件单独存在增加了导航成本，但也保持了“一个域一个文件”的边界。当前建议仅限于改进 README 导航与交叉引用，不建议合并文件。

```bash
rg --files .agent/rules
wc -l .agent/rules/*
python3 scripts/verify_rules.py
rg -n "generated_by|scripts/extract_rules.py|scripts/render_domain_rules.py|rules-verify|07-worktree-rules|RULE-GOALCLI-EXIT-001|exit code|exit_code" .agent/rules .agent/registries/generated-artifacts.yaml README.md scripts/verify_rules.py
```

### 2.8 [LOW] 规则间交叉引用仍可增强

手写文件之间缺少显式“相关规则”导航。例如：

- `evidence-rules.md` 可指向 `release-rules.md` 与 `pr-rules.md`
- `pr-rules.md` 可指向 `evidence-rules.md` 与 `harness-rules.md`
- `security-rules.md` 可指向 `harness-rules.md` 中的 secret-check gate

这属于可发现性增强，不影响当前规则有效性。

---

## 3. 优化建议

### 3.1 [P1] 修复退出码冲突

**动作**: 统一到 `iron-rules.md` 定义的标准退出码。

**方法**:

1. 检查 `RULE-GOALCLI-EXIT-001` 的源数据与 `scripts/render_domain_rules.py` 中的渲染逻辑
2. 修正源数据或渲染脚本，使 `agent-runtime-rules.md` 的映射与 `iron-rules.md` 一致
3. 重跑渲染生成，并运行 `python3 scripts/verify_rules.py` 与 `go run ./cmd/goalcli rules-consistency-check`

**预期**: 消除实现者阅读生成文档时的退出码歧义。

---

### 3.2 [验证项] 保持 `registry.yaml` 新鲜度验证

**动作**: 保留验证动作，不再作为“缺失文件”修复项。

**方法**:

1. 使用 `rg -n 'total_rules|p0_count|p1_count|active_count|indexed_count' .agent/rules/registry.yaml` 确认当前统计
2. 使用 `python3 scripts/verify_rules.py` 或 `make rules-verify` 验证 active 规则的 `enforced_by` 命令
3. 在规则变更后重跑 `scripts/extract_rules.py` / `scripts/render_domain_rules.py`，再提交生成产物

---

### 3.3 [P1] 收敛 `goal-rules.md` 的摘要边界

**动作**: 重构 `goal-rules.md`，将重复叙述降为导航引用。

**方案 A（推荐）**: 保留 `goal-rules.md` 作为“Goal 全生命周期导航文档”，但：

- 将“第一性原理铁律”改为引用 `iron-rules.md`
- 将状态机/ID/Context/Object Model 改为引用 `core-rules.md`
- 保留独有的 Spec/Design/Task/AutoResearch/Change Propagation 规则
- 保留独有的 `xlib-standard/kernel/x.go` 专用规则
- 保留独有的 `goalcli` 最小命令集和 Makefile gate 最小集合

**方案 B**: 将 `goal-rules.md` 明确标记为导航/聚合文件，所有权威细节跳转到各域文件。

**预期**: 降低 SSOT 漂移风险，同时保留全生命周期视角。

---

### 3.4 [P1] 拆分 `agent-runtime-rules.md`

**动作**: 将 1640 行的大文件拆分为 2-3 个子文件。

**建议拆分**:

- `agent-protocol-rules.md` — Agent 执行协议/权限边界/并发/Lease/Heartbeat
- `goalcli-arch-rules.md` — `goalcli` 命令契约/退出码/架构/Checker/事务/Dry-run
- `governance-metrics-rules.md` — Dashboard/度量/治理节奏/Bootstrap/Doctor/Repair

**前提**: 需要调整 `scripts/render_domain_rules.py` 的分区逻辑并重新生成 `registry.yaml`。

---

### 3.5 [P2] 记录或修正规则布局建议差异

**动作**: 不执行批量重命名；先处理 `RULE-REPO-LAYOUT-001` 与当前文件名的差异。

**可选路径**:

1. 如果当前命名是有意设计，则更新源规则/渲染产物，移除数字前缀建议
2. 如果数字前缀仍是目标形态，则单独开迁移计划，覆盖 README、脚本、链接和 `registry.yaml` 再生成

**当前结论**: README 与实际文件名一致；“README 推荐数字前缀”是原报告的错误引用。

---

### 3.6 [P2] 增加交叉引用

**动作**: 在手写文件间添加交叉引用。

**示例**:

- `evidence-rules.md` 末尾添加: `> 相关: `.agent/rules/release-rules.md` · `.agent/rules/pr-rules.md``
- `pr-rules.md` 末尾添加: `> 相关: `.agent/rules/evidence-rules.md` · `.agent/rules/harness-rules.md``

---

## 4. 优先级排序

| 优先级 | 动作 | 影响 | 工作量 |
| --- | --- | --- | --- |
| P1 | 修复退出码冲突 | 消除生成文档与 SSOT 的实现歧义 | 低-中（源数据/渲染脚本 + 再生成） |
| P1 | 收敛 `goal-rules.md` 重叠叙述 | 降低 SSOT 漂移风险 | 中（重构文档边界） |
| P1 | 拆分 `agent-runtime-rules.md` | 提升机器渲染产物可维护性 | 中（调整渲染分区） |
| P2 | 记录或修正规则布局建议差异 | 避免 README/生成规则读者困惑 | 低-中 |
| P2 | 增加交叉引用 | 提升可发现性 | 低 |
| 验证项 | 持续验证 `registry.yaml` | 证明覆盖率与 `enforced_by` 一致 | 低 |

---

## 5. 结论

`.agent/rules/` 体系当前包含 18 个文件、8687 行，其中 `registry.yaml` 已存在并声明 419 条规则（P0=119, P1=300，active=363，indexed=56）。原报告关于 “registry.yaml 缺失” 和 “README 推荐数字前缀命名” 的结论需要删除或改为已纠正事实。

当前真实问题集中在：

1. **一致性**: `iron-rules.md` 与 `agent-runtime-rules.md` 的退出码映射不一致；`goal-rules.md` 与 SSOT/域文件存在概念重叠。
2. **可维护性**: 机器渲染文件过大；历史 `source: §N Lxxxx` 引用只适合作为考古线索，不能作为当前权威依据。
3. **可发现性**: 手写域文件较薄且交叉引用不足；`RULE-REPO-LAYOUT-001` 的布局建议与当前 README/文件名需要对齐。

建议按 P1 → P2 → 验证项持续化的顺序处理。当前不应再投入 P0 精力去“找回 registry.yaml”或批量重命名文件。

---

## 6. 本次复核证据

- `find .agent/rules -maxdepth 1 -type f | sort`：确认当前 18 个文件，包括 `.agent/rules/registry.yaml` 与 `.agent/rules/enforcement-normalization.yaml`。
- `wc -l .agent/rules/*`：确认 `.agent/rules/` 合计 8687 行；`registry.yaml` 3798 行；`agent-runtime-rules.md` 1640 行；`schema-registry-rules.md` 1118 行；`core-rules.md` 1008 行。
- `rg -n 'total_rules|p0_count|p1_count|active_count|indexed_count|generated_at' .agent/rules/registry.yaml`：确认 `generated_at: 2026-06-03`、`total_rules: 419`、`p0_count: 119`、`p1_count: 300`、`active_count: 363`、`indexed_count: 56`。
- `rg -n '标准退出码|RULE-GOALCLI-EXIT-001|POLICY_VIOLATION|WORKTREE_INVALID|TRACEABILITY_BROKEN' .agent/rules/iron-rules.md .agent/rules/agent-runtime-rules.md`：确认退出码冲突仍存在。
- `python3 scripts/verify_rules.py`：用于验证 active 规则的 `enforced_by` 命令；本报告修正后应继续作为规则变更 gate。

---

_报告结束_
