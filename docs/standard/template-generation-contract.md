# 模板生成契约

`scripts/render_template.sh` 是从 `baselib-template` 生成具体基础库的唯一标准入口。

## 输入

必须显式传入：

- `--module-name`
- `--module-path`
- `--package-name`
- `--out`

示例：

```bash
scripts/render_template.sh \
  --module-name kernel \
  --module-path github.com/ZoneCNH/kernel \
  --package-name kernel \
  --out ../kernel
```

中性下游路径必须纳入固定验证集：

```bash
scripts/render_template.sh \
  --module-name corekit \
  --module-path example.com/acme/corekit \
  --package-name corekit \
  --out ../corekit
```

## 输出目录安全

generator 必须保护调用方不覆盖源码或已有仓库：

- `--out` 不得指向 `baselib-template` 仓库根目录。
- `--out` 不得位于 `baselib-template` 仓库内部。
- `--out` 不得包含已有 `.git/` 或 `go.mod`。
- `--out` 必须为空目录；不存在时可以创建。
- 失败时必须在复制前退出，避免留下半渲染仓库。

## 替换规则

generator 必须替换：

- module name token。
- module path token。
- package name token。
- `github.com/ZoneCNH/baselib-template`
- `pkg/templatex`
- `package templatex`
- `templatex` imports、文档引用、测试 fixture 和脚本参数。
- `Templatex` title-case 引用。
- `TEMPLATEX` upper-case 引用。

### Metrics Prefix

Metrics Prefix 必须跟随 package name 替换。模板中的 `templatex_` prefix 在 `kernel` 渲染后必须变为 `kernel_`，在 `example.com/acme/corekit` 渲染后必须变为 `corekit_`。metrics contract、README、docs、examples、测试和 snapshot 中不得残留 `templatex_`，除非某个文件被明确 allowlist 为模板来源说明。

## 排除规则

generator 不得复制：

- `.git/`
- `.omx/`
- `.worktree/`
- `release/manifest/latest.json`
- `release/manifest/latest.json.sha256`
- `docs/adr/`
- `docs/goal.md`
- 临时文件、缓存、coverage 输出、构建目录、本地 Evidence 输出和 editor 产物。

## 输出不变量

生成结果必须满足：

- `go.mod` module path 正确。
- 公共包目录和 package name 正确。
- README、docs、contracts、Makefile、scripts、CI 和 `.agent/` 模板存在。
- 无 template token 未替换残留。
- 无 generic placeholder、TODO-style template marker 或 `templatex_` metrics prefix 残留。
- 无 `baselib-template` module import 残留，除非在文档中作为来源说明出现。
- `GOWORK=off go mod tidy` 后 `go.mod` 和 `go.sum` 保持 clean。

## Scanner

`scripts/check_rendered_template.sh` 必须扫描：

- stale module path：`github.com/ZoneCNH/baselib-template`。
- stale module name：`baselib-template`。
- stale package directory：`pkg/templatex`。
- stale package name：`templatex`、`Templatex`、`TEMPLATEX`。
- stale metrics prefix：`templatex_`。
- unresolved template token、generic placeholder 和 TODO-style template marker。

扫描 unresolved template token 时，应跳过合法包含表达式语法的 GitHub Actions workflow、检查脚本自身和 Go template probe 脚本，避免把 scanner 规则或 `go list -f` 模板语法误判为未替换占位符。

扫描失败时 integration gate 必须失败。

## 验证

```bash
GOWORK=off make integration
GOWORK=off make boundary
GOWORK=off make contracts
GOWORK=off make release-check
```

任何 generator 修改必须附带 integration Evidence。Release 级验证还必须证明渲染出的 `kernel` 和 `corekit` 可以独立运行 `go mod tidy`、`go test ./...`、`make contracts`、`make boundary`、`make evidence` 和 `make release-evidence-check`，且所有命令都在 `GOWORK=off` 下执行。
