# 发布模板

## 占位符

- `{{MODULE_NAME}}`
- `{{MODULE_PATH}}`
- `{{PACKAGE_NAME}}`

## Release Gate

- `go test ./...`
- `go test -race ./...`
- `make boundary`
- `make security`
- `make contracts`
- `make evidence`

## Evidence

发布 Evidence 生成到 `release/manifest/latest.json`。

## 规则

- 没有 Evidence 不得发布。
- 不得在 release manifest、PR、Issue 或变更日志条目中包含原始凭据。
- 不得依赖 `x.go`。
