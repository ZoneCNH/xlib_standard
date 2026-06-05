# .agent/rules/ 深度分析报告（事实核对修订版）

日期：2026-06-05

分析范围：`.agent/rules/` 当前版本、与规则生成/校验强相关的 `README.md`、`.agent/registries/generated-artifacts.yaml`、`scripts/verify_rules.py`。

本报告是在隔离 worktree `codex/rules-report-fix-20260605T062208Z` 中经 Agent Team 初修后做的 leader 终审版。目标不是重写规则，而是核对原报告结论，删除错误判断，保留能被当前仓库证据支撑的问题和修复顺序。

## 0. 终审结论

原报告中最严重的两类结论需要下调或删除：

- “`.agent/rules/registry.yaml` 缺失 / 治理索引不存在”是错误结论。当前目录存在 `.agent/rules/registry.yaml`，且 `python3 scripts/verify_rules.py` 可读取并通过校验。
- “需要恢复 `00-index.md`、`01-core-rules.md` 等数字前缀文件名”是错误方向。当前 README 已声明真实树为 `README.md`、`registry.yaml`、`core-rules.md`、`schema-registry-rules.md`、`agent-runtime-rules.md` 等非数字前缀文件。

当前真实问题按风险排序如下：

1. **[P1] 退出码语义在不同规则文档之间漂移**：`iron-rules.md` 作为铁律表，与生成出的规则条目对 exit code 2、5、6、8、9、10 的含义不一致；`RULE-GOALCLI-EXIT-001` 还额外定义了铁律表未定义的 3、4。
2. **[P1] 生成产物清单把规则生成器误标为 `goalcli rules-verify`**：当前 README 和生成文件头部指向 `scripts/extract_rules.py`、`scripts/render_domain_rules.py`，而 generated-artifacts 清单写成 `goalcli rules-verify`，这会误导后续治理自动化。
3. **[P2] 生成文档中残留历史路径 `07-worktree-rules.md`**：当前真实文件是 `.agent/rules/worktree-rules.md`，但若干生成文档仍引用旧命名。
4. **[P2] 规则文档体量集中在生成文件，人工维护面不应直接拆改**：应修源数据或生成器，而不是手工切开 generated Markdown。
5. **[P3] 小型人工规则文件存在局部重复和可合并空间**：可整理，但没有证据支持“60% 可删除”这类强结论。

## 1. 当前目录事实

`.agent/rules/` 当前共有 18 个文件：16 个 Markdown、2 个 YAML。

核心文件：

- `.agent/rules/README.md`
- `.agent/rules/registry.yaml`
- `.agent/rules/enforcement-normalization.yaml`
- `.agent/rules/iron-rules.md`
- `.agent/rules/core-rules.md`
- `.agent/rules/schema-registry-rules.md`
- `.agent/rules/agent-runtime-rules.md`
- `.agent/rules/goal-rules.md`
- `.agent/rules/worktree-rules.md`
- `.agent/rules/commit-rules.md`
- `.agent/rules/pr-rules.md`
- `.agent/rules/evidence-rules.md`
- `.agent/rules/release-rules.md`
- `.agent/rules/harness-rules.md`
- `.agent/rules/security-rules.md`
- `.agent/rules/issue-rules.md`
- `.agent/rules/risk-decision-rules.md`
- `.agent/rules/self-improving-rules.md`

行数分布：

| 类别 | 文件 | 行数 |
| --- | --- | ---: |
| 生成 registry | `registry.yaml` | 3798 |
| 生成归档/域文档 | `core-rules.md`、`schema-registry-rules.md`、`agent-runtime-rules.md` | 3766 |
| 规范化清单 | `enforcement-normalization.yaml` | 71 |
| README | `README.md` | 77 |
| 铁律 | `iron-rules.md` | 49 |
| 人工小型规则文档 | `goal-rules.md` 等 11 个文件 | 926 |

Markdown 总计 4818 行，YAML 总计 3869 行，目录总计 8687 行。体量主要来自生成文件，这是维护策略问题，不是单纯“文档太长就错误”。

## 2. 真实问题

### 2.1 [P1] 退出码语义漂移

证据：

