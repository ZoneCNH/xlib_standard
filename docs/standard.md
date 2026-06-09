# 基础库标准规范

> `xlib-standard` 定义的 Go 基础库模板标准。所有基座库必须遵循本规范。

最后更新：2026-06-09

---

## 1. 目的

本规范定义 FoundationX 基座库的最小公共 API、配置、错误处理、健康检查、metrics、测试、生成和发布标准。所有基座库（`configx`、`observex`、`redisx`、`kafkax`、`natsx`、`postgresx`、`clickhousex`、`ossx`、`taosx` 等）必须遵循本规范。

---

## 2. 非目标

以下不在本规范范围内：

- 业务逻辑
- 具体中间件实现细节
- 下游模块的领域特定规则
- Goal Runtime、Evidence Runtime、Agent Runtime
- Debt Governance、Branch Governance

---

## 3. 仓库结构

### 3.1 标准目录

```text
.
├── contracts/          # JSON Schema 契约
├── docs/               # 文档
├── examples/           # 最小可运行示例
├── pkg/<package-name>/ # 核心包
├── scripts/            # 门禁脚本
├── testkit/            # 测试辅助工具
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### 3.2 命名规则

| 项目 | 规则 | 示例 |
|------|------|------|
| 包名 | 小写，无下划线 | `templatex`, `configx` |
| 函数/类型 | PascalCase | `NewClient`, `HealthStatus` |
| 常量 | UPPER_SNAKE_CASE | `KIND_VALIDATION` |
| 文件名 | snake_case.go | `health_check.go` |
| 测试文件 | `*_test.go` | `client_test.go` |

### 3.3 go.mod 规则

```text
module github.com/ZoneCNH/<module-name>

go 1.22
```

- 独立 module path，不使用 replace 指令
- Go 最低版本 >= 1.22
- 无第三方运行时依赖

---

## 4. 公共 API

### 4.1 Config

```go
type Config struct {
    Name    string
    Timeout time.Duration
    // ... 具体库扩展
}

func (c Config) Validate() error
func (c Config) Sanitize() SanitizedConfig
```

### 4.2 Client

```go
type Client struct { /* ... */ }

func New(ctx context.Context, cfg Config, opts ...Option) (*Client, error)
func (c *Client) Close(ctx context.Context) error
func (c *Client) HealthCheck(ctx context.Context) HealthStatus
```

### 4.3 Error

```go
type Error struct {
    Kind      ErrorKind
    Op        string
    Message   string
    Err       error
    Retryable bool
}

func NewError(kind ErrorKind, op, message string, err error) *Error
func WrapError(kind ErrorKind, op string, err error) *Error
func IsKind(err error, kinds ...ErrorKind) bool
```

### 4.4 Metrics

```go
type Metrics interface {
    IncrCounter(name string, labels map[string]string)
    SetGauge(name string, value float64, labels map[string]string)
    ObserveHistogram(name string, value float64, labels map[string]string)
}
```

---

## 5. Config

### 5.1 规则

1. 所有配置显式传入，不读取环境变量
2. 不读取文件
3. 不内置生产 endpoint
4. `Validate` 检查必填字段和非法值（如负数 timeout）
5. `Sanitize` 脱敏字段名包含 `secret`、`token`、`password`、`key`、`credential`、`dsn`、`url` 的字段

### 5.2 示例

```go
cfg := Config{
    Name:    "my-service",
    Timeout: 5 * time.Second,
}

if err := cfg.Validate(); err != nil {
    // 处理错误
}

