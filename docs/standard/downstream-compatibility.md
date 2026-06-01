# Downstream Compatibility

`xlib-standard` 的下游兼容性必须通过真实生成库 smoke 证明，而不是只证明模板仓库自身可用。

## 默认下游

- `kernel`：默认 L0 集成目标，module path 为 `github.com/ZoneCNH/kernel`。
- `corekit`：中性组织路径 smoke，避免只对 ZoneCNH 路径成立。

旧 `foundationx` 只作为迁移兼容名出现，不再是默认下游。

## 目标库矩阵

详细 module path、package、layer、allowed deps 和 forbidden deps 见 [`../downstream-matrix.md`](../downstream-matrix.md)。矩阵至少覆盖：`kernel`、`configx`、`observex`、`testkitx`、`postgresx`、`redisx`、`kafkax`、`taosx`、`ossx`、`clickhousex`。

## Gate

`GOWORK=off make integration` 是默认下游兼容 gate。它应覆盖 generator smoke、kernel/corekit 代表路径和关键边界检查。

生成出的每个代表下游必须通过：

- `GOWORK=off go test ./...`
- `GOWORK=off make contracts`
- `GOWORK=off make boundary`
- `CHECK_STATUS=passed GOWORK=off make evidence`
- `RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check`

失败时不得宣称 downstream compatible。
