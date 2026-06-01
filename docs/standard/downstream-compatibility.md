# 下游兼容标准

下游兼容验证用于证明模板可以生成真实基础库，而不是只在自身仓库内通过测试。

## 代表性下游

当前代表性下游：

- `foundationx`：组织内常规 module path，使用 `github.com/ZoneCNH/foundationx`。
- `corekit`：中性 module path，使用 `example.com/acme/corekit`，用于证明 generator 不依赖固定组织、GitHub owner 或 module prefix。

未来 L1/L2 profile 可以扩展：

- `postgresx`
- `redisx`
- `kafkax`
- `taosx`
- `ossx`
- `clickhousex`

## 兼容要求

生成库必须满足：

- 独立 `go test ./...` 可运行。
- module path、package name 和 README 已替换。
- contracts 和 metrics schema 可校验。
- `.agent/`、Issue/PR 模板、release manifest 规则被复制。
- 无 `x.go` 依赖。
- 无模板占位符残留。
- 无 `templatex_` metrics prefix 残留。
- `GOWORK=off go mod tidy` 后 `go.mod` 和 `go.sum` 保持 clean。

## Tool Matrix

| 工具 | 用途 | 要求 |
| --- | --- | --- |
| Go 1.23 | 编译、测试、`go mod tidy`、dependency list | 本地和 CI 一致 |
| `make` | 执行 Harness gate | 必须可运行 required targets |
| `git` | 初始化临时下游、检查 clean diff、计算 commit/tree | integration 和 Evidence 必需 |
| `golangci-lint` | `make lint` | 缺失时必须失败 |
| `govulncheck` | `make security` | 缺失时必须失败 |
| `python3` | docs link checker | `make docs-check` 必需 |
| `sha256sum` | 计算 `latest.json` hash | CI artifact Evidence 必需 |
| GitHub Actions artifact | 保存 `release/manifest/latest.json` | 远端 release Evidence 必需 |

## 验证方式

`GOWORK=off make integration` 是默认下游兼容 gate。它应覆盖 generator smoke、foundationx/corekit 代表路径和关键边界检查。

当新增 profile 时，在不污染默认 `make ci` 的前提下补充 profile-specific smoke 或 extended gate。

代表性下游验证至少覆盖：

```bash
GOWORK=off go mod tidy
GOWORK=off go test ./...
GOWORK=off make contracts
GOWORK=off make boundary
CHECK_STATUS=passed GOWORK=off make evidence
RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

这些命令需要在 `foundationx` 和 `corekit` 的渲染结果中通过，失败时不得宣称 downstream compatible。

## 兼容破坏

以下情况视为兼容破坏：

- 删除 generator 必需输入。
- 破坏 `Config`、`New`、`Close`、`HealthCheck` 的基础语义。
- 移除 required gate。
- 让生成库依赖 `x.go`。
- 让生成库不能独立生成 Evidence。
