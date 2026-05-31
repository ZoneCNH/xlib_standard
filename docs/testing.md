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

生成的基础库必须保持测试独立于 `x.go`。