sanitized := cfg.Sanitize()
// sanitized 中敏感字段被替换为 "***"
```

---

## 6. Error

### 6.1 ErrorKind 枚举（8 种）

| Kind | 含义 | retryable |
|------|------|-----------|
| validation | 输入校验失败 | false |
| config | 配置错误 | false |
| connection | 连接失败 | true |
| auth | 认证/授权失败 | false |
| timeout | 超时 | true |
| unavailable | 服务不可用 | true |
| closed | 客户端已关闭 | false |
| internal | 内部错误 | false |

### 6.2 错误包装

```go
if err != nil {
    return WrapError(KindConnection, "Connect", err)
}
```

- 使用 `fmt.Errorf("op: %w", err)` 包装
- `errors.Is` / `errors.As` 能穿透包装层
- `IsKind(err, KindTimeout, KindUnavailable)` 用于分支判断

---

## 7. Health

### 7.1 HealthStatus 结构

```go
type HealthStatus struct {
    Name      string
    Status    string // "healthy" | "degraded" | "unhealthy"
    Message   string
    CheckedAt time.Time
    LatencyMs int64
    Metadata  map[string]string
}
```

### 7.2 规则

| 场景 | 预期状态 |
|------|----------|
| nil context | unhealthy |
| canceled context | unhealthy |
| zero-value client | unhealthy |
| closed client | unhealthy |
| initialized and open client | healthy |
| degraded | 只作为具体库扩展 |
| Metadata | 默认不填 |

---

## 8. Metrics

### 8.1 P0 指标（5 个）

| 名称 | 类型 | Labels | 说明 |
|------|------|--------|------|
| client_created_total | counter | — | 客户端创建计数 |
| client_closed_total | counter | — | 客户端关闭计数 |
| client_errors_total | counter | op, kind | 错误计数 |
| client_health_status | gauge | status | 健康状态 |
| client_health_latency_ms | histogram | status | 健康检查延迟 |

### 8.2 Label 约束

**允许的 Labels**: `op`, `kind`, `status`

**禁止的 Labels**: `user_id`, `request_id`, `trace_id`, `span_id`, `order_id`, `tenant_id`, `account_id`, `email`, `phone`, `token`, `secret`, `password`, `dsn`, `url`, `endpoint`

---

## 9. Testing

### 9.1 必需测试

每个基座库必须包含以下测试：

| 测试文件 | 覆盖内容 |
|----------|----------|
| config_test.go | Validate 必填字段、Sanitize 脱敏 |
| errors_test.go | NewError、WrapError、IsKind |
| metrics_test.go | NoopMetrics 不 panic、指标名匹配 contract |
| client_test.go | New 幂等性、Close 幂等性 |
| health_test.go | nil/closed/healthy 场景 |

### 9.2 禁止依赖

- 真实基础设施（PostgreSQL、Redis、Kafka、NATS、OSS、ClickHouse）
- 生产网络、生产密钥
- x.go、业务仓库

### 9.3 覆盖率

最低 80% 行覆盖率。

---

## 10. Generator

### 10.1 使用方式

```bash
scripts/render_template.sh \
  --module-path <module-path> \
  --package-name <package-name> \
  --out <dir> \
  [--module-name <module-name>]
```

### 10.2 执行步骤

1. 检查 `--module-path`、`--package-name`、`--out` 参数
2. 拒绝写入非空目录
3. 复制模板源码
4. 替换 module path
5. 将 `pkg/templatex` 移动为 `pkg/<package-name>`
6. 替换包名、README、docs、contracts 中的模板占位符
7. 在生成目录执行 `GOWORK=off go test ./...`

### 10.3 生成库验证

```bash
! grep -R "templatex" <output-dir> --exclude-dir=.git
! grep -R "xlib-standard" <output-dir> --exclude-dir=.git
```

---

## 11. Gate

### 11.1 Makefile 目标

```text
fmt → vet → lint → test → race → contracts → boundary → render-smoke → security
```

### 11.2 Gate 脚本

| 脚本 | 检查内容 |
|------|----------|
| check_contracts.sh | contracts/*.schema.json 存在且合法 |
| check_boundary.sh | 禁止 x.go/internal、生产密钥、旧名残留 |
| check_security.sh | 禁止 AWS key、私钥、明文密码 |
| check_rendered_template.sh | 生成库无模板残留 |

### 11.3 CI 流程

```yaml
on: [pull_request, push:main]
steps:
  - checkout
  - setup-go
  - GOWORK=off make ci
  - GOWORK=off make release-check
```

---

## 12. Release

### 12.1 Release Manifest

```json
{
  "module_path": "github.com/ZoneCNH/xlib-standard",
  "package_name": "templatex",
  "version": "v1.0.0",
  "commit": "<git-commit>",
  "tree_sha": "<git-tree-sha>",
  "go_version": "<go-version>",
  "contracts_sha256": "<sha256>",
  "gates": {
    "fmt": "passed",
    "vet": "passed",
    "lint": "passed",
    "test": "passed",
    "race": "passed",
    "contracts": "passed",
    "boundary": "passed",
    "render_smoke": "passed",
    "security": "passed"
  },
  "generated_at": "<utc-rfc3339>"
}
```

### 12.2 完成声明

```text
DONE with evidence:
- manifest: release/manifest/latest.json
- checksum: release/manifest/latest.json.sha256
- gates: passed
```

### 12.3 Semver 规则

- 主版本号变更：不兼容的 API 变更
- 次版本号变更：向后兼容的功能新增
- 修订号变更：向后兼容的 bug 修复
