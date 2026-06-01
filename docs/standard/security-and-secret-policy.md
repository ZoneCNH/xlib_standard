# 安全和密钥策略

基础库模板必须默认安全。安全策略覆盖源码、文档、测试、CI、manifest、Issue 和 PR。

## 密钥边界

禁止提交：

- 真实 token、password、private key、API key、access key。
- 真实生产连接串。
- 可直接访问生产资源的 endpoint 与凭据组合。
- 从本地环境复制的 `.env` 或 kubeconfig。

示例必须使用占位符，例如 `example-token`、`example-secret` 或 `localhost`。

## Secret Gate

`GOWORK=off make security` 必须执行：

- `govulncheck ./...`
- `scripts/check_secrets.sh`

缺少 `govulncheck` 必须失败，不能跳过。secret scan 发现疑似凭据时必须阻断。

## 日志和 Evidence

- 日志不得输出敏感字段原值。
- release manifest 不得记录 secret。
- PR 描述和 Issue 模板不得要求粘贴真实凭据。

## 依赖安全

- 新依赖必须有明确用途。
- 依赖变更后运行 `GOWORK=off make security` 和 `GOWORK=off make boundary`。
- 发现漏洞时记录影响面、修复版本和验证命令。

## 例外

安全例外必须有 ADR 或 Issue 记录，且只能放宽非 secret、非生产、非凭据相关约束。生产凭据和真实密钥没有例外路径。
