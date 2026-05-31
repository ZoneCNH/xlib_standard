# 变更日志

## v0.1.0 - 2026-06-01

### 新增

- 初始化 `baselib-template` 结构。
- 添加标准 Go 基础库包骨架。
- 添加 Makefile 命令。
- 添加 Harness Gate 脚本。
- 添加 GitHub Actions 工作流。
- 添加 contracts 文件。
- 添加 Agent 运行时模板。
- 添加 release manifest 模板。
- 添加 typed error、错误包装和 `ErrorKind` contract。
- 添加 client 生命周期、健康检查和请求扩展 metrics contract。
- 添加 health JSON contract 与 contracts 回归测试。
- 添加 config schema 到 `Config` 字段映射的 contract 回归测试。
- 添加 `scripts/render_template.sh`，支持生成 `foundationx` 等具体基础库。

### 安全

- 添加 Secret Gate。
- 配置脱敏规则覆盖 release Evidence 和日志可见内容。
- Boundary Gate 同时拦截 `github.com/bytechainx/x.go` 和 `github.com/ZoneCNH/x.go`。

### 治理

- 添加 Evidence 和复盘模板。
- `make release-check` 统一执行 CI、integration 和 manifest 生成。
- `make integration` 通过临时 `foundationx` 渲染和测试验证模板链路。
