# Gate

- Boundary Gate：`make boundary`
- Secret Gate：`make security`，必须包含 `govulncheck ./...` 和密钥扫描
- Contract Gate：`make contracts`
- Format Gate：`make fmt`
- Static Check Gate：`make vet`
- Lint Gate：`make lint`，缺少 `golangci-lint` 时失败
- Unit Test Gate：`make test`
- Race Test Gate：`make race`
- Evidence Gate：`make evidence`
- Release Gate：`make release-check`
