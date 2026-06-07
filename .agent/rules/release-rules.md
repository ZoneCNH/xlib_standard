# Release 规则

> 源自 Goal 完整规则 v1.0 §16

## RULE-RELEASE-001：Release 必须有 Release Manifest

Release Manifest 必须包含：

```text
version
date
goal
commit
tag
included issues
included PRs
changes
evidence summary
test summary
compatibility
migration notes
risks
rollback plan
known issues
retrospective
```

## RULE-RELEASE-002：Release 前必须全部验证

必须通过：

```text
make release-check
make evidence-check
make ci
```

## RULE-RELEASE-003：Release 必须可回滚

必须有：

```text
rollback command
rollback condition
last known good version
affected components
risk note
```

## RULE-RELEASE-004：Main 合并必须递增 patch 版本

每一次成功合并到 `main` 必须对应且仅对应一个新的稳定 semver patch release。
发布版本必须从当前远端最高稳定 `vX.Y.Z` tag 计算为 `vX.Y.(Z+1)`，
并由 `.github/workflows/release-auto-patch.yml` 在同一次 `main` push workflow 内完成：

```text
release-final-check
git tag -a
git push origin "refs/tags/${RELEASE_TAG}"
GitHub Release 发布
gh release view 校验
```

同一 `GITHUB_SHA` 的 workflow rerun 若已存在稳定 release tag，
必须设置 `already_released=true` 并复用该 tag，不得再次递增版本。
`release-auto-patch-main` 并发组必须保持串行，
防止多个 `main` push 抢占同一个 patch 版本。

## Release 验收清单

```text
[ ] Issues 已关闭
[ ] PRs 已合并
[ ] Evidence 已归档
[ ] CHANGELOG 已更新
[ ] VERSION 已按 main 合并 patch +1 策略更新
[ ] Release Manifest 已生成
[ ] Rollback Plan 已记录
```
