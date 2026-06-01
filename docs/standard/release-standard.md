# 发布标准

发布流程必须证明源码、contracts、依赖和 gate 状态一致。`xlib-standard` 的 release 标准同时约束生成基础库；旧 `baselib-template` 仅作为迁移兼容名记录。

## 发布路径

1. 运行 required gate。
2. 运行 integration 和 generator smoke。
3. 生成 Evidence manifest。
4. 校验 Evidence manifest。
5. 在 clean workspace 运行 final check。
6. 使用明确版本运行 preflight。
7. 在 PR 或 release notes 中附上 Evidence 摘要。

## 命令

```bash
GOWORK=off make release-check
GOWORK=off make release-check-extended
GOWORK=off make release-final-check
GOWORK=off make release-preflight VERSION=v1.0.0
```

## Manifest

`release/manifest/latest.json` 是生成产物：

- 可以作为 CI artifact 上传。
- 可以作为本地 Evidence 检查输入。
- 不提交到源码历史。
- `release/manifest/latest.json.sha256` 是对应 checksum 产物，随 CI artifact 上传，并保持在 `.gitignore` 中。

## 版本

- `VERSION` 必须显式传入 release-preflight。
- 版本应与 release notes、tag 和 manifest 一致。
- 未创建 tag 或工作区 dirty 时，不得宣称最终发布完成。

## 变更说明

PR 或 release notes 必须说明：

- 对模板行为的影响。
- 对生成库的影响。
- 已运行命令。
- Evidence artifact。
- known gaps 或 blocked gate。
