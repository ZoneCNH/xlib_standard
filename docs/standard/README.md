# 基础库标准文档索引

本目录是 `baselib-template` 对独立标准仓库 [`xlib-standard`](https://github.com/ZoneCNH/xlib-standard) 的 P0 落地入口。生成或维护任何 `x.go` 基础库时，先以这里的规则判断仓库职责、模块边界、Harness gate、Evidence 和 release 是否完整。

## 必读顺序

1. [基础库总标准](xlib-standard.md)：同步 [`xlib-standard`](https://github.com/ZoneCNH/xlib-standard) 定义的基础库公共语义和禁止项。
2. [仓库角色](repository-roles.md)：区分模板、L1 基础库、适配器库和 `x.go`。
3. [分层规则](layering.md)：规定 Standard、L0、L1、L2 和业务层的依赖方向。
4. [模块边界](module-boundary.md)：明确 `baselib-template` 和生成库的允许/禁止内容。
5. [完成定义](dod.md)：给出 Task、Issue、Goal 和 Release 级别的 DoD。
6. [Harness gate](harness-gates.md)：列出 required、extended、generator 和 final gate。
7. [Evidence 协议](evidence-protocol.md)：规范 `DONE with evidence:` 和 release manifest。
8. [发布标准](release-standard.md)：约束 release-check、final-check 和 preflight。
9. [安全和密钥策略](security-and-secret-policy.md)：约束 Secret Gate、govulncheck 和凭据边界。
10. [模板生成契约](template-generation-contract.md)：约束 `scripts/render_template.sh`。
11. [下游兼容](downstream-compatibility.md)：定义 foundationx/corekit 和未来 L1 profile 的验证。
12. [复盘和补丁](retrospective-and-patches.md)：把失败转化为 Prompt、Harness、Rule 和 CI Gate 补丁。

## 适用范围

- [`xlib-standard`](https://github.com/ZoneCNH/xlib-standard)：独立标准源仓库，承载跨基础库共享的标准文本。
- `baselib-template`：模板、generator、Harness 和 Evidence 的实现仓库，负责把标准源落到可编译模板和验证链路。
- 生成基础库：继承本目录的标准，并在自身仓库中记录 profile-specific 差异。
- `x.go`：只能消费基础库，不得成为基础库模板的依赖前提。

## 完成要求

任何声称完成的变更都必须同时满足：

- 关联的标准文档已更新或明确不适用。
- Required gate 有新鲜输出。
- release Evidence 可生成并可校验。
- 禁止项未被引入。
- 最终说明使用 `DONE with evidence:` 格式。

## ADR 2 落地追踪

[ADR 2](../adr/2.md) 要求模板工厂、标准文档和 Evidence gate 可以被同一套命令验证。本目录用以下入口约束落地状态：

- `GOWORK=off make docs-check`：检查标准文档骨架、README 入口、本地链接、未替换模板占位符、关键文本，以及 ADR 2 要求的 Evidence 协议锚点。它是结构和漂移检查，不做语义审查；文档是否准确表达标准仍由 reviewer 判断。
- `GOWORK=off make release-check`：在 `ci`、`integration` 和 `docs-check` 通过后生成并校验 release Evidence。发布式验证必须显式使用 `GOWORK=off`，避免本地或父级 `go.work` 干扰独立 module 解析。
- `release/manifest/template.json`：提交到源码历史的 manifest 结构契约。
- `release/manifest/latest.json`：由 release/Evidence gate 生成、被 `.gitignore` 排除，并在 CI 中作为 artifact 引用；`release/manifest/latest.json.sha256` 同样是生成产物，必须保持忽略，并由 `release-check` 生成和校验。
- `scripts/run_fuzz_smoke.sh`：默认是快速 smoke，`FUZZ_SMOKE_TIME` 未设置时每个 fuzz target 使用 `10s`；深度 fuzz 需要显式设置更长 `FUZZ_SMOKE_TIME`，并在最终 Evidence/DONE 说明中记录该时间配置。
