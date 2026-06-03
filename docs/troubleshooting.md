# 故障排查

## `GOWORK=off is required`

治理和发布目标要求隔离 Go workspace。重新执行命令时在前面加上：

```sh
GOWORK=off make <target>
```

## Makefile 出现 duplicate target warning

应只有一个同名目标定义。若看到 `overriding recipe` 或 `ignoring old recipe`，先检查是否新增了重复 `.PHONY` 或重复 target，再运行：

```sh
GOWORK=off make --warn-undefined-variables governance-check
```

## 缺少 `golangci-lint` 或启用漏洞扫描时缺少 `govulncheck`

本地 `make lint` 依赖 `golangci-lint`。`make security` 默认只运行 secret scan；只有设置 `XLIB_ENABLE_VULNCHECK=1` 时才依赖 `govulncheck`。CI 会按同一开关安装固定版本；本地可按 CI workflow 中的版本安装，或记录为未运行的本地工具缺口。

## release manifest 缺失

`release/manifest/latest.json` 和 `latest.json.sha256` 是生成产物，不应提交。运行 release/evidence 相关命令重新生成，并通过 artifact 或校验和复核。

## 下游 kernel/configx 未通过

不要用本仓库的占位文件替代下游证据。需要在真实下游仓库运行采纳/兼容命令，并把输出记录到正式证据链。
