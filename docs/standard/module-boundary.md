# 模块边界

本文件约束 `xlib-standard`、`baselib-template` 和生成库的模块边界。边界违规必须在 PR、review 或 release gate 中阻断。

## `baselib-template` 允许内容

- 模板包 `pkg/templatex`。
- 内部辅助示例，例如 validation、sanitize、runtime 说明（仅限示例或模板层）。
- generator、模板元数据和渲染脚本。
- contracts schema、metrics contract、标准文档的实现对齐副本与 gate 校验入口。
- 示例程序、测试夹具、CI、Makefile、scripts。
- Agent 模板、Issue/PR 模板和 release Evidence 生成器实现。

## `baselib-template` 禁止内容

- 真实 `foundationx`、PostgreSQL、Redis、Kafka、OSS、ClickHouse、TDengine runtime。
- `x.go` import 或以 `x.go` 为构建前提。
- 业务领域模型、业务 repository、业务消息 schema。
- 生产密钥、真实连接串或默认生产 endpoint。
- 隐藏全局 client、隐式后台进程或不可关闭资源。
- 替代 `xlib-standard` 定义或维护独立标准正文，包括标准契约规范、标准发布证据规则正文。

## `xlib-standard` 允许内容

- P0 标准文档与版本说明。
- 仓库角色、分层、模块边界和下游兼容性规则。
- 标准化 contracts schema、metrics contract 和 release Evidence 规则。
- 与标准相关的非实现性参考模板。

## `xlib-standard` 禁止内容

- generator、模板渲染脚本或运行时代码。
- 基于业务或 profile 的 runtime/Harness 实现。
- 生产凭据、真实连接串或 profile runtime 实现。
- `baselib-template` 特有 CI、release manifest 生成器或 Evidence 执行脚本。

## 生成库允许内容

- 替换后的 `go.mod` module path、package name、README 和 docs。
- 该库 profile 所需的公共 API、配置、错误、健康检查、metrics 和 tests。
- 必要的 profile-specific contracts 和 examples。

## 生成库禁止内容

- 反向依赖业务层。
- 绕过 `Config` 隐式加载生产凭据。
- 删除 Required gate。
- 把 release manifest 生成产物提交进源码历史。

## 边界验证

每次涉及包结构、依赖或模板渲染时运行：

```bash
GOWORK=off make boundary
GOWORK=off make contracts
GOWORK=off make integration
```
