# 测试模板

## 占位符

- `{{MODULE_NAME}}`
- `{{PACKAGE_NAME}}`

## 必需 Gate

本地执行 gate 前必须可用：

- `golangci-lint`
- `govulncheck`

缺少上述工具时，`make lint` 或 `make security` 必须失败。

- `make fmt`
- `make vet`
- `make lint`
- `make test`
- `make race`
- `make boundary`
- `make security`
- `make contracts`
- `make integration`
- `make evidence`

## 必需覆盖范围

- `go test ./...` 必须覆盖公共包、`internal/`、`contracts/`、`testkit/` 和 `examples/`。
- 配置校验。
- 配置脱敏。
- typed error kind 和 wrapped cause。
- 客户端创建、取消 context、过期 context。
- 幂等关闭、zero-value client、取消 context。
- 健康与非健康状态检查。
- 健康检查 JSON 字段 contract。
- 生命周期 metrics 和健康 metrics。
- `contracts/` 与公共常量同步。
- `contracts/config.schema.json` 与 `Config` 字段映射同步。
- `scripts/render_template.sh` 生成的临时 `foundationx` 可以通过 `GOWORK=off go test ./...`。

## 示例与 testkit Smoke

- `examples/basic` 必须输出当前 module name。
- `examples/config` 必须输出脱敏后的 secret 值。
- `examples/health` 必须输出 `healthy`。
- `testkit` 必须验证 `Config("fixture")` 生成可通过 `Validate` 的测试配置。
- `testkit.RequireNoError` 必须接受 `nil`，作为生成库测试断言的最小契约。

生成的基础库必须保持测试独立于 `x.go`。
