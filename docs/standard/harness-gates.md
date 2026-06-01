# Harness Gate 标准

Harness gate 是完成声明的证据来源。所有命令默认在仓库根目录运行；验证模板独立性时使用 `GOWORK=off`。

## Required Gate

| Gate | 命令 | 目的 |
| --- | --- | --- |
| Format | `GOWORK=off make fmt` | 保持 Go 格式稳定 |
| Vet | `GOWORK=off make vet` | 基础静态检查 |
| Lint | `GOWORK=off make lint` | `golangci-lint` 强制检查，缺失时失败 |
| Unit | `GOWORK=off make test` | 单元和示例 smoke |
| Race | `GOWORK=off make race` | 并发安全基线 |
| Boundary | `GOWORK=off make boundary` | 模块边界和模板禁止项 |
| Security | `GOWORK=off make security` | `govulncheck` 和 secret scan |
| Contracts | `GOWORK=off make contracts` | schema、metrics 和 manifest contract |
| Integration | `GOWORK=off make integration` | generator 和 downstream smoke |
| Evidence | `CHECK_STATUS=passed GOWORK=off make evidence` | 生成 release manifest |
| Release Evidence | `RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check` | 校验 manifest 与仓库事实 |

## Extended Gate

- `GOWORK=off make property`
- `GOWORK=off make golden`
- `GOWORK=off make fuzz-smoke`
- `GOWORK=off make ci-extended`
- `GOWORK=off make release-check-extended`

## Generator Gate

Generator gate 必须证明模板能生成代表性下游，而不是只证明 `baselib-template` 自身可用。

必须渲染：

- `foundationx`
- `corekit`

每个渲染结果必须满足：

- 无 module name、module path、package name 等模板 token 残留。
- 无旧 module path import。
- 无 `pkg/templatex` 或 `package templatex` 残留。
- `GOWORK=off go test ./...` 通过。
- `GOWORK=off make contracts` 通过。
- `GOWORK=off make boundary` 通过。
- `CHECK_STATUS=passed GOWORK=off make evidence` 通过。
- `RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check` 通过。

## Final Gate

- `GOWORK=off make release-final-check`
- `GOWORK=off make release-preflight VERSION=<version>`

Final gate 要求工作区状态、版本参数、release Evidence 和所有 required gate 都满足发布条件。开发中有未提交变更时可以运行前置 gate，但不能宣称 final release ready。

## 失败处理

- 缺工具导致失败时，记录工具名和失败命令。
- gate 失败后先修复根因，再重新运行最小可证明命令。
- 不得把 required gate 降级为 optional。
