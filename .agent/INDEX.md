# .agent Index

本文件是 Agent 上下文恢复的人类入口。机器可读分类索引是 `.agent/index.yaml`；两者冲突时以上游权威为准。

## 权威顺序

1. `CONSTITUTION.md`
2. `.agent/rules/`
3. `.agent/harness/`
4. `contracts/`
5. `docs/architecture/`
6. `AGENTS.md`
7. `CLAUDE.md` / tool-specific files
8. `README.md`
9. Issue / PR / temporary notes

## 恢复顺序

1. `CONSTITUTION.md`
2. `AGENTS.md`
3. `.agent/INDEX.md`
4. `.agent/context/`
5. `.agent/rules/`
6. `.agent/harness/`
7. 当前 Goal / Issue / Task
8. `contracts/`
9. `docs/architecture/`
10. 当前代码、最近 Evidence / Release Manifest / Retrospective

## 规则结构

- `.agent/rules/iron-rules.md`: P0 铁律压缩。
- `.agent/rules/registry.yaml`: 机器规则总账。
- `.agent/harness/harness.md` 与 `.agent/harness/harness.yaml`: gate 契约。
- `.agent/index.yaml`: `.agent` 文件分类与校验索引。
- `.agent/context/README.md`: 上下文恢复材料的边界说明。
- `docs/architecture/README.md`: 架构文档入口。
- `docs/goal/goal.md` 与 `docs/standard/`: 目标与标准叙事基线。