- `.agent/rules/iron-rules.md:20` 到 `.agent/rules/iron-rules.md:35` 定义铁律退出码表。
- `.agent/rules/iron-rules.md:41` 明确 Traceability Gate 应返回 exit 9。
- `.agent/rules/agent-runtime-rules.md:421` 到 `.agent/rules/agent-runtime-rules.md:447` 的 `RULE-GOALCLI-EXIT-001` 定义另一组 exit code 语义。
- `.agent/rules/README.md:62` 声明 `exit_code` 字段应参见 `iron-rules.md`。
- `.agent/rules/registry.yaml:2402` 到 `.agent/rules/registry.yaml:2404` 保留生成后的 `RULE-GOALCLI-EXIT-001` 条目。

主要差异：

| exit code | `iron-rules.md` 含义 | `RULE-GOALCLI-EXIT-001` 含义 | 判断 |
| --- | --- | --- | --- |
| 0 | OK | PASS | 兼容 |
| 1 | 通用失败 | GENERAL_FAILURE | 兼容 |
| 2 | 参数错误 | POLICY_VIOLATION | 冲突 |
| 3 | 未定义 | SCHEMA_INVALID | 额外定义 |
| 4 | 未定义 | EVIDENCE_MISSING | 额外定义 |
| 5 | worktree / main 违规 | TRACEABILITY_BROKEN | 冲突 |
| 6 | schema 校验失败 | WORKTREE_INVALID | 冲突 |
| 7 | secret / 凭据泄漏 | SECRET_DETECTED | 兼容 |
| 8 | Evidence 缺失或伪造 | RELEASE_BLOCKED | 冲突 |
| 9 | Traceability 断链 | NEEDS_HUMAN_APPROVAL | 冲突 |
| 10 | Release 不完整 | INCONSISTENT_STATE | 冲突 |

这不是当前 P0 阻断，因为 `python3 scripts/verify_rules.py` 仍通过，且所有 active rule 的 `enforced_by` 命令都可解析。但它是 P1：一旦 CLI、CI、Agent 或 release gate 按不同文档解释相同 exit code，失败原因会被误分类。

修复建议：

- 以 `.agent/rules/iron-rules.md` 为唯一退出码语义来源，或明确将其降级为 “P0 gate exit table”。二者只能选一，不应并存两套无转换关系的表。
- 修正 registry 源数据中 `RULE-GOALCLI-EXIT-001` 的 `content` / `exit_code` 说明，再重新生成 `registry.yaml` 与域规则 Markdown。
- 为 `scripts/verify_rules.py` 或新增校验补一条 exit code 一致性检查，至少比较 `README.md` 声明的权威表与 registry 中 exit-code 规则条目。

### 2.2 [P1] 生成产物清单误标生成器

证据：

- `.agent/registries/generated-artifacts.yaml:23` 到 `.agent/registries/generated-artifacts.yaml:41` 将 `.agent/rules/registry.yaml`、`core-rules.md`、`schema-registry-rules.md`、`agent-runtime-rules.md` 的 `generated_by` 标为 `goalcli rules-verify`。
- `.agent/rules/README.md:39` 到 `.agent/rules/README.md:46` 写明生成命令是 `scripts/extract_rules.py` 与 `scripts/render_domain_rules.py`。
- `.agent/rules/registry.yaml:3` 头部写明由 `scripts/extract_rules.py` 生成。
- `.agent/rules/core-rules.md:3`、`.agent/rules/schema-registry-rules.md:3`、`.agent/rules/agent-runtime-rules.md:3` 头部写明由 `scripts/render_domain_rules.py` 生成。
- `scripts/verify_rules.py:2` 的说明是校验 registry 中的 `enforced_by` 命令，不是生成规则文件。

影响：

- 后续 Agent 或 CI 若按 `generated-artifacts.yaml` 决定再生成命令，会把 verifier 当 generator。
- 报告、release evidence、自动修复流程可能产生“验证通过但未重建产物”的假阳性。

修复建议：

- 将 `.agent/registries/generated-artifacts.yaml` 中相关条目的 `generated_by` 改为真实脚本。
- 若项目希望统一经 `goalcli rules-verify` 入口管理，应新增真实的 `goalcli rules-generate` 或 `goalcli rules-sync`，不要把 verify 命令写成 generator。

### 2.3 [P2] 生成文档残留历史文件名

证据：

