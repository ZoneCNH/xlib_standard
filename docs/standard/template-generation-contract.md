# Template Generation Contract

`scripts/render_template.sh` 是从 `xlib-standard` 生成具体基础库的唯一标准入口。旧 `baselib-template` module path 只作为迁移扫描项和兼容说明保留，不能作为新主身份。

## 默认示例

```bash
scripts/render_template.sh \
  --module-name kernel \
  --module-path github.com/ZoneCNH/kernel \
  --package-name kernel \
  --out ../kernel
```

## 必须替换

- Module name、module path 和 package name 占位符必须替换为目标库值。
- `github.com/ZoneCNH/xlib-standard` 模板自身 import 到目标 module path。
- 旧迁移扫描项：`github.com/ZoneCNH/baselib-template`、`baselib-template`、`foundationx`。
- `pkg/templatex` 目录名到 `pkg/<package-name>`。
- README、docs、contracts、examples、scripts 和 manifest 中的模板占位。

## 不变量

- `--out` 不得指向 `xlib-standard` 仓库根目录，也不得位于仓库内部。
- 生成库不得依赖 `x.go` 或业务仓库。
- 生成库不得读取 `/home/k8s/secrets/env/*`；该路径只属于调用方部署配置。
- 生成后的 module 必须在 `GOWORK=off` 下运行测试、contracts、boundary 和 release Evidence gate。
- 旧名只可在迁移文档或兼容说明中出现，不得作为生成库主标题、module name、package name 或 release 主体。

## Metrics Prefix

Metrics Prefix 必须跟随 package name 替换。模板中的 `templatex_` prefix 在 `kernel` 渲染后必须变为 `kernel_`，在 `example.com/acme/corekit` 渲染后必须变为 `corekit_`。metrics contract、README、docs、examples、测试和 snapshot 中不得残留 `templatex_`，除非某个文件被明确 allowlist 为模板来源说明。

## Release 验证

任何 generator 修改必须附带 integration Evidence。Release 级验证必须证明渲染出的 `kernel` 和中性路径 `corekit` 可以独立运行：

- `GOWORK=off go test ./...`
- `GOWORK=off make contracts`
- `GOWORK=off make boundary`
- `CHECK_STATUS=passed GOWORK=off make evidence`
- `RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check`
