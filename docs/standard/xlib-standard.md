# 基础库总标准

本文定义 `x.go` 基础库体系的最小生产标准。`xlib-standard` 是独立标准仓库，规范地址为 `https://github.com/ZoneCNH/xlib-standard`；`baselib-template` 是该标准在 Go 基础库生态中的模板、generator、Harness 和 Evidence 实现仓库。生成库必须继承这些约束，除非在自身 ADR 中写明更严格的 profile-specific 规则。

## 仓库定位

- `xlib-standard`：独立标准仓库，沉淀跨语言、跨实现的基础库规则、角色边界、分层、Evidence 和 release 要求。
- `baselib-template`：模板/generator/Harness/Evidence 实现仓库，负责把 `xlib-standard` 的规则落到 Go module 模板、contracts、CI 和生成验证中。
- 生成库：由 `baselib-template` 渲染得到的具体基础库，必须遵守 `xlib-standard`，并通过自身 Harness 和 Evidence 证明符合性。
- 标准内容以 `https://github.com/ZoneCNH/xlib-standard` 为长期来源；`baselib-template` 中的标准文档用于实现对齐和本仓库 gate 校验，不改变独立标准仓库的权威定位。

## 标准目标

- 提供独立 Go module，可在没有 `x.go` 的情况下构建、测试和发布。
- 提供稳定公共 API、显式配置、可验证错误模型、健康检查、metrics contract 和 release Evidence。
- 让新基础库从创建开始就具备 Harness gate、contracts、CI、文档、评审和复盘入口。

## 分层

基础库体系按职责分层：

- Standard：由独立仓库 `xlib-standard` 定义；`baselib-template` 只承载 Go 模板、generator、contracts、Harness 和 Evidence 的实现副本与校验入口。
- L0：语言级、无业务依赖的公共能力。
- L1：面向具体中间件或基础设施的库，例如 `postgresx`、`redisx`、`kafkax`。
- L2：组合多个基础能力的技术组件。
- Business：业务服务，只消费基础库，不向基础库反向注入业务模型。

## 依赖方向

- 依赖方向只能从 Business 指向 L2/L1/L0/Standard，或从具体库继承 `xlib-standard`。
- `xlib-standard` 不得依赖 `x.go`、业务仓库或 profile-specific runtime。
- `baselib-template` 只能把 `xlib-standard` 落成模板/generator/Harness/Evidence，不得成为替代标准源。
- `baselib-template` 不得依赖 `x.go`、业务仓库、profile-specific runtime 或生成库的真实 runtime。
- L0/L1/L2 不得依赖 Business。
- profile-specific 扩展必须在自身 ADR、contracts 或 Harness profile 中声明，不能弱化本标准。

## 公共 API

公共能力放在渲染后的 `pkg/<package-name>`。导出 API 必须满足：

- 使用显式 `Config`，不得隐式读取生产环境凭据。
- `New` 返回可关闭的 client 或服务对象，并明确 validation failure。
- `Close` 必须幂等，重复调用不应 panic。
- `HealthCheck` 必须接收 `context.Context`，输出稳定 health status。
- 错误必须可分类，`ErrorKind` 由 `contracts/errors.schema.json` 约束。

## Config

- 配置必须可脱敏，脱敏规则覆盖 token、secret、password、key 等敏感字段。
- 默认值必须安全、可测试、可解释。
- 缺失必需字段时返回 typed validation error。
- 示例配置不得包含真实凭据或生产地址。

## Error

- validation、closed、unhealthy 和 internal failure 必须可区分。
- 公共错误分类必须稳定，并由 contract 约束。
- 错误消息不得包含 secret、token、password、private key 或完整连接凭据。

## Logging

- 日志不得输出 secret、token、password、private key 或连接凭据。
- 日志字段必须可审计，不得隐式暴露配置原文。
- 示例和测试日志只能使用假数据或脱敏数据。

## Metrics

- metrics 名称和 label 由 `contracts/metrics.json` 约束。
- 新增 metrics 必须更新 contract、测试和 Evidence。
- label 不能包含高基数字段、用户凭据或业务私有标识。

## Lifecycle

- 初始化、关闭、健康检查和错误路径必须可单元测试。
- 并发或共享状态必须通过 race gate 验证。
- 不得创建隐藏全局 client、后台 goroutine 或不可关闭资源。

## Health

- health JSON 字段由 `contracts/health.schema.json` 约束。
- `HealthCheck` 必须接收 `context.Context`，并在关闭后返回稳定 unhealthy 状态。
- 健康状态不得依赖生产网络或真实凭据才能测试。

## Contract

- contracts 是公共行为的一部分，不是文档附录。
- 错误、health、metrics 和 manifest 的 schema 变更必须经过 contract gate。
- contract 变更必须说明兼容性影响，并在 PR 中列出验证命令。

## Testing

Required gate：

- `GOWORK=off make ci`
- `GOWORK=off make integration`
- `CHECK_STATUS=passed GOWORK=off make evidence`
- `RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check`

Extended gate：

- `GOWORK=off make ci-extended`
- `GOWORK=off make release-check-extended`

Final gate：

- `GOWORK=off make release-final-check`
- `GOWORK=off make release-preflight VERSION=<version>`

## Security

- `make security` 必须包含 `govulncheck ./...` 和 `scripts/check_secrets.sh`。
- 缺少 `govulncheck` 时必须失败，不得记录为跳过。
- 仓库、CI、manifest、PR 和文档不得提交真实凭据。

## Release

- `release/manifest/latest.json` 是生成产物，只作为 CI artifact 和本地 Evidence 使用。
- 发布前必须通过 release Evidence 校验，并记录 commit、tree SHA、源码摘要、contract 指纹、依赖清单和 gate 状态。
- dirty workspace 不能通过 `release-final-check`，也不能宣称 final release ready。

## Evidence

- Evidence 是完成声明的一部分；没有命令或 artifact 支撑时不得宣称完成。
- 目标、Issue、Task 和 Release 必须记录 scope、gate、artifact、known gap。
- release Evidence 必须包含 manifest、source digest、contract fingerprint、dependency list、tool versions 和 gate status。

## Retrospective

- 标准、generator、Harness 或 release 流程变更后必须记录复盘入口。
- 复盘输出可以形成 Prompt Patch、Harness Patch、Rule Patch 或 CI Gate Suggestion。
- 复盘不得绕过 required gate；只能补强规则或记录明确 backlog。

## 禁止项

- 禁止依赖 `x.go`。
- 禁止承载业务模型或业务流程。
- 禁止隐式读取生产密钥。
- 禁止隐藏全局客户端。
- 禁止无 Evidence 声称 `DONE`。

## 下游生成兼容

- generator 必须能渲染 `foundationx` 和 `corekit` 代表下游。
- 生成库必须无旧 module path、`pkg/templatex`、`package templatex` 和模板占位符残留。
- 生成库必须在 `GOWORK=off` 下通过 `go test ./...`、contracts、boundary、evidence 和 release Evidence 校验。

## 完成声明格式

完成声明必须使用：

```text
DONE with evidence:
- scope: <task|issue|goal|release>
- gates:
  - <command>: <passed|failed|blocked> <short evidence>
- artifacts:
  - <path>: <purpose>
- known gaps:
  - <none or explicit blocker>
```
