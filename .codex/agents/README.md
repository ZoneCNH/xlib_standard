# Project Subagents

本目录保存 `xlib-standard` 的项目级 Codex native subagent 配置。Codex 会在项目范围内优先加载这里的 `.toml` 代理定义；这些代理必须服从根级 `AGENTS.md`、`CONSTITUTION.md` 与 `.agent/` 控制面的约束。

## Agents

| Agent | Use for | Default stance |
| --- | --- | --- |
| `xlib-explore` | 仓库事实查找、规则/契约/分层关系映射 | 只读 |
| `xlib-executor` | 有界实现、重构、配置修改 | 最小变更 + Gate |
| `xlib-verifier` | 完成声明、Evidence、Gate 与分层边界核验 | 只读优先 |
| `xlib-release-reviewer` | Release、manifest、checksum、downstream adoption 就绪性审查 | 只读 |
| `xlib-harness-selector` | 根据变更类型选择必须运行的 Harness Gate | 只读 |
| `xlib-harness-runner` | 执行已选本地 Gate 并整理 exact command / result 证据 | Gate 执行，不改源码 |
| `xlib-harness-auditor` | 审计 Harness、Makefile、registry 与规则之间的一致性 | 只读优先 |
| `xlib-claude-reviewer` | 通过 Claude CLI 做有界外部 review，并保留失败/配额/认证证据边界 | 只读优先 |
| `xlib-docs-contract-drift-auditor` | 审计 docs、contracts、examples、templates 与实现之间的漂移 | 只读 |
| `xlib-downstream-adoption-auditor` | 审计 downstream adoption 证明强度、治理包传播与本地/外部证据边界 | 只读 |
| `xlib-layer-boundary-reviewer` | 审查分层依赖方向、模块边界、契约绕过和 L2 耦合风险 | 只读 |
| `xlib-security-dependency-auditor` | 审计 secret、配置暴露、依赖风险与 security/dependency Gate 证据 | 只读 |

## Routing Notes

- 简单仓库查找优先交给 `xlib-explore`。
- 代码或配置实现交给 `xlib-executor`，但必须在非 `main` / `master` 的独立 worktree 中执行。
- 完成前交给 `xlib-verifier` 检查证据、Gate 和未验证风险。
- 不确定 Gate 覆盖面时先交给 `xlib-harness-selector`，输出必须区分 required、recommended、not applicable、not run。
- 需要执行本地 Gate 时交给 `xlib-harness-runner`；它只运行已批准范围内的非破坏性本地命令，并报告 exact command、exit status 和关键输出。
- 修改 `.agent/harness/`、`.agent/rules/`、`.agent/registries/`、`Makefile` 或命令入口时，交给 `xlib-harness-auditor` 检查 Harness 契约、registry、Makefile target 和 Evidence 要求是否同步。
- 需要外部 Claude CLI 复核时交给 `xlib-claude-reviewer`；它必须记录 exact command、target、exit status 和认证/配额/工具阻塞，不得把 Claude 输出替代 Harness Gate。
- Release 或发布声明相关工作交给 `xlib-release-reviewer`；没有 release manifest、checksum、score、context release 和 final gate 证据时，不得声明 release-ready。
- 声称 downstream adoption、下游同步或模板采纳时交给 `xlib-downstream-adoption-auditor`；本地 template/contract 通过只能声明 local proof。
- 文档、contract、schema、example、template、CLI 契约或公共行为同步风险交给 `xlib-docs-contract-drift-auditor`。
- 分层、模块边界、L2 adapter、公共依赖方向或 contract bypass 风险交给 `xlib-layer-boundary-reviewer`。
- secret、配置默认值、依赖、网络行为、security gate 或 vulnerability proof 风险交给 `xlib-security-dependency-auditor`。
- `worker` 仍只用于 OMX team/swarm 运行时，不作为普通 subagent 使用。

修改本目录后，重新打开或重载 Codex 会话以确保新 agent 定义生效。
