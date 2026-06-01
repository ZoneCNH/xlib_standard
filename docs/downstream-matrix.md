# 下游矩阵

本矩阵定义 Full Goal Runtime v3.1 的下游生成目标。所有下游库必须保留 Required gate、release Evidence 和 `x.go` 反向依赖禁令。

| 库 | module path | package | layer | 允许依赖 | 禁止依赖 |
| --- | --- | --- | --- | --- | --- |
| `kernel` | `github.com/ZoneCNH/kernel` | `kernel` | L0 | Go 标准库、稳定 contracts | `x.go`、业务模型、profile runtime、生产密钥 |
| `configx` | `github.com/ZoneCNH/configx` | `configx` | L1 | `kernel`、显式配置来源 adapter | `x.go`、隐式 `/home/k8s/secrets/env/*` 读取、业务配置语义 |
| `observex` | `github.com/ZoneCNH/observex` | `observex` | L1 | `kernel`、metrics/logging/tracing contracts | 业务指标模型、应用告警策略 |
| `testkitx` | `github.com/ZoneCNH/testkitx` | `testkitx` | L1 | `kernel`、fake runtime、contract helpers | 真实生产连接、业务 fixture 默认值 |
| `postgresx` | `github.com/ZoneCNH/postgresx` | `postgresx` | L2 | `kernel`、`configx`、`observex` | 业务 repository、应用 transaction 编排 |
| `redisx` | `github.com/ZoneCNH/redisx` | `redisx` | L2 | `kernel`、`configx`、`observex` | 业务 key 语义、应用缓存策略 |
| `kafkax` | `github.com/ZoneCNH/kafkax` | `kafkax` | L2 | `kernel`、`configx`、`observex` | 业务 topic 设计、业务消息 schema |
| `taosx` | `github.com/ZoneCNH/taosx` | `taosx` | L2 | `kernel`、`configx`、`observex` | 业务指标模型、应用时序策略 |
| `ossx` | `github.com/ZoneCNH/ossx` | `ossx` | L2 | `kernel`、`configx`、`observex` | 业务文件生命周期策略 |
| `clickhousex` | `github.com/ZoneCNH/clickhousex` | `clickhousex` | L2 | `kernel`、`configx`、`observex` | 产品报表语义、业务查询模型 |

## Evidence 要求

- 每个库必须由 generator 产出可编译 module path、package name、README 和 docs。
- downstream integration 必须至少覆盖 `kernel`，完整 release Evidence 应覆盖本矩阵中的目标库或记录未覆盖原因。
- 任何库不得导入 `x.go` 或读取 `/home/k8s/secrets/env/*`。
