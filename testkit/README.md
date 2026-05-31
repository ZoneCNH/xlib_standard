# Testkit 测试工具

为生成的基础库提供可复用测试夹具和断言。

## 契约

- `Config(name string)` 返回带 `Name` 和 `Timeout` 的最小有效配置。
- `RequireNoError(t, err)` 在 `err == nil` 时保持静默，在非空错误时终止当前测试。

## 回归覆盖

`fixture_test.go` 锁定 `Config("fixture")` 的字段和 `Validate` 结果，并验证 `RequireNoError(t, nil)` 可用。生成后的基础库需要保留这组最小测试，以防测试夹具随包名替换或配置 contract 漂移。

生成的库应保持此包独立于 `x.go` 和业务特定模型。
