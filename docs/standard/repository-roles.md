# 仓库角色

本文件区分标准仓库、模板仓库、生成基础库、适配器库和业务应用的职责。角色边界用于防止基础库向上依赖业务层，也用于判断文档、CI 和 release Evidence 是否完整。

| 仓库 | 角色 | 必须包含 | 禁止包含 |
| --- | --- | --- | --- |
| `xlib-standard` | 独立标准仓库（`https://github.com/ZoneCNH/xlib-standard`） | P0 标准文档、仓库角色、分层、模块边界、contracts 规范源、release Evidence 规则 | `baselib-template` 特有实现代码、generator 运行时代码、真实 runtime 实现、业务模型、生产凭据 |
| `baselib-template` | 模板、generator、Harness、Evidence 实现仓库 | 标准对齐文档副本、模板包、generator、contracts 实现副本、CI、release manifest、`.agent/` 模板、Evidence 生成器 | 真实 L1 runtime 实现、业务模型、`x.go` 依赖、生产凭据、替代 `xlib-standard` 的标准正文 |
| `foundationx` | L1 通用基础能力库 | context、error、config、logging、metrics、lifecycle 等通用 runtime primitives | 存储、消息、业务语义、应用框架耦合 |
| `postgresx` | L1/L2 数据库适配器库 | PostgreSQL profile、连接配置、健康检查、错误分类、测试夹具 | 业务 repository、应用 transaction 编排 |
| `redisx` | L1/L2 cache 适配器库 | Redis profile、连接、health、metrics、testkit | 业务缓存 key 语义、应用级策略 |
| `kafkax` | L1/L2 messaging 适配器库 | producer/consumer profile、contract、health、metrics | 业务 topic 设计、业务消息 schema |
| `taosx` | L1/L2 时序数据适配器库 | TDengine profile、连接、health、contract | 业务指标模型 |
| `ossx` | L1/L2 object storage 适配器库 | bucket/profile、上传下载 contract、health | 业务文件生命周期策略 |
| `clickhousex` | L1/L2 分析数据库适配器库 | ClickHouse profile、连接、query contract、health | 产品报表语义 |
| `x.go` | 应用或框架组合层 | 组合基础库、业务 wiring、应用生命周期 | 作为基础库模板的依赖前提 |

## 判定规则

- `xlib-standard` 是标准权威源；`baselib-template` 中的标准文档用于实现对齐、生成和 gate 校验，不得覆盖或替代 `xlib-standard`。
- 模板仓库只承担可复制结构、generator、Harness 与 Evidence 生成实现，不实现具体基础设施能力。
- L1 基础库可以复用 `foundationx`，但不得反向依赖 `x.go`。
- L2/profile 库可以依赖更低层基础库，但不得向业务层取配置或模型。
- 业务层可以组合所有基础库，但不得把业务规则下沉到基础库。
