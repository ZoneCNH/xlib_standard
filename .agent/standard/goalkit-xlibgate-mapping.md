# goalkit ↔ xlibgate 合约-实现映射

> 解决 goal-patch.md v2.0+ 全文使用 `goalkit` 命名、而本仓库实现是 `xlibgate` 的命名冲突。  
> **核心原则**：合约 ≠ 实现。**不做物理改名。**

---

## 命名分工

| 概念 | 含义 | 归属 |
|---|---|---|
| **goalkit** | Goal Runtime 标准合约：CheckResult schema、Exit Code 表、命令语义 | 标准层（面向下游 / 文档） |
| **xlibgate** | xlib-standard 仓库的 goalkit 合约 Go 实现 | 实现层（本仓库 CLI） |
| 下游 `<repo>gate` | 下游库可选的另一种 goalkit 实现 | 各下游库 |

---

## 命令映射表

| goalkit 合约（原文 §276） | xlibgate 实现 | 状态 |
|---|---|---|
| `goalkit doctor` | `xlibgate doctor` | ✅ |
| `goalkit worktree check` | `xlibgate worktree-guard` | ✅ |
| `goalkit secret check` | `xlibgate secrets` / `xlibgate security` | ✅ |
| `goalkit schema check` | （未独立实现，融入各 check） | ⚠️ |
| `goalkit evidence check` | `xlibgate evidence-check` | ✅ |
| `goalkit traceability check` | `xlibgate goal-acceptance` 等 | ⚠️ |
| `goalkit release check` | `xlibgate release-evidence-check` | ✅ |
| `goalkit retro check` | （隐含于 `xlibgate score` 计分中） | ⚠️ |
| `goalkit audit goal` | `xlibgate score` | ✅ |
| `goalkit bootstrap repo` | `make bootstrap`（待补） | 🔴 |

---

## 退出码（沿用原文 §280）

| Code | 含义 | xlibgate 是否兼容 |
|---|---|---|
| 0 | 通过 | ✅ |
| 1 | 业务失败（违规） | ✅ |
| 2 | 用法错误 | ✅ |
| 3 | 配置缺失 | ⚠️（部分） |
| 4 | 环境异常 | ⚠️（部分） |

---

## 下游采用

下游库可以选择：
1. **复制 xlibgate**（最简单，直接 vendor）
2. **自写 `<repo>gate`**（需声明实现 goalkit-v0.1.0 合约）
3. **依赖标准 `goalkit` 二进制**（若未来真的需要独立内核，再创建 `tools/goalkit`）

**当前阶段**：所有下游优先用方案 1。


---

> 本文件被 `scripts/check_docs.sh` 列为 required，删除会阻断 `make docs-check` / `make governance-check` / release pipeline。
