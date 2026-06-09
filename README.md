# xlib-standard

Go 基础库模板标准源。

## 定位

`xlib-standard` 定义并验证 FoundationX 基础库模板的最小能力：

- 独立 Go module
- 稳定公共 API（Config / Client / Error / Health / Metrics）
- 显式配置，不读取环境变量
- 可分类错误（8 种 ErrorKind）
- 健康检查（nil/closed/unhealthy 全覆盖）
- 低基数 metrics contract（5 个 P0 指标）
- 最小测试与示例
- 模板生成器
- release manifest

## 四项职责

| 职责 | 说明 |
|------|------|
| 标准源 | `docs/standard.md` 定义全部规则 |
| Go 参考模板 | `pkg/templatex/` 提供可编译的参考实现 |
| 生成器 | `scripts/render_template.sh` 从模板创建独立基座库 |
| 门禁 | `make ci` 验证模板和生成库合规性 |

## 快速开始

### 使用模板生成基座库

```bash
scripts/render_template.sh \
  --module-path github.com/your-org/your-lib \
  --package-name yourlib \
  --out /tmp/your-lib

cd /tmp/your-lib
GOWORK=off go test ./...
```

### 验证模板

```bash
GOWORK=off make ci
GOWORK=off make release-check
```

## 公共 API

```go
// Config — 显式配置，Validate 检查合法性，Sanitize 脱敏敏感字段
type Config struct { /* ... */ }
func (c Config) Validate() error
func (c Config) Sanitize() SanitizedConfig

// Client — 基座库客户端
type Client struct { /* ... */ }
func New(ctx context.Context, cfg Config, opts ...Option) (*Client, error)
func (c *Client) Close(ctx context.Context) error    // 幂等
func (c *Client) HealthCheck(ctx context.Context) HealthStatus
```

## ErrorKind（8 种）

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

## Metrics（5 个 P0 指标）

| 名称 | 类型 | Labels |
|------|------|--------|
| client_created_total | counter | — |
| client_closed_total | counter | — |
| client_errors_total | counter | op, kind |
| client_health_status | gauge | status |
| client_health_latency_ms | histogram | status |

## 目录结构

```text
.
├── contracts/          # JSON Schema 契约
├── docs/               # 文档（9 个文件）
├── examples/           # 最小可运行示例
├── pkg/templatex/      # Go 参考模板
├── scripts/            # 门禁脚本
├── testkit/            # 测试辅助工具
├── Makefile            # 构建目标
└── README.md
```

## 相关文档

| 文档 | 用途 |
|------|------|
| [docs/standard.md](docs/standard.md) | 完整标准规范 |
| [docs/api.md](docs/api.md) | 公共 API 说明 |
| [docs/config.md](docs/config.md) | 配置规范 |
| [docs/errors.md](docs/errors.md) | 错误处理规范 |
| [docs/health.md](docs/health.md) | 健康检查规范 |
| [docs/metrics.md](docs/metrics.md) | Metrics 规范 |
| [docs/testing.md](docs/testing.md) | 测试规范 |
| [docs/generation.md](docs/generation.md) | 生成器使用说明 |
| [docs/release.md](docs/release.md) | 发布流程 |

## License

[MIT](LICENSE)
