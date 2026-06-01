# 分层规则

依赖方向只能从上层指向下层，不能反向依赖。

```text
Business / x.go
  -> L2 profile libraries
  -> L1 reusable infrastructure libraries
  -> L0 minimal primitives
  -> Standard / template contracts
```

`xlib-standard` 是 Standard 规则的独立来源；`baselib-template` 依赖该标准，提供 template、generator、Harness 和 Evidence 的 Go 实现载体。

## Standard

`xlib-standard` 所在层。它提供：

- 标准文档。
- 角色边界与分层规则。
- 标准 contracts 和 release Evidence 规则。
- 可被实现仓库复用的规范性模板元定义。

`baselib-template` 不是 Standard 的权威来源；它是 `xlib-standard` 在 Go 基础库模板中的实现仓库，负责模板目录、generator 契约、Harness gate、Evidence 生成和 CI 校验。Standard 与 `baselib-template` 都不得提供真实基础设施 runtime，不依赖业务仓库，也不依赖 `x.go`。

## L0

极小、稳定、无业务含义的 primitive，例如通用错误类型、脱敏规则或测试断言。L0 必须可独立测试。

## L1

可复用基础库，例如 `foundationx`、`postgresx`、`redisx`。L1 可以依赖 Standard 产物和 L0，必要时依赖 `foundationx`，但不得依赖 `x.go`。

## L2

面向具体基础设施 profile 的组合层，例如带特定驱动、协议或部署 profile 的 adapter。L2 可以依赖 L1，但不得包含业务流程。

## Business

业务应用、服务和 `x.go` 组合层。业务层消费基础库，负责业务配置、领域模型和流程编排。

## 违规示例

- L1 引入 `x.go` 包。
- `xlib-standard` 或 `baselib-template` 实现 PostgreSQL、Kafka 或 Redis 的真实 runtime。
- 基础库读取生产环境变量来创建默认 client。
- 基础库定义订单、账户、交易等业务模型。
