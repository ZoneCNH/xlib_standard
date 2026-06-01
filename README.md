# baselib-template

`baselib-template` 是一个独立的 Go 基础库模板模块，定位为生产级共享基础库基座。它提供可复用的包结构、稳定错误模型、可观测性 contract、release Evidence、Harness gate 和边界检查，供 `foundationx`、`configx`、`observex`、`postgresx`、`kafkax`、`redisx`、`taosx`、`ossx` 和 `testkitx` 等库生成时复用。

标准源仓库是 [`xlib-standard`](https://github.com/ZoneCNH/xlib-standard)，用于承载跨基础库共享的标准文本；本仓库只把该标准落到可编译模板、generator、Harness gate 和 Evidence 实现中。

## 目标

本模块的目标不是只生成目录，而是让新基础库从第一天就具备可验证的最小生产语义：

- 显式 `Config`，带稳定 validation error。
- `New`、`Close` 和 `HealthCheck` 具备上下文处理、幂等关闭和生命周期 metrics。
- `ErrorKind`、health status 和 metrics 名称由 `contracts/` 锁定。
- `make ci` 和 `make release-check` 生成可追溯 Evidence。

本仓库默认可编译模板包为 `templatex`。生成具体库时，优先使用 `scripts/render_template.sh` 统一替换 module name、module path、package name、Go module path、包目录、package 声明、imports、文档引用和 metrics prefix。

## 非目标

- 不依赖 `x.go`。
- 不包含 `x.go` 业务模型。
- 不隐式读取生产密钥。
- 不创建隐藏全局客户端。
- 不实现真实 `foundationx`、PostgreSQL、Redis、Kafka、OSS、ClickHouse 或 TDengine runtime。

## 标准边界

本仓库同时承担五个角色：

- Standard implementation：同步并落地独立标准仓库 [`xlib-standard`](https://github.com/ZoneCNH/xlib-standard) 的基础库标准、仓库角色、分层、DoD、Harness、Evidence、release、安全和 generator 契约。
- Template：提供可复制的 Go module、公共包、内部辅助、examples、contracts、scripts 和 CI。
- Generator：通过 [docs/generation.md](docs/generation.md) 和 `scripts/render_template.sh` 渲染具体基础库。
- Harness：通过 Makefile、scripts、CI 和 [.agent/harness.yaml](.agent/harness.yaml) 固化 required、extended 和 final gate。
- Evidence：通过 [docs/standard/evidence-protocol.md](docs/standard/evidence-protocol.md)、[docs/release.md](docs/release.md) 和 `release/manifest/latest.json` 记录可验证完成状态。

边界规则见 [模块边界](docs/standard/module-boundary.md)：`baselib-template` 只提供模板、generator、Harness 和 Evidence 实现，不承载业务语义；标准源以 [`xlib-standard`](https://github.com/ZoneCNH/xlib-standard) 为准；本仓库不依赖 `x.go`，不读取生产密钥，不创建隐藏全局客户端。

## 标准结构

- `pkg/templatex`：公共包 API 的可编译参考实现；渲染后会移动到 `pkg/<package-name>`。
- `internal/`：脱敏、校验和运行时说明等内部辅助代码。
- `testkit/`：可复用测试夹具和断言。
- `examples/`：最小使用示例。
- `contracts/`：JSON schema 和 metrics contract。
- `docs/`：规格、设计、API、配置、测试和发布模板。
- `scripts/`：Harness Gate 脚本。
- `.agent/`：Goal Runtime 工件、Evidence、评审、发布和复盘模板。
- `release/manifest/`：release manifest 模板；`latest.json` 由 release gate 生成并作为 Evidence artifact 保存。

## 文档入口

- [基础库标准索引](docs/standard/README.md)：P0 标准入口，覆盖仓库角色、分层、DoD、Harness、Evidence、release、安全和 generator 契约。
- [基础库总标准](docs/standard/xlib-standard.md)：同步 [`xlib-standard`](https://github.com/ZoneCNH/xlib-standard) 的公共 API、配置、错误、健康检查、metrics、测试、安全和发布规则。
- [仓库角色](docs/standard/repository-roles.md)：区分 `baselib-template`、`foundationx`、适配器库和 `x.go`。
- [Harness gate](docs/standard/harness-gates.md)：required、extended、generator 和 final gate 命令。
- [Evidence 协议](docs/standard/evidence-protocol.md)：`DONE with evidence:` 和 release manifest 要求。
- [规格](docs/spec.md)：模板能力、验收标准和可追踪性。
- [设计](docs/design.md)：模块边界、公共 API、错误、健康检查和指标设计。
- [API](docs/api.md)：`Config`、`Client`、typed error、health JSON 和 metrics contract。
- [配置](docs/config.md)：显式配置、validation 和脱敏规则。
- [生成](docs/generation.md)：从模板渲染 `foundationx` 等具体基础库。
- [错误模型](docs/errors.md)：`ErrorKind`、`NewError`、`WrapError` 和重试语义。
- [可观测性](docs/observability.md)：指标名、健康状态和 JSON 字段。
- [测试策略母版](docs/test-strategy.md)：Required、Extended 和 profile-specific gates。
- [测试](docs/testing.md)：单元、race、contracts、boundary 和 release 验证要求。
- [供应链](docs/supply-chain.md)：可校验 release Evidence、源码摘要、contract 指纹、依赖清单和 CI artifact。
- [发布](docs/release.md)：`release-check`、manifest 字段和 Evidence 规则。

## 命令

本地运行完整 gate 前需要安装 `golangci-lint` 和 `govulncheck`；CI 会显式安装这两个工具。缺少任一工具时，`make lint` 或 `make security` 必须失败，不允许把必需 gate 记录为跳过。

```bash
make ci
make ci-extended
make docs-check
GOWORK=off make release-check
make release-preflight VERSION=v0.1.0
make evidence
```

`release-check` 和 `release-check-extended` 已依赖 `docs-check`，用于在生成 Evidence 前确认标准文档入口、链接、模板占位符、关键文本和 release manifest 协议没有漂移。`docs-check` 是结构性 gate，不替代人工语义审查；标准变更仍需 reviewer 判断文案是否准确覆盖实际行为。

生成 `foundationx` 示例：

```bash
scripts/render_template.sh \
  --module-name foundationx \
  --module-path github.com/ZoneCNH/foundationx \
  --package-name foundationx \
  --out ../foundationx
```

发布式验证必须使用 `GOWORK=off`，避免父级或本地 `go.work` 改写 module 解析并掩盖模板独立性问题：

```bash
GOWORK=off make docs-check
GOWORK=off make release-check
```

## Evidence

完成需要 release manifest 和 CI Evidence。`release/manifest/latest.json` 是生成产物，不提交到源码历史；对应的 `release/manifest/latest.json.sha256` 也是生成产物，两者都必须保持在 `.gitignore` 中。manifest 会记录 module、commit、tree SHA、源码摘要、contract 指纹、依赖清单、工具版本、生成时间、工作区状态、gate 结果和这两个 Evidence artifact；`release-check` 会生成并校验 checksum，CI 会上传两者作为 artifact。`make release-evidence-check` 会验证 manifest 与当前仓库事实一致，`make release-final-check` 会额外要求工作区为 `clean`。最终完成声明必须包含 `DONE with evidence:`。

Agent 运行时模板位于 [.agent](.agent/)，其中 [goal-runtime](.agent/goal-runtime.md)、[agent-teams](.agent/agent-teams.md)、[review-template](.agent/review-template.md)、[evidence-template](.agent/evidence-template.md)、[release-template](.agent/release-template.md) 和 [retrospective-template](.agent/retrospective-template.md) 用于把标准、执行、评审和复盘连接到同一套 Evidence 协议。

## Smoke 覆盖

`go test ./...` 必须覆盖公共包、`internal/`、`contracts/`、`testkit/` 和 `examples/`。当前示例 smoke 测试会验证 `examples/basic` 输出模块名、`examples/config` 输出脱敏值、`examples/health` 输出健康状态，防止文档示例和模板行为漂移。

`scripts/run_fuzz_smoke.sh` 默认执行快速 fuzz smoke，`FUZZ_SMOKE_TIME` 未设置时每个 fuzz target 使用 `10s`。需要深度 fuzz 时显式设置更长时间，例如 `FUZZ_SMOKE_TIME=2m make fuzz-smoke`，并在最终 Evidence/DONE 说明中记录该时间配置。
