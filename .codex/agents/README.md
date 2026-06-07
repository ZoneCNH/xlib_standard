# Project Subagents

本目录保存 `xlib-standard` 的项目级 Codex native subagent 配置。Codex 会在项目范围内优先加载这里的 `.toml` 代理定义；这些代理必须服从根级 `AGENTS.md`、`CONSTITUTION.md` 与 `.agent/` 控制面的约束。

## Agents

| Agent | Use for | Default stance |
| --- | --- | --- |
| `xlib-explore` | 仓库事实查找、规则/契约/分层关系映射 | 只读 |
| `xlib-executor` | 有界实现、重构、配置修改 | 最小变更 + Gate |
| `xlib-verifier` | 完成声明、Evidence、Gate 与分层边界核验 | 只读优先 |
| `xlib-release-reviewer` | Release、manifest、checksum、downstream adoption 就绪性审查 | 只读 |

## Routing Notes

- 简单仓库查找优先交给 `xlib-explore`。
- 代码或配置实现交给 `xlib-executor`，但必须在非 `main` / `master` 的独立 worktree 中执行。
- 完成前交给 `xlib-verifier` 检查证据、Gate 和未验证风险。
- Release 或发布声明相关工作交给 `xlib-release-reviewer`；没有 release manifest、checksum、score、context release 和 final gate 证据时，不得声明 release-ready。
- `worker` 仍只用于 OMX team/swarm 运行时，不作为普通 subagent 使用。

修改本目录后，重新打开或重载 Codex 会话以确保新 agent 定义生效。
