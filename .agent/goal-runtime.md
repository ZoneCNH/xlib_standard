# Goal Runtime

本文件定义 Agent 执行 `baselib-template` 目标时的运行时契约。目标文本见 [.agent/goal.md](goal.md)，标准入口见 [docs/standard/README.md](../docs/standard/README.md)。

## 阶段

1. Context：读取 ADR、README、docs、Makefile、contracts、scripts 和 `.agent/`。
2. Goal：确认目标、非目标、禁止项和完成声明格式。
3. Spec：对照标准文档检查缺口。
4. Plan：列出最小可验证变更。
5. Implementation：只修改目标要求的文件。
6. Test：运行相关 gate。
7. Evidence：生成或校验 release manifest。
8. Review：按 `.agent/review-template.md` 检查。
9. Retrospective：失败时按 `.agent/retrospective-template.md` 记录补丁。

## 停止条件

- 所有 ADR 必需项已实现。
- 相关 gate 有新鲜 Evidence。
- known gaps 已明确记录。
- 最终声明使用 `DONE with evidence:`。

未满足上述条件时不得声称完成。
