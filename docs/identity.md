# xlib-standard 身份

## 我是谁

`xlib-standard` 是 FoundationX 的 **唯一标准源**（Standard Source）。它是 xlib 体系所有工程规则的权威定义者，也是 Go Reference Template、Generator、Harness Gate 和 Evidence Runtime 的唯一承载者。

> xlib-standard 是 FoundationX 的**标准锚点**。所有其他 15 个模块的一致性最终追溯到本模块。

## 五类职责

| 职责 | 说明 |
|------|------|
| **Standard Source** | xlib 体系工程规则的唯一真源 |
| **Go Reference Template** | 可编译、可测试的 Go library 骨架 |
| **Generator** | `render_template.sh` — 确定性渲染独立 Go module |
| **Harness Gate** | 9 个最小 CI 门禁串联执行 |
| **Evidence Runtime** | 可复现的 release manifest + checksum |

## 我不做什么

| 不是 | 原因 |
|------|------|
| **不承载业务运行** | 不作为其他模块的 runtime import 依赖 |
| **不包含业务域逻辑** | 交易/行情/风控逻辑属于分析域/决策域/执行域 |
| **不替代下游模块测试** | 业务测试/集成测试由各模块自己负责 |
| **不提供跨语言模板** | v1 仅覆盖 Go module |

## 宪法合规

| 条款 | 遵循方式 |
|------|----------|
| §1 P1 | Foundation 先边界后功能 — xlib-standard 是边界锚点 |
| §1 P2 | xlib-standard 不是运行时依赖 — 是标准/模板/Gate/Evidence 输入 |
| §3.3 | 门禁模块，无运行时依赖 |
