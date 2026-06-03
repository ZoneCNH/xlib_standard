# `.agent/rules/` — Goal Runtime 规则源

## 层级结构

```
.agent/rules/
├── iron-rules.md          # 7 条铁律 SSOT (P0 规则压缩)
├── registry.yaml          # 417 条规则机器化索引 (P0=119, P1=298)
├── README.md              # 本文件
├── goal-rules.md          # 域规则: Goal 对象模型
├── worktree-rules.md      # 域规则: Worktree-only
├── commit-rules.md        # 域规则: Commit
├── pr-rules.md            # 域规则: PR
├── issue-rules.md         # 域规则: Issue
├── evidence-rules.md      # 域规则: Evidence
├── release-rules.md       # 域规则: Release
├── harness-rules.md       # 域规则: Harness
├── security-rules.md      # 域规则: Security
├── risk-decision-rules.md # 域规则: Risk / Decision / Rollback
└── self-improving-rules.md # 域规则: Retrospective / Patch
```

## 权威顺序（裁决冲突时按本顺序）

1. `iron-rules.md` — 第一性铁律，不可豁免
2. `registry.yaml` — 规则唯一 ID + level + enforced_by + 退出码
3. 各域 `*-rules.md` — 详细叙述（人类阅读）
4. `docs/adr/ADR-*` — 决策记录
5. `.worktree/goal-patch.md` — 历史推导，**仅供考古，不可作为依据**

任何文档若与 `iron-rules.md` 或 `registry.yaml` 冲突，**以本目录为准**。

## 重新生成 `registry.yaml`

```bash
python3 scripts/extract_rules.py
```

`registry.yaml` 是 derived artifact，从 `.worktree/goal-patch.md` 提取并经人工分级映射生成。修改流程：

1. 调整 `scripts/extract_rules.py` 中的 `P0_PREFIXES` 或 `ENFORCED_BY` 映射
2. 重跑脚本
3. 提交时附说明: `Constraint:` 说明分级变化, `Tested:` 给出 P0/active 计数变化

## 字段说明

| 字段 | 含义 |
|---|---|
| `id` | 唯一规则 ID, 形如 `RULE-CORE-001` |
| `level` | `P0` (release-blocking) / `P1` (governance) / `P2` (advisory) |
| `title` | 一句话标题（中文） |
| `source_section` | `.worktree/goal-patch.md` 一级章节编号 |
| `source_section_title` | 章节标题 |
| `source_line` | 定义所在行号 |
| `enforced_by` | 实际执行该规则的命令（空 = 仅登记，尚未机器化） |
| `exit_code` | 违规时该命令应返回的标准退出码（见 `iron-rules.md`） |
| `status` | `active` (有 enforced_by) / `indexed` (仅登记) / `deprecated` |
| `duplicate_at` | 该 id 在源文档中被二次定义的行号（若存在） |

## 当前覆盖率

- 总规则数: 417
- P0 规则: 119
- P1 规则: 298
- 已 active (有 enforced_by): 173 (41%)
- 仅 indexed (待机器化): 244 (59%)

提高 active 比例是后续 Goal 的核心 KPI。