- `.agent/rules/core-rules.md:738` 引用 `.agent/rules/07-worktree-rules.md`。
- `.agent/rules/schema-registry-rules.md:218`、`.agent/rules/schema-registry-rules.md:282`、`.agent/rules/schema-registry-rules.md:902` 引用 `.agent/rules/07-worktree-rules.md`。
- `.agent/rules/agent-runtime-rules.md:580` 引用 `.agent/rules/07-worktree-rules.md`。
- 当前目录真实文件为 `.agent/rules/worktree-rules.md`。

影响：

- 新贡献者按生成文档追溯规则时会找不到文件。
- 自动化 traceability 检查若未来启用路径存在性校验，会产生失败。

修复建议：

- 找到这些条目的源 registry 或原始规则定义，改为 `.agent/rules/worktree-rules.md` 后重新生成。
- 若这些引用用于历史迁移说明，应在内容里明确“历史路径”，而不是看起来像当前路径。

### 2.4 [P2] 生成文件体量集中，维护入口需要更明确

`core-rules.md`、`schema-registry-rules.md`、`agent-runtime-rules.md` 合计 3766 行，占 Markdown 规则文件约 78%。它们是生成产物，不应通过手工拆分解决。

风险：

- 人工修改 generated Markdown 会与下一次生成冲突。
- 真正需要维护的是 registry 源、生成脚本、渲染模板和校验策略。

修复建议：

- 在 README 中把“人改哪里 / 机器生成哪里 / 如何重建 / 如何验证”写成明确流程。
- 对 generated Markdown 增加更醒目的头部约束：不要手改，改源数据后重新生成。

### 2.5 [P3] 小型人工规则文档可做轻量去重

`commit-rules.md`、`pr-rules.md`、`evidence-rules.md`、`release-rules.md`、`harness-rules.md`、`security-rules.md`、`issue-rules.md`、`risk-decision-rules.md`、`self-improving-rules.md`、`worktree-rules.md` 等文件短小，部分与 registry/铁律存在重复表述。

当前只能定为 P3，因为缺少逐条重复率证据。建议后续用 rule ID 和段落相似度做量化，再决定是否合并。不要直接按“删除 60%”执行。

## 3. 对原报告建议的修正

应删除或改写的建议：

- 删除“补建 `.agent/rules/registry.yaml`”建议；该文件存在。
- 删除“恢复数字前缀规则文档”建议；当前 README 与目录结构都不支持这种命名。
- 删除“把退出码问题列为 P0 且声明当前规则系统不可用”的说法；当前 verifier 通过，问题是语义漂移而非现时不可执行。
- 删除“立即删除 60% 内容”的说法；没有量化证据，且多数体量来自 generated 文件。
- 删除“直接手工拆分 generated Markdown”的方向；应修源数据/生成器再重建。

应保留并重排的建议：

- 保留 exit code 一致性修复，但优先级调整为 P1。
- 保留路径漂移修复，但定位为 P2。
- 保留维护入口清晰化，但应围绕 generated/source contract，而不是围绕文件名重排。

## 4. 建议修复顺序

1. 修 `.agent/registries/generated-artifacts.yaml` 的 `generated_by` 字段，避免自动化继续使用错误入口。
2. 修 registry 源数据中的 exit code 规则，与 `iron-rules.md` 或新指定的唯一权威表对齐。
3. 重新生成 `.agent/rules/registry.yaml` 与域规则 Markdown。
4. 为 exit code 语义增加一致性校验，防止再次漂移。
5. 修正 `07-worktree-rules.md` 历史路径引用。
6. 最后再做人工小规则文件的去重和 README 流程补强。

## 5. 验证证据

已运行命令：

```bash
rg --files .agent/rules
wc -l .agent/rules/*
python3 scripts/verify_rules.py
rg -n "generated_by|scripts/extract_rules.py|scripts/render_domain_rules.py|rules-verify|07-worktree-rules|RULE-GOALCLI-EXIT-001|exit code|exit_code" .agent/rules .agent/registries/generated-artifacts.yaml README.md scripts/verify_rules.py
```

关键输出：

```text
rules total: 419
rules active: 363
goalcli subcommands available: 122
makefile targets available: 116
all active rules have valid enforced_by commands
```

验证结论：

- 当前 `.agent/rules/registry.yaml` 存在，且 active rule 的 `enforced_by` 命令校验通过。
- 当前报告已移除“registry 缺失”“恢复数字前缀命名”“P0 不可用”等错误判断。
- 当前仍保留的 P1/P2/P3 问题均有本地文件证据支撑。
