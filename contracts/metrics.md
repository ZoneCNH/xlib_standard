# Metrics Contract

标准指标用于描述 `pkg/templatex` 暴露给调用方的最小可观测面。实现可以接入任意 metrics 后端，但指标名、类型和标签语义必须保持兼容。

## P0 指标（5 个）

| 指标 | 类型 | 标签 | 说明 |
| --- | --- | --- | --- |
| `client_created_total` | counter | — | 成功创建 client 的次数。 |
| `client_closed_total` | counter | — | 成功关闭 client 的次数；重复关闭不重复计数。 |
| `client_errors_total` | counter | `op`, `kind` | client 生命周期错误次数，`kind` 必须来自 error contract。 |
| `client_health_status` | gauge | `status` | 健康状态数值，healthy 为 `1`，其他状态为 `0`。 |
| `client_health_latency_ms` | histogram | `status` | 单次健康检查耗时，单位为毫秒。 |

## Label 约束

**允许的 Labels**: `op`, `kind`, `status`

**禁止的 Labels**: `user_id`, `request_id`, `trace_id`, `span_id`, `order_id`, `tenant_id`, `account_id`, `email`, `phone`, `token`, `secret`, `password`, `dsn`, `url`, `endpoint`
