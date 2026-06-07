# baselib-template 到 xlib-standard 迁移指南

## 目标

把旧 `baselib-template` 身份迁移为 `xlib-standard`，并把旧默认下游 `foundationx` 迁移为 `kernel`。迁移完成后，旧名只允许出现在迁移 ADR、迁移指南、历史变更记录和兼容性说明中。

## 名称规则

| 旧名 | 新名 | 允许保留位置 |
| --- | --- | --- |
| `baselib-template` | `xlib-standard` | ADR、迁移指南、CHANGELOG、兼容性说明 |
| `foundationx` | `kernel` | ADR、迁移指南、历史 compatibility note |

## 迁移要求

- README 主标题和主叙事必须使用 `xlib-standard`。
- 生成示例默认使用 `kernel` / `github.com/ZoneCNH/kernel`。
- 标准文档必须声明 `xlib-standard` 同时承担 Standard Source、Go Reference Template、Generator、Harness 和 Evidence Runtime。
- Go module path、包名、render script 和 CI 迁移由实现 worker 完成；本文档只定义语义约束。

## 稳定版迁移 Gate 覆盖

Full release 前，legacy-reference scan 不能只覆盖文档和 `.agent`/`.xlib` YAML/JSON。稳定版检查必须覆盖所有可能携带旧身份、下游默认名或旧配置根的手写/模板/执行入口：

- 源码与脚本：`cmd/`、`internal/`、`pkg/`、`scripts/`、`testkit/`、`*.go`、`*.sh`、`Makefile`、`Dockerfile`、`docker-compose.yml`。
- 治理与运行入口：`.agent/`、`.xlib/`、`.github/`、`.githooks/`、`.devcontainer/`、`.gitignore`、`.dockerignore`。
- 下游与模板表面：`templates/`、`contracts/`、`examples/`、`docs/standard/`、`docs/migration/`。
- 运行时上下文证据：已纳入版本控制的 `.agent/context/` 示例与 schema；本地 `.omx/`、`.worktree/` 只作为运行时输入，不作为 release source-of-truth。

推荐扫描入口：

```bash
rg -n "baselib-template|foundationx|\.agent|\.xlib|\.config" \
  README.md docs .agent .xlib .github .githooks .devcontainer templates contracts examples \
  cmd internal pkg scripts testkit Makefile Dockerfile docker-compose.yml .gitignore .dockerignore \
  --glob '*.go' --glob '*.sh' --glob '*.md' --glob '*.yaml' --glob '*.yml' --glob '*.json' --glob '*.tmpl'
```

命中分类规则：

- `baselib-template` / `foundationx` 只能留在迁移 ADR、本文档、历史 changelog、兼容性说明或明确的下游 adoption proof 中。
- `.agent` / `.xlib` 允许作为当前运行时和事实目录引用，但迁移到稳定配置根前必须在 runtime ownership、physical migration manifest 或 downstream/template 文档中被分类。
- `.config` 命中必须表示目标配置根或迁移阻断项；在 `.config/` 统一完成前，不能声明 `v1.0.0` 或 release-ready。
- `.github`、`.githooks`、`.devcontainer`、ignore 文件和 generated artifacts 不是配置 source-of-truth；它们是平台入口、保护性忽略规则或 gate 生成输出，必须由 inventory/ownership gate 覆盖。

## 验证

- 上述 `rg` legacy/config scan 的命中必须能按“命中分类规则”归类。
- `GOWORK=off make docs-check` 必须通过。
- `GOWORK=off make p2-runtime-check` 必须通过，覆盖 runtime ownership、downstream baseline 和 downstream adoption。
- Full release 前必须补充 downstream integration Evidence，且 `.config/` 统一完成前不得声称 release-ready。
