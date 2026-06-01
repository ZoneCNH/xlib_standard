# Harness Gates

Harness Gate 把 `xlib-standard` 的标准、模板、generator、Evidence 和 release 要求变成可执行检查。

## Required Gates

- `GOWORK=off make fmt`
- `GOWORK=off make vet`
- `GOWORK=off make lint`
- `GOWORK=off make test`
- `GOWORK=off make race`
- `GOWORK=off make boundary`
- `GOWORK=off make security`
- `GOWORK=off make contracts`
- `GOWORK=off make docs-check`
- `GOWORK=off make integration`
- `CHECK_STATUS=passed GOWORK=off make evidence`
- `RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check`

## Generator Gate

Generator gate 必须证明模板能生成代表性下游，而不是只证明 `xlib-standard` 自身可用。

代表下游：

- `kernel`
- `corekit`

旧 `foundationx` 只作为迁移兼容扫描项，不再作为默认下游。

## Final Gates

- `GOWORK=off make release-final-check`
- `GOWORK=off make release-preflight VERSION=<version>`
- `xlibgate score --min 9.8`

## Secret Gate

Secret Gate 必须确认源码、README、测试日志、release manifest、PR 描述和 Evidence 不包含 `/home/k8s/secrets/env/*` 的真实内容。该路径只能在文档中作为调用方部署路径名出现。
